package servicelayer

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_servicelayer_mc_SetNodeId(t *testing.T) {
	s := &McData{}
	s.SetNodeId("node1")
	require.Equal(t, "node1", s.GetNodeId())
}

func Test_servicelayer_mc_SetLogFileName(t *testing.T) {
	s := &McData{}
	s.SetLogFileName("console.log")
	require.Equal(t, "console.log", s.GetLogFileName())
}

func Test_servicelayer_mc_SetLogsRefreshRate(t *testing.T) {
	s := &McData{}
	s.SetLogsRefreshRate(time.Minute)
	require.Equal(t, time.Minute, s.GetLogsRefreshRate())
}

func Test_servicelayer_mc_SetLastPageMarker(t *testing.T) {
	s := &McData{}
	s.SetLastPageMarker(1231122)
	var expected int64
	expected = 1231122
	require.Equal(t, expected, s.GetLastPageMarker())
}

