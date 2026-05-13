package dbutils

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type cityInfo struct {
	city       string
	state_name string
	lat        float64
	lng        float64
}

func main() {
	//	err := citiesdb_setup()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	err := query_city_table()
	if err != nil {
		log.Fatal(err)
	}
}

func Get_city_coordinates(city string) error {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT lat, lng FROM cities WHERE city = '%s';", city)
	rows, err := db.Query(query)

	var lat float64
	var lng float64
	var coords [2]float64

	rows.Next()
	rows.Scan(&lat, &lng)
	coords[0] = lat
	coords[1] = lng

	fmt.Printf("Coordinates: lat = %v, lng = %v\n", coords[0], coords[1])
	return nil
}

func Dbinstall() error {
	err := citiesdb_setup()
	if err != nil {
		return err
	}
	return nil
}

func Dbcheck() error {
	err := query_city_table()
	if err != nil {
		return err
	}
	return nil
}

func query_city_table() error {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		return err
	}
	query := `SELECT city, state_name, lat, lng FROM cities`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	var city_data cityInfo
	for range 10 {
		rows.Next()
		rows.Scan(&city_data.city, &city_data.state_name, &city_data.lat, &city_data.lng)
		fmt.Printf("City Info: %s, %s, %v, %v\n", city_data.city, city_data.state_name,
			city_data.lat, city_data.lng)
	}

	return nil
}

func citiesdb_setup() error {
	fmt.Println("Database utils")
	file, err := os.Open("./data/uscities.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	fmt.Printf("Setting up U.S. Cities Database:\n")
	err = create_city_table(records)
	if err != nil {
		return err
	}
	fmt.Printf("Database Setup Complete!")
	return nil
}

func create_city_table(records [][]string) error {
	db, err := sql.Open("sqlite3", "./data/database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	create_table := `CREATE TABLE IF NOT EXISTS cities (			
city,city_ascii,state_id,state_name,county_fips,county_name,lat,lng,population,density,source,military,incorporated,timezone,ranking,zips,id
	);`
	_, err = db.Exec(create_table)
	if err != nil {
		return err
	}

	delete_contents := `DELETE FROM cities;`
	_, err = db.Exec(delete_contents)
	if err != nil {
		return err
	}

	insert_statement := `INSERT INTO cities (city,city_ascii,state_id,state_name,county_fips,county_name,lat,lng,population,density,source,military,incorporated,timezone,ranking,zips,id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := db.Prepare(insert_statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert a row
	for _, r := range records {
		_, err := stmt.Exec(r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8],
			r[9], r[10], r[11], r[12], r[13], r[14], r[15], r[16])
		fmt.Printf("Recorded city: %s, %s\n", r[0], r[3])
		if err != nil {
			return err
		}
	}

	return nil
}
