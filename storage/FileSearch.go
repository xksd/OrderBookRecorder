package storage

import (
	"os"
)

// Check if settings files exist
// Search for alternative
func SearchAlternativeFilePath(filePath *string, numberOfTries int) bool {
	for i := 0; i < numberOfTries; i++ {
		if !FileExists(*filePath) {
			*filePath = "." + *filePath
		}
	}
	return FileExists(*filePath)

}

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
