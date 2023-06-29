# Edge-Vertex

Edge-URID aggregator

It runs as a daemon, queries a list of `edge-urid` instances and pushes their contents into DDM for deal management


## Setup

1. specify an `edge.json` file containing a list of the edge URIs to be aggregated into ddm
ex) 
```json
[
  "https://dev-hivemapper.edge.estuary.tech",
  "https://dev-site2.edge.estuary.tech"
]
```

> Note: This file can be edited at any time to add or remove edge URIs, a daemon restart is not required

2. run edge-vertex and provide env vars or flags so it can access DDM api
```bash
edge-vertex \
  --ddm-api https://delta-dm.estuary.tech \
  --ddm-token <token> \
  --interval 60 \ # optional, defaults to 60 seconds
  --edge-file edge.json \ # optional, defaults to edge.json
  --debug  # optional, more verbose logging
```

3. Check `log output` for log output if anything fails

## DDM Set-Up
- Tag that is present on `edge` must also exist as a `dataset name` in DDM, otherwise it will simply not add the content
- untagged content is not supported yet