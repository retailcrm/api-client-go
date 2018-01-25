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

type ApiKey struct {
	Current bool `json:"current,omitempty"`
}

type Property struct {
	Code  string `json:"code,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Customer type
type Customer struct {
	Id                           int                 `json:"id,omitempty"`
	ExternalId                   string              `json:"externalId,omitempty"`
	FirstName                    string              `json:"firstName,omitempty"`
	LastName                     string              `json:"lastName,omitempty"`
	Patronymic                   string              `json:"patronymic,omitempty"`
	Sex                          string              `json:"sex,omitempty"`
	Email                        string              `json:"email,omitempty"`
	Phones                       []CustomerPhone     `json:"phones,brackets,omitempty"`
	Address                      *Address            `json:"address,omitempty"`
	CreatedAt                    string              `json:"createdAt,omitempty"`
	Birthday                     string              `json:"birthday,omitempty"`
	ManagerId                    int                 `json:"managerId,omitempty"`
	Vip                          bool                `json:"vip,omitempty"`
	Bad                          bool                `json:"bad,omitempty"`
	Site                         string              `json:"site,omitempty"`
	Source                       *Source             `json:"source,omitempty"`
	Contragent                   *Contragent         `json:"contragent,omitempty"`
	PersonalDiscount             float32             `json:"personalDiscount,omitempty"`
	CumulativeDiscount           float32             `json:"cumulativeDiscount,omitempty"`
	DiscountCardNumber           string              `json:"discountCardNumber,omitempty"`
	EmailMarketingUnsubscribedAt string              `json:"emailMarketingUnsubscribedAt,omitempty"`
	AvgMarginSumm                float32             `json:"avgMarginSumm,omitempty"`
	MarginSumm                   float32             `json:"marginSumm,omitempty"`
	TotalSumm                    float32             `json:"totalSumm,omitempty"`
	AverageSumm                  float32             `json:"averageSumm,omitempty"`
	OrdersCount                  int                 `json:"ordersCount,omitempty"`
	CostSumm                     float32             `json:"costSumm,omitempty"`
	MaturationTime               int                 `json:"maturationTime,omitempty"`
	FirstClientId                string              `json:"firstClientId,omitempty"`
	LastClientId                 string              `json:"lastClientId,omitempty"`
	BrowserId                    string              `json:"browserId,omitempty"`
	CustomFields                 []map[string]string `json:"customFields,omitempty,brackets"`
}

// CustomerPhone type
type CustomerPhone struct {
	Number string `json:"number,omitempty"`
}

// CustomerIdentifiers type
type CustomerIdentifiers struct {
	Id         int    `json:"id,omitempty"`
	ExternalId string `json:"externalId,omitempty"`
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
	//Payments                      []OrderPayment    `json:"payments,omitempty,brackets"`
}

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

type OrderDeliveryTime struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Custom string `json:"custom,omitempty"`
}

type OrderDeliveryService struct {
	Name   string `json:"name,omitempty"`
	Code   string `json:"code,omitempty"`
	Active bool   `json:"active,omitempty"`
}

type OrderDeliveryData struct {
	TrackNumber        string `json:"trackNumber,omitempty"`
	Status             string `json:"status,omitempty"`
	PickuppointAddress string `json:"pickuppointAddress,omitempty"`
	PayerType          string `json:"payerType,omitempty"`
}

type OrderMarketplace struct {
	Code    string `json:"code,omitempty"`
	OrderId string `json:"orderId,omitempty"`
}

type OrderPayment struct {
	Id         int     `json:"id,omitempty"`
	ExternalId string  `json:"externalId,omitempty"`
	Type       string  `json:"type,omitempty"`
	Status     string  `json:"status,omitempty"`
	PaidAt     string  `json:"paidAt,omitempty"`
	Amount     float32 `json:"amount,omitempty"`
	Comment    string  `json:"comment,omitempty"`
}

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
	Offer                 *Offer      `json:"offer,omitempty"`
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

type Offer struct {
	Id         int    `json:"id,omitempty"`
	ExternalId string `json:"externalId,omitempty"`
	XmlId      string `json:"xmlId,omitempty"`
	VatRate    string `json:"vatRate,omitempty"`
}

type User struct {
	Id         int    `json:"id,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
}

type PriceType struct {
	Name             string `json:"name,omitempty"`
	Code             string `json:"code,omitempty"`
	Description      string `json:"description,omitempty"`
	FilterExpression string `json:"filterExpression,omitempty"`
	Active           bool   `json:"active,omitempty"`
	Ordering         int    `json:"ordering,omitempty"`
}
