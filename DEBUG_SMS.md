# Debugging SMS Issues

## Quick Checklist:

### 1. ✅ Is the phone number verified in Twilio?
- Go to: https://console.twilio.com/us1/develop/phone-numbers/manage/verified
- Check if the vendor's phone number is listed there
- **Trial accounts can ONLY send to verified numbers!**

### 2. ✅ Check Backend Logs
When you click "Pay", look for these messages in your backend terminal:

**If SMS is working:**
```
SMS sent successfully to +91XXXXXXXXXX. Message SID: SM...
Payment SMS sent successfully to vendor: +91XXXXXXXXXX
```

**If SMS failed:**
```
ERROR: Failed to send payment SMS to vendor: +91XXXXXXXXXX: [error message]
```

**If SMS not configured:**
```
SMS not configured. Would send to +91XXXXXXXXXX: [message]
```

### 3. ✅ Check Vendor Mobile Number
- Make sure the vendor has a mobile number in the database
- Check the format (should be 10 digits for India, e.g., `9876543210`)

### 4. ✅ Test SMS Directly
Run this test command (replace with your verified phone number):

```bash
cd go-kkt-backend
go run test_sms.go
```

Edit `test_sms.go` and change the `testNumber` to your verified phone number first!

### 5. ✅ Check Twilio Console Logs
- Go to: https://console.twilio.com/us1/monitor/logs/sms
- Check if there are any failed attempts
- Look for error messages

## Common Issues:

### Issue: "SMS not configured"
**Solution:** Make sure `run-dev.sh` has the Twilio credentials uncommented and restart backend

### Issue: "Phone number not verified"
**Solution:** Verify the phone number in Twilio Console first

### Issue: "Invalid phone number format"
**Solution:** Make sure vendor mobile number is correct (10 digits for India)

### Issue: "Payment works but no SMS"
**Solution:** Check backend logs - SMS might be failing silently. The payment still works even if SMS fails.

## Next Steps:
1. Check backend logs when you click "Pay"
2. Verify the phone number in Twilio
3. Test SMS directly using `test_sms.go`
4. Share the error message from logs if you see one

