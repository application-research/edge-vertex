package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Edge struct {
	URI string
}

func NewEdge(uri string) *Edge {
	return &Edge{URI: uri}
}

func (e *Edge) GetOpenBuckets() (*OpenBucketsResponse, error) {
	resp, err := http.Get(e.URI + "/buckets/get-open")
	if err != nil {
		return nil, fmt.Errorf("could not reach edge api: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err)
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("error in edge api call %d : %s", resp.StatusCode, body)
	}

	defer resp.Body.Close()

	result, err := UnmarshalOpenBucketsResponse(body)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal open buckets response %s : %s", err, string(body))
	}

	return &result, nil
}

func UnmarshalOpenBucketsResponse(data []byte) (OpenBucketsResponse, error) {
	var r OpenBucketsResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

type OpenBucketsResponse []BucketContent
type BucketContent struct {
	BucketUUID  string        `json:"bucket_uuid"`
	PieceCID    string        `json:"piece_cid"`
	PayloadCID  string        `json:"payload_cid"`
	DirCID      string        `json:"dir_cid"`
	PieceSize   uint          `json:"piece_size"`
	DownloadURL string        `json:"download_url"`
	Status      ContentStatus `json:"status"`
	Size        uint          `json:"size"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Collection  string        `json:"collection"`
}

type ContentStatus string

const (
	Processing ContentStatus = "processing"
	Open       ContentStatus = "open"
	Ready      ContentStatus = "ready"
)

func (o *OpenBucketsResponse) StatusReady() OpenBucketsResponse {
	var ready OpenBucketsResponse

	for _, content := range *o {
		if content.Status == Ready {
			ready = append(ready, content)
		}
	}

	return ready
}
