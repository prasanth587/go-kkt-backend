# Testing SMS Locally

## ✅ Step 1: Verify a Phone Number in Twilio (IMPORTANT!)

Since you're on a **trial account**, you can only send SMS to **verified phone numbers**.

1. Go to: https://console.twilio.com/us1/develop/phone-numbers/manage/verified
2. Click **"Add a new number"**
3. Enter the phone number you want to test with (e.g., your own number or a vendor's number)
4. Select **"SMS"** as the verification method
5. Click **"Send Verification Code"**
6. Enter the code you receive via SMS
7. Click **"Verify"**

**Note:** You can verify multiple numbers for testing.

## ✅ Step 2: Make Sure Backend is Running

The backend should be running with Twilio credentials. Check the terminal output for:
```
Starting backend server...
Server will run on: http://localhost:9005
```

## ✅ Step 3: Test SMS by Making a Payment

1. **Open your frontend:**
   - Go to: http://localhost:3000 (or your frontend URL)
   - Login to the application

2. **Navigate to Payments:**
   - Go to **Settings** → **Payments**
   - Click on **"Base Amount Alerts"** or **"Advance Amount Alerts"** tab

3. **Make a Payment:**
   - Find a trip sheet with pending payment
   - Click the **"Pay"** button
   - Confirm the payment

4. **Check Backend Logs:**
   - Look for: `SMS sent successfully to +91XXXXXXXXXX. Message SID: SM...`
   - If there's an error, you'll see: `ERROR: Failed to send payment SMS...`

5. **Check Phone:**
   - The verified phone number should receive an SMS like:
   ```
   Dear [Vendor Name], Payment of ₹ 10,000.00 has been processed for Trip Sheet 000001/2025-2026. Payment Type: Base Amount. Thank you! - KK Transport
   ```

## 🔍 Troubleshooting

### SMS not sending?
- **Check if phone number is verified** in Twilio Console
- **Check backend logs** for error messages
- **Verify vendor has mobile number** in the database
- **Check Twilio Console → Monitor → Logs** for API errors

### "SMS not configured" message?
- Make sure `run-dev.sh` has the Twilio credentials uncommented
- Restart the backend after adding credentials

### Payment works but no SMS?
- Check backend logs for SMS errors
- Verify the phone number format (should be +91XXXXXXXXXX for India)
- Check Twilio account has credits/balance

## 📱 Testing with Your Own Number

1. Verify your own phone number in Twilio Console
2. Update a vendor in your database to use your phone number
3. Make a payment for that vendor
4. You should receive SMS on your phone

## 🎯 Expected Backend Log Output

**Success:**
```
SMS sent successfully to +919876543210. Message SID: SM1234567890abcdef1234567890abcdef
Payment SMS sent successfully to vendor: +919876543210
```

**Error (if phone not verified):**
```
ERROR: Failed to send payment SMS to vendor: +919876543210: Twilio API returned error: ... (status: 400)
```

