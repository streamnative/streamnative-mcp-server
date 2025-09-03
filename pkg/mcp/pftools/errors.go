package pftools

import (
	"errors"
	"strings"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/rest"
)

var (
	ErrFunctionNotFound = errors.New("function not found")
	ErrNotOurMessage    = errors.New("not our message")
)

// IsClusterUnhealthy checks if an error indicates cluster health issues
func IsClusterUnhealthy(err error) bool {
	if restErr, ok := err.(rest.Error); ok {
		return restErr.Code == 503 && strings.Contains(restErr.Reason, "no healthy upstream")
	}
	return false
}
