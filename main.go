package main

import (
	"os"

	"github.com/paperclicks/gomember/amember"
)

func main() {

	//
	am := amember.New("https://membership.theoptimizer.io", "HUuL1WiRgA3ISZZZyyzo", os.Stdout)

	p := amember.Params{
		Filter: map[string]string{"user_id": "8"},
	}
	//am.Users(p)
	//am.Invoices(p)
	am.Accesses(p)
	//am.Payments(p)
}
