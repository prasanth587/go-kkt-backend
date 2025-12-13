# Quick Guide: Test SMS Locally

## Step 1: Set Twilio Environment Variables

Open your terminal and set these environment variables:

```bash
export TWILIO_ACCOUNT_SID="your_account_sid_here"
export TWILIO_AUTH_TOKEN="your_auth_token_here"
export TWILIO_PHONE_NUMBER="+1234567890"  # Your Twilio phone number with country code
export FRONTEND_URL="http://localhost:3000"  # For local testing
```

**OR** create a `.env.local` file in `go-kkt-backend/` directory:

```bash
cd go-kkt-backend
cat > .env.local << EOF
export TWILIO_ACCOUNT_SID="your_account_sid_here"
export TWILIO_AUTH_TOKEN="your_auth_token_here"
export TWILIO_PHONE_NUMBER="+1234567890"
export FRONTEND_URL="http://localhost:3000"
EOF
```

Then source it:
```bash
source .env.local
```

## Step 2: Test SMS Directly (Quick Test)

Test SMS without running the full server:

```bash
cd go-kkt-backend
go run test_sms.go
```

**Important:** 
- Update `test_sms.go` line 28 with your verified phone number
- Twilio trial accounts can only send to verified numbers
- Verify your number at: https://console.twilio.com/us1/develop/phone-numbers/manage/verified

## Step 3: Run Backend Locally

```bash
cd go-kkt-backend
./run-dev.sh
```

The server will start on `http://localhost:9005`

## Step 4: Test Payment SMS

1. Open your frontend: `http://localhost:3000`
2. Go to Payments → Base Amount Alerts
3. Click "Pay" on a trip sheet
4. Check backend console logs for SMS debug messages:
   - `SMS DEBUG: Attempting to send SMS to: +91XXXXXXXXXX`
   - `SUCCESS: Payment SMS sent successfully to vendor: +91XXXXXXXXXX`
   - OR `ERROR: Failed to send payment SMS...`

## Step 5: Check Logs

Look for these messages in the backend console:

**If SMS is working:**
```
SMS DEBUG: Attempting to send SMS to: +91XXXXXXXXXX
SMS DEBUG: Message: Dear ABC Transport,...
SUCCESS: Payment SMS sent successfully to vendor: +91XXXXXXXXXX
```

**If SMS failed:**
```
ERROR: Failed to send payment SMS to vendor: +91XXXXXXXXXX: [error message]
```

**If vendor contact info missing:**
```
SMS DEBUG: Vendor mobile number not found in contact info, skipping SMS. Vendor ID: X
```

**If Twilio not configured:**
```
ERROR: SMS not configured. Missing: AccountSID=true, AuthToken=true, FromNumber=true
```

## Troubleshooting

### Issue: "SMS not configured"
**Solution:** Make sure all 3 environment variables are set:
- `TWILIO_ACCOUNT_SID`
- `TWILIO_AUTH_TOKEN`
- `TWILIO_PHONE_NUMBER`

### Issue: "Vendor mobile number not found"
**Solution:** 
1. Go to Vendors → Edit Vendor
2. Add Contact Information
3. Fill `contact_number1` or `contact_number2` with format: `+91XXXXXXXXXX`

### Issue: "Twilio API returned error: 400"
**Solution:**
- Check if phone number is verified in Twilio Console
- Check if Twilio account has credits
- Verify phone number format includes country code

### Issue: Payment shows "Paid" but no SMS
**Solution:**
- Check backend logs for 400 errors (update might be failing)
- Check if vendor has contact info with phone number
- Check if Twilio credentials are set correctly

## Quick Test Command

```bash
# Set variables
export TWILIO_ACCOUNT_SID="ACxxxxx"
export TWILIO_AUTH_TOKEN="your_token"
export TWILIO_PHONE_NUMBER="+1234567890"

# Test SMS
cd go-kkt-backend
go run test_sms.go
```

