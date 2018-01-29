package v5

import "encoding/json"

// ErrorResponse type
type ErrorResponse struct {
	ErrorMsg string            `json:"errorMsg,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

// ErrorResponse method
func (c *Client) ErrorResponse(data []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	err := json.Unmarshal(data, &resp)

	return &resp, err
}

// SucessfulResponse type
type SucessfulResponse struct {
	Success bool `json:"success,omitempty"`
}

// VersionResponse return available API versions
type VersionResponse struct {
	Success  bool     `json:"success,omitempty"`
	Versions []string `json:"versions,brackets,omitempty"`
}

// CredentialResponse return available API methods
type CredentialResponse struct {
	Success        bool     `json:"success,omitempty"`
	Credentials    []string `json:"credentials,brackets,omitempty"`
	SiteAccess     string   `json:"siteAccess,omitempty"`
	SitesAvailable []string `json:"sitesAvailable,brackets,omitempty"`
}

// CustomerResponse type
type CustomerResponse struct {
	Success  bool      `json:"success"`
	Customer *Customer `json:"customer,omitempty,brackets"`
}

// CustomersResponse type
type CustomersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Customers  []Customer  `json:"customers,omitempty,brackets"`
}

// CustomerChangeResponse type
type CustomerChangeResponse struct {
	Success bool   `json:"success"`
	Id      int    `json:"id,omitempty"`
	State   string `json:"state,omitempty"`
}

// CustomersUploadResponse type
type CustomersUploadResponse struct {
	Success           bool                  `json:"success"`
	UploadedCustomers []CustomerIdentifiers `json:"uploadedCustomers,omitempty,brackets"`
}

// CustomersHistoryResponse type
type CustomersHistoryResponse struct {
	Success     bool                    `json:"success,omitempty"`
	GeneratedAt string                  `json:"generatedAt,omitempty"`
	History     []CustomerHistoryRecord `json:"history,omitempty,brackets"`
	Pagination  *Pagination             `json:"pagination,omitempty"`
}

// OrderResponse type
type OrderResponse struct {
	Success bool   `json:"success"`
	Order   *Order `json:"order,omitempty,brackets"`
}

// OrdersResponse type
type OrdersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Orders     []Order     `json:"orders,omitempty,brackets"`
}

// OrdersHistoryResponse type
type OrdersHistoryResponse struct {
	Success     bool                  `json:"success,omitempty"`
	GeneratedAt string                `json:"generatedAt,omitempty"`
	History     []OrdersHistoryRecord `json:"history,omitempty,brackets"`
	Pagination  *Pagination           `json:"pagination,omitempty"`
}

// UserResponse type
type UserResponse struct {
	Success bool  `json:"success"`
	User    *User `json:"user,omitempty,brackets"`
}

// UsersResponse type
type UsersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Users      []User      `json:"users,omitempty,brackets"`
}

// UserGroupsResponse type
type UserGroupsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Groups     []UserGroup `json:"groups,omitempty,brackets"`
}

// TaskResponse type
type TaskResponse struct {
	Success bool  `json:"success"`
	Task    *Task `json:"task,omitempty,brackets"`
}

// TaskChangeResponse type
type TaskChangeResponse struct {
	Success bool `json:"success"`
	Id      int  `json:"id,omitempty"`
}

// TasksResponse type
type TasksResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Tasks      []Task      `json:"tasks,omitempty,brackets"`
}
