package noa

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Forecast struct {
	Properties struct {
		Periods []struct {
			Number          int    `json:"number"`
			Name            string `json:"name"`
			Time            string `json:"startTime"`
			IsDaytime       bool   `json:"isDaytime"`
			Temperature     int    `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
			ProbOfPrecip    struct {
				UnitCode string `json:"unitCode"`
				Value    int    `json:"value"`
			} `json:"probabilityOfPrecipitation"`
			Windspeed        string `json:"windSpeed"`
			WindDirection    string `json:"windDirection"`
			Icon             string `json:"icon"`
			ShortForecast    string `json:"shortForecast"`
			DetailedForecast string `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func display_forecast(f *Forecast, k int) error {
	if k > len(f.Properties.Periods) {
		s := fmt.Sprintf("Requested forecast period %v, must be less than %v",
			k, len(f.Properties.Periods))
		return errors.New(s)
	}

	fmt.Println("Time:", f.Properties.Periods[k].Time)
	if f.Properties.Periods[k].IsDaytime {
		fmt.Println("Day")
	} else {
		fmt.Println("Night")
	}
	fmt.Println("Temperature:", f.Properties.Periods[k].Temperature, f.Properties.Periods[k].TemperatureUnit)
	fmt.Println("Precipitation:", f.Properties.Periods[k].ProbOfPrecip.Value, "percent")
	fmt.Println("Windspeed:", f.Properties.Periods[k].Windspeed, f.Properties.Periods[k].WindDirection)
	fmt.Println("Forecast:", f.Properties.Periods[k].ShortForecast)

	return nil
}

func test_forecast() {
	coords := []float64{41.8832, -87.6324}
	download_and_save_pointdata(coords)
	// Read in station data to var body []byte
	body, err := os.ReadFile("./data/points.json")
	check_err(err)

	// Parse point data
	var w PointData
	err = json.Unmarshal(body, &w)
	check_err(err)
	//err = print_point_data(&w)
	//check_err(err)

	// Make request for forecast data
	download_and_save_forecast(&w)
	body, err = os.ReadFile("./data/forecast.json")
	check_err(err)

	// Parse forecast data
	var f Forecast
	err = json.Unmarshal(body, &f)
	check_err(err)
	err = display_forecast(&f, 0)
	check_err(err)
}

func display_forecast_entries(f *Forecast, k int) {
	for i := range min(k, len(f.Properties.Periods)) {
		display_forecast(f, i)
		fmt.Println("--------------------------------")
	}
}

func unmarshall_forecast(body []byte) (*Forecast, error) {
	var f *Forecast
	err := json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func request_hourly_forecast(w *PointData) (*Forecast, error) {
	if w == nil {
		return nil, errors.New("Point data is nil")
	}
	s := w.Properties.ForecastHourly
	body, err := make_api_request(s)
	if err != nil {
		return nil, err
	}

	f, err := unmarshall_forecast(body)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func request_forecast(w *PointData) (*Forecast, error) {
	if w == nil {
		return nil, errors.New("Point data is nil")
	}
	s := w.Properties.Forecast
	body, err := make_api_request(s)
	if err != nil {
		return nil, err
	}

	f, err := unmarshall_forecast(body)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func download_and_save_forecast(w *PointData) error {
	if w == nil {
		return errors.New("Nil pointer provided")
	}
	if _, err := os.Stat("./data/forecast.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Downloading forecast...")
		s := w.Properties.Forecast
		body, err := make_api_request(s)
		if err != nil {
			return err
		}
		err = os.WriteFile("./data/forecast.json", body, 0644)
		if err != nil {
			return err
		}
		fmt.Println("Forecast data saved successfully!")
	} else {
		fmt.Println("Warning: Data already exists")
	}

	return nil
}
