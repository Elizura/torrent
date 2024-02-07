package model

import (
	"encoding/binary"
	"errors"
	"simplebittorrent/common"

	"log"
	"net"
	"syscall"
)

type Client struct {
	Conn        net.Conn
	Peer        Peer
	BitField    []byte
	ChokedState uint8
}

func (client *Client) Interested() {
	msg := Message{MessageID: common.INTERESTED, Payload: []byte{}}
	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending interested message to peer: %s", err)
	}
}

func (client *Client) Choke() {
	msg := Message{MessageID: common.CHOKE, Payload: []byte{}}
	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending choke message to peer: %s", err)
	}
}

func (client *Client) UnChoke() {
	msg := Message{MessageID: common.UN_CHOKE, Payload: []byte{}}
	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		log.Fatalf("Error sending unchoke message to peer: %s", err)
	}

}

func (client *Client) Request(index uint32, begin uint32, length uint32) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], index)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], length)
	msg := Message{MessageID: common.REQUEST, Payload: payload}

	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {

		}

		return err
	}

	return nil
}

func (client *Client) Have(index uint32) {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload[0:4], index)
	msg := Message{MessageID: common.HAVE, Payload: payload}

	_, err := client.Conn.Write(msg.Serialize())
	if err != nil {
	}
}
