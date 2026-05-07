package loader

import (
	"os"
	"strconv"
)

func SaveProgress(i int) error {
	if err := os.WriteFile("progress.txt", []byte(strconv.Itoa(i)), 0644); err != nil {
		return err
	}
	return nil
}
