package model

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP `bencode:"ip"`
	Port uint16   `bencode:"port"`
}

func (peer *Peer) String() string {
	return peer.IP.String() + ":" + strconv.Itoa(int(peer.Port))
}

func PeerParser(peerByte []byte) ([]Peer, error) {

	numOfPeers := len(peerByte) / 6
	peers := make([]Peer, numOfPeers)

	if len(peerByte)%6 != 0 {
		return []Peer{}, errors.New("invalid peer byte array")
	}

	for i := 0; i < numOfPeers; i++ {

		ip := net.IP(peerByte[i*6 : i*6+4])

		port := binary.BigEndian.Uint16([]byte(peerByte[i*6+4 : i*6+6]))
		peers[i] = Peer{IP: ip, Port: port}
	}

	return peers, nil
}