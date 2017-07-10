// +build !cgo

package torrent_nicolov

import (
	"github.com/anacrolix/utp"
)

func NewUtpSocket(network, addr string) (utpSocket, error) {
	return utp.NewSocket(network, addr)
}
