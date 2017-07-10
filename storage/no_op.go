package storage

import (
	"github.com/anacrolix/torrent/storage"
	"github.com/anacrolix/torrent/metainfo"
)

type noOpClient struct{}

type noOpTorrent struct{}

type noOpPiece struct {
}

func NewNoOpStorage() storage.ClientImpl {
	ret := &noOpClient{}
	return ret
}

func (self *noOpClient) Close() error {
	return nil
}

func (self *noOpClient) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
	return &noOpTorrent{}, nil
}

func (self *noOpTorrent) Piece(p metainfo.Piece) storage.PieceImpl {
	ret := &noOpPiece{}
	return ret
}

func (self *noOpTorrent) Close() error {
	return nil
}

func (self *noOpPiece) GetIsComplete() (complete bool) {
	complete = false
	return
}

func (self *noOpPiece) MarkComplete() error {
	return nil
}

func (self *noOpPiece) MarkNotComplete() error {
	return nil
}

func (self *noOpPiece) ReadAt(b []byte, off int64) (n int, err error) {
	return len(b), nil
}

func (me *noOpPiece) WriteAt(b []byte, off int64) (n int, err error) {
	return len(b), nil
}
