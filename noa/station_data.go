package noa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Data structure for parsing JSON from /stations GET request
// Go's json.Unmarshal automatically populates struct with
// data that matches the json tag after each variable if the
// JSON file follows a compatible structure.
// NOA's JSON files are deeply nested, so this struct contains
// several nested structs. Oof.
type WeatherStations struct {
	T        string `json:"type"`
	Features []struct {
		Id         string `json:"id"`
		Properties struct {
			StatId   string `json:"stationIdentifier"`
			Name     string `json:"name"`
			Forecast string `json:"forecast"`
		} `json:"properties"`
		Geometry struct {
			T           string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func test_nearest_station() {
	test_coord := []float64{-87.63, 41.88}

	// Read in station data to var body []byte
	if _, err := os.Stat("./data/stations.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Downloading weather station data...")
		err := download_and_save_stations()
		check_err(err)
		fmt.Println("Station data saved successfully!")
	}
	body, err := os.ReadFile("./data/stations.json")

	// Parse body json to weatherstations struct
	var w WeatherStations
	err = json.Unmarshal(body, &w)
	check_err(err)

	// Make slice of all coordinates
	coords, err := coordinates_from_stationdata(&w)
	check_err(err)

	// Find nearest coordinate
	nearest_idx, err := argmin_distance(test_coord, coords)
	check_err(err)
	fmt.Println("Nearest index: ", nearest_idx)
	fmt.Println(w.Features[nearest_idx].Properties.Name)
	fmt.Println("Forecast: ", w.Features[nearest_idx].Properties.Forecast)

}

func coordinates_from_stationdata(w *WeatherStations) ([][]float64, error) {
	if len(w.Features) == 0 {
		err := errors.New("Invalid station data")
		return nil, err
	}

	if len(w.Features[0].Geometry.Coordinates) != 2 {
		err := errors.New("Invalid coordinates")
		return nil, err
	}

	coords := make([][]float64, 0, len(w.Features))
	for _, f := range w.Features {
		coords = append(coords, f.Geometry.Coordinates)
	}
	return coords, nil
}

func download_and_save_stations() error {
	body, err := get_weather_stations()
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile("stations.json", body, 0644)
	return err
}

func get_weather_stations() ([]byte, error) {
	// Get request from NOA api
	resp, err := http.Get("https://api.weather.gov/stations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
