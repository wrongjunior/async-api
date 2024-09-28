package main

import (
	"async-api/internal/api"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// Ввод от пользователя
	var query string
	fmt.Println("Enter location name (e.g., Friedrichstraße, Berlin):")
	fmt.Scanln(&query)

	// Контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Поиск локаций по запросу
	locations, err := api.SearchLocations(ctx, query)
	if err != nil {
		log.Fatal("Error fetching locations:", err)
	}

	// Проверяем, есть ли результаты
	if len(locations) == 0 {
		fmt.Println("No locations found")
		return
	}

	// Выводим список найденных локаций
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
}
