package core

import (
	"fmt"
	"os"
	"time"

	s "github.com/application-research/edge-vertex/services"
	logging "github.com/ipfs/go-log/v2"
)

type EdgeDaemon struct {
	interval         int
	edgeListFilename string
	attemptedCids    map[string]bool // Map of CIDs that have already been queried - to avoid duplicate queries to DDM
	// TODO: this could use a lot of memory as it's indexing all CIDs across all edges - maybe use a bloom filter or some other approach instead?
}

var (
	log = logging.Logger("edge-daemon")
)

func NewEdgeDaemon(interval int, edgeListFilename string) *EdgeDaemon {
	// Check to ensure that the edgeListFilename exists
	if _, err := os.Stat(edgeListFilename); os.IsNotExist(err) {
		log.Fatalf("Edge list file %s does not exist", edgeListFilename)
	}

	return &EdgeDaemon{
		interval:         interval,
		edgeListFilename: edgeListFilename,
		attemptedCids:    make(map[string]bool),
	}
}

func (ed *EdgeDaemon) Run() {
	for {
		cnt, err := ed.aggregateContent()
		if err != nil {
			log.Errorf("error aggregating content: %s", err)
		}

		ed.publishContent(cnt)
		time.Sleep(time.Duration(ed.interval))
	}
}

// TODO: run in a separate goroutine for each Edge for better performance
func (ed *EdgeDaemon) aggregateContent() (map[string]s.BucketContent, error) {
	result := make(map[string]s.BucketContent)

	// Read the edge.json file from disk at the beginning of each run, so it can be modified without restarting the daemon
	edgeAddrs, err := s.ReadEdgeList(ed.edgeListFilename)

	if err != nil {
		return result, fmt.Errorf("error reading edge list: %s", err)
	}

	var edges []s.Edge
	for _, e := range edgeAddrs {
		edges = append(edges, *s.NewEdge(e))
	}

	// Loop through the list of []edge URIs, for each one query it for available contents
	// Aggregate them all into a single list of all contents - removing duplicates that have the same CID
	for _, e := range edges {
		cnt, err := e.GetOpenBuckets()

		if err != nil {
			log.Errorf("error getting open buckets for edge %s: %s", e.URI, err)
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

func (ed *EdgeDaemon) publishContent(map[string]s.BucketContent) {
	//TODO
	// for each content, POST it to DDM /contents endpoint for the given tag

	// write output to specified log file for each of the contents
}
