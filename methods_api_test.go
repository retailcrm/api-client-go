package retailcrm

import (
	"testing"
	"net/http"
)

func TestVersions(t *testing.T) {
	c := UnversionedClient()

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

func TestCredentials(t *testing.T) {
	c := UnversionedClient()

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
