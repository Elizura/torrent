package server

import (
	"bytes"
	"net"
	"time"

	"simplebittorrent/common"
	"simplebittorrent/model"
)

func ClientFactory(peer model.Peer, torrent model.Torrent) (*model.Client, error) {
	client, err := createClient(peer, torrent)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createClient(peer model.Peer, torrent model.Torrent) (*model.Client, error) {
	conn, err := connectToPeer(peer, torrent)
	if err != nil {
		return nil, err
	}

	err = ShakeHandWithPeer(torrent, peer, common.CLIENT_ID, conn)
	if err != nil {
		return nil, err
	}

	bitFieldMessage, err := ReceiveBitFieldMessage(conn)
	if err != nil {
		return &model.Client{}, err
	}

	client := &model.Client{
		Peer:        peer,
		BitField:    bitFieldMessage.Payload,
		Conn:        conn,
		ChokedState: common.CHOKE,
	}

	return client, nil
}

func connectToPeer(peer model.Peer, torrent model.Torrent) (net.Conn, error) {

	conn, err := net.DialTimeout("tcp", peer.String(), common.CONNECTION_TIMEOUT)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ShakeHandWithPeer(torrent model.Torrent, peer model.Peer, clientID string, conn net.Conn) error {

	conn.SetDeadline(time.Now().Add(common.CONNECTION_TIMEOUT))
	defer conn.SetDeadline(time.Time{})

	clientIDByte := [20]byte{}
	copy(clientIDByte[:], []byte(clientID))

	handshakeRequest := model.HandShake{
		Pstr:     "BitTorrent protocol",
		InfoHash: torrent.InfoHash,
		PeerID:   clientIDByte,
	}

	handshakeResponse, err := handshakeRequest.Send(conn)
	if err != nil {
		return err
	}

	if !bytes.Equal(handshakeResponse.InfoHash[:], torrent.InfoHash[:]) {
		return err
	}

	if bytes.Equal(handshakeResponse.PeerID[:], common.ConvertStringToByteArray(common.CLIENT_ID)[:]) {
		return err
	}

	return nil
}

func ReceiveBitFieldMessage(conn net.Conn) (*model.Message, error) {
	conn.SetDeadline(time.Now().Add(common.CONNECTION_TIMEOUT))
	defer conn.SetDeadline(time.Time{})

	bitFieldMessageResponse, err := model.DeserializeMessage(conn)
	if err != nil {

		return nil, err
	}

	if bitFieldMessageResponse.MessageID != common.BIT_FIELD {

		return &model.Message{}, nil

	}

	return bitFieldMessageResponse, nil
}
