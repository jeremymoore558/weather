package main

import (
	"log"
	"weather/noa"
)

func main() {
	coords := []float64{41.8832, -87.6324} // Chicago coordinates
	err := noa.Noa_request(coords)
	if err != nil {
		log.Fatal(err)
	}
}
