package api

import (
	"async-api/config"
	"async-api/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// GetWeather получает погоду по координатам
func GetWeather(ctx context.Context, lat, lon float64) (models.WeatherResponse, error) {
	var result models.WeatherResponse
	url := fmt.Sprintf(config.OpenWeatherAPIURL, lat, lon, config.OpenWeatherAPIKey)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error fetching weather:", err)
		return result, err
	}
	defer resp.Body.Close()

	// Изменение проверки на Content-Type, чтобы допускать application/json с параметрами
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		log.Println("Invalid content type received:", resp.Header.Get("Content-Type"))
		return result, fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
