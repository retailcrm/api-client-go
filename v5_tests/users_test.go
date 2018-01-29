package v5_tests

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

var user, _ = strconv.Atoi(os.Getenv("RETAILCRM_USER"))

func TestClient_UsersUsers(t *testing.T) {
	c := client()
	f := v5.UsersRequest{
		Filter: v5.UsersFilter{
			Active: 1,
		},
		Page: 1,
	}

	data, status, err := c.Users(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_UsersUser(t *testing.T) {
	c := client()

	data, st, err := c.User(user)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_UsersGroups(t *testing.T) {
	c := client()
	f := v5.UserGroupsRequest{
		Page: 1,
	}

	data, status, err := c.UserGroups(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_UsersUpdate(t *testing.T) {
	c := client()

	data, st, err := c.UserStatus(user, "busy")
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}
