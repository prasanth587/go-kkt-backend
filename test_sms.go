package main

import (
	"fmt"
	"go-transport-hub/utils"
	"os"
)

func main() {
	// Set environment variables for testing
	// IMPORTANT: Set these as environment variables before running, or update with your actual credentials
	// Example: export TWILIO_ACCOUNT_SID="your_account_sid" before running this test
	if os.Getenv("TWILIO_ACCOUNT_SID") == "" {
		fmt.Println("ERROR: TWILIO_ACCOUNT_SID environment variable not set")
		fmt.Println("Please set TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_PHONE_NUMBER before running")
		os.Exit(1)
	}
	// Use environment variables - do not hardcode credentials

	// Test phone number (replace with your verified number)
	// IMPORTANT: This must be a verified number in Twilio Console!
	testNumber := "+919362196294" // Using your verified number from Twilio
	testMessage := "Test SMS from KK Transport - Payment notification system"

	fmt.Printf("Testing SMS to: %s\n", testNumber)
	fmt.Printf("Message: %s\n\n", testMessage)

	err := utils.SendSMS(testNumber, testMessage)
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
	} else {

		fmt.Printf("✅ SMS sent successfully!\n")
	}
}
