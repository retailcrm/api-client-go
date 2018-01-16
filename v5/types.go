package v5

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

// Customer type
type Customer struct {
	Id                           int               `json:"id,omitempty"`
	ExternalId                   string            `json:"externalId,omitempty"`
	FirstName                    string            `json:"firstName,omitempty"`
	LastName                     string            `json:"lastName,omitempty"`
	Patronymic                   string            `json:"patronymic,omitempty"`
	Sex                          string            `json:"sex,omitempty"`
	Email                        string            `json:"email,omitempty"`
	Phones                       []CustomerPhone   `json:"phones,brackets,omitempty"`
	Address                      *Address          `json:"address,omitempty"`
	CreatedAt                    string            `json:"createdAt,omitempty"`
	Birthday                     string            `json:"birthday,omitempty"`
	ManagerId                    int               `json:"managerId,omitempty"`
	Vip                          bool              `json:"vip,omitempty"`
	Bad                          bool              `json:"bad,omitempty"`
	Site                         string            `json:"site,omitempty"`
	Source                       *Source           `json:"source,omitempty"`
	Contragent                   *Contragent       `json:"contragent,omitempty"`
	PersonalDiscount             float32           `json:"personalDiscount,omitempty"`
	CumulativeDiscount           float32           `json:"cumulativeDiscount,omitempty"`
	DiscountCardNumber           string            `json:"discountCardNumber,omitempty"`
	EmailMarketingUnsubscribedAt string            `json:"emailMarketingUnsubscribedAt,omitempty"`
	AvgMarginSumm                float32           `json:"avgMarginSumm,omitempty"`
	MarginSumm                   float32           `json:"marginSumm,omitempty"`
	TotalSumm                    float32           `json:"totalSumm,omitempty"`
	AverageSumm                  float32           `json:"averageSumm,omitempty"`
	OrdersCount                  int               `json:"ordersCount,omitempty"`
	CostSumm                     float32           `json:"costSumm,omitempty"`
	MaturationTime               int               `json:"maturationTime,omitempty"`
	FirstClientId                string            `json:"firstClientId,omitempty"`
	LastClientId                 string            `json:"lastClientId,omitempty"`
	//CustomFields                 map[string]string `json:"customFields,omitempty"`
}

// CustomerPhone type
type CustomerPhone struct {
	Number string `json:"number,omitempty"`
}

// CustomerGetFilter for get customer request
type CustomerGetFilter struct {
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
	Filter CustomersFilter  `url:"filter,omitempty"`
	Limit  int              `url:"limit,omitempty"`
	Page   int              `url:"page,omitempty"`
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
