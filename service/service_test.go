package service

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestReadFile(t *testing.T) {
	filename := "./testdata/sample.csv"
	svc := NewService(filename)
	if err := svc.ReadFile(); err != nil {
		t.Error(err)
	}
	svc.Output()
}
