package livelog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jfrog/live-logs/internal/constants"
	"github.com/jfrog/live-logs/internal/model"
	"github.com/jfrog/live-logs/internal/servicelayer"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_LiveLogs_SetProductId(t *testing.T) {
	s := &Data{}
	s.SetProductId("rt")
	require.Equal(t, "rt", s.GetProductId())
}

func Test_LiveLogs_SetServiceId(t *testing.T) {
	s := &Data{}
	s.SetServiceId("test-rt")
	require.Equal(t, "test-rt", s.GetServiceId())
}

func Test_LiveLogs_SetLogsRefreshRate(t *testing.T) {
	s := &Data{}
	s.SetLogsRefreshRate(time.Minute)
	require.Equal(t, time.Minute, s.GetLogsRefreshRate())
}

func Test_LiveLogs_CatLog(t *testing.T) {
	tests := []struct {
		name            string
		logFileName     string
		nodeId          string
		productId       string
		serviceId       string
		logsRefreshRate time.Duration
		lastPageMarker  int64
		mockGetLogResponse model.Data
		mockGetConfigResponse *model.Config
		mockGetErr      error
		want            string
		wantErr         bool
	}{
		{
			name:    "missing node id",
			wantErr: true,
		},
		{
			name:    "missing log file name",
			nodeId:  "node-1",
			wantErr: true,
		},
		{
			name:        "error response",
			nodeId:      "node-1",
			logFileName: "one.log",
			mockGetErr:  fmt.Errorf("some-error"),
			wantErr:     true,
		},
		{
			name:            "cat log response",
			nodeId:          "node-1",
			logFileName:     "one.log",
			mockGetLogResponse: model.Data{ Content : "some log content", PageMarker: 123},
			want:            "some log content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Data{
				serviceLayerClient: &mockServiceLayer{
					t:                  t,
					expectNodeId:       tt.nodeId,
					expectLogFileName:  tt.logFileName,
					getLogResponse:     tt.mockGetLogResponse,
					getConfigResponse:  tt.mockGetConfigResponse,
					getErr:             tt.mockGetErr,
				},
				productId:       tt.productId,
				serviceId:        tt.serviceId,
				logsRefreshRate: time.Second,
			}
			out := &bytes.Buffer{}
			err := s.CatLog(context.Background(), out)
			if (err != nil) != tt.wantErr {
				t.Errorf("CatLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(out.String(), tt.want) {
				t.Errorf("CatLog() got = %v, want %v", out.String(), tt.want)
			}
		})
	}
}

func Test_LiveLogs_PrintLogs(t *testing.T) {
	tests := []struct {
		name            string
		logFileName     string
		nodeId          string
		productId       string
		serviceId       string
		logsRefreshRate time.Duration
		lastPageMarker  int64
		mockGetLogResponse model.Data
		mockGetConfigResponse *model.Config
		mockGetErr      error
		want            string
		wantErr         bool
	}{
		{
			name:    "missing node id",
			wantErr: true,
		},
		{
			name:    "missing log file name",
			nodeId:  "node-1",
			wantErr: true,
		},
		{
			name:        "error response",
			nodeId:      "node-1",
			logFileName: "one.log",
			mockGetErr:  fmt.Errorf("some-error"),
			wantErr:     true,
		},
		{
			name:            "cat log response",
			nodeId:          "node-1",
			logFileName:     "one.log",
			mockGetLogResponse: model.Data{ Content : "some log content", PageMarker: 123},
			want:            "some log content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Data{
				serviceLayerClient: &mockServiceLayer{
					t:                  t,
					expectNodeId:       tt.nodeId,
					expectLogFileName:  tt.logFileName,
					getLogResponse:     tt.mockGetLogResponse,
					getConfigResponse:  tt.mockGetConfigResponse,
					getErr:             tt.mockGetErr,
				},
				productId:       tt.productId,
				serviceId:        tt.serviceId,
				logsRefreshRate: time.Second,
			}
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := s.PrintLogs(context.Background(), tt.nodeId, tt.logFileName, false)

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(out), tt.want) {
				t.Errorf("PrintLogs() got = %v, want %v", string(out), tt.want)
			}
		})
	}
}

func Test_LiveLogs_DisplayConfig(t *testing.T) {
	tests := []struct {
		name            string
		logFileName     string
		nodeId          string
		productId       string
		serviceId       string
		logsRefreshRate time.Duration
		lastPageMarker  int64
		mockGetLogResponse model.Data
		mockGetConfigResponse *model.Config
		mockGetErr      error
		want            model.ConfigDisplayData
		wantErr         bool
	}{
		{
			name:            "Display Config response",
			nodeId:          "node-1",
			logFileName:     "one.log",
			mockGetConfigResponse: &model.Config{ Nodes : []string{"node1","node2"}, LogFileNames: []string{"log1","log2"}, RefreshRateMillis: 10000},
			want:            model.ConfigDisplayData{Nodes : []string{"node1","node2"}, Logs: []string{"log1","log2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Data{
				productId:       tt.productId,
				serviceId:        tt.serviceId,
				logsRefreshRate: time.Second,
			}
			realServiceLayer := newServiceLayer
			newServiceLayer = func(productId string) (servicelayer.ServiceLayer, error) {
				return &mockServiceLayer{
					t:                  t,
					expectNodeId:       tt.nodeId,
					expectLogFileName:  tt.logFileName,
					getLogResponse:     tt.mockGetLogResponse,
					getConfigResponse:  tt.mockGetConfigResponse,
					getErr:             tt.mockGetErr,
				},nil
			}
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := s.DisplayConfig(context.Background())

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout
			newServiceLayer = realServiceLayer
			if (err != nil) != tt.wantErr {
				t.Errorf("DisplayConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var actualOut model.ConfigDisplayData
			err = json.Unmarshal(out, &actualOut)

			if !reflect.DeepEqual(actualOut.Nodes, tt.want.Nodes ) {
				t.Errorf("DisplayConfig() got = %v, want %v", actualOut.Nodes, tt.want.Nodes)
			}
			if !reflect.DeepEqual(actualOut.Logs, tt.want.Logs ) {
				t.Errorf("DisplayConfig() got = %v, want %v", actualOut.Logs, tt.want.Logs)
			}
		})
	}
}

func Test_LiveLogs_ConfigNonInteractive(t *testing.T) {
	tests := []struct {
		name            string
		logFileName     string
		nodeId          string
		productId       string
		serviceId       string
		logsRefreshRate time.Duration
		lastPageMarker  int64
		mockGetLogResponse model.Data
		mockGetConfigResponse *model.Config
		mockGetErr      error
		want            model.ConfigDisplayData
		wantErr         bool
	}{
		{
			name:            "ConfigNonInteractive no product id",
			nodeId:          "node1",
			wantErr:          true,
			logFileName:     "log1",
		},
		{
			name:            "ConfigNonInteractive response",
			productId:        constants.ArtifactoryId,
			nodeId:          "node-1",
			logFileName:     "one.log",
			mockGetConfigResponse: &model.Config{ Nodes : []string{"node1","node2"}, LogFileNames: []string{"log1","log2"}, RefreshRateMillis: 10000},
			want:            model.ConfigDisplayData{Nodes : []string{"node1","node2"}, Logs: []string{"log1","log2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Data{
				serviceLayerClient: &mockServiceLayer{
					t:                  t,
					expectNodeId:       tt.nodeId,
					expectLogFileName:  tt.logFileName,
					getLogResponse:     tt.mockGetLogResponse,
					getConfigResponse:  tt.mockGetConfigResponse,
					getErr:             tt.mockGetErr,
				},
				productId:       tt.productId,
				serviceId:        tt.serviceId,
				logsRefreshRate: time.Second,
			}
			realServiceLayer := newServiceLayer
			newServiceLayer = func(productId string) (servicelayer.ServiceLayer, error) {
				return &mockServiceLayer{
					t:                  t,
					expectNodeId:       tt.nodeId,
					expectLogFileName:  tt.logFileName,
					getLogResponse:     tt.mockGetLogResponse,
					getConfigResponse:  tt.mockGetConfigResponse,
					getErr:             tt.mockGetErr,
				},nil
			}
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := s.ConfigNonInteractive(context.Background(),tt.productId, tt.serviceId)

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout
			newServiceLayer = realServiceLayer
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigNonInteractive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var actualOut model.ConfigDisplayData
			err = json.Unmarshal(out, &actualOut)

			if !reflect.DeepEqual(actualOut.Nodes, tt.want.Nodes ) {
				t.Errorf("ConfigNonInteractive() got = %v, want %v", actualOut.Nodes, tt.want.Nodes)
			}
			if !reflect.DeepEqual(actualOut.Logs, tt.want.Logs ) {
				t.Errorf("ConfigNonInteractive() got = %v, want %v", actualOut.Logs, tt.want.Logs)
			}
		})
	}
}

func Test_LiveLogs_TailLog(t *testing.T) {
	tests := []struct {
		name            string
		nodeId          string
		logFileName     string
		productId       string
		lastPageMarker  int64
		mockGetLogResponse model.Data
		mockGetConfigResponse *model.Config
		mockGetErr      error
		want            string
		wantErr         bool
	}{
		{
			name:    "missing node id",
			wantErr: true,
		},
		{
			name:    "missing log file name",
			nodeId:  "node-1",
			wantErr: true,
		},
		{
			name:        "error response",
			nodeId:      "node-1",
			logFileName: "one.log",
			mockGetErr:  fmt.Errorf("some-error"),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Data{
				serviceLayerClient: &mockServiceLayer{
					t:                  t,
					expectNodeId:       tt.nodeId,
					expectLogFileName:  tt.logFileName,
					getLogResponse:     tt.mockGetLogResponse,
					getConfigResponse:  tt.mockGetConfigResponse,
					getErr:             tt.mockGetErr,
				},
				productId:       tt.productId,
				logsRefreshRate: time.Second,
			}
			out := &bytes.Buffer{}
			err := s.tailLog(context.Background(), out)
			if (err != nil) != tt.wantErr {
				t.Errorf("tailLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(out.String(), tt.want) {
				t.Errorf("tailLog() got = %v, want %v", out.String(), tt.want)
			}
		})
	}
}

type mockServiceLayer struct {
	t                  *testing.T
	getLogResponse     model.Data
	getConfigResponse  *model.Config
	getErr             error
	expectNodeId       string
	expectLogFileName  string
	expectLastPageMarker int64
	expectLogsRefreshRate time.Duration
	logFileName        string
	lastPageMarker     int64
}

func (s *mockServiceLayer) GetLogData (_ context.Context,serviceId string) (logData model.Data, err error) {
	return s.getLogResponse, s.getErr
}
func (s *mockServiceLayer) GetConfig(ctx context.Context, serviceId string) (*model.Config, error) {
	return s.getConfigResponse, s.getErr
}
func (s *mockServiceLayer) GetConfigData(ctx context.Context, productId, serviceId string) (*model.Config, error) {
	return s.getConfigResponse, s.getErr
}
func (s *mockServiceLayer) GetNodeId () string {
	return s.expectNodeId
}
func (s *mockServiceLayer) SetNodeId (nodeId string) () {
	s.expectNodeId=nodeId
}
func (s *mockServiceLayer) GetLogFileName () string {
	return s.expectLogFileName
}
func (s *mockServiceLayer) SetLastPageMarker (lastPageMarker int64) {
	s.lastPageMarker=lastPageMarker
}
func (s *mockServiceLayer) GetLastPageMarker () (lastPageMarker int64) {
	return s.lastPageMarker
}
func (s *mockServiceLayer) SetLogFileName (logFileName string) {
	s.logFileName=logFileName
}
func (s *mockServiceLayer) SetLogsRefreshRate (logRefreshRate time.Duration) {
	s.expectLogsRefreshRate=logRefreshRate
}

func (s *mockServiceLayer) GetLogsRefreshRate () (logRefreshRate time.Duration) {
	return s.expectLogsRefreshRate
}

func (s *mockServiceLayer) GetServiceLayer () *mockServiceLayer{
	return s
}

func (s *mockServiceLayer) SetServiceLayer (productId string) error {
	return nil
}
