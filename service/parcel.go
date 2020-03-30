package service

import (
	courier "../courier"
)

// GetParcelData func
func GetParcelData(parcelNumber string) (string, bool) {
	return courier.GetGlobalCanaioData(parcelNumber)
}

// ResolveCourier func
func ResolveCourier(parcelNumber string) ([]string, bool) {
	return courier.ResolveCourier(parcelNumber)
}
