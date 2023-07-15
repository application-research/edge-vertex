package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DDMApi struct {
	baseUrl string
	authKey string
}

type ContentsRequest []DDMContent
type DDMContent struct {
	PayloadCID      string `json:"payload_cid"`
	CommP           string `json:"commp"`
	PaddedSize      uint   `json:"padded_size"`
	Size            uint   `json:"size"`
	ContentLocation string `json:"content_location"`
	Collection      string `json:"collection"`
}

type ContentsResponse struct {
	Success []string `json:"success"`
	Fail    []string `json:"fail"`
}

func UnmarshalContentsResponse(data []byte) (ContentsResponse, error) {
	var r ContentsResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func NewDDMApi(baseUrl string, authKey string) *DDMApi {
	d := &DDMApi{baseUrl: baseUrl, authKey: authKey}

	err := d.healthCheck()
	if err != nil {
		log.Fatalf("could not connect to ddm %s", err)
	}

	return d
}

func (d *DDMApi) healthCheck() error {
	req, err := http.NewRequest(http.MethodGet, d.baseUrl+"/api/v1/health", nil)
	if err != nil {
		return fmt.Errorf("could not construct http request %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+d.authKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not make http request %s", err)
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf(string(body))
	}

	return err
}

// Publishes a list of contents to DDM
func (d *DDMApi) PublishContent(c *ContentsRequest) (*ContentsResponse, error) {

	rawCnt, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("could not marshal content %s", err)
	}

	res, closer, err := d.postRequest("/api/v1/contents", rawCnt, d.authKey)
	if err != nil {
		return nil, fmt.Errorf("could not make request %s", err)
	}
	defer closer()

	result, err := UnmarshalContentsResponse(res)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response %s", err)
	}

	return &result, nil

}

func (d *DDMApi) postRequest(url string, raw []byte, authKey string) ([]byte, func() error, error) {
	if authKey == "" {
		return nil, nil, fmt.Errorf("auth token must be provided")
	}

	req, err := http.NewRequest(http.MethodPost, d.baseUrl+url, bytes.NewBuffer(raw))
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct http request %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("could not make http request %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("error in ddm call %d : %s", resp.StatusCode, body)
	}

	if err != nil {
		return nil, nil, err
	}

	return body, resp.Body.Close, nil
}
