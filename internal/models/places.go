package models

type PlacesResponse struct {
	Features []Place `json:"features"`
}

type Place struct {
	Properties struct {
		Xid  string `json:"xid"`
		Name string `json:"name"`
		Kind string `json:"kinds"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

type PlaceDescriptionResponse struct {
	Name              string `json:"name"`
	WikipediaExtracts struct {
		Text string `json:"text"`
	} `json:"wikipedia_extracts"`
}
