package retailcrm

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
