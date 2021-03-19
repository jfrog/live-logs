package commands

import (
	"context"
	"fmt"
	"github.com/jfrog/live-logs/internal/model"
	"github.com/jfrog/live-logs/internal/servicelayer"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)


func Test_terminal_ConfigInteractive(t *testing.T) {
	tests := []struct {
		name            string
		mockConfigResponse string
		want            string
		wantErr         bool
		mockGetErr      string
		selectedKey     string
	}{
		{
			name:            "ConfigInteractive",
			mockConfigResponse: "['nodes': {'node1','node2'}, 'logs' : {'log1','log2'}]",
			want:             "['nodes': {'node1','node2'}, 'logs' : {'log1','log2'}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockLiveLog{
					mockConfigResponse: tt.mockConfigResponse,
			}

			var origPromptForAnyKey=PromptForAnyKey
			var origPromptSelectMenu=PromptSelectMenu

			PromptForAnyKey = func(promptPrefix string)  {
				return
			}
			PromptSelectMenu = func(selectionHeader string, selectionLabel string, values []string) (string, error) {
				if tt.wantErr {
					return tt.selectedKey, fmt.Errorf("%s",tt.mockGetErr)
				} else {
					return tt.selectedKey, nil
				}
			}
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ConfigInteractive(context.Background(),s)

			w.Close()
			out, _ := ioutil.ReadAll(r)

			os.Stdout = rescueStdout
			PromptForAnyKey=origPromptForAnyKey
			PromptSelectMenu=origPromptSelectMenu

			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigInteractive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(strings.Trim(string(out),"\n"), tt.want) {
				t.Errorf("ConfigInteractive() got = %v, want %v", strings.Trim(string(out),"\n"), tt.want)
			}
		})
	}
}

func Test_terminal_LogInteractive(t *testing.T) {
	tests := []struct {
		name            string
		nodeId          string
		logName         string
		mockGetErr      string
		mockExpectErr   string
		want            string
		wantErr         bool
		selectedKey     string
	}{
		{
			name:            "log Interactive",
			nodeId:          "node1",
			logName:          "artifactory-service.log",
			want:             "nodeId is node1, logName is artifactory-service.log",
		},
		{
			name:            "log Interactive error",
			nodeId:          "node1",
			logName:          "artifactory-service.log",
			want:             "nodeId is node1, logName is artifactory-service.log",
			wantErr:          true,
			mockGetErr:       "some error",
			mockExpectErr:    "some error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockLiveLog{
				LogName: tt.logName,
				nodeId: tt.nodeId,
			}

			var origPromptForAnyKey=PromptForAnyKey
			var origPromptSelectMenu=PromptSelectMenu

			PromptForAnyKey = func(promptPrefix string)  {
				return
			}
			PromptSelectMenu = func(selectionHeader string, selectionLabel string, values []string) (string, error) {
				if tt.wantErr {
					return tt.selectedKey, fmt.Errorf("%s",tt.mockGetErr)
				} else {
					if strings.Contains(selectionHeader,"Node") {
						return tt.nodeId,nil
					} else {
						return tt.logName,nil
					}
				}
			}
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := LogInteractiveMenu(context.Background(), true, s)

			w.Close()
			out, _ := ioutil.ReadAll(r)

			os.Stdout = rescueStdout
			PromptForAnyKey=origPromptForAnyKey
			PromptSelectMenu=origPromptSelectMenu

			if tt.wantErr && !strings.Contains(err.Error(),tt.mockExpectErr) {
				t.Errorf("LogInteractiveMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(strings.Trim(string(out),"\n"), tt.want) {
				t.Errorf("LogInteractiveMenu() got = %v, want %v", strings.Trim(string(out),"\n"), tt.want)
			}
		})
	}
}

func Test_terminal_selectCliServerId(t *testing.T) {
	tests := []struct {
		name            string
		want            string
		wantErr         bool
		mockGetErr      string
		selectedKey     string
	}{
		{
			name:            "selectCliServerId",
			selectedKey:      "local-artifactory",
			want:             "local-artifactory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var origPromptSelectMenu=PromptSelectMenu
			var origCliServerIds=CliServerIds

			CliServerIds = func() []string {
				return []string{tt.selectedKey}
			}
			PromptSelectMenu = func(selectionHeader string, selectionLabel string, values []string) (string, error) {
				if tt.wantErr {
					return tt.selectedKey, fmt.Errorf("%s",tt.mockGetErr)
				} else {
					return tt.selectedKey, nil
				}
			}

			serverId,err := selectCliServerId()

			PromptSelectMenu=origPromptSelectMenu
			CliServerIds=origCliServerIds

			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigInteractive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(serverId, tt.want) {
				t.Errorf("ConfigInteractive() got = %v, want %v", serverId, tt.want)
			}
		})
	}
}

func Test_terminal_selectLogDetails(t *testing.T) {
	tests := []struct {
		name            string
		want            string
		wantErr         bool
		nodeId          string
		mockGetErr      string
		logName         string
		selectedKey     string
	}{
		{
			name:            "log Interactive",
			nodeId:          "node1",
			logName:          "artifactory-service.log",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockLiveLog{}
			var origPromptSelectMenu=PromptSelectMenu

			PromptSelectMenu = func(selectionHeader string, selectionLabel string, values []string) (string, error) {
				if tt.wantErr {
					return tt.selectedKey, fmt.Errorf("%s",tt.mockGetErr)
				} else {
					if strings.Contains(selectionHeader,"Node") {
						return tt.nodeId,nil
					} else {
						return tt.logName,nil
					}
				}
			}

			nodeId,logName,_,err := selectLogDetails(context.Background(), s)

			PromptSelectMenu=origPromptSelectMenu

			if (err != nil) != tt.wantErr {
				t.Errorf("selectLogDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(nodeId, tt.nodeId) {
				t.Errorf("selectLogDetails() got = %v, want %v", nodeId, tt.nodeId)
			}
			if !reflect.DeepEqual(logName, tt.logName) {
				t.Errorf("selectLogDetails() got = %v, want %v", logName, tt.logName)
			}
		})
	}
}

func Test_terminal_selectProductId(t *testing.T) {
	tests := []struct {
		name            string
		productId       string
		serviceId       string
		logsRefreshRate time.Duration
		lastPageMarker  int64
		mockConfigResponse string
		mockGetErr      error
		want            string
		wantErr         bool
		selectedKey     string
	}{
		{
			name:            "selectProductId",
			selectedKey:      "rt",
			want:             "rt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var origPromptSelectMenu=PromptSelectMenu

			PromptSelectMenu = func(selectionHeader string, selectionLabel string, values []string) (string, error) {
				if tt.wantErr {
					return tt.selectedKey, fmt.Errorf("%s",tt.mockGetErr)
				} else {
					return tt.selectedKey, nil
				}
			}

			productId,err := selectProductId()

			PromptSelectMenu=origPromptSelectMenu

			if (err != nil) != tt.wantErr {
				t.Errorf("selectProductId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(productId, tt.want) {
				t.Errorf("selectProductId() got = %v, want %v", productId, tt.want)
			}
		})
	}
}

type mockLiveLog struct {
	productId       string
	serviceId       string
	logsRefreshRate time.Duration
	mockLogResponse string
	mockConfigResponse string
	LogName         string
	nodeId          string
}

func (s *mockLiveLog) SetProductId(productId string) {
	s.productId = productId
}

func (s *mockLiveLog) SetServiceLayer(productId string) error {
	return nil
}

func (s *mockLiveLog) GetServiceLayer() servicelayer.ServiceLayer {
	return nil
}

func (s *mockLiveLog) SetServiceId(serviceId string) {
	s.serviceId = serviceId
}

func (s *mockLiveLog) GetLogsRefreshRate() (logsRefreshRate time.Duration) {
	return s.logsRefreshRate
}

func (s *mockLiveLog) GetProductId()  string {
	return s.productId
}

func (s *mockLiveLog) GetServiceId() string {
	return s.serviceId
}

func (s *mockLiveLog) SetLogsRefreshRate(logsRefreshRate time.Duration) {
	s.logsRefreshRate = logsRefreshRate
}

func (s *mockLiveLog) CatLog(ctx context.Context, output io.Writer) error {
	return nil
}

func (s *mockLiveLog) tailLog(ctx context.Context, output io.Writer) error {
	return nil
}

func (s *mockLiveLog) doCatLog(ctx context.Context) (logReader io.Reader, err error) {
	return nil, nil
}

func (s *mockLiveLog) LogNonInteractive(ctx context.Context, cliProductId, cliServerId, nodeId, logName string, isStreaming bool) error {
	return nil
}

func (s *mockLiveLog) GetConfigData (ctx context.Context, productId, serviceId string) (srvConfig *model.Config, err error) {
	return &model.Config{RefreshRateMillis: 100,LogFileNames: []string{s.LogName}}, nil
}

func (s *mockLiveLog) PrintLogs (ctx context.Context, nodeId, logName  string, isStreaming bool) error {
	fmt.Printf("nodeId is %s, logName is %s", nodeId,logName)
	return nil
}

func (s *mockLiveLog) DisplayConfig(ctx context.Context)  error {
	fmt.Println(s.mockConfigResponse)
	return nil
}

func (s *mockLiveLog)  ConfigNonInteractive(ctx context.Context, cliProductId, cliServerId string) error {
	return nil
}
