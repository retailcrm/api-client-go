package v5

// CustomerRequest type
type CustomerRequest struct {
	By   string `url:"by,omitempty"`
	Site string `url:"site,omitempty"`
}

// CustomersRequest type
type CustomersRequest struct {
	Filter CustomersFilter `url:"filter,omitempty"`
	Limit  int             `url:"limit,omitempty"`
	Page   int             `url:"page,omitempty"`
}

// CustomersUploadRequest type
type CustomersUploadRequest struct {
	Customers []Customer `url:"customers,omitempty,brackets"`
	Site      string     `url:"site,omitempty"`
}

// CustomersHistoryRequest type
type CustomersHistoryRequest struct {
	Filter CustomersHistoryFilter `url:"filter,omitempty"`
	Limit  int                    `url:"limit,omitempty"`
	Page   int                    `url:"page,omitempty"`
}

// OrderRequest type
type OrderRequest struct {
	By   string `url:"by,omitempty"`
	Site string `url:"site,omitempty"`
}

// OrdersRequest type
type OrdersRequest struct {
	Filter OrdersFilter `url:"filter,omitempty"`
	Limit  int          `url:"limit,omitempty"`
	Page   int          `url:"page,omitempty"`
}

// CustomersUploadRequest type
type OrdersUploadRequest struct {
	Customers []Order `url:"orders,omitempty,brackets"`
	Site      string  `url:"site,omitempty"`
}

// OrdersHistoryRequest type
type OrdersHistoryRequest struct {
	Filter OrdersHistoryFilter `url:"filter,omitempty"`
	Limit  int                 `url:"limit,omitempty"`
	Page   int                 `url:"page,omitempty"`
}
