package service

import (
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestReadFile(t *testing.T) {
	filename := "./testdata/sample.csv"
	var wg sync.WaitGroup
	svc := NewService(filename, &wg)
	if err := svc.ReadFile(); err != nil {
		t.Error(err)
	}
	svc.Output()
}
