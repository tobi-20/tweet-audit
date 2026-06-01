package loader

import (
	"os"
	"strconv"
)

func LoadProgress(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	i, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return i, nil
}
