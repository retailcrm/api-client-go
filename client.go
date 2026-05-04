package retailcrm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// New initialize client.
func New(url string, key string) *Client {
	return &Client{
		URL:        strings.TrimRight(url, "/"),
		Key:        key,
		httpClient: &http.Client{Timeout: time.Minute},
	}
}

// WithLogger sets the provided logger instance into the Client.
func (c *Client) WithLogger(logger BasicLogger) *Client {
	c.logger = logger
	return c
}

// WithHTTPClient sets the provided HTTP client instance into the Client.
func (c *Client) WithHTTPClient(client *http.Client) *Client {
	c.httpClient = client
	return c
}

// WithContext sets context.Context which, if canceled, will terminate all active API calls.
func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

// EnableRateLimiter activates rate limiting with specified retry attempts.
func (c *Client) EnableRateLimiter(maxAttempts uint) *Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.limiter = NewSingleKeyLimiter()
	c.maxAttempts = maxAttempts

	return c
}

// EnableCustomRateLimiter activates rate limiting with custom limiter logic and specified retry attempts.
func (c *Client) EnableCustomRateLimiter(limiter Limiter, maxAttempts uint) *Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if limiter == nil {
		c.limiter = nil
		c.maxAttempts = 0

		return c
	}

	c.limiter = limiter
	c.maxAttempts = maxAttempts

	return c
}

// applyRateLimit applies rate limiting before sending a request.
func (c *Client) applyRateLimit(uri string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.limiter == nil {
		return nil
	}

	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	return c.limiter.Limit(ctx, uri, c.Key)
}

// triggerResponseAwareLimiter sends *http.Response to the limiter for further reading (e.g. headers).
func (c *Client) triggerResponseAwareLimiter(resp *http.Response) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.limiter == nil || resp == nil {
		return
	}

	if val, ok := c.limiter.(ResponseAware); ok {
		val.ProcessResponse(resp)
	}
}

func (c *Client) executeWithRetryBytes(
	uri string,
	executeFunc func() (interface{}, *http.Response, int, error),
) ([]byte, int, error) {
	res, status, err := c.executeWithRetry(uri, executeFunc)
	if res == nil {
		return nil, status, err
	}
	return res.([]byte), status, err
}

func (c *Client) executeWithRetryReadCloser(
	uri string,
	executeFunc func() (interface{}, *http.Response, int, error),
) (io.ReadCloser, int, error) {
	res, status, err := c.executeWithRetry(uri, executeFunc)
	if res == nil {
		return nil, status, err
	}
	return res.(io.ReadCloser), status, err
}

// executeWithRetry executes a request with retry logic for rate limiting.
func (c *Client) executeWithRetry(
	uri string,
	executeFunc func() (interface{}, *http.Response, int, error),
) (interface{}, int, error) {
	if c.limiter == nil {
		resp, _, st, err := executeFunc()
		return resp, st, err
	}

	var (
		res           interface{}
		httpResp      *http.Response
		statusCode    int
		err           error
		lastAttempt   bool
		attempt       uint = 1
		maxAttempts        = c.maxAttempts
		totalAttempts      = "∞"
		infinite           = maxAttempts == 0
	)

	var baseDelay time.Duration
	if strings.HasPrefix(uri, "/telephony") {
		baseDelay = telephonyDelay
	} else {
		baseDelay = regularDelay
	}

	if !infinite {
		totalAttempts = strconv.FormatUint(uint64(maxAttempts), 10)
	}

	for infinite || attempt <= maxAttempts {
		if err := c.applyRateLimit(uri); err != nil {
			return nil, 0, err
		}

		res, httpResp, statusCode, err = executeFunc()
		c.triggerResponseAwareLimiter(httpResp)
		lastAttempt = !infinite && attempt == maxAttempts
		isRateLimited := statusCode == http.StatusServiceUnavailable || statusCode == http.StatusTooManyRequests

		// If rate limited on final attempt, set error to ErrRateLimited. Return results otherwise.
		if isRateLimited && lastAttempt {
			return res, statusCode, ErrRateLimited
		}

		// If not rate limited or on final attempt, return result.
		if !isRateLimited || lastAttempt {
			return res, statusCode, err
		}

		// Calculate exponential backoff delay: baseDelay * 2^(attempt-1).
		backoffDelay := baseDelay * (1 << (attempt - 1))
		if c.Debug {
			c.writeLog("API Error: rate limited (503), retrying in %v (attempt %d/%s)",
				backoffDelay, attempt, totalAttempts)
		}

		time.Sleep(backoffDelay)
		attempt++
	}

	return res, statusCode, err
}

// writeLog writes to the log.
func (c *Client) writeLog(format string, v ...interface{}) {
	if c.logger != nil {
		c.logger.Printf(format, v...)
		return
	}

	log.Printf(format, v...)
}

// GetRequest implements GET Request.
func (c *Client) GetRequest(urlWithParameters string, versioned ...bool) ([]byte, int, error) {
	var prefix = "/api/v5"

	if len(versioned) > 0 {
		if !versioned[0] {
			prefix = "/api"
		}
	}

	uri := urlWithParameters

	return c.executeWithRetryBytes(uri, func() (interface{}, *http.Response, int, error) {
		var res []byte

		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), nil)
		if err != nil {
			return res, nil, 0, err
		}

		req.Header.Set("X-API-KEY", c.Key)

		if c.Debug {
			c.writeLog("API Request: %s %s", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), c.Key)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return res, resp, 0, err
		}

		if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != http.StatusServiceUnavailable {
			return res, resp, resp.StatusCode, CreateGenericAPIError(
				fmt.Sprintf("HTTP request error. Status code: %d.", resp.StatusCode))
		}

		res, err = buildRawResponse(resp)
		if err != nil {
			return res, resp, 0, err
		}

		if resp.StatusCode >= http.StatusBadRequest &&
			resp.StatusCode < http.StatusInternalServerError &&
			resp.StatusCode != http.StatusServiceUnavailable {
			return res, resp, resp.StatusCode, CreateAPIError(res)
		}

		if c.Debug {
			c.writeLog("API Response: %s", res)
		}

		return res, resp, resp.StatusCode, nil
	})
}

// PostRequest implements POST Request with generic body data.
func (c *Client) PostRequest(
	uri string,
	postData interface{},
	contType ...string,
) ([]byte, int, error) {
	var contentType string

	if len(contType) > 0 {
		contentType = contType[0]
	} else {
		contentType = "application/x-www-form-urlencoded"
	}

	prefix := "/api/v5"

	return c.executeWithRetryBytes(uri, func() (interface{}, *http.Response, int, error) {
		var res []byte

		reader, err := getReaderForPostData(postData)
		if err != nil {
			return res, nil, 0, err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", c.URL, prefix, uri), reader)
		if err != nil {
			return res, nil, 0, err
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("X-API-KEY", c.Key)

		if c.Debug {
			c.writeLog("API Request: %s %s", uri, c.Key)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return res, resp, 0, err
		}

		if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != http.StatusServiceUnavailable {
			return res, resp, resp.StatusCode, CreateGenericAPIError(
				fmt.Sprintf("HTTP request error. Status code: %d.", resp.StatusCode))
		}

		res, err = buildRawResponse(resp)
		if err != nil {
			return res, resp, 0, err
		}

		if resp.StatusCode >= http.StatusBadRequest &&
			resp.StatusCode < http.StatusInternalServerError &&
			resp.StatusCode != http.StatusServiceUnavailable {
			return res, resp, resp.StatusCode, CreateAPIError(res)
		}

		if c.Debug {
			c.writeLog("API Response: %s", res)
		}

		return res, resp, resp.StatusCode, nil
	})
}

func getReaderForPostData(postData interface{}) (io.Reader, error) {
	var reader io.Reader

	switch d := postData.(type) {
	case url.Values:
		reader = strings.NewReader(d.Encode())
	default:
		if i, ok := d.(io.Reader); ok {
			reader = i
		} else {
			err := errors.New("postData should be url.Values or implement io.Reader")
			return nil, err
		}
	}

	return reader, nil
}

func buildRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}

// checkBy select identifier type.
func checkBy(by string) string {
	var context = ByID

	if by != ByID {
		context = ByExternalID
	}

	return context
}

// fillSite add site code to parameters if present.
func fillSite(p *url.Values, site []string) {
	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}
}

func marshalToString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// APIVersions get all available API versions for exact account
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-api-versions
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.APIVersions()
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.versions {
//		log.Printf("%v\n", value)
//	}
func (c *Client) APIVersions() (VersionResponse, int, error) {
	var resp VersionResponse

	data, status, err := c.GetRequest("/api-versions", false)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// APICredentials get all available API methods for exact account
//
// For more information see https://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-credentials
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.APICredentials()
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Scopes {
//		log.Printf("%v\n", value)
//	}
func (c *Client) APICredentials() (CredentialResponse, int, error) {
	var resp CredentialResponse

	data, status, err := c.GetRequest("/credentials", false)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// APISystemInfo get info about system
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-system-info
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.APISystemInfo()
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	log.Printf("%v\n", data)
func (c *Client) APISystemInfo() (SystemInfoResponse, int, error) {
	var resp SystemInfoResponse

	data, status, err := c.GetRequest("/system-info", false)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Customers returns list of customers matched the specified filter
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-customers
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Customers(retailcrm.CustomersRequest{
//		Filter: CustomersFilter{
//			City: "Moscow",
//		},
//		Page: 3,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Customers {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Customers(parameters CustomersRequest) (CustomersResponse, int, error) {
	var resp CustomersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomersCombine combine given customers
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-combine
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomersCombine([]retailcrm.Customer{{ID: 1}, {ID: 2}}, Customer{ID: 3})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CustomersCombine(customers []Customer, resultCustomer Customer) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	combineJSONIn, err := marshalToString(&customers)
	if err != nil {
		return resp, 0, err
	}

	combineJSONOut, err := marshalToString(&resultCustomer)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customers":      {combineJSONIn},
		"resultCustomer": {combineJSONOut},
	}

	data, status, err := c.PostRequest("/customers/combine", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomerCreate creates customer
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomersCombine(retailcrm.Customer{
//		FirstName:  "Ivan",
//		LastName:   "Ivanov",
//		Patronymic: "Ivanovich",
//		ExternalID: 1,
//		Email:      "ivanov@example.com",
//		Address: &retailcrm.Address{
//			City:     "Moscow",
//			Street:   "Kutuzovsky",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v", data.ID)
//	}
func (c *Client) CustomerCreate(customer Customer, site ...string) (CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse

	customerJSON, err := marshalToString(&customer)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customer": {customerJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomersFixExternalIds customers external ID
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-fix-external-ids
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomersFixExternalIds([]retailcrm.IdentifiersPair{{
//		ID:         1,
//		ExternalID: 12,
//	}})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CustomersFixExternalIds(customers []IdentifiersPair) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	customersJSON, err := marshalToString(&customers)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customers": {customersJSON},
	}

	data, status, err := c.PostRequest("/customers/fix-external-ids", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomersHistory returns customer's history
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-customers-history
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomersHistory(retailcrm.CustomersHistoryRequest{
//		Filter: retailcrm.CustomersHistoryFilter{
//			SinceID: 20,
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.History {
//		log.Printf("%v\n", value)
//	}
func (c *Client) CustomersHistory(parameters CustomersHistoryRequest) (CustomersHistoryResponse, int, error) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/history?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomerNotes returns customer related notes
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-customers-notes
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomerNotes(retailcrm.NotesRequest{
//		Filter: retailcrm.NotesFilter{
//			CustomerIds: []int{1,2,3}
//		},
//		Page: 1,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Notes {
//		log.Printf("%v\n", value)
//	}
func (c *Client) CustomerNotes(parameters NotesRequest) (NotesResponse, int, error) {
	var resp NotesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/notes?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomerNoteCreate note creation
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-notes-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomerNoteCreate(retailcrm.Note{
//		Text:      "some text",
//		ManagerID: 12,
//		Customer: &retailcrm.Customer{
//			ID: 1,
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) CustomerNoteCreate(note Note, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse

	noteJSON, err := marshalToString(&note)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"note": {noteJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/notes/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomerNoteDelete remove customer related note
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-notes-id-delete
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomerNoteDelete(12)
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CustomerNoteDelete(id int) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {strconv.Itoa(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/customers/notes/%d/delete", id), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomersUpload customers batch upload
//
// # This method can return response together with error if http status is equal 460
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomersUpload([]retailcrm.Customer{
//		{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			ExternalID: 1,
//			Email:      "ivanov@example.com",
//		},
//		{
//			FirstName:  "Petr",
//			LastName:   "Petrov",
//			Patronymic: "Petrovich",
//			ExternalID: 2,
//			Email:      "petrov@example.com",
//		},
//	}}
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.UploadedCustomers)
//	}
func (c *Client) CustomersUpload(customers []Customer, site ...string) (CustomersUploadResponse, int, error) {
	var resp CustomersUploadResponse

	uploadJSON, err := marshalToString(&customers)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customers": {uploadJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/upload", p)
	if err != nil && status != HTTPStatusUnknown {
		return resp, status, err
	}

	errJSON := json.Unmarshal(data, &resp)
	if errJSON != nil {
		return resp, status, errJSON
	}

	if status == HTTPStatusUnknown {
		return resp, status, err
	}

	return resp, status, nil
}

// Customer returns information about customer
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-customers-externalId
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Customer(12, retailcrm.ByExternalID, "")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Customer)
//	}
func (c *Client) Customer(customerID, by, site string) (CustomerResponse, int, error) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/%s?%s", customerID, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomerEdit edit exact customer
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-customers-externalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomerEdit(
//		retailcrm.Customer{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			ID: 		1,
//			Email:      "ivanov@example.com",
//		},
//		retailcrm.ByID,
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Customer)
//	}
func (c *Client) CustomerEdit(customer Customer, by string, site ...string) (CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse
	var uid = strconv.Itoa(customer.ID)
	var context = checkBy(by)

	if context == ByExternalID {
		uid = customer.ExternalID
	}

	customerJSON, err := marshalToString(&customer)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":       {context},
		"customer": {customerJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers/%s/edit", uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomerSubscriptions changes customer's mailing subscriptions
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-customers-externalId-subscriptions
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomerSubscriptions(
//		"customer-external-id",
//		retailcrm.ByExternalID,
//		[]retailcrm.CustomerSubscriptionEdit{
//			{
//				Channel:      "email",
//				Subscription: "news",
//				Active:       false,
//				MessageID:    123,
//			},
//			{
//				Channel: "sms",
//				Active:  true,
//			},
//		},
//		"site",
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Success)
//	}
func (c *Client) CustomerSubscriptions(
	customerID, by string,
	subscriptions []CustomerSubscriptionEdit,
	site ...string,
) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse
	var context = checkBy(by)

	subscriptionsJSON, err := marshalToString(&subscriptions)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":            {context},
		"subscriptions": {subscriptionsJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers/%s/subscriptions", customerID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomers returns list of corporate customers matched the specified filter
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-customers-corporate
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomers(retailcrm.CorporateCustomersRequest{
//		Filter: CorporateCustomersFilter{
//			City: "Moscow",
//		},
//		Page: 3,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.CustomersCorporate {
//		log.Printf("%v\n", value)
//	}
func (c *Client) CorporateCustomers(parameters CorporateCustomersRequest) (CorporateCustomersResponse, int, error) {
	var resp CorporateCustomersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerCreate creates corporate customer
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomerCreate(retailcrm.CorporateCustomer{
//		Nickname:  "Company",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v", data.ID)
//	}
func (c *Client) CorporateCustomerCreate(customer CorporateCustomer, site ...string) (
	CorporateCustomerChangeResponse, int, error,
) {
	var resp CorporateCustomerChangeResponse

	customerJSON, err := marshalToString(&customer)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customerCorporate": {customerJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers-corporate/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomersFixExternalIds will fix corporate customers external ID's
//
// For more information see  http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-fix-external-ids
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomersFixExternalIds([]retailcrm.IdentifiersPair{{
//		ID:         1,
//		ExternalID: 12,
//	}})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CorporateCustomersFixExternalIds(customers []IdentifiersPair) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	customersJSON, err := marshalToString(&customers)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customersCorporate": {customersJSON},
	}

	data, status, err := c.PostRequest("/customers-corporate/fix-external-ids", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomersHistory returns corporate customer's history
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-fix-external-ids
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomersHistory(retailcrm.CorporateCustomersHistoryRequest{
//		Filter: retailcrm.CorporateCustomersHistoryFilter{
//			SinceID: 20,
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.History {
//		log.Printf("%v\n", value)
//	}
func (c *Client) CorporateCustomersHistory(parameters CorporateCustomersHistoryRequest) (
	CorporateCustomersHistoryResponse, int, error,
) {
	var resp CorporateCustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate/history?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomersNotes returns list of corporate customers notes matched the specified filter
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-customers-corporate-notes
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomersNotes(retailcrm.CorporateCustomersNotesRequest{
//		Filter: CorporateCustomersNotesFilter{
//			Text: "text",
//		},
//		Page: 3,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Notes {
//		log.Printf("%v\n", value)
//	}
func (c *Client) CorporateCustomersNotes(parameters CorporateCustomersNotesRequest) (
	CorporateCustomersNotesResponse, int, error,
) {
	var resp CorporateCustomersNotesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate/notes?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerNoteCreate creates corporate customer note
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-notes-create
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//		data, status, err := client.CorporateCustomerNoteCreate(retailcrm.CorporateCustomerNote{
//			Text:  "text",
//	     Customer: &retailcrm.IdentifiersPair{
//				ID: 1,
//			}
//		})
//
//		if err != nil {
//			if apiErr, ok := retailcrm.AsAPIError(err); ok {
//				log.Fatalf("http status: %d, %s", status, apiErr.String())
//			}
//
//			log.Fatalf("http status: %d, error: %s", status, err)
//		}
//
//		if data.Success == true {
//			fmt.Printf("%v", data.ID)
//		}
func (c *Client) CorporateCustomerNoteCreate(note CorporateCustomerNote, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse

	noteJSON, err := marshalToString(&note)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"note": {noteJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers-corporate/notes/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerNoteDelete removes note from corporate customer
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-notes-id-delete
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomerNoteDelete(12)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CorporateCustomerNoteDelete(id int) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {strconv.Itoa(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/notes/%d/delete", id), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomersUpload corporate customers batch upload
//
// # This method can return response together with error if http status is equal 460
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomersUpload([]retailcrm.CorporateCustomer{
//		{
//			Nickname:  "Company",
//			ExternalID: 1,
//		},
//		{
//			Nickname:  "Company 2",
//			ExternalID: 2,
//		},
//	}}
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.UploadedCustomers)
//	}
func (c *Client) CorporateCustomersUpload(
	customers []CorporateCustomer, site ...string,
) (CorporateCustomersUploadResponse, int, error) {
	var resp CorporateCustomersUploadResponse

	uploadJSON, err := marshalToString(&customers)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customersCorporate": {uploadJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers-corporate/upload", p)
	if err != nil && status != HTTPStatusUnknown {
		return resp, status, err
	}

	errJSON := json.Unmarshal(data, &resp)
	if errJSON != nil {
		return resp, status, err
	}

	if status == HTTPStatusUnknown {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomer returns information about corporate customer
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-customers-corporate-externalId
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomer(12, retailcrm.ByExternalID, "")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.CorporateCustomer)
//	}
func (c *Client) CorporateCustomer(id, by, site string) (CorporateCustomerResponse, int, error) {
	var resp CorporateCustomerResponse
	var context = checkBy(by)

	fw := CustomerRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate/%s?%s", id, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerAddresses returns information about corporate customer addresses
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-customers-corporate-externalId-addresses
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomerAddresses("ext-id", retailcrm.CorporateCustomerAddressesRequest{
//		Filter: v5,CorporateCustomerAddressesFilter{
//			Name: "Main Address",
//		},
//		By:    retailcrm.ByExternalID,
//		Site:  "site",
//		Limit: 20,
//		Page:  1,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Addresses)
//	}
func (c *Client) CorporateCustomerAddresses(
	id string, parameters CorporateCustomerAddressesRequest,
) (CorporateCustomersAddressesResponse, int, error) {
	var resp CorporateCustomersAddressesResponse

	parameters.By = checkBy(parameters.By)
	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate/%s/addresses?%s", id, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerAddressesCreate creates corporate customer address
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-addresses-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := c.CorporateCustomerAddressesCreate("ext-id", retailcrm.ByExternalID, retailcrm.CorporateCustomerAddress{
//		Text: "this is new address",
//		Name: "New Address",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v", data.ID)
//	}
func (c *Client) CorporateCustomerAddressesCreate(
	customerID string, customerBy string, address CorporateCustomerAddress, site ...string,
) (CreateResponse, int, error) {
	var resp CreateResponse

	addressJSON, err := marshalToString(&address)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"address": {addressJSON},
		"by":      {checkBy(customerBy)},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/addresses/create", customerID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerAddressesEdit edit exact corporate customer address
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-addresses-entityExternalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := c.CorporateCustomerAddressesEdit(
//		"customer-ext-id",
//		retailcrm.ByExternalID,
//		retailcrm.ByExternalID,
//		CorporateCustomerAddress{
//			ExternalID: "addr-ext-id",
//			Name:       "Main Address 2",
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Customer)
//	}
func (c *Client) CorporateCustomerAddressesEdit(
	customerID, customerBy, entityBy string, address CorporateCustomerAddress, site ...string,
) (CreateResponse, int, error) {
	var (
		resp CreateResponse
		uid  string
	)

	customerBy = checkBy(customerBy)
	entityBy = checkBy(entityBy)

	switch entityBy {
	case ByID:
		uid = strconv.Itoa(address.ID)
	case ByExternalID:
		uid = address.ExternalID
	}

	addressJSON, err := marshalToString(&address)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":       {customerBy},
		"entityBy": {entityBy},
		"address":  {addressJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/addresses/%s/edit", customerID, uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerCompanies returns information about corporate customer companies
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-customers-corporate-externalId-companies
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomerCompanies("ext-id", retailcrm.IdentifiersPairRequest{
//		Filter: v5,IdentifiersPairFilter{
//			Ids: []string{"1"},
//		},
//		By:    retailcrm.ByExternalID,
//		Site:  "site",
//		Limit: 20,
//		Page:  1,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Companies)
//	}
func (c *Client) CorporateCustomerCompanies(
	id string, parameters IdentifiersPairRequest,
) (CorporateCustomerCompaniesResponse, int, error) {
	var resp CorporateCustomerCompaniesResponse

	parameters.By = checkBy(parameters.By)
	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate/%s/companies?%s", id, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerCompaniesCreate creates corporate customer company
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-companies-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := c.CorporateCustomerCompaniesCreate("ext-id", retailcrm.ByExternalID, retailcrm.Company{
//		Name: "Company name",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v", data.ID)
//	}
func (c *Client) CorporateCustomerCompaniesCreate(
	customerID string, customerBy string, company Company, site ...string,
) (CreateResponse, int, error) {
	var resp CreateResponse

	companyJSON, err := marshalToString(&company)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"company": {companyJSON},
		"by":      {checkBy(customerBy)},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/companies/create", customerID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerCompaniesEdit edit exact corporate customer company
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-companies-entityExternalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := c.CorporateCustomerCompaniesEdit(
//		"customer-ext-id",
//		retailcrm.ByExternalID,
//		retailcrm.ByExternalID,
//		Company{
//			ExternalID: "company-ext-id",
//			Name:       "New Company Name",
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) CorporateCustomerCompaniesEdit(
	customerID, customerBy, entityBy string, company Company, site ...string,
) (CreateResponse, int, error) {
	var (
		resp CreateResponse
		uid  string
	)

	customerBy = checkBy(customerBy)
	entityBy = checkBy(entityBy)

	switch entityBy {
	case ByID:
		uid = strconv.Itoa(company.ID)
	case ByExternalID:
		uid = company.ExternalID
	}

	addressJSON, err := marshalToString(&company)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":       {customerBy},
		"entityBy": {entityBy},
		"company":  {addressJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/companies/%s/edit", customerID, uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerContacts returns information about corporate customer contacts
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-customers-corporate-externalId-contacts
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomerContacts("ext-id", retailcrm.IdentifiersPairRequest{
//		Filter: retailcrm.IdentifiersPairFilter{
//			Ids: []string{"1"},
//		},
//		By:    retailcrm.ByExternalID,
//		Site:  "site",
//		Limit: 20,
//		Page:  1,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Contacts)
//	}
func (c *Client) CorporateCustomerContacts(
	id string, parameters IdentifiersPairRequest,
) (CorporateCustomerContactsResponse, int, error) {
	var resp CorporateCustomerContactsResponse

	parameters.By = checkBy(parameters.By)
	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers-corporate/%s/contacts?%s", id, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerContactsCreate creates corporate customer contact
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-contacts-create
//
// Example (customer with specified id or externalId should exist in specified site):
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := c.CorporateCustomerContactsCreate("ext-id", retailcrm.ByExternalID, retailcrm.CorporateCustomerContact{
//		IsMain: false,
//		Customer: retailcrm.CorporateCustomerContactCustomer{
//			ExternalID: "external_id",
//			Site:       "site",
//		},
//		Companies: []IdentifiersPair{},
//	}, "site")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v", data.ID)
//	}
func (c *Client) CorporateCustomerContactsCreate(
	customerID string, customerBy string, contact CorporateCustomerContact, site ...string,
) (CreateResponse, int, error) {
	var resp CreateResponse

	companyJSON, err := marshalToString(&contact)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"contact": {companyJSON},
		"by":      {checkBy(customerBy)},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/contacts/create", customerID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerContactsEdit edit exact corporate customer contact
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-contacts-entityExternalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := c.CorporateCustomerContactsEdit("ext-id", retailcrm.ByExternalID, retailcrm.ByID, retailcrm.CorporateCustomerContact{
//		IsMain: false,
//		Customer: retailcrm.CorporateCustomerContactCustomer{
//			ID: 2350,
//		},
//	}, "site")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) CorporateCustomerContactsEdit(
	customerID, customerBy, entityBy string, contact CorporateCustomerContact, site ...string,
) (CreateResponse, int, error) {
	var (
		resp CreateResponse
		uid  string
	)

	customerBy = checkBy(customerBy)
	entityBy = checkBy(entityBy)

	switch entityBy {
	case ByID:
		uid = strconv.Itoa(contact.Customer.ID)
	case ByExternalID:
		uid = contact.Customer.ExternalID
	}

	addressJSON, err := marshalToString(&contact)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":       {customerBy},
		"entityBy": {entityBy},
		"contact":  {addressJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/contacts/%s/edit", customerID, uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CorporateCustomerEdit edit exact corporate customer
//
// For more information see http://help.retailcrm.pro/Developers/ApiVersion5#post--api-v5-customers-corporate-externalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CorporateCustomerEdit(
//		retailcrm.CorporateCustomer{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			ID: 		1,
//			Email:      "ivanov@example.com",
//		},
//		retailcrm.ByID,
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Customer)
//	}
func (c *Client) CorporateCustomerEdit(customer CorporateCustomer, by string, site ...string) (
	CustomerChangeResponse, int, error,
) {
	var resp CustomerChangeResponse
	var uid = strconv.Itoa(customer.ID)
	var context = checkBy(by)

	if context == ByExternalID {
		uid = customer.ExternalID
	}

	customerJSON, err := marshalToString(&customer)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":                {context},
		"customerCorporate": {customerJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers-corporate/%s/edit", uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ClearCart clears the current customer's shopping cart
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-customer-interaction-site-cart-clear
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.ClearCart("site_id", retailcrm.SiteFilter{SiteBy: "id"},
//		retailcrm.ClearCartRequest{
//			ClearedAt: time.Now().String(),
//			Customer: retailcrm.CartCustomer{
//				ID:         1,
//				ExternalID: "ext_id",
//				Site:       "site",
//				BrowserID:  "browser_id",
//				GaClientID: "ga_client_id",
//			},
//			Order: retailcrm.ClearCartOrder{
//				ID:         1,
//				ExternalID: "ext_id",
//				Number:     "abc123",
//			},
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) ClearCart(site string, filter SiteFilter, req ClearCartRequest) (
	SuccessfulResponse, int, error,
) {
	var resp SuccessfulResponse

	updateJSON, err := json.Marshal(&req)
	if err != nil {
		return SuccessfulResponse{}, 0, err
	}

	p := url.Values{
		"cart": {string(updateJSON)},
	}

	params, _ := query.Values(filter)

	data, status, err := c.PostRequest(fmt.Sprintf("/customer-interaction/%s/cart/clear?%s", site, params.Encode()), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// SetCart creates or overwrites shopping cart data
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-customer-interaction-site-cart-set
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.SetCart("site_id", retailcrm.SiteFilter{SiteBy: "id"},
//		retailcrm.SetCartRequest{
//			ExternalID: "ext_id",
//			DroppedAt:  time.Now().String(),
//			Link:       "link",
//			Customer: retailcrm.CartCustomer{
//				ID:         1,
//				ExternalID: "ext_id",
//				Site:       "site",
//				BrowserID:  "browser_id",
//				GaClientID: "ga_client_id",
//			},
//			Items: []retailcrm.SetCartItem{
//				{
//					Quantity: 1,
//					Price:    1.0,
//					Offer: retailcrm.SetCartOffer{
//						ID: 1,
//						ExternalID: "ext_id",
//						XMLID: "xml_id",
//					},
//				},
//			},
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) SetCart(site string, filter SiteFilter, req SetCartRequest) (
	SuccessfulResponse, int, error,
) {
	var resp SuccessfulResponse

	updateJSON, err := json.Marshal(&req)
	if err != nil {
		return SuccessfulResponse{}, 0, err
	}

	p := url.Values{
		"cart": {string(updateJSON)},
	}

	params, _ := query.Values(filter)

	data, status, err := c.PostRequest(fmt.Sprintf("/customer-interaction/%s/cart/set?%s", site, params.Encode()), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// GetCart returns the current customer's shopping cart
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-customer-interaction-site-cart-customerId
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.GetCart("site_id","customer_id",
//		retailcrm.GetCartFilter{ SiteBy: "code", By: "externalId"})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) GetCart(site, customer string, filter GetCartFilter) (CartResponse, int, error) {
	var resp CartResponse

	params, _ := query.Values(filter)

	data, status, err := c.GetRequest(fmt.Sprintf("/customer-interaction/%s/cart/%s?%s", site, customer, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// GetFavorites returns the current customer's list of favorite offers
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-customer-interaction-site-favorites-customerId
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.GetFavorites("site_id","customer_id",
//		retailcrm.GetFavoritesFilter{ SiteBy: "code", By: "externalId"})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) GetFavorites(site, customer string, filter FavoritesFilter) (FavoritesResponse, int, error) {
	var resp FavoritesResponse

	params, _ := query.Values(filter)

	data, status, err := c.GetRequest(fmt.Sprintf("/customer-interaction/%s/favorites/%s?%s", site, customer, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// AddFavorite adds favorite products for customer
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-customer-interaction-site-favorites-customerId-add
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.AddFavorite("site_id", "customer_external_id", retailcrm.FavoritesFilter{SiteBy: "id", By: "externalId"},
//		retailcrm.ChangeFavoritesRequest{
//			Offer: retailcrm.SerializedRelationOffer{
//			    ID: 1,
//			    ExternalID: "ext_id",
//			    XMLID: "xml_id",
//			},
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) AddFavorite(site, customer string, filter FavoritesFilter, req ChangeFavoritesRequest) (
	SuccessfulResponse, int, error,
) {
	var resp SuccessfulResponse

	updateJSON, err := json.Marshal(&req)
	if err != nil {
		return SuccessfulResponse{}, 0, err
	}

	p := url.Values{
		"favorite": {string(updateJSON)},
	}

	params, _ := query.Values(filter)

	data, status, err := c.PostRequest(fmt.Sprintf("/customer-interaction/%s/favorites/%s/add?%s", site, customer, params.Encode()), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// RemoveFavorite removes favorite products from customer
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-customer-interaction-site-favorites-customerId-remove
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.RemoveFavorite("site_id", "customer_external_id", retailcrm.FavoritesFilter{SiteBy: "id", By: "externalId"},
//		retailcrm.ChangeFavoritesRequest{
//			Offer: retailcrm.SerializedRelationOffer{
//			    ID: 1,
//			    ExternalID: "ext_id",
//			    XMLID: "xml_id",
//			},
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) RemoveFavorite(site, customer string, filter FavoritesFilter, req ChangeFavoritesRequest) (
	SuccessfulResponse, int, error,
) {
	var resp SuccessfulResponse

	updateJSON, err := json.Marshal(&req)
	if err != nil {
		return SuccessfulResponse{}, 0, err
	}

	p := url.Values{
		"favorite": {string(updateJSON)},
	}

	params, _ := query.Values(filter)

	data, status, err := c.PostRequest(fmt.Sprintf("/customer-interaction/%s/favorites/%s/remove?%s", site, customer, params.Encode()), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryTracking updates tracking data
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-delivery-generic-subcode-tracking
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//	 t, _ := time.Parse("2006-01-02 15:04:05", "2012-12-12 12:12:12")
//
//		data, status, err := client.DeliveryTracking(
//			[]retailcrm.DeliveryTrackingRequest{{
//				DeliveryID: "1",
//				TrackNumber: "123",
//				History: []retailcrm.DeliveryHistoryRecord{
//					{
//						Code: "cancel",
//						UpdatedAt: t.Format(time.RFC3339),
//					},
//				},
//			}},
//			"delivery-1",
//		)
//
//		if err != nil {
//			if apiErr, ok := retailcrm.AsAPIError(err); ok {
//				log.Fatalf("http status: %d, %s", status, apiErr.String())
//			}
//
//			log.Fatalf("http status: %d, error: %s", status, err)
//		}
func (c *Client) DeliveryTracking(parameters []DeliveryTrackingRequest, subcode string) (
	SuccessfulResponse, int, error,
) {
	var resp SuccessfulResponse

	updateJSON, err := marshalToString(&parameters)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"statusUpdate": {updateJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/delivery/generic/%s/tracking", subcode), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryShipments returns list of shipments
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-delivery-shipments
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryShipments(retailcrm.DeliveryShipmentsRequest{
//		Limit: 12,
//		Filter: retailcrm.ShipmentFilter{
//			DateFrom: "2012-12-12",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.DeliveryShipments {
//		log.Printf("%v\n", value)
//	}
func (c *Client) DeliveryShipments(parameters DeliveryShipmentsRequest) (DeliveryShipmentsResponse, int, error) {
	var resp DeliveryShipmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/delivery/shipments?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryCalculate calculates delivery cost
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIMethods#post--api-v5-delivery-calculate
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryCalculate(retailcrm.DeliveryCalculateRequest{
//		DeliveryTypeCodes: []string{"courier"},
//		Order: retailcrm.DeliveryCalculateOrder{
//			Weight: 1200,
//			Delivery: &retailcrm.DeliveryCalculateSerializedDelivery{
//				Address: &retailcrm.Address{
//					City:     "Москва",
//					Building: "1",
//				},
//			},
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Calculations {
//		log.Printf("%v\n", value)
//	}
func (c *Client) DeliveryCalculate(parameters DeliveryCalculateRequest) (DeliveryCalculateResponse, int, error) {
	var resp DeliveryCalculateResponse

	orderJSON, err := marshalToString(&parameters.Order)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"deliveryTypeCodes[]": parameters.DeliveryTypeCodes,
		"order":               {orderJSON},
	}

	data, status, err := c.PostRequest("/delivery/calculate", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryShipmentCreate creates shipment
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-delivery-shipments-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryShipmentCreate(
//		retailcrm.DeliveryShipment{
//			Date: "2012-12-12",
//			Time: retailcrm.DeliveryTime{
//				From: "18:00",
//				To: "20:00",
//			},
//			Orders: []retailcrm.Order{{Number: "12"}},
//		},
//		"sdek",
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) DeliveryShipmentCreate(
	shipment DeliveryShipment, deliveryType string, site ...string,
) (DeliveryShipmentUpdateResponse, int, error) {
	var resp DeliveryShipmentUpdateResponse
	updateJSON, err := marshalToString(&shipment)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"deliveryType":     {deliveryType},
		"deliveryShipment": {updateJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/delivery/shipments/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryShipment get information about shipment
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-delivery-shipments-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryShipment(12)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.DeliveryShipment)
//	}
func (c *Client) DeliveryShipment(id int) (DeliveryShipmentResponse, int, error) {
	var resp DeliveryShipmentResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/delivery/shipments/%d", id))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryShipmentEdit shipment editing
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-delivery-shipments-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryShipmentEdit(retailcrm.DeliveryShipment{
//		ID: "12",
//		Time: retailcrm.DeliveryTime{
//			From: "14:00",
//			To: "18:00",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) DeliveryShipmentEdit(shipment DeliveryShipment, site ...string) (
	DeliveryShipmentUpdateResponse, int, error,
) {
	var resp DeliveryShipmentUpdateResponse
	updateJSON, err := marshalToString(&shipment)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"deliveryShipment": {updateJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/delivery/shipments/%s/edit", strconv.Itoa(shipment.ID)), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// IntegrationModule returns integration module
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-integration-modules-code
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.IntegrationModule("moysklad3")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.IntegrationModule)
//	}
func (c *Client) IntegrationModule(code string) (IntegrationModuleResponse, int, error) {
	var resp IntegrationModuleResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/integration-modules/%s", code))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// LinksCreate creates a link
//
// For more information see https://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-links-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.LinksCreate(retailcrm.SerializedOrderLink{
//		Comment:  "comment for link",
//		Orders: []retailcrm.LinkedOrder{{ID: 10}, {ID: 12}},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Println("Creating a link")
//	}
func (c *Client) LinksCreate(link SerializedOrderLink, site ...string) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	linkJSON, err := json.Marshal(link)

	if err != nil {
		return resp, http.StatusBadRequest, err
	}

	p := url.Values{
		"link": {string(linkJSON)},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/links/create", p)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ClientIdsUpload uploading of web analytics clientId
//
// For more information see https://docs.simla.com/Developers/API/APIVersions/APIv5#post--api-v5-web-analytics-client-ids-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.ClientIdsUpload([]retailcrm.ClientID{
//		{
//			Value:    "value",
//			Order:    LinkedOrder{ID: 10, ExternalID: "externalID", Number: "number"},
//			Customer: SerializedEntityCustomer{ID: 10, ExternalID: "externalID"},
//			Site: "site",
//		},
//		{
//			Value:    "value",
//			Order:    LinkedOrder{ID: 12, ExternalID: "externalID2", Number: "number2"},
//			Customer: SerializedEntityCustomer{ID: 12, ExternalID: "externalID2"},
//			Site: "site2",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Println("Upload is successful")
//	}
func (c *Client) ClientIdsUpload(clientIds []ClientID) (ClientIDResponse, int, error) {
	var resp ClientIDResponse
	clientIdsJSON, err := json.Marshal(&clientIds)

	if err != nil {
		return resp, http.StatusBadRequest, err
	}

	p := url.Values{
		"clientIds": {string(clientIdsJSON)},
	}

	data, status, err := c.PostRequest("/web-analytics/client-ids/upload", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// SourcesUpload uploading of sources
//
// For more information see https://docs.simla.com/Developers/API/APIVersions/APIv5#post--api-v5-web-analytics-sources-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.SourcesUpload([]retailcrm.Source{
//		{
//			Source:   "source",
//			Medium:   "medium",
//			Campaign: "campaign",
//			Keyword:  "keyword",
//			Content:  "content",
//			ClientID: "10",
//			Order:    LinkedOrder{ID: 10, ExternalID: "externalId", Number: "number"},
//			Customer: SerializedEntityCustomer{ID: 10, ExternalID: "externalId"},
//			Site:     "site",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Println("Upload is successful!")
//	}
func (c *Client) SourcesUpload(sources []Source) (SourcesResponse, int, error) {
	var resp SourcesResponse
	sourcesJSON, err := json.Marshal(&sources)

	if err != nil {
		return resp, http.StatusBadRequest, err
	}

	p := url.Values{
		"sources": {string(sourcesJSON)},
	}

	data, status, err := c.PostRequest("/web-analytics/sources/upload", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Currencies returns a list of currencies
//
// For more information see https://docs.simla.com/Developers/API/APIVersions/APIv5#get--api-v5-reference-currencies
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Currencies()
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Currencies {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Currencies() (CurrencyResponse, int, error) {
	var resp CurrencyResponse

	data, status, err := c.GetRequest("/reference/currencies")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CurrenciesCreate create currency
//
// For more information see https://docs.simla.com/Developers/API/APIVersions/APIv5#post--api-v5-reference-currencies-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CurrenciesCreate(retailcrm.Currency{
//		Code:                    "RUB",
//		IsBase:                  true,
//		IsAutoConvert:           true,
//		AutoConvertExtraPercent: 1,
//		ManualConvertNominal:    1,
//		ManualConvertValue:      1,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Println("Create currency")
//	}
func (c *Client) CurrenciesCreate(currency Currency) (CurrencyCreateResponse, int, error) {
	var resp CurrencyCreateResponse
	currencyJSON, err := json.Marshal(&currency)

	if err != nil {
		return resp, http.StatusBadRequest, err
	}

	p := url.Values{
		"currency": {string(currencyJSON)},
	}

	data, status, err := c.PostRequest("/reference/currencies/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CurrenciesEdit edit an currency
//
// For more information see https://docs.simla.com/Developers/API/APIVersions/APIv5#post--api-v5-reference-currencies-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CurrenciesEdit(
//		retailcrm.Currency{
//			ID:                      10,
//			Code:                    "RUB",
//			IsBase:                  true,
//			IsAutoConvert:           true,
//			AutoConvertExtraPercent: 1,
//			ManualConvertNominal:    1,
//			ManualConvertValue:      1,
//		},
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//	if data.Success == true {
//		log.Println("Currency was edit")
//	}
func (c *Client) CurrenciesEdit(currency Currency) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse
	var uid = strconv.Itoa(currency.ID)

	currencyJSON, err := json.Marshal(&currency)

	if err != nil {
		return resp, http.StatusBadRequest, err
	}

	p := url.Values{
		"currency": {string(currencyJSON)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/currencies/%s/edit", uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// IntegrationModuleEdit integration module create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-integration-modules-code
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	name := "MS"
//	code := "moysklad3"
//
//	data, status, err := client.IntegrationModuleEdit(retailcrm.IntegrationModule{
//		Code:            code,
//		IntegrationCode: code,
//		Active:          false,
//		Name:            fmt.Sprintf("Integration module %s", name),
//		AccountURL:      fmt.Sprintf("http://example.com/%s/account", name),
//		BaseURL:         fmt.Sprintf("http://example.com/%s", name),
//		ClientID:        "123",
//		Logo:            "https://cdn.worldvectorlogo.com/logos/github-icon.svg",
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v\n", data.Info)
//	}
func (c *Client) IntegrationModuleEdit(integrationModule IntegrationModule) (
	IntegrationModuleEditResponse, int, error,
) {
	var resp IntegrationModuleEditResponse
	updateJSON, err := marshalToString(&integrationModule)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{"integrationModule": {updateJSON}}

	data, status, err := c.PostRequest(fmt.Sprintf("/integration-modules/%s/edit", integrationModule.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// UpdateScopes updates permissions for the API key
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-integration-modules-code-update-scopes
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.UpdateScopes("moysklad3", retailcrm.ScopesRequired{
//		Scopes: []string{"scope1", "scope2"},
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v\n", data.APIKey)
//	}
func (c *Client) UpdateScopes(code string, request ScopesRequired) (UpdateScopesResponse, int, error) {
	var resp UpdateScopesResponse
	updateJSON, err := marshalToString(&request)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"requires": {updateJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/integration-modules/%s/update-scopes", code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Orders returns list of orders matched the specified filters
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Orders(retailcrm.OrdersRequest{Filter: retailcrm.OrdersFilter{City: "Moscow"}, Page: 1})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Orders {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Orders(parameters OrdersRequest) (OrdersResponse, int, error) {
	var resp OrdersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrdersCombine combines given orders
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-combine
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrdersCombine("ours", retailcrm.Order{ID: 1}, retailcrm.Order{ID: 1})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrdersCombine(technique string, order, resultOrder Order) (OperationResponse, int, error) {
	var resp OperationResponse

	combineJSONIn, err := marshalToString(&order)
	if err != nil {
		return resp, 0, err
	}

	combineJSONOut, err := marshalToString(&resultOrder)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"technique":   {technique},
		"order":       {combineJSONIn},
		"resultOrder": {combineJSONOut},
	}

	data, status, err := c.PostRequest("/orders/combine", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderCreate creates an order
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderCreate(retailcrm.Order{
//		FirstName:  "Ivan",
//		LastName:   "Ivanov",
//		Patronymic: "Ivanovich",
//		Email:      "ivanov@example.com",
//		Items:      []retailcrm.OrderItem{{Offer: retailcrm.Offer{ID: 12}, Quantity: 5}},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) OrderCreate(order Order, site ...string) (OrderCreateResponse, int, error) {
	var resp OrderCreateResponse
	orderJSON, err := marshalToString(&order)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"order": {orderJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrdersFixExternalIds set order external ID
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-fix-external-ids
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrdersFixExternalIds(([]retailcrm.IdentifiersPair{{
//		ID:         1,
//		ExternalID: 12,
//	}})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v\n", data.ID)
//	}
func (c *Client) OrdersFixExternalIds(orders []IdentifiersPair) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	ordersJSON, err := marshalToString(&orders)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"orders": {ordersJSON},
	}

	data, status, err := c.PostRequest("/orders/fix-external-ids", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrdersHistory returns orders history
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders-history
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrdersHistory(retailcrm.OrdersHistoryRequest{Filter: retailcrm.OrdersHistoryFilter{SinceID: 20}})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.History {
//		log.Printf("%v\n", value)
//	}
func (c *Client) OrdersHistory(parameters OrdersHistoryRequest) (OrdersHistoryResponse, int, error) {
	var resp OrdersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/history?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderPaymentCreate creates payment
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-payments-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderPaymentCreate(retailcrm.Payment{
//		Order: &retailcrm.Order{
//			ID: 12,
//		},
//		Amount: 300,
//		Type:   "cash",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) OrderPaymentCreate(payment Payment, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse

	paymentJSON, err := marshalToString(&payment)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"payment": {paymentJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/payments/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderPaymentDelete payment removing
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-payments-id-delete
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderPaymentDelete(12)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrderPaymentDelete(id int) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {strconv.Itoa(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/payments/%d/delete", id), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderPaymentEdit edit payment
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-payments-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderPaymentEdit(
//		retailcrm.Payment{
//			ID:     12,
//			Amount: 500,
//		},
//		retailcrm.ByID,
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrderPaymentEdit(payment Payment, by string, site ...string) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse
	var uid = strconv.Itoa(payment.ID)
	var context = checkBy(by)

	if context == ByExternalID {
		uid = payment.ExternalID
	}

	paymentJSON, err := marshalToString(&payment)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":      {context},
		"payment": {paymentJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/payments/%s/edit", uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrdersStatuses returns orders statuses
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders-statuses
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrdersStatuses(retailcrm.OrdersStatusesRequest{
//		IDs:         []int{1},
//		ExternalIDs: []string{"2"},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrdersStatuses(request OrdersStatusesRequest) (OrdersStatusesResponse, int, error) {
	var resp OrdersStatusesResponse

	params, _ := query.Values(request)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/statuses?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrdersUpload batch orders uploading
//
// # This method can return response together with error if http status is equal 460
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrdersUpload([]retailcrm.Order{
//		{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			Email:      "ivanov@example.com",
//			Items:      []retailcrm.OrderItem{{Offer: retailcrm.Offer{ID: 12}, Quantity: 5}},
//		},
//		{
//			FirstName:  "Pert",
//			LastName:   "Petrov",
//			Patronymic: "Petrovich",
//			Email:      "petrov@example.com",
//			Items:      []retailcrm.OrderItem{{Offer: retailcrm.Offer{ID: 13}, Quantity: 1}},
//		}
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.UploadedOrders)
//	}
func (c *Client) OrdersUpload(orders []Order, site ...string) (OrdersUploadResponse, int, error) {
	var resp OrdersUploadResponse

	uploadJSON, err := marshalToString(&orders)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"orders": {uploadJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/upload", p)
	if err != nil && status != HTTPStatusUnknown {
		return resp, status, err
	}

	errJSON := json.Unmarshal(data, &resp)
	if errJSON != nil {
		return resp, status, err
	}

	if status == HTTPStatusUnknown {
		return resp, status, err
	}

	return resp, status, nil
}

// Order returns information about order
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders-externalId
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Order(12, retailcrm.ByExternalID, "")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Order)
//	}
func (c *Client) Order(id, by, site string) (OrderResponse, int, error) {
	var resp OrderResponse
	var context = checkBy(by)

	fw := OrderRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/%s?%s", id, params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderEdit edit an order
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-externalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderEdit(
//		retailcrm.Order{
//			ID:    12,
//			Items: []retailcrm.OrderItem{{Offer: retailcrm.Offer{ID: 13}, Quantity: 6}},
//		},
//		retailcrm.ByID,
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrderEdit(order Order, by string, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse
	var uid = strconv.Itoa(order.ID)
	var context = checkBy(by)

	if context == ByExternalID {
		uid = order.ExternalID
	}

	orderJSON, err := marshalToString(&order)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"by":    {context},
		"order": {orderJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/%s/edit", uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Packs returns list of packs matched the specified filters
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders-packs
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Packs(retailcrm.PacksRequest{Filter: retailcrm.PacksFilter{OrderID: 12}})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Packs {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Packs(parameters PacksRequest) (PacksResponse, int, error) {
	var resp PacksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PackCreate creates a pack
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-packs-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.PackCreate(Pack{
//		Store:    "store-1",
//		ItemID:   12,
//		Quantity: 1,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) PackCreate(pack Pack) (CreateResponse, int, error) {
	var resp CreateResponse
	packJSON, err := marshalToString(&pack)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"pack": {packJSON},
	}

	data, status, err := c.PostRequest("/orders/packs/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PacksHistory returns a history of order packing
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders-packs-history
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.PacksHistory(retailcrm.PacksHistoryRequest{Filter: retailcrm.OrdersHistoryFilter{SinceID: 5}})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.History {
//		log.Printf("%v\n", value)
//	}
func (c *Client) PacksHistory(parameters PacksHistoryRequest) (PacksHistoryResponse, int, error) {
	var resp PacksHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs/history?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Pack returns a pack info
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-orders-packs-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Pack(112)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Pack)
//	}
func (c *Client) Pack(id int) (PackResponse, int, error) {
	var resp PackResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs/%d", id))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PackDelete removes a pack
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-packs-id-delete
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.PackDelete(112)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) PackDelete(id int) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/packs/%d/delete", id), url.Values{})
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PackEdit edit a pack
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-orders-packs-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.PackEdit(Pack{ID: 12, Quantity: 2})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) PackEdit(pack Pack) (CreateResponse, int, error) {
	var resp CreateResponse

	packJSON, err := marshalToString(&pack)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"pack": {packJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/packs/%d/edit", pack.ID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Countries returns list of available country codes
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-countries
func (c *Client) Countries() (CountriesResponse, int, error) {
	var resp CountriesResponse

	data, status, err := c.GetRequest("/reference/countries")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostGroups returns costs groups list
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-cost-groups
func (c *Client) CostGroups() (CostGroupsResponse, int, error) {
	var resp CostGroupsResponse

	data, status, err := c.GetRequest("/reference/cost-groups")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostGroupEdit edit costs groups
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-cost-groups-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostGroupEdit(retailcrm.CostGroup{
//		Code:   "group-1",
//		Color:  "#da5c98",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CostGroupEdit(costGroup CostGroup) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&costGroup)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"costGroup": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/cost-groups/%s/edit", costGroup.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostItems returns costs items list
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-cost-items
func (c *Client) CostItems() (CostItemsResponse, int, error) {
	var resp CostItemsResponse

	data, status, err := c.GetRequest("/reference/cost-items")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostItemEdit edit costs items
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-cost-items-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostItemEdit(retailcrm.CostItem{
//		Code:            "seo",
//		Active:          false,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CostItemEdit(costItem CostItem) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&costItem)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"costItem": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/cost-items/%s/edit", costItem.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Couriers returns list of couriers
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-couriers
func (c *Client) Couriers() (CouriersResponse, int, error) {
	var resp CouriersResponse

	data, status, err := c.GetRequest("/reference/couriers")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CourierCreate creates a courier
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-couriers
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostItemEdit(retailcrm.Courier{
//		Active:    true,
//		Email:     "courier1@example.com",
//		FirstName: "Ivan",
//		LastName:  "Ivanov",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", data.ID)
//	}
func (c *Client) CourierCreate(courier Courier) (CreateResponse, int, error) {
	var resp CreateResponse

	objJSON, err := marshalToString(&courier)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"courier": {objJSON},
	}

	data, status, err := c.PostRequest("/reference/couriers/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CourierEdit edit a courier
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-couriers-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostItemEdit(retailcrm.Courier{
//		ID:    	  1,
//		Patronymic: "Ivanovich",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CourierEdit(courier Courier) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&courier)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"courier": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/couriers/%d/edit", courier.ID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryServices returns list of delivery services
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-delivery-services
func (c *Client) DeliveryServices() (DeliveryServiceResponse, int, error) {
	var resp DeliveryServiceResponse

	data, status, err := c.GetRequest("/reference/delivery-services")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryServiceEdit delivery service create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-delivery-services-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryServiceEdit(retailcrm.DeliveryService{
//		Active: false,
//		Code:   "delivery-1",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) DeliveryServiceEdit(deliveryService DeliveryService) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&deliveryService)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"deliveryService": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/delivery-services/%s/edit", deliveryService.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryTypes returns list of delivery types
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-delivery-types
func (c *Client) DeliveryTypes() (DeliveryTypesResponse, int, error) {
	var resp DeliveryTypesResponse

	data, status, err := c.GetRequest("/reference/delivery-types")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// DeliveryTypeEdit delivery type create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-delivery-types-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryTypeEdit(retailcrm.DeliveryType{
//		Active:        false,
//		Code:          "type-1",
//		DefaultCost:   300,
//		DefaultForCrm: false,
//	}
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) DeliveryTypeEdit(deliveryType DeliveryType) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&deliveryType)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"deliveryType": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/delivery-types/%s/edit", deliveryType.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// LegalEntities returns list of legal entities
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-legal-entities
func (c *Client) LegalEntities() (LegalEntitiesResponse, int, error) {
	var resp LegalEntitiesResponse

	data, status, err := c.GetRequest("/reference/legal-entities")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// LegalEntityEdit change information about legal entity
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-legal-entities-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.LegalEntityEdit(retailcrm.LegalEntity{
//		Code:          "legal-entity-1",
//		CertificateDate:   "2012-12-12",
//	}
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) LegalEntityEdit(legalEntity LegalEntity) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&legalEntity)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"legalEntity": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/legal-entities/%s/edit", legalEntity.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderMethods returns list of order methods
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-order-methods
func (c *Client) OrderMethods() (OrderMethodsResponse, int, error) {
	var resp OrderMethodsResponse

	data, status, err := c.GetRequest("/reference/order-methods")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderMethodEdit order method create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-order-methods-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderMethodEdit(retailcrm.OrderMethod{
//		Code:          "method-1",
//		Active:        false,
//		DefaultForCRM: false,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrderMethodEdit(orderMethod OrderMethod) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&orderMethod)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"orderMethod": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/order-methods/%s/edit", orderMethod.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderTypes return list of order types
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-order-types
func (c *Client) OrderTypes() (OrderTypesResponse, int, error) {
	var resp OrderTypesResponse

	data, status, err := c.GetRequest("/reference/order-types")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// OrderTypeEdit create/edit order type
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-order-methods-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrderTypeEdit(retailcrm.OrderType{
//		Code:          "order-type-1",
//		Active:        false,
//		DefaultForCRM: false,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) OrderTypeEdit(orderType OrderType) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&orderType)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"orderType": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/order-types/%s/edit", orderType.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PaymentStatuses returns list of payment statuses
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-payment-statuses
func (c *Client) PaymentStatuses() (PaymentStatusesResponse, int, error) {
	var resp PaymentStatusesResponse

	data, status, err := c.GetRequest("/reference/payment-statuses")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PaymentStatusEdit payment status creation/editing
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-payment-statuses-code-edit
func (c *Client) PaymentStatusEdit(paymentStatus PaymentStatus) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&paymentStatus)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"paymentStatus": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/payment-statuses/%s/edit", paymentStatus.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PaymentTypes returns list of payment types
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-payment-types
func (c *Client) PaymentTypes() (PaymentTypesResponse, int, error) {
	var resp PaymentTypesResponse

	data, status, err := c.GetRequest("/reference/payment-types")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PaymentTypeEdit payment type create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-payment-types-code-edit
func (c *Client) PaymentTypeEdit(paymentType PaymentType) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&paymentType)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"paymentType": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/payment-types/%s/edit", paymentType.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PriceTypes returns list of price types
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-price-types
func (c *Client) PriceTypes() (PriceTypesResponse, int, error) {
	var resp PriceTypesResponse

	data, status, err := c.GetRequest("/reference/price-types")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PriceTypeEdit price type create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-price-types-code-edit
func (c *Client) PriceTypeEdit(priceType PriceType) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&priceType)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"priceType": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/price-types/%s/edit", priceType.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ProductStatuses returns list of item statuses in order
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-product-statuses
func (c *Client) ProductStatuses() (ProductStatusesResponse, int, error) {
	var resp ProductStatusesResponse

	data, status, err := c.GetRequest("/reference/product-statuses")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ProductStatusEdit order item status create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-product-statuses-code-edit
func (c *Client) ProductStatusEdit(productStatus ProductStatus) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&productStatus)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"productStatus": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/product-statuses/%s/edit", productStatus.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Sites returns the sites list
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-sites
func (c *Client) Sites() (SitesResponse, int, error) {
	var resp SitesResponse

	data, status, err := c.GetRequest("/reference/sites")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// SiteEdit site create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-sites-code-edit
func (c *Client) SiteEdit(site Site) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&site)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"site": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/sites/%s/edit", site.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// StatusGroups returns list of order status groups
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-status-groups
func (c *Client) StatusGroups() (StatusGroupsResponse, int, error) {
	var resp StatusGroupsResponse

	data, status, err := c.GetRequest("/reference/status-groups")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Statuses returns list of order statuses
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-statuses
func (c *Client) Statuses() (StatusesResponse, int, error) {
	var resp StatusesResponse

	data, status, err := c.GetRequest("/reference/statuses")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// StatusEdit order status create/edit
//
// For more information see www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-sites-code-edit.
func (c *Client) StatusEdit(statusItem Status) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&statusItem)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"status": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/statuses/%s/edit", statusItem.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Stores returns list of warehouses
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-stores
func (c *Client) Stores() (StoresResponse, int, error) {
	var resp StoresResponse

	data, status, err := c.GetRequest("/reference/stores")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// StoreEdit warehouse create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-stores-code-edit
func (c *Client) StoreEdit(store Store) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&store)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"store": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/stores/%s/edit", store.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Units returns units list
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-reference-units
func (c *Client) Units() (UnitsResponse, int, error) {
	var resp UnitsResponse

	data, status, err := c.GetRequest("/reference/units")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// UnitEdit unit create/edit
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-reference-units-code-edit
func (c *Client) UnitEdit(unit Unit) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	objJSON, err := marshalToString(&unit)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"unit": {objJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/units/%s/edit", unit.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Segments returns segments
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-segments
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Segments(SegmentsRequest{
//		Filter: retailcrm.SegmentsFilter{
//			Ids: []int{1,2,3}
//		}
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	for _, value := range data.Segments {
//		fmt.Printf("%v\n", value)
//	}
func (c *Client) Segments(parameters SegmentsRequest) (SegmentsResponse, int, error) {
	var resp SegmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/segments?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Settings returns system settings
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-settings
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Settings()
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	fmt.Printf("%#v\n", data)
func (c *Client) Settings() (SettingsResponse, int, error) {
	var resp SettingsResponse

	data, status, err := c.GetRequest("/settings")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Inventories returns leftover stocks and purchasing prices
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-store-inventories
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Inventories(retailcrm.InventoriesRequest{Filter: retailcrm.InventoriesFilter{Details: 1, ProductActive: 1}, Page: 1})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Offers {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Inventories(parameters InventoriesRequest) (InventoriesResponse, int, error) {
	var resp InventoriesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/inventories?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// InventoriesUpload updates the leftover stocks and purchasing prices
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-store-inventories-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.InventoriesUpload(
//	   []retailcrm.InventoryUpload{
//		   {
//			   XMLID: "pT22K9YzX21HTdzFCe1",
//			   Stores: []InventoryUploadStore{
//				   {Code: "test-store-v5", Available: 10, PurchasePrice: 1500},
//				   {Code: "test-store-v4", Available: 20, PurchasePrice: 1530},
//				   {Code: "test-store", Available: 30, PurchasePrice: 1510},
//			   },
//		   },
//		   {
//			   XMLID: "JQICtiSpOV3AAfMiQB3",
//			   Stores: []InventoryUploadStore{
//				   {Code: "test-store-v5", Available: 45, PurchasePrice: 1500},
//				   {Code: "test-store-v4", Available: 32, PurchasePrice: 1530},
//				   {Code: "test-store", Available: 46, PurchasePrice: 1510},
//			   },
//		   },
//	   },
//	)
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	fmt.Printf("%v\n", data.NotFoundOffers)
func (c *Client) InventoriesUpload(inventories []InventoryUpload, site ...string) (StoreUploadResponse, int, error) {
	var resp StoreUploadResponse

	uploadJSON, err := marshalToString(&inventories)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"offers": {uploadJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/store/inventories/upload", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// PricesUpload updates prices
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-store-prices-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.PricesUpload([]retailcrm.OfferPriceUpload{
//		{
//			ID         1
//			Site       "store-1"
//			Prices     []PriceUpload{{Code:  "price-1"}}
//		},
//		{
//			ID         2
//			Site       "store-1"
//			Prices     []PriceUpload{{Code:  "price-2"}}
//		},
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	fmt.Printf("%v\n", data.NotFoundOffers)
func (c *Client) PricesUpload(prices []OfferPriceUpload) (StoreUploadResponse, int, error) {
	var resp StoreUploadResponse

	uploadJSON, err := marshalToString(&prices)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"prices": {uploadJSON},
	}

	data, status, err := c.PostRequest("/store/prices/upload", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ProductsGroup returns list of product groups
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-store-product-groups
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.ProductsGroup(retailcrm.ProductsGroupsRequest{
//		Filter: retailcrm.ProductsGroupsFilter{
//			Active: 1,
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.ProductGroup {
//		log.Printf("%v\n", value)
//	}
func (c *Client) ProductsGroup(parameters ProductsGroupsRequest) (ProductsGroupsResponse, int, error) {
	var resp ProductsGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/product-groups?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Products returns list of products and SKU
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-store-products
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Products(retailcrm.ProductsRequest{
//		Filter: retailcrm.ProductsFilter{
//			Active:   1,
//			MinPrice: 1000,
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Products {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Products(parameters ProductsRequest) (ProductsResponse, int, error) {
	var resp ProductsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/products?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ProductsProperties returns list of item properties, matching the specified filters
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-store-products-properties
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.ProductsProperties(retailcrm.ProductsPropertiesRequest{
//		Filter: retailcrm.ProductsPropertiesFilter{
//			Sites: []string["store"],
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Properties {
//		log.Printf("%v\n", value)
//	}
func (c *Client) ProductsProperties(parameters ProductsPropertiesRequest) (ProductsPropertiesResponse, int, error) {
	var resp ProductsPropertiesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/products/properties?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Tasks returns task list
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-tasks
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Tasks(retailcrm.TasksRequest{
//		Filter: TasksFilter{
//			DateFrom: "2012-12-12",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Tasks {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Tasks(parameters TasksRequest) (TasksResponse, int, error) {
	var resp TasksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/tasks?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// TaskCreate create a task
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-tasks-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Tasks(retailcrm.Task{
//		Text:        "task №1",
//		PerformerID: 12,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.ID)
//	}
func (c *Client) TaskCreate(task Task, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse
	taskJSON, err := marshalToString(&task)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"task": {taskJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/tasks/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Task returns task
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-tasks-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Task(12)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.Task)
//	}
func (c *Client) Task(id int) (TaskResponse, int, error) {
	var resp TaskResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/tasks/%d", id))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// TaskEdit edit a task
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-tasks-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Task(retailcrm.Task{
//		ID:   12
//		Text: "task №2",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) TaskEdit(task Task, site ...string) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse
	var uid = strconv.Itoa(task.ID)

	taskJSON, err := marshalToString(&task)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"task": {taskJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/tasks/%s/edit", uid), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	if !resp.Success {
		a := CreateAPIError(data)
		return resp, status, a
	}

	return resp, status, nil
}

// UserGroups returns list of user groups
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-user-groups
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.UserGroups(retailcrm.UserGroupsRequest{Page: 1})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Groups {
//		log.Printf("%v\n", value)
//	}
func (c *Client) UserGroups(parameters UserGroupsRequest) (UserGroupsResponse, int, error) {
	var resp UserGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/user-groups?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Users returns list of users matched the specified filters
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-users
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Users(retailcrm.UsersRequest{Filter: retailcrm.UsersFilter{Active: 1}, Page: 1})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Users {
//		log.Printf("%v\n", value)
//	}
func (c *Client) Users(parameters UsersRequest) (UsersResponse, int, error) {
	var resp UsersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/users?%s", params.Encode()))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// User returns information about user
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-users-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.User(12)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.User)
//	}
func (c *Client) User(id int) (UserResponse, int, error) {
	var resp UserResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/users/%d", id))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// UserStatus change user status
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-users
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.UserStatus(12, "busy")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) UserStatus(id int, status string) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	p := url.Values{
		"status": {status},
	}

	data, statusCode, err := c.PostRequest(fmt.Sprintf("/users/%d/status", id), p)
	if err != nil {
		return resp, statusCode, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, statusCode, err
	}

	if !resp.Success {
		return resp, statusCode, CreateAPIError(data)
	}

	return resp, statusCode, nil
}

// StaticticsUpdate updates statistics
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-statistic-update
func (c *Client) StaticticsUpdate() (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	data, status, err := c.GetRequest("/statistic/update")
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Costs returns costs list
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-costs
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Costs(CostsRequest{
//		Filter: CostsFilter{
//			IDs: []string{"1","2","3"},
//			MinSumm: "1000"
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.Costs {
//		log.Printf("%v\n", value.Summ)
//	}
func (c *Client) Costs(costs CostsRequest) (CostsResponse, int, error) {
	var resp CostsResponse

	params, _ := query.Values(costs)

	data, status, err := c.GetRequest(fmt.Sprintf("/costs?%s", params.Encode()))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostCreate create the cost
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-costs-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostCreate(
//		retailcrm.CostRecord{
//			DateFrom:  "2012-12-12",
//			DateTo:    "2012-12-12",
//			Summ:      12,
//			CostItem:  "calculation-of-costs",
//			Order: Order{
//				Number: "1"
//			},
//			Sites:    []string{"store"},
//		},
//		"store"
//	)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.ID)
//	}
func (c *Client) CostCreate(cost CostRecord, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse

	costJSON, err := marshalToString(&cost)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"cost": {costJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/costs/create", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostsDelete removes a cost
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-costs-delete
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.client.CostsDelete([]int{1, 2, 3, 48, 49, 50})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("Not removed costs: %v", data.NotRemovedIds)
//	}
func (c *Client) CostsDelete(ids []int) (CostsDeleteResponse, int, error) {
	var resp CostsDeleteResponse

	costJSON, err := marshalToString(&ids)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"ids": {costJSON},
	}

	data, status, err := c.PostRequest("/costs/delete", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostsUpload batch costs upload
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-costs-upload
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostCreate([]retailcrm.CostRecord{
//		{
//			DateFrom:  "2012-12-12",
//			DateTo:    "2012-12-12",
//			Summ:      12,
//			CostItem:  "calculation-of-costs",
//			Order: Order{
//				Number: "1"
//			},
//			Sites:    []string{"store"},
//		},
//		{
//			DateFrom:  "2012-12-13",
//			DateTo:    "2012-12-13",
//			Summ:      13,
//			CostItem:  "seo",
//		}
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("Uploaded costs: %v", data.UploadedCosts)
//	}
func (c *Client) CostsUpload(cost []CostRecord) (CostsUploadResponse, int, error) {
	var resp CostsUploadResponse

	costJSON, err := marshalToString(&cost)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"costs": {costJSON},
	}

	data, status, err := c.PostRequest("/costs/upload", p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Cost returns information about specified cost
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-costs-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Cost(1)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.Cost)
//	}
func (c *Client) Cost(id int) (CostResponse, int, error) {
	var resp CostResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/costs/%d", id))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostDelete removes cost
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-costs-id-delete
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostDelete(1)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) CostDelete(costID int) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	costJSON, err := marshalToString(&costID)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"costs": {costJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/costs/%d/delete", costID), p)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CostEdit edit a cost
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-costs-id-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostEdit(1, retailcrm.Cost{
//		DateFrom:  "2012-12-12",
//		DateTo:    "2018-12-13",
//		Summ:      321,
//		CostItem:  "seo",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.ID)
//	}
func (c *Client) CostEdit(costID int, cost CostRecord, site ...string) (CreateResponse, int, error) {
	var resp CreateResponse

	costJSON, err := marshalToString(&cost)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"cost": {costJSON},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/costs/%d/edit", costID), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// Files returns files list
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-files
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Files(FilesRequest{
//		Filter: FilesFilter{
//			Filename: "image.jpeg",
//		},
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) Files(files FilesRequest) (FilesResponse, int, error) {
	var resp FilesResponse

	params, _ := query.Values(files)

	data, status, err := c.GetRequest(fmt.Sprintf("/files?%s", params.Encode()))

	if err != nil && err.Error() != "" {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// FileUpload uploads file to RetailCRM
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-files
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	 file, err := os.Open("file.jpg")
//	 if err != nil {
//		    fmt.Print(err)
//	 }
//
//	 data, status, err := client.FileUpload(file)
//
//	 if err.Error() != "" {
//		    fmt.Printf("%v", err.Error())
//	 }
//
//	 if status >= http.StatusBadRequest {
//		    fmt.Printf("%v", err.Error())
//	 }
func (c *Client) FileUpload(reader io.Reader) (FileUploadResponse, int, error) {
	var resp FileUploadResponse

	data, status, err := c.PostRequest("/files/upload", reader, "application/octet-stream")

	if err != nil && err.Error() != "" {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// File returns a file info
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-files
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.File(112)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v\n", data.File)
//	}
func (c *Client) File(id int) (FileResponse, int, error) {
	var resp FileResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/files/%d", id))
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// FileDelete removes file from RetailCRM
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-files
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	 data, status, err := client.FileDelete(123)
//
//	 if err.Error() != "" {
//		    fmt.Printf("%v", err.Error())
//	 }
//
//	 if status >= http.StatusBadRequest {
//		    fmt.Printf("%v", err.Error())
//	 }
func (c *Client) FileDelete(id int) (SuccessfulResponse, int, error) {
	var resp SuccessfulResponse

	data, status, err := c.PostRequest(fmt.Sprintf("/files/%d/delete", id), strings.NewReader(""))

	if err != nil && err.Error() != "" {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// FileDownload downloads file from RetailCRM
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-files
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	 fileData, status, err := client.FileDownload(123)
//
//	 if err.Error() != "" {
//		    fmt.Printf("%v", err.Error())
//	 }
//
//	 if status >= http.StatusBadRequest {
//		    fmt.Printf("%v", err.Error())
//	 }
func (c *Client) FileDownload(id int) (io.ReadCloser, int, error) {
	data, status, err := c.GetRequest(fmt.Sprintf("/files/%d/download", id))
	if status != http.StatusOK {
		return nil, status, err
	}

	closer := ioutil.NopCloser(bytes.NewReader(data))
	return closer, status, nil
}

// FileEdit edits file name and relations with orders and customers in RetailCRM
//
// For more information see https://help.retailcrm.pro/Developers/ApiVersion5#get--api-v5-files
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	 data, status, err := client.FileEdit(123, File{Filename: "image2.jpg"})
//
//	 if err.Error() != "" {
//		    fmt.Printf("%v", err.Error())
//	 }
//
//	 if status >= http.StatusBadRequest {
//		    fmt.Printf("%v", err.Error())
//	 }
func (c *Client) FileEdit(id int, file File) (FileResponse, int, error) {
	var resp FileResponse

	fileJSON, err := marshalToString(file)
	if err != nil {
		return resp, 0, err
	}
	data, status, err := c.PostRequest(
		fmt.Sprintf("/files/%d/edit", id), url.Values{
			"file": {fileJSON},
		},
	)

	if err != nil && err.Error() != "" {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomFields returns list of custom fields
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-custom-fields
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomFields(retailcrm.CustomFieldsRequest{
//		Type: "string",
//		Entity: "customer",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.CustomFields {
//		fmt.Printf("%v\n", value)
//	}
func (c *Client) CustomFields(customFields CustomFieldsRequest) (CustomFieldsResponse, int, error) {
	var resp CustomFieldsResponse

	params, _ := query.Values(customFields)

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields?%s", params.Encode()))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomDictionaries returns list of custom directory
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-custom-fields-dictionaries
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionaries(retailcrm.CustomDictionariesRequest{
//		Filter: retailcrm.CustomDictionariesFilter{
//			Name: "Dictionary-1",
//		},
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	for _, value := range data.CustomDictionaries {
//		fmt.Printf("%v\n", value.Elements)
//	}
func (c *Client) CustomDictionaries(customDictionaries CustomDictionariesRequest) (
	CustomDictionariesResponse, int, error,
) {
	var resp CustomDictionariesResponse

	params, _ := query.Values(customDictionaries)

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields/dictionaries?%s", params.Encode()))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomDictionariesCreate creates a custom dictionary
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-custom-fields-dictionaries-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionariesCreate(retailcrm.CustomDictionary{
//		Name: "Courier profiles",
//		Code: "courier-profiles",
//		Elements: []Element{
//			{
//				Name: "Name",
//				Code: "name",
//			},
//			{
//				Name: "Lastname",
//				Code: "lastname",
//			}
//		},
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	If data.Success == true {
//		fmt.Printf("%v", data.Code)
//	}
func (c *Client) CustomDictionariesCreate(customDictionary CustomDictionary) (CustomResponse, int, error) {
	var resp CustomResponse

	costJSON, err := marshalToString(&customDictionary)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customDictionary": {costJSON},
	}

	data, status, err := c.PostRequest("/custom-fields/dictionaries/create", p)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomDictionary returns information about dictionary
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-custom-fields-entity-code
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionary("courier-profiles")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.CustomDictionary.Name)
//	}
func (c *Client) CustomDictionary(code string) (CustomDictionaryResponse, int, error) {
	var resp CustomDictionaryResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields/dictionaries/%s", code))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomDictionaryEdit edit custom dictionary
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-custom-fields-dictionaries-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionaryEdit(retailcrm.CustomDictionary{
//	Name: "Courier profiles",
//		Code: "courier-profiles",
//		Elements: []Element{
//			{
//				Name: "Name",
//				Code: "name",
//			},
//			{
//				Name: "Lastname",
//				Code: "lastname",
//			}
//		},
//	})
//
//	if err.Error() != "" {
//		fmt.Printf("%v", err.Error())
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.Error())
//	}
//
//	If data.Success == true {
//		fmt.Printf("%v", data.Code)
//	}
func (c *Client) CustomDictionaryEdit(customDictionary CustomDictionary) (CustomResponse, int, error) {
	var resp CustomResponse

	costJSON, err := marshalToString(&customDictionary)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customDictionary": {costJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/custom-fields/dictionaries/%s/edit", customDictionary.Code), p)
	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomFieldsCreate creates custom field
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-custom-fields-entity-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomFieldsCreate(CustomFields{
//		Name:        "First order",
//		Code:        "first-order",
//		Type:        "bool",
//		Entity:      "order",
//		DisplayArea: "customer",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.Code)
//	}
func (c *Client) CustomFieldsCreate(customFields CustomFields) (CustomResponse, int, error) {
	var resp CustomResponse

	costJSON, err := marshalToString(&customFields)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customField": {costJSON},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/custom-fields/%s/create", customFields.Entity), p)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomField returns information about custom fields
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-custom-fields-entity-code
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomField("order", "first-order")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.CustomField)
//	}
func (c *Client) CustomField(entity, code string) (CustomFieldResponse, int, error) {
	var resp CustomFieldResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields/%s/%s", entity, code))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// CustomFieldEdit list method
//
// For more information see http://www.simla.com/docs/Developers/API/APIVersions/APIv5#post--api-v5-custom-fields-entity-code-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomFieldEdit(CustomFields{
//		Code:        "first-order",
//		Entity:      "order",
//		DisplayArea: "delivery",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	If data.Success == true {
//		log.Printf("%v", data.Code)
//	}
func (c *Client) CustomFieldEdit(customFields CustomFields) (CustomResponse, int, error) {
	var resp CustomResponse

	costJSON, err := marshalToString(&customFields)
	if err != nil {
		return resp, 0, err
	}

	p := url.Values{
		"customField": {costJSON},
	}

	data, status, err := c.PostRequest(
		fmt.Sprintf("/custom-fields/%s/%s/edit", customFields.Entity, customFields.Code), p,
	)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// BonusOperations returns history of the bonus account for all participations
//
// For more information see https://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-loyalty-bonus-operations
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.BonusOperations(retailcrm.BonusOperationsRequest{
//		Filter: BonusOperationsFilter{
//			LoyaltyId: []int{123, 456},
//		},
//		Limit: 20,
//		Cursor: "123456",
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.BonusOperations {
//		log.Printf("%v\n", value)
//	}
func (c *Client) BonusOperations(parameters BonusOperationsRequest) (BonusOperationsResponse, int, error) {
	var resp BonusOperationsResponse

	params, _ := query.Values(parameters)
	data, status, err := c.GetRequest(fmt.Sprintf("/loyalty/bonus/operations?%s", params.Encode()))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// AccountBonusOperations returns history of the bonus account for a specific participation
//
// For more information see https://www.simla.com/docs/Developers/API/APIVersions/APIv5#get--api-v5-loyalty-account-id-bonus-operations
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.AccountBonusOperations(retailcrm.AccountBonusOperationsRequest{
//		Filter: AccountBonusOperationsFilter{
//			CreatedAtFrom: "2012-12-12",
//			CreatedAtTo: "2012-12-13",
//		},
//		Limit: 20,
//		Page: 3,
//	})
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	for _, value := range data.BonusOperations {
//		log.Printf("%v\n", value)
//	}
func (c *Client) AccountBonusOperations(id int, parameters AccountBonusOperationsRequest) (BonusOperationsResponse, int, error) {
	var resp BonusOperationsResponse

	if id == 0 {
		c.writeLog("cannot get loyalty bonus operations for user with id %d", id)
	}

	params, _ := query.Values(parameters)
	data, status, err := c.GetRequest(fmt.Sprintf(
		"/loyalty/account/%d/bonus/operations?%s",
		id, params.Encode(),
	))

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ProductsBatchEdit perform editing products batch
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-store-products-batch-edit
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	 	products := []ProductEdit{{
//				CatalogID: 3,
//				BaseProduct: BaseProduct{
//					Name:         "Product 1",
//					URL:          "https://example.com/p/1",
//					Article:      "p1",
//					ExternalID:   "ext1",
//				},
//				ID:     194,
//				Site:   "second",
//				Groups: []ProductEditGroupInput{{ID: 19}},
//		}}
//
//		data, status, err := client.ProductsBatchEdit(products)
//
//		if err != nil {
//			if apiErr, ok := retailcrm.AsAPIError(err); ok {
//				log.Fatalf("http status: %d, %s", status, apiErr.String())
//			}
//
//			log.Fatalf("http status: %d, error: %s", status, err)
//		}
//
//		if data.Success == true {
//			log.Printf("%v", data.ProcessedProductsCount)
//		}
func (c *Client) ProductsBatchEdit(products []ProductEdit) (ProductsBatchEditResponse, int, error) {
	var resp ProductsBatchEditResponse

	productsEditJSON, err := marshalToString(products)
	if err != nil {
		return resp, 0, err
	}
	p := url.Values{
		"products": {productsEditJSON},
	}

	data, status, err := c.PostRequest("/store/products/batch/edit", p)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// ProductsBatchCreate perform adding products batch
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-store-products-batch-create
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	 	products := []ProductCreate{{
//				CatalogID: 3,
//				BaseProduct: BaseProduct{
//					Name:         "Product 1",
//					URL:          "https://example.com/p/1",
//					Article:      "p1",
//					ExternalID:   "ext1",
//				},
//				Groups: []ProductEditGroupInput{{ID: 19}},
//		}}
//
//		data, status, err := client.ProductsBatchCreate(products)
//
//		if err != nil {
//			if apiErr, ok := retailcrm.AsAPIError(err); ok {
//				log.Fatalf("http status: %d, %s", status, apiErr.String())
//			}
//
//			log.Fatalf("http status: %d, error: %s", status, err)
//		}
//
//		if data.Success == true {
//			log.Printf("%v", data.AddedProducts)
//		}
func (c *Client) ProductsBatchCreate(products []ProductCreate) (ProductsBatchEditResponse, int, error) {
	var resp ProductsBatchEditResponse

	productsEditJSON, err := marshalToString(products)
	if err != nil {
		return resp, 0, err
	}
	p := url.Values{
		"products": {productsEditJSON},
	}

	data, status, err := c.PostRequest("/store/products/batch/create", p)

	if err != nil {
		return resp, status, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, status, err
	}

	return resp, status, nil
}

// LoyaltyAccountCreate аdd a client to the loyalty program
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-loyalty-account-create
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//	 	acc := SerializedCreateLoyaltyAccount{
//			SerializedBaseLoyaltyAccount: SerializedBaseLoyaltyAccount{
//				PhoneNumber:  "89151005004",
//				CustomFields: []string{"dog"},
//			},
//			Customer: SerializedEntityCustomer{
//				ID: 123,
//			},
//		}
//
//		data, status, err := client.LoyaltyAccountCreate("site", acc)
//
//		if err != nil {
//			if apiErr, ok := retailcrm.AsAPIError(err); ok {
//				log.Fatalf("http status: %d, %s", status, apiErr.String())
//			}
//
//			log.Fatalf("http status: %d, error: %s", status, err)
//		}
//
//		if data.Success == true {
//			log.Printf("%v", data.LoyaltyAccount.ID)
//		}
func (c *Client) LoyaltyAccountCreate(site string, loyaltyAccount SerializedCreateLoyaltyAccount) (CreateLoyaltyAccountResponse, int, error) {
	var result CreateLoyaltyAccountResponse

	loyaltyAccountJSON, err := marshalToString(loyaltyAccount)
	if err != nil {
		return result, 0, err
	}
	p := url.Values{
		"site":           {site},
		"loyaltyAccount": {loyaltyAccountJSON},
	}

	resp, status, err := c.PostRequest("/loyalty/account/create", p)

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyAccountEdit edit a client in the loyalty program
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-loyalty-account-id-edit
//
// Example:
//
//		var client = retailcrm.New("https://demo.url", "09jIJ")
//	 	acc := SerializedEditLoyaltyAccount{
//			SerializedBaseLoyaltyAccount: SerializedBaseLoyaltyAccount{
//				PhoneNumber:  "89151005004",
//				CustomFields: []string{"dog"},
//			},
//		}
//
//		data, status, err := client.LoyaltyAccountEdit(13, acc)
//
//		if err != nil {
//			if apiErr, ok := retailcrm.AsAPIError(err); ok {
//				log.Fatalf("http status: %d, %s", status, apiErr.String())
//			}
//
//			log.Fatalf("http status: %d, error: %s", status, err)
//		}
//
//		if data.Success == true {
//			log.Printf("%v", data.LoyaltyAccount.PhoneNumber)
//		}
func (c *Client) LoyaltyAccountEdit(accountID int, loyaltyAccount SerializedEditLoyaltyAccount) (EditLoyaltyAccountResponse, int, error) {
	var result EditLoyaltyAccountResponse

	loyaltyAccountJSON, err := marshalToString(loyaltyAccount)
	if err != nil {
		return result, 0, err
	}
	p := url.Values{
		"loyaltyAccount": {loyaltyAccountJSON},
	}

	resp, status, err := c.PostRequest(fmt.Sprintf("/loyalty/account/%d/edit", accountID), p)

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyAccount return information about client in the loyalty program
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-loyalty-account-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.LoyaltyAccount(13)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", data.LoyaltyAccount.PhoneNumber)
//	}
func (c *Client) LoyaltyAccount(id int) (LoyaltyAccountResponse, int, error) {
	var result LoyaltyAccountResponse

	resp, status, err := c.GetRequest(fmt.Sprintf("/loyalty/account/%d", id))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyAccountActivate activate participation in the loyalty program for client
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-loyalty-account-id-activate
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.LoyaltyAccountActivate(13)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", data.LoyaltyAccount.Active)
//	}
func (c *Client) LoyaltyAccountActivate(id int) (LoyaltyAccountActivateResponse, int, error) {
	var result LoyaltyAccountActivateResponse

	resp, status, err := c.PostRequest(fmt.Sprintf("/loyalty/account/%d/activate", id), strings.NewReader(""))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyBonusCredit accrue bonus participation in the program of loyalty
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-loyalty-account-id-bonus-credit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	req := LoyaltyBonusCreditRequest{
//		Amount:      120,
//		ExpiredDate: "2023-11-24 12:39:37",
//		Comment:     "Test",
//	}
//
// data, status, err := client.LoyaltyBonusCredit(13, req)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", data.LoyaltyBonus.ActivationDate)
//	}
func (c *Client) LoyaltyBonusCredit(id int, req LoyaltyBonusCreditRequest) (LoyaltyBonusCreditResponse, int, error) {
	var result LoyaltyBonusCreditResponse
	p, _ := query.Values(req)

	resp, status, err := c.PostRequest(fmt.Sprintf("/loyalty/account/%d/bonus/credit", id), p)

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyBonusStatusDetails get details on the bonus account
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-loyalty-account-id-bonus-status-details
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	req := LoyaltyBonusStatusDetailsRequest{
//		Limit: 20,
//		Page:  3,
//	}
//
// data, status, err := client.LoyaltyBonusStatusDetails(13, "waiting_activation", req)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		for _, bonus := range data.Bonuses {
//			log.Printf("%v", bonus.Amount)
//		}
//	}
func (c *Client) LoyaltyBonusStatusDetails(
	id int, statusType string, request LoyaltyBonusStatusDetailsRequest,
) (LoyaltyBonusDetailsResponse, int, error) {
	var result LoyaltyBonusDetailsResponse

	p, _ := query.Values(request)

	resp, status, err := c.GetRequest(fmt.Sprintf("/loyalty/account/%d/bonus/%s/details?%s", id, statusType, p.Encode()))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyAccounts return list of participations in the loyalty program
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-loyalty-accounts
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	req := LoyaltyAccountsRequest{
//		Filter: LoyaltyAccountAPIFilter{
//			Status:      "activated",
//			PhoneNumber: "89185556363",
//			Ids:         []int{14},
//			Level:       5,
//			Loyalties:   []int{2},
//			CustomerID:  "109",
//		},
//	}
//
// data, status, err := client.LoyaltyAccounts(req)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		for _, account := range data.LoyaltyAccounts {
//			log.Printf("%v", account.Status)
//		}
//	}
func (c *Client) LoyaltyAccounts(req LoyaltyAccountsRequest) (LoyaltyAccountsResponse, int, error) {
	var result LoyaltyAccountsResponse

	p, _ := query.Values(req)

	resp, status, err := c.GetRequest(fmt.Sprintf("/loyalty/accounts?%s", p.Encode()))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// LoyaltyCalculate calculations of the maximum discount
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-loyalty-calculate
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	req := LoyaltyCalculateRequest{
//		Site: "main",
//		Order: Order{
//			PrivilegeType: "loyalty_level",
//			Customer: &Customer{
//				ID: 123,
//			},
//			Items: []OrderItem{
//				{
//					InitialPrice: 10000,
//					Quantity:     1,
//					Offer:        Offer{ID: 214},
//					PriceType:    &PriceType{Code: "base"},
//				},
//			},
//		},
//		Bonuses: 10,
//	}
//
// data, status, err := client.LoyaltyCalculate(req)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", data.Order.PrivilegeType)
//		log.Printf("%v", data.Order.BonusesCreditTotal)
//	}
func (c *Client) LoyaltyCalculate(req LoyaltyCalculateRequest) (LoyaltyCalculateResponse, int, error) {
	var result LoyaltyCalculateResponse

	orderJSON, err := marshalToString(req.Order)
	if err != nil {
		return result, 0, err
	}

	p := url.Values{
		"site":    {req.Site},
		"order":   {orderJSON},
		"bonuses": {fmt.Sprintf("%f", req.Bonuses)},
	}

	resp, status, err := c.PostRequest("/loyalty/calculate", p)

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// GetLoyalties returns list of loyalty programs
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-loyalty-loyalties
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	req := LoyaltiesRequest{
//		Filter: LoyaltyAPIFilter{
//			Active: active,
//			Ids:    []int{2},
//			Sites:  []string{"main"},
//		},
//	}
//
// data, status, err := client.GetLoyalties(req)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		for _, l := range data.Loyalties {
//			log.Printf("%v", l.ID)
//			log.Printf("%v", l.Active)
//		}
//	}
func (c *Client) GetLoyalties(req LoyaltiesRequest) (LoyaltiesResponse, int, error) {
	var result LoyaltiesResponse

	p, _ := query.Values(req)

	resp, status, err := c.GetRequest(fmt.Sprintf("/loyalty/loyalties?%s", p.Encode()))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// GetLoyaltyByID return program of loyalty by id
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#get--api-v5-loyalty-loyalties-id
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
// data, status, err := client.GetLoyaltyByID(2)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", res.Loyalty.ID)
//		log.Printf("%v", res.Loyalty.Active)
//	}
func (c *Client) GetLoyaltyByID(id int) (LoyaltyResponse, int, error) {
	var result LoyaltyResponse

	resp, status, err := c.GetRequest(fmt.Sprintf("/loyalty/loyalties/%d", id))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// OrderIntegrationDeliveryCancel cancels of integration delivery
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-orders-externalId-delivery-cancel
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
// data, status, err := client.OrderIntegrationDeliveryCancel("externalId", false, "1001C")
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", res.Success)
//	}
func (c *Client) OrderIntegrationDeliveryCancel(by string, force bool, id string) (SuccessfulResponse, int, error) {
	var result SuccessfulResponse

	p := url.Values{
		"by":    {checkBy(by)},
		"force": {fmt.Sprintf("%t", force)},
	}

	resp, status, err := c.PostRequest(fmt.Sprintf("/orders/%s/delivery/cancel?%s", id, p.Encode()), strings.NewReader(""))

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// CreateProductsGroup adds a product group
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-store-product-groups-create
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	group := ProductGroup{
//		ParentID:    125,
//		Name:        "Фрукты",
//		Site:        "main",
//		Active:      true,
//		ExternalID:  "abc22",
//	}
//
// data, status, err := client.CreateProductsGroup(group)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", res.ID)
//	}
func (c *Client) CreateProductsGroup(group ProductGroup) (ActionProductsGroupResponse, int, error) {
	var result ActionProductsGroupResponse

	groupJSON, err := marshalToString(group)
	if err != nil {
		return result, 0, err
	}

	p := url.Values{
		"productGroup": {groupJSON},
	}

	resp, status, err := c.PostRequest("/store/product-groups/create", p)

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// EditProductsGroup edits a product group
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-store-product-groups-externalId-edit
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	group := ProductGroup{
//		Name:        "Овощи",
//		Active:      true,
//		ExternalID:  "abc22",
//	}
//
// data, status, err := client.EditProductsGroup("id", "125", "main", group)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data.Success == true {
//		log.Printf("%v", res.ID)
//	}
func (c *Client) EditProductsGroup(by, groupID, site string, group ProductGroup) (ActionProductsGroupResponse, int, error) {
	var result ActionProductsGroupResponse

	groupJSON, err := marshalToString(group)
	if err != nil {
		return result, 0, err
	}

	p := url.Values{
		"by":           {checkBy(by)},
		"site":         {site},
		"productGroup": {groupJSON},
	}

	resp, status, err := c.PostRequest(fmt.Sprintf("/store/product-groups/%s/edit", groupID), p)

	if err != nil {
		return result, status, err
	}

	err = json.Unmarshal(resp, &result)

	if err != nil {
		return result, status, err
	}

	return result, status, nil
}

// GetOrderPlate receives a print form file for the order
//
// # Body of response is already closed
//
// For more information see https://help.retailcrm.ru/api_v5_ru.html#get--api-v5-orders-externalId-plates-plateId-print
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
// data, status, err := client.GetOrderPlate("id", "107", "main", 1)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
//
//	if data != nil {
//		fileData, err := io.ReadAll(data)
//		if err != nil {
//			return
//		}
//
//		log.Printf("%s", fileData)
//	}
func (c *Client) GetOrderPlate(by, orderID, site string, plateID int) (io.ReadCloser, int, error) {
	requestURL := fmt.Sprintf("%s/api/v5/orders/%s/plates/%d/print?%s", c.URL, orderID, plateID, url.Values{
		"by":   {checkBy(by)},
		"site": {site},
	}.Encode())

	return c.executeWithRetryReadCloser(requestURL, func() (interface{}, *http.Response, int, error) {
		req, err := http.NewRequest("GET", requestURL, nil)
		if err != nil {
			return nil, nil, 0, err
		}

		req.Header.Set("X-API-KEY", c.Key)

		if c.Debug {
			c.writeLog("API Request: %s %s", requestURL, c.Key)
		}

		resp, err := c.httpClient.Do(req)

		if err != nil {
			return nil, resp, 0, err
		}

		if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != http.StatusServiceUnavailable {
			return nil, resp, resp.StatusCode, CreateGenericAPIError(
				fmt.Sprintf("HTTP request error. Status code: %d.", resp.StatusCode))
		}

		if resp.StatusCode >= http.StatusBadRequest &&
			resp.StatusCode < http.StatusInternalServerError &&
			resp.StatusCode != http.StatusServiceUnavailable {
			res, err := buildRawResponse(resp)

			if err != nil {
				return nil, resp, 0, err
			}

			return nil, resp, resp.StatusCode, CreateAPIError(res)
		}

		if err != nil {
			return nil, resp, 0, err
		}

		return resp.Body, resp, resp.StatusCode, nil
	})
}

// NotificationsSend send a notification
//
// For more information see https://docs.retailcrm.ru/Developers/API/APIv5#post--api-v5-notifications-send
//
// Example:
//
//	var client = retailcrm.New("https://demo.url", "09jIJ")
//
//	req := retailcrm.NotificationsSendRequest{
//		UserGroups: []retailcrm.UserGroupType{retailcrm.UserGroupSuperadmins},
//		Type:       retailcrm.NotificationTypeInfo,
//		Message:    "Hello everyone!",
//	}
//
//	status, err := client.NotificationsSend(req)
//
//	if err != nil {
//		if apiErr, ok := retailcrm.AsAPIError(err); ok {
//			log.Fatalf("http status: %d, %s", status, apiErr.String())
//		}
//
//		log.Fatalf("http status: %d, error: %s", status, err)
//	}
func (c *Client) NotificationsSend(req NotificationsSendRequest) (int, error) {
	marshaled, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}

	_, status, err := c.PostRequest("/notifications/send",
		url.Values{"notification": {string(marshaled)}})
	if err != nil {
		return status, err
	}

	return status, nil
}

func (c *Client) ListMGChannelTemplates(channelID, page, limit int) (MGChannelTemplatesResponse, int, error) {
	var resp MGChannelTemplatesResponse

	values := url.Values{
		"page":       {fmt.Sprintf("%d", page)},
		"limit":      {fmt.Sprintf("%d", limit)},
		"channel_id": {fmt.Sprintf("%d", channelID)},
	}

	data, code, err := c.GetRequest(fmt.Sprintf("/reference/mg-channels/templates?%s", values.Encode()))

	if err != nil {
		return resp, code, err
	}

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return resp, code, err
	}

	return resp, code, nil
}

func (c *Client) EditMGChannelTemplate(req EditMGChannelTemplateRequest) (int, error) {
	templates, err := json.Marshal(req.Templates)

	if err != nil {
		return 0, err
	}

	if string(templates) == "null" {
		templates = []byte(`[]`)
	}

	removed, err := json.Marshal(req.Removed)

	if err != nil {
		return 0, err
	}

	if string(removed) == "null" {
		removed = []byte(`[]`)
	}

	values := url.Values{
		"templates": {string(templates)},
		"removed":   {string(removed)},
	}

	_, code, err := c.PostRequest("/reference/mg-channels/templates/edit", values)

	if err != nil {
		return code, err
	}

	return code, nil
}

func (c *Client) StoreOffers(req OffersRequest) (StoreOffersResponse, int, error) {
	var result StoreOffersResponse

	filter, err := query.Values(req)

	if err != nil {
		return StoreOffersResponse{}, 0, err
	}

	resp, status, err := c.GetRequest(fmt.Sprintf("/store/offers?%s", filter.Encode()))

	if err != nil {
		return StoreOffersResponse{}, status, err
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return StoreOffersResponse{}, status, err
	}

	return result, status, nil
}
