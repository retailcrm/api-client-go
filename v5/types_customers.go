package v5

// Customer type
type Customer struct {
	Id                           int             `json:"id,omitempty"`
	ExternalId                   string          `json:"externalId,omitempty"`
	FirstName                    string          `json:"firstName,omitempty"`
	LastName                     string          `json:"lastName,omitempty"`
	Patronymic                   string          `json:"patronymic,omitempty"`
	Sex                          string          `json:"sex,omitempty"`
	Email                        string          `json:"email,omitempty"`
	Phones                       []CustomerPhone `json:"phones,brackets,omitempty"`
	Address                      *Address        `json:"address,omitempty"`
	CreatedAt                    string          `json:"createdAt,omitempty"`
	Birthday                     string          `json:"birthday,omitempty"`
	ManagerId                    int             `json:"managerId,omitempty"`
	Vip                          bool            `json:"vip,omitempty"`
	Bad                          bool            `json:"bad,omitempty"`
	Site                         string          `json:"site,omitempty"`
	Source                       *Source         `json:"source,omitempty"`
	Contragent                   *Contragent     `json:"contragent,omitempty"`
	PersonalDiscount             float32         `json:"personalDiscount,omitempty"`
	CumulativeDiscount           float32         `json:"cumulativeDiscount,omitempty"`
	DiscountCardNumber           string          `json:"discountCardNumber,omitempty"`
	EmailMarketingUnsubscribedAt string          `json:"emailMarketingUnsubscribedAt,omitempty"`
	AvgMarginSumm                float32         `json:"avgMarginSumm,omitempty"`
	MarginSumm                   float32         `json:"marginSumm,omitempty"`
	TotalSumm                    float32         `json:"totalSumm,omitempty"`
	AverageSumm                  float32         `json:"averageSumm,omitempty"`
	OrdersCount                  int             `json:"ordersCount,omitempty"`
	CostSumm                     float32         `json:"costSumm,omitempty"`
	MaturationTime               int             `json:"maturationTime,omitempty"`
	FirstClientId                string          `json:"firstClientId,omitempty"`
	LastClientId                 string          `json:"lastClientId,omitempty"`
	//CustomFields                 map[string]string `json:"customFields,omitempty"`
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

// CustomerFilter for get customer request
type CustomerFilter struct {
	By   string `url:"by,omitempty"`
	Site string `url:"site,omitempty"`
}

// CustomerResponse type
type CustomerResponse struct {
	Success  bool      `json:"success"`
	Customer *Customer `json:"customer,omitempty,brackets"`
}

type CustomersFilter struct {
	ExternalIds []string `url:"externalIds,omitempty,brackets"`
	City        string   `url:"city,omitempty"`
}

type CustomersParameters struct {
	Filter CustomersFilter `url:"filter,omitempty"`
	Limit  int             `url:"limit,omitempty"`
	Page   int             `url:"page,omitempty"`
}

type CustomersResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Customers  []Customer  `json:"customers,omitempty,brackets"`
}

type CustomerChangeResponse struct {
	Success bool   `json:"success"`
	Id      int    `json:"id,omitempty"`
	State   string `json:"state,omitempty"`
}

type CustomersUploadParameters struct {
	Customers []Customer `url:"customers,omitempty,brackets"`
	Site      string     `url:"site,omitempty"`
}

type CustomersUploadResponse struct {
	Success           bool                  `json:"success"`
	UploadedCustomers []CustomerIdentifiers `json:"uploadedCustomers,omitempty,brackets"`
}

type CustomersHistoryFilter struct {
	CustomerId         int    `url:"customerId,omitempty"`
	SinceId            int    `url:"sinceId,omitempty"`
	CustomerExternalId string `url:"customerExternalId,omitempty"`
	StartDate          string `url:"startDate,omitempty"`
	EndDate            string `url:"endDate,omitempty"`
}

type CustomersHistoryParameters struct {
	Filter CustomersHistoryFilter `url:"filter,omitempty"`
	Limit  int                    `url:"limit,omitempty"`
	Page   int                    `url:"page,omitempty"`
}

type CustomersHistoryResponse struct {
	Success     bool                    `json:"success,omitempty"`
	GeneratedAt string                  `json:"generatedAt,omitempty"`
	History     []CustomerHistoryRecord `json:"history,omitempty,brackets"`
	Pagination  *Pagination             `json:"pagination,omitempty"`
}
