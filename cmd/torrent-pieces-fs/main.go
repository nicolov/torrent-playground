package main

/*
Exposes a directory
 */

import (
	//"flag"
	"log"
	//"net"
	//"net/http"
	_ "net/http/pprof"
	"os"
	//"os/signal"
	//"os/user"
	//"path/filepath"
	//"syscall"
	//"time"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/tagflag"

	//"github.com/anacrolix/torrent"
	//"github.com/anacrolix/torrent/fs"
	//"github.com/anacrolix/torrent/util/dirwatch"

	"torrent-object-storage/pieces_fs"
	"os/signal"
	"syscall"
)

func exitSignalHandlers(fs *pieces_fs.PiecesFS, mountDir string) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for {
		<-c
		//fs.Destroy()
		err := fuse.Unmount(mountDir)
		if err != nil {
			log.Print(err)
		}
	}
}

func main() {
	os.Exit(mainExitCode())
}

func mainExitCode() int {
	log.SetFlags(log.Flags() | log.Lshortfile)
	var flags = struct {
		TorrentsDir string // Directory with .torrent files that make up the FS
		DataDir     string // Directory with pieces data
		MountDir    string // Location where the FS will be mounted
	}{}
	tagflag.Parse(&flags)

	conn, err := fuse.Mount(flags.MountDir)
	if err != nil {
		log.Fatal(err)
	}
	defer fuse.Unmount(flags.MountDir)
	// TODO: Think about the ramifications of exiting not due to a signal.
	defer conn.Close()

	//client, err := torrent.NewClient(&torrent.Config{
	//	DataDir:         *downloadDir,
	//	DisableTrackers: *disableTrackers,
	//	ListenAddr:      *listenAddr,
	//	NoUpload:        true, // Ensure that downloads are responsive.
	//})
	//if err != nil {
	//	log.Print(err)
	//	return 1
	//}
	//// This is naturally exported via GOPPROF=http.
	//http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//	client.WriteStatus(w)
	//})
	//dw, err := dirwatch.New(*torrentPath)
	//if err != nil {
	//	log.Printf("error watching torrent dir: %s", err)
	//	return 1
	//}
	//go func() {
	//	for ev := range dw.Events {
	//		switch ev.Change {
	//		case dirwatch.Added:
	//			if ev.TorrentFilePath != "" {
	//				_, err := client.AddTorrentFromFile(ev.TorrentFilePath)
	//				if err != nil {
	//					log.Printf("error adding torrent to client: %s", err)
	//				}
	//			} else if ev.MagnetURI != "" {
	//				_, err := client.AddMagnet(ev.MagnetURI)
	//				if err != nil {
	//					log.Printf("error adding magnet: %s", err)
	//				}
	//			}
	//		case dirwatch.Removed:
	//			T, ok := client.Torrent(ev.InfoHash)
	//			if !ok {
	//				break
	//			}
	//			T.Drop()
	//		}
	//	}
	//}()
	//resolveTestPeerAddr()

	fs := pieces_fs.New(&pieces_fs.PiecesFSConfig{
		DataDir:     flags.DataDir,
		TorrentsDir: flags.TorrentsDir,
	})

	//fs := torrentfs.New(client)
	go exitSignalHandlers(fs, flags.MountDir)
	//
	//
	if err := fusefs.Serve(conn, fs); err != nil {
		log.Fatal(err)
	}
	<-conn.Ready
	if err := conn.MountError; err != nil {
		log.Fatal(err)
	}
	return 0
}
