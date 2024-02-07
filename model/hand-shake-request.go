package model

import (
	"fmt"
	"net"
)

type HandShake struct {
	Pstr     string   `bencode:"pstr"`
	InfoHash [20]byte `bencode:"info_hash"`
	PeerID   [20]byte `bencode:"peer_id"` 
}



func (handShake *HandShake) Serialize() []byte {
	
	buffer := make([]byte, 49+len(handShake.Pstr))
	buffer[0] = byte(len(handShake.Pstr))
	copy(buffer[1:], []byte(handShake.Pstr))
	copy(buffer[1+len(handShake.Pstr):], make([]byte, 8))
	copy(buffer[1+len(handShake.Pstr)+8:], handShake.InfoHash[:])
	copy(buffer[1+len(handShake.Pstr)+8+20:], handShake.PeerID[:])

	return buffer
}


func DeserializeHandShake(buffer []byte) (*HandShake, error) {
	
	handShake := &HandShake{}
	pstrLength := int(buffer[0])
	
	if pstrLength != 19 {
		fmt.Println("pstr length is not 19")
		return &HandShake{}, fmt.Errorf("pstr length is not 19")
	}
	handShake.Pstr = string(buffer[1 : pstrLength+1])
	
	copy(handShake.InfoHash[:], buffer[28:48])
	copy(handShake.PeerID[:], buffer[48:68])

	return handShake, nil
}

func (h *HandShake) Send(conn net.Conn) (*HandShake, error) {
	
	buffer := h.Serialize()

	
	_, err := conn.Write(buffer)
	if err != nil {
		return &HandShake{}, err
	}

	
	
	buffer = make([]byte, 68)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading handshake response")
		return &HandShake{}, err
	}

	
	handShake, err := DeserializeHandShake(buffer)
	if err != nil {
		fmt.Println("Error deserializing handshake response")
		return &HandShake{}, err
	}

	fmt.Println("Handshake sent successfully")
	return handShake, nil
}
