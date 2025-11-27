package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/prabha303-vi/log-util/log"
)

// MinusIDs get array of IDs
func MinusIDs(from, to []int64) []int64 {
	res := []int64{}

	for _, a := range from {
		find := false
		for _, b := range to {
			if a == b {
				find = true
				break
			}
		}

		if !find {
			res = append(res, a)
		}
	}
	return res
}

// GetClientIP is an utility function for getting actual IP of end user
func GetClientIP(r *http.Request) string {
	// Header X-Forwarded-For
	hdrForwardedFor := http.CanonicalHeaderKey("X-Forwarded-For")
	if fwdFor := strings.TrimSpace(r.Header.Get(hdrForwardedFor)); fwdFor != "" {
		index := strings.Index(fwdFor, ",")
		if index == -1 {
			return fwdFor
		}
		return fwdFor[:index]
	}

	// Header X-Real-Ip
	hdrRealIP := http.CanonicalHeaderKey("X-Real-Ip")
	if realIP := strings.TrimSpace(r.Header.Get(hdrRealIP)); realIP != "" {
		return realIP
	}

	return "10.82.33.161"
}

func ValidateMobileNumber(mobileNo string) error {
	// Regular expression to match exactly 10 digits
	re := regexp.MustCompile(`^\d{10}$`)
	if !re.MatchString(mobileNo) {
		return fmt.Errorf("mobile number must be a 10-digit number")
	}
	return nil
}

func ValidateDriverExperience(experience float64) error {
	if experience < 0 || experience > 80 {
		return errors.New("driver experience is not valide")
	}
	return nil
}

func ValidateLicenseExpiryDate(dateStr string) error {
	// Parse the date string into a time.Time object
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	fmt.Println("date - ", date)

	// Check if the date is not in the past or today
	now := time.Now().In(TimeLoc())
	if date.Before(now) || date.Equal(now) {
		return errors.New("license expiry date must be a future date")
	}

	return nil
}

func ValidateDateStr(dateStr string) error {
	// Parse the date string into a time.Time object
	layout := "2006-01-02"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	return nil
}

func CheckSizeImage(file *multipart.FileHeader, limit int64, ul *log.Logger) bool {
	size := file.Size / 1024
	ul.Info("Employee image size kb - ", size)
	return size <= limit
}

func CheckFileExists(filePath string) bool {
	_, err := os.Open(filePath)
	return err == nil
}

func MustMarshal(v interface{}, logMsd string) string {
	b, _err := json.Marshal(v)
	if _err != nil {
		fmt.Println("logMsd ", logMsd, " : error: ", _err)
	}
	return string(b)
}
