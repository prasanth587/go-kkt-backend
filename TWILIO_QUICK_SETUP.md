# Twilio Quick Setup Guide

## ✅ What You Need to Do:

### 1. Get Your Twilio Credentials

1. **Sign up/Login to Twilio:**
   - Go to: https://console.twilio.com/
   - If new, sign up (free trial available)

2. **Get Account SID and Auth Token:**
   - On the dashboard, you'll see:
     - **Account SID**: `ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
     - **Auth Token**: Click "View" to see it (keep it secret!)

3. **Get a Phone Number:**
   - Go to: Phone Numbers → Manage → Buy a number
   - Choose a number that supports SMS
   - Copy the number (e.g., `+1234567890`)

### 2. Set Environment Variables

**For Local Development:**

Edit `run-dev.sh` and uncomment/add these lines:

```bash
export TWILIO_ACCOUNT_SID="ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export TWILIO_AUTH_TOKEN="your_auth_token_here"
export TWILIO_PHONE_NUMBER="+1234567890"
```

**For Production (Railway):**

1. Go to your Railway project
2. Click "Variables" tab
3. Add these 3 variables:
   - `TWILIO_ACCOUNT_SID`
   - `TWILIO_AUTH_TOKEN`
   - `TWILIO_PHONE_NUMBER`

### 3. Restart Your Backend

```bash
# Stop the backend (Ctrl+C)
# Then restart:
cd go-kkt-backend
./run-dev.sh
```

### 4. Test It!

1. Go to Payments section
2. Click "Pay" on any Base Amount or Advance Amount
3. Check backend logs - you should see: "SMS sent successfully..."
4. Vendor should receive SMS on their phone

## 📱 Example SMS Message:

```
Dear MG Transport, Payment of ₹ 10,000.00 has been processed for Trip Sheet 000001/2025-2026. Payment Type: Base Amount. Thank you! - KK Transport
```

## ⚠️ Important Notes:

- **Trial Account:** Twilio trial accounts can only send SMS to verified numbers. Verify your number in Twilio Console.
- **Phone Format:** Indian numbers are automatically formatted (e.g., `9876543210` → `+919876543210`)
- **Costs:** Check Twilio pricing for SMS costs in your region
- **Testing:** If SMS fails, payment still works - check backend logs for errors

## 🆘 Troubleshooting:

**SMS not sending?**
- Check all 3 environment variables are set correctly
- Verify Twilio phone number format (must include + and country code)
- Check Twilio Console → Monitor → Logs for errors
- Ensure vendor has valid mobile number in database

**Need help?**
- Twilio Docs: https://www.twilio.com/docs/sms
- Check backend logs for detailed error messages

