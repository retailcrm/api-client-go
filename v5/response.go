package v5

// SuccessfulResponse type
type SuccessfulResponse struct {
	Success bool `json:"success,omitempty"`
}

// CreateResponse type
type CreateResponse struct {
	Success bool `json:"success"`
	ID      int  `json:"id,omitempty"`
}

// OrderCreateResponse type
type OrderCreateResponse struct {
	CreateResponse
	Order Order `json:"order,omitempty"`
}

// OperationResponse type
type OperationResponse struct {
	Success bool              `json:"success"`
	Errors  map[string]string `json:"errors,omitempty,brackets"`
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

// CorporateCustomerResponse type
type CorporateCustomerResponse struct {
	Success           bool               `json:"success"`
	CorporateCustomer *CorporateCustomer `json:"customerCorporate,omitempty,brackets"`
}

// CustomersResponse type
type CustomersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Customers  []Customer  `json:"customers,omitempty,brackets"`
}

// CorporateCustomersResponse type
type CorporateCustomersResponse struct {
	Success            bool                `json:"success"`
	Pagination         *Pagination         `json:"pagination,omitempty"`
	CustomersCorporate []CorporateCustomer `json:"customersCorporate,omitempty,brackets"`
}

// CorporateCustomersNotesResponse type
type CorporateCustomersNotesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Notes      []Note      `json:"notes,omitempty,brackets"`
}

// CorporateCustomersAddressesResponse type
type CorporateCustomersAddressesResponse struct {
	Success   bool                       `json:"success"`
	Addresses []CorporateCustomerAddress `json:"addresses"`
}

// CorporateCustomerCompaniesResponse type
type CorporateCustomerCompaniesResponse struct {
	Success   bool      `json:"success"`
	Companies []Company `json:"companies"`
}

// CorporateCustomerContactsResponse type
type CorporateCustomerContactsResponse struct {
	Success  bool                       `json:"success"`
	Contacts []CorporateCustomerContact `json:"contacts"`
}

// CustomerChangeResponse type
type CustomerChangeResponse struct {
	Success bool   `json:"success"`
	ID      int    `json:"id,omitempty"`
	State   string `json:"state,omitempty"`
}

// CorporateCustomerChangeResponse type
type CorporateCustomerChangeResponse CustomerChangeResponse

// CustomersUploadResponse type
type CustomersUploadResponse struct {
	Success           bool              `json:"success"`
	UploadedCustomers []IdentifiersPair `json:"uploadedCustomers,omitempty,brackets"`
}

// CorporateCustomersUploadResponse type
type CorporateCustomersUploadResponse CustomersUploadResponse

// CustomersHistoryResponse type
type CustomersHistoryResponse struct {
	Success     bool                    `json:"success,omitempty"`
	GeneratedAt string                  `json:"generatedAt,omitempty"`
	History     []CustomerHistoryRecord `json:"history,omitempty,brackets"`
	Pagination  *Pagination             `json:"pagination,omitempty"`
}

// CorporateCustomersHistoryResponse type
type CorporateCustomersHistoryResponse struct {
	Success     bool                             `json:"success,omitempty"`
	GeneratedAt string                           `json:"generatedAt,omitempty"`
	History     []CorporateCustomerHistoryRecord `json:"history,omitempty,brackets"`
	Pagination  *Pagination                      `json:"pagination,omitempty"`
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

// OrdersStatusesResponse type
type OrdersStatusesResponse struct {
	Success bool           `json:"success"`
	Orders  []OrdersStatus `json:"orders"`
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

// PackResponse type
type PackResponse struct {
	Success bool  `json:"success"`
	Pack    *Pack `json:"pack,omitempty,brackets"`
}

// PacksResponse type
type PacksResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Packs      []Pack      `json:"packs,omitempty,brackets"`
}

// PacksHistoryResponse type
type PacksHistoryResponse struct {
	Success     bool                 `json:"success,omitempty"`
	GeneratedAt string               `json:"generatedAt,omitempty"`
	History     []PacksHistoryRecord `json:"history,omitempty,brackets"`
	Pagination  *Pagination          `json:"pagination,omitempty"`
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

// SettingsResponse type
type SettingsResponse struct {
	Success    bool        `json:"success"`
	Settings   Settings    `json:"settings,omitempty,brackets"`
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

// InventoriesResponse type
type InventoriesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Offers     []Offer     `json:"offers,omitempty"`
}

// StoreUploadResponse type
type StoreUploadResponse struct {
	Success              bool    `json:"success"`
	ProcessedOffersCount int     `json:"processedOffersCount,omitempty"`
	NotFoundOffers       []Offer `json:"notFoundOffers,omitempty"`
}

// ProductsGroupsResponse type
type ProductsGroupsResponse struct {
	Success      bool           `json:"success"`
	Pagination   *Pagination    `json:"pagination,omitempty"`
	ProductGroup []ProductGroup `json:"productGroup,omitempty,brackets"`
}

// ProductsResponse type
type ProductsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Products   []Product   `json:"products,omitempty,brackets"`
}

// ProductsPropertiesResponse type
type ProductsPropertiesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Properties []Property  `json:"properties,omitempty,brackets"`
}

// DeliveryShipmentsResponse type
type DeliveryShipmentsResponse struct {
	Success           bool               `json:"success"`
	Pagination        *Pagination        `json:"pagination,omitempty"`
	DeliveryShipments []DeliveryShipment `json:"deliveryShipments,omitempty,brackets"`
}

// DeliveryShipmentResponse type
type DeliveryShipmentResponse struct {
	Success          bool              `json:"success"`
	DeliveryShipment *DeliveryShipment `json:"deliveryShipment,omitempty,brackets"`
}

// DeliveryShipmentUpdateResponse type
type DeliveryShipmentUpdateResponse struct {
	Success bool   `json:"success"`
	ID      int    `json:"id,omitempty"`
	Status  string `json:"status,omitempty"`
}

// IntegrationModuleResponse type
type IntegrationModuleResponse struct {
	Success           bool               `json:"success"`
	IntegrationModule *IntegrationModule `json:"integrationModule,omitempty"`
}

// IntegrationModuleEditResponse type
type IntegrationModuleEditResponse struct {
	Success bool         `json:"success"`
	Info    ResponseInfo `json:"info,omitempty,brackets"`
}

// ResponseInfo type
type ResponseInfo struct {
	MgTransportInfo MgInfo `json:"mgTransport,omitempty,brackets"`
	MgBotInfo       MgInfo `json:"mgBot,omitempty,brackets"`
}

// MgInfo type
type MgInfo struct {
	EndpointUrl string `json:"endpointUrl"`
	Token       string `json:"token"`
}

// CostsResponse type
type CostsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Costs      []Cost      `json:"costs,omitempty,brackets"`
}

// CostsUploadResponse type
type CostsUploadResponse struct {
	Success       bool  `json:"success"`
	UploadedCosts []int `json:"uploadedCosts,omitempty,brackets"`
}

// CostsDeleteResponse type
type CostsDeleteResponse struct {
	Success       bool  `json:"success"`
	Count         int   `json:"count,omitempty,brackets"`
	NotRemovedIds []int `json:"notRemovedIds,omitempty,brackets"`
}

// CostResponse type
type CostResponse struct {
	Success bool  `json:"success"`
	Cost    *Cost `json:"cost,omitempty,brackets"`
}

// FilesResponse type
type FilesResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Files      []File      `json:"files,omitempty"`
}

// FileUpload response
type FileUploadResponse struct {
	Success bool  `json:"success"`
	File    *File `json:"file,omitempty"`
}

// FileResponse type
type FileResponse struct {
	Success bool  `json:"success"`
	File    *File `json:"file,omitempty"`
}

// CustomFieldsResponse type
type CustomFieldsResponse struct {
	Success      bool           `json:"success"`
	Pagination   *Pagination    `json:"pagination,omitempty"`
	CustomFields []CustomFields `json:"customFields,omitempty,brackets"`
}

// CustomDictionariesResponse type
type CustomDictionariesResponse struct {
	Success            bool                `json:"success"`
	Pagination         *Pagination         `json:"pagination,omitempty"`
	CustomDictionaries *[]CustomDictionary `json:"customDictionaries,omitempty,brackets"`
}

// CustomResponse type
type CustomResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
}

// CustomDictionaryResponse type
type CustomDictionaryResponse struct {
	Success          bool              `json:"success"`
	CustomDictionary *CustomDictionary `json:"CustomDictionary,omitempty,brackets"`
}

// CustomFieldResponse type
type CustomFieldResponse struct {
	Success     bool         `json:"success"`
	CustomField CustomFields `json:"customField,omitempty,brackets"`
}

// UnitsResponse type
type UnitsResponse struct {
	Success bool    `json:"success"`
	Units   *[]Unit `json:"units,omitempty,brackets"`
}
