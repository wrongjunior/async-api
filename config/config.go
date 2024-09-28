package config

const (
	GraphhopperAPIKey = "74e86a83-1a91-4f80-9291-f56a2f0378b3"
	OpenWeatherAPIKey = "ecacbd3f749199d911e168b9d5b86af3"
	OpenTripMapAPIKey = "5ae2e3f221c38a28845f05b6ed5af70c89b75a6de9898e0c5c4fd2af"

	GraphhopperAPIURL = "https://graphhopper.com/api/1/geocode?q=%s&locale=ru&key=%s"
	OpenWeatherAPIURL = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s"

	OpenTripMapAPIURL         = "https://api.opentripmap.com/0.1/en/places/radius?radius=%d&lon=%f&lat=%f&kinds=interesting_places&limit=%d&format=geojson&apikey=%s"
	OpenTripMapDescriptionURL = "https://api.opentripmap.com/0.1/en/places/xid/%s?apikey=%s"
)
