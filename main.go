package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

type inputs struct {
	hourly      bool
	num_entries int
}

func main() {
	// In the future, find a way to autodetect coordinates.
	coords := []float64{41.8832, -87.6324} // Chicago coordinates

	// Handle command line inputs
	input_params := parse_inputs()
	num_entries := input_params.num_entries
	hourly := input_params.hourly

	// Retrieve and print weather data
	fmt.Printf("\033[2J")
	fmt.Printf("NOA Weather forecast:\n--------------------------------\n")
	point_data, err := request_pointdata(coords)
	check_err(err)

	if !hourly {
		fmt.Printf("Daily:\n--------------------------------\n")
		forecast, err := request_forecast(point_data)
		check_err(err)
		display_forecast_entries(forecast, num_entries)
	} else if hourly {
		fmt.Printf("Hourly:\n--------------------------------\n")
		forecast, err := request_hourly_forecast(point_data)
		check_err(err)
		display_forecast_entries(forecast, num_entries)
	}
}

func parse_inputs() *inputs {
	var p inputs
	var nFlag = flag.Int("n", 3, "Number of entries to display")
	var hFlag = flag.Bool("h", false, "")
	flag.Parse()
	p.hourly = *hFlag
	p.num_entries = *nFlag
	return &p
}

func make_api_request(s string) ([]byte, error) {
	resp, err := http.Get(s)
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
