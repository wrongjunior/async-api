package models

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}
