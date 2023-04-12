package helpers

import (
	"os"
	"strconv"
)

func FetchIndex() (int, error) {
	index := "-1"

	if os.Getenv("CF_INSTANCE_INDEX") != "" {
		index = os.Getenv("CF_INSTANCE_INDEX")
	} else if os.Getenv("INSTANCE_INDEX") != "" {
		index = os.Getenv("INSTANCE_INDEX")
	}
	return strconv.Atoi(index)
}
