package amember

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	_ "github.com/go-sql-driver/mysql"
	"github.com/paperclicks/golog"
)

type Amember struct {
	APIKey   string
	APIURL   string
	Gologger *golog.Golog
	client   *http.Client
	DB       *sql.DB
}

type Params struct {
	Filter map[string]string
	Nested []string
	Count  int
	Page   int
}

type Condition struct {
	Column   string
	Operator string
	Values   []interface{}
	SubQuery string
}

func New(apiURL string, apiKey string, gl *golog.Golog) *Amember {

	//create a custom timout dialer
	dialer := &net.Dialer{Timeout: 30 * time.Second}

	//create a custom transport layer to use during API calls
	tr := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cli := &http.Client{Transport: tr}

	return &Amember{APIURL: apiURL, APIKey: apiKey, Gologger: gl, client: cli}
}

func NewWithDb(apiURL string, apiKey string, dburi string, gl *golog.Golog) (*Amember, error) {

	//create a custom timout dialer
	dialer := &net.Dialer{Timeout: 30 * time.Second}

	//create a custom transport layer to use during API calls
	tr := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cli := &http.Client{Transport: tr}

	db, err := sql.Open("mysql", dburi)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		gl.Log(err.Error(), golog.ERROR)
		return nil, err
	}

	return &Amember{APIURL: apiURL, APIKey: apiKey, DB: db, Gologger: gl, client: cli}, nil
}

// Memberships return a map of Membership having username as key.
// If activeAccessOnly=true only accesses that have not expired yet will be attached to memberships
func (am *Amember) MembershipsFromDB() (map[string]Membership, error) {

	start := time.Now()
	memberships := make(map[string]Membership)

	users := make(map[int]User)
	accesses := make(map[int][]Access)

	//get users from amember DB
	usersQuery := `select user_id, login as username from am_user where status in
	(1,2) and user_id in (select user_id from am_access where expire_date>= DATE_SUB(NOW(), INTERVAL 30 DAY))`

	rows, err := am.DB.Query(usersQuery)
	if err != nil {
		return memberships, err
	}

	for rows.Next() {

		user := User{}
		rows.Scan(&user.UserID, &user.Login)

		users[user.UserID] = user
	}

	//get access from amember DB
	accessQuery := `select access_id,
	coalesce(invoice_id,0) as invoice_id,
	user_id,
	product_id,
	begin_date,
	expire_date
from am_access where expire_date>= DATE_SUB(NOW(), INTERVAL 30 DAY)`

	rows, err = am.DB.Query(accessQuery)
	if err != nil {
		return memberships, err
	}

	for rows.Next() {

		access := Access{}
		err := rows.Scan(&access.AccessID, &access.InvoiceID, &access.UserID, &access.ProductID, &access.BeginDate, &access.ExpireDate)
		if err != nil {
			return memberships, err
		}
		accesses[access.UserID] = append(accesses[access.UserID], access)

	}

	//build memberships from users and access records
	for userID, user := range users {

		membership := Membership{}
		membership.User = user
		membership.Accesses = accesses[userID]

		memberships[user.Login] = membership
	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] memberships in [%f] seconds", len(memberships), time.Since(start).Seconds()), golog.DEBUG)

	return memberships, nil
}

// UsersFromDB return a map of Users having the userID as key.
func (am *Amember) UsersFromDB(status int, addedFrom time.Time, addedTo time.Time) (map[int]DBUser, error) {

	start := time.Now()

	users := make(map[int]DBUser)

	//get users from amember DB
	usersQuery := `select user_id, login, name_f, name_l,email,status,added from am_user where status=? and added between ? and ?`

	rows, err := am.DB.Query(usersQuery, status, addedFrom, addedTo)
	if err != nil {
		return users, err
	}

	for rows.Next() {

		user := DBUser{}
		err := rows.Scan(&user.UserID, &user.Login, &user.NameF, &user.NameL, &user.Email, &user.Status, &user.Added)
		if err != nil {
			return users, err
		}
		users[user.UserID] = user
	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] users in [%f] seconds", len(users), time.Since(start).Seconds()), golog.DEBUG)

	return users, nil
}

// UsersFromDB return a map of Users having the userID as key.
func (am *Amember) GetUserFromView(email string) (ViewUser, error) {

	user := ViewUser{}

	//get users from amember DB
	usersQuery := `select userId,
	username,
	first_name,
	last_name,
	email,
	signup_date,
	subscriptionStatus,
	click_id,
	mobile_phone,
	subscription_plan,
	product_name,
	expiration_date,
	total_months,
	total_days,
	total_days_excluding_trial,
	total_payments,
	first_payment,
	last_payment,
	how_did_you_hear
from users where email=?`

	rows, err := am.DB.Query(usersQuery, email)
	if err != nil {
		return user, err
	}

	for rows.Next() {

		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName,
			&user.Email, &user.SignupDate, &user.SubscriptionStatus, &user.ClickID, &user.MobilePhone, &user.SubscriptionPlan,
			&user.ProductName, &user.ExpirationDate, &user.TotalMonths, &user.TotalDays, &user.TotalDaysExcludingTrial,
			&user.TotalPayments, &user.FirstPayment, &user.LastPayment, &user.HowDidYouHear)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

// UsersFromDB return a map of Users having the userID as key.
func (am *Amember) AccessesFromDB(expiredFrom time.Time, expiredTo time.Time) (map[int][]DBAccess, error) {

	start := time.Now()

	accesses := make(map[int][]DBAccess)

	//get users from amember DB
	query := `select coalesce(access_id,0) as access_id,
	coalesce(invoice_id,0) as invoice_id,
	coalesce(invoice_public_id,0) as invoice_public_id,
	coalesce(invoice_payment_id,0) as invoice_payment_id,
	coalesce(invoice_item_id,0) as invoice_item_id,
	coalesce(user_id,0) as user_id,
	coalesce(product_id,0) as product_id,
	coalesce(transaction_id,0) as transaction_id,
	begin_date,
	expire_date,
	coalesce(qty,0) as qty,
	coalesce(comment,"") as comment
from am_access where expire_date>= DATE_SUB(NOW(), INTERVAL 30 DAY)?`
	rows, err := am.DB.Query(query)
	if err != nil {
		return accesses, err
	}

	for rows.Next() {

		access := DBAccess{}
		err := rows.Scan(&access.AccessID, &access.InvoiceID, &access.InvoicePublicID, &access.InvoicePaymentID, &access.InvoiceItemID, &access.UserID,
			&access.ProductID, &access.TransactionID, &access.BeginDate, &access.ExpireDate, &access.Qty, &access.Comment)
		if err != nil {
			return accesses, err
		}

		accesses[access.UserID] = append(accesses[access.UserID], access)
	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] accesses in [%f] seconds", len(accesses), time.Since(start).Seconds()), golog.DEBUG)

	return accesses, nil
}

// Users returns a map of User having username as key
func (am *Amember) Users(p Params) map[string]User {

	start := time.Now()

	users := make(map[string]User)

	page := p.Page
	count := p.Count
	if p.Count == 0 {
		count = 100
	}

	//reange over all the pages
	for {
		//update page value and parse params
		p.Page = page
		params := am.parseParams(p)

		//add page param the url
		url := fmt.Sprintf("%s/api/users?_key=%s%s", am.APIURL, am.APIKey, params)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Log(err, golog.ERROR)
			return users
		}

		//fmt.Printf("%#v", response)

		//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
		for k, v := range response {

			if k == "_total" {
				continue
			}

			m := v.(map[string]interface{})

			u := User{}
			//try to parse the map into the struct fields
			am.mapToStruct(m, &u)

			users[u.Login] = u
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] users in [%f] seconds", len(users), time.Since(start).Seconds()), golog.DEBUG)

	return users
}

func (am *Amember) Invoices(p Params) map[int]Invoice {

	start := time.Now()

	invoices := make(map[int]Invoice)

	page := p.Page
	count := p.Count
	if p.Count == 0 {
		count = 100
	}

	//reange over all the pages
	for {

		//update page value and parse params
		p.Page = page
		params := am.parseParams(p)

		//add page param the url
		url := fmt.Sprintf("%s/api/invoices?_key=%s%s", am.APIURL, am.APIKey, params)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Log(err, golog.ERROR)
			return invoices
		}

		for k, v := range response {

			if k == "_total" {
				continue
			}

			rawInvoice := v.(map[string]interface{})

			invoice := Invoice{}
			nested := InvoiceNested{}

			invoicePayments := []Payment{}
			invoiceItems := []Item{}
			invoiceAccess := []Access{}

			//try to parse the map into the struct fields
			am.mapToStruct(rawInvoice, &invoice)

			rawNested := rawInvoice["nested"].(map[string]interface{})

			for k2, v2 := range rawNested {

				switch k2 {
				case "invoice-payments":

					rawPayments := v2.([]interface{})

					for _, v3 := range rawPayments {
						var payment Payment
						m := v3.(map[string]interface{})

						am.mapToStruct(m, &payment)
						invoicePayments = append(invoicePayments, payment)
					}
				case "access":

					rawAccess := v2.([]interface{})

					for _, v3 := range rawAccess {
						var access Access
						m := v3.(map[string]interface{})

						am.mapToStruct(m, &access)
						invoiceAccess = append(invoiceAccess, access)
					}

				case "invoice-items":

					rawItems := v2.([]interface{})

					for _, v3 := range rawItems {
						var item Item

						m := v3.(map[string]interface{})

						am.mapToStruct(m, &item)
						invoiceItems = append(invoiceItems, item)
					}

				}
			}

			nested.InvoicePayments = invoicePayments
			nested.Access = invoiceAccess
			nested.InvoiceItems = invoiceItems

			invoice.Nested = nested
			invoices[invoice.InvoiceID] = invoice
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] invoices in [%f] seconds", len(invoices), time.Since(start).Seconds()), golog.DEBUG)

	return invoices
}

// Accesses returns a map of Access slices. The map has user_id as key
func (am *Amember) Accesses(p Params, activeOnly bool) map[int][]Access {

	start := time.Now()

	accesses := make(map[int][]Access)

	page := p.Page
	count := p.Count
	if p.Count == 0 {
		count = 100
	}

	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	//reange over all the pages
	for {

		//update page value and parse params
		p.Page = page
		params := am.parseParams(p)

		//add page param the url
		url := fmt.Sprintf("%s/api/access?_key=%s%s", am.APIURL, am.APIKey, params)

		am.Gologger.Log(fmt.Sprintf("GET: %s", url), golog.DEBUG)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Log(err, golog.ERROR)
			return accesses
		}

		//fmt.Printf("%#v", response)

		//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
		for k, v := range response {

			if k == "_total" {
				continue
			}

			m := v.(map[string]interface{})

			i := Access{}
			//try to parse the map into the struct fields
			am.mapToStruct(m, &i)

			expires := time.Date(i.ExpireDate.Year(), i.ExpireDate.Month(), i.ExpireDate.Day(), 0, 0, 0, 0, i.ExpireDate.Location())

			//skip any expired acces if we are requesting only active ones
			if expires.Before(today) && activeOnly {
				continue
			}

			accesses[i.UserID] = append(accesses[i.UserID], i)
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] accesses in [%f] seconds", len(accesses), time.Since(start).Seconds()), golog.DEBUG)

	return accesses
}

func (am *Amember) Payments(p Params) map[int]Payment {
	start := time.Now()

	payments := make(map[int]Payment)

	page := p.Page
	count := p.Count
	if p.Count == 0 {
		count = 100
	}

	//reange over all the pages
	for {
		//update page value and parse params
		p.Page = page
		params := am.parseParams(p)

		//add page param the url
		url := fmt.Sprintf("%s/api/invoice-payments?_key=%s%s", am.APIURL, am.APIKey, params)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Log(err, golog.ERROR)
			return payments
		}

		//fmt.Printf("%#v", response)

		//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
		for k, v := range response {

			if k == "_total" {
				continue
			}

			m := v.(map[string]interface{})

			i := Payment{}
			//try to parse the map into the struct fields
			am.mapToStruct(m, &i)

			payments[i.InvoicePaymentID] = i
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] payments in [%f] seconds", len(payments), time.Since(start).Seconds()), golog.DEBUG)

	return payments
}

// Memberships return a map of Membership having username as key.
// If activeAccessOnly=true only accesses that have not expired yet will be attached to memberships
func (am *Amember) Memberships(p Params, activeAccessOnly bool) map[string]Membership {

	start := time.Now()
	memberships := make(map[string]Membership)

	page := p.Page
	count := p.Count
	if p.Count == 0 {
		count = 100
	}

	//reange over all the pages
	for {
		//update page value and parse params
		p.Page = page
		params := am.parseParams(p)

		//add page param the url
		url := fmt.Sprintf("%s/api/users?_key=%s%s", am.APIURL, am.APIKey, params)
		fmt.Println(url)
		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Log(err.Error(), golog.ERROR)
			fmt.Println(err)
			return memberships
		}

		//fmt.Printf("%#v", response)

		//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
		for k, v := range response {

			membership := Membership{}

			if k == "_total" {
				continue
			}

			uMap := v.(map[string]interface{})

			u := User{}
			//parse user data and add to the current membership
			am.mapToStruct(uMap, &u)
			membership.User = u

			//if this user got no nested element, and activeAccessOnly=true skip this user
			if uMap["nested"] == nil && activeAccessOnly {
				continue
			}

			//if this user got no nested element, but activeAccessOnly=false, add user and skip the rest
			if uMap["nested"] == nil && !activeAccessOnly {
				memberships[membership.User.Login] = membership
				continue
			}

			nMap := uMap["nested"].(map[string]interface{})

			aMap := nMap["access"].([]interface{})

			//parse all access data for current user and add to the current membership
			accesses := []Access{}
			for _, a := range aMap {
				m := a.(map[string]interface{})
				access := Access{}

				am.mapToStruct(m, &access)

				//add only if access is valid (not expired)
				if activeAccessOnly && validAccess(access) {
					accesses = append(accesses, access)
					continue
				}

				//otherwise add all accesses
				accesses = append(accesses, access)

			}

			membership.Accesses = accesses

			memberships[membership.User.Login] = membership

		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] memberships in [%f] seconds", len(memberships), time.Since(start).Seconds()), golog.DEBUG)

	return memberships
}

// ProductCategories returns a map of products having product id as key, and the corresponding map of categories as value
func (am *Amember) ProductCategories() map[int]map[int]int {

	start := time.Now()

	pc := make(map[int]map[int]int)

	//add page param the url
	url := fmt.Sprintf("%s/api/product-product-category?_key=%s", am.APIURL, am.APIKey)

	response, err := am.doGet(url)
	if err != nil {
		am.Gologger.Log(err, golog.ERROR)
		return pc
	}

	//fmt.Printf("%#v", response)

	//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
	for k, v := range response {

		if k == "_total" {
			continue
		}

		prod := v.([]interface{})

		cid, err := strconv.Atoi(k)
		if err != nil {
			panic(err)
		}

		//range over the slice of product ids and build the final response
		for _, pi := range prod {

			id, err := strconv.Atoi(pi.(string))
			if err != nil {

				panic(err)
			}

			//if categories map is nil, first initialize the map
			if pc[id] == nil {

				pc[id] = make(map[int]int)
			}

			pc[id][cid] = cid

		}

	}
	am.Gologger.Log(fmt.Sprintf("Returned [%d] products with categories in [%f] seconds", len(pc), time.Since(start).Seconds()), golog.DEBUG)

	return pc
}

func (am *Amember) Products(p Params) map[int]Product {
	start := time.Now()

	products := make(map[int]Product)

	page := p.Page
	count := p.Count
	if p.Count == 0 {
		count = 100
	}

	//reange over all the pages
	for {

		//update page value and parse params
		p.Page = page
		params := am.parseParams(p)

		//add page param the url
		url := fmt.Sprintf("%s/api/products?_key=%s%s", am.APIURL, am.APIKey, params)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Log(err, golog.ERROR)
			return products
		}

		//fmt.Printf("%#v", response)

		//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
		for k, v := range response {

			if k == "_total" {
				continue
			}

			m := v.(map[string]interface{})

			i := Product{}
			//try to parse the map into the struct fields
			am.mapToStruct(m, &i)

			products[i.ProductID] = i
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Log(fmt.Sprintf("Returned [%d] products in [%f] seconds", len(products), time.Since(start).Seconds()), golog.DEBUG)

	return products
}

func (am *Amember) doGet(url string) (map[string]interface{}, error) {

	response := make(map[string]interface{})

	am.Gologger.Log(fmt.Sprintf("GET: %s", url), golog.DEBUG)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return response, err
	}

	resp, err := am.client.Do(req)

	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return response, err
	}

	var responseError bool
	var responseMessage string

	if _, ok := response["error"]; ok {
		responseError = response["error"].(bool)
	}

	if _, ok := response["message"]; ok {
		responseMessage = response["message"].(string)

	}

	if responseError == true {
		return response, errors.New(responseMessage)
	}

	return response, nil
}

func (am *Amember) mapToStruct(m map[string]interface{}, s interface{}) {

	//uValue := reflect.ValueOf(u)
	elem := reflect.ValueOf(s).Elem()

	//range over all the fields of the struct
	for i := 0; i < elem.NumField(); i++ {

		f := elem.Field(i)
		jsonTag := elem.Type().Field(i).Tag.Get("json")

		val, ok := m[jsonTag]
		if !ok {
			am.Gologger.Log(fmt.Sprintf("Key [%s] not found in map %#v", jsonTag, m), golog.DEBUG)
		}

		//check inf the underlaying value of interface{} is nil.
		//In this case the value will be ignored and the struct will have the zero value for that field
		if val == nil {
			continue
		}

		//switch on the type of the i-th struct field, and perform the right conversion
		switch t := f.Interface().(type) {

		case string:
			f.SetString(val.(string))

		case int:

			kind := reflect.ValueOf(val).Kind()
			//switch on the type of underlying value of the interface, and attempt a conversion to int
			switch kind {
			case reflect.String:
				//if the string is empty leave the 0 value
				if val.(string) == "" {
					break
				}

				v, err := strconv.Atoi(val.(string))
				if err != nil {
					am.Gologger.Log(fmt.Sprintf("strconv error for field [%s]: %v", jsonTag, err), golog.ERROR)
					break
				}
				f.SetInt(int64(v))

			case reflect.Int:

				if _, ok := val.(int); !ok {
					am.Gologger.Log(fmt.Sprintf("assertion to int failed for field [%s]: %v", jsonTag, val), golog.ERROR)
					break
				}
				f.SetInt(int64(val.(int)))

			case reflect.Int32:
				if _, ok := val.(int32); !ok {
					am.Gologger.Log(fmt.Sprintf("assertion to int32 failed: %v", val), golog.ERROR)
					break
				}
				f.SetInt(int64(val.(int32)))

			case reflect.Int64:
				if _, ok := val.(int64); !ok {
					am.Gologger.Log(fmt.Sprintf("assertion to int64 failed: %v", val), golog.ERROR)
					break
				}
				f.SetInt(int64(val.(int64)))

			case reflect.Float64:
				if _, ok := val.(float64); !ok {
					am.Gologger.Log(fmt.Sprintf("assertion to float64 failed: %v", val), golog.ERROR)
					break
				}

				f.SetInt(int64(val.(float64)))

			default:
				am.Gologger.Log(fmt.Sprintf("%v", reflect.ValueOf(val).Kind()), golog.DEBUG)
			}

		case float32:
			v, err := strconv.ParseFloat(val.(string), 32)
			if err != nil {
				break
				//panic(err)
			}
			f.SetFloat(v)
		case CustomTime:
			ct := CustomTime{}
			err := ct.UnmarshalJSON([]byte(val.(string)))
			if err != nil {
				am.Gologger.Log(fmt.Sprintf("Error parsing date string %s - %v", val.(string), err), golog.ERROR)
				break
			}
			f.Set(reflect.ValueOf(ct))
		case time.Time:

			t, err := dateparse.ParseAny(val.(string))
			if err != nil {
				am.Gologger.Log(fmt.Sprintf("Error parsing date string %s - %v", val.(string), err), golog.ERROR)
				break
			}

			f.Set(reflect.ValueOf(t))
		default:
			am.Gologger.Log(fmt.Sprintf("Type for field [%s] not found: %v", jsonTag, t), golog.DEBUG)

		}
	}
}

func (am *Amember) parseParams(p Params) string {

	qs := ""

	//add all eventual filters
	for k, v := range p.Filter {

		qs = fmt.Sprintf("%s&_filter[%s]=%s", qs, k, v)
	}

	//add all eventual nested
	for _, v := range p.Nested {

		qs = fmt.Sprintf("%s&_nested[]=%s", qs, v)
	}

	//add "page" param; if not set it starts from page=0
	qs = fmt.Sprintf("%s&_page=%d", qs, p.Page)

	//add "count" param; if not set the default value is 100
	if p.Count > 0 {
		qs = fmt.Sprintf("%s&_count=%d", qs, p.Count)
	} else {
		qs = fmt.Sprintf("%s&_count=%d", qs, 100)
	}

	return qs
}

func validAccess(a Access) bool {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	expires := time.Date(a.ExpireDate.Year(), a.ExpireDate.Month(), a.ExpireDate.Day(), 0, 0, 0, 0, a.ExpireDate.Location())

	//skip any expired acces if we are requesting only active ones
	if expires.Before(today) {
		return false
	}

	return true
}

func (am *Amember) ExpiredUsers(expiredSince int) map[string]User {

	expiredUsers := make(map[string]User)
	//1 get users
	users := am.Users(Params{})
	//2 get accesses
	accesses := am.Accesses(Params{}, false)

	// for each user
	//// for each access
	//// if access=active break, if expire is not before the wanted interval continue, else (if for each was not breaked by an active subscription or a subscriptio with an earlier expire date) add this user to the list

	for _, u := range users {

		excludeUser := false
		var expired time.Time

		for _, a := range accesses[u.UserID] {

			//if the expire_date of at least one access is earlier than expiredSince, then break the foreach here. This user must not be added to the list
			if time.Since(a.ExpireDate).Hours() < float64(expiredSince*24) {
				excludeUser = true
				break
			}

			//get the last access
			if a.ExpireDate.After(expired) {
				expired = a.ExpireDate
			}

		}

		if !excludeUser {
			u.ExpiredAt.Time = expired
			expiredUsers[u.Login] = u

			am.Gologger.Log(fmt.Sprintf("Adding user to expired list: [username: %s] [expired: %s]  [days: %f]", u.Login, expired.Format("2006-01-02"), time.Since(expired).Hours()/24), golog.INFO)
		}
	}
	return expiredUsers
}

// PaymentsByDay returns a map having username as key and Payment object as value, for a given date
// For native payments: itemTitle=Native, itemDescription=""
// For mobile payments: itemTitle=Mobile, itemDescription=""
// For overage payments: itemTitle=Overage, itemDescription=[Native,Mobile]
func (am *Amember) PaymentsByDate(datetime time.Time, itemTitle string, itemDescription string) (map[string]Payment, error) {

	paymets := make(map[string]Payment)

	//check connection and reconnet if necessary
	err := am.DB.Ping()
	if err != nil {
		return paymets, err
	}

	//because overages do not have an item description anymore, but all the description is inside the item title, we need to check for the keyword Native
	q := `select ip.user_id, u.login as username,  ip.dattm, ip.amount
       	from am_invoice_payment ip
		left join am_invoice_item ii on ii.invoice_id=ip.invoice_id
		left join am_user u on ip.user_id=u.user_id
		where ip.amount>0.0 and (refund_amount is  null or refund_amount=0.0)
		and ip.dattm between ? and ?
		and ii.item_title  like ? and (ii.item_title  like ? OR ii.item_description like ?)`

	from := fmt.Sprintf("%s 00:00:01", datetime.Format("2006-01-02"))
	to := fmt.Sprintf("%s 23:59:59", datetime.Format("2006-01-02"))

	rows, err := am.DB.Query(q, from, to, "%"+itemTitle+"%", "%"+itemTitle+"%", "%"+itemDescription+"%")
	defer rows.Close()

	if err != nil {
		return paymets, err
	}

	for rows.Next() {
		var (
			userID      int
			username    string
			paymentdate sql.NullTime
			amount      float32
		)

		err := rows.Scan(&userID, &username, &paymentdate, &amount)
		if err != nil {
			return paymets, err
		}

		payment := Payment{}
		payment.Username = username
		payment.Amount = amount
		payment.Dattm = paymentdate.Time

		paymets[username] = payment

	}
	return paymets, nil
}

// RefundsByDate returns a map having username as key and Payment object as value, for refunds in a given date
func (am *Amember) RefundsByDate(datetime time.Time, itemTitle string) (map[string]Payment, error) {

	paymets := make(map[string]Payment)

	//check connection and reconnet if necessary
	err := am.DB.Ping()
	if err != nil {
		return paymets, err
	}

	q := `select ip.user_id, u.login as username, ip.refund_dattm, ip.refund_amount
		 from am_invoice_payment ip
		left join am_invoice_item ii on ii.invoice_id=ip.invoice_id
		left join am_user u on ip.user_id=u.user_id
		where ip.refund_amount>0.0
		and ip.refund_dattm between ? and ? and (ii.item_title like ? or ii.item_description like ?)`

	//it considers a refunds to belong to a certain platform if itemTitle is found in either the description or the title (for example overage refunds)
	from := fmt.Sprintf("%s 00:00:01", datetime.Format("2006-01-02"))
	to := fmt.Sprintf("%s 23:59:59", datetime.Format("2006-01-02"))

	rows, err := am.DB.Query(q, from, to, "%"+itemTitle+"%", "%"+itemTitle+"%")
	defer rows.Close()

	if err != nil {
		return paymets, err
	}

	for rows.Next() {
		var (
			userID     int
			username   string
			refundDate sql.NullTime
			amount     float32
		)

		err := rows.Scan(&userID, &username, &refundDate, &amount)
		if err != nil {
			return paymets, err
		}

		payment := Payment{}
		payment.Username = username
		payment.RefundAmount = amount
		payment.RefundDattm = refundDate.Time

		paymets[username] = payment

	}
	return paymets, nil
}

// UsersFromDB return a map of Users having the userID as key.
func (am *Amember) GetUsersFromView(conditions []Condition, limit int) ([]ViewUser, error) {

	users := []ViewUser{}

	whereConditions, conditionValues := BuildWhereConditions(conditions, 10000)

	//get users from amember DB
	usersQuery := fmt.Sprintf(`select userId,
	username,
	first_name,
	last_name,
	email,
	signup_date,
	subscriptionStatus,
	click_id,
	mobile_phone,
	subscription_plan,
	product_name,
	expiration_date,
	total_months,
	total_days,
	total_days_excluding_trial,
	total_payments,
	first_payment,
	last_payment,
	how_did_you_hear,
	coalesce(preferred_contact_method,'not_specified') as preferred_contact_method ,
	coalesce(preferred_contact,'not_specified') as preferred_contact,
	payments_last_3_months,
	is_top_paying_user,
	cancellation_date,
	last_updated
from users %s`, whereConditions)

	//test
	rows, err := am.DB.Query(usersQuery, conditionValues...)
	if err != nil {
		panic(err)
	}

	for rows.Next() {

		user := ViewUser{}
		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName,
			&user.Email, &user.SignupDate, &user.SubscriptionStatus, &user.ClickID, &user.MobilePhone, &user.SubscriptionPlan,
			&user.ProductName, &user.ExpirationDate, &user.TotalMonths, &user.TotalDays, &user.TotalDaysExcludingTrial,
			&user.TotalPayments, &user.FirstPayment, &user.LastPayment, &user.HowDidYouHear, &user.PreferredContactMethod, &user.PreferredContact,
			&user.PaymentsLast3Months, &user.IsTopPayingUser, &user.CancellationDate, &user.LastUpdated)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users, nil
}

// SelectQuery gets in input a Model and a slice of Condition, and returns a parametrized query and a slice of values to pass to the Exec method
func BuildWhereConditions(conditions []Condition, limit int) (string, []interface{}) {

	//whereConditions is a slice with already formated pieces of conditions
	//these conditions will be finally joined with "AND" operator
	whereConditions := []string{}

	//keeps track of the total number of params in the final query
	paramsIndex := 1

	values := []interface{}{}

	limitCondition := ""

	if limit > 0 {
		limitCondition = fmt.Sprintf("LIMIT %d", limit)
	}

	//if conditions are empty return a query without where conditions, and empty values
	if len(conditions) == 0 {
		q := fmt.Sprintf("WHERE ? %s", limitCondition)

		values = append(values, true)

		return q, values
	}

	//for each element in the conditions map
	//check if the lenght of the values for the condition
	//if there is only 1 value then it is a simple operator (=,=>, <=, !=)
	//if the length is 2 or higher it might be one of IN, BETWEEN operators
	for _, v := range conditions {

		switch v.Operator {
		case "IN":

			//if a subquery was passed for the IN condition the simply append "column IN (subquery)"
			if v.SubQuery != "" {

				whereConditions = append(whereConditions, fmt.Sprintf("%s IN (%s)", v.Column, v.SubQuery))
				break
			}

			params := []string{}
			//add a list of params for the IN condition in the form of $1, $2, $3...
			for _, val := range v.Values {

				params = append(params, "?")
				values = append(values, val)
				paramsIndex++
			}
			//append "column IN ($1, $2, $3,...)"
			whereConditions = append(whereConditions, fmt.Sprintf("%s %s (%s)", v.Column, v.Operator, strings.Join(params, ",")))

		case "BETWEEN":

			//append "column BETWEEN $1 AND $2"
			whereConditions = append(whereConditions, fmt.Sprintf("%s %s $? AND ?", v.Column, v.Operator))

			values = append(values, v.Values[0])
			values = append(values, v.Values[1])

			//increase index by 2 because 2 values were used here
			paramsIndex = paramsIndex + 2

		case "=", ">", "<", "<=", ">=", "<>", "!=":

			whereConditions = append(whereConditions, fmt.Sprintf("%s %s ?", v.Column, v.Operator))

			values = append(values, v.Values[0])

			paramsIndex++
		case "LIKE":

			whereConditions = append(whereConditions, fmt.Sprintf("%s %s ?", v.Column, v.Operator))

			values = append(values, v.Values[0])

			paramsIndex++

		default:
			fmt.Printf("Missing or unknown operator: %s\n", v.Operator)
		}

	}

	q := fmt.Sprintf("WHERE %s %s", strings.Join(whereConditions, " AND "), limitCondition)

	return q, values
}
