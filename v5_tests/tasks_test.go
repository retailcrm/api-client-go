package v5_tests

import (
	"net/http"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

func TestClient_TasksTasks(t *testing.T) {
	c := client()

	f := v5.TasksRequest{
		Filter: v5.TasksFilter{
			Creators: []int{GetUser()},
		},
		Page: 1,
	}

	data, status, err := c.Tasks(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_TaskChange(t *testing.T) {
	c := client()

	f := v5.Task{
		Text:        RandomString(15),
		PerformerId: GetUser(),
	}

	cr, sc, err := c.TaskCreate(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%s", err)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	f.Id = cr.Id
	f.Commentary = RandomString(20)

	gt, sg, err := c.Task(f.Id)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if sg != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if gt.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	data, status, err := c.TaskEdit(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}
