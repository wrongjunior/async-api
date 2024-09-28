package main

import (
	"async-api/internal/api"
	"async-api/internal/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// PlaceDetail представляет детальную информацию об интересном месте
type PlaceDetail struct {
	Xid         string
	Name        string
	Kind        string
	Description string
}

// FinalResult объединяет все данные для отображения
type FinalResult struct {
	Location models.LocationResult
	Weather  models.WeatherResponse
	Places   []PlaceDetail
}

// SearchLocationsResult представляет результат поиска локаций
type SearchLocationsResult struct {
	Locations []models.LocationResult
	Err       error
}

// FetchLocationDetailsResult представляет результат получения деталей локации
type FetchLocationDetailsResult struct {
	Result FinalResult
	Err    error
}

// SearchLocationsAsync выполняет асинхронный поиск локаций
func SearchLocationsAsync(ctx context.Context, query string) <-chan SearchLocationsResult {
	resultChan := make(chan SearchLocationsResult, 1)
	go func() {
		locations, err := api.SearchLocations(ctx, query)
		resultChan <- SearchLocationsResult{
			Locations: locations,
			Err:       err,
		}
	}()
	return resultChan
}

// FetchLocationDetailsAsync выполняет асинхронное получение деталей выбранной локации
func FetchLocationDetailsAsync(ctx context.Context, location models.LocationResult) <-chan FetchLocationDetailsResult {
	resultChan := make(chan FetchLocationDetailsResult, 1)
	go func() {
		var wg sync.WaitGroup
		var weather models.WeatherResponse
		var weatherErr error

		// Асинхронно получить погоду
		wg.Add(1)
		go func() {
			defer wg.Done()
			weather, weatherErr = api.GetWeather(ctx, location.Lat, location.Lon)
		}()

		// Асинхронно получить интересные места
		placesResponse, placesErr := api.GetPlaces(ctx, location.Lat, location.Lon, 1000, 10)
		if placesErr != nil {
			resultChan <- FetchLocationDetailsResult{Err: placesErr}
			return
		}

		// Асинхронно получить описания для каждого места
		places := make([]PlaceDetail, len(placesResponse.Features))
		var placeWg sync.WaitGroup
		for i, feature := range placesResponse.Features {
			placeWg.Add(1)
			go func(i int, feature models.Place) {
				defer placeWg.Done()
				descResponse, err := api.GetPlaceDescription(ctx, feature.Properties.Xid)
				description := "Description not available"
				if err == nil && descResponse.WikipediaExtracts.Text != "" {
					description = descResponse.WikipediaExtracts.Text
				}
				places[i] = PlaceDetail{
					Xid:         feature.Properties.Xid,
					Name:        feature.Properties.Name,
					Kind:        feature.Properties.Kind,
					Description: description,
				}
			}(i, feature)
		}
		placeWg.Wait()

		// Ждем завершения получения погоды
		wg.Wait()
		if weatherErr != nil {
			resultChan <- FetchLocationDetailsResult{Err: weatherErr}
			return
		}

		finalResult := FinalResult{
			Location: location,
			Weather:  weather,
			Places:   places,
		}
		resultChan <- FetchLocationDetailsResult{Result: finalResult}
	}()
	return resultChan
}

func main() {
	// Ввод от пользователя
	var query string
	fmt.Println("Enter location name (e.g., Friedrichstraße, Berlin):")
	fmt.Scanln(&query)

	// Контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Асинхронный поиск локаций
	searchResultChan := SearchLocationsAsync(ctx, query)

	// Ожидание результата поиска
	searchResult := <-searchResultChan
	if searchResult.Err != nil {
		log.Fatal("Error fetching locations:", searchResult.Err)
	}
	locations := searchResult.Locations
	if len(locations) == 0 {
		fmt.Println("No locations found")
		return
	}

	// Вывод найденных локаций
	fmt.Println("Found locations:")
	for i, location := range locations {
		fmt.Printf("%d. %s, %s, %s, %s %s (Lat: %.6f, Lon: %.6f)\n",
			i+1,
			location.Name,
			location.City,
			location.Street,
			location.HouseNumber,
			location.Postcode,
			location.Lat,
			location.Lon)
	}

	// Выбор локации пользователем
	var choice int
	fmt.Printf("Select a location by number (1-%d): ", len(locations))
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(locations) {
		log.Fatal("Invalid selection")
	}
	selectedLocation := locations[choice-1]

	// Асинхронное получение деталей выбранной локации
	detailsResultChan := FetchLocationDetailsAsync(ctx, selectedLocation)

	// Ожидание результата получения деталей
	detailsResult := <-detailsResultChan
	if detailsResult.Err != nil {
		log.Fatal("Error fetching location details:", detailsResult.Err)
	}
	finalResult := detailsResult.Result

	// Конвертация температуры из Кельвинов в Цельсии
	tempCelsius := finalResult.Weather.Main.Temp - 273.15

	// Вывод окончательных результатов
	fmt.Printf("\nLocation: %s, %s, %s, %s %s\n", finalResult.Location.Name, finalResult.Location.City, finalResult.Location.Street, finalResult.Location.HouseNumber, finalResult.Location.Postcode)
	fmt.Printf("Coordinates: Lat: %.6f, Lon: %.6f\n", finalResult.Location.Lat, finalResult.Location.Lon)
	fmt.Printf("Weather: %.2f°C\n", tempCelsius)
	fmt.Println("Interesting Places:")
	for _, place := range finalResult.Places {
		fmt.Printf("- %s (%s): %s\n", place.Name, place.Kind, place.Description)
	}
}
