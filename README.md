[![Build Status](https://github.com/retailcrm/api-client-go/workflows/ci/badge.svg)](https://github.com/retailcrm/api-client-go/actions)
[![Covarage](https://img.shields.io/codecov/c/gh/retailcrm/api-client-go/master.svg?logo=codecov&logoColor=white)](https://codecov.io/gh/retailcrm/api-client-go)
[![GitHub release](https://img.shields.io/github/release/retailcrm/api-client-go.svg?logo=github&logoColor=white)](https://github.com/retailcrm/api-client-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/retailcrm/api-client-go)](https://goreportcard.com/report/github.com/retailcrm/api-client-go)
[![GoLang version](https://img.shields.io/badge/go->=1.8-blue)](https://golang.org/dl/)
[![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/retailcrm/api-client-go)


# RetailCRM API Go client

This is golang RetailCRM API client.

## Installation

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

## Rate limits

This client can work with default rate limits but doesn't do that unless specified explicitly. You can enable default 
rate limits like this:

```go
client := retailcrm.New("https://demo.retailcrm.pro", "09jIJ09j0JKhgyfvyuUIKhiugF").
	EnableRateLimiter(0)
```

This `client` will automatically apply rate limiting. Requests will block until they are finished. You can also provide 
different retry parameter instead of `0` if you don't want the client to block completely and would like to process 
rate limit by yourself after several attempts.

Custom rate limiter can be provided like this:

```go
limiter := retailcrm.NewSingleKeyLimiter()
client := retailcrm.New("https://demo.retailcrm.pro", "09jIJ09j0JKhgyfvyuUIKhiugF").
	EnableCustomRateLimiter(limiter, 0)
```

You can use your own version of limiter by implementing `retailcrm.Limiter`. Also, any instance of 
`retailcrm.Limiter` which implements `retailcrm.ResponseAware` interface will be able to read response for each 
request attempt (`(*http.Response).Body` is not guaranteed to be accessible). This feature can be used to control 
rate limits for distributed applications using the same key.

## Upgrading

Please check the [UPGRADING.md](UPGRADING.md) to learn how to upgrade to the new version.

## Knowledge base

The API key needs `knowledge_base_read` / `knowledge_base_write` scopes.
Sections and documents are account-wide and do not take a `site` parameter.

```go
client := retailcrm.New("https://example.retailcrm.ru", "apiKey")
title, content, externalID := "Shipping", "# Shipping options", "shipping"
created, status, err := client.KnowledgeBaseDocumentCreate(retailcrm.KnowledgeBaseDocumentInput{
    Title: &title, Content: &content, ExternalID: &externalID,
})
if err != nil {
    log.Fatalf("create failed (%d): %v", status, err)
}

// by="id" selects the internal ID; by="externalId" or by="" selects the external ID.
document, status, err := client.KnowledgeBaseDocument(strconv.Itoa(created.ID), "id")
if err != nil {
    log.Fatalf("read failed (%d): %v", status, err)
}
log.Println(document.Document.Title)

// Non-nil pointers send the value even when it is false, zero or an empty string.
disabled := false
_, status, err = client.KnowledgeBaseDocumentEdit("shipping", retailcrm.KnowledgeBaseDocumentInput{
    AvailableForCopilot: &disabled,
}, "externalId")
if err != nil {
    log.Fatalf("edit failed (%d): %v", status, err)
}
```

The example uses the standard `log` and `strconv` packages.
Nil input fields are omitted. Set `SectionID` or `ParentID` to a pointer to zero to move a document
or section to the root; use a pointer to an empty string to clear content.
Set `ClearExternalID: true` to send `externalId: null` and release the external identifier.

Methods:

- `KnowledgeBaseSections`, `KnowledgeBaseSectionCreate`, `KnowledgeBaseSectionEdit`, `KnowledgeBaseSectionDelete`.
- `KnowledgeBaseDocuments`, `KnowledgeBaseDocument`, `KnowledgeBaseDocumentCreate`,
  `KnowledgeBaseDocumentEdit`, `KnowledgeBaseDocumentDelete`, `KnowledgeBaseDocumentsHistory`.

Sections have no pagination. Document lists accept `KnowledgeBaseDocumentsRequest` with pagination
and filters (`IDs`, `ExternalIDs`, `SectionID`, `DirectOnly`, `AvailableForCopilot`, `HasFile`, `Search`).
`DirectOnly` selects direct children; without `SectionID`, it selects root documents.
History accepts `KnowledgeBaseDocumentsHistoryRequest`. Store the last record ID and pass it as
`Filter.SinceID` to retrieve subsequent changes. Content history contains hashes, not the text.

To attach an original file, use the existing `FileUpload` and `FileEdit` methods and set
`Attachment.KnowledgeBaseDocument` to `&retailcrm.KnowledgeBaseDocumentReference{ID: documentID}`.
The file API needs the corresponding file scopes. Each document supports one original file;
its URL is available in `Document.File.URL` when `Document.File != nil`.

Deleting a section also deletes **all its subsections and documents**.
