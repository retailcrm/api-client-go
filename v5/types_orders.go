package v5

// Order type
type Order struct {
	Id                            int     `json:"id,omitempty"`
	ExternalId                    string  `json:"externalId,omitempty"`
	Number                        string  `json:"number,omitempty"`
	FirstName                     string  `json:"firstName,omitempty"`
	LastName                      string  `json:"lastName,omitempty"`
	Patronymic                    string  `json:"patronymic,omitempty"`
	Email                         string  `json:"email,omitempty"`
	Phone                         string  `json:"phone,omitempty"`
	AdditionalPhone               string  `json:"additionalPhone,omitempty"`
	CreatedAt                     string  `json:"createdAt,omitempty"`
	StatusUpdatedAt               string  `json:"statusUpdatedAt,omitempty"`
	ManagerId                     int     `json:"managerId,omitempty"`
	Mark                          int     `json:"mark,omitempty"`
	Call                          bool    `json:"call,omitempty"`
	Expired                       bool    `json:"expired,omitempty"`
	FromApi                       bool    `json:"fromApi,omitempty"`
	MarkDatetime                  string  `json:"markDatetime,omitempty"`
	CustomerComment               string  `json:"customerComment,omitempty"`
	ManagerComment                string  `json:"managerComment,omitempty"`
	Status                        string  `json:"status,omitempty"`
	StatusComment                 string  `json:"statusComment,omitempty"`
	FullPaidAt                    string  `json:"fullPaidAt,omitempty"`
	Site                          string  `json:"site,omitempty"`
	OrderType                     string  `json:"orderType,omitempty"`
	OrderMethod                   string  `json:"orderMethod,omitempty"`
	CountryIso                    string  `json:"countryIso,omitempty"`
	Summ                          float32 `json:"summ,omitempty"`
	TotalSumm                     float32 `json:"totalSumm,omitempty"`
	PrepaySum                     float32 `json:"prepaySum,omitempty"`
	PurchaseSumm                  float32 `json:"purchaseSumm,omitempty"`
	Weight                        float32 `json:"weight,omitempty"`
	Length                        int     `json:"length,omitempty"`
	Width                         int     `json:"width,omitempty"`
	Height                        int     `json:"height,omitempty"`
	ShipmentStore                 string  `json:"shipmentStore,omitempty"`
	ShipmentDate                  string  `json:"shipmentDate,omitempty"`
	ClientId                      string  `json:"clientId,omitempty"`
	Shipped                       bool    `json:"shipped,omitempty"`
	UploadedToExternalStoreSystem bool    `json:"uploadedToExternalStoreSystem,omitempty"`

	Source      *Source           `json:"source,omitempty"`
	Contragent  *Contragent       `json:"contragent,omitempty"`
	Customer    *Customer         `json:"customer,omitempty"`
	Delivery    *OrderDelivery    `json:"delivery,omitempty"`
	Marketplace *OrderMarketplace `json:"marketplace,omitempty"`
	Items       []OrderItem       `json:"items,omitempty"`
	Payments    []OrderPayments   `json:"payments,omitempty"`
}
