package application

import (
	"os"
)

func isFile(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}
