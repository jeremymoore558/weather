package noa

import (
	"errors"
	"log"
	"math"
)

func argmin_distance(t []float64, M [][]float64) (int, error) {
	if M == nil {
		err := errors.New("Empty coordinate list provided")
		return 0, err
	}
	if len(t) != len(M[0]) {
		err := errors.New("Dimensions do not match.")
		return 0, err
	}

	var argmin int = 0
	var mindistance float64 = 9999999

	for i, c := range M {
		distance_temp, err := distance(t, c)
		if err != nil {
			return 0, err
		}
		if distance_temp < mindistance {
			mindistance = distance_temp
			argmin = i
		}
	}
	return argmin, nil
}

func distance(a []float64, b []float64) (float64, error) {
	if len(a) != len(b) {
		err := errors.New("Slices not of same length")
		return -1, err
	}

	var squaresum float64 = 0
	for i := range len(a) {
		squaresum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(squaresum), nil
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
