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
	if err != nil {
		fmt.Printf("%v", err)
	}

	if status >= http.StatusBadRequest {
		fmt.Printf("%v", err)
	}

	for _, value := range data.Orders {
		fmt.Printf("%v\n", value.Email)
	}

	fmt.Println(data.Orders[1].FirstName)
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