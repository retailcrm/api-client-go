# Upgrading to the v2

### Install the new version

```bash
go get -u github.com/retailcrm/api-client-go/v2
```

### Update all imports

Before:
```go
package main

import v5 "github.com/retailcrm/api-client-go/v5"
```

After:  
```go
package main

import "github.com/retailcrm/api-client-go/v2"
```

You can use package alias `v5` to skip the second step.

### Replace package name for all imported symbols

Before:

```go
package main

import v5 "github.com/retailcrm/api-client-go/v5"

func main() {
    client := v5.New("https://test.retailcrm.pro", "key")
	data, status, err := client.Orders(v5.OrdersRequest{
		Filter: v5.OrdersFilter{
			City: "Moscow",
		},
		Page: 1,
	})
	...
}
```

After:

```go
package main

import "github.com/retailcrm/api-client-go/v2"

func main() {
    client := retailcrm.New("https://test.retailcrm.pro", "key")
	data, status, err := client.Orders(retailcrm.OrdersRequest{
		Filter: retailcrm.OrdersFilter{
			City: "Moscow",
		},
		Page: 1,
	})
	...
}
```

### Upgrade client usages

This major release contains some breaking changes regarding field names and fully redesigned error handling. Use the second example from 
the readme to learn how to process errors correctly.


## WithContext now returns a clone

Always use the return value of `WithContext`. It no longer modifies the original client:

```go
// Before:
client.WithContext(ctx)
response, status, err := client.Customers(retailcrm.CustomersRequest{})

// After:
response, status, err := client.WithContext(ctx).Customers(retailcrm.CustomersRequest{})
// Or save the clone for multiple calls:
requestClient := client.WithContext(ctx)
```

The clone uses the context for HTTP requests, response body reads, rate limiting and retry waits.
Cancellation no longer affects other clones or the original client.
Configure the client before cloning: later configuration changes on the original do not propagate to existing clones.
The HTTP client, logger and rate limiter are shared; cloning does not reset rate limits.
Do not modify client settings concurrently with cloning or requests.
