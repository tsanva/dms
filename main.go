package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Timer int    `json:"timer"`
	Dir   string `json:"dir"`
}

func main() {

	// Read the JSON file
	filePath := "config.json" // Replace with the actual file path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON into Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	usbDir := "dummy/dev/bus/usb"
	var initialCount int

	// Get initial count of files in usbDir
	initialCount = countFiles(usbDir)

	fmt.Println("dms active")

	// Loop every 5 seconds to check if the number of files in usbDir has changed
	for {
		time.Sleep(time.Duration(config.Timer) * time.Second)
		currentCount := countFiles(usbDir)
		if currentCount != initialCount {
			// Delete secretsDir if the number of files in usbDir has changed
			err := os.RemoveAll(config.Dir)
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	fmt.Println("dms triggered")
}

// countFiles counts the number of files in a directory
func countFiles(dir string) int {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		fmt.Println(err)
	}
	return len(files)
}
