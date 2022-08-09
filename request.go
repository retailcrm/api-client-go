package retailcrm

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// CustomerRequest type.
type CustomerRequest struct {
	By   string `url:"by,omitempty"`
	Site string `url:"site,omitempty"`
}

// CustomersRequest type.
type CustomersRequest struct {
	Filter CustomersFilter `url:"filter,omitempty"`
	Limit  int             `url:"limit,omitempty"`
	Page   int             `url:"page,omitempty"`
}

// CorporateCustomersRequest type.
type CorporateCustomersRequest struct {
	Filter CorporateCustomersFilter `url:"filter,omitempty"`
	Limit  int                      `url:"limit,omitempty"`
	Page   int                      `url:"page,omitempty"`
}

// CorporateCustomersNotesRequest type.
type CorporateCustomersNotesRequest struct {
	Filter CorporateCustomersNotesFilter `url:"filter,omitempty"`
	Limit  int                           `url:"limit,omitempty"`
	Page   int                           `url:"page,omitempty"`
}

// CorporateCustomerAddressesRequest type.
type CorporateCustomerAddressesRequest struct {
	Filter CorporateCustomerAddressesFilter `url:"filter,omitempty"`
	By     string                           `url:"by,omitempty"`
	Site   string                           `url:"site,omitempty"`
	Limit  int                              `url:"limit,omitempty"`
	Page   int                              `url:"page,omitempty"`
}

// IdentifiersPairRequest type.
type IdentifiersPairRequest struct {
	Filter IdentifiersPairFilter `url:"filter,omitempty"`
	By     string                `url:"by,omitempty"`
	Site   string                `url:"site,omitempty"`
	Limit  int                   `url:"limit,omitempty"`
	Page   int                   `url:"page,omitempty"`
}

// CustomersUploadRequest type.
type CustomersUploadRequest struct {
	Customers []Customer `url:"customers,omitempty,brackets"`
	Site      string     `url:"site,omitempty"`
}

// CustomersHistoryRequest type.
type CustomersHistoryRequest struct {
	Filter CustomersHistoryFilter `url:"filter,omitempty"`
	Limit  int                    `url:"limit,omitempty"`
	Page   int                    `url:"page,omitempty"`
}

// CorporateCustomersHistoryRequest type.
type CorporateCustomersHistoryRequest struct {
	Filter CorporateCustomersHistoryFilter `url:"filter,omitempty"`
	Limit  int                             `url:"limit,omitempty"`
	Page   int                             `url:"page,omitempty"`
}

// OrderRequest type.
type OrderRequest struct {
	By   string `url:"by,omitempty"`
	Site string `url:"site,omitempty"`
}

// OrdersRequest type.
type OrdersRequest struct {
	Filter OrdersFilter `url:"filter,omitempty"`
	Limit  int          `url:"limit,omitempty"`
	Page   int          `url:"page,omitempty"`
}

// OrdersStatusesRequest type.
type OrdersStatusesRequest struct {
	IDs         []int    `url:"ids,omitempty,brackets"`
	ExternalIDs []string `url:"externalIds,omitempty,brackets"`
}

// OrdersUploadRequest type.
type OrdersUploadRequest struct {
	Orders []Order `url:"orders,omitempty,brackets"`
	Site   string  `url:"site,omitempty"`
}

// OrdersHistoryRequest type.
type OrdersHistoryRequest struct {
	Filter OrdersHistoryFilter `url:"filter,omitempty"`
	Limit  int                 `url:"limit,omitempty"`
	Page   int                 `url:"page,omitempty"`
}

// PacksRequest type.
type PacksRequest struct {
	Filter PacksFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// PacksHistoryRequest type.
type PacksHistoryRequest struct {
	Filter OrdersHistoryFilter `url:"filter,omitempty"`
	Limit  int                 `url:"limit,omitempty"`
	Page   int                 `url:"page,omitempty"`
}

// UsersRequest type.
type UsersRequest struct {
	Filter UsersFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// UserGroupsRequest type.
type UserGroupsRequest struct {
	Limit int `url:"limit,omitempty"`
	Page  int `url:"page,omitempty"`
}

// TasksRequest type.
type TasksRequest struct {
	Filter TasksFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// NotesRequest type.
type NotesRequest struct {
	Filter NotesFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// SegmentsRequest type.
type SegmentsRequest struct {
	Filter SegmentsFilter `url:"filter,omitempty"`
	Limit  int            `url:"limit,omitempty"`
	Page   int            `url:"page,omitempty"`
}

// InventoriesRequest type.
type InventoriesRequest struct {
	Filter InventoriesFilter `url:"filter,omitempty"`
	Limit  int               `url:"limit,omitempty"`
	Page   int               `url:"page,omitempty"`
}

// ProductsGroupsRequest type.
type ProductsGroupsRequest struct {
	Filter ProductsGroupsFilter `url:"filter,omitempty"`
	Limit  int                  `url:"limit,omitempty"`
	Page   int                  `url:"page,omitempty"`
}

// ProductsRequest type.
type ProductsRequest struct {
	Filter ProductsFilter `url:"filter,omitempty"`
	Limit  int            `url:"limit,omitempty"`
	Page   int            `url:"page,omitempty"`
}

// ProductsPropertiesRequest type.
type ProductsPropertiesRequest struct {
	Filter ProductsPropertiesFilter `url:"filter,omitempty"`
	Limit  int                      `url:"limit,omitempty"`
	Page   int                      `url:"page,omitempty"`
}

// DeliveryTrackingRequest type.
type DeliveryTrackingRequest struct {
	DeliveryID  string                  `json:"deliveryId,omitempty"`
	TrackNumber string                  `json:"trackNumber,omitempty"`
	History     []DeliveryHistoryRecord `json:"history,omitempty"`
	ExtraData   map[string]string       `json:"extraData,omitempty"`
}

// DeliveryShipmentsRequest type.
type DeliveryShipmentsRequest struct {
	Filter ShipmentFilter `url:"filter,omitempty"`
	Limit  int            `url:"limit,omitempty"`
	Page   int            `url:"page,omitempty"`
}

// CostsRequest type.
type CostsRequest struct {
	Filter CostsFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// FilesRequest type.
type FilesRequest struct {
	Filter FilesFilter `url:"filter,omitempty"`
	Limit  int         `url:"limit,omitempty"`
	Page   int         `url:"page,omitempty"`
}

// CustomFieldsRequest type.
type CustomFieldsRequest struct {
	Filter CustomFieldsFilter `url:"filter,omitempty"`
	Limit  int                `url:"limit,omitempty"`
	Page   int                `url:"page,omitempty"`
}

// CustomDictionariesRequest type.
type CustomDictionariesRequest struct {
	Filter CustomDictionariesFilter `url:"filter,omitempty"`
	Limit  int                      `url:"limit,omitempty"`
	Page   int                      `url:"page,omitempty"`
}

// ConnectRequest contains information about the system connection that is requested to be created.
type ConnectRequest struct {
	// Token is used to verify the request. Do not use directly; use Verify instead.
	Token string `json:"token"`
	// APIKey that was generated for the module.
	APIKey string `json:"apiKey"`
	// URL of the system. Do not use directly; use SystemURL instead.
	URL string `json:"systemUrl"`
}

// BonusOperationsRequest type.
type BonusOperationsRequest struct {
	Filter BonusOperationsFilter `url:"filter,omitempty"`
	Limit  int                   `url:"limit,omitempty"`
	Cursor string                `url:"cursor,omitempty"`
}

// AccountBonusOperationsRequest type.
type AccountBonusOperationsRequest struct {
	Filter AccountBonusOperationsFilter `url:"filter,omitempty"`
	Limit  int                          `url:"limit,omitempty"`
	Page   int                          `url:"page,omitempty"`
}

// SystemURL returns system URL from the connection request without trailing slash.
func (r ConnectRequest) SystemURL() string {
	if r.URL == "" {
		return ""
	}

	if r.URL[len(r.URL)-1:] == "/" {
		return r.URL[:len(r.URL)-1]
	}

	return r.URL
}

// Verify returns true if connection request is legitimate. Application secret should be provided to this method.
func (r ConnectRequest) Verify(secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(r.APIKey)); err != nil {
		panic(err)
	}
	return hmac.Equal([]byte(r.Token), []byte(hex.EncodeToString(mac.Sum(nil))))
}
