package pieces_fs

import (
	//"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"bazil.org/fuse"
	"os"
	"golang.org/x/net/context"
)

type PiecesFSConfig struct {
	// Config needed to set up the file system (location of torrent files, etc..)
	TorrentsDir string
	DataDir     string
}

type PiecesFS struct {
	config *PiecesFSConfig
}

type rootNode struct {
	fs *PiecesFS
}

func (tfs PiecesFS) Root() (fusefs.Node, error) {
	return rootNode{&tfs}, nil
}

func (rn rootNode) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir
	return nil
}

//func (rn rootNode) Lookup(ctx context.Context, name string) (_node fusefs.Node, err error) {
//	//for _, t := range rn.fs.Client.Torrents() {
//	//	info := t.Info()
//	//	if t.Name() != name || info == nil {
//	//		continue
//	//	}
//	//	__node := node{
//	//		metadata: info,
//	//		FS:       rn.fs,
//	//		t:        t,
//	//	}
//	//	if !info.IsDir() {
//	//		_node = fileNode{__node, uint64(info.Length), 0}
//	//	} else {
//	//		_node = dirNode{__node}
//	//	}
//	//	break
//	//}
//	//if _node == nil {
//	//	err = fuse.ENOENT
//	//}
//	return
//}
//
//func (rn rootNode) ReadDirAll(ctx context.Context) (dirents []fuse.Dirent, err error) {
//	//for _, t := range rn.fs.Client.Torrents() {
//	//	info := t.Info()
//	//	if info == nil {
//	//		continue
//	//	}
//	//	dirents = append(dirents, fuse.Dirent{
//	//		Name: info.Name,
//	//		Type: func() fuse.DirentType {
//	//			if !info.IsDir() {
//	//				return fuse.DT_File
//	//			} else {
//	//				return fuse.DT_Dir
//	//			}
//	//		}(),
//	//	})
//	//}
//	return
//}

//func New(cl *torrent.Client) *TorrentFS {
//	fs := &TorrentFS{
//		Client:    cl,
//		destroyed: make(chan struct{}),
//	}
//	fs.event.L = &fs.mu
//	return fs
//}

func New(config *PiecesFSConfig) * PiecesFS {
	fs := &PiecesFS{
		config: config,
	}
	return fs
}