package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// GetSMSConfig retrieves SMS configuration from environment variables
func GetSMSConfig() *SMSConfig {
	return &SMSConfig{
		AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		FromNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
	}
}

// SendSMS sends an SMS message using Twilio API
func SendSMS(mobileNumber, message string) error {
	config := GetSMSConfig()
	
	fmt.Printf("=== SMS UTILITY DEBUG ===\n")
	fmt.Printf("Mobile Number: %s\n", mobileNumber)
	fmt.Printf("Account SID: %s\n", config.AccountSID)
	fmt.Printf("Auth Token: %s (length: %d)\n", maskToken(config.AuthToken), len(config.AuthToken))
	fmt.Printf("From Number: %s\n", config.FromNumber)
	
	// If SMS is not configured, just log and return (don't fail)
	if config.AccountSID == "" || config.AuthToken == "" || config.FromNumber == "" {
		fmt.Printf("ERROR: SMS not configured. Missing: AccountSID=%v, AuthToken=%v, FromNumber=%v\n", 
			config.AccountSID == "", config.AuthToken == "", config.FromNumber == "")
		return nil
	}

	// Format mobile number (ensure it starts with country code, e.g., +91 for India)
	toNumber := formatPhoneNumber(mobileNumber)
	fmt.Printf("Formatted number: %s\n", toNumber)

	// Twilio API endpoint
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", config.AccountSID)
	fmt.Printf("Twilio API URL: %s\n", apiURL)

	// Prepare form data
	data := url.Values{}
	data.Set("From", config.FromNumber)
	data.Set("To", toNumber)
	data.Set("Body", message)

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	// Set Basic Authentication header (Twilio uses Account SID and Auth Token)
	auth := base64.StdEncoding.EncodeToString([]byte(config.AccountSID + ":" + config.AuthToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %w", err)
	}

	// Check if request was successful (Twilio returns 201 for created)
	fmt.Printf("Twilio Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Twilio Response Body: %s\n", string(body))
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Twilio API returned error: %s (status: %d)", string(body), resp.StatusCode)
	}

	// Parse response to get message SID
	var twilioResponse map[string]interface{}
	if err := json.Unmarshal(body, &twilioResponse); err == nil {
		if sid, ok := twilioResponse["sid"].(string); ok {
			fmt.Printf("SUCCESS: SMS sent successfully to %s. Message SID: %s\n", toNumber, sid)
		}
	}

	fmt.Printf("=== SMS UTILITY DEBUG END ===\n")
	return nil
}

// maskToken masks the auth token for logging (shows first 4 and last 4 chars)
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

// formatPhoneNumber ensures phone number has country code
// If number doesn't start with +, assumes it's an Indian number and adds +91
func formatPhoneNumber(number string) string {
	// Remove any spaces or dashes
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")
	
	// If already has +, return as is
	if strings.HasPrefix(number, "+") {
		return number
	}
	
	// If starts with 0, remove it
	if strings.HasPrefix(number, "0") {
		number = number[1:]
	}
	
	// If doesn't start with country code, assume India (+91)
	if !strings.HasPrefix(number, "91") && len(number) == 10 {
		return "+91" + number
	}
	
	// If already has 91 prefix, add +
	if strings.HasPrefix(number, "91") {
		return "+" + number
	}
	
	// Default: add +91 for Indian numbers
	return "+91" + number
}

// FormatIndianCurrency formats amount in Indian currency format
func FormatIndianCurrency(amount float64) string {
	return fmt.Sprintf("₹ %s", formatNumber(amount))
}

func formatNumber(num float64) string {
	str := fmt.Sprintf("%.2f", num)
	// Add comma separators for Indian number format
	integerPart := ""
	decimalPart := ""
	
	if idx := len(str) - 3; idx > 0 {
		integerPart = str[:idx]
		decimalPart = str[idx:]
	} else {
		integerPart = str
	}

	// Reverse the integer part to add commas
	runes := []rune(integerPart)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Add commas every 2 digits (Indian numbering system)
	result := ""
	for i, r := range runes {
		if i > 0 && i%2 == 0 {
			result += ","
		}
		result += string(r)
	}

	// Reverse back
	runes = []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes) + decimalPart
}

