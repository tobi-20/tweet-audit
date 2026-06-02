package loader

import (
	"os"
	"strconv"
)

func SaveProgress(i int, filepath string) error {
	if err := os.WriteFile(filepath, []byte(strconv.Itoa(i)), 0644); err != nil {
		return err
	}
	return nil
}
