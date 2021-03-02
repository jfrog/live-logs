package servicelayer

import (
	"context"
	"encoding/json"
	"fmt"
	cliCommands "github.com/jfrog/jfrog-cli-core/common/commands"
	"github.com/jfrog/live-logs/internal/clientlayer"
	"github.com/jfrog/live-logs/internal/constants"
	"github.com/jfrog/live-logs/internal/model"
	"time"
)
const (
	defaultRequestTimeout    = 15 * time.Second
	defaultLogRequestTimeout = time.Minute
	defaultLogsRefreshRate   = time.Second
)

type ArtifactoryData struct {
	ServiceData
}

func (s *ArtifactoryData) GetConfig(ctx context.Context, serverId string) (*model.Config, error) {

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelTimeout()
	baseUrl, err := s.getUrl(serverId)
	if err != nil {
		return nil, err
	}
	resBody, err := clientlayer.SendGet(timeoutCtx, serverId, constants.ConfigEndpoint,constants.EmptyNodeId, baseUrl, nil)

	if err != nil {
		return nil, err
	}

	logConfig := model.Config{}
	err = json.Unmarshal(resBody, &logConfig)
	if err != nil {
		return nil, err
	}
	if len(logConfig.LogFileNames) == 0 {
		return nil, fmt.Errorf("no log file names were found")
	}
	if len(logConfig.Nodes) == 0 {
		return nil, fmt.Errorf("no node names were found")
	}
	return &logConfig, nil
}

func (s *ArtifactoryData) GetLogData(ctx context.Context, serverId string) (logData model.Data, err error) {
	if s.nodeId == "" {
		return logData, fmt.Errorf("node id must be set")
	}
	if s.logFileName == "" {
		return logData, fmt.Errorf("log file name must be set")
	}

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultLogRequestTimeout)
	defer cancelTimeout()

	var endpoint string
	endpoint = fmt.Sprintf("%s?file_size=%d&id=%s", constants.DataEndpoint, s.lastPageMarker, s.logFileName)
	baseUrl, err := s.getUrl(serverId)
	if err != nil {
		return logData, err
	}

	resBody, err := clientlayer.SendGet(timeoutCtx, serverId, endpoint, s.nodeId, baseUrl, nil)

	if err != nil {
		return logData, err
	}

	if err := json.Unmarshal(resBody, &logData); err != nil {
		return logData, err
	}

	return logData, nil
}

func (s *ArtifactoryData) getUrl(serverId string)(url string,_ error){
	confDetails, err := cliCommands.GetConfig(serverId, false)
	if err != nil {
		return "", err
	}
	url = confDetails.GetArtifactoryUrl()
	if url == "" {
		return "", fmt.Errorf("no url found in serverId : %s", serverId)
	}
	return url, nil
}

func (s *ArtifactoryData) SetNodeId(nodeId string) {
	s.nodeId = nodeId
}

func (s *ArtifactoryData) SetLogFileName(logFileName string) {
	s.logFileName = logFileName
}

func (s *ArtifactoryData) SetLogsRefreshRate(logsRefreshRate time.Duration) {
	s.logsRefreshRate = logsRefreshRate
}

func (s *ArtifactoryData) SetLastPageMarker(pageMarker int64) {
	s.lastPageMarker = pageMarker
}

func (s *ArtifactoryData) GetLastPageMarker() int64 {
	return s.lastPageMarker
}

func (s *ArtifactoryData) GetNodeId() string {
	return s.nodeId
}

func (s *ArtifactoryData) GetLogFileName() string {
	return s.logFileName
}

func (s *ArtifactoryData) GetLogsRefreshRate() time.Duration {
	return s.logsRefreshRate
}
