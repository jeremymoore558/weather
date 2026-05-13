package main

import (
	"fmt"
	"log"
	"os"
	"weather/dbutils"
	"weather/noa"
)

func main() {
	coords := []float64{41.8832, -87.6324} // Chicago coordinates

	if len(os.Args) == 1 {
		help_page()
		return
	}

	if os.Args[1] == "dbinstall" {
		err := dbutils.Dbinstall()
		log_err(err)
		err = dbutils.Dbcheck()
		log_err(err)
		return
	}

	if os.Args[1] == "forecast" {
		err := noa.Noa_request(coords)
		log_err(err)
	}

	dbutils.Get_city_coordinates("chicago")
}

func log_err(err error) {
	log.Fatal(err)
}

func help_page() {
	fmt.Println("<weather> help page:")
	fmt.Println("Commands:")
	fmt.Println("dbinstall: reinstalls cities database. Run if data become corrupted")
}
