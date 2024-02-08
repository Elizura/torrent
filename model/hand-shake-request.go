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

//format to a similar format
// \x13BitTorrent protocol\x00\x00\x00\x00\x00\x00\x00\x00\x86\xd4\xc8\x00\x24\xa4\x69\xbe\x4c\x50\xbc\x5a\x10\x2c\xf7\x17\x80\x31\x00\x74-TR2940-k8hj0wgej6ch
func (handShake *HandShake) Serialize() []byte {
	
	buffer := make([]byte, 49+len(handShake.Pstr))
	//put the pstr len -> 1 byte
	buffer[0] = byte(len(handShake.Pstr))
	//copy pstr -> BitTorrent protocol
	copy(buffer[1:], []byte(handShake.Pstr))
	// 8 bytes reserve
	copy(buffer[1+len(handShake.Pstr):], make([]byte, 8))
	//copy Info Hash
	copy(buffer[1+len(handShake.Pstr)+8:], handShake.InfoHash[:])
	// copy the peer Id
	copy(buffer[1+len(handShake.Pstr)+28:], handShake.PeerID[:])

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
