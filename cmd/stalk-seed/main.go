// Downloads torrents from the command-line.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/tagflag"
	"github.com/gosuri/uiprogress"
	"golang.org/x/time/rate"

	"github.com/anacrolix/torrent/metainfo"
	"time"

	"torrent-object-storage/storage"
	"torrent-object-storage"
	"github.com/anacrolix/torrent"
)

func addTorrents(client *torrent_nicolov.Client) {
	for _, arg := range flags.Torrent {
		t := func() *torrent_nicolov.Torrent {
			if strings.HasPrefix(arg, "magnet:") {
				t, err := client.AddMagnet(arg)
				if err != nil {
					log.Fatalf("error adding magnet: %s", err)
				}
				return t
			} else if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
				response, err := http.Get(arg)
				if err != nil {
					log.Fatalf("Error downloading torrent file: %s", err)
				}

				metaInfo, err := metainfo.Load(response.Body)
				defer response.Body.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error loading torrent file %q: %s\n", arg, err)
					os.Exit(1)
				}
				t, err := client.AddTorrent(metaInfo)
				if err != nil {
					log.Fatal(err)
				}
				return t
			} else if strings.HasPrefix(arg, "infohash:") {
				t, _ := client.AddTorrentInfoHash(metainfo.NewHashFromHex(strings.TrimPrefix(arg, "infohash:")))
				return t
			} else {
				metaInfo, err := metainfo.LoadFromFile(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error loading torrent file %q: %s\n", arg, err)
					os.Exit(1)
				}
				t, err := client.AddTorrent(metaInfo)
				if err != nil {
					log.Fatal(err)
				}
				return t
			}
		}()
		t.AddPeers(func() (ret []torrent_nicolov.Peer) {
			for _, ta := range flags.TestPeer {
				ret = append(ret, torrent_nicolov.Peer{
					IP:   ta.IP,
					Port: ta.Port,
				})
			}
			return
		}())
		go func() {
			<-t.GotInfo()
			t.DownloadAll()

			for {
				time.Sleep(500 * time.Millisecond)
				client.WriteSwarmHealth(os.Stdout)
			}
		}()
	}
}

var flags = struct {
	TestPeer     []*net.TCPAddr `help:"addresses of some starting peers"`
	Seed         bool           `help:"seed after download is complete"`
	Addr         *net.TCPAddr   `help:"network listen addr"`
	UploadRate   tagflag.Bytes  `help:"max piece bytes to send per second"`
	DownloadRate tagflag.Bytes  `help:"max bytes per second down from peers"`
	tagflag.StartPos
	Torrent      []string `arity:"+" help:"torrent file path or magnet uri"`
}{
	UploadRate:   -1,
	DownloadRate: -1,
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tagflag.Parse(&flags)
	var clientConfig torrent.Config

	clientConfig.DefaultStorage = storage.NewNoOpStorage()

	if flags.Addr != nil {
		clientConfig.ListenAddr = flags.Addr.String()
	}
	if flags.Seed {
		clientConfig.Seed = true
	}
	if flags.UploadRate != -1 {
		clientConfig.UploadRateLimiter = rate.NewLimiter(rate.Limit(flags.UploadRate), 256<<10)
	}
	if flags.DownloadRate != -1 {
		clientConfig.DownloadRateLimiter = rate.NewLimiter(rate.Limit(flags.DownloadRate), 1<<20)
	}

	client, err := torrent_nicolov.NewClient(&clientConfig)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	defer client.Close()
	// Write status on the root path on the default HTTP muxer. This will be
	// bound to localhost somewhere if GOPPROF is set, thanks to the envpprof
	// import.
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		client.WriteStatus(w)
	})
	uiprogress.Start()
	addTorrents(client)
	if client.WaitAll() {
		log.Print("downloaded ALL the torrents")
	} else {
		log.Fatal("y u no complete torrents?!")
	}
	if flags.Seed {
		select {}
	}
}
