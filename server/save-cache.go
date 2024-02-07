package server

import (
	"encoding/json"
	"fmt"
	"os"

	"simplebittorrent/model"
)

func SaveCache(filename string, cache *model.PiecesCache) error {
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("error different from non existent")
		return err
	}

	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error opening cache while non existent")
			return err
		}
		file.Close()
	}
	file, err := os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Error opening cache")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cache)
	if err != nil {
		fmt.Println("Error encoding cache", err)
		return err
	}

	return nil
}

func LoadCache(filename string) (*model.PiecesCache, error) {

	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error while creating file", filename)
			return nil, err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.Encode(&model.PiecesCache{Pieces: map[int]bool{}})
		return &model.PiecesCache{Pieces: map[int]bool{}}, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}

	decoder := json.NewDecoder(file)
	var cache model.PiecesCache
	err = decoder.Decode(&cache)
	if err != nil {
		fmt.Println("error decoding file")
		return nil, err
	}

	return &cache, nil
}
