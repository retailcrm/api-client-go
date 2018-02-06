package v5

import "net/http"

// Client type
type Client struct {
	Url        string
	Key        string
	httpClient *http.Client
}

// Pagination type
type Pagination struct {
	Limit          int `json:"limit,omitempty"`
	TotalCount     int `json:"totalCount,omitempty"`
	CurrentPage    int `json:"currentPage,omitempty"`
	TotalPageCount int `json:"totalPageCount,omitempty"`
}

// Address type
type Address struct {
	Index        string `json:"index,omitempty"`
	CountryIso   string `json:"countryIso,omitempty"`
	Region       string `json:"region,omitempty"`
	RegionId     int    `json:"regionId,omitempty"`
	City         string `json:"city,omitempty"`
	CityId       int    `json:"cityId,omitempty"`
	CityType     string `json:"cityType,omitempty"`
	Street       string `json:"street,omitempty"`
	StreetId     int    `json:"streetId,omitempty"`
	StreetType   string `json:"streetType,omitempty"`
	Building     string `json:"building,omitempty"`
	Flat         string `json:"flat,omitempty"`
	IntercomCode string `json:"intercomCode,omitempty"`
	Floor        int    `json:"floor,omitempty"`
	Block        int    `json:"block,omitempty"`
	House        string `json:"house,omitempty"`
	Metro        string `json:"metro,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Text         string `json:"text,omitempty"`
}

// GeoHierarchyRow type
type GeoHierarchyRow struct {
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	RegionId int    `json:"regionId,omitempty"`
	City     string `json:"city,omitempty"`
	CityId   int    `json:"cityId,omitempty"`
}

// Source type
type Source struct {
	Source   string `json:"source,omitempty"`
	Medium   string `json:"medium,omitempty"`
	Campaign string `json:"campaign,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	Content  string `json:"content,omitempty"`
}

// Contragent type
type Contragent struct {
	ContragentType    string `json:"contragentType,omitempty"`
	LegalName         string `json:"legalName,omitempty"`
	LegalAddress      string `json:"legalAddress,omitempty"`
	INN               string `json:"INN,omitempty"`
	OKPO              string `json:"OKPO,omitempty"`
	KPP               string `json:"KPP,omitempty"`
	OGRN              string `json:"OGRN,omitempty"`
	OGRNIP            string `json:"OGRNIP,omitempty"`
	CertificateNumber string `json:"certificateNumber,omitempty"`
	CertificateDate   string `json:"certificateDate,omitempty"`
	BIK               string `json:"BIK,omitempty"`
	Bank              string `json:"bank,omitempty"`
	BankAddress       string `json:"bankAddress,omitempty"`
	CorrAccount       string `json:"corrAccount,omitempty"`
	BankAccount       string `json:"bankAccount,omitempty"`
}

// ApiKey type
type ApiKey struct {
	Current bool `json:"current,omitempty"`
}

// Property type
type Property struct {
	Code  string `json:"code,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// IdentifiersPair type
type IdentifiersPair struct {
	Id         int    `json:"id,omitempty"`
	ExternalId string `json:"externalId,omitempty"`
}

/**
Customer related types
*/

// Customer type
type Customer struct {
	Id                           int         `json:"id,omitempty"`
	ExternalId                   string      `json:"externalId,omitempty"`
	FirstName                    string      `json:"firstName,omitempty"`
	LastName                     string      `json:"lastName,omitempty"`
	Patronymic                   string      `json:"patronymic,omitempty"`
	Sex                          string      `json:"sex,omitempty"`
	Email                        string      `json:"email,omitempty"`
	Phones                       []Phone     `json:"phones,brackets,omitempty"`
	Address                      *Address    `json:"address,omitempty"`
	CreatedAt                    string      `json:"createdAt,omitempty"`
	Birthday                     string      `json:"birthday,omitempty"`
	ManagerId                    int         `json:"managerId,omitempty"`
	Vip                          bool        `json:"vip,omitempty"`
	Bad                          bool        `json:"bad,omitempty"`
	Site                         string      `json:"site,omitempty"`
	Source                       *Source     `json:"source,omitempty"`
	Contragent                   *Contragent `json:"contragent,omitempty"`
	PersonalDiscount             float32     `json:"personalDiscount,omitempty"`
	CumulativeDiscount           float32     `json:"cumulativeDiscount,omitempty"`
	DiscountCardNumber           string      `json:"discountCardNumber,omitempty"`
	EmailMarketingUnsubscribedAt string      `json:"emailMarketingUnsubscribedAt,omitempty"`
	AvgMarginSumm                float32     `json:"avgMarginSumm,omitempty"`
	MarginSumm                   float32     `json:"marginSumm,omitempty"`
	TotalSumm                    float32     `json:"totalSumm,omitempty"`
	AverageSumm                  float32     `json:"averageSumm,omitempty"`
	OrdersCount                  int         `json:"ordersCount,omitempty"`
	CostSumm                     float32     `json:"costSumm,omitempty"`
	MaturationTime               int         `json:"maturationTime,omitempty"`
	FirstClientId                string      `json:"firstClientId,omitempty"`
	LastClientId                 string      `json:"lastClientId,omitempty"`
	BrowserId                    string      `json:"browserId,omitempty"`
	// CustomFields                 []map[string]string `json:"customFields,omitempty,brackets"`
}

// Phone type
type Phone struct {
	Number string `json:"number,omitempty"`
}

// CustomerHistoryRecord type
type CustomerHistoryRecord struct {
	Id        int       `json:"id,omitempty"`
	CreatedAt string    `json:"createdAt,omitempty"`
	Created   bool      `json:"created,omitempty"`
	Deleted   bool      `json:"deleted,omitempty"`
	Source    string    `json:"source,omitempty"`
	Field     string    `json:"field,omitempty"`
	User      *User     `json:"user,omitempty,brackets"`
	ApiKey    *ApiKey   `json:"apiKey,omitempty,brackets"`
	Customer  *Customer `json:"customer,omitempty,brackets"`
}

/**
Order related types
*/

// Order type
type Order struct {
	Id                            int                 `json:"id,omitempty"`
	ExternalId                    string              `json:"externalId,omitempty"`
	Number                        string              `json:"number,omitempty"`
	FirstName                     string              `json:"firstName,omitempty"`
	LastName                      string              `json:"lastName,omitempty"`
	Patronymic                    string              `json:"patronymic,omitempty"`
	Email                         string              `json:"email,omitempty"`
	Phone                         string              `json:"phone,omitempty"`
	AdditionalPhone               string              `json:"additionalPhone,omitempty"`
	CreatedAt                     string              `json:"createdAt,omitempty"`
	StatusUpdatedAt               string              `json:"statusUpdatedAt,omitempty"`
	ManagerId                     int                 `json:"managerId,omitempty"`
	Mark                          int                 `json:"mark,omitempty"`
	Call                          bool                `json:"call,omitempty"`
	Expired                       bool                `json:"expired,omitempty"`
	FromApi                       bool                `json:"fromApi,omitempty"`
	MarkDatetime                  string              `json:"markDatetime,omitempty"`
	CustomerComment               string              `json:"customerComment,omitempty"`
	ManagerComment                string              `json:"managerComment,omitempty"`
	Status                        string              `json:"status,omitempty"`
	StatusComment                 string              `json:"statusComment,omitempty"`
	FullPaidAt                    string              `json:"fullPaidAt,omitempty"`
	Site                          string              `json:"site,omitempty"`
	OrderType                     string              `json:"orderType,omitempty"`
	OrderMethod                   string              `json:"orderMethod,omitempty"`
	CountryIso                    string              `json:"countryIso,omitempty"`
	Summ                          float32             `json:"summ,omitempty"`
	TotalSumm                     float32             `json:"totalSumm,omitempty"`
	PrepaySum                     float32             `json:"prepaySum,omitempty"`
	PurchaseSumm                  float32             `json:"purchaseSumm,omitempty"`
	DiscountManualAmount          float32             `json:"discountManualAmount,omitempty"`
	DiscountManualPercent         float32             `json:"discountManualPercent,omitempty"`
	Weight                        float32             `json:"weight,omitempty"`
	Length                        int                 `json:"length,omitempty"`
	Width                         int                 `json:"width,omitempty"`
	Height                        int                 `json:"height,omitempty"`
	ShipmentStore                 string              `json:"shipmentStore,omitempty"`
	ShipmentDate                  string              `json:"shipmentDate,omitempty"`
	ClientId                      string              `json:"clientId,omitempty"`
	Shipped                       bool                `json:"shipped,omitempty"`
	UploadedToExternalStoreSystem bool                `json:"uploadedToExternalStoreSystem,omitempty"`
	Source                        *Source             `json:"source,omitempty"`
	Contragent                    *Contragent         `json:"contragent,omitempty"`
	Customer                      *Customer           `json:"customer,omitempty"`
	Delivery                      *OrderDelivery      `json:"delivery,omitempty"`
	Marketplace                   *OrderMarketplace   `json:"marketplace,omitempty"`
	Items                         []OrderItem         `json:"items,omitempty,brackets"`
	CustomFields                  []map[string]string `json:"customFields,omitempty,brackets"`
	// Payments                      []OrderPayment    `json:"payments,omitempty,brackets"`
}

// OrderDelivery type
type OrderDelivery struct {
	Code            string                `json:"code,omitempty"`
	IntegrationCode string                `json:"integrationCode,omitempty"`
	Cost            float32               `json:"cost,omitempty"`
	NetCost         float32               `json:"netCost,omitempty"`
	VatRate         string                `json:"vatRate,omitempty"`
	Date            string                `json:"date,omitempty"`
	Time            *OrderDeliveryTime    `json:"time,omitempty"`
	Address         *Address              `json:"address,omitempty"`
	Service         *OrderDeliveryService `json:"service,omitempty"`
	Data            *OrderDeliveryData    `json:"data,omitempty"`
}

// OrderDeliveryTime type
type OrderDeliveryTime struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Custom string `json:"custom,omitempty"`
}

// OrderDeliveryService type
type OrderDeliveryService struct {
	Name   string `json:"name,omitempty"`
	Code   string `json:"code,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// OrderDeliveryData type
type OrderDeliveryData struct {
	TrackNumber        string `json:"trackNumber,omitempty"`
	Status             string `json:"status,omitempty"`
	PickuppointAddress string `json:"pickuppointAddress,omitempty"`
	PayerType          string `json:"payerType,omitempty"`
}

// OrderMarketplace type
type OrderMarketplace struct {
	Code    string `json:"code,omitempty"`
	OrderId string `json:"orderId,omitempty"`
}

// OrderPayment type
type OrderPayment struct {
	Id         int     `json:"id,omitempty"`
	ExternalId string  `json:"externalId,omitempty"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
	PaidAt     string  `json:"paidAt,omitempty"`
	Amount     float32 `json:"amount,omitempty"`
	Comment    string  `json:"comment,omitempty"`
}

// OrderItem type
type OrderItem struct {
	Id                    int         `json:"id,omitempty"`
	InitialPrice          float32     `json:"initialPrice,omitempty"`
	PurchasePrice         float32     `json:"purchasePrice,omitempty"`
	DiscountTotal         float32     `json:"discountTotal,omitempty"`
	DiscountManualAmount  float32     `json:"discountManualAmount,omitempty"`
	DiscountManualPercent float32     `json:"discountManualPercent,omitempty"`
	ProductName           string      `json:"productName,omitempty"`
	VatRate               string      `json:"vatRate,omitempty"`
	CreatedAt             string      `json:"createdAt,omitempty"`
	Quantity              float32     `json:"quantity,omitempty"`
	Status                string      `json:"status,omitempty"`
	Comment               string      `json:"comment,omitempty"`
	IsCanceled            bool        `json:"isCanceled,omitempty"`
	Offer                 Offer       `json:"offer,omitempty"`
	Properties            []*Property `json:"properties,omitempty,brackets"`
	PriceType             *PriceType  `json:"priceType,omitempty"`
}

// OrdersHistoryRecord type
type OrdersHistoryRecord struct {
	Id        int     `json:"id,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	Created   bool    `json:"created,omitempty"`
	Deleted   bool    `json:"deleted,omitempty"`
	Source    string  `json:"source,omitempty"`
	Field     string  `json:"field,omitempty"`
	User      *User   `json:"user,omitempty,brackets"`
	ApiKey    *ApiKey `json:"apiKey,omitempty,brackets"`
	Order     *Order  `json:"order,omitempty,brackets"`
}

// Pack type
type Pack struct {
	Id                 int      `json:"id,omitempty"`
	PurchasePrice      float32  `json:"purchasePrice,omitempty"`
	Quantity           float32  `json:"quantity,omitempty"`
	Store              string   `json:"store,omitempty"`
	ShipmentDate       string   `json:"shipmentDate,omitempty"`
	InvoiceNumber      string   `json:"invoiceNumber,omitempty"`
	DeliveryNoteNumber string   `json:"deliveryNoteNumber,omitempty"`
	Item               PackItem `json:"item,omitempty"`
	ItemId             int      `json:"itemId,omitempty"`
}

// PackItem type
type PackItem struct {
	Id    int    `json:"id,omitempty"`
	Order *Order `json:"order,omitempty"`
	Offer *Offer `json:"offer,omitempty"`
}

// PacksHistoryRecord type
type PacksHistoryRecord struct {
	Id        int    `json:"id,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Created   bool   `json:"created,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
	Source    string `json:"source,omitempty"`
	Field     string `json:"field,omitempty"`
	User      *User  `json:"user,omitempty,brackets"`
	Pack      *Pack  `json:"pack,omitempty,brackets"`
}

// Offer type
type Offer struct {
	Id            int                 `json:"id,omitempty"`
	ExternalId    string              `json:"externalId,omitempty"`
	Name          string              `json:"name,omitempty"`
	XmlId         string              `json:"xmlId,omitempty"`
	Article       string              `json:"article,omitempty"`
	VatRate       string              `json:"vatRate,omitempty"`
	Price         float32             `json:"price,omitempty"`
	PurchasePrice float32             `json:"purchasePrice,omitempty"`
	Quantity      float32             `json:"quantity,omitempty"`
	Height        float32             `json:"height,omitempty"`
	Width         float32             `json:"width,omitempty"`
	Length        float32             `json:"length,omitempty"`
	Weight        float32             `json:"weight,omitempty"`
	Stores        []Inventory         `json:"stores,omitempty,brackets"`
	Properties    []map[string]string `json:"properties,omitempty,brackets"`
	Prices        []OfferPrice        `json:"prices,omitempty,brackets"`
	Images        []string            `json:"images,omitempty,brackets"`
}

// Inventory type
type Inventory struct {
	PurchasePrice float32 `json:"purchasePrice,omitempty"`
	Quantity      float32 `json:"quantity,omitempty"`
	Store         string  `json:"store,omitempty"`
}

// InventoryUpload type
type InventoryUpload struct {
	Id         int                    `json:"id,omitempty"`
	ExternalId string                 `json:"externalId,omitempty"`
	XmlId      string                 `json:"xmlId,omitempty"`
	Stores     []InventoryUploadStore `json:"stores,omitempty"`
}

// InventoryUploadStore type
type InventoryUploadStore struct {
	PurchasePrice float32 `json:"purchasePrice,omitempty"`
	Available     float32 `json:"available,omitempty"`
	Code          string  `json:"code,omitempty"`
}

// OfferPrice type
type OfferPrice struct {
	Price     float32 `json:"price,omitempty"`
	Ordering  int     `json:"ordering,omitempty"`
	PriceType string  `json:"priceType,omitempty"`
}

/**
User related types
*/

// User type
type User struct {
	Id         int         `json:"id,omitempty"`
	FirstName  string      `json:"firstName,omitempty"`
	LastName   string      `json:"lastName,omitempty"`
	Patronymic string      `json:"patronymic,omitempty"`
	CreatedAt  string      `json:"createdAt,omitempty"`
	Active     bool        `json:"active,omitempty"`
	Online     bool        `json:"online,omitempty"`
	IsAdmin    bool        `json:"isAdmin,omitempty"`
	IsManager  bool        `json:"isManager,omitempty"`
	Email      string      `json:"email,omitempty"`
	Phone      string      `json:"phone,omitempty"`
	Status     string      `json:"status,omitempty"`
	Groups     []UserGroup `json:"groups,omitempty,brackets"`
}

// UserGroup type
type UserGroup struct {
	Name                  string   `json:"name,omitempty"`
	Code                  string   `json:"code,omitempty"`
	SignatureTemplate     string   `json:"signatureTemplate,omitempty"`
	IsManager             bool     `json:"isManager,omitempty"`
	IsDeliveryMen         bool     `json:"isDeliveryMen,omitempty"`
	DeliveryTypes         []string `json:"deliveryTypes,omitempty,brackets"`
	BreakdownOrderTypes   []string `json:"breakdownOrderTypes,omitempty,brackets"`
	BreakdownSites        []string `json:"breakdownSites,omitempty,brackets"`
	BreakdownOrderMethods []string `json:"breakdownOrderMethods,omitempty,brackets"`
	GrantedOrderTypes     []string `json:"grantedOrderTypes,omitempty,brackets"`
	GrantedSites          []string `json:"grantedSites,omitempty,brackets"`
}

/**
Task related types
*/

// Task type
type Task struct {
	Id          int       `json:"id,omitempty"`
	PerformerId int       `json:"performerId,omitempty"`
	Text        string    `json:"text,omitempty"`
	Commentary  string    `json:"commentary,omitempty"`
	Datetime    string    `json:"datetime,omitempty"`
	Complete    bool      `json:"complete,omitempty"`
	CreatedAt   string    `json:"createdAt,omitempty"`
	Creator     int       `json:"creator,omitempty"`
	Performer   int       `json:"performer,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	PhoneSite   string    `json:"phoneSite,omitempty"`
	Customer    *Customer `json:"customer,omitempty"`
	Order       *Order    `json:"order,omitempty"`
}

/*
	Notes related types
*/

// Note type
type Note struct {
	Id        int       `json:"id,omitempty"`
	ManagerId int       `json:"managerId,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt string    `json:"createdAt,omitempty"`
	Customer  *Customer `json:"customer,omitempty"`
}

/*
	Payments related types
*/

// Payment type
type Payment struct {
	Id         int     `json:"id,omitempty"`
	ExternalId string  `json:"externalId,omitempty"`
	PaidAt     string  `json:"paidAt,omitempty"`
	Amount     float32 `json:"amount,omitempty"`
	Comment    string  `json:"comment,omitempty"`
	Status     string  `json:"status,omitempty"`
	Type       string  `json:"type,omitempty"`
	Order      *Order  `json:"order,omitempty"`
}

/*
	Segment related types
*/

// Segment type
type Segment struct {
	Id             int    `json:"id,omitempty"`
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	CustomersCount int    `json:"customersCount,omitempty"`
	IsDynamic      bool   `json:"isDynamic,omitempty"`
	Active         bool   `json:"active,omitempty"`
}

/**
Reference related types
*/

// CostGroup type
type CostGroup struct {
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Color    string `json:"color,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Ordering int    `json:"ordering,omitempty"`
}

// CostItem type
type CostItem struct {
	Name            string  `json:"name,omitempty"`
	Code            string  `json:"code,omitempty"`
	Group           string  `json:"group,omitempty"`
	Type            string  `json:"type,omitempty"`
	Active          bool    `json:"active,omitempty"`
	AppliesToOrders bool    `json:"appliesToOrders,omitempty"`
	AppliesToUsers  bool    `json:"appliesToUsers,omitempty"`
	Ordering        int     `json:"ordering,omitempty"`
	Source          *Source `json:"source,omitempty"`
}

// Courier type
type Courier struct {
	Id          int    `json:"id,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Email       string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active,omitempty"`
	Phone       *Phone `json:"phone,omitempty"`
}

// DeliveryService type
type DeliveryService struct {
	Name   string `json:"name,omitempty"`
	Code   string `json:"code,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// DeliveryType type
type DeliveryType struct {
	Name             string   `json:"name,omitempty"`
	Code             string   `json:"code,omitempty"`
	Active           bool     `json:"active,omitempty"`
	DefaultCost      float32  `json:"defaultCost,omitempty"`
	DefaultNetCost   float32  `json:"defaultNetCost,omitempty"`
	Description      string   `json:"description,omitempty"`
	IntegrationCode  string   `json:"integrationCode,omitempty"`
	VatRate          string   `json:"vatRate,omitempty"`
	DefaultForCrm    bool     `json:"defaultForCrm,omitempty"`
	DeliveryServices []string `json:"deliveryServices,omitempty"`
	PaymentTypes     []string `json:"paymentTypes,omitempty"`
}

// LegalEntity type
type LegalEntity struct {
	Code              string `json:"code,omitempty"`
	VatRate           string `json:"vatRate,omitempty"`
	CountryIso        string `json:"countryIso,omitempty"`
	ContragentType    string `json:"contragentType,omitempty"`
	LegalName         string `json:"legalName,omitempty"`
	LegalAddress      string `json:"legalAddress,omitempty"`
	INN               string `json:"INN,omitempty"`
	OKPO              string `json:"OKPO,omitempty"`
	KPP               string `json:"KPP,omitempty"`
	OGRN              string `json:"OGRN,omitempty"`
	OGRNIP            string `json:"OGRNIP,omitempty"`
	CertificateNumber string `json:"certificateNumber,omitempty"`
	CertificateDate   string `json:"certificateDate,omitempty"`
	BIK               string `json:"BIK,omitempty"`
	Bank              string `json:"bank,omitempty"`
	BankAddress       string `json:"bankAddress,omitempty"`
	CorrAccount       string `json:"corrAccount,omitempty"`
	BankAccount       string `json:"bankAccount,omitempty"`
}

// OrderMethod type
type OrderMethod struct {
	Name          string `json:"name,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        bool   `json:"active,omitempty"`
	DefaultForCrm bool   `json:"defaultForCrm,omitempty"`
	DefaultForApi bool   `json:"defaultForApi,omitempty"`
}

// OrderType type
type OrderType struct {
	Name          string `json:"name,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        bool   `json:"active,omitempty"`
	DefaultForCrm bool   `json:"defaultForCrm,omitempty"`
	DefaultForApi bool   `json:"defaultForApi,omitempty"`
}

// PaymentStatus type
type PaymentStatus struct {
	Name            string   `json:"name,omitempty"`
	Code            string   `json:"code,omitempty"`
	Active          bool     `json:"active,omitempty"`
	DefaultForCrm   bool     `json:"defaultForCrm,omitempty"`
	DefaultForApi   bool     `json:"defaultForApi,omitempty"`
	PaymentComplete bool     `json:"paymentComplete,omitempty"`
	Description     string   `json:"description,omitempty"`
	Ordering        int      `json:"ordering,omitempty"`
	PaymentTypes    []string `json:"paymentTypes,omitempty,brackets"`
}

// PaymentType type
type PaymentType struct {
	Name            string   `json:"name,omitempty"`
	Code            string   `json:"code,omitempty"`
	Active          bool     `json:"active,omitempty"`
	DefaultForCrm   bool     `json:"defaultForCrm,omitempty"`
	DefaultForApi   bool     `json:"defaultForApi,omitempty"`
	Description     string   `json:"description,omitempty"`
	DeliveryTypes   []string `json:"deliveryTypes,omitempty,brackets"`
	PaymentStatuses []string `json:"PaymentStatuses,omitempty,brackets"`
}

// PriceType type
type PriceType struct {
	Id               int               `json:"id,omitempty"`
	Code             string            `json:"code,omitempty"`
	Name             string            `json:"name,omitempty"`
	Active           bool              `json:"active,omitempty"`
	Default          bool              `json:"default,omitempty"`
	Description      string            `json:"description,omitempty"`
	FilterExpression string            `json:"filterExpression,omitempty"`
	Ordering         int               `json:"ordering,omitempty"`
	Groups           []string          `json:"groups,omitempty,brackets"`
	Geo              []GeoHierarchyRow `json:"geo,omitempty,brackets"`
}

// ProductStatus type
type ProductStatus struct {
	Name                        string `json:"name,omitempty"`
	Code                        string `json:"code,omitempty"`
	Active                      bool   `json:"active,omitempty"`
	Ordering                    int    `json:"ordering,omitempty"`
	CreatedAt                   string `json:"createdAt,omitempty"`
	CancelStatus                bool   `json:"cancelStatus,omitempty"`
	OrderStatusByProductStatus  string `json:"orderStatusByProductStatus,omitempty"`
	OrderStatusForProductStatus string `json:"orderStatusForProductStatus,omitempty"`
}

// Status type
type Status struct {
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Ordering int    `json:"ordering,omitempty"`
	Group    string `json:"group,omitempty"`
}

// StatusGroup type
type StatusGroup struct {
	Name     string   `json:"name,omitempty"`
	Code     string   `json:"code,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Ordering int      `json:"ordering,omitempty"`
	Process  bool     `json:"process,omitempty"`
	Statuses []string `json:"statuses,omitempty,brackets"`
}

// Site type
type Site struct {
	Name             string       `json:"name,omitempty"`
	Code             string       `json:"code,omitempty"`
	Url              string       `json:"url,omitempty"`
	Description      string       `json:"description,omitempty"`
	Phones           string       `json:"phones,omitempty"`
	Zip              string       `json:"zip,omitempty"`
	Address          string       `json:"address,omitempty"`
	CountryIso       string       `json:"countryIso,omitempty"`
	YmlUrl           string       `json:"ymlUrl,omitempty"`
	LoadFromYml      bool         `json:"loadFromYml,omitempty"`
	CatalogUpdatedAt string       `json:"catalogUpdatedAt,omitempty"`
	CatalogLoadingAt string       `json:"catalogLoadingAt,omitempty"`
	Contragent       *LegalEntity `json:"contragent,omitempty"`
}

// Store type
type Store struct {
	Name          string   `json:"name,omitempty"`
	Code          string   `json:"code,omitempty"`
	ExternalId    string   `json:"externalId,omitempty"`
	Description   string   `json:"description,omitempty"`
	XmlId         string   `json:"xmlId,omitempty"`
	Email         string   `json:"email,omitempty"`
	Type          string   `json:"type,omitempty"`
	InventoryType string   `json:"inventoryType,omitempty"`
	Active        bool     `json:"active,omitempty"`
	Phone         *Phone   `json:"phone,omitempty"`
	Address       *Address `json:"address,omitempty"`
}
