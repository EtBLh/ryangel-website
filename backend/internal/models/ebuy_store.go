package models

// EbuyStore represents an ebuy store.
type EbuyStore struct {
	StoreID     string  `json:"store_id"`
	StoreName   string  `json:"store_name"`
	Type        string  `json:"type"`
	OfficeHours string  `json:"office_hours"`
	Address     string  `json:"address"`
	AddressEn   string  `json:"address_en"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}