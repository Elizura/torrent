package server

import (
	"net/http"
	"net/url"
	"simplebittorrent/common"
	"simplebittorrent/model"

	// "net"
	"github.com/zeebo/bencode"
)

func GetPeersFromTrackers(torrent *model.Torrent) ([]model.Peer, error) {

	httpTrackerURLs, err := getTrackerUrl(torrent)
	if err != nil {
		return nil, err
	}

	peers, err := getPeersFromTrackers(httpTrackerURLs)
	if err != nil {
		return nil, err
	}
	// peers = []model.Peer{
	// 	{IP: net.IP([]byte{192, 168, 82, 111}), Port: 6881},
	// 	{IP: net.IP([]byte{192, 168, 202, 94}), Port: 6881},
	// }

	return peers, nil
}

func getTrackerUrl(torrent *model.Torrent) ([]string, error) {

	requestParams := model.TrackerRequestParams{
		Info_hash:  torrent.InfoHash,
		Peer_id:    common.CLIENT_ID,
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Left:       torrent.Info.Length,
		Compact:    1,
		Event:      "started",
	}

	URLs := []string{}

	URL, err := url.Parse(torrent.Announce)
	if err != nil {
		return []string{}, err
	}
	URL.RawQuery = requestParams.Encode()
	URLs = append(URLs, URL.String())

	//check if there are other trackers and collect
	for _, tracker := range torrent.AnnounceList {

		URL, err := url.Parse(tracker[0])
		if err != nil {
			return []string{}, err
		}

		URL.RawQuery = requestParams.Encode()
		URLs = append(URLs, URL.String())

	}

	return URLs, nil
}

func getPeersFromTrackers(URLs []string) ([]model.Peer, error) {
	peers := []model.Peer{}
	for _, URL := range URLs {
		response, err := getPeerFromURL(URL)
		if err == nil && len(response) > 0 {
			for _, p := range response {
				peers = append(peers, p)
			}
		}
	}

	return peers, nil
}

func getPeerFromURL(URL string) ([]model.Peer, error) {

	response, err := http.Get(URL)
	if err != nil {
		return []model.Peer{}, err
	}

	defer response.Body.Close()

	trackerResponse := model.TrackerResponse{}
	err = bencode.NewDecoder(response.Body).Decode(&trackerResponse)
	if err != nil {
		return []model.Peer{}, err
	}

	peers, err := model.PeerParser([]byte(trackerResponse.Peers))
	if err != nil {
		return []model.Peer{}, err
	}

	return peers, nil
}
