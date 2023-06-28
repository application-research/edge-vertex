package core

func EdgePollingDaemon(interval int) {
	// Read the edge.json file from disk at the beginning of each run, so it can be modified without restarting the daemon

	// Loop through the list of []edge URIs, for each one query it for available contents

	// for each content, POST it to DDM /contents endpoint for the given tag
	// TODO:DDM- add ability to POST contents by dataset (tag) name instead of id

	// write output to specified log file for each of the contents

	// sleep for interval
}
