package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EdgeApi interface {
	GetContents(edgeUri string) (*OpenBucketsResponse, error)
}

type EdgeUrid struct{}

func (e *EdgeUrid) GetContents(edgeUri string) (*OpenBucketsResponse, error) {
	resp, err := http.Get(edgeUri + "/buckets/get/open")
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

	rawOpenBuckets, err := UnmarshalOpenBucketsResponse(body)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal open buckets response %s : %s", err, string(body))
	}

	result := prependDownloadUrl(edgeUri, &rawOpenBuckets)

	return result, nil
}

// Includes the edge URL in the download URL, ex) /gw/baga123 -> http://edge.com/gw/baga123
func prependDownloadUrl(baseUrl string, content *OpenBucketsResponse) *OpenBucketsResponse {
	for i := range *content {
		(*content)[i].DownloadURL = baseUrl + (*content)[i].DownloadURL
	}

	return content
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
	Collection  string        `json:"collection_name"`
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
