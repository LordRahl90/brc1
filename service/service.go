package service

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Station struct {
	City        string
	Measurement float32
	Sum         float32
	Mean        float32
	Min         float32
	Max         float32
	Count       int32
}

// Service is the service interface.
type Service struct {
	fileName string
	stations map[string]*Station
	output   strings.Builder
}

func NewService(fileName string) *Service {
	return &Service{
		fileName: fileName,
		stations: make(map[string]*Station),
	}
}

// ReadFile reads the file and send the data to the channel.
func (s *Service) ReadFile() error {
	file, err := os.Open(s.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		station, err := newStation(line)
		if err != nil {
			slog.Error("Error while reading file", err, slog.String("line", line))
			continue
		}
		// send this to the channel
		s.Compute(station)
	}
	println("Scan completed!")
	return nil
}

// Compute computes the station data.
func (s *Service) Compute(station *Station) {
	// get the current station
	currentStation, ok := s.stations[station.City]
	if !ok {
		station.Count = 1
		station.Mean = station.Measurement
		station.Min = station.Measurement
		station.Max = station.Measurement
		station.Sum = station.Measurement
		s.stations[station.City] = station
		return
	}
	currentStation.Count++
	currentStation.Sum += station.Measurement
	currentStation.Mean = currentStation.Sum / float32(currentStation.Count)
	if currentStation.Min > station.Measurement {
		currentStation.Min = station.Measurement
	}
	if currentStation.Max < station.Measurement {
		currentStation.Max = station.Measurement
	}
	s.stations[station.City] = currentStation
}

func (s *Service) Output() {
	s.output.WriteString("{")
	for _, v := range s.stations {
		out := fmt.Sprintf("%s=%.2f/%.2f/%.2f\n",
			v.City, v.Min, v.Mean, v.Max)
		s.output.WriteString(out)
	}
	s.output.WriteString("}")
	println(s.output.String())
}

func newStation(line string) (*Station, error) {
	val := strings.Split(line, ";")
	if len(val) != 2 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}
	temp, err := strconv.ParseFloat(val[1], 32)
	if err != nil {
		return nil, err
	}
	return &Station{
		City:        val[0],
		Measurement: float32(temp),
	}, nil
}
