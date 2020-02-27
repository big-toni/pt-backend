package service

import (
	courier "../courier"
)

func GetParcelData(parcelNumber string) (string, bool) {
	return courier.GetGlobalCanaioData(parcelNumber)
}

func ResolveCourier(parcelNumber string) ([]string, bool) {
	return courier.ResolveCourier(parcelNumber)
}
