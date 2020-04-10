package couriers

import (
	"fmt"
	"regexp"
	"strings"
)

type inquiry struct {
	name    string
	regex   string
	approve func(b string) (bool, bool)
}

func approveUps(trk string) (bool, bool) {
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

	if checkdigit == int(runes[17]-'0') {
		return true, true
	}
	return false, false
}

func checkDigit(trk string, multipliers []int, mod int) bool {
	midx := 0
	sum := 0
	runes := []rune(trk)
	var checkdigit int

	for _, char := range runes[0 : len(runes)-1] {
		sum += int(char-'0') * multipliers[midx]
		if midx == len(multipliers)-1 {
			midx = 0
		} else {
			midx++
		}
	}
	if mod == 11 {
		checkdigit = sum % 11
		if checkdigit == 10 {
			checkdigit = 0
		}
	}
	if mod == 10 {
		checkdigit = 0
		if (sum % 10) > 0 {
			checkdigit = (10 - sum%10)
		}
	}
	value := int(runes[len(runes)-1] - '0')
	return checkdigit == value
}

func approveUpsFreight(trk string) (bool, bool) {
	runes := []rune(trk)
	firstChar := int(runes[0]-63) % 10
	remaining := runes[1:]
	newtrk := fmt.Sprintf("%v%s", firstChar, string(remaining))
	if checkDigit(newtrk, []int{3, 1, 7}, 10) {
		return true, true
	}
	return false, false
}

func approveFedex12(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1, 7}, 11) {
		return true, true
	}
	return false, false
}

func approveFedexDoorTag(trk string) (bool, bool) {
	re := regexp.MustCompile(`^DT(\d{12})$`)
	newtrk := re.FindStringSubmatch(trk)[1]
	if checkDigit(newtrk, []int{3, 1, 7}, 11) {
		return true, false
	}
	return false, false
}

func approveFedexSmartPost(trk string) (bool, bool) {
	if checkDigit(fmt.Sprintf("%s%s", "91", trk), []int{3, 1}, 10) {
		return true, false
	}
	return false, false
}

func approveFedex15(trk string) (bool, bool) {
	if checkDigit(trk, []int{1, 2}, 10) {
		return true, true
	}
	return false, false
}

func approveFedex20(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1, 7}, 11) {
		return true, true
	} else if checkDigit(fmt.Sprintf("%s%s", "92", trk), []int{3, 1}, 10) {
		return true, true
	}
	return false, false
}

func approveUsps20(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1}, 10) {
		return true, true
	}
	return false, false
}

func approveFedex9622(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1, 7}, 11) {
		return true, true
	} else if checkDigit(trk[7:], []int{1, 3}, 10) {
		return true, true
	}
	return false, false
}

func approveUsps22(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1}, 10) {
		return true, false
	}
	return false, false
}

func approveUsps26(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1}, 10) {
		return true, false
	}
	return false, false
}

func approveUsps420Zip(trk string) (bool, bool) {
	re := regexp.MustCompile(`^420\d{5}(\d{22})$`)
	newtrk := re.FindStringSubmatch(trk)[1]
	if checkDigit(newtrk, []int{3, 1}, 10) {
		return true, false
	}
	return false, false
}

func approveUsps420ZipPlus4(trk string) (bool, bool) {
	re1 := regexp.MustCompile(`^420\d{9}(\d{22})$`)
	newtrk1 := re1.FindStringSubmatch(trk)[1]
	if checkDigit(newtrk1, []int{3, 1}, 10) {
		return true, false
	}

	re2 := regexp.MustCompile(`^420\d{5}(\d{26})$`)
	newtrk2 := re2.FindStringSubmatch(trk)[1]
	if checkDigit(newtrk2[7:], []int{3, 1}, 10) {
		return true, false
	}
	return false, false
}

func approveCanadaPost16(trk string) (bool, bool) {
	if checkDigit(trk, []int{3, 1}, 10) {
		return true, false
	}
	return false, false
}

func approveA1International(trk string) (bool, bool) {
	if len(trk) == 9 || len(trk) == 13 {
		return true, false
	}
	return false, false
}

var courierInquiries = []inquiry{
	{name: "ups", regex: `^1Z[0-9A-Z]{16}$`, approve: approveUps},
	{name: "ups", regex: `^(H|T|J|K|F|W|M|Q|A)\d{10}$`, approve: approveUpsFreight},
	{name: "amazon", regex: `^1\d{2}-\d{7}-\d{7}:\d{13}$`},
	{name: "fedex", regex: `^\d{12}$`, approve: approveFedex12},
	{name: "fedex", regex: `^\d{15}$`, approve: approveFedex15},
	{name: "fedex", regex: `^\d{20}$`, approve: approveFedex20},
	{name: "usps", regex: `^\d{20}$`, approve: approveUsps20},
	{name: "usps", regex: `^02\d{18}$`, approve: approveFedexSmartPost},
	{name: "fedex", regex: `^02\d{18}$`, approve: approveFedexSmartPost},
	{name: "fedex", regex: `^DT\d{12}$`, approve: approveFedexDoorTag},
	{name: "fedex", regex: `^927489\d{16}$`},
	{name: "fedex", regex: `^926129\d{16}$`},
	{name: "upsmi", regex: `^927489\d{16}$`},
	{name: "upsmi", regex: `^926129\d{16}$`},
	{name: "upsmi", regex: `^927489\d{20}$`},
	{name: "fedex", regex: `^96\d{20}$`, approve: approveFedex9622},
	{name: "usps", regex: `^927489\d{16}$`},
	{name: "usps", regex: `^926129\d{16}$`},
	{name: "fedex", regex: `^7489\d{16}$`},
	{name: "fedex", regex: `^6129\d{16}$`},
	{name: "usps", regex: `^(91|92|93|94|95|96)\d{20}$`, approve: approveUsps22},
	{name: "usps", regex: `^\d{26}$`, approve: approveUsps26},
	{name: "usps", regex: `^420\d{27}$`, approve: approveUsps420Zip},
	{name: "usps", regex: `^420\d{31}$`, approve: approveUsps420ZipPlus4},
	{name: "dhlgm", regex: `^420\d{27}$`, approve: approveUsps420Zip},
	{name: "dhlgm", regex: `^420\d{31}$`, approve: approveUsps420ZipPlus4},
	{name: "dhlgm", regex: `^94748\d{17}$`, approve: approveUsps22},
	{name: "dhlgm", regex: `^93612\d{17}$`, approve: approveUsps22},
	{name: "dhlgm", regex: `^GM\d{16}`},
	{name: "usps", regex: `^[A-Z]{2}\d{9}[A-Z]{2}$`},
	{name: "canadapost", regex: `^\d{16}$`, approve: approveCanadaPost16},
	{name: "lasership", regex: `^L[A-Z]\d{8}$`},
	{name: "lasership", regex: `^1LS\d{12}`},
	{name: "lasership", regex: `^Q\d{8}[A-Z]`},
	{name: "ontrac", regex: `^(C|D)\d{14}$`},
	{name: "prestige", regex: `^P[A-Z]{1}\d{8}`},
	{name: "a1intl", regex: `^AZ.\d+`, approve: approveA1International},
}

// ResolveCourier resolves courier
func ResolveCourier(trackingNumber string) ([]string, bool) {
	couriers := []string{}
	trk := strings.Replace(trackingNumber, " ", "", -1)
	for _, i := range courierInquiries {
		matched, _ := regexp.MatchString(i.regex, trk)
		if matched {
			if i.approve == nil {
				couriers = append(couriers, i.name)
			} else {
				ok, stop := i.approve(trk)
				if ok {
					couriers = append(couriers, i.name)
				}
				if stop {
					return couriers, true
				}
			}
		}
	}

	if len(couriers) > 0 {
		return couriers, true
	}
	return couriers, false
}
