# retailCRM API Go client

Go client for [retailCRM API](http://www.retailcrm.pro/docs/Developers/ApiVersion5).

## Installation

```bash
go get -x github.com/retailcrm/api-client-go
```

## Usage

```golang
import (
	c "github.com/retailcrm/api-client-go"
)

var client = c.Version5("https://demo.retailcrm.ru", "09jIJ09j0JKhgyfvyuUIKhiugF")

data, status, err := c.Customers(CustomersRequest{
    Filter: CustomersFilter{
        MinCostSumm: 500,
    },
    Page: 2,
})

if err != nil {
    t.Errorf("%s", err)
    t.Fail()
}

if status >= http.StatusBadRequest {
    t.Errorf("%s", err)
    t.Fail()
}

var email = data.Customers[0].Email
```

## Testing

```bash
export RETAILCRM_URL="https://demo.retailcrm.ru"
export RETAILCRM_KEY="09jIJ09j0JKhgyfvyuUIKhiugF"
export RETAILCRM_USER="1"

cd $GOPATH/src/github.com/retailcrm/api-client-go

go test -v ./...

```