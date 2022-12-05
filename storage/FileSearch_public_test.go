package storage_test

import (
	"testing"

	"github.com/xksd/OrderBookRecorder/storage"
)

func TestSeachAlternativeFilePath(t *testing.T) {
	testPath := "./main.go"
	if !storage.SearchAlternativeFilePath(&testPath, 3) {
		t.Error("File exists, but NOT found!")
		t.Error(testPath)
	}
}
