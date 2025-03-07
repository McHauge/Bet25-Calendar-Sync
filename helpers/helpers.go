package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/s00500/env_logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ##########################
// #### Helper Functions ####
// ##########################

func ToUpperCase(str string) string {
	caser := cases.Upper(language.English)
	return caser.String(str)
}

func ToTitleCase(str string) string {
	caser := cases.Title(language.English)
	return caser.String(str)
}

func ToLowerCase(str string) string {
	caser := cases.Lower(language.English)
	return caser.String(str)
}

func ParseInt(value string) int {
	resp, err := strconv.ParseInt(value, 10, 32)
	log.Should(err)
	return int(resp)
}

func IntToStr(val int, padding int) (H string) {
	switch padding {
	case 1:
		H = fmt.Sprintf("%01d", val)
	case 2:
		H = fmt.Sprintf("%02d", val)
	case 3:
		H = fmt.Sprintf("%03d", val)
	case 4:
		H = fmt.Sprintf("%04d", val)
	case 5:
		H = fmt.Sprintf("%05d", val)
	case 6:
		H = fmt.Sprintf("%06d", val)
	default:
		H = fmt.Sprintf("%d", val)
	}

	return H
}

// Create a folder on disk if it dows not already exsist
func CreateFolder(path string, folder string) string {
	folderPath := strings.Join([]string{path, folder}, "/")

	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			log.Should(err)
		}
	}
	return folder
}

// Delete all folders and files in a folder
func DeleteFolder(path string, folder string) (err error) {
	folderPath := strings.Join([]string{path, folder}, "/")

	err = os.RemoveAll(folderPath)
	if err != nil {
		return err
	}

	return nil
}

// Delete a saved file
func DeleteFile(path string, fileName string) (err error) {
	filePath := strings.Join([]string{path, fileName}, "/")

	err = os.Remove(filePath)
	if err == nil {
		return err
	}

	return nil
}

// updates the Json
func UpdateJson(data interface{}, path string, fileName string) error {
	filePath := strings.Join([]string{path, fileName}, "/")

	file, err := json.Marshal(data)
	if !log.Should(err) {
		err = os.WriteFile(filePath, file, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// reads the Json and creates a new if not found
func ReadJson(data interface{}, path string, fileName string) error {
	filePath := strings.Join([]string{path, fileName}, "/")

	file, err := os.ReadFile(filePath) // Read File
	if err != nil {
		// Write a new file
		file, err := json.Marshal(data)
		log.Debug(err)
		if err == nil {
			err = os.WriteFile(filePath, file, 0755)
		}
		return err
	}

	err = json.Unmarshal(file, data)
	if err != nil {
		return err
	}
	return nil
}

// Write a text file
func WriteFile(filePath string, data string) error {
	return os.WriteFile(filePath, []byte(data), 0666)
}

// Read a text file
func ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

