package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/LordRahl90/brc1/service"
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup

	// read file content
	// get the line
	// split the line to city and temperature
	// check if the city exists in the map
	// if it doesn't exist, add it to the map and set the count to 1
	// if it exists, increment the count. and check if the latest temperature is less that the current one.
	// if it is less, update, otherwise, proceed.

	filename := "./data/measurements_1b.txt"
	// filename := "./data/measurements.txt"
	// filename := "./data/weather_stations.csv"
	// println(filename)

	svc := service.NewService(filename, &wg)
	if err := svc.ReadFile(); err != nil {
		log.Fatal(err)
	}

	// svc.Output()
	dur := time.Since(start)
	fmt.Printf("\nAll done Within %.2f seconds\n", dur.Seconds())
}
