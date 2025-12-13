package tripmanagement

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"go-transport-hub/constant"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	"go-transport-hub/utils"
)

func (trp *TripSheetObj) UpdateTripSheetHeader(tripSheetId int64, tripSheetUpdateReq dtos.UpdateTripSheetHeader) (*dtos.Messge, error) {

	err := trp.validateUpdateTrip(tripSheetUpdateReq)
	if err != nil {
		trp.l.Error("ERROR: UpdateTripSheetHeader", err)
		return nil, err
	}

	tripSheetInfo, errV := trp.tripSheetDao.GetTripSheet(tripSheetId)
	if errV != nil {
		trp.l.Error("ERROR: TripSheet not found", tripSheetInfo, errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(tripSheetInfo)
	trp.l.Info("TripSheet: ******* ", string(jsonBytes))

	// if tripSheetUpdateReq.PodRequired == 0 {

	// }

	tripSheetUpdateReq.TripSubmittedDate = tripSheetInfo.TripSubmittedDate
	tripSheetUpdateReq.LoadStatus = tripSheetInfo.LoadStatus
	if tripSheetInfo.LoadStatus == constant.STATUS_CREATED {
		tripSheetUpdateReq.LoadStatus = constant.STATUS_SUBMITTED
		tripSheetUpdateReq.TripSubmittedDate = utils.GetCurrentDatetimeStr()
	}

	if tripSheetInfo.LoadStatus == constant.STATUS_SUBMITTED {
		tripSheetUpdateReq.TripSubmittedDate = utils.GetCurrentDatetimeStr()
	}

	err1 := trp.tripSheetDao.UpdateTripSheetHeader(tripSheetId, tripSheetUpdateReq)
	if err1 != nil {
		trp.l.Error("ERROR: UpdateTripSheetHeader ", tripSheetUpdateReq.TripSheetNum, err1)
		return nil, err1
	}

	// Send SMS to vendor if payment was made
	trp.sendPaymentSMS(tripSheetInfo, tripSheetUpdateReq)

	loadUnloads, errL := trp.tripSheetDao.GetTripSheetLoadUnLoadPoints(tripSheetId)
	if errL != nil {
		trp.l.Error("ERROR: GetTripSheetLoadUnLoadPoints ", tripSheetUpdateReq.TripSheetNum, errL)
		return nil, errL
	}
	loadingPoints := make([]int64, 0)
	unloadingPoints := make([]int64, 0)

	for _, point := range *loadUnloads {
		switch point.Type {
		case constant.LOADING_POINT:
			loadingPoints = append(loadingPoints, point.LoadingPointID)
		case constant.UN_LOADING_POINT:
			unloadingPoints = append(unloadingPoints, point.LoadingPointID)
		}
	}

	updateLoading := IsNeedsToUpdateLoadUnLoadPoints(loadingPoints, tripSheetUpdateReq.LoadingPointIDs)
	updateUnLoading := IsNeedsToUpdateLoadUnLoadPoints(unloadingPoints, tripSheetUpdateReq.UnLoadingPointIDs)

	if updateLoading {
		errD := trp.tripSheetDao.DeleteTripSheetLoadUnLoadPoints(tripSheetId, constant.LOADING_POINT)
		if errD != nil {
			trp.l.Error("ERROR: DeleteTripSheetLoadUnLoadPoints ", tripSheetUpdateReq.TripSheetNum, errD)
			return nil, errD
		}

		for _, loadingPointId := range tripSheetUpdateReq.LoadingPointIDs {
			errD := trp.tripSheetDao.SaveTripSheetLoadingPoint(uint64(tripSheetId), uint64(loadingPointId), constant.LOADING_POINT)
			if errD != nil {
				trp.l.Error("ERROR: SaveTripSheetLoadingPoint", errD)
				return nil, errD
			}
		}
		trp.l.Info("upadated loading points : ", updateLoading, loadingPoints, tripSheetUpdateReq.LoadingPointIDs)
	}

	if updateUnLoading {
		errD := trp.tripSheetDao.DeleteTripSheetLoadUnLoadPoints(tripSheetId, constant.UN_LOADING_POINT)
		if errD != nil {
			trp.l.Error("ERROR: DeleteTripSheetLoadUnLoadPoints ", tripSheetUpdateReq.TripSheetNum, errD)
			return nil, errD
		}

		for _, loadingPointId := range tripSheetUpdateReq.UnLoadingPointIDs {
			errD := trp.tripSheetDao.SaveTripSheetLoadingPoint(uint64(tripSheetId), uint64(loadingPointId), constant.UN_LOADING_POINT)
			if errD != nil {
				trp.l.Error("ERROR: SaveTripSheetLoadingPoint", errD)
				return nil, errD
			}
		}
		trp.l.Info("upadated unloading points : ", updateLoading, loadingPoints, tripSheetUpdateReq.LoadingPointIDs)
	}

	// if oldStatus != tripSheetUpdateReq.LoadStatus {
	// 	var err error

	// 	switch tripSheetUpdateReq.LoadStatus {
	// 	case constant.STATUS_SUBMITTED:
	// 	case constant.STATUS_DELIVERED:
	// 	case constant.STATUS_CLOSED:
	// 	case constant.STATUS_COMPLETED:
	// 	default:
	// 	}

	// 	if err != nil {
	// 	}
	// }

	trp.l.Info("tripsheet updated successfully! : ", tripSheetUpdateReq.TripSheetNum)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("tripsheet updated successfully: %v", tripSheetUpdateReq.TripSheetNum)
	return &roleResponse, nil
}

func IsNeedsToUpdateLoadUnLoadPoints(existingArray, reqArray []int64) bool {
	if len(existingArray) == len(reqArray) {

		slices.Sort(existingArray)
		slices.Sort(reqArray)

		for i := range existingArray {
			if existingArray[i] != reqArray[i] {
				return true
			}
		}

		return false
	}

	return true
}

// sendPaymentSMS sends SMS to vendor when payment is made
func (trp *TripSheetObj) sendPaymentSMS(oldTripSheet *dtos.TripSheet, newTripSheet dtos.UpdateTripSheetHeader) {
	oldPaidDate := oldTripSheet.VendorPaidDate
	newPaidDate := newTripSheet.VendorPaidDate

	trp.l.Info("=== SMS DEBUG: Checking payment SMS ===")
	trp.l.Info("Old paid date: ", oldPaidDate)
	trp.l.Info("New paid date: ", newPaidDate)
	trp.l.Info("Vendor ID: ", newTripSheet.VendorID)

	// Check if payment was just made (BASE: or ADVANCE: format)
	if newPaidDate == "" || newPaidDate == oldPaidDate {
		trp.l.Info("SMS DEBUG: No payment change detected, skipping SMS")
		return // No payment change
	}

	// Check if it's a payment (not declined)
	if strings.Contains(newPaidDate, "DECLINED") {
		trp.l.Info("SMS DEBUG: Payment declined, skipping SMS")
		return // Don't send SMS for declined payments
	}

	var paymentType string
	var amount float64
	var isBasePaid, isAdvancePaid bool

	// Check if base payment was made
	if strings.Contains(newPaidDate, "BASE:") && !strings.Contains(newPaidDate, "BASE:DECLINED") {
		if !strings.Contains(oldPaidDate, "BASE:") || strings.Contains(oldPaidDate, "BASE:DECLINED") {
			isBasePaid = true
			paymentType = "Base Amount"
			amount = newTripSheet.VendorBaseRate
		}
	}

	// Check if advance payment was made
	if strings.Contains(newPaidDate, "ADVANCE:") && !strings.Contains(newPaidDate, "ADVANCE:DECLINED") {
		if !strings.Contains(oldPaidDate, "ADVANCE:") || strings.Contains(oldPaidDate, "ADVANCE:DECLINED") {
			isAdvancePaid = true
			if paymentType == "" {
				paymentType = "Advance Amount"
				amount = newTripSheet.VendorAdvance
			} else {
				// Both payments made
				paymentType = "Base Amount and Advance Amount"
				amount = newTripSheet.VendorBaseRate + newTripSheet.VendorAdvance
			}
		}
	}

	if !isBasePaid && !isAdvancePaid {
		trp.l.Info("SMS DEBUG: No new payment detected (isBasePaid: ", isBasePaid, ", isAdvancePaid: ", isAdvancePaid, "), skipping SMS")
		return // No new payment made
	}

	trp.l.Info("SMS DEBUG: Payment detected! Type: ", paymentType, ", Amount: ", amount)

	// Get vendor info and mobile number
	vendorDao := daos.NewVendorObj(trp.l, trp.dbConnMSSQL)

	// Get vendor basic info (for vendor name)
	vendor, err := vendorDao.GetVendorV1(int64(newTripSheet.VendorID))
	if err != nil {
		trp.l.Error("ERROR: Failed to get vendor for SMS: ", err)
		return
	}
	trp.l.Info("SMS DEBUG: Vendor name: ", vendor.VendorName)

	// Get vendor contact info (for mobile number - contact person's mobile number)
	contactInfo, err := vendorDao.GetVendorContactInfo(int64(newTripSheet.VendorID))
	if err != nil {
		trp.l.Error("ERROR: Failed to get vendor contact info for SMS: ", err)
		return
	}

	trp.l.Info("SMS DEBUG: Contact info count: ", len(*contactInfo))

	// Get mobile number from contact info (use contact_number1, fallback to contact_number2)
	var mobileNumber string
	if contactInfo != nil && len(*contactInfo) > 0 {
		firstContact := (*contactInfo)[0]
		trp.l.Info("SMS DEBUG: First contact - Name: ", firstContact.ContactPersonName, ", Number1: ", firstContact.ContactNumber1, ", Number2: ", firstContact.ContactNumber2)
		if firstContact.ContactNumber1 != "" {
			mobileNumber = firstContact.ContactNumber1
			trp.l.Info("SMS DEBUG: Using contact_number1 for SMS: ", mobileNumber)
		} else if firstContact.ContactNumber2 != "" {
			mobileNumber = firstContact.ContactNumber2
			trp.l.Info("SMS DEBUG: Using contact_number2 for SMS: ", mobileNumber)
		}
	}

	if mobileNumber == "" {
		trp.l.Info("SMS DEBUG: Vendor mobile number not found in contact info, skipping SMS. Vendor ID: ", newTripSheet.VendorID)
		return
	}

	// Format amount in Indian currency
	formattedAmount := utils.FormatIndianCurrency(amount)

	// Get frontend URL from environment variable
	// Production: https://kktransport.netlify.app
	// Development: http://localhost:3000
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		// Default to production URL if not set
		frontendURL = "https://kktransport.netlify.app"
	}

	// Create route link
	tripSheetId := oldTripSheet.TripSheetID
	routeLink := fmt.Sprintf("%s/route/%d", frontendURL, tripSheetId)

	// Create SMS message with route link
	message := fmt.Sprintf("Dear %s,\n\nPayment of %s has been processed for Trip Sheet %s.\nPayment Type: %s\n\nView Route: %s\n\nThank you!\nKK Transport",
		vendor.VendorName,
		formattedAmount,
		newTripSheet.TripSheetNum,
		paymentType,
		routeLink,
	)

	// Send SMS (non-blocking - don't fail the payment if SMS fails)
	trp.l.Info("SMS DEBUG: Attempting to send SMS to: ", mobileNumber)
	trp.l.Info("SMS DEBUG: Message: ", message)
	go func() {
		err := utils.SendSMS(mobileNumber, message)
		if err != nil {
			trp.l.Error("ERROR: Failed to send payment SMS to vendor: ", mobileNumber, " Error: ", err)
		} else {
			trp.l.Info("SUCCESS: Payment SMS sent successfully to vendor: ", mobileNumber)
		}
	}()
}
