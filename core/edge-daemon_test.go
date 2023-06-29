package core

import (
	"testing"

	s "github.com/application-research/edge-vertex/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockListReader struct {
	mock.Mock
}

func (m *MockListReader) Read() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

type MockEdgeApi struct {
	mock.Mock
}

func (m *MockEdgeApi) GetContents(edgeUri string) (*s.OpenBucketsResponse, error) {
	args := m.Called(edgeUri)
	return args.Get(0).(*s.OpenBucketsResponse), args.Error(1)
}

func TestAggregateContent(t *testing.T) {
	// Create a mock list reader
	mockListReader := new(MockListReader)
	mockListReader.On("Read").Return([]string{"http://example.com/edge1", "http://example.com/edge2"}, nil)

	// Create a mock edge API
	mockEdgeApi := new(MockEdgeApi)
	mockEdgeApi.On("GetContents", "http://example.com/edge1").Return(&s.OpenBucketsResponse{
		{
			PieceCID:   "abc123",
			PayloadCID: "def456",
			Status:     "ready",
		},
	}, nil)
	mockEdgeApi.On("GetContents", "http://example.com/edge2").Return(&s.OpenBucketsResponse{
		{
			PieceCID:   "abc123",
			PayloadCID: "def456",
			Status:     "ready",
		}, // Will get ignored as it's a duplicate
		{
			PieceCID:   "ghi789",
			PayloadCID: "jkl012",
			Status:     "ready",
		},
	}, nil)

	// Create an EdgeDaemon with the mock dependencies
	ed := &EdgeDaemon{
		interval:         1,
		edgeListFilename: "edge.json",
		attemptedCids:    make(map[string]bool),
		DDM:              nil,
		listReader:       mockListReader,
		edgeApi:          mockEdgeApi,
		totalSuccess:     0,
		totalFail:        0,
	}

	// Call the function being tested
	result, err := ed.aggregateContent()

	// Check the result
	assert.NoError(t, err)
	assert.Equal(t, map[string]s.BucketContent{
		"abc123": {
			PieceCID:   "abc123",
			PayloadCID: "def456",
			Status:     "ready",
		},
		"ghi789": {
			PieceCID:   "ghi789",
			PayloadCID: "jkl012",
			Status:     "ready",
		},
	}, result)

	// Check that the mock objects were called as expected
	mockListReader.AssertCalled(t, "Read")
	mockEdgeApi.AssertCalled(t, "GetContents", "http://example.com/edge1")
	mockEdgeApi.AssertCalled(t, "GetContents", "http://example.com/edge2")
}
