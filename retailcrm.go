package retailcrm

import (
	"github.com/retailcrm/api-client-go/v5"
)

// Version5 API client for v5
func Version5(url string, key string) *v5.Client {
	var client = v5.New(url, key)

	return client
}
