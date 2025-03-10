package amember

import (
	"database/sql"
	"regexp"
	"strings"
	"time"
)

type DBUser struct {
	Added              sql.NullTime `json:"added"`
	AffAdded           interface{}  `json:"aff_added"`
	AffCustomRedirect  int          `json:"aff_custom_redirect"`
	AffID              interface{}  `json:"aff_id"`
	AffPayoutType      interface{}  `json:"aff_payout_type"`
	City               string       `json:"city"`
	Comment            string       `json:"comment"`
	Country            string       `json:"country"`
	DisableLockUntil   sql.NullTime `json:"disable_lock_until"`
	Email              string       `json:"email"`
	IAgree             int          `json:"i_agree"`
	IsAffiliate        int          `json:"is_affiliate"`
	IsApproved         int          `json:"is_approved"`
	IsLocked           int          `json:"is_locked"`
	Lang               string       `json:"lang"`
	LastIP             string       `json:"last_ip"`
	LastLogin          sql.NullTime `json:"last_login"`
	LastSession        string       `json:"last_session"`
	LastUserAgent      string       `json:"last_user_agent"`
	Login              string       `json:"login"`
	NameF              string       `json:"name_f"`
	NameL              string       `json:"name_l"`
	NeedSessionRefresh int          `json:"need_session_refresh"`
	Pass               string       `json:"pass"`
	PassDattm          sql.NullTime `json:"pass_dattm"`
	Phone              string       `json:"phone"`
	RememberKey        string       `json:"remember_key"`
	RemoteAddr         string       `json:"remote_addr"`
	ResellerID         int          `json:"reseller_id"`
	SavedFormID        int          `json:"saved_form_id"`
	SignupEmailSent    int          `json:"signup_email_sent"`
	State              string       `json:"state"`
	Status             int          `json:"status"`
	Street             string       `json:"street"`
	Street2            string       `json:"street2"`
	Unsubscribed       int          `json:"unsubscribed"`
	UserAgent          string       `json:"user_agent"`
	UserID             int          `json:"user_id"`
	Zip                string       `json:"zip"`
	//StripeCcExpires    string      `json:"stripe_cc_expires"`
	//StripeCcMasked string      `json:"stripe_cc_masked"`
	//StripeToken    string      `json:"stripe_token"`
	CompanyName    string       `json:"company_name"`
	CompanyAddress string       `json:"company_address"`
	TaxID          string       `json:"taxid"`
	ExpiredAt      sql.NullTime `json:"expired_at"`
}

type ViewUser struct {
	UserID                  int          `json:"userId"`
	Username                string       `json:"username"`
	FirstName               string       `json:"first_name"`
	LastName                string       `json:"last_name"`
	Email                   string       `json:"email"`
	SignupDate              sql.NullTime `json:"signup_date"`
	SubscriptionStatus      string       `json:"subscriptionStatus"`
	ClickID                 string       `json:"click_id"`
	MobilePhone             string       `json:"mobile_phone"`
	SubscriptionPlan        string       `json:"subscription_plan"`
	ProductName             string       `json:"product_name"`
	ExpirationDate          sql.NullTime `json:"expiration_date"`
	TotalMonths             int          `json:"total_months"`
	TotalDays               int          `json:"total_days"`
	TotalDaysExcludingTrial int          `json:"total_days_excluding_trial"`
	TotalPayments           float32      `json:"total_payments"`
	FirstPayment            sql.NullTime `json:"first_payment"`
	LastPayment             sql.NullTime `json:"last_payment"`
	HowDidYouHear           string       `json:"how_did_you_hear"`
	PreferredContactMethod  string       `json:"preferred_contact_method"`
	PreferredContact        string       `json:"preferred_contact"`
	PaymentsLast3Months     float32      `json:"payments_last_3_months"`
	IsTopPayingUser         string       `json:"is_top_paying_user"`
	CancellationDate        sql.NullTime `json:"cancellation_date"`
	LastUpdated             sql.NullTime `json:"last_updated"`
}

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
	Unsubscribed       int         `json:"unsubscribed"`
	UserAgent          string      `json:"user_agent"`
	UserID             int         `json:"user_id"`
	Zip                string      `json:"zip"`
	//StripeCcExpires    string      `json:"stripe_cc_expires"`
	//StripeCcMasked string      `json:"stripe_cc_masked"`
	//StripeToken    string      `json:"stripe_token"`
	CompanyName    string     `json:"company_name"`
	CompanyAddress string     `json:"company_address"`
	TaxID          string     `json:"taxid"`
	ExpiredAt      CustomTime `json:"expired_at"`
}

type Invoice struct {
	InvoiceID         int           `json:"invoice_id"`
	UserID            int           `json:"user_id"`
	PaysysID          string        `json:"paysys_id"`
	Currency          string        `json:"currency"`
	FirstSubtotal     float32       `json:"first_subtotal"`
	FirstDiscount     float32       `json:"first_discount"`
	FirstTax          float32       `json:"first_tax"`
	FirstShipping     float32       `json:"first_shipping"`
	FirstTotal        float32       `json:"first_total"`
	FirstPeriod       string        `json:"first_period"`
	RebillTimes       int           `json:"rebill_times"`
	SecondSubtotal    float32       `json:"second_subtotal"`
	SecondDiscount    float32       `json:"second_discount"`
	SecondTax         float32       `json:"second_tax"`
	SecondShipping    float32       `json:"second_shipping"`
	SecondTotal       float32       `json:"second_total"`
	SecondPeriod      string        `json:"second_period"`
	TaxRate           interface{}   `json:"tax_rate"`
	TaxType           interface{}   `json:"tax_type"`
	TaxTitle          interface{}   `json:"tax_title"`
	Status            int           `json:"status"`
	CouponID          int           `json:"coupon_id"`
	CouponCode        string        `json:"coupon_code"`
	DiscountFirst     float32       `json:"discount_first"`
	DiscountSecond    float32       `json:"discount_second"`
	IsConfirmed       int           `json:"is_confirmed"`
	PublicID          string        `json:"public_id"`
	InvoiceKey        string        `json:"invoice_key"`
	TmAdded           *CustomTime   `json:"tm_added"`
	TmStarted         *CustomTime   `json:"tm_started"`
	TmCancelled       *CustomTime   `json:"tm_cancelled"`
	RebillDate        *CustomTime   `json:"rebill_date"`
	DueDate           *CustomTime   `json:"due_date"`
	Terms             interface{}   `json:"terms"`
	Comment           interface{}   `json:"comment"`
	BaseCurrencyMulti string        `json:"base_currency_multi"`
	SavedFormID       int           `json:"saved_form_id"`
	AffID             interface{}   `json:"aff_id"`
	KeywordID         interface{}   `json:"keyword_id"`
	RemoteAddr        interface{}   `json:"remote_addr"`
	Nested            InvoiceNested `json:"nested"`
}

type InvoiceNested struct {
	Access          []Access  `json:"access"`
	InvoiceItems    []Item    `json:"invoice-items"`
	InvoicePayments []Payment `json:"invoice-payments"`
}

type Access struct {
	AccessID           int       `json:"access_id"`
	InvoiceID          int       `json:"invoice_id"`
	InvoicePublicID    string    `json:"invoice_public_id"`
	InvoicePaymentID   int       `json:"invoice_payment_id"`
	InvoiceItemID      int       `json:"invoice_item_id"`
	UserID             int       `json:"user_id"`
	ProductID          int       `json:"product_id"`
	TransactionID      string    `json:"transaction_id"`
	BeginDate          time.Time `json:"begin_date"`
	ExpireDate         time.Time `json:"expire_date"`
	Qty                int       `json:"qty"`
	Comment            string    `json:"comment"`
	ProductTitle       string    `json:"product_title"`
	Status             bool      `json:"status"`
	ProductDescription string    `json:"product_description"`
	Spend              float32   `json:"spend"`
	SpendCoveredByPlan float32   `json:"spend_covered_by_plan"`
	Overage            float32   `json:"overage"`
	ProjectedSpend     float32   `json:"projected_spend"`
	ProjectedOverage   float32   `json:"projected_overage"`
}

type DBAccess struct {
	AccessID           int            `json:"access_id"`
	InvoiceID          int            `json:"invoice_id"`
	InvoicePublicID    sql.NullString `json:"invoice_public_id"`
	InvoicePaymentID   int            `json:"invoice_payment_id"`
	InvoiceItemID      int            `json:"invoice_item_id"`
	UserID             int            `json:"user_id"`
	ProductID          int            `json:"product_id"`
	TransactionID      string         `json:"transaction_id"`
	BeginDate          sql.NullTime   `json:"begin_date"`
	ExpireDate         sql.NullTime   `json:"expire_date"`
	Qty                int            `json:"qty"`
	Comment            string         `json:"comment"`
	ProductTitle       string         `json:"product_title"`
	Status             bool           `json:"status"`
	ProductDescription string         `json:"product_description"`
	Spend              float32        `json:"spend"`
	SpendCoveredByPlan float32        `json:"spend_covered_by_plan"`
	Overage            float32        `json:"overage"`
	ProjectedSpend     float32        `json:"projected_spend"`
	ProjectedOverage   float32        `json:"projected_overage"`
}

type Payment struct {
	ConversionTrackDone    int       `json:"conversion-track-done"`
	GoogleAnalyticsDone    int       `json:"google-analytics-done"`
	InvoicePaymentID       int       `json:"invoice_payment_id"`
	InvoiceID              int       `json:"invoice_id"`
	InvoicePublicID        string    `json:"invoice_public_id"`
	UserID                 int       `json:"user_id"`
	PaysysID               string    `json:"paysys_id"`
	ReceiptID              string    `json:"receipt_id"`
	TransactionID          string    `json:"transaction_id"`
	Dattm                  time.Time `json:"dattm"`
	Currency               string    `json:"currency"`
	Amount                 float32   `json:"amount"`
	Discount               float32   `json:"discount"`
	Tax                    float32   `json:"tax"`
	Shipping               float32   `json:"shipping"`
	RefundDattm            time.Time `json:"refund_dattm"`
	RefundAmount           float32   `json:"refund_amount"`
	BaseCurrencyMulti      float32   `json:"base_currency_multi"`
	DisplayInvoiceID       string    `json:"display_invoice_id"`
	Username               string    `json:"username"`
	PaymentItemDescription string    `json:"payment_item_description"`
	PaymentItemTitle       string    `json:"payment_item_title"`
	Refunded               bool      `json:"refunded"`
}

type Product struct {
	CartDescription      string     `json:"cart_description"`
	Comment              string     `json:"comment"`
	Currency             string     `json:"currency"`
	DefaultBillingPlanID int        `json:"default_billing_plan_id"`
	Description          string     `json:"description"`
	Img                  int        `json:"img"`
	ImgCartPath          string     `json:"img_cart_path"`
	ImgDetailPath        string     `json:"img_detail_path"`
	ImgOrigPath          string     `json:"img_orig_path"`
	ImgPath              string     `json:"img_path"`
	IsArchived           int        `json:"is_archived"`
	IsDisabled           int        `json:"is_disabled"`
	IsTangible           int        `json:"is_tangible"`
	MetaDescription      string     `json:"meta_description"`
	MetaKeywords         string     `json:"meta_keywords"`
	MetaRobots           string     `json:"meta_robots"`
	MetaTitle            string     `json:"meta_title"`
	Path                 string     `json:"path"`
	PaysysID             string     `json:"paysys_id"`
	PreventIfOther       string     `json:"prevent_if_other"`
	ProductID            int        `json:"product_id"`
	RenewalGroup         string     `json:"renewal_group"`
	RequireOther         string     `json:"require_other"`
	SortOrder            int        `json:"sort_order"`
	StartDate            CustomTime `json:"start_date"`
	StartDateFixed       CustomTime `json:"start_date_fixed"`
	Tags                 string     `json:"tags"`
	TaxDigital           string     `json:"tax_digital"`
	TaxGroup             string     `json:"tax_group"`
	TaxRateGroup         string     `json:"tax_rate_group"`
	ThanksRedirectURL    string     `json:"thanks_redirect_url"`
	Title                string     `json:"title"`
	TrialGroup           string     `json:"trial_group"`
	URL                  string     `json:"url"`
}

type Item struct {
	BillingPlanData string      `json:"billing_plan_data"`
	BillingPlanID   string      `json:"billing_plan_id"`
	Currency        string      `json:"currency"`
	FirstDiscount   string      `json:"first_discount"`
	FirstPeriod     string      `json:"first_period"`
	FirstPrice      string      `json:"first_price"`
	FirstShipping   string      `json:"first_shipping"`
	FirstTax        string      `json:"first_tax"`
	FirstTotal      string      `json:"first_total"`
	InvoiceID       string      `json:"invoice_id"`
	InvoiceItemID   int64       `json:"invoice_item_id"`
	InvoicePublicID string      `json:"invoice_public_id"`
	IsCountable     string      `json:"is_countable"`
	IsTangible      interface{} `json:"is_tangible"`
	ItemDescription string      `json:"item_description"`
	ItemID          string      `json:"item_id"`
	ItemTitle       string      `json:"item_title"`
	ItemType        string      `json:"item_type"`
	Option1         interface{} `json:"option1"`
	Option2         interface{} `json:"option2"`
	Option3         interface{} `json:"option3"`
	Options         interface{} `json:"options"`
	Qty             string      `json:"qty"`
	RebillTimes     string      `json:"rebill_times"`
	SecondDiscount  string      `json:"second_discount"`
	SecondPeriod    string      `json:"second_period"`
	SecondPrice     string      `json:"second_price"`
	SecondShipping  string      `json:"second_shipping"`
	SecondTax       string      `json:"second_tax"`
	SecondTotal     string      `json:"second_total"`
	TaxGroup        string      `json:"tax_group"`
	TaxRate         interface{} `json:"tax_rate"`
	VariableQty     string      `json:"variable_qty"`
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
	case "xxxx-xx-xx xx:xx:xx +xxxx", "xxxx-xx-xx xx:xx:xx -xxxx":
		t, err := time.Parse("2006-01-02 15:04:05 -0700", s)
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
