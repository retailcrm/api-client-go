[![Build Status](https://img.shields.io/travis/retailcrm/api-client-go/master.svg?style=flat-square)](https://travis-ci.org/retailcrm/api-client-go)
[![Covarage](https://img.shields.io/codecov/c/gh/retailcrm/api-client-go/master.svg?style=flat-square)](https://codecov.io/gh/retailcrm/api-client-go)
[![GitHub release](https://img.shields.io/github/release/retailcrm/api-client-go.svg?style=flat-square)](https://github.com/retailcrm/api-client-go/releases)
[![GoLang version](https://img.shields.io/badge/go->=1.8-blue?style=flat-square)](https://golang.org/dl/)
[![Godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/retailcrm/api-client-go)


# retailCRM API Go client

This is golang retailCRM API client.

## Install

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
	},)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		fmt.Printf("%v", err.ApiError())
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
					{
						Code: "test-store-v5",
						Available: 10,
						PurchasePrice: 1500,
					},
					{
						Code: "test-store-v4",
						Available: 20,
						PurchasePrice: 1530,
					},
					{
						Code: "test-store",
						Available: 30,
						PurchasePrice: 1510,
					},
				},
			},
			{
				XMLID: "JQIvcrCtiSpOV3AAfMiQB3",
				Stores: []InventoryUploadStore{
					{
						Code: "test-store-v5",
						Available: 45,
						PurchasePrice: 1500,
					},
					{
						Code: "test-store-v4",
						Available: 32,
						PurchasePrice: 1530,
					},
					{
						Code: "test-store",
						Available: 46,
						PurchasePrice: 1510,
					},
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		fmt.Printf("%v", err.ApiError())
	}

	fmt.Println(idata.processedOffersCount)
}
```

## Documentation

* [English](http://www.retailcrm.pro/docs/Developers/Index)
* [Russian](http://www.retailcrm.ru/docs/Developers/Index)
