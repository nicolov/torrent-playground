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
	"github.com/anacrolix/torrent/metainfo"
	"time"

	"torrent-object-storage/storage"
	"torrent-object-storage"
	"github.com/anacrolix/torrent"
	"encoding/json"
	"io/ioutil"
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
			// Wait for torrent info
			<-t.GotInfo()
			// Start downloading
			//t.DownloadPieces(0, 1)
			t.DownloadAll()
		}()
	}

	go func() {
		// Wait and print/save swarm stats
		waitTimeSec := 10
		log.Printf("Waiting %d seconds before saving stats", waitTimeSec)
		time.Sleep(time.Duration(waitTimeSec) * time.Second)
		swarmHealth := client.SwarmHealth()

		jsonDump, err := json.MarshalIndent(swarmHealth, "", "  ")

		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", jsonDump)

		err = ioutil.WriteFile("swarm_health.json", jsonDump, 0644)
	}()
}

var flags = struct {
	TestPeer []*net.TCPAddr `help:"addresses of some starting peers"`
	Addr     *net.TCPAddr   `help:"network listen addr"`
	tagflag.StartPos
	Torrent  []string `arity:"+" help:"torrent file path or magnet uri"`
}{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tagflag.Parse(&flags)
	var clientConfig torrent.Config

	clientConfig.DefaultStorage = storage.NewNoOpStorage()

	if flags.Addr != nil {
		clientConfig.ListenAddr = flags.Addr.String()
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
}
