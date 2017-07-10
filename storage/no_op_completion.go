package storage

import (
	"github.com/anacrolix/torrent/storage"
	"github.com/anacrolix/torrent/metainfo"
)

type noOpPieceCompletion struct{}

func NewNoOpPieceCompletion() (ret storage.PieceCompletion, err error) {
	ret = &noOpPieceCompletion{}
	return
}

func (self *noOpPieceCompletion) Get(pk metainfo.PieceKey) (ret bool, err error) {
	ret = false
	return
}

func (self *noOpPieceCompletion) Set(pk metainfo.PieceKey, b bool) error {
	return nil
}

func (self *noOpPieceCompletion) Close() error {
	return nil
}
