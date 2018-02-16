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

// OrdersUploadRequest type
type OrdersUploadRequest struct {
	Orders []Order `url:"orders,omitempty,brackets"`
	Site   string  `url:"site,omitempty"`
}

// OrdersHistoryRequest type
type OrdersHistoryRequest struct {
	Filter OrdersHistoryFilter `url:"filter,omitempty"`
	Limit  int                 `url:"limit,omitempty"`
	Page   int                 `url:"page,omitempty"`
}

// PacksRequest type
type PacksRequest struct {
	Filter PacksFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// PacksHistoryRequest type
type PacksHistoryRequest struct {
	Filter OrdersHistoryFilter `url:"filter,omitempty"`
	Limit  int                 `url:"limit,omitempty"`
	Page   int                 `url:"page,omitempty"`
}

// UsersRequest type
type UsersRequest struct {
	Filter UsersFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// UserGroupsRequest type
type UserGroupsRequest struct {
	Limit int `url:"limit,omitempty"`
	Page  int `url:"page,omitempty"`
}

// TasksRequest type
type TasksRequest struct {
	Filter TasksFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// NotesRequest type
type NotesRequest struct {
	Filter TasksFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// SegmentsRequest type
type SegmentsRequest struct {
	Filter SegmentsFilter `url:"filter,omitempty"`
	Limit  int            `url:"limit,omitempty"`
	Page   int            `url:"page,omitempty"`
}

// InventoriesRequest type
type InventoriesRequest struct {
	Filter InventoriesFilter `url:"filter,omitempty"`
	Limit  int               `url:"limit,omitempty"`
	Page   int               `url:"page,omitempty"`
}

// ProductsGroupsRequest type
type ProductsGroupsRequest struct {
	Filter ProductsGroupsFilter `url:"filter,omitempty"`
	Limit  int                  `url:"limit,omitempty"`
	Page   int                  `url:"page,omitempty"`
}

// ProductsRequest type
type ProductsRequest struct {
	Filter ProductsFilter `url:"filter,omitempty"`
	Limit  int            `url:"limit,omitempty"`
	Page   int            `url:"page,omitempty"`
}

// ProductsPropertiesRequest type
type ProductsPropertiesRequest struct {
	Filter ProductsPropertiesFilter `url:"filter,omitempty"`
	Limit  int                      `url:"limit,omitempty"`
	Page   int                      `url:"page,omitempty"`
}

// DeliveryTrackingRequest type
type DeliveryTrackingRequest struct {
	DeliveryId  string                  `url:"deliveryId,omitempty"`
	TrackNumber string                  `url:"trackNumber,omitempty"`
	History     []DeliveryHistoryRecord `url:"history,omitempty,brackets"`
	ExtraData   map[string]string       `url:"extraData,omitempty,brackets"`
}

// DeliveryShipmentsRequest type
type DeliveryShipmentsRequest struct {
	Filter ShipmentFilter `url:"filter,omitempty"`
	Limit  int            `url:"limit,omitempty"`
	Page   int            `url:"page,omitempty"`
}
