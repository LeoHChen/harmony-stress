package transactions

import (
	"strings"

	"github.com/SebastianJ/harmony-stress/utils"
)

// FetchReceivers - fetch a list of proxies from a specified file
func FetchReceivers(filePath string) (lines []string, err error) {
	data, err := utils.ReadFileToString(filePath)

	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		lines = strings.Split(string(data), "\n")

		if strings.Contains(data, "\n") {
			lines = lines[:len(lines)-1]
		}
	}

	return lines, nil
}
