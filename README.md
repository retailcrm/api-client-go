[![Build Status](https://img.shields.io/travis/retailcrm/api-client-go/master.svg?style=flat-square)](https://travis-ci.org/retailcrm/api-client-go)

# retailCRM API Go client

Go client for [retailCRM API](http://www.retailcrm.pro/docs/Developers/ApiVersion5).

## Installation

```bash
go get -x github.com/retailcrm/api-client-go
```

## Usage

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/retailcrm/api-client-go/v5"
)

func main() {
	var client = v5.New("https://demo.retailcrm.pro", "09jIJ09j0JKhgyfvyuUIKhiugF")

	data, status, err := client.Orders(v5.OrdersRequest{
		Filter: v5.OrdersFilter{},
		Limit: 20,
		Page: 1,
	})
	if err.ErrorMsg != "" {
		fmt.Printf("%v", err.ErrorMsg)
	}

	if status >= http.StatusBadRequest {
		fmt.Printf("%v", err.ErrorMsg)
	}

	for _, value := range data.Orders {
		fmt.Printf("%v\n", value.Email)
	}

	fmt.Println(data.Orders[1].FirstName)

	idata, status, err := c.InventoriesUpload(
        []InventoryUpload{
            {
                XMLID: "pTKIKAeghYzX21HTdzFCe1",
                Stores: []InventoryUploadStore{
                    {Code: "test-store-v5", Available: 10, PurchasePrice: 1500},
                    {Code: "test-store-v4", Available: 20, PurchasePrice: 1530},
                    {Code: "test-store", Available: 30, PurchasePrice: 1510},
                },
            },
            {
                XMLID: "JQIvcrCtiSpOV3AAfMiQB3",
                Stores: []InventoryUploadStore{
                    {Code: "test-store-v5", Available: 45, PurchasePrice: 1500},
                    {Code: "test-store-v4", Available: 32, PurchasePrice: 1530},
                    {Code: "test-store", Available: 46, PurchasePrice: 1510},
                },
            },
        },
    )
    if err.ErrorMsg != "" {
        fmt.Printf("%v", err.ErrorMsg)
    }

    if status >= http.StatusBadRequest {
        fmt.Printf("%v", err.ErrorMsg)
    }

    fmt.Println(idata.processedOffersCount)
}
```

## Testing

```bash
export RETAILCRM_URL="https://demo.retailcrm.pro"
export RETAILCRM_KEY="09jIJ09j0JKhgyfvyuUIKhiugF"
export RETAILCRM_USER="1"

cd $GOPATH/src/github.com/retailcrm/api-client-go

go test -v ./...

```
