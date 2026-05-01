package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type PointData struct {
	Id       string `json:"id"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Forecast         string `json:"forecast"`
		ForecastHourly   string `json:"forecastHourly"`
		Timezone         string `json:"timeZone"`
		AstronomicalData struct {
			Sunrise string `json:"sunrise"`
			Sunset  string `json:"sunset"`
		} `json:"astronomicalData"`
	} `json:"properties"`
}

func request_pointdata(coords []float64) (*PointData, error) {
	if coords[0] < 0 || coords[1] > 0 {
		return nil, errors.New("Provide coordinates in latitude, longitude order")
	}
	s := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", coords[0], coords[1])
	body, err := make_api_request(s)
	if err != nil {
		return nil, err
	}
	var w *PointData
	err = json.Unmarshal(body, &w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func download_and_save_pointdata(coords []float64) error {
	if len(coords) != 2 {
		return errors.New("Coordinates not 2-tuple")
	}
	if _, err := os.Stat("./data/points.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Downloading point data...")
		s := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", coords[0], coords[1])
		body, err := make_api_request(s)
		if err != nil {
			return err
		}
		err = os.WriteFile("./data/points.json", body, 0644)
		if err != nil {
			return err
		}
		fmt.Println("Station data saved successfully!")
	} else {
		fmt.Println("Warning: Data already exists")
	}
	return nil
}

func print_point_data(w *PointData) error {
	if w == nil {
		return errors.New("Pointer to PointData is nil")
	}
	fmt.Println("ID:", w.Id)
	fmt.Println("Coordinates:", w.Geometry.Coordinates)
	fmt.Println("Forecast url:", w.Properties.Forecast)
	fmt.Println("Sunrise Time:", w.Properties.AstronomicalData.Sunrise)
	fmt.Println("Sunset Time:", w.Properties.AstronomicalData.Sunset)

	return nil
}
