package v5

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
)

// Customer get method
func (c *Client) Customer(id, by, site string) (*CustomerResponse, int, error) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerFilter{context, site}
	params, _ := query.Values(fw)
	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/%s?%s", versionedPrefix, id, params.Encode()))
	if err != nil {
		return &resp, status, err
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

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerCreate method
func (c *Client) CustomerCreate(customer Customer, site ...string) (*CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse
	customerJson, _ := json.Marshal(&customer)

	p := url.Values{
		"customer": {string(customerJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerEdit method
func (c *Client) CustomerEdit(customer Customer, by string, site ...string) (*CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse
	var uid = strconv.Itoa(customer.Id)
	var context = checkBy(by)

	if context == "externalId" {
		uid = customer.ExternalId
	}

	customerJson, _ := json.Marshal(&customer)

	p := url.Values{
		"by":       {string(context)},
		"customer": {string(customerJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersUpload method
func (c *Client) CustomersUpload(customers []Customer, site ...string) (*CustomersUploadResponse, int, error) {
	var resp CustomersUploadResponse

	uploadJson, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(uploadJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/upload", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersFixExternalIds method
func (c *Client) CustomersFixExternalIds(customers []CustomerIdentifiers) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	customersJson, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(customersJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/fix-external-ids", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersHistory method
func (c *Client) CustomersHistory(filter CustomersHistoryFilter, limit, page int) (*CustomersHistoryResponse, int, error) {
	var resp CustomersHistoryResponse

	if limit == 0 {
		limit = 20
	}

	if page == 0 {
		page = 1
	}

	fw := CustomersHistoryParameters{filter, limit, page}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/history?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}
