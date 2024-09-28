package models

type GeocodeResponse struct {
	Hits []struct {
		Point struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lng"` // Исправленный тег
		} `json:"point"`
		Name        string `json:"name"`
		Country     string `json:"country"`
		City        string `json:"city"`
		Street      string `json:"street"`
		HouseNumber string `json:"housenumber"`
		Postcode    string `json:"postcode"`
	} `json:"hits"`
}

type LocationResult struct {
	Name        string
	Country     string
	City        string
	Street      string
	HouseNumber string
	Postcode    string
	Lat         float64
	Lon         float64
}
