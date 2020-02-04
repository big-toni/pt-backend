package courier

import (
	"regexp"
)

type inquiry struct {
	name    string
	regex   string
	approve func(b string) (bool, bool)
}

func aproveUps(trk string) (bool, bool) {
	runes := []rune(trk)
	sum := 0

	for _, rune := range runes[2:16] {
		sum += int(rune - '0')
	}

	var checkdigit int
	if sum%10 > 0 {
		checkdigit = 10 - sum%10
	} else {
		checkdigit = 0
	}

	if checkdigit == int(runes[17] - '0') {
		return true, true
	}
	return false, false
}

var courierInquiries = []inquiry{
	{name: "ups", regex: `^1Z[0-9A-Z]{16}$`, approve: aproveUps},
	{name: "ups", regex: `^(H|T|J|K|F|W|M|Q|A)\d{10}$`},
	{name: "amazon", regex: `^1\d{2}-\d{7}-\d{7}:\d{13}$`},
	{name: "fedex", regex: `^\d{12}$`},
	{name: "fedex", regex: `^\d{15}$`},
	{name: "fedex", regex: `^\d{20}$`},
	{name: "usps", regex: `^\d{20}$`},
	{name: "usps", regex: `^02\d{18}$`},
	{name: "fedex", regex: `^02\d{18}$`},
	{name: "fedex", regex: `^DT\d{12}$`},
	{name: "fedex", regex: `^927489\d{16}$`},
	{name: "fedex", regex: `^926129\d{16}$`},
	{name: "upsmi", regex: `^927489\d{16}$`},
	{name: "upsmi", regex: `^926129\d{16}$`},
	{name: "upsmi", regex: `^927489\d{20}$`},
	{name: "fedex", regex: `^96\d{20}$`},
	{name: "usps", regex: `^927489\d{16}$`},
	{name: "usps", regex: `^926129\d{16}$`},
	{name: "fedex", regex: `^7489\d{16}$`},
	{name: "fedex", regex: `^6129\d{16}$`},
	{name: "usps", regex: `^(91|92|93|94|95|96)\d{20}$`},
	{name: "usps", regex: `^\d{26}$`},
	{name: "usps", regex: `^420\d{27}$`},
	{name: "usps", regex: `^420\d{31}$`},
	{name: "dhlgm", regex: `^420\d{27}$`},
	{name: "dhlgm", regex: `^420\d{31}$`},
	{name: "dhlgm", regex: `^94748\d{17}$`},
	{name: "dhlgm", regex: `^93612\d{17}$`},
	{name: "dhlgm", regex: `^GM\d{16}`},
	{name: "usps", regex: `^[A-Z]{2}\d{9}[A-Z]{2}$`},
	{name: "canadapost", regex: `^\d{16}$`},
	{name: "lasership", regex: `^L[A-Z]\d{8}$`},
	{name: "lasership", regex: `^1LS\d{12}`},
	{name: "lasership", regex: `^Q\d{8}[A-Z]`},
	{name: "ontrac", regex: `^(C|D)\d{14}$`},
	{name: "prestige", regex: `^P[A-Z]{1}\d{8}`},
	{name: "a1intl", regex: `^AZ.\d+`},
}

// ResolveCourier resolves courier
func ResolveCourier(trackingNumber string) (string, bool) {
	for _, i := range courierInquiries {
		matched, _ := regexp.MatchString(i.regex, trackingNumber)
		if matched {
			if i.approve != nil {
				ok, _ := i.approve(trackingNumber)
				if ok { 
					return i.name, true
				}
				return "", false
			}
			return i.name, true
		}
	}
	return "", false
}
