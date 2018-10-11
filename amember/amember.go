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
	fmt.Printf("%#v", users)

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

func (am *Amember) Accesses(p Params, activeOnly bool) map[int]Access {

	start := time.Now()

	accesses := make(map[int]Access)

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

			accesses[i.AccessID] = i
		}

		if len(response) < count+1 {
			break
		}

		page++

	}

	fmt.Printf("%#v", accesses)

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
			am.Gologger.Error("Type not found: %v", t)
			panic(fmt.Sprintf("Type not found: %v", t))
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
