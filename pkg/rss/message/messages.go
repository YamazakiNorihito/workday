package message

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/YamazakiNorihito/workday/internal/domain/rss"
)

const MaxMessageSize = 256 * 1024

type Subscribe struct {
	FeedURL        string `json:"feed_url"`
	Language       string `json:"language"`
	rss.ItemFilter `json:"item_filter"`
}

type Write struct {
	RssFeed    rss.Rss `json:"rss,omitempty"`
	Compressed bool    `json:"compressed"`
	Data       []byte  `json:"data,omitempty"`
}

func NewWriteMessage(entryRss rss.Rss) (writeMessage Write, err error) {
	serializedRss, _ := json.Marshal(entryRss)
	if len(serializedRss) > MaxMessageSize {
		compressedRssData, err := compressAndEncodeData(serializedRss)
		if err != nil {
			return Write{}, err
		}

		writeMessage = Write{
			Compressed: true,
			Data:       compressedRssData,
		}
	} else {
		writeMessage = Write{
			Compressed: false,
			RssFeed:    entryRss,
		}
	}
	return writeMessage, nil
}

func compressAndEncodeData(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	if _, err := gzipWriter.Write(data); err != nil {
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return []byte(base64.StdEncoding.EncodeToString(buffer.Bytes())), nil
}

func DecodeAndDecompressData(compressedData []byte) (rss.Rss, error) {
	decodedData, err := base64.StdEncoding.DecodeString(string(compressedData))
	if err != nil {
		return rss.Rss{}, err
	}

	buffer := bytes.NewBuffer(decodedData)
	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return rss.Rss{}, err
	}
	defer gzipReader.Close()

	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return rss.Rss{}, err
	}

	var rssFeed rss.Rss
	err = json.Unmarshal(decompressedData, &rssFeed)
	if err != nil {
		return rss.Rss{}, err
	}

	return rssFeed, nil
}
