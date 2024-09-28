package api

import (
	"async-api/config"
	"async-api/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// SearchLocations выполняет поиск локации и возвращает список найденных результатов
func SearchLocations(ctx context.Context, query string) ([]models.LocationResult, error) {
	var result models.GeocodeResponse
	url := fmt.Sprintf(config.GraphhopperAPIURL, query, config.GraphhopperAPIKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error fetching location:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Invalid response status: %d", resp.StatusCode)
		return nil, fmt.Errorf("invalid response status: %d", resp.StatusCode)
	}

	// Проверяем тип контента
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		log.Println("Invalid content type received:", resp.Header.Get("Content-Type"))
		return nil, fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
	}

	// Читаем тело ответа для отладки
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Geocode API Response:", string(body)) // Для отладки
	//
	// Декодируем JSON
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Преобразуем результат в удобный формат
	var locations []models.LocationResult
	for _, hit := range result.Hits {
		locations = append(locations, models.LocationResult{
			Name:        hit.Name,
			Country:     hit.Country,
			City:        hit.City,
			Street:      hit.Street,
			HouseNumber: hit.HouseNumber,
			Postcode:    hit.Postcode,
			Lat:         hit.Point.Lat,
			Lon:         hit.Point.Lon, // Теперь Lon корректно маппится
		})
	}

	return locations, nil
}
