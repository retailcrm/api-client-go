package retailcrm

import (
	"net/http"
	"net/url"
	"testing"
)

const (
	TestUrl        = "https://demo.retailcrm.ru"
	TestApiKey     = "111"
	WrongApiKeyMsg = "Wrong \"apiKey\" value."
)

func client() *Client {
	return New(TestUrl, TestApiKey)
}

func TestGetRequest(t *testing.T) {
	c := client()

	data, status, _ := c.GetRequest("/store/products")
	if status != http.StatusForbidden {
		t.Fail()
	}

	resp, _ := c.ErrorResponse(data)
	if resp.ErrorMsg != WrongApiKeyMsg {
		t.Fail()
	}
}

func TestPostRequest(t *testing.T) {
	c := client()

	data, status, _ := c.PostRequest("/orders/create", url.Values{})
	if status != http.StatusForbidden {
		t.Fail()
	}

	resp, _ := c.ErrorResponse(data)
	if resp.ErrorMsg != WrongApiKeyMsg {
		t.Fail()
	}
}
