package debug

import (
	"os"
	"testing"
)

type testStruct struct {
	Hello string `json:"hello"`
	World string `json:"world"`
}

var testData = testStruct{
	Hello: "hello",
	World: "world",
}

var testDataFile = ".debug_test-data.json"

func TestWriteDataFile(t *testing.T) {
	if _, err := os.Stat(testDataFile); err == nil {
		os.Remove(testDataFile)
	}
	WriteDataFile(testDataFile, testData)
	if _, err := os.Stat(testDataFile); err != nil {
		t.Error("test did not write a data file")
	}
}

func TestLoadDataFileSuccess(t *testing.T) {
	want := testData
	got := LoadDataFile[testStruct](testDataFile)
	if got.Hello != want.Hello || got.World != want.World {
		t.Errorf("got %v but want %v", got, want)
	}
}

func TestLoadDataFileFail(t *testing.T) {
	got := LoadDataFile[testStruct]("non-exist")
	if got != nil {
		t.Error("loaded non-exist file")
	}
}
