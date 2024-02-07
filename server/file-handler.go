package server

import (
	"fmt"
	"os"
	"path/filepath"

	"simplebittorrent/model"
)

func CreateFile(torrent *model.Torrent) (*os.File, string, error) {
	outFile, err := CreateOrOpenFile("downloads" + "/" + torrent.Info.Name)
	if err != nil {
		return nil, "", err
	}
	return outFile, "", err
}

func CreateOrOpenFile(filename string) (*os.File, error) {

	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {

		dir := filepath.Dir(filename)
		if _, err := os.Stat(dir); os.IsNotExist(err) {

			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				fmt.Println("Error while creating directory:", dir)
				return nil, err
			}
		}

		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error while creating file", filename)
			return nil, err
		}
		return file, nil
	}

	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}

	return file, nil
}
