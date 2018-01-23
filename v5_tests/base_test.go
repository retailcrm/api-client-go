package v5_tests

import (
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/retailcrm/api-client-go/v5"
)

type Configuration struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

func buildConfiguration() *Configuration {
	return &Configuration{
		os.Getenv("RETAILCRM_URL"),
		os.Getenv("RETAILCRM_KEY"),
	}
}

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

func client() *v5.Client {
	configuration := buildConfiguration()
	return v5.New(configuration.Url, configuration.Key)
}

func TestGetRequest(t *testing.T) {
	c := client()
	_, status, _ := c.GetRequest("/fake-method")

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestPostRequest(t *testing.T) {
	c := client()
	_, status, _ := c.PostRequest("/fake-method", url.Values{})

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestClient_ApiVersionsVersions(t *testing.T) {
	c := client()
	data, status, err := c.ApiVersions()
	if err != nil {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Fail()
	}

	if data.Success != true {
		t.Fail()
	}
}

func TestClient_ApiCredentialsCredentials(t *testing.T) {
	c := client()
	data, status, err := c.ApiCredentials()
	if err != nil {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Fail()
	}

	if data.Success != true {
		t.Fail()
	}
}
