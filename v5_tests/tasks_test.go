package v5_tests

import (
	"net/http"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

func TestClient_TasksTask(t *testing.T) {
	c := client()

	data, st, err := c.Task(88)
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

func TestClient_TasksTasks(t *testing.T) {
	c := client()
	f := v5.TasksRequest{
		Filter: v5.TasksFilter{
			Creators: []int{6},
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

	random1 := RandomString(15)
	random2 := RandomString(20)

	f := v5.Task{
		Text:        random1,
		PerformerId: 6,
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
	f.Commentary = random2

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
