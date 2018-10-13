package amember

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/paperclicks/golog"
)

type Amember struct {
	APIKey   string
	APIURL   string
	Gologger *golog.Golog
	client   *http.Client
}

type Params struct {
	Filter map[string]string
	Nested []string
}

func New(apiURL string, apiKey string, output io.Writer) *Amember {
	gologger := golog.New(output)
	gologger.ShowCallerInfo = true

	//create a custom timout dialer
	dialer := &net.Dialer{Timeout: 30 * time.Second}

	//create a custom transport layer to use during API calls
	tr := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cli := &http.Client{Transport: tr}

	return &Amember{APIURL: apiURL, APIKey: apiKey, Gologger: gologger, client: cli}
}

//Users returns a map of User having username as key
func (am *Amember) Users(p Params) map[string]User {

	start := time.Now()

	users := make(map[string]User)

	page := 0
	count := 100
	params := am.parseParams(p)

	//reange over all the pages
	for {

		//add page param the url
		url := fmt.Sprintf("%s/api/users?_key=%s%s&_count=%d&_page=%d", am.APIURL, am.APIKey, params, count, page)

		am.Gologger.Debug("URL: %s", url)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Error("doGet error: %v", err)
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

	am.Gologger.Info("Returned [%d] users in [%f] seconds", len(users), time.Since(start).Seconds())

	return users
}

func (am *Amember) Invoices(p Params) map[int]Invoice {

	start := time.Now()

	invoices := make(map[int]Invoice)

	page := 0
	count := 100
	params := am.parseParams(p)
	//reange over all the pages
	for {

		//add page param the url
		url := fmt.Sprintf("%s/api/invoices?_key=%s%s&_page=%d&_count=%d", am.APIURL, am.APIKey, params, page, count)

		am.Gologger.Debug("URL: %s", url)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Error("doGet error: %v", err)
			return invoices
		}

		//fmt.Printf("%#v", response)

		//perform again a marshall for every element of the response, and attempt to unmarshall into User struct
		for k, v := range response {

			if k == "_total" {
				continue
			}

			m := v.(map[string]interface{})

			i := Invoice{}
			//try to parse the map into the struct fields
			am.mapToStruct(m, &i)

			invoices[i.InvoiceID] = i
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	fmt.Printf("%#v", invoices)

	am.Gologger.Info("Returned [%d] invoices in [%f] seconds", len(invoices), time.Since(start).Seconds())

	return invoices
}

//Accesses returns a map of Access slices. The map has user_id as key
func (am *Amember) Accesses(p Params, activeOnly bool) map[int][]Access {

	start := time.Now()

	accesses := make(map[int][]Access)

	page := 0
	count := 100
	params := am.parseParams(p)

	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	//reange over all the pages
	for {

		//add page param the url
		url := fmt.Sprintf("%s/api/access?_key=%s%s&_page=%d&_count=%d", am.APIURL, am.APIKey, params, page, count)

		am.Gologger.Debug("URL: %s", url)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Error("doGet error: %v", err)
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

			expires := time.Date(i.ExpireDate.Time.Year(), i.ExpireDate.Time.Month(), i.ExpireDate.Time.Day(), 0, 0, 0, 0, i.ExpireDate.Time.Location())

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

	am.Gologger.Info("Returned [%d] accesses in [%f] seconds", len(accesses), time.Since(start).Seconds())

	return accesses
}

func (am *Amember) Payments(p Params) map[int]Payment {
	start := time.Now()

	payments := make(map[int]Payment)

	page := 0
	count := 100
	params := am.parseParams(p)

	//reange over all the pages
	for {

		//add page param the url
		url := fmt.Sprintf("%s/api/invoice-payments?_key=%s%s&_page=%d&_count=%d", am.APIURL, am.APIKey, params, page, count)

		am.Gologger.Debug("URL: %s", url)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Error("doGet error: %v", err)
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

	fmt.Printf("%#v", payments)

	am.Gologger.Info("Returned [%d] payments in [%f] seconds", len(payments), time.Since(start).Seconds())

	return payments
}

func (am *Amember) doGet(url string) (map[string]interface{}, error) {

	response := make(map[string]interface{})

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
			am.Gologger.Error("Key [%s] not found in map %#v", jsonTag, m)
		}

		if jsonTag == "added" {

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

				v, err := strconv.Atoi(val.(string))
				if err != nil {
					am.Gologger.Debug("strconv error for field [%s]: %v", jsonTag, err)
					break
				}
				f.SetInt(int64(v))

			case reflect.Int:

				if _, ok := val.(int); !ok {
					am.Gologger.Debug("assertion to int failed for field [%s]: %v", jsonTag, val)
					break
				}
				f.SetInt(int64(val.(int)))

			case reflect.Int32:
				if _, ok := val.(int32); !ok {
					am.Gologger.Debug("assertion to int32 failed: %v", val)
					break
				}
				f.SetInt(int64(val.(int32)))

			case reflect.Int64:
				if _, ok := val.(int64); !ok {
					am.Gologger.Debug("assertion to int64 failed: %v", val)
					break
				}
				f.SetInt(int64(val.(int64)))

			case reflect.Float64:
				if _, ok := val.(float64); !ok {
					am.Gologger.Debug("assertion to float64 failed: %v", val)
					break
				}

				f.SetInt(int64(val.(float64)))

			default:
				am.Gologger.Debug("%v", reflect.ValueOf(val).Kind())
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

				panic(err)
			}
			f.Set(reflect.ValueOf(ct))
		default:
			am.Gologger.Error("Type for field [%s] not found: %v", jsonTag, t)
			//panic(fmt.Sprintf("Type not found: %v", t))
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

	return qs
}

func (am *Amember) Memberships(p Params) map[string]Membership {

	start := time.Now()
	memberships := make(map[string]Membership)

	page := 0
	count := 100
	params := am.parseParams(p)

	//reange over all the pages
	for {

		//add page param the url
		url := fmt.Sprintf("%s/api/users?_key=%s%s&_count=%d&_page=%d", am.APIURL, am.APIKey, params, count, page)

		am.Gologger.Debug("URL: %s", url)

		response, err := am.doGet(url)
		if err != nil {
			am.Gologger.Error("doGet error: %v", err)
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

			nMap := uMap["nested"].(map[string]interface{})
			aMap := nMap["access"].([]interface{})

			//parse all access data for current user and add to the current membership
			accesses := []Access{}
			for _, a := range aMap {
				m := a.(map[string]interface{})
				access := Access{}

				am.mapToStruct(m, &access)

				//add only if access is valid (not expired)
				if validAccess(access) {
					accesses = append(accesses, access)
				}

			}

			membership.Accesses = accesses

			memberships[membership.User.Login] = membership

		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	am.Gologger.Info("Returned [%d] memberships in [%f] seconds", len(memberships), time.Since(start).Seconds())

	return memberships
}

func (am *Amember) ProductCategories() map[int]map[int]int {

	start := time.Now()

	pc := make(map[int]map[int]int)

	//add page param the url
	url := fmt.Sprintf("%s/api/product-product-category?_key=%s", am.APIURL, am.APIKey)

	am.Gologger.Debug("URL: %s", url)

	response, err := am.doGet(url)
	if err != nil {
		am.Gologger.Error("doGet error: %v", err)
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
			am.Gologger.Debug("strconv error converting category id [%s] to int: %v", k, err)
			break
		}

		//range over the slice of product ids and build the final response
		for _, pi := range prod {

			id, err := strconv.Atoi(pi.(string))
			if err != nil {
				am.Gologger.Debug("strconv error converting product id [%s]: %v", pi.(string), err)
				break
			}

			//if categories map is nil, first initialize the map
			if pc[id] == nil {

				pc[id] = make(map[int]int)
			}

			pc[id][cid] = cid

		}

	}

	fmt.Printf("%#v", pc)

	am.Gologger.Info("Returned [%d] products with categories in [%f] seconds", len(pc), time.Since(start).Seconds())

	return pc
}
func validAccess(a Access) bool {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	expires := time.Date(a.ExpireDate.Time.Year(), a.ExpireDate.Time.Month(), a.ExpireDate.Time.Day(), 0, 0, 0, 0, a.ExpireDate.Time.Location())

	//skip any expired acces if we are requesting only active ones
	if expires.Before(today) {
		return false
	}

	return true
}
