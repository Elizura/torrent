package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"simplebittorrent/common"
	"simplebittorrent/model"
)

func Seeder() {

	torrent, err := ParseTorrentFile("torrent-files/debian-11.6.0-amd64-netinst.iso.torrent")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", ":6881")

	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	for {
		if conn, err := ln.Accept(); err == nil {
			fmt.Println("Accepted connection")
			go handleSeedConnection(conn, torrent)
		}
	}
}

func handleSeedConnection(conn net.Conn, torrent model.Torrent) {
	conn.SetDeadline(time.Now().Add(common.PIECE_UPLOAD_TIMEOUT))
	defer conn.Close()
	_, err := ReceiveHandShake(conn)
	fmt.Println("Successfully Received handshake")
	if err != nil {
		fmt.Println("Error receiving handshake")
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	SendHandShake(conn, torrent)
	SendBitField(conn)
	ReceiveUnchoke(conn)
	ReceiveInterested(conn)
	SendUnchoke(conn)
	for {
		requestMsg, err := ReceiveRequest(conn)
		if err != nil {
			log.Fatal("found an error while receiving a seeder request")
			return
		}
		go handleRequest(*requestMsg, conn)
	}
}

func ReceiveHandShake(conn net.Conn) (*model.HandShake, error) {
	buffer := make([]byte, 68)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading handshake response")
		return &model.HandShake{}, err
	}
	handShake, err := model.DeserializeHandShake(buffer)
	if err != nil {
		fmt.Println("Error deserializing handshake response")
		return &model.HandShake{}, err
	}
	fmt.Println("Handshake sent successfully")
	return handShake, nil
}

func SendHandShake(conn net.Conn, torrent model.Torrent) error {

	clientIDByte := [20]byte{}
	copy(clientIDByte[:], []byte(common.CLIENT_ID))

	handshakeRequest := model.HandShake{
		Pstr:     "BitTorrent protocol",
		InfoHash: torrent.InfoHash,
		PeerID:   clientIDByte,
	}

	buffer := handshakeRequest.Serialize()

	_, err := conn.Write(buffer)
	if err != nil {
		fmt.Println("Error sending handshake request SEEDER")
		return err
	}
	return nil
}

func SendBitField(conn net.Conn) error {

	bitField := make([]byte, 255)
	for i := 0; i < len(bitField); i++ {
		bitField[i] = 255
	}

	msg := model.Message{MessageID: common.BIT_FIELD, Payload: bitField}
	_, err := conn.Write(msg.Serialize())
	if err != nil {
		return err
	}
	return nil
}

func ReceiveUnchoke(conn net.Conn) error {
	buffer := make([]byte, 5)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading unchoke message")
		return err
	}

	return nil
}

func ReceiveInterested(conn net.Conn) error {
	buffer := make([]byte, 5)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading interested message")
		return err
	}

	return nil
}

func ReceiveRequest(conn net.Conn) (*model.Message, error) {

	conn.SetDeadline(time.Now().Add(common.PIECE_UPLOAD_TIMEOUT))
	defer conn.SetDeadline(time.Time{})

	requestMsg, err := model.DeserializeMessage(conn)
	if err != nil {
		fmt.Println("Error reading request message")
		return &model.Message{}, err
	}

	if err != nil {
		fmt.Println("Error opening file")
		return &model.Message{}, err
	}

	return requestMsg, nil

}

func handleRequest(requestMsg model.Message, conn net.Conn) error {

	file, err := os.Open("downloads/debian-11.6.0-amd64-netinst.iso")
	defer file.Close()

	if requestMsg.MessageID != common.REQUEST {
		fmt.Println("Error: received message is not a request")
		return nil
	}
	index, begin, size, blockStart := ParseRequestPayload(requestMsg.Payload)

	piece := make([]byte, int64(size))
	_, err = file.ReadAt(piece, int64(begin))
	if err != nil {
		fmt.Println("Error reading piece from file")
		return err
	}

	err = SendPiece(conn, piece, index, blockStart)
	if err != nil {
		return err
	}

	return nil
}

func SendPiece(conn net.Conn, piece []byte, index int, blockStart int) error {
	payload := make([]byte, 8+len(piece))
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(blockStart))
	copy(payload[8:], piece[:])

	msg := model.Message{MessageID: common.PIECE, Payload: payload}
	_, err := conn.Write(msg.Serialize())
	if err != nil {
		fmt.Println("Error sending piece")
		return err
	}

	return nil
}

func ParseRequestPayload(payload []byte) (int, int, int, int) {
	index := int(binary.BigEndian.Uint32(payload[0:4]))
	blockStart := int(binary.BigEndian.Uint32(payload[4:8]))
	blockSize := int(binary.BigEndian.Uint32(payload[8:12]))
	pieceSize := 262144
	fileSize := 471859200
	begin := index*pieceSize + blockStart

	end := common.CalcMin(fileSize, begin+blockSize)

	if blockSize == 0 {
		fmt.Println("Error: block size is 0")
	}
	if end > fileSize {
		end = fileSize
		blockSize = end - begin
	}

	return index, begin, blockSize, blockStart
}

func SendUnchoke(conn net.Conn) {
	msg := model.Message{MessageID: common.UN_CHOKE, Payload: []byte{}}
	_, err := conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending unchoke message to peer: %s", err)
	}

}
