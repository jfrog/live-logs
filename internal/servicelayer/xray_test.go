package servicelayer

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_servicelayer_xray_SetNodeId(t *testing.T) {
	s := &XrayData{}
	s.SetNodeId("node1")
	require.Equal(t, "node1", s.GetNodeId())
}

func Test_servicelayer_xray_SetLogFileName(t *testing.T) {
	s := &XrayData{}
	s.SetLogFileName("console.log")
	require.Equal(t, "console.log", s.GetLogFileName())
}

func Test_servicelayer_xray_SetLogsRefreshRate(t *testing.T) {
	s := &XrayData{}
	s.SetLogsRefreshRate(time.Minute)
	require.Equal(t, time.Minute, s.GetLogsRefreshRate())
}

func Test_servicelayer_xray_SetLastPageMarker(t *testing.T) {
	s := &XrayData{}
	s.SetLastPageMarker(1231122)
	var expected int64
	expected = 1231122
	require.Equal(t, expected, s.GetLastPageMarker())
}

