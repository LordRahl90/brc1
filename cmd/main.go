package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/LordRahl90/brc1/service"
)

func main() {
	start := time.Now()
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	f, err = os.Create("mem.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal(err)
	}

	// read file content
	// get the line
	// split the line to city and temperature
	// check if the city exists in the map
	// if it doesn't exist, add it to the map and set the count to 1
	// if it exists, increment the count. and check if the latest temperature is less that the current one.
	// if it is less, update, otherwise, proceed.
	filename := "./data/weather_stations.csv"
	println(filename)
	svc := service.NewService(filename)
	if err := svc.ReadFile(); err != nil {
		log.Fatal(err)
	}
	svc.Output()
	dur := time.Since(start)
	fmt.Printf("All done Within %.2f seconds\n", dur.Seconds())
}
