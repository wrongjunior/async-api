package models

type PlacesResponse struct {
	Features []struct {
		Properties struct {
			Xid  string `json:"xid"`
			Name string `json:"name"`
			Kind string `json:"kinds"`
		} `json:"properties"`
	} `json:"features"`
}

type PlaceDescriptionResponse struct {
	Name        string `json:"name"`
	Description string `json:"wikipedia_extracts"`
}
