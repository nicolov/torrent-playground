package torrent_nicolov

import (
	"crypto"
	"time"
	"expvar"
)

const (
	pieceHash        = crypto.SHA1
	maxRequests      = 250    // Maximum pending requests we allow peers to send us.
	defaultChunkSize = 0x4000 // 16KiB

	// Updated occasionally to when there's been some changes to client
	// behaviour in case other clients are assuming anything of us. See also
	// `bep20`.
	extendedHandshakeClientVersion = "go.torrent dev 20150624"
	// Peer ID client identifier prefix. We'll update this occasionally to
	// reflect changes to client behaviour that other clients may depend on.
	// Also see `extendedHandshakeClientVersion`.
	bep20 = "-GT0001-"

	nominalDialTimeout = time.Second * 30
	minDialTimeout     = 5 * time.Second

	// Justification for set bits follows.
	//
	// Extension protocol ([5]|=0x10):
	// http://www.bittorrent.org/beps/bep_0010.html
	//
	// Fast Extension ([7]|=0x04):
	// http://bittorrent.org/beps/bep_0006.html.
	// Disabled until AllowedFast is implemented.
	//
	// DHT ([7]|=1):
	// http://www.bittorrent.org/beps/bep_0005.html
	defaultExtensionBytes = "\x00\x00\x00\x00\x00\x10\x00\x01"

	defaultEstablishedConnsPerTorrent = 80
	defaultHalfOpenConnsPerTorrent    = 80
	torrentPeersHighWater             = 200
	torrentPeersLowWater              = 50

	// Limit how long handshake can take. This is to reduce the lingering
	// impact of a few bad apples. 4s loses 1% of successful handshakes that
	// are obtained with 60s timeout, and 5% of unsuccessful handshakes.
	handshakesTimeout = 20 * time.Second

	// These are our extended message IDs. Peers will use these values to
	// select which extension a message is intended for.
	metadataExtendedId = iota + 1 // 0 is reserved for deleting keys
	pexExtendedId
)

// I could move a lot of these counters to their own file, but I suspect they
// may be attached to a Client someday.
var (
	unwantedChunksReceived   = expvar.NewInt("chunksReceivedUnwanted_2")
	unexpectedChunksReceived = expvar.NewInt("chunksReceivedUnexpected_2")
	chunksReceived           = expvar.NewInt("chunksReceived_2")

	peersAddedBySource = expvar.NewMap("peersAddedBySource_2")

	uploadChunksPosted = expvar.NewInt("uploadChunksPosted_2")
	unexpectedCancels  = expvar.NewInt("unexpectedCancels_2")
	postedCancels      = expvar.NewInt("postedCancels_2")

	pieceHashedCorrect    = expvar.NewInt("pieceHashedCorrect_2")
	pieceHashedNotCorrect = expvar.NewInt("pieceHashedNotCorrect_2")

	unsuccessfulDials = expvar.NewInt("dialSuccessful_2")
	successfulDials   = expvar.NewInt("dialUnsuccessful_2")

	acceptUTP    = expvar.NewInt("acceptUTP_2")
	acceptTCP    = expvar.NewInt("acceptTCP_2")
	acceptReject = expvar.NewInt("acceptReject_2")

	peerExtensions                    = expvar.NewMap("peerExtensions_2")
	completedHandshakeConnectionFlags = expvar.NewMap("completedHandshakeConnectionFlags_2")
	// Count of connections to peer with same client ID.
	connsToSelf = expvar.NewInt("connsToSelf_2")
	// Number of completed connections to a client we're already connected with.
	duplicateClientConns       = expvar.NewInt("duplicateClientConns_2")
	receivedMessageTypes       = expvar.NewMap("receivedMessageTypes_2")
	receivedKeepalives         = expvar.NewInt("receivedKeepalives_2")
	supportedExtensionMessages = expvar.NewMap("supportedExtensionMessages_2")
	postedMessageTypes         = expvar.NewMap("postedMessageTypes_2")
	postedKeepalives           = expvar.NewInt("postedKeepalives_2")
	// Requests received for pieces we don't have.
	requestsReceivedForMissingPieces = expvar.NewInt("requestsReceivedForMissingPieces_2")

	// Track the effectiveness of Torrent.connPieceInclinationPool.
	pieceInclinationsReused = expvar.NewInt("pieceInclinationsReused_2")
	pieceInclinationsNew    = expvar.NewInt("pieceInclinationsNew_2")
	pieceInclinationsPut    = expvar.NewInt("pieceInclinationsPut_2")
)
