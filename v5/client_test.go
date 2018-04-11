package v5

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

var r *rand.Rand // Rand for this package.
var user, _ = strconv.Atoi(os.Getenv("RETAILCRM_USER"))
var statuses = map[int]bool{http.StatusOK: true, http.StatusCreated: true}
var id int
var ids []int

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

func client() *Client {
	return New(os.Getenv("RETAILCRM_URL"), os.Getenv("RETAILCRM_KEY"))
}

func badurlclient() *Client {
	return New("https://qwertypoiu.retailcrm.ru", os.Getenv("RETAILCRM_KEY"))
}

func badkeyclient() *Client {
	return New(os.Getenv("RETAILCRM_URL"), "1234567890")
}

func TestGetRequest(t *testing.T) {
	c := client()
	_, status, _ := c.GetRequest("/fake-method")

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestPostRequest(t *testing.T) {
	c := client()
	_, status, _ := c.PostRequest("/fake-method", url.Values{})

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestClient_ApiVersionsVersions(t *testing.T) {
	c := client()

	data, status, err := c.APIVersions()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ApiVersionsVersionsBadUrl(t *testing.T) {
	c := badurlclient()

	data, status, err := c.APIVersions()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Logf("%v", err.ApiError())
	}
}

func TestClient_ApiVersionsVersionsBadKey(t *testing.T) {
	c := badkeyclient()

	data, status, err := c.APIVersions()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Logf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Logf("%v", err.ApiError())
	}
}

func TestClient_ApiCredentialsCredentials(t *testing.T) {
	c := client()

	data, status, err := c.APICredentials()
	if err.RuntimeErr != nil {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersCustomers(t *testing.T) {
	c := client()

	data, status, err := c.Customers(CustomersRequest{
		Filter: CustomersFilter{
			City: "Москва",
		},
		Page: 3,
	})

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomerChange(t *testing.T) {
	c := client()

	f := Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		Address: &Address{
			City:     "Москва",
			Street:   "Кутузовский",
			Building: "14",
		},
	}

	cr, sc, err := c.CustomerCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f.ID = cr.ID
	f.Vip = true

	ed, se, err := c.CustomerEdit(f, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if se != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if ed.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	data, status, err := c.Customer(f.ExternalID, "externalId", "")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if data.Customer.ExternalID != f.ExternalID {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersUpload(t *testing.T) {
	c := client()
	customers := make([]Customer, 3)

	for i := range customers {
		customers[i] = Customer{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalID: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	data, status, err := c.CustomersUpload(customers)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersCombine(t *testing.T) {
	c := client()

	dataFirst, status, err := c.CustomerCreate(Customer{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if dataFirst.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	dataSecond, status, err := c.CustomerCreate(Customer{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if dataSecond.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	dataThird, status, err := c.CustomerCreate(Customer{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if dataThird.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	data, status, err := c.CustomersCombine([]Customer{{ID: dataFirst.ID}, {ID: dataSecond.ID}}, Customer{ID: dataThird.ID})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersFixExternalIds(t *testing.T) {
	c := client()
	f := Customer{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.CustomerCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	customers := []IdentifiersPair{{
		ID:         cr.ID,
		ExternalID: RandomString(8),
	}}

	fx, fe, err := c.CustomersFixExternalIds(customers)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if fx.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersHistory(t *testing.T) {
	c := client()
	f := CustomersHistoryRequest{
		Filter: CustomersHistoryFilter{
			SinceID: 20,
		},
	}

	data, status, err := c.CustomersHistory(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.History) == 0 {
		t.Errorf("%v", "Empty history")
	}
}

func TestClient_NotesNotes(t *testing.T) {
	c := client()

	data, status, err := c.CustomerNotes(NotesRequest{Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_NotesCreateDelete(t *testing.T) {
	c := client()

	createCustomerResponse, createCustomerStatus, err := c.CustomerCreate(Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if createCustomerStatus != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if createCustomerResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	noteCreateResponse, noteCreateStatus, err := c.CustomerNoteCreate(Note{
		Text:      "some text",
		ManagerID: user,
		Customer: &Customer{
			ID: createCustomerResponse.ID,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if noteCreateStatus != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if noteCreateResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	noteDeleteResponse, noteDeleteStatus, err := c.CustomerNoteDelete(noteCreateResponse.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if noteDeleteStatus != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if noteDeleteResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersOrders(t *testing.T) {
	c := client()

	data, status, err := c.Orders(OrdersRequest{Filter: OrdersFilter{City: "Москва"}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderChange(t *testing.T) {
	c := client()

	random := RandomString(8)

	f := Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: random,
		Email:      fmt.Sprintf("%s@example.com", random),
	}

	cr, sc, err := c.OrderCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f.ID = cr.ID
	f.CustomerComment = "test comment"

	ed, se, err := c.OrderEdit(f, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if se != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if ed.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	data, status, err := c.Order(f.ExternalID, "externalId", "")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if data.Order.ExternalID != f.ExternalID {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersUpload(t *testing.T) {
	c := client()
	orders := make([]Order, 3)

	for i := range orders {
		orders[i] = Order{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalID: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	data, status, err := c.OrdersUpload(orders)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersCombine(t *testing.T) {
	c := client()

	dataFirst, status, err := c.OrderCreate(Order{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if dataFirst.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	dataSecond, status, err := c.OrderCreate(Order{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if dataSecond.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	data, status, err := c.OrdersCombine("ours", Order{ID: dataFirst.ID}, Order{ID: dataSecond.ID})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersFixExternalIds(t *testing.T) {
	c := client()
	f := Order{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.OrderCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	orders := []IdentifiersPair{{
		ID:         cr.ID,
		ExternalID: RandomString(8),
	}}

	fx, fe, err := c.OrdersFixExternalIds(orders)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if fx.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersHistory(t *testing.T) {
	c := client()

	data, status, err := c.OrdersHistory(OrdersHistoryRequest{Filter: OrdersHistoryFilter{SinceID: 20}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.History) == 0 {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentCreateEditDelete(t *testing.T) {
	c := client()

	order := Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	createOrderResponse, status, err := c.OrderCreate(order)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if createOrderResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f := Payment{
		Order: &Order{
			ID: createOrderResponse.ID,
		},
		Amount: 300,
		Type:   "cash",
	}

	paymentCreateResponse, status, err := c.OrderPaymentCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if paymentCreateResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	k := Payment{
		ID:     paymentCreateResponse.ID,
		Amount: 500,
	}

	paymentEditResponse, status, err := c.OrderPaymentEdit(k, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if paymentEditResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	paymentDeleteResponse, status, err := c.OrderPaymentDelete(paymentCreateResponse.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if paymentDeleteResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_TasksTasks(t *testing.T) {
	c := client()

	f := TasksRequest{
		Filter: TasksFilter{
			Creators: []int{user},
		},
		Page: 1,
	}

	data, status, err := c.Tasks(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_TaskChange(t *testing.T) {
	c := client()

	f := Task{
		Text:        RandomString(15),
		PerformerID: user,
	}

	cr, sc, err := c.TaskCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f.ID = cr.ID
	f.Commentary = RandomString(20)

	gt, sg, err := c.Task(f.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sg != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if gt.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	data, status, err := c.TaskEdit(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUsers(t *testing.T) {
	c := client()

	data, status, err := c.Users(UsersRequest{Filter: UsersFilter{Active: 1}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUser(t *testing.T) {
	c := client()

	data, st, err := c.User(user)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersGroups(t *testing.T) {
	c := client()

	data, status, err := c.UserGroups(UserGroupsRequest{Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUpdate(t *testing.T) {
	c := client()

	data, st, err := c.UserStatus(user, "busy")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StaticticsUpdate(t *testing.T) {
	c := client()

	data, st, err := c.StaticticsUpdate()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Countries(t *testing.T) {
	c := client()

	data, st, err := c.Couriers()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostGroups(t *testing.T) {
	c := client()

	data, st, err := c.CostGroups()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostItems(t *testing.T) {
	c := client()

	data, st, err := c.CostItems()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Couriers(t *testing.T) {
	c := client()

	data, st, err := c.Couriers()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryService(t *testing.T) {
	c := client()

	data, st, err := c.DeliveryServices()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryTypes(t *testing.T) {
	c := client()

	data, st, err := c.DeliveryTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_LegalEntities(t *testing.T) {
	c := client()

	data, st, err := c.LegalEntities()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderMethods(t *testing.T) {
	c := client()

	data, st, err := c.OrderMethods()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderTypes(t *testing.T) {
	c := client()

	data, st, err := c.OrderTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentStatuses(t *testing.T) {
	c := client()

	data, st, err := c.PaymentStatuses()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentTypes(t *testing.T) {
	c := client()

	data, st, err := c.PaymentTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PriceTypes(t *testing.T) {
	c := client()

	data, st, err := c.PriceTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ProductStatuses(t *testing.T) {
	c := client()

	data, st, err := c.ProductStatuses()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Statuses(t *testing.T) {
	c := client()

	data, st, err := c.Statuses()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StatusGroups(t *testing.T) {
	c := client()

	data, st, err := c.StatusGroups()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Sites(t *testing.T) {
	c := client()

	data, st, err := c.Sites()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Stores(t *testing.T) {
	c := client()

	data, st, err := c.Stores()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostGroupItemEdit(t *testing.T) {
	c := client()

	uid := RandomString(4)

	data, st, err := c.CostGroupEdit(CostGroup{
		Code:   fmt.Sprintf("cost-gr-%s", uid),
		Active: false,
		Color:  "#da5c98",
		Name:   fmt.Sprintf("CostGroup-%s", uid),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	idata, st, err := c.CostItemEdit(CostItem{
		Code:            fmt.Sprintf("cost-it-%s", uid),
		Name:            fmt.Sprintf("CostItem-%s", uid),
		Group:           fmt.Sprintf("cost-gr-%s", uid),
		Type:            "const",
		AppliesToOrders: true,
		Active:          false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if idata.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Courier(t *testing.T) {
	c := client()

	cur := Courier{
		Active:    true,
		Email:     fmt.Sprintf("%s@example.com", RandomString(5)),
		FirstName: fmt.Sprintf("%s", RandomString(5)),
		LastName:  fmt.Sprintf("%s", RandomString(5)),
	}

	data, st, err := c.CourierCreate(cur)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	cur.ID = data.ID
	cur.Patronymic = fmt.Sprintf("%s", RandomString(5))

	idata, st, err := c.CourierEdit(cur)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if idata.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryServiceEdit(t *testing.T) {
	c := client()

	data, st, err := c.DeliveryServiceEdit(DeliveryService{
		Active: false,
		Code:   RandomString(5),
		Name:   RandomString(5),
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryTypeEdit(t *testing.T) {
	c := client()

	x := []string{"cash", "bank-card"}

	data, st, err := c.DeliveryTypeEdit(DeliveryType{
		Active:        false,
		Code:          RandomString(5),
		Name:          RandomString(5),
		DefaultCost:   300,
		PaymentTypes:  x,
		DefaultForCrm: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderMethodEdit(t *testing.T) {
	c := client()

	data, st, err := c.OrderMethodEdit(OrderMethod{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderTypeEdit(t *testing.T) {
	c := client()

	data, st, err := c.OrderTypeEdit(OrderType{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentStatusEdit(t *testing.T) {
	c := client()

	data, st, err := c.PaymentStatusEdit(PaymentStatus{
		Code:            RandomString(5),
		Name:            RandomString(5),
		Active:          false,
		DefaultForCRM:   false,
		PaymentTypes:    []string{"cash"},
		PaymentComplete: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentTypeEdit(t *testing.T) {
	c := client()

	data, st, err := c.PaymentTypeEdit(PaymentType{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PriceTypeEdit(t *testing.T) {
	c := client()

	data, st, err := c.PriceTypeEdit(PriceType{
		Code:   RandomString(5),
		Name:   RandomString(5),
		Active: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ProductStatusEdit(t *testing.T) {
	c := client()

	data, st, err := c.ProductStatusEdit(ProductStatus{
		Code:         RandomString(5),
		Name:         RandomString(5),
		Active:       false,
		CancelStatus: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StatusEdit(t *testing.T) {
	c := client()

	data, st, err := c.StatusEdit(Status{
		Code:   RandomString(5),
		Name:   RandomString(5),
		Active: false,
		Group:  "new",
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_SiteEdit(t *testing.T) {
	c := client()

	data, _, err := c.SiteEdit(Site{
		Code:        RandomString(5),
		Name:        RandomString(5),
		URL:         fmt.Sprintf("https://%s.example.com", RandomString(4)),
		LoadFromYml: false,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if data.Success != false {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StoreEdit(t *testing.T) {
	c := client()

	data, st, err := c.StoreEdit(Store{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		Type:          "store-type-warehouse",
		InventoryType: "integer",
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PackChange(t *testing.T) {
	c := client()

	o, status, err := c.OrderCreate(Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		Items:      []OrderItem{{Offer: Offer{ID: 1609}, Quantity: 5}},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if o.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	g, status, err := c.Order(strconv.Itoa(o.ID), "id", "")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if o.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	p, status, err := c.PackCreate(Pack{
		Store:    "test-store",
		ItemID:   g.Order.Items[0].ID,
		Quantity: 1,
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if p.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	s, status, err := c.Pack(p.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if s.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	e, status, err := c.PackEdit(Pack{ID: p.ID, Quantity: 2})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if e.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	d, status, err := c.PackDelete(p.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if d.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PacksHistory(t *testing.T) {
	c := client()

	data, status, err := c.PacksHistory(PacksHistoryRequest{Filter: OrdersHistoryFilter{SinceID: 5}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.History) == 0 {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Packs(t *testing.T) {
	c := client()

	data, status, err := c.Packs(PacksRequest{Filter: PacksFilter{}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Inventories(t *testing.T) {
	c := client()

	data, status, err := c.Inventories(InventoriesRequest{Filter: InventoriesFilter{Details: 1}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Segments(t *testing.T) {
	c := client()

	data, status, err := c.Segments(SegmentsRequest{})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_IntegrationModule(t *testing.T) {
	c := client()

	name := RandomString(5)
	code := RandomString(8)

	m, status, err := c.IntegrationModuleEdit(IntegrationModule{
		Code:            code,
		IntegrationCode: code,
		Active:          false,
		Name:            fmt.Sprintf("Integration module %s", name),
		AccountURL:      fmt.Sprintf("http://example.com/%s/account", name),
		BaseURL:         fmt.Sprintf("http://example.com/%s", name),
		ClientID:        RandomString(10),
		Logo:            "https://cdn.worldvectorlogo.com/logos/github-icon.svg",
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if m.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	g, status, err := c.IntegrationModule(code)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_IntegrationModuleFail(t *testing.T) {
	c := client()
	code := RandomString(8)

	m, status, err := c.IntegrationModuleEdit(IntegrationModule{
		Code: code,
	})
	if err.RuntimeErr == nil {
		t.Fail()
	}

	if status < http.StatusBadRequest {
		t.Fail()
	}

	if m.Success != false {
		t.Fail()
	}

	g, status, err := c.IntegrationModule(RandomString(12))
	if err.RuntimeErr == nil {
		t.Fail()
	}

	if status < http.StatusBadRequest {
		t.Fail()
	}

	if g.Success != false {
		t.Fail()
	}
}

func TestClient_ProductsGroup(t *testing.T) {
	c := client()

	g, status, err := c.ProductsGroup(ProductsGroupsRequest{
		Filter: ProductsGroupsFilter{
			Active: 1,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Products(t *testing.T) {
	c := client()

	g, status, err := c.Products(ProductsRequest{
		Filter: ProductsFilter{
			Active:   1,
			MinPrice: 1,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ProductsProperties(t *testing.T) {
	c := client()

	sites := make([]string, 1)
	sites[0] = os.Getenv("RETAILCRM_SITE")

	g, status, err := c.ProductsProperties(ProductsPropertiesRequest{
		Filter: ProductsPropertiesFilter{
			Sites: sites,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryShipments(t *testing.T) {
	c := client()

	g, status, err := c.DeliveryShipments(DeliveryShipmentsRequest{
		Filter: ShipmentFilter{
			DateFrom: "2017-10-10",
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostCreate(t *testing.T) {
	c := client()

	data, status, err := c.CostCreate(CostRecord{
		DateFrom: "2018-04-02",
		DateTo:   "2018-04-02",
		Summ:     124,
		CostItem: "seo",
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	id = data.ID
}

func TestClient_Costs(t *testing.T) {
	c := client()

	data, status, err := c.Costs(CostsRequest{
		Filter: CostsFilter{
			Ids: []string{strconv.Itoa(id)},
		},
		Limit: 20,
		Page:  1,
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_Cost(t *testing.T) {
	c := client()

	data, status, err := c.Cost(id)

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CostEdit(t *testing.T) {
	c := client()

	data, status, err := c.CostEdit(id, CostRecord{
		DateFrom: "2018-04-09",
		DateTo:   "2018-04-09",
		Summ:     421,
		CostItem: "seo",
		Order:    nil,
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CostDelete(t *testing.T) {
	c := client()

	data, status, err := c.CostDelete(id)

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CostsUpload(t *testing.T) {
	c := client()

	data, status, err := c.CostsUpload([]CostRecord{
		{
			Source:   nil,
			DateFrom: "2018-04-02",
			DateTo:   "2018-04-02",
			Summ:     124,
			CostItem: "seo",
			Order:    nil,
		},
		{
			Source:   nil,
			DateFrom: "2018-04-03",
			DateTo:   "2018-04-03",
			Summ:     125,
			CostItem: "seo",
			Order:    nil,
			Sites:    []string{"catalog-test"},
		},
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	ids = data.UploadedCosts
}

func TestClient_CostsDelete(t *testing.T) {
	c := client()
	data, status, err := c.CostsDelete(ids)

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomFields(t *testing.T) {
	c := client()

	data, status, err := c.CustomFields(CustomFieldsRequest{})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}
func TestClient_CustomDictionaries(t *testing.T) {
	c := client()

	data, status, err := c.CustomDictionaries(CustomDictionariesRequest{
		Filter: CustomDictionariesFilter{
			Name: "test",
		},
		Limit: 10,
		Page:  1,
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomDictionary(t *testing.T) {
	c := client()

	data, status, err := c.CustomDictionary("test2")

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomDictionariesCreate(t *testing.T) {
	c := client()

	data, status, err := c.CustomDictionariesCreate(CustomDictionary{
		Name: "test2",
		Code: "test2",
		Elements: []Element{
			{
				Name: "test",
				Code: "test",
			},
		},
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomDictionaryEdit(t *testing.T) {
	c := client()

	data, status, err := c.CustomDictionaryEdit(CustomDictionary{
		Name: "test223",
		Code: "test2",
		Elements: []Element{
			{
				Name: "test3",
				Code: "test3",
			},
		},
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomFieldsCreate(t *testing.T) {
	c := client()

	data, status, err := c.CustomFieldsCreate(CustomFieldsEditRequest{
		CustomField: CustomFields{
			Name:        "test4",
			Code:        "test4",
			Type:        "text",
			Entity:      "order",
			DisplayArea: "customer",
		},
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomField(t *testing.T) {
	c := client()

	data, status, err := c.CustomField("customer", "testtest")

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}

func TestClient_CustomFieldEdit(t *testing.T) {
	c := client()

	data, status, err := c.CustomFieldEdit("customer",  CustomFieldsEditRequest{
		CustomField: CustomFields{
			Name: "testtesttest",
		},
	})

	if err.ErrorMsg != "" {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err.ErrorMsg)
		t.Fail()
	}
}
