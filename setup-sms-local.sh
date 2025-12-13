#!/bin/bash

# Quick setup script for local SMS testing
# This sets up your Twilio credentials for local testing
# IMPORTANT: Set your actual Twilio credentials as environment variables before running this script

echo "Setting up Twilio SMS for local testing..."
echo ""

# Get Twilio credentials from environment variables or use placeholders
# Set these before running: export TWILIO_ACCOUNT_SID="your_sid"
export TWILIO_ACCOUNT_SID="${TWILIO_ACCOUNT_SID:-your_account_sid_here}"
export TWILIO_AUTH_TOKEN="${TWILIO_AUTH_TOKEN:-your_auth_token_here}"
export TWILIO_PHONE_NUMBER="${TWILIO_PHONE_NUMBER:-your_phone_number_here}"
export FRONTEND_URL="${FRONTEND_URL:-http://localhost:3000}"

if [ "$TWILIO_ACCOUNT_SID" = "your_account_sid_here" ]; then
    echo "⚠️  WARNING: Twilio credentials not set!"
    echo "   Please set these environment variables:"
    echo "   export TWILIO_ACCOUNT_SID=\"your_account_sid\""
    echo "   export TWILIO_AUTH_TOKEN=\"your_auth_token\""
    echo "   export TWILIO_PHONE_NUMBER=\"your_phone_number\""
    echo ""
fi

echo "✅ Twilio credentials configured:"
echo "   Account SID: ${TWILIO_ACCOUNT_SID:0:10}... (hidden)"
echo "   Auth Token: ${TWILIO_AUTH_TOKEN:0:10}... (hidden)"
echo "   Phone Number: $TWILIO_PHONE_NUMBER"
echo "   Frontend URL: $FRONTEND_URL"
echo ""
echo "⚠️  IMPORTANT: You're on a Twilio Trial account"
echo "   - You can only send SMS to VERIFIED phone numbers"
echo "   - Verify numbers at: https://console.twilio.com/us1/develop/phone-numbers/manage/verified"
echo ""
echo "To test SMS, run:"
echo "  1. Quick test: go run test_sms.go"
echo "  2. Full server: ./run-dev.sh"
echo ""

