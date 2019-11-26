package v5

// CustomersFilter type
type CustomersFilter struct {
	Ids                        []string          `url:"ids,omitempty,brackets"`
	ExternalIds                []string          `url:"externalIds,omitempty,brackets"`
	City                       string            `url:"city,omitempty"`
	Region                     string            `url:"region,omitempty"`
	Name                       string            `url:"name,omitempty"`
	Email                      string            `url:"email,omitempty"`
	Notes                      string            `url:"notes,omitempty"`
	MinOrdersCount             int               `url:"minOrdersCount,omitempty"`
	MaxOrdersCount             int               `url:"maxOrdersCount,omitempty"`
	MinAverageSumm             float32           `url:"minAverageSumm,omitempty"`
	MaxAverageSumm             float32           `url:"maxAverageSumm,omitempty"`
	MinTotalSumm               float32           `url:"minTotalSumm,omitempty"`
	MaxTotalSumm               float32           `url:"maxTotalSumm,omitempty"`
	MinCostSumm                float32           `url:"minCostSumm,omitempty"`
	MaxCostSumm                float32           `url:"maxCostSumm,omitempty"`
	ClassSegment               string            `url:"classSegment,omitempty"`
	Vip                        int               `url:"vip,omitempty"`
	Bad                        int               `url:"bad,omitempty"`
	Attachments                int               `url:"attachments,omitempty"`
	Online                     int               `url:"online,omitempty"`
	EmailMarketingUnsubscribed int               `url:"emailMarketingUnsubscribed,omitempty"`
	Sex                        string            `url:"sex,omitempty"`
	Segment                    string            `url:"segment,omitempty"`
	DiscountCardNumber         string            `url:"discountCardNumber,omitempty"`
	ContragentName             string            `url:"contragentName,omitempty"`
	ContragentInn              string            `url:"contragentInn,omitempty"`
	ContragentKpp              string            `url:"contragentKpp,omitempty"`
	ContragentBik              string            `url:"contragentBik,omitempty"`
	ContragentCorrAccount      string            `url:"contragentCorrAccount,omitempty"`
	ContragentBankAccount      string            `url:"contragentBankAccount,omitempty"`
	ContragentTypes            []string          `url:"contragentTypes,omitempty,brackets"`
	Sites                      []string          `url:"sites,omitempty,brackets"`
	Managers                   []string          `url:"managers,omitempty,brackets"`
	ManagerGroups              []string          `url:"managerGroups,omitempty,brackets"`
	DateFrom                   string            `url:"dateFrom,omitempty"`
	DateTo                     string            `url:"dateTo,omitempty"`
	FirstWebVisitFrom          string            `url:"firstWebVisitFrom,omitempty"`
	FirstWebVisitTo            string            `url:"firstWebVisitTo,omitempty"`
	LastWebVisitFrom           string            `url:"lastWebVisitFrom,omitempty"`
	LastWebVisitTo             string            `url:"lastWebVisitTo,omitempty"`
	FirstOrderFrom             string            `url:"firstOrderFrom,omitempty"`
	FirstOrderTo               string            `url:"firstOrderTo,omitempty"`
	LastOrderFrom              string            `url:"lastOrderFrom,omitempty"`
	LastOrderTo                string            `url:"lastOrderTo,omitempty"`
	BrowserID                  string            `url:"browserId,omitempty"`
	Commentary                 string            `url:"commentary,omitempty"`
	SourceName                 string            `url:"sourceName,omitempty"`
	MediumName                 string            `url:"mediumName,omitempty"`
	CampaignName               string            `url:"campaignName,omitempty"`
	KeywordName                string            `url:"keywordName,omitempty"`
	AdContentName              string            `url:"adContentName,omitempty"`
	MgCustomerID               string            `url:"mgCustomerId,omitempty"`
	CustomFields               map[string]string `url:"customFields,omitempty,brackets"`
}

// CorporateCustomersFilter type
type CorporateCustomersFilter struct {
	ContragentName        string            `url:"contragentName,omitempty"`
	ContragentInn         string            `url:"contragentInn,omitempty"`
	ContragentKpp         string            `url:"contragentKpp,omitempty"`
	ContragentBik         string            `url:"contragentBik,omitempty"`
	ContragentCorrAccount string            `url:"contragentCorrAccount,omitempty"`
	ContragentBankAccount string            `url:"contragentBankAccount,omitempty"`
	ContragentTypes       []string          `url:"contragentTypes,omitempty,brackets"`
	ExternalIds           []string          `url:"externalIds,omitempty,brackets"`
	Name                  string            `url:"name,omitempty"`
	City                  string            `url:"city,omitempty"`
	Region                string            `url:"region,omitempty"`
	Email                 string            `url:"email,omitempty"`
	Notes                 string            `url:"notes,omitempty"`
	MinOrdersCount        int               `url:"minOrdersCount,omitempty"`
	MaxOrdersCount        int               `url:"maxOrdersCount,omitempty"`
	MinAverageSumm        float32           `url:"minAverageSumm,omitempty"`
	MaxAverageSumm        float32           `url:"maxAverageSumm,omitempty"`
	MinTotalSumm          float32           `url:"minTotalSumm,omitempty"`
	MaxTotalSumm          float32           `url:"maxTotalSumm,omitempty"`
	ClassSegment          string            `url:"classSegment,omitempty"`
	DiscountCardNumber    string            `url:"discountCardNumber,omitempty"`
	Attachments           int               `url:"attachments,omitempty"`
	MinCostSumm           float32           `url:"minCostSumm,omitempty"`
	MaxCostSumm           float32           `url:"maxCostSumm,omitempty"`
	Vip                   int               `url:"vip,omitempty"`
	Bad                   int               `url:"bad,omitempty"`
	TasksCount            int               `url:"tasksCounts,omitempty"`
	Ids                   []string          `url:"ids,omitempty,brackets"`
	Sites                 []string          `url:"sites,omitempty,brackets"`
	Managers              []string          `url:"managers,omitempty,brackets"`
	ManagerGroups         []string          `url:"managerGroups,omitempty,brackets"`
	DateFrom              string            `url:"dateFrom,omitempty"`
	DateTo                string            `url:"dateTo,omitempty"`
	FirstOrderFrom        string            `url:"firstOrderFrom,omitempty"`
	FirstOrderTo          string            `url:"firstOrderTo,omitempty"`
	LastOrderFrom         string            `url:"lastOrderFrom,omitempty"`
	LastOrderTo           string            `url:"lastOrderTo,omitempty"`
	CustomFields          map[string]string `url:"customFields,omitempty,brackets"`
}

// CorporateCustomersNotesFilter type
type CorporateCustomersNotesFilter struct {
	Ids                 []string `url:"ids,omitempty,brackets"`
	CustomerIds         []string `url:"ids,omitempty,brackets"`
	CustomerExternalIds []string `url:"customerExternalIds,omitempty,brackets"`
	ManagerIds          []string `url:"managerIds,omitempty,brackets"`
	Text                string   `url:"text,omitempty"`
	CreatedAtFrom       string   `url:"createdAtFrom,omitempty"`
	CreatedAtTo         string   `url:"createdAtTo,omitempty"`
}

// CorporateCustomerAddressesFilter type
type CorporateCustomerAddressesFilter struct {
	Ids    []string `url:"ids,omitempty,brackets"`
	Name   string   `url:"name,omitempty"`
	City   string   `url:"city,omitempty"`
	Region string   `url:"region,omitempty"`
}

// IdentifiersPairFilter type
type IdentifiersPairFilter struct {
	Ids         []string `url:"ids,omitempty,brackets"`
	ExternalIds []string `url:"externalIds,omitempty,brackets"`
}

// CustomersHistoryFilter type
type CustomersHistoryFilter struct {
	CustomerID         int    `url:"customerId,omitempty"`
	SinceID            int    `url:"sinceId,omitempty"`
	CustomerExternalID string `url:"customerExternalId,omitempty"`
	StartDate          string `url:"startDate,omitempty"`
	EndDate            string `url:"endDate,omitempty"`
}

// CorporateCustomersHistoryFilter type
type CorporateCustomersHistoryFilter struct {
	CustomerID         int      `url:"customerId,omitempty"`
	SinceID            int      `url:"sinceId,omitempty"`
	CustomerExternalID string   `url:"customerExternalId,omitempty"`
	ContactIds         []string `url:"contactIds,omitempty,brackets"`
	StartDate          string   `url:"startDate,omitempty"`
	EndDate            string   `url:"endDate,omitempty"`
}

// OrdersFilter type
type OrdersFilter struct {
	Ids                            []int             `url:"ids,omitempty,brackets"`
	ExternalIds                    []string          `url:"externalIds,omitempty,brackets"`
	Numbers                        []string          `url:"numbers,omitempty,brackets"`
	Customer                       string            `url:"customer,omitempty"`
	CustomerID                     string            `url:"customerId,omitempty"`
	CustomerExternalID             string            `url:"customerExternalId,omitempty"`
	Countries                      []string          `url:"countries,omitempty,brackets"`
	City                           string            `url:"city,omitempty"`
	Region                         string            `url:"region,omitempty"`
	Index                          string            `url:"index,omitempty"`
	Metro                          string            `url:"metro,omitempty"`
	Email                          string            `url:"email,omitempty"`
	DeliveryTimeFrom               string            `url:"deliveryTimeFrom,omitempty"`
	DeliveryTimeTo                 string            `url:"deliveryTimeTo,omitempty"`
	MinPrepaySumm                  string            `url:"minPrepaySumm,omitempty"`
	MaxPrepaySumm                  string            `url:"maxPrepaySumm,omitempty"`
	MinPrice                       string            `url:"minPrice,omitempty"`
	MaxPrice                       string            `url:"maxPrice,omitempty"`
	Product                        string            `url:"product,omitempty"`
	Vip                            int               `url:"vip,omitempty"`
	Bad                            int               `url:"bad,omitempty"`
	Attachments                    int               `url:"attachments,omitempty"`
	Expired                        int               `url:"expired,omitempty"`
	Call                           int               `url:"call,omitempty"`
	Online                         int               `url:"online,omitempty"`
	Shipped                        int               `url:"shipped,omitempty"`
	UploadedToExtStoreSys          int               `url:"uploadedToExtStoreSys,omitempty"`
	ReceiptFiscalDocumentAttribute int               `url:"receiptFiscalDocumentAttribute,omitempty"`
	ReceiptStatus                  int               `url:"receiptStatus,omitempty"`
	ReceiptOperation               int               `url:"receiptOperation,omitempty"`
	MinDeliveryCost                string            `url:"minDeliveryCost,omitempty"`
	MaxDeliveryCost                string            `url:"maxDeliveryCost,omitempty"`
	MinDeliveryNetCost             string            `url:"minDeliveryNetCost,omitempty"`
	MaxDeliveryNetCost             string            `url:"maxDeliveryNetCost,omitempty"`
	ManagerComment                 string            `url:"managerComment,omitempty"`
	CustomerComment                string            `url:"customerComment,omitempty"`
	MinMarginSumm                  string            `url:"minMarginSumm,omitempty"`
	MaxMarginSumm                  string            `url:"maxMarginSumm,omitempty"`
	MinPurchaseSumm                string            `url:"minPurchaseSumm,omitempty"`
	MaxPurchaseSumm                string            `url:"maxPurchaseSumm,omitempty"`
	MinCostSumm                    string            `url:"minCostSumm,omitempty"`
	MaxCostSumm                    string            `url:"maxCostSumm,omitempty"`
	TrackNumber                    string            `url:"trackNumber,omitempty"`
	ContragentName                 string            `url:"contragentName,omitempty"`
	ContragentInn                  string            `url:"contragentInn,omitempty"`
	ContragentKpp                  string            `url:"contragentKpp,omitempty"`
	ContragentBik                  string            `url:"contragentBik,omitempty"`
	ContragentCorrAccount          string            `url:"contragentCorrAccount,omitempty"`
	ContragentBankAccount          string            `url:"contragentBankAccount,omitempty"`
	ContragentTypes                []string          `url:"contragentTypes,omitempty,brackets"`
	OrderTypes                     []string          `url:"orderTypes,omitempty,brackets"`
	PaymentStatuses                []string          `url:"paymentStatuses,omitempty,brackets"`
	PaymentTypes                   []string          `url:"paymentTypes,omitempty,brackets"`
	DeliveryTypes                  []string          `url:"deliveryTypes,omitempty,brackets"`
	OrderMethods                   []string          `url:"orderMethods,omitempty,brackets"`
	ShipmentStores                 []string          `url:"shipmentStores,omitempty,brackets"`
	Couriers                       []string          `url:"couriers,omitempty,brackets"`
	Managers                       []string          `url:"managers,omitempty,brackets"`
	ManagerGroups                  []string          `url:"managerGroups,omitempty,brackets"`
	Sites                          []string          `url:"sites,omitempty,brackets"`
	CreatedAtFrom                  string            `url:"createdAtFrom,omitempty"`
	CreatedAtTo                    string            `url:"createdAtTo,omitempty"`
	FullPaidAtFrom                 string            `url:"fullPaidAtFrom,omitempty"`
	FullPaidAtTo                   string            `url:"fullPaidAtTo,omitempty"`
	DeliveryDateFrom               string            `url:"deliveryDateFrom,omitempty"`
	DeliveryDateTo                 string            `url:"deliveryDateTo,omitempty"`
	StatusUpdatedAtFrom            string            `url:"statusUpdatedAtFrom,omitempty"`
	StatusUpdatedAtTo              string            `url:"statusUpdatedAtTo,omitempty"`
	DpdParcelDateFrom              string            `url:"dpdParcelDateFrom,omitempty"`
	DpdParcelDateTo                string            `url:"dpdParcelDateTo,omitempty"`
	FirstWebVisitFrom              string            `url:"firstWebVisitFrom,omitempty"`
	FirstWebVisitTo                string            `url:"firstWebVisitTo,omitempty"`
	LastWebVisitFrom               string            `url:"lastWebVisitFrom,omitempty"`
	LastWebVisitTo                 string            `url:"lastWebVisitTo,omitempty"`
	FirstOrderFrom                 string            `url:"firstOrderFrom,omitempty"`
	FirstOrderTo                   string            `url:"firstOrderTo,omitempty"`
	LastOrderFrom                  string            `url:"lastOrderFrom,omitempty"`
	LastOrderTo                    string            `url:"lastOrderTo,omitempty"`
	ShipmentDateFrom               string            `url:"shipmentDateFrom,omitempty"`
	ShipmentDateTo                 string            `url:"shipmentDateTo,omitempty"`
	ExtendedStatus                 []string          `url:"extendedStatus,omitempty,brackets"`
	SourceName                     string            `url:"sourceName,omitempty"`
	MediumName                     string            `url:"mediumName,omitempty"`
	CampaignName                   string            `url:"campaignName,omitempty"`
	KeywordName                    string            `url:"keywordName,omitempty"`
	AdContentName                  string            `url:"adContentName,omitempty"`
	CustomFields                   map[string]string `url:"customFields,omitempty,brackets"`
}

// OrdersHistoryFilter type
type OrdersHistoryFilter struct {
	OrderID         int    `url:"orderId,omitempty"`
	SinceID         int    `url:"sinceId,omitempty"`
	OrderExternalID string `url:"orderExternalId,omitempty"`
	StartDate       string `url:"startDate,omitempty"`
	EndDate         string `url:"endDate,omitempty"`
}

// UsersFilter type
type UsersFilter struct {
	Email         string   `url:"email,omitempty"`
	Status        string   `url:"status,omitempty"`
	Online        int      `url:"online,omitempty"`
	Active        int      `url:"active,omitempty"`
	IsManager     int      `url:"isManager,omitempty"`
	IsAdmin       int      `url:"isAdmin,omitempty"`
	CreatedAtFrom string   `url:"createdAtFrom,omitempty"`
	CreatedAtTo   string   `url:"createdAtTo,omitempty"`
	Groups        []string `url:"groups,omitempty,brackets"`
}

// TasksFilter type
type TasksFilter struct {
	OrderNumber string `url:"orderNumber,omitempty"`
	Status      string `url:"status,omitempty"`
	Customer    string `url:"customer,omitempty"`
	Text        string `url:"text,omitempty"`
	DateFrom    string `url:"dateFrom,omitempty"`
	DateTo      string `url:"dateTo,omitempty"`
	Creators    []int  `url:"creators,omitempty,brackets"`
	Performers  []int  `url:"performers,omitempty,brackets"`
}

// NotesFilter type
type NotesFilter struct {
	Ids                 []int    `url:"ids,omitempty,brackets"`
	CustomerIds         []int    `url:"customerIds,omitempty,brackets"`
	CustomerExternalIds []string `url:"customerExternalIds,omitempty,brackets"`
	ManagerIds          []int    `url:"managerIds,omitempty,brackets"`
	Text                string   `url:"text,omitempty"`
	CreatedAtFrom       string   `url:"createdAtFrom,omitempty"`
	CreatedAtTo         string   `url:"createdAtTo,omitempty"`
}

// SegmentsFilter type
type SegmentsFilter struct {
	Ids               []int  `url:"ids,omitempty,brackets"`
	Active            int    `url:"active,omitempty"`
	Name              string `url:"name,omitempty"`
	Type              string `url:"type,omitempty"`
	MinCustomersCount int    `url:"minCustomersCount,omitempty"`
	MaxCustomersCount int    `url:"maxCustomersCount,omitempty"`
	DateFrom          string `url:"dateFrom,omitempty"`
	DateTo            string `url:"dateTo,omitempty"`
}

// PacksFilter type
type PacksFilter struct {
	Ids                []int    `url:"ids,omitempty,brackets"`
	Stores             []string `url:"stores,omitempty"`
	ItemID             int      `url:"itemId,omitempty"`
	OfferXMLID         string   `url:"offerXmlId,omitempty"`
	OfferExternalID    string   `url:"offerExternalId,omitempty"`
	OrderID            int      `url:"orderId,omitempty"`
	OrderExternalID    string   `url:"orderExternalId,omitempty"`
	ShipmentDateFrom   string   `url:"shipmentDateFrom,omitempty"`
	ShipmentDateTo     string   `url:"shipmentDateTo,omitempty"`
	InvoiceNumber      string   `url:"invoiceNumber,omitempty"`
	DeliveryNoteNumber string   `url:"deliveryNoteNumber,omitempty"`
}

// InventoriesFilter type
type InventoriesFilter struct {
	Ids               []int    `url:"ids,omitempty,brackets"`
	ProductExternalID string   `url:"productExternalId,omitempty"`
	ProductArticle    string   `url:"productArticle,omitempty"`
	OfferExternalID   string   `url:"offerExternalId,omitempty"`
	OfferXMLID        string   `url:"offerXmlId,omitempty"`
	OfferArticle      string   `url:"offerArticle,omitempty"`
	ProductActive     int      `url:"productActive,omitempty"`
	Details           int      `url:"details,omitempty"`
	Sites             []string `url:"sites,omitempty,brackets"`
}

// ProductsGroupsFilter type
type ProductsGroupsFilter struct {
	Ids           []int    `url:"ids,omitempty,brackets"`
	Sites         []string `url:"sites,omitempty,brackets"`
	Active        int      `url:"active,omitempty"`
	ParentGroupID string   `url:"parentGroupId,omitempty"`
}

// ProductsFilter type
type ProductsFilter struct {
	Ids              []int             `url:"ids,omitempty,brackets"`
	OfferIds         []int             `url:"offerIds,omitempty,brackets"`
	Active           int               `url:"active,omitempty"`
	Recommended      int               `url:"recommended,omitempty"`
	Novelty          int               `url:"novelty,omitempty"`
	Stock            int               `url:"stock,omitempty"`
	Popular          int               `url:"popular,omitempty"`
	MaxQuantity      float32           `url:"maxQuantity,omitempty"`
	MinQuantity      float32           `url:"minQuantity,omitempty"`
	MaxPurchasePrice float32           `url:"maxPurchasePrice,omitempty"`
	MinPurchasePrice float32           `url:"minPurchasePrice,omitempty"`
	MaxPrice         float32           `url:"maxPrice,omitempty"`
	MinPrice         float32           `url:"minPrice,omitempty"`
	Groups           string            `url:"groups,omitempty"`
	Name             string            `url:"name,omitempty"`
	ClassSegment     string            `url:"classSegment,omitempty"`
	XMLID            string            `url:"xmlId,omitempty"`
	ExternalID       string            `url:"externalId,omitempty"`
	Manufacturer     string            `url:"manufacturer,omitempty"`
	URL              string            `url:"url,omitempty"`
	PriceType        string            `url:"priceType,omitempty"`
	OfferExternalID  string            `url:"offerExternalId,omitempty"`
	Sites            []string          `url:"sites,omitempty,brackets"`
	Properties       map[string]string `url:"properties,omitempty,brackets"`
}

// ProductsPropertiesFilter type
type ProductsPropertiesFilter struct {
	Code  string   `url:"code,omitempty"`
	Name  string   `url:"name,omitempty"`
	Sites []string `url:"sites,omitempty,brackets"`
}

// ShipmentFilter type
type ShipmentFilter struct {
	Ids           []int    `url:"ids,omitempty,brackets"`
	ExternalID    string   `url:"externalId,omitempty"`
	OrderNumber   string   `url:"orderNumber,omitempty"`
	DateFrom      string   `url:"dateFrom,omitempty"`
	DateTo        string   `url:"dateTo,omitempty"`
	Stores        []string `url:"stores,omitempty,brackets"`
	Managers      []string `url:"managers,omitempty,brackets"`
	DeliveryTypes []string `url:"deliveryTypes,omitempty,brackets"`
	Statuses      []string `url:"statuses,omitempty,brackets"`
}

// CostsFilter type
type CostsFilter struct {
	MinSumm          string   `url:"minSumm,omitempty"`
	MaxSumm          string   `url:"maxSumm,omitempty"`
	OrderNumber      string   `url:"orderNumber,omitempty"`
	Comment          string   `url:"orderNumber,omitempty"`
	Ids              []string `url:"ids,omitempty,brackets"`
	Sites            []string `url:"sites,omitempty,brackets"`
	CreatedBy        []string `url:"createdBy,omitempty,brackets"`
	CostGroups       []string `url:"costGroups,omitempty,brackets"`
	CostItems        []string `url:"costItems,omitempty,brackets"`
	Users            []string `url:"users,omitempty,brackets"`
	DateFrom         string   `url:"dateFrom,omitempty"`
	DateTo           string   `url:"dateTo,omitempty"`
	CreatedAtFrom    string   `url:"createdAtFrom,omitempty"`
	CreatedAtTo      string   `url:"createdAtTo,omitempty"`
	OrderIds         []string `url:"orderIds,omitempty,brackets"`
	OrderExternalIds []string `url:"orderIds,omitempty,brackets"`
}

// FilesFilter type
type FilesFilter struct {
	Ids                 []int    `url:"ids,omitempty,brackets"`
	OrderIds            []int    `url:"orderIds,omitempty,brackets"`
	OrderExternalIds    []string `url:"orderExternalIds,omitempty,brackets"`
	CustomerIds         []int    `url:"customerIds,omitempty,brackets"`
	CustomerExternalIds []string `url:"customerExternalIds,omitempty,brackets"`
	CreatedAtFrom       string   `url:"createdAtFrom,omitempty"`
	CreatedAtTo         string   `url:"createdAtTo,omitempty"`
	SizeFrom            int      `url:"sizeFrom,omitempty"`
	SizeTo              int      `url:"sizeTo,omitempty"`
	Type                []string `url:"type,omitempty,brackets"`
	Filename            string   `url:"filename,omitempty"`
	IsAttached          string   `url:"isAttached,omitempty"`
	Sites               []string `url:"sites,omitempty,brackets"`
}

// CustomFieldsFilter type
type CustomFieldsFilter struct {
	Name        string `url:"name,omitempty"`
	Code        string `url:"code,omitempty"`
	Type        string `url:"type,omitempty"`
	Entity      string `url:"entity,omitempty"`
	ViewMode    string `url:"viewMode,omitempty"`
	DisplayArea string `url:"displayArea,omitempty"`
}

// CustomDictionariesFilter type
type CustomDictionariesFilter struct {
	Name string `url:"name,omitempty"`
	Code string `url:"code,omitempty"`
}
