package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func GetFilenameWithoutExtension(path string) string {
	split := strings.Split(path, "/")
	filename := split[len(split)-1]
	return filename[:len(filename)-4]
}

func RemoveNilStringInArray(array []string) []string {
	for index, value := range array {
		if value == "" {
			array = append(array[:index], array[index+1:]...)
		}
	}

	return array
}

func ConvertArrayStringToBool(array []string) []bool {
	boolArray := []bool{}

	for _, value := range array {
		if value == "true" {
			boolArray = append(boolArray, true)
		} else {
			boolArray = append(boolArray, false)
		}
	}

	return boolArray
}

func ReplaceUnderScoreWithSpace(text string) string {
	return strings.ReplaceAll(text, "_", " ")
}
