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

// CreateResponse type
type CreateResponse struct {
	Success bool `json:"success"`
	Id      int  `json:"id,omitempty"`
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
	Success           bool              `json:"success"`
	UploadedCustomers []IdentifiersPair `json:"uploadedCustomers,omitempty,brackets"`
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

// OrdersUploadResponse type
type OrdersUploadResponse struct {
	Success        bool              `json:"success"`
	UploadedOrders []IdentifiersPair `json:"uploadedOrders,omitempty,brackets"`
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

// TasksResponse type
type TasksResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Tasks      []Task      `json:"tasks,omitempty,brackets"`
}

// NotesResponse type
type NotesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Notes      []Note      `json:"notes,omitempty,brackets"`
}

// SegmentsResponse type
type SegmentsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Segments   []Segment   `json:"segments,omitempty,brackets"`
}

// CountriesResponse type
type CountriesResponse struct {
	Success      bool     `json:"success"`
	CountriesIso []string `json:"countriesIso,omitempty,brackets"`
}

// CostGroupsResponse type
type CostGroupsResponse struct {
	Success    bool        `json:"success"`
	CostGroups []CostGroup `json:"costGroups,omitempty,brackets"`
}

// CostItemsResponse type
type CostItemsResponse struct {
	Success   bool       `json:"success"`
	CostItems []CostItem `json:"costItems,omitempty,brackets"`
}

// CouriersResponse type
type CouriersResponse struct {
	Success  bool      `json:"success"`
	Couriers []Courier `json:"couriers,omitempty,brackets"`
}

// DeliveryServiceResponse type
type DeliveryServiceResponse struct {
	Success          bool                       `json:"success"`
	DeliveryServices map[string]DeliveryService `json:"deliveryServices,omitempty,brackets"`
}

// DeliveryTypesResponse type
type DeliveryTypesResponse struct {
	Success       bool                    `json:"success"`
	DeliveryTypes map[string]DeliveryType `json:"deliveryTypes,omitempty,brackets"`
}

// LegalEntitiesResponse type
type LegalEntitiesResponse struct {
	Success       bool          `json:"success"`
	LegalEntities []LegalEntity `json:"legalEntities,omitempty,brackets"`
}

// OrderMethodsResponse type
type OrderMethodsResponse struct {
	Success      bool                   `json:"success"`
	OrderMethods map[string]OrderMethod `json:"orderMethods,omitempty,brackets"`
}

// OrderTypesResponse type
type OrderTypesResponse struct {
	Success    bool                 `json:"success"`
	OrderTypes map[string]OrderType `json:"orderTypes,omitempty,brackets"`
}

// PaymentStatusesResponse type
type PaymentStatusesResponse struct {
	Success         bool                     `json:"success"`
	PaymentStatuses map[string]PaymentStatus `json:"paymentStatuses,omitempty,brackets"`
}

// PaymentTypesResponse type
type PaymentTypesResponse struct {
	Success      bool                   `json:"success"`
	PaymentTypes map[string]PaymentType `json:"paymentTypes,omitempty,brackets"`
}

// PriceTypesResponse type
type PriceTypesResponse struct {
	Success    bool        `json:"success"`
	PriceTypes []PriceType `json:"priceTypes,omitempty,brackets"`
}

// ProductStatusesResponse type
type ProductStatusesResponse struct {
	Success         bool                     `json:"success"`
	ProductStatuses map[string]ProductStatus `json:"productStatuses,omitempty,brackets"`
}

// StatusesResponse type
type StatusesResponse struct {
	Success  bool              `json:"success"`
	Statuses map[string]Status `json:"statuses,omitempty,brackets"`
}

// StatusGroupsResponse type
type StatusGroupsResponse struct {
	Success      bool                   `json:"success"`
	StatusGroups map[string]StatusGroup `json:"statusGroups,omitempty,brackets"`
}

// SitesResponse type
type SitesResponse struct {
	Success bool            `json:"success"`
	Sites   map[string]Site `json:"sites,omitempty,brackets"`
}

// StoresResponse type
type StoresResponse struct {
	Success bool    `json:"success"`
	Stores  []Store `json:"stores,omitempty,brackets"`
}
