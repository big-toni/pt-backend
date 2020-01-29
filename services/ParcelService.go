package services

import (
	providers "../providers"
)

func GetParcelData(parcelNumber string) (string, bool) {
	return providers.GetGlobalCanaioData(parcelNumber)
}
