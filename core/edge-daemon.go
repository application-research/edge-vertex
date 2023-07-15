package core

import (
	"fmt"
	"os"
	"time"

	s "github.com/application-research/edge-vertex/services"
	"github.com/application-research/edge-vertex/util"
	log "github.com/sirupsen/logrus"
)

type EdgeDaemon struct {
	interval         time.Duration
	edgeListFilename string
	attemptedCids    map[string]bool // Map of CIDs that have already been queried - to avoid duplicate queries to DDM
	DDM              *s.DDMApi
	listReader       s.ListReader
	edgeApi          s.EdgeApi
	totalSuccess     int
	totalFail        int
	// TODO: this could use a lot of memory as it's indexing all CIDs across all edges - maybe use a bloom filter or some other approach instead?
}

func NewEdgeDaemon(interval int, edgeListFilename string, ddmUrl string, ddmKey string) *EdgeDaemon {
	// Check to ensure that the edgeListFilename exists
	if _, err := os.Stat(edgeListFilename); os.IsNotExist(err) {
		log.Fatalf("Edge list file %s does not exist", edgeListFilename)
	}

	return &EdgeDaemon{
		interval:         time.Duration(time.Second) * time.Duration(interval),
		edgeListFilename: edgeListFilename,
		listReader:       &s.JsonEdgeListReader{Filename: edgeListFilename},
		edgeApi:          &s.EdgeUrid{},
		attemptedCids:    make(map[string]bool),
		DDM:              s.NewDDMApi(ddmUrl, ddmKey),
	}
}

func (ed *EdgeDaemon) Run() {
	for {
		cnt, err := ed.aggregateContent()
		if err != nil {
			log.Errorf("error aggregating content: %s", err)
			continue
		}

		success, fail := ed.publishContent(cnt)
		log.Infof(util.Gray+"published %d contents to DDM, %d succeeded, %d failed"+util.Reset, len(cnt), success, fail)

		ed.totalSuccess += success
		ed.totalFail += fail

		log.Infof(util.Green+"Total count success: %d, fail: %d"+util.Reset, ed.totalSuccess, ed.totalFail)
		time.Sleep(ed.interval)
	}
}

// Collects a list of contents from all Edges, deduplicates them into a map
// TODO: run in a separate goroutine for each Edge for better performance
func (ed *EdgeDaemon) aggregateContent() (map[string]s.BucketContent, error) {
	result := make(map[string]s.BucketContent)

	// Read the edge.json file from disk at the beginning of each run, so it can be modified without restarting the daemon
	edgeAddrs, err := ed.listReader.Read()

	if err != nil {
		return result, fmt.Errorf("error reading edge list: %s", err)
	}

	// Loop through the list of []edge URIs, for each one query it for available contents
	// Aggregate them all into a single list of all contents - removing duplicates that have the same CID
	for _, e := range edgeAddrs {
		cnt, err := ed.edgeApi.GetContents(e)

		if err != nil {
			log.Errorf("error getting open buckets for edge %s: %s", e, err)
			continue
		}

		filteredCnt := cnt.StatusReady()

		// insert into `result` if not already there (map key is PieceCID)
		for _, cnt := range filteredCnt {
			if _, ok := ed.attemptedCids[cnt.PieceCID]; ok {
				continue
			}

			if _, ok := result[cnt.PieceCID]; !ok {
				result[cnt.PieceCID] = cnt
				ed.attemptedCids[cnt.PieceCID] = true
			}
		}
	}

	return result, nil
}

// Publishes the list of contents to DDM
func (ed *EdgeDaemon) publishContent(contents map[string]s.BucketContent) (int, int) {
	var request s.ContentsRequest

	// Create a DDM content list from the edge content list
	for _, cnt := range contents {
		request = append(request, s.DDMContent{
			CommP:           cnt.PieceCID,
			PayloadCID:      cnt.PayloadCID,
			PaddedSize:      cnt.PieceSize,
			Size:            cnt.Size,
			Collection:      cnt.Collection,
			ContentLocation: cnt.DownloadURL,
		})
	}

	// Publish the list of contents to DDM
	res, err := ed.DDM.PublishContent(&request)
	if err != nil {
		log.Errorf("error publishing content to DDM: %s", err)
		return 0, 0
	}

	return len(res.Success), len(res.Fail)
}
