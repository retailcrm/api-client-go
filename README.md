[![Build Status](https://github.com/retailcrm/api-client-go/workflows/ci/badge.svg)](https://github.com/retailcrm/api-client-go/actions)
[![Covarage](https://img.shields.io/codecov/c/gh/retailcrm/api-client-go/master.svg?logo=codecov&logoColor=white)](https://codecov.io/gh/retailcrm/api-client-go)
[![GitHub release](https://img.shields.io/github/release/retailcrm/api-client-go.svg?logo=github&logoColor=white)](https://github.com/retailcrm/api-client-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/retailcrm/api-client-go)](https://goreportcard.com/report/github.com/retailcrm/api-client-go)
[![GoLang version](https://img.shields.io/badge/go->=1.8-blue)](https://golang.org/dl/)
[![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/retailcrm/api-client-go)


# RetailCRM API Go client

This is golang RetailCRM API client.

## Install

```bash
go get -u github.com/retailcrm/api-client-go/v2
```

## Usage

Example:

```go
package main

import (
	"log"

	"github.com/retailcrm/api-client-go/v2"
)

func main() {
	var client = retailcrm.New("https://demo.retailcrm.pro", "09jIJ09j0JKhgyfvyuUIKhiugF")

	data, status, err := client.Orders(retailcrm.OrdersRequest{
		Filter: retailcrm.OrdersFilter{},
		Limit: 20,
		Page: 1,
	})
	if err != nil {
		if apiErr, ok := retailcrm.AsAPIError(err); ok {
			log.Fatalf("http status: %d, %s", status, apiErr.String())
        }

		log.Fatalf("http status: %d, error: %s", status, err)
	}

	for _, value := range data.Orders {
		log.Printf("%v\n", value.Email)
	}

	log.Println(data.Orders[1].FirstName)

	inventories, status, err := client.InventoriesUpload([]retailcrm.InventoryUpload{
			{
				XMLID: "pTKIKAeghYzX21HTdzFCe1",
				Stores: []retailcrm.InventoryUploadStore{
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
				Stores: []retailcrm.InventoryUploadStore{
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
		if apiErr, ok := retailcrm.AsAPIError(err); ok {
			log.Fatalf("http status: %d, %s", status, apiErr.String())
		}

		log.Fatalf("http status: %d, error: %s", status, err)
	}

	log.Println(inventories.ProcessedOffersCount)
}
```

You can use different error types and `retailcrm.AsAPIError` to process client errors. Example:

```go
package retailcrm

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/retailcrm/api-client-go/v2"
)

func main() {
	var client = retailcrm.New("https://demo.retailcrm.pro", "09jIJ09j0JKhgyfvyuUIKhiugF")

	resp, status, err := client.APICredentials()
	if err != nil {
		apiErr, ok := retailcrm.AsAPIError(err)
		if !ok {
			log.Fatalf("http status: %d, error: %s", status, err)
		}

		switch {
		case errors.Is(apiErr, retailcrm.ErrMissingCredentials):
			log.Fatalln("No API key provided.")
		case errors.Is(apiErr, retailcrm.ErrInvalidCredentials):
			log.Fatalln("Invalid API key.")
		case errors.Is(apiErr, retailcrm.ErrAccessDenied):
			log.Fatalln("Access denied. Please check that the provided key has access to the credentials info.")
		case errors.Is(apiErr, retailcrm.ErrAccountDoesNotExist):
			log.Fatalln("There is no RetailCRM at the provided URL.")
		case errors.Is(apiErr, retailcrm.ErrMissingParameter):
			// retailcrm.APIError in this case will always contain "Name" key in the errors list with the parameter name.
			log.Fatalln("This parameter should be present:", apiErr.Errors()["Name"])
		case errors.Is(apiErr, retailcrm.ErrValidation):
			log.Println("Validation errors from the API:")

			for name, value := range apiErr.Errors() {
				log.Printf(" - %s: %s\n", name, value)
			}

			os.Exit(1)
		case errors.Is(apiErr, retailcrm.ErrGeneric):
			log.Fatalf("failure from the API. %s", apiErr.String())
		}
	}

	log.Println("Available scopes:", strings.Join(resp.Scopes, ", "))
}
```
