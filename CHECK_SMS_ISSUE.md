# Debugging SMS Not Received

## Quick Checks:

### 1. Check Backend Logs
When you clicked "Pay", what did you see in the backend terminal?

Look for one of these messages:

**✅ Success:**
```
SMS sent successfully to +91XXXXXXXXXX. Message SID: SM...
Payment SMS sent successfully to vendor: +91XXXXXXXXXX
```

**❌ Error - Phone not verified:**
```
ERROR: Failed to send payment SMS to vendor: +91XXXXXXXXXX: Twilio API returned error: ... (status: 400)
```

**❌ Error - No mobile number:**
```
Vendor mobile number not found, skipping SMS
```

**❌ Error - SMS not configured:**
```
SMS not configured. Would send to +91XXXXXXXXXX: [message]
```

### 2. Check Vendor Mobile Number
- Does the vendor have a mobile number in the database?
- What is the vendor's mobile number?
- Is it in the correct format? (10 digits for India, e.g., `9876543210`)

### 3. Check if Number is Verified
- Go to: https://console.twilio.com/us1/develop/phone-numbers/manage/verified
- Is the vendor's phone number listed there?
- **Trial accounts can ONLY send to verified numbers!**

### 4. Check Twilio Console Logs
- Go to: https://console.twilio.com/us1/monitor/logs/sms
- Look for any failed attempts
- Check error messages

## Common Issues:

### Issue: "Vendor mobile number not found"
**Solution:** Make sure the vendor has a mobile number in the database

### Issue: "Phone number not verified"
**Solution:** Verify the vendor's phone number in Twilio Console first

### Issue: "SMS not configured"
**Solution:** Restart backend after setting Twilio credentials

### Issue: Payment works but no SMS
**Solution:** Check backend logs - SMS might be failing silently. Payment still works even if SMS fails.

## Next Steps:
1. **Check backend logs** - What error message do you see?
2. **Verify vendor number** - Is it verified in Twilio?
3. **Check vendor database** - Does vendor have mobile number?

