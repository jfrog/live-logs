package clientlayer

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	cliCommands "github.com/jfrog/jfrog-cli-core/common/commands"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/live-logs/internal/constants"
)

func newPlatformHttpClient(cliServerId string) (*platformHttpClient, error) {
	platformDetails, err := cliCommands.GetConfig(cliServerId, false)
	if err != nil {
		return nil, err
	}
	platform, err := utils.CreateServiceManager(platformDetails, false)
	if err != nil {
		return nil, err
	}

	return &platformHttpClient{
		platform: platform,
	},nil
}

type platformHttpClient struct {
	platform artifactory.ArtifactoryServicesManager
}

func SendGet(_ context.Context, cliServerId, endpoint, nodeId, baseUrl string, extraHeaders map[string]string) ([]byte, error) {

	platformClient, err := newPlatformHttpClient(cliServerId)
	if err != nil {
		return nil, err
	}
	client := platformClient.platform.Client()
	httpClientDetails := (*client.JfrogServiceDetails).CreateHttpClientDetails()

	if nodeId != constants.EmptyNodeId && nodeId != "" {
		httpClientDetails.Headers[constants.NodeIdHeader] = nodeId
	}

	for key, value := range extraHeaders {
		httpClientDetails.Headers[key] = value
	}

	res, resBody, _, err := client.SendGet(baseUrl+endpoint, true, &httpClientDetails)

	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected response; status code: %d, message: %s", res.StatusCode, resBody)
	}
	return resBody, nil
}

