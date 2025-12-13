# Debug: SMS Not Sending When Clicking Pay

## Quick Checklist

### 1. Check Railway Backend Logs
After clicking "Pay", look for these log messages:

**✅ If SMS function is being called:**
```
Payment update successful, checking if SMS should be sent...
=== SMS DEBUG: Checking payment SMS ===
Trip Sheet ID: X
Old paid date: [old value]
New paid date: BASE:2025-12-13
```

**❌ If payment update failed:**
```
ERROR: UpdateTripSheetHeader [error message]
```
→ SMS won't be sent if update fails

**❌ If no payment change detected:**
```
SMS DEBUG: No payment change detected, skipping SMS
```
→ Old and new paid dates are the same

**❌ If vendor contact missing:**
```
SMS DEBUG: Vendor mobile number not found in contact info, skipping SMS. Vendor ID: X
```
→ Vendor needs contact info with phone number

**❌ If Twilio not configured:**
```
ERROR: SMS not configured. Missing: AccountSID=true, AuthToken=true, FromNumber=true
```
→ Set Twilio credentials in Railway

**✅ If SMS is being sent:**
```
SMS DEBUG: Attempting to send SMS to: +91XXXXXXXXXX
SUCCESS: Payment SMS sent successfully to vendor: +91XXXXXXXXXX
```

## Common Issues & Fixes

### Issue 1: "No payment change detected"
**Cause:** Old and new `vendor_paid_date` are the same
**Fix:** Make sure you're actually changing the payment status when clicking "Pay"

### Issue 2: "Vendor mobile number not found"
**Cause:** Vendor has no contact info in database
**Fix:** 
1. Go to Vendors → Edit Vendor
2. Add Contact Information
3. Fill `contact_number1` or `contact_number2` (format: `+91XXXXXXXXXX`)

### Issue 3: "SMS not configured"
**Cause:** Twilio credentials not set in Railway
**Fix:** Set these in Railway:
- `TWILIO_ACCOUNT_SID`
- `TWILIO_AUTH_TOKEN`
- `TWILIO_PHONE_NUMBER`

### Issue 4: Payment update failing (400 error)
**Cause:** Validation error (missing LR number, loading points, etc.)
**Fix:** Check Railway logs for exact validation error

## How to Debug Step by Step

1. **Click "Pay" button**
2. **Check Railway logs immediately** - look for:
   - `Payment update successful` → Good, update worked
   - `=== SMS DEBUG: Checking payment SMS ===` → Good, SMS function called
   - Then check what happens next...

3. **If you see "No payment change detected":**
   - The old and new paid dates are the same
   - This means the payment wasn't actually updated
   - Check if the update request is being sent correctly

4. **If you see "Vendor mobile number not found":**
   - Go to database or UI and add vendor contact info
   - Make sure phone number has country code (+91 for India)

5. **If you see "SMS not configured":**
   - Set Twilio credentials in Railway environment variables

6. **If you see "ERROR: Failed to send payment SMS":**
   - Check Twilio account balance
   - Verify phone number is verified (for trial accounts)
   - Check phone number format

## Test Locally

Run backend locally and check console output:

```bash
cd go-kkt-backend
source setup-sms-local.sh
./run-dev.sh
```

Then click "Pay" and watch the console for SMS debug messages.

