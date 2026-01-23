package pftools

import (
	"errors"
	"testing"

	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/rest"
)

func TestIsClusterUnhealthy(t *testing.T) {
	t.Run("no healthy upstream", func(t *testing.T) {
		err := rest.Error{Code: 503, Reason: "no healthy upstream"}
		if !IsClusterUnhealthy(err) {
			t.Fatalf("expected cluster unhealthy for 503 no healthy upstream")
		}
	})

	t.Run("network error text", func(t *testing.T) {
		err := errors.New("read tcp 127.0.0.1:1234->127.0.0.1:443: read: connection reset by peer")
		if !IsClusterUnhealthy(err) {
			t.Fatalf("expected cluster unhealthy for network error text")
		}
	})

	t.Run("non cluster error", func(t *testing.T) {
		err := rest.Error{Code: 500, Reason: "internal error"}
		if IsClusterUnhealthy(err) {
			t.Fatalf("did not expect cluster unhealthy for 500 internal error")
		}
	})
}

func TestIsAuthError(t *testing.T) {
	t.Run("rest auth error", func(t *testing.T) {
		err := rest.Error{Code: 401, Reason: "unauthorized"}
		if !IsAuthError(err) {
			t.Fatalf("expected auth error for 401")
		}
	})

	t.Run("text auth error", func(t *testing.T) {
		err := errors.New("token expired")
		if !IsAuthError(err) {
			t.Fatalf("expected auth error for token expired")
		}
	})
}

func TestIsNotFoundError(t *testing.T) {
	t.Run("rest not found", func(t *testing.T) {
		err := rest.Error{Code: 404, Reason: "Not Found"}
		if !IsNotFoundError(err) {
			t.Fatalf("expected not found for 404")
		}
	})

	t.Run("text not found", func(t *testing.T) {
		err := errors.New("code: 404 reason: 404 Not Found")
		if !IsNotFoundError(err) {
			t.Fatalf("expected not found for 404 text")
		}
	})
}
