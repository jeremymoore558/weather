module weather

go 1.26.1

replace weather/noa => ./noa

require (
	weather/dbutils v0.0.0-00010101000000-000000000000
	weather/noa v0.0.0-00010101000000-000000000000
)

require github.com/mattn/go-sqlite3 v1.14.44 // indirect

replace weather/dbutils => ./dbutils
