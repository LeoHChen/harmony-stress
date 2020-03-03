package utils

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// RandomStringSliceItem - fetches a random string from a given string slice
func RandomStringSliceItem(r *rand.Rand, items []string) string {
	return items[r.Intn(len(items))]
}

// RandomItemFromMap - select a random item from a map
func RandomItemFromMap(itemMap map[string]string) (string, string) {
	var keys []string

	for key := range itemMap {
		keys = append(keys, key)
	}

	randKey := RandomItemFromSlice(keys)
	randItem := itemMap[randKey]

	return randKey, randItem
}

// RandomItemFromSlice - select a random item from a slice
func RandomItemFromSlice(items []string) string {
	rand.Seed(time.Now().Unix())
	item := items[rand.Intn(len(items))]

	return item
}

// CreateDirectory - creates a directory (mkdir -p) if it already doesn't exist
func CreateDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadFileToString - check if a file exists, proceed to read it to memory if it does
func ReadFileToString(filePath string) (string, error) {
	if FileExists(filePath) {
		data, err := ioutil.ReadFile(filePath)

		if err != nil {
			return "", err
		}

		return string(data), nil
	} else {
		return "", fmt.Errorf("file %s doesn't exist - make sure it exists or that you've specified the correct path for it", filePath)
	}
}

// FileExists - checks if a given file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ParseYaml - parses yaml into a specific type
func ParseYaml(path string, entity interface{}) error {
	yamlData, err := ReadFileToString(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal([]byte(yamlData), entity); err != nil {
		return err
	}

	return nil
}
