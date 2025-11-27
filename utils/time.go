package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	//TFhhmm time format hour:minute
	TFhhmm = "15:04"
)

// DateToHHMMFormat convert date to hh:mm
func DateToHHMMFormat(time time.Time) string {
	return time.Format(TFhhmm)
}

// CheckTimeFormatHHMM function is use to check given time is valid or not
func CheckTimeFormatHHMM(timeStr string) bool {
	_, err := time.ParseInLocation(TFhhmm, timeStr, time.UTC)
	if err != nil {
		return false
	}
	return true
}

// HHMMToTime function is use to convert hhmm to date format
func HHMMToTime(date time.Time, timeStr string) time.Time {

	strArr := strings.Split(timeStr, ":")

	hour, _ := strconv.Atoi(strArr[0])
	min, _ := strconv.Atoi(strArr[1])
	// t := time.Now()

	t1 := time.Date(date.Year(), date.Month(), date.Day(), hour, min, 00, 0, time.UTC)

	return t1.UTC()
}

func TimeLoc() *time.Location {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil
	}
	return loc
}

func HandleZeroPadding(dateInput string) (string, error) {
	t, err := time.Parse("2-1-2006", dateInput)
	if err != nil {
		fmt.Println("Invalid date:", err)
		return "", errors.New("date format issue")
	}
	dateFormat := t.Format("2006-01-02")
	fmt.Println(dateFormat) // Output: 2025-07-05
	return dateFormat, nil
}
func HandleZeroPaddingWithTime(dateInput string) (string, error) {
	t, err := time.Parse("2-1-2006 15:04", dateInput)
	if err != nil {
		fmt.Println("Invalid date:", err)
		return "", errors.New("date format issue")
	}
	dateFormat := t.Format("2006-01-02 15:04")
	fmt.Println(dateFormat) // Output: 2025-07-05
	return dateFormat, nil
}
