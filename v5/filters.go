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
	MinOrdersCount             string            `url:"minOrdersCount,omitempty"`
	MaxOrdersCount             string            `url:"maxOrdersCount,omitempty"`
	MinAverageSumm             string            `url:"minAverageSumm,omitempty"`
	MaxAverageSumm             string            `url:"maxAverageSumm,omitempty"`
	MinTotalSumm               string            `url:"minTotalSumm,omitempty"`
	MaxTotalSumm               string            `url:"maxTotalSumm,omitempty"`
	MinCostSumm                string            `url:"minCostSumm,omitempty"`
	MaxCostSumm                string            `url:"maxCostSumm,omitempty"`
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
	BrowserId                  string            `url:"browserId,omitempty"`
	Commentary                 string            `url:"commentary,omitempty"`
	SourceName                 string            `url:"sourceName,omitempty"`
	MediumName                 string            `url:"mediumName,omitempty"`
	CampaignName               string            `url:"campaignName,omitempty"`
	KeywordName                string            `url:"keywordName,omitempty"`
	AdContentName              string            `url:"adContentName,omitempty"`
	CustomFields               map[string]string `url:"customFields,omitempty,brackets"`
}

// CustomersHistoryFilter type
type CustomersHistoryFilter struct {
	CustomerId         int    `url:"customerId,omitempty"`
	SinceId            int    `url:"sinceId,omitempty"`
	CustomerExternalId string `url:"customerExternalId,omitempty"`
	StartDate          string `url:"startDate,omitempty"`
	EndDate            string `url:"endDate,omitempty"`
}

// OrdersFilter type
type OrdersFilter struct {
	Ids                            []string          `url:"ids,omitempty,brackets"`
	ExternalIds                    []string          `url:"externalIds,omitempty,brackets"`
	Numbers                        []string          `url:"numbers,omitempty,brackets"`
	Customer                       string            `url:"customer,omitempty"`
	CustomerId                     string            `url:"customerId,omitempty"`
	CustomerExternalId             string            `url:"customerExternalId,omitempty"`
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
	OrderId         int    `url:"orderId,omitempty"`
	SinceId         int    `url:"sinceId,omitempty"`
	OrderExternalId string `url:"orderExternalId,omitempty"`
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
	Groups        []string `url:"groups,omitempty"`
}

// TasksFilter type
type TasksFilter struct {
	OrderNumber string `url:"orderNumber,omitempty"`
	Status      string `url:"status,omitempty"`
	Customer    string `url:"customer,omitempty"`
	Text        string `url:"text,omitempty"`
	DateFrom    string `url:"dateFrom,omitempty"`
	DateTo      string `url:"dateTo,omitempty"`
	Creators    []int  `url:"creators,omitempty"`
	Performers  []int  `url:"performers,omitempty"`
}
