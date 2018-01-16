package retailcrm

import (
	"fmt"
	"net/http"
	"errors"
	"encoding/json"
	"github.com/google/go-querystring/query"
)

// Customer get method
func (c *Client) Customer(id, by, site string) (*CustomerResponse, int, error) {
	var resp CustomerResponse
	var context = CheckBy(by)

	fw := CustomerGetFilter{context, site}
	params, _ := query.Values(fw)
	data, status, err := c.getRequest(fmt.Sprintf("/customers/%s?%s", id, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	if status >= http.StatusBadRequest {
		return &resp, status, errors.New(fmt.Sprintf("HTTP request error. Status code: %d.\n", status))
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Customers list method
func (c *Client) Customers(filter CustomersFilter, limit, page int) (*CustomersResponse, int, error) {
	var resp CustomersResponse

	if limit == 0 {
		limit = 20
	}

	if page == 0 {
		page = 1
	}

	fw := CustomersParameters{filter, limit, page}
	params, _ := query.Values(fw)

	data, status, err := c.getRequest(fmt.Sprintf("/customers?%s", params.Encode()))

	if err != nil {
		return &resp, status, err
	}

	if status >= http.StatusBadRequest {
		return &resp, status, errors.New(fmt.Sprintf("HTTP request error. Status code: %d.\n", status))
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}
