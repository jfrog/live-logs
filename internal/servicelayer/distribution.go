package servicelayer

import (
	"context"
	"encoding/json"
	"fmt"
	cliCommands "github.com/jfrog/jfrog-cli-core/common/commands"
	cliVersionHelper "github.com/jfrog/jfrog-client-go/utils/version"
	"github.com/jfrog/live-logs/internal/clientlayer"
	"github.com/jfrog/live-logs/internal/constants"
	"github.com/jfrog/live-logs/internal/model"
	"os"
	"strings"
	"time"
)

const (
	distributionVersionEndPoint = "api/v1/system/info"
	distributionMinVersionSupport = "2.7.0"
	distributionConfigEndpoint = "api/v1/system/logs/config"
	distributionDataEndpoint   = "api/v1/system/logs/data"
)

type DistributionData struct {
	nodeId          string
	logFileName     string
	lastPageMarker  int64
	logsRefreshRate time.Duration
}

type distributionVersionData struct {
	Version string `json:"version,omitempty"`
}

func (s *DistributionData) GetConfig(ctx context.Context, serverId string) (*model.Config, error) {

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelTimeout()

	baseUrl, headers, err := s.getConnectionDetails(serverId)
	if err != nil {
		return nil, err
	}

	res, resBody, err := clientlayer.SendGet(timeoutCtx, serverId, distributionConfigEndpoint,constants.EmptyNodeId,baseUrl,headers)
	if err != nil {
		return nil, err
	}

	err = errorHandle(res.StatusCode, resBody)
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
	endpoint = fmt.Sprintf("%s?file_size=%d&id=%s", distributionDataEndpoint, s.lastPageMarker, s.logFileName)
	baseUrl, headers, err := s.getConnectionDetails(serverId)
	if err != nil {
		return logData, err
	}
	res, resBody, err := clientlayer.SendGet(timeoutCtx, serverId, endpoint, s.nodeId,baseUrl,headers)
	if err != nil {
		return logData, err
	}

	err = errorHandle(res.StatusCode, resBody)
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
		return "",nil, fmt.Errorf("the Distribution url was not found in the serverId : %s; verify that you are using the latest version of the JFrog CLI",serverId)
	}
	if accessToken == "" {
		return "",nil, fmt.Errorf("no access token found in the serverId : %s; the tokens mandatory for connecting to Distribution",serverId)
	}

	headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + accessToken

	return url,headers, nil
}

func (s *DistributionData) getVersion(ctx context.Context, serverId string) (string, error) {
	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancelTimeout()

	baseUrl, headers, err := s.getConnectionDetails(serverId)
	if err != nil {
		return "", err
	}
	res, resBody, err := clientlayer.SendGet(timeoutCtx, serverId, distributionVersionEndPoint,constants.EmptyNodeId, baseUrl, headers)
	if err != nil {
		return "", err
	}

	err = errorHandle(res.StatusCode, resBody)
	if err != nil {
		return "", err
	}

	versionData := distributionVersionData{}
	err = json.Unmarshal(resBody, &versionData)
	if err != nil {
		return "", err
	}
	if versionData.Version == "" {
		return "", fmt.Errorf("could not retreive version information from Distribution")
	}
	return strings.TrimSpace(versionData.Version), nil
}

func (s *DistributionData) checkVersion(ctx context.Context, serverId string) error {
	if os.Getenv(constants.VersionCheckEnv) == "false" {
		return nil
	}

	currentVersion, err := s.getVersion(ctx, serverId)
	if err != nil {
		return err
	}
	versionHelper := cliVersionHelper.NewVersion(currentVersion)

	if versionHelper.Compare(distributionMinVersionSupport) < 0 {
		return fmt.Errorf("found distribution version as %s, minimum supported version is %s", currentVersion, distributionMinVersionSupport)
	}
	return nil
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

