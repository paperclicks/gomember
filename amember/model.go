package amember

import (
	"regexp"
	"strings"
	"time"
)

type User struct {
	Added              CustomTime  `json:"added"`
	AffAdded           interface{} `json:"aff_added"`
	AffCustomRedirect  int         `json:"aff_custom_redirect"`
	AffID              interface{} `json:"aff_id"`
	AffPayoutType      interface{} `json:"aff_payout_type"`
	City               string      `json:"city"`
	Comment            string      `json:"comment"`
	Country            string      `json:"country"`
	DisableLockUntil   CustomTime  `json:"disable_lock_until"`
	Email              string      `json:"email"`
	IAgree             int         `json:"i_agree"`
	IsAffiliate        int         `json:"is_affiliate"`
	IsApproved         int         `json:"is_approved"`
	IsLocked           int         `json:"is_locked"`
	Lang               string      `json:"lang"`
	LastIP             string      `json:"last_ip"`
	LastLogin          CustomTime  `json:"last_login"`
	LastSession        string      `json:"last_session"`
	LastUserAgent      string      `json:"last_user_agent"`
	Login              string      `json:"login"`
	NameF              string      `json:"name_f"`
	NameL              string      `json:"name_l"`
	NeedSessionRefresh int         `json:"need_session_refresh"`
	Pass               string      `json:"pass"`
	PassDattm          CustomTime  `json:"pass_dattm"`
	Phone              string      `json:"phone"`
	RememberKey        string      `json:"remember_key"`
	RemoteAddr         string      `json:"remote_addr"`
	ResellerID         int         `json:"reseller_id"`
	SavedFormID        int         `json:"saved_form_id"`
	SignupEmailSent    int         `json:"signup_email_sent"`
	State              string      `json:"state"`
	Status             int         `json:"status"`
	Street             string      `json:"street"`
	Street2            string      `json:"street2"`
	TaxID              interface{} `json:"tax_id"`
	Unsubscribed       int         `json:"unsubscribed"`
	UserAgent          string      `json:"user_agent"`
	UserID             int         `json:"user_id"`
	Zip                string      `json:"zip"`
	//StripeCcExpires    string      `json:"stripe_cc_expires"`
	//StripeCcMasked string      `json:"stripe_cc_masked"`
	//StripeToken    string      `json:"stripe_token"`
}

type Invoice struct {
	InvoiceID         int         `json:"invoice_id"`
	UserID            int         `json:"user_id"`
	PaysysID          string      `json:"paysys_id"`
	Currency          string      `json:"currency"`
	FirstSubtotal     float32     `json:"first_subtotal"`
	FirstDiscount     float32     `json:"first_discount"`
	FirstTax          float32     `json:"first_tax"`
	FirstShipping     float32     `json:"first_shipping"`
	FirstTotal        float32     `json:"first_total"`
	FirstPeriod       string      `json:"first_period"`
	RebillTimes       int         `json:"rebill_times"`
	SecondSubtotal    float32     `json:"second_subtotal"`
	SecondDiscount    float32     `json:"second_discount"`
	SecondTax         float32     `json:"second_tax"`
	SecondShipping    float32     `json:"second_shipping"`
	SecondTotal       float32     `json:"second_total"`
	SecondPeriod      string      `json:"second_period"`
	TaxRate           interface{} `json:"tax_rate"`
	TaxType           interface{} `json:"tax_type"`
	TaxTitle          interface{} `json:"tax_title"`
	Status            int         `json:"status"`
	CouponID          int         `json:"coupon_id"`
	CouponCode        string      `json:"coupon_code"`
	DiscountFirst     float32     `json:"discount_first"`
	DiscountSecond    float32     `json:"discount_second"`
	IsConfirmed       int         `json:"is_confirmed"`
	PublicID          string      `json:"public_id"`
	InvoiceKey        string      `json:"invoice_key"`
	TmAdded           *CustomTime `json:"tm_added"`
	TmStarted         *CustomTime `json:"tm_started"`
	TmCancelled       *CustomTime `json:"tm_cancelled"`
	RebillDate        *CustomTime `json:"rebill_date"`
	DueDate           *CustomTime `json:"due_date"`
	Terms             interface{} `json:"terms"`
	Comment           interface{} `json:"comment"`
	BaseCurrencyMulti float32     `json:"base_currency_multi"`
	SavedFormID       int         `json:"saved_form_id"`
	AffID             interface{} `json:"aff_id"`
	KeywordID         interface{} `json:"keyword_id"`
	RemoteAddr        interface{} `json:"remote_addr"`
}

type Access struct {
	AccessID         int        `json:"access_id"`
	InvoiceID        int        `json:"invoice_id"`
	InvoicePublicID  string     `json:"invoice_public_id"`
	InvoicePaymentID int        `json:"invoice_payment_id"`
	InvoiceItemID    int        `json:"invoice_item_id"`
	UserID           int        `json:"user_id"`
	ProductID        int        `json:"product_id"`
	TransactionID    string     `json:"transaction_id"`
	BeginDate        CustomTime `json:"begin_date"`
	ExpireDate       CustomTime `json:"expire_date"`
	Qty              int        `json:"qty"`
	Comment          string     `json:"comment"`
}

type Payment struct {
	ConversionTrackDone int         `json:"conversion-track-done"`
	GoogleAnalyticsDone int         `json:"google-analytics-done"`
	InvoicePaymentID    int         `json:"invoice_payment_id"`
	InvoiceID           int         `json:"invoice_id"`
	InvoicePublicID     string      `json:"invoice_public_id"`
	UserID              int         `json:"user_id"`
	PaysysID            string      `json:"paysys_id"`
	ReceiptID           string      `json:"receipt_id"`
	TransactionID       string      `json:"transaction_id"`
	Dattm               *CustomTime `json:"dattm"`
	Currency            string      `json:"currency"`
	Amount              float32     `json:"amount"`
	Discount            float32     `json:"discount"`
	Tax                 float32     `json:"tax"`
	Shipping            float32     `json:"shipping"`
	RefundDattm         *CustomTime `json:"refund_dattm"`
	RefundAmount        float32     `json:"refund_amount"`
	BaseCurrencyMulti   float32     `json:"base_currency_multi"`
	DisplayInvoiceID    string      `json:"display_invoice_id"`
}

type APIResponseUser struct {
	Users map[int]User
}

type CustomTime struct {
	time.Time
}

type Membership struct {
	User     User     `json:"user`
	Accesses []Access `json:"accesses"`
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {

	//remove any extra " from the date string
	s := strings.Trim(string(b), "\"")

	regex := regexp.MustCompile("[0-9]")
	mask := regex.ReplaceAllString(s, "x")

	//try to parse from different formats
	switch mask {
	case "xxxx-xx-xx xx:xx:xx":
		t, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
		ct.Time = t
	case "xxxx-xx-xx":
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
		ct.Time = t

	}

	return nil
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {

	return []byte(ct.Time.Format("2006-01-02 15:04:05")), nil
}
