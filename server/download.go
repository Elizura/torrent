package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"simplebittorrent/common"
	"simplebittorrent/model"
	"strings"
	"time"

	// "github.com/inancgumus/screen"
)

type PieceResult struct {
	Index int    `bencode:"index"`
	Begin int    `bencode:"begin"`
	Block []byte `bencode:"block"`
}

type PieceRequest struct {
	Index  int      `bencode:"index"`
	Hash   [20]byte `bencode:"hash"`
	Length int      `bencode:"length"`
}

func PrepareDownload(filename string) (model.Torrent, []model.Peer) {

	torrent, err := ParseTorrentFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	peers, err := GetPeersFromTrackers(&torrent)
	if err != nil {
		log.Fatal(err)
	}

	return torrent, peers
}

func StartDownload(filename string) {

	torrent, peers := PrepareDownload(filename)
	

	outFile, filename, err := CreateFile(&torrent)
	if err != nil {
		log.Fatal("Error creating output file: ", err)
	}
	defer outFile.Close()

	piecesCache, err := LoadCache(outFile.Name() + ".json")
	if err != nil {
		panic(err)
	}

	piecesHashList := torrent.Info.PiecesToByteArray()
	downloadChannel := make(chan *PieceRequest, len(piecesHashList))
	resultChannel := make(chan *PieceResult)

	for idx, hash := range piecesHashList {
		length := torrent.CalculateRange(idx)

		downloadChannel <- &PieceRequest{Index: idx, Hash: hash, Length: length}
	}

	for _, peer := range peers {
		go DownloadFromPeer(peer, torrent, downloadChannel, resultChannel, piecesCache)
	}

	buf := make([]byte, torrent.Info.Length)
	donePieces := 0

	StoreDownloadedPieces(donePieces, torrent, resultChannel, err, outFile, piecesCache, buf)

	fmt.Println("Done downloading all pieces")
	close(downloadChannel)

}

func StoreDownloadedPieces(donePieces int, torrent model.Torrent, resultChannel chan *PieceResult, err error, outFile *os.File, piecesCache *model.PiecesCache, buf []byte) {

	for len(piecesCache.Pieces) < len(torrent.Info.PiecesToByteArray()) {
		res := <-resultChannel

		pieceSize := int(torrent.Info.PieceLength)
		pieceStartIdx := res.Index * pieceSize
		pieceEndIdx := common.CalcMin(pieceStartIdx+pieceSize, int(torrent.Info.Length))

		_, err = outFile.WriteAt(res.Block, int64(pieceStartIdx))
		if err != nil {
			log.Fatalf("Failed to write to file: %s", "downloaded_file.iso")
		}
		piecesCache.Pieces[res.Index] = true
		SaveCache(outFile.Name()+".json", piecesCache)
		copy(buf[pieceStartIdx:pieceEndIdx], res.Block)
		donePieces++

		percent := float64(len(piecesCache.Pieces)) / float64(len(torrent.Info.PiecesToByteArray())) * 100
		numWorkers := runtime.NumGoroutine() - 1
		// screen.Clear()
		// screen.MoveTopLeft()
		fmt.Println(strings.Repeat("=", int(percent)) + ">")
		log.Printf("Downloading... (%0.2f%%) Active Peers: %d\n", percent, numWorkers)
	}
	return
}

func DownloadFromPeer(peer model.Peer, torrent model.Torrent, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult, piecesCache *model.PiecesCache) {

	client, err := ClientFactory(peer, torrent)
	if err != nil {
		fmt.Printf("Failed to create a client with peer %s %s", peer.String(), err)
		return
	}

	client.UnChoke()
	client.Interested()

	for piece := range downloadChannel {
		fmt.Println("Found from cache: ", !piecesCache.Pieces[piece.Index])

		if common.BitOn(client.BitField, piece.Index) {

			_, err = DownloadPiece(piece, client, downloadChannel, resultChannel, &torrent)
			if err != nil {
				downloadChannel <- piece
				return
			}
		} else {
			downloadChannel <- piece
		}
	}
}

func DownloadPiece(piece *PieceRequest, client *model.Client, downloadChannel chan *PieceRequest, resultChannel chan *PieceResult, torrent *model.Torrent) (PieceResult, error) {

	client.Conn.SetDeadline(time.Now().Add(common.PIECE_DOWNLOAD_TIMEOUT))
	defer client.Conn.SetDeadline(time.Time{})

	totalDownloaded := 0
	requested := 0
	blockDownloadCount := 0
	blockLength := common.MAX_BLOCK_LENGTH

	buffer := make([]byte, piece.Length)

	for totalDownloaded < piece.Length {
		if client.ChokedState != common.CHOKE {
			for blockDownloadCount < common.MAX_BATCH_DOWNLOAD && requested < piece.Length {
				length := blockLength

				if piece.Length-requested < blockLength {
					length = piece.Length - requested
				}

				err := client.Request(uint32(piece.Index), uint32(requested), uint32(length))
				if err != nil {
					downloadChannel <- piece
					return PieceResult{}, err
				}
				requested += length
				blockDownloadCount++
			}
		}

		message, err := model.DeserializeMessage(client.Conn)
		if err != nil {
			downloadChannel <- piece
			return PieceResult{}, err
		}

		if message == nil {
			downloadChannel <- piece
			return PieceResult{}, err
		}

		switch message.MessageID {
		case common.CHOKE:
			client.ChokedState = common.CHOKE
		case common.UN_CHOKE:
			client.ChokedState = common.UN_CHOKE
		case common.INTERESTED:
			ParseInterested(message)
		case common.NOT_INTERESTED:
			ParseNotInterested(message)
		case common.HAVE:
			index, err := ParseHave(message)
			if err != nil {
				fmt.Println("Error parsing have message from peer: ", client.Peer.String())
				return PieceResult{}, err
			}
			common.TurnBitOn(client.BitField, index)
		case common.REQUEST:
			ParseRequest(message)
		case common.PIECE:
			n, err := ParsePiece(piece.Index, buffer, message)
			if err != nil {
				fmt.Println("Error parsing piece message from peer: ", client.Peer.String())
				downloadChannel <- piece
				return PieceResult{}, err
			}
			totalDownloaded += n
			blockDownloadCount--
		case common.CANCEL:
			ParseCancel(message)
		}

	}

	if !common.BitHashChecker(buffer, piece.Hash) {
		return PieceResult{}, fmt.Errorf("Piece hash verification failed for piece: %d", piece.Index)
	}

	resultChannel <- &PieceResult{Index: piece.Index, Block: buffer}

	return PieceResult{}, nil
}

func ParseInterested(msg *model.Message) {

	if msg.MessageID != common.INTERESTED {
		fmt.Printf("Expected INTERESTED (ID %d), got ID %d", common.INTERESTED, msg.MessageID)
	}
}

func ParseNotInterested(msg *model.Message) {

	if msg.MessageID != common.NOT_INTERESTED {
		fmt.Printf("Expected NOT_INTERESTED (ID %d), got ID %d", common.NOT_INTERESTED, msg.MessageID)
	}
}

func ParseHave(msg *model.Message) (int, error) {
	if msg.MessageID != common.HAVE {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", common.HAVE, msg.MessageID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload))

	return index, nil
}

func ParseRequest(msg *model.Message) {

}

func ParsePiece(index int, buf []byte, msg *model.Message) (int, error) {

	if msg.MessageID != common.PIECE {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", common.PIECE, msg.MessageID)
	}

	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}

	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		fmt.Println("begin problem")
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}

	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		fmt.Println("data problem: ", begin+len(data), " - ", len(buf))
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}

	copy(buf[begin:], data)

	return len(data), nil
}

func ParseCancel(msg *model.Message) {

	if msg.MessageID != common.CANCEL {
		fmt.Errorf("Expected CANCEL (ID %d), got ID %d", common.CANCEL, msg.MessageID)
	}
}
