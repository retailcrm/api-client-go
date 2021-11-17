package retailcrm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectResponse(t *testing.T) {
	assert.Equal(t, ConnectResponse{
		SuccessfulResponse: SuccessfulResponse{
			Success: true,
		},
		AccountURL: "https://example.com",
	}, NewConnectResponse("https://example.com"))
}
