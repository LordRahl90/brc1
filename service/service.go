package service

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
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

	wg *sync.WaitGroup
}

func NewService(fileName string, wg *sync.WaitGroup) *Service {
	return &Service{
		fileName: fileName,
		stations: make(map[string]*Station),
		wg:       wg,
	}
}

// ReadFile reads the file and send the data to the channel.
func (s *Service) ReadFile() error {
	var wg sync.WaitGroup
	file, err := os.Open(s.fileName)
	if err != nil {
		return err
	}
	fs, err := file.Stat()
	if err != nil {
		return err
	}
	data, err := syscall.Mmap(int(file.Fd()), 0, int(fs.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	defer func() {
		if err := syscall.Munmap(data); err != nil {
			log.Fatal(err)
		}
	}()
	// var line strings.Builder
	nChunks := 10000
	chunkSize := len(data) / nChunks
	if chunkSize == 0 {
		slog.Error("chunksize is less than 0, will default to 1000")
		chunkSize = 1000
	}
	chunks := make([]int, 0, nChunks)
	offset := 0
	for {
		offset += chunkSize
		if offset >= len(data) {
			break
		}

		// find the first newline position between the offset and the end of the data
		nlPos := bytes.IndexByte(data[offset:], '\n')
		if nlPos == -1 {
			chunks = append(chunks, len(data))
		} else {
			offset += nlPos + 1
			chunks = append(chunks, offset)
		}
	}
	start := 0
	results := make([]map[string]*Station, len(chunks))
	for i, chunk := range chunks {
		wg.Add(1)
		go func() {
			// println("Chunk: ", chunk, "Size: ", len(data))
			chunkData := data[start:chunk]
			result := s.process(&wg, chunkData)
			if result != nil {
				results[i] = result
			}
		}()
		start = chunk
	}

	wg.Wait()
	println("All chunks read successfully")

	for _, result := range results {
		for _, v := range result {
			s.Compute(v)
		}
	}

	return nil
}

func (s *Service) process(wg *sync.WaitGroup, data []byte) map[string]*Station {
	defer func() {
		wg.Done()
	}()
	if data == nil {
		return nil
	}
	result := make(map[string]*Station)
	var (
		line    strings.Builder
		content []byte
	)

	for _, v := range data {
		if v != '\n' {
			content = append(content, v)
			continue
		}
		station, err := newStation(content)
		if err != nil {
			slog.Error("an error occurred", "error", err, "line", line.String())
			content = nil
			continue

		}
		curr, ok := result[station.City]
		if !ok {
			curr = &Station{
				City:        station.City,
				Measurement: station.Measurement,
				Sum:         station.Measurement,
				Mean:        station.Measurement,
				Min:         station.Measurement,
				Max:         station.Measurement,
				Count:       1,
			}
			result[station.City] = curr
		} else {
			curr.Count++
			curr.Sum += station.Measurement
			curr.Mean = curr.Sum / float32(curr.Count)
			curr.Min = min(curr.Min, station.Measurement)
			curr.Max = max(station.Measurement)
			result[station.City] = curr
		}
		content = nil
	}
	return result
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
	currentStation.Count += station.Count
	currentStation.Sum += station.Sum
	currentStation.Mean = currentStation.Sum / float32(currentStation.Count)
	currentStation.Min = min(currentStation.Min, station.Min)
	currentStation.Max = max(currentStation.Max, station.Max)
	s.stations[station.City] = currentStation
}

func (s *Service) Output() {
	s.output.WriteString("{")
	stationNames := make([]string, 0, len(s.stations))

	for _, v := range s.stations {
		stationNames = append(stationNames, v.City)
	}
	sort.Strings(stationNames)

	for _, st := range stationNames {
		v := s.stations[st]
		out := fmt.Sprintf("%s=%.1f/%.1f/%.1f/%d\n",
			v.City, v.Min, v.Mean, v.Max, v.Count)
		s.output.WriteString(out)
	}

	s.output.WriteString("}")
	fmt.Print(s.output.String())
}

func newStation(line []byte) (*Station, error) {
	offset := bytes.IndexByte(line, ';')
	if offset == -1 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}

	temp, err := strconv.ParseFloat(string(line[offset+1:]), 32)
	if err != nil {
		return nil, err
	}
	return &Station{
		City:        string(line[:offset]),
		Measurement: float32(temp),
	}, nil
}
