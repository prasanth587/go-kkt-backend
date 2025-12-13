# SMS Setup Guide - Twilio

This guide explains how to configure SMS notifications for vendor payments using Twilio.

## Overview

When a payment is made (Base Amount or Advance Amount) from the Payments section, an SMS is automatically sent to the vendor's mobile number notifying them about the payment.

## Step 1: Create a Twilio Account

1. Go to https://www.twilio.com/
2. Sign up for a free account (includes trial credits)
3. Verify your email and phone number

## Step 2: Get Your Twilio Credentials

1. Log in to your Twilio Console: https://console.twilio.com/
2. On the dashboard, you'll see:
   - **Account SID** (starts with `AC...`)
   - **Auth Token** (click "View" to reveal it)

3. **Get a Twilio Phone Number:**
   - Go to "Phone Numbers" → "Manage" → "Buy a number"
   - Select a number that supports SMS
   - For India, you may need to verify your account first
   - Copy the phone number (format: `+1234567890`)

## Step 3: Configure Environment Variables

### For Local Development

Add these to your `.env.local` file or `run-dev.sh`:

```bash
# Twilio SMS Configuration
export TWILIO_ACCOUNT_SID="ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export TWILIO_AUTH_TOKEN="your_auth_token_here"
export TWILIO_PHONE_NUMBER="+1234567890"  # Your Twilio phone number
```

### For Production (Railway)

1. Go to your Railway project
2. Navigate to "Variables" tab
3. Add these environment variables:
   - `TWILIO_ACCOUNT_SID` = Your Account SID
   - `TWILIO_AUTH_TOKEN` = Your Auth Token
   - `TWILIO_PHONE_NUMBER` = Your Twilio phone number (with + and country code)

## Step 4: Test the Setup

1. Start your backend server
2. Make a payment from the Payments section
3. Check the backend logs for SMS confirmation
4. Vendor should receive SMS on their mobile number

## Phone Number Format

The system automatically formats Indian phone numbers:
- If number is `9876543210`, it becomes `+919876543210`
- If number already has country code, it's used as-is
- Always include country code (e.g., `+91` for India)

## SMS Message Format

When a payment is made, vendors receive an SMS in this format:

```
Dear [Vendor Name], Payment of ₹ [Amount] has been processed for Trip Sheet [Trip Sheet Number]. Payment Type: [Base Amount/Advance Amount]. Thank you! - KK Transport
```

### Example Messages

**Base Amount Payment:**
```
Dear MG Transport, Payment of ₹ 10,000.00 has been processed for Trip Sheet 000001/2025-2026. Payment Type: Base Amount. Thank you! - KK Transport
```

**Advance Amount Payment:**
```
Dear MG Transport, Payment of ₹ 5,000.00 has been processed for Trip Sheet 000001/2025-2026. Payment Type: Advance Amount. Thank you! - KK Transport
```

## Testing Without SMS Provider

If SMS is not configured (environment variables not set), the system will:
- Log the SMS message to the console
- Continue processing the payment normally
- Not fail the payment operation

This allows you to test the payment functionality without setting up SMS immediately.

## Troubleshooting

1. **SMS not sending:**
   - Check that all environment variables are set correctly
   - Verify the SMS API URL is correct
   - Check backend logs for error messages
   - Ensure vendor has a valid mobile number in the database

2. **SMS sending but not received:**
   - Verify the mobile number format is correct
   - Check SMS provider dashboard for delivery status
   - Ensure SMS provider account has sufficient credits

3. **Payment works but SMS fails:**
   - SMS failures don't affect payment processing
   - Check backend logs for SMS error details
   - Verify SMS provider credentials

## Implementation Details

- SMS is sent asynchronously (non-blocking) using goroutines
- SMS is only sent when payment status changes to "paid" (BASE: or ADVANCE: format)
- SMS is NOT sent for declined payments
- Vendor mobile number is fetched from the `vendors` table


