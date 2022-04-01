package retailcrm

// SuccessfulResponse type.
type SuccessfulResponse struct {
	Success bool `json:"success,omitempty"`
}

// CreateResponse type.
type CreateResponse struct {
	Success bool `json:"success"`
	ID      int  `json:"id,omitempty"`
}

// OrderCreateResponse type.
type OrderCreateResponse struct {
	CreateResponse
	Order Order `json:"order,omitempty"`
}

// OperationResponse type.
type OperationResponse struct {
	Success bool              `json:"success"`
	Errors  map[string]string `json:"ErrorsList,omitempty"`
}

// VersionResponse return available API versions.
type VersionResponse struct {
	Success  bool     `json:"success,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

// CredentialResponse return available API methods.
type CredentialResponse struct {
	Success bool `json:"success,omitempty"`
	// deprecated
	Credentials    []string `json:"credentials,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	SiteAccess     string   `json:"siteAccess,omitempty"`
	SitesAvailable []string `json:"sitesAvailable,omitempty"`
}

// SystemInfoResponse return system info.
type SystemInfoResponse struct {
	Success       bool   `json:"success,omitempty"`
	SystemVersion string `json:"systemVersion,omitempty"`
	PublicURL     string `json:"publicUrl,omitempty"`
	TechnicalURL  string `json:"technicalUrl,omitempty"`
}

// CustomerResponse type.
type CustomerResponse struct {
	Success  bool      `json:"success"`
	Customer *Customer `json:"customer,omitempty"`
}

// CorporateCustomerResponse type.
type CorporateCustomerResponse struct {
	Success           bool               `json:"success"`
	CorporateCustomer *CorporateCustomer `json:"customerCorporate,omitempty"`
}

// CustomersResponse type.
type CustomersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Customers  []Customer  `json:"customers,omitempty"`
}

// CorporateCustomersResponse type.
type CorporateCustomersResponse struct {
	Success            bool                `json:"success"`
	Pagination         *Pagination         `json:"pagination,omitempty"`
	CustomersCorporate []CorporateCustomer `json:"customersCorporate,omitempty"`
}

// CorporateCustomersNotesResponse type.
type CorporateCustomersNotesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Notes      []Note      `json:"notes,omitempty"`
}

// CorporateCustomersAddressesResponse type.
type CorporateCustomersAddressesResponse struct {
	Success   bool                       `json:"success"`
	Addresses []CorporateCustomerAddress `json:"addresses"`
}

// CorporateCustomerCompaniesResponse type.
type CorporateCustomerCompaniesResponse struct {
	Success   bool      `json:"success"`
	Companies []Company `json:"companies"`
}

// CorporateCustomerContactsResponse type.
type CorporateCustomerContactsResponse struct {
	Success  bool                       `json:"success"`
	Contacts []CorporateCustomerContact `json:"contacts"`
}

// CustomerChangeResponse type.
type CustomerChangeResponse struct {
	Success bool   `json:"success"`
	ID      int    `json:"id,omitempty"`
	State   string `json:"state,omitempty"`
}

// CorporateCustomerChangeResponse type.
type CorporateCustomerChangeResponse CustomerChangeResponse

// CustomersUploadResponse type.
type CustomersUploadResponse struct {
	Success           bool              `json:"success"`
	UploadedCustomers []IdentifiersPair `json:"uploadedCustomers,omitempty"`
}

// CorporateCustomersUploadResponse type.
type CorporateCustomersUploadResponse CustomersUploadResponse

// CustomersHistoryResponse type.
type CustomersHistoryResponse struct {
	Success     bool                    `json:"success,omitempty"`
	GeneratedAt string                  `json:"generatedAt,omitempty"`
	History     []CustomerHistoryRecord `json:"history,omitempty"`
	Pagination  *Pagination             `json:"pagination,omitempty"`
}

// CorporateCustomersHistoryResponse type.
type CorporateCustomersHistoryResponse struct {
	Success     bool                             `json:"success,omitempty"`
	GeneratedAt string                           `json:"generatedAt,omitempty"`
	History     []CorporateCustomerHistoryRecord `json:"history,omitempty"`
	Pagination  *Pagination                      `json:"pagination,omitempty"`
}

// OrderResponse type.
type OrderResponse struct {
	Success bool   `json:"success"`
	Order   *Order `json:"order,omitempty"`
}

// OrdersResponse type.
type OrdersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Orders     []Order     `json:"orders,omitempty"`
}

// OrdersStatusesResponse type.
type OrdersStatusesResponse struct {
	Success bool           `json:"success"`
	Orders  []OrdersStatus `json:"orders"`
}

// OrdersUploadResponse type.
type OrdersUploadResponse struct {
	Success        bool              `json:"success"`
	UploadedOrders []IdentifiersPair `json:"uploadedOrders,omitempty"`
}

// OrdersHistoryResponse type.
type OrdersHistoryResponse struct {
	Success     bool                  `json:"success,omitempty"`
	GeneratedAt string                `json:"generatedAt,omitempty"`
	History     []OrdersHistoryRecord `json:"history,omitempty"`
	Pagination  *Pagination           `json:"pagination,omitempty"`
}

// PackResponse type.
type PackResponse struct {
	Success bool  `json:"success"`
	Pack    *Pack `json:"pack,omitempty"`
}

// PacksResponse type.
type PacksResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Packs      []Pack      `json:"packs,omitempty"`
}

// PacksHistoryResponse type.
type PacksHistoryResponse struct {
	Success     bool                 `json:"success,omitempty"`
	GeneratedAt string               `json:"generatedAt,omitempty"`
	History     []PacksHistoryRecord `json:"history,omitempty"`
	Pagination  *Pagination          `json:"pagination,omitempty"`
}

// UserResponse type.
type UserResponse struct {
	Success bool  `json:"success"`
	User    *User `json:"user,omitempty"`
}

// UsersResponse type.
type UsersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Users      []User      `json:"users,omitempty"`
}

// UserGroupsResponse type.
type UserGroupsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Groups     []UserGroup `json:"groups,omitempty"`
}

// TaskResponse type.
type TaskResponse struct {
	Success bool  `json:"success"`
	Task    *Task `json:"task,omitempty"`
}

// TasksResponse type.
type TasksResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Tasks      []Task      `json:"tasks,omitempty"`
}

// NotesResponse type.
type NotesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Notes      []Note      `json:"notes,omitempty"`
}

// SegmentsResponse type.
type SegmentsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Segments   []Segment   `json:"segments,omitempty"`
}

// SettingsResponse type.
type SettingsResponse struct {
	Success  bool     `json:"success"`
	Settings Settings `json:"settings,omitempty"`
}

// CountriesResponse type.
type CountriesResponse struct {
	Success      bool     `json:"success"`
	CountriesIso []string `json:"countriesIso,omitempty"`
}

// CostGroupsResponse type.
type CostGroupsResponse struct {
	Success    bool        `json:"success"`
	CostGroups []CostGroup `json:"costGroups,omitempty"`
}

// CostItemsResponse type.
type CostItemsResponse struct {
	Success   bool       `json:"success"`
	CostItems []CostItem `json:"costItems,omitempty"`
}

// CouriersResponse type.
type CouriersResponse struct {
	Success  bool      `json:"success"`
	Couriers []Courier `json:"couriers,omitempty"`
}

// DeliveryServiceResponse type.
type DeliveryServiceResponse struct {
	Success          bool                       `json:"success"`
	DeliveryServices map[string]DeliveryService `json:"deliveryServices,omitempty"`
}

// DeliveryTypesResponse type.
type DeliveryTypesResponse struct {
	Success       bool                    `json:"success"`
	DeliveryTypes map[string]DeliveryType `json:"deliveryTypes,omitempty"`
}

// LegalEntitiesResponse type.
type LegalEntitiesResponse struct {
	Success       bool          `json:"success"`
	LegalEntities []LegalEntity `json:"legalEntities,omitempty"`
}

// OrderMethodsResponse type.
type OrderMethodsResponse struct {
	Success      bool                   `json:"success"`
	OrderMethods map[string]OrderMethod `json:"orderMethods,omitempty"`
}

// OrderTypesResponse type.
type OrderTypesResponse struct {
	Success    bool                 `json:"success"`
	OrderTypes map[string]OrderType `json:"orderTypes,omitempty"`
}

// PaymentStatusesResponse type.
type PaymentStatusesResponse struct {
	Success         bool                     `json:"success"`
	PaymentStatuses map[string]PaymentStatus `json:"paymentStatuses,omitempty"`
}

// PaymentTypesResponse type.
type PaymentTypesResponse struct {
	Success      bool                   `json:"success"`
	PaymentTypes map[string]PaymentType `json:"paymentTypes,omitempty"`
}

// PriceTypesResponse type.
type PriceTypesResponse struct {
	Success    bool        `json:"success"`
	PriceTypes []PriceType `json:"priceTypes,omitempty"`
}

// ProductStatusesResponse type.
type ProductStatusesResponse struct {
	Success         bool                     `json:"success"`
	ProductStatuses map[string]ProductStatus `json:"productStatuses,omitempty"`
}

// StatusesResponse type.
type StatusesResponse struct {
	Success  bool              `json:"success"`
	Statuses map[string]Status `json:"statuses,omitempty"`
}

// StatusGroupsResponse type.
type StatusGroupsResponse struct {
	Success      bool                   `json:"success"`
	StatusGroups map[string]StatusGroup `json:"statusGroups,omitempty"`
}

// SitesResponse type.
type SitesResponse struct {
	Success bool            `json:"success"`
	Sites   map[string]Site `json:"sites,omitempty"`
}

// StoresResponse type.
type StoresResponse struct {
	Success bool    `json:"success"`
	Stores  []Store `json:"stores,omitempty"`
}

// InventoriesResponse type.
type InventoriesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Offers     []Offer     `json:"offers,omitempty"`
}

// StoreUploadResponse type.
type StoreUploadResponse struct {
	Success              bool    `json:"success"`
	ProcessedOffersCount int     `json:"processedOffersCount,omitempty"`
	NotFoundOffers       []Offer `json:"notFoundOffers,omitempty"`
}

// ProductsGroupsResponse type.
type ProductsGroupsResponse struct {
	Success      bool           `json:"success"`
	Pagination   *Pagination    `json:"pagination,omitempty"`
	ProductGroup []ProductGroup `json:"productGroup,omitempty"`
}

// ProductsResponse type.
type ProductsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Products   []Product   `json:"products,omitempty"`
}

// ProductsPropertiesResponse type.
type ProductsPropertiesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Properties []Property  `json:"properties,omitempty"`
}

// DeliveryShipmentsResponse type.
type DeliveryShipmentsResponse struct {
	Success           bool               `json:"success"`
	Pagination        *Pagination        `json:"pagination,omitempty"`
	DeliveryShipments []DeliveryShipment `json:"deliveryShipments,omitempty"`
}

// DeliveryShipmentResponse type.
type DeliveryShipmentResponse struct {
	Success          bool              `json:"success"`
	DeliveryShipment *DeliveryShipment `json:"deliveryShipment,omitempty"`
}

// DeliveryShipmentUpdateResponse type.
type DeliveryShipmentUpdateResponse struct {
	Success bool   `json:"success"`
	ID      int    `json:"id,omitempty"`
	Status  string `json:"status,omitempty"`
}

// IntegrationModuleResponse type.
type IntegrationModuleResponse struct {
	Success           bool               `json:"success"`
	IntegrationModule *IntegrationModule `json:"integrationModule,omitempty"`
}

// UpdateScopesResponse update scopes response.
type UpdateScopesResponse struct {
	ErrorResponse
	APIKey string `json:"apiKey"`
}

// IntegrationModuleEditResponse type.
type IntegrationModuleEditResponse struct {
	Success bool         `json:"success"`
	Info    ResponseInfo `json:"info,omitempty"`
}

// ResponseInfo type.
type ResponseInfo struct {
	MgTransportInfo MgInfo       `json:"mgTransport,omitempty"`
	MgBotInfo       MgInfo       `json:"mgBot,omitempty"`
	BillingInfo     *BillingInfo `json:"billingInfo,omitempty"`
}

type BillingInfo struct {
	Price             float64              `json:"price,omitempty"`
	PriceWithDiscount float64              `json:"priceWithDiscount,omitempty"`
	BillingType       string               `json:"billingType,omitempty"`
	Currency          *BillingInfoCurrency `json:"currency,omitempty"`
}

type BillingInfoCurrency struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	Code      string `json:"code,omitempty"`
}

// MgInfo type.
type MgInfo struct {
	EndpointURL string `json:"endpointUrl"`
	Token       string `json:"token"`
}

// CostsResponse type.
type CostsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Costs      []Cost      `json:"costs,omitempty"`
}

// CostsUploadResponse type.
type CostsUploadResponse struct {
	Success       bool  `json:"success"`
	UploadedCosts []int `json:"uploadedCosts,omitempty"`
}

// CostsDeleteResponse type.
type CostsDeleteResponse struct {
	Success       bool  `json:"success"`
	Count         int   `json:"count,omitempty"`
	NotRemovedIds []int `json:"notRemovedIds,omitempty"`
}

// CostResponse type.
type CostResponse struct {
	Success bool  `json:"success"`
	Cost    *Cost `json:"cost,omitempty"`
}

// FilesResponse type.
type FilesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Files      []File      `json:"files,omitempty"`
}

// FileUpload response.
type FileUploadResponse struct {
	Success bool  `json:"success"`
	File    *File `json:"file,omitempty"`
}

// FileResponse type.
type FileResponse struct {
	Success bool  `json:"success"`
	File    *File `json:"file,omitempty"`
}

// CustomFieldsResponse type.
type CustomFieldsResponse struct {
	Success      bool           `json:"success"`
	Pagination   *Pagination    `json:"pagination,omitempty"`
	CustomFields []CustomFields `json:"customFields,omitempty"`
}

// CustomDictionariesResponse type.
type CustomDictionariesResponse struct {
	Success            bool                `json:"success"`
	Pagination         *Pagination         `json:"pagination,omitempty"`
	CustomDictionaries *[]CustomDictionary `json:"customDictionaries,omitempty"`
}

// CustomResponse type.
type CustomResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
}

// CustomDictionaryResponse type.
type CustomDictionaryResponse struct {
	Success          bool              `json:"success"`
	CustomDictionary *CustomDictionary `json:"CustomDictionary,omitempty"`
}

// CustomFieldResponse type.
type CustomFieldResponse struct {
	Success     bool         `json:"success"`
	CustomField CustomFields `json:"customField,omitempty"`
}

// UnitsResponse type.
type UnitsResponse struct {
	Success bool    `json:"success"`
	Units   *[]Unit `json:"units,omitempty"`
}

// ErrorResponse should be returned to the one-step connection request in case of failure.
type ErrorResponse struct {
	SuccessfulResponse
	ErrorMessage string `json:"errorMsg"`
}

// ConnectResponse should be returned to the one-step connection request in case of successful connection.
type ConnectResponse struct {
	SuccessfulResponse
	AccountURL string `json:"accountUrl"`
}

// ConnectionConfigResponse contains connection configuration for one-step connection.
type ConnectionConfigResponse struct {
	SuccessfulResponse
	Scopes      []string `json:"scopes"`
	RegisterURL string   `json:"registerUrl"`
}

// NewConnectResponse returns ConnectResponse with the provided account URL.
func NewConnectResponse(accountURL string) ConnectResponse {
	return ConnectResponse{
		SuccessfulResponse: SuccessfulResponse{Success: true},
		AccountURL:         accountURL,
	}
}
