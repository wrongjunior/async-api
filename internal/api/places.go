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

// GetPlaces получает список интересных мест по координатам
func GetPlaces(ctx context.Context, lat, lon float64, radius int, limit int) (models.PlacesResponse, error) {
	var result models.PlacesResponse
	url := fmt.Sprintf(config.OpenTripMapAPIURL, radius, lon, lat, limit, config.OpenTripMapAPIKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return result, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error fetching places:", err)
		return result, err
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		log.Println("Invalid content type received:", resp.Header.Get("Content-Type"))
		return result, fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

// GetPlaceDescription получает описание места по его xid
func GetPlaceDescription(ctx context.Context, xid string) (models.PlaceDescriptionResponse, error) {
	var result models.PlaceDescriptionResponse
	url := fmt.Sprintf(config.OpenTripMapDescriptionURL, xid, config.OpenTripMapAPIKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return result, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error fetching place description:", err)
		return result, err
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		log.Println("Invalid content type received:", resp.Header.Get("Content-Type"))
		return result, fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
