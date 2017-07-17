package main

/*
Takes in a torrent and its (completed) files, and saves each piece in a
separate file. The data is saved under OutDir/{torrent_infohash}.
 */

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/anacrolix/tagflag"
	"github.com/bradfitz/iter"
	"github.com/edsrzf/mmap-go"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/mmap_span"
	"path"
)

func mmapFile(name string) (mm mmap.MMap, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return
	}
	if fi.Size() == 0 {
		return
	}
	return mmap.MapRegion(f, -1, mmap.RDONLY, mmap.COPY, 0)
}

func saveTorrentPieces(info *metainfo.Info, inputDir string, outputDir string) error {
	span := new(mmap_span.MMapSpan)
	for _, file := range info.UpvertedFiles() {
		filename := filepath.Join(append([]string{inputDir, info.Name}, file.Path...)...)
		mm, err := mmapFile(filename)
		if err != nil {
			return err
		}
		if int64(len(mm)) != file.Length {
			return fmt.Errorf("file %q has wrong length", filename)
		}
		span.Append(mm)
	}

	os.MkdirAll(outputDir, os.ModePerm)

	for i := range iter.N(info.NumPieces()) {
		p := info.Piece(i)
		hash := sha1.New()
		_, err := io.Copy(hash, io.NewSectionReader(span, p.Offset(), p.Length()))
		if err != nil {
			return err
		}
		good := bytes.Equal(hash.Sum(nil), p.Hash().Bytes())
		if !good {
			log.Fatalf("hash mismatch at piece %d", i)
		}

		pieceOutputPath := path.Join(outputDir, p.Hash().String())

		fmt.Println(pieceOutputPath)

		reader := io.NewSectionReader(span, p.Offset(), p.Length())

		f, err := os.Create(pieceOutputPath)
		if err != nil {
			log.Panic(err)
		}

		io.Copy(f, reader)
		f.Close()

		fmt.Printf("%d: %x: %v\n", i, p.Hash(), good)
	}
	return nil
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	var flags = struct {
		DataDir       string
		BaseOutputDir string
		tagflag.StartPos
		TorrentFile   string
	}{}
	tagflag.Parse(&flags)

	if flags.DataDir == "" {
		log.Fatal("Missing dataDir")
	}

	if flags.BaseOutputDir == "" {
		log.Fatal("Missing outputDir")
	}

	metaInfo, err := metainfo.LoadFromFile(flags.TorrentFile)
	if err != nil {
		log.Fatal(err)
	}
	info, err := metaInfo.UnmarshalInfo()
	if err != nil {
		log.Fatalf("error unmarshalling info: %s", err)
	}

	outputDir, err := filepath.Abs(
		path.Join(flags.BaseOutputDir, metaInfo.HashInfoBytes().String()))

	log.Print("Saving pieces to: ", outputDir)

	err = saveTorrentPieces(&info, flags.DataDir, outputDir)
	if err != nil {
		log.Fatalf("torrent failed verification: %s", err)
	}
}
