package servicelayer

import (
	"context"
	"fmt"
	"github.com/jfrog/live-logs/internal/constants"
	"github.com/jfrog/live-logs/internal/model"
	"github.com/jfrog/live-logs/internal/util"
	"time"
)

type ServiceLayer interface {
	// Queries and returns the livelog configuration from the remote service, based on the set node id.
	GetConfig(ctx context.Context, serverId string) (*model.Config, error)

	// Queries and returns the livelog data from the remote service, based on the set node id and log file name.
	GetLogData(ctx context.Context, serverId string) (model.Data, error)

	// Sets the node id to use when querying the remote service for log data.
	SetNodeId(nodeId string)
	GetNodeId() string

	// Sets the log file name to use when querying the remote service for log data.
	SetLogFileName(logFileName string)
	GetLogFileName() string

	// Sets the refresh rate interval between each log request.
	SetLogsRefreshRate(logsRefreshRate time.Duration)
	GetLogsRefreshRate() time.Duration

	// Sets the file size between each log request.
	SetLastPageMarker(pageMarker int64)
	GetLastPageMarker() int64
}

func NewService(productId string) (serviceLayer ServiceLayer, err error) {
	if productId == "" {
		return nil, fmt.Errorf("service id must be set")
	}

	switch productId {
		case constants.ArtifactoryId :
			serviceLayer = new (ArtifactoryData)

		case constants.McId :
			serviceLayer = new (McData)

		case constants.PipelinesId :
			serviceLayer = new (PipelinesData)

		case constants.DistributionId :
			serviceLayer = new (DistributionData)

		case constants.XrayId :
			serviceLayer = new (XrayData)

		default :
			err = fmt.Errorf("invalid product id '%s' provided, valid values are %v", productId, util.FetchAllProductIds())
	}
	return serviceLayer, err
}

func errorHandle(statusCode int, resBody []byte) error{
	if statusCode == 200 {
		return nil
	}
	if statusCode == 404 || statusCode == 400 || statusCode == 429 {
		return fmt.Errorf("status code: %d; message: %s", statusCode, resBody)
	}
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("unexpected response; status code: %d, message: %s", statusCode, resBody)
	}
	return nil
}
