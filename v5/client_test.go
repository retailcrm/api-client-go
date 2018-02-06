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

	data, status, err := c.ApiVersions()
	if err.ErrorMsg != "" {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Fail()
	}

	if data.Success != true {
		t.Fail()
	}
}

func TestClient_ApiCredentialsCredentials(t *testing.T) {
	c := client()

	data, status, err := c.ApiCredentials()
	if err.ErrorMsg != "" {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Fail()
	}

	if data.Success != true {
		t.Fail()
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

	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_CustomerChange(t *testing.T) {
	c := client()

	f := Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.CustomerCreate(f)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	f.Id = cr.Id
	f.Vip = true

	ed, se, err := c.CustomerEdit(f, "id")
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if se != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if ed.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	data, status, err := c.Customer(f.ExternalId, "externalId", "")
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Customer.ExternalId != f.ExternalId {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_CustomersUpload(t *testing.T) {
	c := client()
	customers := make([]Customer, 3)

	for i := range customers {
		customers[i] = Customer{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalId: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	data, status, err := c.CustomersUpload(customers)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_CustomersFixExternalIds(t *testing.T) {
	c := client()
	f := Customer{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.CustomerCreate(f)
	if err.ErrorMsg != "" {
		t.Errorf("%s", sc)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%s", sc)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%s", sc)
		t.Fail()
	}

	customers := []IdentifiersPair{{
		Id:         cr.Id,
		ExternalId: RandomString(8),
	}}

	fx, fe, err := c.CustomersFixExternalIds(customers)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if fx.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_CustomersHistory(t *testing.T) {
	c := client()
	f := CustomersHistoryRequest{
		Filter: CustomersHistoryFilter{
			SinceId: 20,
		},
	}

	data, status, err := c.CustomersHistory(f)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	if len(data.History) == 0 {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_NotesNotes(t *testing.T) {
	c := client()

	data, status, err := c.Notes(NotesRequest{Page: 1})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_NotesCreateDelete(t *testing.T) {
	c := client()

	createCustomerResponse, createCustomerStatus, err := c.CustomerCreate(Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if createCustomerStatus != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if createCustomerResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	noteCreateResponse, noteCreateStatus, err := c.NoteCreate(Note{
		Text:      "some text",
		ManagerId: user,
		Customer: &Customer{
			Id: createCustomerResponse.Id,
		},
	})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if noteCreateStatus != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if noteCreateResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	noteDeleteResponse, noteDeleteStatus, err := c.NoteDelete(noteCreateResponse.Id)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if noteDeleteStatus != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if noteDeleteResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrdersOrders(t *testing.T) {
	c := client()

	data, status, err := c.Orders(OrdersRequest{Filter: OrdersFilter{City: "Москва"}, Page: 1})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrderChange(t *testing.T) {
	c := client()

	random := RandomString(8)

	f := Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: random,
		Email:      fmt.Sprintf("%s@example.com", random),
	}

	cr, sc, err := c.OrderCreate(f)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	f.Id = cr.Id
	f.CustomerComment = "test comment"

	ed, se, err := c.OrderEdit(f, "id")
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if se != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if ed.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	data, status, err := c.Order(f.ExternalId, "externalId", "")
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Order.ExternalId != f.ExternalId {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrdersUpload(t *testing.T) {
	c := client()
	orders := make([]Order, 3)

	for i := range orders {
		orders[i] = Order{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalId: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	data, status, err := c.OrdersUpload(orders)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrdersFixExternalIds(t *testing.T) {
	c := client()
	f := Order{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.OrderCreate(f)
	if err.ErrorMsg != "" {
		t.Errorf("%s", sc)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%s", sc)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%s", sc)
		t.Fail()
	}

	orders := []IdentifiersPair{{
		Id:         cr.Id,
		ExternalId: RandomString(8),
	}}

	fx, fe, err := c.OrdersFixExternalIds(orders)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if fx.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrdersHistory(t *testing.T) {
	c := client()

	data, status, err := c.OrdersHistory(OrdersHistoryRequest{Filter: OrdersHistoryFilter{SinceId: 20}})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	if len(data.History) == 0 {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_PaymentCreateEditDelete(t *testing.T) {
	c := client()

	order := Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	createOrderResponse, status, err := c.OrderCreate(order)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if createOrderResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	f := Payment{
		Order: &Order{
			Id: createOrderResponse.Id,
		},
		Amount: 300,
		Type:   "cash",
	}

	paymentCreateResponse, status, err := c.PaymentCreate(f)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if paymentCreateResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	k := Payment{
		Id:     paymentCreateResponse.Id,
		Amount: 500,
	}

	paymentEditResponse, status, err := c.PaymentEdit(k, "id")
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if paymentEditResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	paymentDeleteResponse, status, err := c.PaymentDelete(paymentCreateResponse.Id)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if paymentDeleteResponse.Success != true {
		t.Errorf("%v", err)
		t.Fail()
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
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_TaskChange(t *testing.T) {
	c := client()

	f := Task{
		Text:        RandomString(15),
		PerformerId: user,
	}

	cr, sc, err := c.TaskCreate(f)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	f.Id = cr.Id
	f.Commentary = RandomString(20)

	gt, sg, err := c.Task(f.Id)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if sg != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if gt.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}

	data, status, err := c.TaskEdit(f)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_UsersUsers(t *testing.T) {
	c := client()

	data, status, err := c.Users(UsersRequest{Filter: UsersFilter{Active: 1}, Page: 1})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_UsersUser(t *testing.T) {
	c := client()

	data, st, err := c.User(user)
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_UsersGroups(t *testing.T) {
	c := client()

	data, status, err := c.UserGroups(UserGroupsRequest{Page: 1})
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_UsersUpdate(t *testing.T) {
	c := client()

	data, st, err := c.UserStatus(user, "busy")
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_StaticticUpdate(t *testing.T) {
	c := client()

	data, st, err := c.StaticticUpdate()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_Countries(t *testing.T) {
	c := client()

	data, st, err := c.Couriers()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_CostGroups(t *testing.T) {
	c := client()

	data, st, err := c.CostGroups()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_CostItems(t *testing.T) {
	c := client()

	data, st, err := c.CostItems()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_Couriers(t *testing.T) {
	c := client()

	data, st, err := c.Couriers()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_DeliveryService(t *testing.T) {
	c := client()

	data, st, err := c.DeliveryService()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_DeliveryTypes(t *testing.T) {
	c := client()

	data, st, err := c.DeliveryTypes()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_LegalEntities(t *testing.T) {
	c := client()

	data, st, err := c.LegalEntities()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrderMethods(t *testing.T) {
	c := client()

	data, st, err := c.OrderMethods()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_OrderTypes(t *testing.T) {
	c := client()

	data, st, err := c.OrderTypes()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_PaymentStatuses(t *testing.T) {
	c := client()

	data, st, err := c.PaymentStatuses()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_PaymentTypes(t *testing.T) {
	c := client()

	data, st, err := c.PaymentTypes()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_PriceTypes(t *testing.T) {
	c := client()

	data, st, err := c.PriceTypes()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_ProductStatuses(t *testing.T) {
	c := client()

	data, st, err := c.ProductStatuses()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_Statuses(t *testing.T) {
	c := client()

	data, st, err := c.Statuses()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_StatusGroups(t *testing.T) {
	c := client()

	data, st, err := c.StatusGroups()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_Sites(t *testing.T) {
	c := client()

	data, st, err := c.Sites()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}

func TestClient_Stores(t *testing.T) {
	c := client()

	data, st, err := c.Stores()
	if err.ErrorMsg != "" {
		t.Errorf("%v", err)
		t.Fail()
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%v", err)
		t.Fail()
	}
}
