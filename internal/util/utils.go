package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/jfrog/live-logs/internal/constants"
	"strings"
	"time"
)

func InSlice(values []string, wantedVal string) bool {
	for _, val := range values {
		if val == wantedVal {
			return true
		}
	}
	return false
}

func SliceToCsv(values []string) string {
	var buf bytes.Buffer
	wr := csv.NewWriter(&buf)
	err := wr.Write(values)
	if err != nil {
		return ""
	}

	wr.Flush()
	return strings.TrimSuffix(buf.String(), "\n")
}

func MillisToDuration(timeInMillis int64) time.Duration {
	return time.Duration(timeInMillis) * time.Millisecond
}

func ValidateArgument(argumentName string, wantedVal string, allValues []string) error {
	values := allValues
	if len(values) == 0 {
		return fmt.Errorf("no %v found", argumentName)
	}
	if !InSlice(values, wantedVal) {
		return fmt.Errorf("%v not found [%v], consider using one of the following %v values [%v]", argumentName, wantedVal, argumentName, SliceToCsv(values))
	}
	return nil
}

func FetchAllProductIds() []string {
    productIds := []string{constants.ArtifactoryId, constants.XrayId, constants.McId, constants.PipelinesId, constants.DistributionId}
	return productIds
}
