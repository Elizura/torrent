package server

import (
	"fmt"
	"os"

	bencode "github.com/zeebo/bencode"

	"simplebittorrent/model"
)

func ParseTorrentFile(filename string) (model.Torrent, error) {

	file, err := os.Open(filename)
	if err != nil {
		return model.Torrent{}, err
	}
	defer file.Close()

	var torrent = model.Torrent{}
	err = bencode.NewDecoder(file).Decode(&torrent)
	fmt.Println("the announce", torrent.Announce)
	fmt.Println("the announce list", torrent.AnnounceList)
	fmt.Println("the Encoding", torrent.Encoding)
	fmt.Println("the Info files", torrent.Info.Files)
	fmt.Println("the Info length", torrent.Info.Length)
	fmt.Println("the Info Name", torrent.Info.Name)
	fmt.Println("the Info Length", torrent.Info.PieceLength)
	fmt.Println("the Info Private", torrent.Info.Private)
	// fmt.Println("the Info hash", torrent.InfoHash)
	if err != nil {
		fmt.Println("Encountered error while decoding")
		return model.Torrent{}, err
	}

	torrent.GenerateInfoHash()
	fmt.Println("Torrent file parsed successfully")
	return torrent, nil
}
