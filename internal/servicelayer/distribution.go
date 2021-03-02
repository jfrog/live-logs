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

type DistributionData struct {
	ServiceData
}

func (s *DistributionData) GetConfig(ctx context.Context, serverId string) (*model.Config, error) {

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelTimeout()

	baseUrl, headers, err := s.getConnectionDetails(serverId)
	if err != nil {
		return nil, err
	}

	resBody, err := clientlayer.SendGet(timeoutCtx, serverId, constants.ConfigEndpoint,constants.EmptyNodeId,baseUrl,headers)

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

func (s *DistributionData) GetLogData(ctx context.Context, serverId string) (logData model.Data, err error) {
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
	baseUrl, headers, err := s.getConnectionDetails(serverId)
	if err != nil {
		return logData, err
	}
	resBody, err := clientlayer.SendGet(timeoutCtx, serverId, endpoint, s.nodeId,baseUrl,headers)

	if err != nil {
		return logData, err
	}

	if err := json.Unmarshal(resBody, &logData); err != nil {
		return logData, err
	}

	return logData, nil
}

func (s *DistributionData) getConnectionDetails(serverId string)(url string, headers map[string]string,_ error){
	confDetails, err := cliCommands.GetConfig(serverId, false)
	if err != nil {
		return "",nil, err
	}
	url = confDetails.GetDistributionUrl()
	accessToken := confDetails.GetAccessToken()
	if url == "" {
		return "",nil, fmt.Errorf("no url found in serverId : %s",serverId)
	}
	if accessToken == "" {
		return "",nil, fmt.Errorf("no access token found in serverId : %s, this is mandatory to connect to Distribution product",serverId)
	}

	headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + accessToken

	return url,headers, nil
}

func (s *DistributionData) SetNodeId(nodeId string) {
	s.nodeId = nodeId
}

func (s *DistributionData) SetLogFileName(logFileName string) {
	s.logFileName = logFileName
}

func (s *DistributionData) SetLogsRefreshRate(logsRefreshRate time.Duration) {
	s.logsRefreshRate = logsRefreshRate
}

func (s *DistributionData) SetLastPageMarker(pageMarker int64) {
	s.lastPageMarker = pageMarker
}

func (s *DistributionData) GetLastPageMarker() int64 {
	return s.lastPageMarker
}

func (s *DistributionData) GetNodeId() string {
	return s.nodeId
}

func (s *DistributionData) GetLogFileName() string {
	return s.logFileName
}

func (s *DistributionData) GetLogsRefreshRate() time.Duration {
	return s.logsRefreshRate
}

