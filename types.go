package retailcrm

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

// ByID is "id" constant to use as `by` property in methods.
const ByID = "id"

// ByExternalId is "externalId" constant to use as `by` property in methods.
const ByExternalID = "externalId"

// Client type.
type Client struct {
	URL        string
	Key        string
	Debug      bool
	httpClient *http.Client
	logger     BasicLogger
}

// Pagination type.
type Pagination struct {
	Limit          int `json:"limit,omitempty"`
	TotalCount     int `json:"totalCount,omitempty"`
	CurrentPage    int `json:"currentPage,omitempty"`
	TotalPageCount int `json:"totalPageCount,omitempty"`
}

// Address type.
type Address struct {
	Index      string `json:"index,omitempty"`
	CountryIso string `json:"countryIso,omitempty"`
	Region     string `json:"region,omitempty"`
	RegionID   int    `json:"regionId,omitempty"`
	City       string `json:"city,omitempty"`
	CityID     int    `json:"cityId,omitempty"`
	CityType   string `json:"cityType,omitempty"`
	Street     string `json:"street,omitempty"`
	StreetID   int    `json:"streetId,omitempty"`
	StreetType string `json:"streetType,omitempty"`
	Building   string `json:"building,omitempty"`
	Flat       string `json:"flat,omitempty"`
	Floor      int    `json:"floor,omitempty"`
	Block      int    `json:"block,omitempty"`
	House      string `json:"house,omitempty"`
	Housing    string `json:"housing,omitempty"`
	Metro      string `json:"metro,omitempty"`
	Notes      string `json:"notes,omitempty"`
	Text       string `json:"text,omitempty"`
}

// GeoHierarchyRow type.
type GeoHierarchyRow struct {
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	RegionID int    `json:"regionId,omitempty"`
	City     string `json:"city,omitempty"`
	CityID   int    `json:"cityId,omitempty"`
}

// Source type.
type Source struct {
	Source   string `json:"source,omitempty"`
	Medium   string `json:"medium,omitempty"`
	Campaign string `json:"campaign,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	Content  string `json:"content,omitempty"`
}

// Contragent type.
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

// APIKey type.
type APIKey struct {
	Current bool `json:"current,omitempty"`
}

// Property type.
type Property struct {
	Code  string   `json:"code,omitempty"`
	Name  string   `json:"name,omitempty"`
	Value string   `json:"value,omitempty"`
	Sites []string `json:"Sites,omitempty"`
}

// IdentifiersPair type.
type IdentifiersPair struct {
	ID         int    `json:"id,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}

// DeliveryTime type.
type DeliveryTime struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Custom string `json:"custom,omitempty"`
}

/**
Customer related types
*/

// Customer type.
type Customer struct {
	ID                           int         `json:"id,omitempty"`
	ExternalID                   string      `json:"externalId,omitempty"`
	FirstName                    string      `json:"firstName,omitempty"`
	LastName                     string      `json:"lastName,omitempty"`
	Patronymic                   string      `json:"patronymic,omitempty"`
	Sex                          string      `json:"sex,omitempty"`
	Email                        string      `json:"email,omitempty"`
	Phones                       []Phone     `json:"phones,omitempty"`
	Address                      *Address    `json:"address,omitempty"`
	CreatedAt                    string      `json:"createdAt,omitempty"`
	Birthday                     string      `json:"birthday,omitempty"`
	ManagerID                    int         `json:"managerId,omitempty"`
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
	FirstClientID                string      `json:"firstClientId,omitempty"`
	LastClientID                 string      `json:"lastClientId,omitempty"`
	BrowserID                    string      `json:"browserId,omitempty"`
	MgCustomerID                 string      `json:"mgCustomerId,omitempty"`
	PhotoURL                     string      `json:"photoUrl,omitempty"`
	CustomFields                 StringMap   `json:"customFields,omitempty"`
	Tags                         []Tag       `json:"tags,omitempty"`
}

// CorporateCustomer type.
type CorporateCustomer struct {
	ID                 int                        `json:"id,omitempty"`
	ExternalID         string                     `json:"externalId,omitempty"`
	Nickname           string                     `json:"nickName,omitempty"`
	CreatedAt          string                     `json:"createdAt,omitempty"`
	Vip                bool                       `json:"vip,omitempty"`
	Bad                bool                       `json:"bad,omitempty"`
	CustomFields       StringMap                  `json:"customFields,omitempty"`
	PersonalDiscount   float32                    `json:"personalDiscount,omitempty"`
	DiscountCardNumber string                     `json:"discountCardNumber,omitempty"`
	ManagerID          int                        `json:"managerId,omitempty"`
	Source             *Source                    `json:"source,omitempty"`
	CustomerContacts   []CorporateCustomerContact `json:"customerContacts,omitempty"`
	Companies          []Company                  `json:"companies,omitempty"`
	Addresses          []CorporateCustomerAddress `json:"addresses,omitempty"`
}

type CorporateCustomerContact struct {
	IsMain    bool                             `json:"isMain,omitempty"`
	Customer  CorporateCustomerContactCustomer `json:"customer,omitempty"`
	Companies []IdentifiersPair                `json:"companies,omitempty"`
}

// CorporateCustomerAddress type. Address didn't inherited in order to simplify declaration.
type CorporateCustomerAddress struct {
	ID           int    `json:"id,omitempty"`
	Index        string `json:"index,omitempty"`
	CountryISO   string `json:"countryIso,omitempty"`
	Region       string `json:"region,omitempty"`
	RegionID     int    `json:"regionId,omitempty"`
	City         string `json:"city,omitempty"`
	CityID       int    `json:"cityId,omitempty"`
	CityType     string `json:"cityType,omitempty"`
	Street       string `json:"street,omitempty"`
	StreetID     int    `json:"streetId,omitempty"`
	StreetType   string `json:"streetType,omitempty"`
	Building     string `json:"building,omitempty"`
	Flat         string `json:"flat,omitempty"`
	IntercomCode string `json:"intercomCode,omitempty"`
	Floor        int    `json:"floor,omitempty"`
	Block        int    `json:"block,omitempty"`
	House        string `json:"house,omitempty"`
	Housing      string `json:"housing,omitempty"`
	Metro        string `json:"metro,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Text         string `json:"text,omitempty"`
	ExternalID   string `json:"externalId,omitempty"`
	Name         string `json:"name,omitempty"`
}

type CorporateCustomerContactCustomer struct {
	ID         int    `json:"id,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
	BrowserID  string `json:"browserId,omitempty"`
	Site       string `json:"site,omitempty"`
}

type Company struct {
	ID           int              `json:"id,omitempty"`
	IsMain       bool             `json:"isMain,omitempty"`
	ExternalID   string           `json:"externalId,omitempty"`
	Active       bool             `json:"active,omitempty"`
	Name         string           `json:"name,omitempty"`
	Brand        string           `json:"brand,omitempty"`
	Site         string           `json:"site,omitempty"`
	CreatedAt    string           `json:"createdAt,omitempty"`
	Contragent   *Contragent      `json:"contragent,omitempty"`
	Address      *IdentifiersPair `json:"address,omitempty"`
	CustomFields StringMap        `json:"customFields,omitempty"`
}

// CorporateCustomerNote type.
type CorporateCustomerNote struct {
	ManagerID int              `json:"managerId,omitempty"`
	Text      string           `json:"text,omitempty"`
	Customer  *IdentifiersPair `json:"customer,omitempty"`
}

// Phone type.
type Phone struct {
	Number string `json:"number,omitempty"`
}

// CustomerHistoryRecord type.
type CustomerHistoryRecord struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt string    `json:"createdAt,omitempty"`
	Created   bool      `json:"created,omitempty"`
	Deleted   bool      `json:"deleted,omitempty"`
	Source    string    `json:"source,omitempty"`
	Field     string    `json:"field,omitempty"`
	User      *User     `json:"user,omitempty"`
	APIKey    *APIKey   `json:"apiKey,omitempty"`
	Customer  *Customer `json:"customer,omitempty"`
}

// CorporateCustomerHistoryRecord type.
type CorporateCustomerHistoryRecord struct {
	ID                int                `json:"id,omitempty"`
	CreatedAt         string             `json:"createdAt,omitempty"`
	Created           bool               `json:"created,omitempty"`
	Deleted           bool               `json:"deleted,omitempty"`
	Source            string             `json:"source,omitempty"`
	Field             string             `json:"field,omitempty"`
	User              *User              `json:"user,omitempty"`
	APIKey            *APIKey            `json:"apiKey,omitempty"`
	CorporateCustomer *CorporateCustomer `json:"corporateCustomer,omitempty"`
}

/**
Order related types
*/

type OrderPayments map[string]OrderPayment
type StringMap map[string]string
type Properties map[string]Property

// Order type.
type Order struct {
	ID                            int               `json:"id,omitempty"`
	ExternalID                    string            `json:"externalId,omitempty"`
	Number                        string            `json:"number,omitempty"`
	FirstName                     string            `json:"firstName,omitempty"`
	LastName                      string            `json:"lastName,omitempty"`
	Patronymic                    string            `json:"patronymic,omitempty"`
	Email                         string            `json:"email,omitempty"`
	Phone                         string            `json:"phone,omitempty"`
	AdditionalPhone               string            `json:"additionalPhone,omitempty"`
	CreatedAt                     string            `json:"createdAt,omitempty"`
	StatusUpdatedAt               string            `json:"statusUpdatedAt,omitempty"`
	ManagerID                     int               `json:"managerId,omitempty"`
	Mark                          int               `json:"mark,omitempty"`
	Call                          bool              `json:"call,omitempty"`
	Expired                       bool              `json:"expired,omitempty"`
	FromAPI                       bool              `json:"fromApi,omitempty"`
	MarkDatetime                  string            `json:"markDatetime,omitempty"`
	CustomerComment               string            `json:"customerComment,omitempty"`
	ManagerComment                string            `json:"managerComment,omitempty"`
	Status                        string            `json:"status,omitempty"`
	StatusComment                 string            `json:"statusComment,omitempty"`
	FullPaidAt                    string            `json:"fullPaidAt,omitempty"`
	Site                          string            `json:"site,omitempty"`
	OrderType                     string            `json:"orderType,omitempty"`
	OrderMethod                   string            `json:"orderMethod,omitempty"`
	CountryIso                    string            `json:"countryIso,omitempty"`
	Summ                          float32           `json:"summ,omitempty"`
	TotalSumm                     float32           `json:"totalSumm,omitempty"`
	PrepaySum                     float32           `json:"prepaySum,omitempty"`
	PurchaseSumm                  float32           `json:"purchaseSumm,omitempty"`
	DiscountManualAmount          float32           `json:"discountManualAmount,omitempty"`
	DiscountManualPercent         float32           `json:"discountManualPercent,omitempty"`
	Weight                        float32           `json:"weight,omitempty"`
	Length                        int               `json:"length,omitempty"`
	Width                         int               `json:"width,omitempty"`
	Height                        int               `json:"height,omitempty"`
	ShipmentStore                 string            `json:"shipmentStore,omitempty"`
	ShipmentDate                  string            `json:"shipmentDate,omitempty"`
	ClientID                      string            `json:"clientId,omitempty"`
	Shipped                       bool              `json:"shipped,omitempty"`
	UploadedToExternalStoreSystem bool              `json:"uploadedToExternalStoreSystem,omitempty"`
	Source                        *Source           `json:"source,omitempty"`
	Contragent                    *Contragent       `json:"contragent,omitempty"`
	Customer                      *Customer         `json:"customer,omitempty"`
	Delivery                      *OrderDelivery    `json:"delivery,omitempty"`
	Marketplace                   *OrderMarketplace `json:"marketplace,omitempty"`
	Items                         []OrderItem       `json:"items,omitempty"`
	CustomFields                  StringMap         `json:"customFields,omitempty"`
	Payments                      OrderPayments     `json:"payments,omitempty"`
}

// OrdersStatus type.
type OrdersStatus struct {
	ID         int    `json:"id"`
	ExternalID string `json:"externalId,omitempty"`
	Status     string `json:"status"`
	Group      string `json:"group"`
}

// OrderDelivery type.
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

// OrderDeliveryTime type.
type OrderDeliveryTime struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Custom string `json:"custom,omitempty"`
}

// OrderDeliveryService type.
type OrderDeliveryService struct {
	Name   string `json:"name,omitempty"`
	Code   string `json:"code,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// OrderDeliveryDataBasic type.
type OrderDeliveryDataBasic struct {
	TrackNumber        string `json:"trackNumber,omitempty"`
	Status             string `json:"status,omitempty"`
	PickuppointAddress string `json:"pickuppointAddress,omitempty"`
	PayerType          string `json:"payerType,omitempty"`
}

// OrderDeliveryData type.
type OrderDeliveryData struct {
	OrderDeliveryDataBasic
	AdditionalFields map[string]interface{}
}

// UnmarshalJSON method.
func (v *OrderDeliveryData) UnmarshalJSON(b []byte) error {
	var additionalData map[string]interface{}
	err := json.Unmarshal(b, &additionalData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &v.OrderDeliveryDataBasic)
	if err != nil {
		return err
	}
	object := reflect.TypeOf(v.OrderDeliveryDataBasic)

	for i := 0; i < object.NumField(); i++ {
		field := object.Field(i)

		if i, ok := field.Tag.Lookup("json"); ok {
			name := strings.Split(i, ",")[0]
			delete(additionalData, strings.TrimSpace(name))
		} else {
			delete(additionalData, field.Name)
		}
	}

	v.AdditionalFields = additionalData
	return nil
}

// MarshalJSON method.
func (v OrderDeliveryData) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	data, _ := json.Marshal(v.OrderDeliveryDataBasic)
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	for key, value := range v.AdditionalFields {
		result[key] = value
	}

	return json.Marshal(result)
}

// OrderMarketplace type.
type OrderMarketplace struct {
	Code    string `json:"code,omitempty"`
	OrderID string `json:"orderId,omitempty"`
}

// OrderPayment type.
type OrderPayment struct {
	ID         int     `json:"id,omitempty"`
	ExternalID string  `json:"externalId,omitempty"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
	PaidAt     string  `json:"paidAt,omitempty"`
	Amount     float32 `json:"amount,omitempty"`
	Comment    string  `json:"comment,omitempty"`
}

// OrderItem type.
type OrderItem struct {
	ID                    int        `json:"id,omitempty"`
	InitialPrice          float32    `json:"initialPrice,omitempty"`
	PurchasePrice         float32    `json:"purchasePrice,omitempty"`
	DiscountTotal         float32    `json:"discountTotal,omitempty"`
	DiscountManualAmount  float32    `json:"discountManualAmount,omitempty"`
	DiscountManualPercent float32    `json:"discountManualPercent,omitempty"`
	ProductName           string     `json:"productName,omitempty"`
	VatRate               string     `json:"vatRate,omitempty"`
	CreatedAt             string     `json:"createdAt,omitempty"`
	Quantity              float32    `json:"quantity,omitempty"`
	Status                string     `json:"status,omitempty"`
	Comment               string     `json:"comment,omitempty"`
	IsCanceled            bool       `json:"isCanceled,omitempty"`
	Offer                 Offer      `json:"offer,omitempty"`
	Properties            Properties `json:"properties,omitempty"`
	PriceType             *PriceType `json:"priceType,omitempty"`
}

// OrdersHistoryRecord type.
type OrdersHistoryRecord struct {
	ID        int     `json:"id,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	Created   bool    `json:"created,omitempty"`
	Deleted   bool    `json:"deleted,omitempty"`
	Source    string  `json:"source,omitempty"`
	Field     string  `json:"field,omitempty"`
	User      *User   `json:"user,omitempty"`
	APIKey    *APIKey `json:"apiKey,omitempty"`
	Order     *Order  `json:"order,omitempty"`
}

// Pack type.
type Pack struct {
	ID                 int       `json:"id,omitempty"`
	PurchasePrice      float32   `json:"purchasePrice,omitempty"`
	Quantity           float32   `json:"quantity,omitempty"`
	Store              string    `json:"store,omitempty"`
	ShipmentDate       string    `json:"shipmentDate,omitempty"`
	InvoiceNumber      string    `json:"invoiceNumber,omitempty"`
	DeliveryNoteNumber string    `json:"deliveryNoteNumber,omitempty"`
	Item               *PackItem `json:"item,omitempty"`
	ItemID             int       `json:"itemId,omitempty"`
	Unit               *Unit     `json:"unit,omitempty"`
}

// PackItem type.
type PackItem struct {
	ID    int    `json:"id,omitempty"`
	Order *Order `json:"order,omitempty"`
	Offer *Offer `json:"offer,omitempty"`
}

// PacksHistoryRecord type.
type PacksHistoryRecord struct {
	ID        int    `json:"id,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	Created   bool   `json:"created,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
	Source    string `json:"source,omitempty"`
	Field     string `json:"field,omitempty"`
	User      *User  `json:"user,omitempty"`
	Pack      *Pack  `json:"pack,omitempty"`
}

// Offer type.
type Offer struct {
	ID            int          `json:"id,omitempty"`
	ExternalID    string       `json:"externalId,omitempty"`
	Name          string       `json:"name,omitempty"`
	XMLID         string       `json:"xmlId,omitempty"`
	Article       string       `json:"article,omitempty"`
	VatRate       string       `json:"vatRate,omitempty"`
	Price         float32      `json:"price,omitempty"`
	PurchasePrice float32      `json:"purchasePrice,omitempty"`
	Quantity      float32      `json:"quantity,omitempty"`
	Height        float32      `json:"height,omitempty"`
	Width         float32      `json:"width,omitempty"`
	Length        float32      `json:"length,omitempty"`
	Weight        float32      `json:"weight,omitempty"`
	Stores        []Inventory  `json:"stores,omitempty"`
	Properties    StringMap    `json:"properties,omitempty"`
	Prices        []OfferPrice `json:"prices,omitempty"`
	Images        []string     `json:"images,omitempty"`
	Unit          *Unit        `json:"unit,omitempty"`
}

// Inventory type.
type Inventory struct {
	PurchasePrice float32 `json:"purchasePrice,omitempty"`
	Quantity      float32 `json:"quantity,omitempty"`
	Store         string  `json:"store,omitempty"`
}

// InventoryUpload type.
type InventoryUpload struct {
	ID         int                    `json:"id,omitempty"`
	ExternalID string                 `json:"externalId,omitempty"`
	XMLID      string                 `json:"xmlId,omitempty"`
	Stores     []InventoryUploadStore `json:"stores,omitempty"`
}

// InventoryUploadStore type.
type InventoryUploadStore struct {
	PurchasePrice float32 `json:"purchasePrice,omitempty"`
	Available     float32 `json:"available,omitempty"`
	Code          string  `json:"code,omitempty"`
}

// OfferPrice type.
type OfferPrice struct {
	Price     float32 `json:"price,omitempty"`
	Ordering  int     `json:"ordering,omitempty"`
	PriceType string  `json:"priceType,omitempty"`
}

// OfferPriceUpload type.
type OfferPriceUpload struct {
	ID         int           `json:"id,omitempty"`
	ExternalID string        `json:"externalId,omitempty"`
	XMLID      string        `json:"xmlId,omitempty"`
	Site       string        `json:"site,omitempty"`
	Prices     []PriceUpload `json:"prices,omitempty"`
}

// PriceUpload type.
type PriceUpload struct {
	Code  string  `json:"code,omitempty"`
	Price float32 `json:"price,omitempty"`
}

// Unit type.
type Unit struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Sym     string `json:"sym"`
	Default bool   `json:"default,omitempty"`
	Active  bool   `json:"active,omitempty"`
}

/**
User related types
*/

// User type.
type User struct {
	ID         int         `json:"id,omitempty"`
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
	Groups     []UserGroup `json:"groups,omitempty"`
	MGUserID   uint64      `json:"mgUserId,omitempty"`
}

// UserGroup type.
type UserGroup struct {
	Name                  string   `json:"name,omitempty"`
	Code                  string   `json:"code,omitempty"`
	SignatureTemplate     string   `json:"signatureTemplate,omitempty"`
	IsManager             bool     `json:"isManager,omitempty"`
	IsDeliveryMen         bool     `json:"isDeliveryMen,omitempty"`
	DeliveryTypes         []string `json:"deliveryTypes,omitempty"`
	BreakdownOrderTypes   []string `json:"breakdownOrderTypes,omitempty"`
	BreakdownSites        []string `json:"breakdownSites,omitempty"`
	BreakdownOrderMethods []string `json:"breakdownOrderMethods,omitempty"`
	GrantedOrderTypes     []string `json:"grantedOrderTypes,omitempty"`
	GrantedSites          []string `json:"grantedSites,omitempty"`
}

/**
Task related types
*/

// Task type.
type Task struct {
	ID          int       `json:"id,omitempty"`
	PerformerID int       `json:"performerId,omitempty"`
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

// Note type.
type Note struct {
	ID        int       `json:"id,omitempty"`
	ManagerID int       `json:"managerId,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt string    `json:"createdAt,omitempty"`
	Customer  *Customer `json:"customer,omitempty"`
}

/*
	Payments related types
*/

// Payment type.
type Payment struct {
	ID         int     `json:"id,omitempty"`
	ExternalID string  `json:"externalId,omitempty"`
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

// Segment type.
type Segment struct {
	ID             int    `json:"id,omitempty"`
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	CustomersCount int    `json:"customersCount,omitempty"`
	IsDynamic      bool   `json:"isDynamic,omitempty"`
	Active         bool   `json:"active,omitempty"`
}

/*
 * Settings related types
 */

// SettingsNode represents an item in settings. All settings nodes contains only string value and update time for now.
type SettingsNode struct {
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at"`
}

// Settings type. Contains retailCRM configuration.
type Settings struct {
	DefaultCurrency SettingsNode `json:"default_currency"`
	SystemLanguage  SettingsNode `json:"system_language"`
	Timezone        SettingsNode `json:"timezone"`
}

/**
Reference related types
*/

// CostGroup type.
type CostGroup struct {
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Color    string `json:"color,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Ordering int    `json:"ordering,omitempty"`
}

// CostItem type.
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

// Courier type.
type Courier struct {
	ID          int    `json:"id,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Email       string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active,omitempty"`
	Phone       *Phone `json:"phone,omitempty"`
}

// DeliveryService type.
type DeliveryService struct {
	Name   string `json:"name,omitempty"`
	Code   string `json:"code,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// DeliveryType type.
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

// LegalEntity type.
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

// OrderMethod type.
type OrderMethod struct {
	Name          string `json:"name,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        bool   `json:"active,omitempty"`
	DefaultForCRM bool   `json:"defaultForCrm,omitempty"`
	DefaultForAPI bool   `json:"defaultForApi,omitempty"`
}

// OrderType type.
type OrderType struct {
	Name          string `json:"name,omitempty"`
	Code          string `json:"code,omitempty"`
	Active        bool   `json:"active,omitempty"`
	DefaultForCRM bool   `json:"defaultForCrm,omitempty"`
	DefaultForAPI bool   `json:"defaultForApi,omitempty"`
}

// PaymentStatus type.
type PaymentStatus struct {
	Name            string   `json:"name,omitempty"`
	Code            string   `json:"code,omitempty"`
	Active          bool     `json:"active,omitempty"`
	DefaultForCRM   bool     `json:"defaultForCrm,omitempty"`
	DefaultForAPI   bool     `json:"defaultForApi,omitempty"`
	PaymentComplete bool     `json:"paymentComplete,omitempty"`
	Description     string   `json:"description,omitempty"`
	Ordering        int      `json:"ordering,omitempty"`
	PaymentTypes    []string `json:"paymentTypes,omitempty"`
}

// PaymentType type.
type PaymentType struct {
	Name            string   `json:"name,omitempty"`
	Code            string   `json:"code,omitempty"`
	Active          bool     `json:"active,omitempty"`
	DefaultForCRM   bool     `json:"defaultForCrm,omitempty"`
	DefaultForAPI   bool     `json:"defaultForApi,omitempty"`
	Description     string   `json:"description,omitempty"`
	DeliveryTypes   []string `json:"deliveryTypes,omitempty"`
	PaymentStatuses []string `json:"PaymentStatuses,omitempty"`
}

// PriceType type.
type PriceType struct {
	ID               int               `json:"id,omitempty"`
	Code             string            `json:"code,omitempty"`
	Name             string            `json:"name,omitempty"`
	Active           bool              `json:"active,omitempty"`
	Default          bool              `json:"default,omitempty"`
	Description      string            `json:"description,omitempty"`
	FilterExpression string            `json:"filterExpression,omitempty"`
	Ordering         int               `json:"ordering,omitempty"`
	Groups           []string          `json:"groups,omitempty"`
	Geo              []GeoHierarchyRow `json:"geo,omitempty"`
}

// ProductStatus type.
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

// Status type.
type Status struct {
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Ordering int    `json:"ordering,omitempty"`
	Group    string `json:"group,omitempty"`
}

// StatusGroup type.
type StatusGroup struct {
	Name     string   `json:"name,omitempty"`
	Code     string   `json:"code,omitempty"`
	Active   bool     `json:"active,omitempty"`
	Ordering int      `json:"ordering,omitempty"`
	Process  bool     `json:"process,omitempty"`
	Statuses []string `json:"statuses,omitempty"`
}

// Site type.
type Site struct {
	Name             string       `json:"name,omitempty"`
	Code             string       `json:"code,omitempty"`
	URL              string       `json:"url,omitempty"`
	Description      string       `json:"description,omitempty"`
	Phones           string       `json:"phones,omitempty"`
	Zip              string       `json:"zip,omitempty"`
	Address          string       `json:"address,omitempty"`
	CountryIso       string       `json:"countryIso,omitempty"`
	YmlURL           string       `json:"ymlUrl,omitempty"`
	LoadFromYml      bool         `json:"loadFromYml,omitempty"`
	CatalogUpdatedAt string       `json:"catalogUpdatedAt,omitempty"`
	CatalogLoadingAt string       `json:"catalogLoadingAt,omitempty"`
	Contragent       *LegalEntity `json:"contragent,omitempty"`
}

// Store type.
type Store struct {
	Name          string   `json:"name,omitempty"`
	Code          string   `json:"code,omitempty"`
	ExternalID    string   `json:"externalId,omitempty"`
	Description   string   `json:"description,omitempty"`
	XMLID         string   `json:"xmlId,omitempty"`
	Email         string   `json:"email,omitempty"`
	Type          string   `json:"type,omitempty"`
	InventoryType string   `json:"inventoryType,omitempty"`
	Active        bool     `json:"active,omitempty"`
	Phone         *Phone   `json:"phone,omitempty"`
	Address       *Address `json:"address,omitempty"`
}

// ProductGroup type.
type ProductGroup struct {
	ID       int    `json:"id,omitempty"`
	ParentID int    `json:"parentId,omitempty"`
	Name     string `json:"name,omitempty"`
	Site     string `json:"site,omitempty"`
	Active   bool   `json:"active,omitempty"`
}

// Product type.
type Product struct {
	ID           int            `json:"id,omitempty"`
	MaxPrice     float32        `json:"maxPrice,omitempty"`
	MinPrice     float32        `json:"minPrice,omitempty"`
	Name         string         `json:"name,omitempty"`
	URL          string         `json:"url,omitempty"`
	Article      string         `json:"article,omitempty"`
	ExternalID   string         `json:"externalId,omitempty"`
	Manufacturer string         `json:"manufacturer,omitempty"`
	ImageURL     string         `json:"imageUrl,omitempty"`
	Description  string         `json:"description,omitempty"`
	Popular      bool           `json:"popular,omitempty"`
	Stock        bool           `json:"stock,omitempty"`
	Novelty      bool           `json:"novelty,omitempty"`
	Recommended  bool           `json:"recommended,omitempty"`
	Active       bool           `json:"active,omitempty"`
	Quantity     float32        `json:"quantity,omitempty"`
	Offers       []Offer        `json:"offers,omitempty"`
	Groups       []ProductGroup `json:"groups,omitempty"`
	Properties   StringMap      `json:"properties,omitempty"`
}

// DeliveryHistoryRecord type.
type DeliveryHistoryRecord struct {
	Code      string `json:"code,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// DeliveryShipment type.
type DeliveryShipment struct {
	IntegrationCode string        `json:"integrationCode,omitempty"`
	ID              int           `json:"id,omitempty"`
	ExternalID      string        `json:"externalId,omitempty"`
	DeliveryType    string        `json:"deliveryType,omitempty"`
	Store           string        `json:"store,omitempty"`
	ManagerID       int           `json:"managerId,omitempty"`
	Status          string        `json:"status,omitempty"`
	Date            string        `json:"date,omitempty"`
	Time            *DeliveryTime `json:"time,omitempty"`
	LunchTime       string        `json:"lunchTime,omitempty"`
	Comment         string        `json:"comment,omitempty"`
	Orders          []Order       `json:"orders,omitempty"`
	ExtraData       StringMap     `json:"extraData,omitempty"`
}

// IntegrationModule type.
type IntegrationModule struct {
	Code               string        `json:"code,omitempty"`
	IntegrationCode    string        `json:"integrationCode,omitempty"`
	Active             bool          `json:"active,omitempty"`
	Freeze             bool          `json:"freeze,omitempty"`
	Native             bool          `json:"native,omitempty"`
	Name               string        `json:"name,omitempty"`
	Logo               string        `json:"logo,omitempty"`
	ClientID           string        `json:"clientId,omitempty"`
	BaseURL            string        `json:"baseUrl,omitempty"`
	AccountURL         string        `json:"accountUrl,omitempty"`
	AvailableCountries []string      `json:"availableCountries,omitempty"`
	Actions            StringMap     `json:"actions,omitempty"`
	Integrations       *Integrations `json:"integrations,omitempty"`
}

// UpdateScopesRequest type.
type UpdateScopesRequest struct {
	Requires ScopesRequired `json:"requires"`
}

type ScopesRequired struct {
	Scopes []string
}

// Integrations type.
type Integrations struct {
	Telephony   *Telephony   `json:"telephony,omitempty"`
	Delivery    *Delivery    `json:"delivery,omitempty"`
	Store       *Warehouse   `json:"store,omitempty"`
	MgTransport *MgTransport `json:"mgTransport,omitempty"`
	MgBot       *MgBot       `json:"mgBot,omitempty"`
}

// Delivery type.
type Delivery struct {
	Description           string              `json:"description,omitempty"`
	Actions               StringMap           `json:"actions,omitempty"`
	PayerType             []string            `json:"payerType,omitempty"`
	PlatePrintLimit       int                 `json:"platePrintLimit,omitempty"`
	RateDeliveryCost      bool                `json:"rateDeliveryCost,omitempty"`
	AllowPackages         bool                `json:"allowPackages,omitempty"`
	CodAvailable          bool                `json:"codAvailable,omitempty"`
	SelfShipmentAvailable bool                `json:"selfShipmentAvailable,omitempty"`
	AllowTrackNumber      bool                `json:"allowTrackNumber,omitempty"`
	AvailableCountries    []string            `json:"availableCountries,omitempty"`
	RequiredFields        []string            `json:"requiredFields,omitempty"`
	StatusList            []DeliveryStatus    `json:"statusList,omitempty"`
	PlateList             []Plate             `json:"plateList,omitempty"`
	DeliveryDataFieldList []DeliveryDataField `json:"deliveryDataFieldList,omitempty"`
	ShipmentDataFieldList []DeliveryDataField `json:"shipmentDataFieldList,omitempty"`
}

// DeliveryStatus type.
type DeliveryStatus struct {
	Code       string `json:"code,omitempty"`
	Name       string `json:"name,omitempty"`
	IsEditable bool   `json:"isEditable,omitempty"`
}

// Plate type.
type Plate struct {
	Code  string `json:"code,omitempty"`
	Label string `json:"label,omitempty"`
}

// DeliveryDataField type.
type DeliveryDataField struct {
	Code            string `json:"code,omitempty"`
	Label           string `json:"label,omitempty"`
	Hint            string `json:"hint,omitempty"`
	Type            string `json:"type,omitempty"`
	AutocompleteURL string `json:"autocompleteUrl,omitempty"`
	Multiple        bool   `json:"multiple,omitempty"`
	Required        bool   `json:"required,omitempty"`
	AffectsCost     bool   `json:"affectsCost,omitempty"`
	Editable        bool   `json:"editable,omitempty"`
}

// Telephony type.
type Telephony struct {
	MakeCallURL          string           `json:"makeCallUrl,omitempty"`
	AllowEdit            bool             `json:"allowEdit,omitempty"`
	InputEventSupported  bool             `json:"inputEventSupported,omitempty"`
	OutputEventSupported bool             `json:"outputEventSupported,omitempty"`
	HangupEventSupported bool             `json:"hangupEventSupported,omitempty"`
	ChangeUserStatusURL  string           `json:"changeUserStatusUrl,omitempty"`
	AdditionalCodes      []AdditionalCode `json:"additionalCodes,omitempty"`
	ExternalPhones       []ExternalPhone  `json:"externalPhones,omitempty"`
}

// AdditionalCode type.
type AdditionalCode struct {
	Code   string `json:"code,omitempty"`
	UserID string `json:"userId,omitempty"`
}

// ExternalPhone type.
type ExternalPhone struct {
	SiteCode      string `json:"siteCode,omitempty"`
	ExternalPhone string `json:"externalPhone,omitempty"`
}

// Warehouse type.
type Warehouse struct {
	Actions []Action `json:"actions,omitempty"`
}

// Action type.
type Action struct {
	Code       string   `json:"code,omitempty"`
	URL        string   `json:"url,omitempty"`
	CallPoints []string `json:"callPoints,omitempty"`
}

// MgTransport type.
type MgTransport struct {
	WebhookURL string `json:"webhookUrl,omitempty"`
}

// MgBot type.
type MgBot struct{}

/**
Cost related types
*/

// CostRecord type.
type CostRecord struct {
	Source   *Source  `json:"source,omitempty"`
	Comment  string   `json:"comment,omitempty"`
	DateFrom string   `json:"dateFrom,omitempty"`
	DateTo   string   `json:"dateTo,omitempty"`
	Summ     float32  `json:"summ,omitempty"`
	CostItem string   `json:"costItem,omitempty"`
	UserID   int      `json:"userId,omitempty"`
	Order    *Order   `json:"order,omitempty"`
	Sites    []string `json:"sites,omitempty"`
}

// Cost type.
type Cost struct {
	Source    *Source  `json:"source,omitempty"`
	ID        int      `json:"id,omitempty"`
	DateFrom  string   `json:"dateFrom,omitempty"`
	DateTo    string   `json:"dateTo,omitempty"`
	Summ      float32  `json:"summ,omitempty"`
	CostItem  string   `json:"costItem,omitempty"`
	Comment   string   `json:"comment,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	CreatedBy string   `json:"createdBy,omitempty"`
	Order     *Order   `json:"order,omitempty"`
	UserID    int      `json:"userId,omitempty"`
	Sites     []string `json:"sites,omitempty"`
}

// File type.
type File struct {
	ID         int          `json:"id,omitempty"`
	Filename   string       `json:"filename,omitempty"`
	Type       string       `json:"type,omitempty"`
	CreatedAt  string       `json:"createdAt,omitempty"`
	Size       int          `json:"size,omitempty"`
	Attachment []Attachment `json:"attachment,omitempty"`
}

// Attachment type.
type Attachment struct {
	Customer *Customer `json:"customer,omitempty"`
	Order    *Order    `json:"order,omitempty"`
}

// CustomFields type.
type CustomFields struct {
	Name           string `json:"name,omitempty"`
	Code           string `json:"code,omitempty"`
	Required       bool   `json:"required,omitempty"`
	InFilter       bool   `json:"inFilter,omitempty"`
	InList         bool   `json:"inList,omitempty"`
	InGroupActions bool   `json:"inGroupActions,omitempty"`
	Type           string `json:"type,omitempty"`
	Entity         string `json:"entity,omitempty"`
	Default        string `json:"default,omitempty"`
	Ordering       int    `json:"ordering,omitempty"`
	DisplayArea    string `json:"displayArea,omitempty"`
	ViewMode       string `json:"viewMode,omitempty"`
	Dictionary     string `json:"dictionary,omitempty"`
}

/**
CustomDictionaries related types
*/

// CustomDictionary type.
type CustomDictionary struct {
	Name     string    `json:"name,omitempty"`
	Code     string    `json:"code,omitempty"`
	Elements []Element `json:"elements,omitempty"`
}

// Element type.
type Element struct {
	Name     string `json:"name,omitempty"`
	Code     string `json:"code,omitempty"`
	Ordering int    `json:"ordering,omitempty"`
}

// Activity struct.
type Activity struct {
	Active bool `json:"active"`
	Freeze bool `json:"freeze"`
}

// Tag struct.
type Tag struct {
	Name     string `json:"name,omitempty"`
	Color    string `json:"color,omitempty"`
	Attached bool   `json:"attached,omitempty"`
}
