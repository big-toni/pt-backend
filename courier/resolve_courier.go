package courier

import (
	"regexp"
)

type inquiry struct {
	Name  string
	Regex string
}

var couriers = []inquiry{
	{Name: "ups", Regex: "/^1Z[0-9A-Z]{16}$/"},
	{Name: "ups", Regex: "/^(H|T|J|K|F|W|M|Q|A)\\d{10}$/"},
	{Name: "amazon", Regex: "/^1\\d{2}-\\d{7}-\\d{7}:\\d{13}$/"},
	{Name: "fedex", Regex: "/^\\d{12}$/"},
	{Name: "fedex", Regex: "/^\\d{15}$/"},
	{Name: "fedex", Regex: "/^\\d{20}$/"},
	{Name: "usps", Regex: "/^\\d{20}$/"},
	{Name: "usps", Regex: "/^02\\d{18}$/"},
	{Name: "fedex", Regex: "/^02\\d{18}$/"},
	{Name: "fedex", Regex: "/^DT\\d{12}$/"},
	{Name: "fedex", Regex: "/^927489\\d{16}$/"},
	{Name: "fedex", Regex: "/^926129\\d{16}$/"},
	{Name: "upsmi", Regex: "/^927489\\d{16}$/"},
	{Name: "upsmi", Regex: "/^926129\\d{16}$/"},
	{Name: "upsmi", Regex: "/^927489\\d{20}$/"},
	{Name: "fedex", Regex: "/^96\\d{20}$/"},
	{Name: "usps", Regex: "/^927489\\d{16}$/"},
	{Name: "usps", Regex: "/^926129\\d{16}$/"},
	{Name: "fedex", Regex: "/^7489\\d{16}$/"},
	{Name: "fedex", Regex: "/^6129\\d{16}$/"},
	{Name: "usps", Regex: "/^(91|92|93|94|95|96)\\d{20}$/"},
	{Name: "usps", Regex: "/^\\d{26}$/"},
	{Name: "usps", Regex: "/^420\\d{27}$/"},
	{Name: "usps", Regex: "/^420\\d{31}$/"},
	{Name: "dhlgm", Regex: "/^420\\d{27}$/"},
	{Name: "dhlgm", Regex: "/^420\\d{31}$/"},
	{Name: "dhlgm", Regex: "/^94748\\d{17}$/"},
	{Name: "dhlgm", Regex: "/^93612\\d{17}$/"},
	{Name: "dhlgm", Regex: "/^GM\\d{16}/"},
	{Name: "usps", Regex: "/^[A-Z]{2}\\d{9}[A-Z]{2}$/"},
	{Name: "canadapost", Regex: "/^\\d{16}$/"},
	{Name: "lasership", Regex: "/^L[A-Z]\\d{8}$/"},
	{Name: "lasership", Regex: "/^1LS\\d{12}/"},
	{Name: "lasership", Regex: "/^Q\\d{8}[A-Z]/"},
	{Name: "ontrac", Regex: "/^(C|D)\\d{14}$/"},
	{Name: "prestige", Regex: "/^P[A-Z]{1}\\d{8}/"},
	{Name: "a1intl", Regex: "/^AZ.\\d+/"},
	{Name: "fake", Regex: "([0-9]+)"},
}

func ResolveCourier(trackingNumber string) (string, bool) {
	for _, i := range couriers {
		matched, _ := regexp.MatchString(i.Regex, trackingNumber)
		if matched {
			return i.Name, true
		}
	}
	return "", false
}
