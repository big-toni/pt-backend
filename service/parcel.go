package service

import (
	courier "../courier"
)

func GetParcelData(parcelNumber string) (string, bool) {
	return courier.GetGlobalCanaioData(parcelNumber)
}
