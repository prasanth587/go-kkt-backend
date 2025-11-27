package tripmanagement

import (
	"errors"
	"fmt"
	"strings"

	"go-transport-hub/constant"
	"go-transport-hub/dtos"
)

func (trp *TripSheetObj) validateTrip(tripSheetReq dtos.CreateTripSheetHeader) error {

	tripSheetTypeList := []string{constant.LOCAL_SCHEDULED_TRIP, constant.LOCAL_ADHOC_TRIP, constant.LINE_HUAL_SCHEDULED_TRIP, constant.LINE_HUAL_ADHOC_TRIP}

	exists := false
	for _, imgType := range tripSheetTypeList {
		if imgType == tripSheetReq.TripSheetType {
			exists = true
			break
		}
	}
	if !exists {
		return errors.New("trip sheet type not valid")
	}

	if tripSheetReq.OrgId == 0 {
		trp.l.Error("Error OrgId: org should not empty")
		return errors.New("org should not empty")
	}
	if tripSheetReq.TripSheetType == "" {
		trp.l.Error("Error TripSheetType: seelct trip sheet type")
		return errors.New("select trip sheet type")
	}
	if tripSheetReq.BranchID == 0 {
		trp.l.Error("Error BranchID: select brach")
		return errors.New("select brach")
	}
	if tripSheetReq.TripSheetNum == "" {
		trp.l.Error("Error TripSheetNum: trip sheet number should not empty")
		return errors.New("trip sheet number should not empty")
	}
	if tripSheetReq.CustomerID == 0 {
		trp.l.Error("Error CustomerID: select customer")
		return errors.New("select customer")
	}
	trp.l.Info("tripSheetReq.LoadingPointIDs : ", tripSheetReq.LoadingPointIDs, tripSheetReq.UnLoadingPointIDs)
	if len(tripSheetReq.LoadingPointIDs) == 0 {
		trp.l.Error("Error LoadingPointID: select from (loading point)")
		return errors.New("select from (loading point)")
	}
	if len(tripSheetReq.UnLoadingPointIDs) == 0 {
		trp.l.Error("Error UnLoadingPointID: select to (unloading point)")
		return errors.New("select to (unloading point)")
	}
	if tripSheetReq.TripType == "" {
		trp.l.Error("Error TripType: select trip type")
		return errors.New("select trip type")
	}
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_ADHOC_TRIP) {
		if tripSheetReq.LoadHoursType == "" {
			trp.l.Error("Error TripType: Load hours type should not empty")
			return errors.New("load hours type should not empty")
		}
	}
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_ADHOC_TRIP) {
		if tripSheetReq.ZonalName == "" {
			trp.l.Error("Error ZonalName: zonal name should not empty")
			return errors.New("zonal name should not empty")
		}
	}

	if tripSheetReq.OpenTripDateTime == "" {
		trp.l.Error("Error OpenTripDateTime: open trip date should not empty")
		return errors.New("open trip date should not empty")
	}
	if tripSheetReq.VehicleCapacityTon == "" {
		trp.l.Error("Error VehicleCapacityTon: vehicle capacity ton should not empty")
		return errors.New("vehicle capacity ton should not empty")
	}
	if tripSheetReq.VendorID == 0 {
		trp.l.Error("Error VendorID: vendor required")
		return errors.New("vendor required")
	}
	if tripSheetReq.VehicleNumber == "" {
		trp.l.Error("Error VehicleNumber: vehicle number should not empty")
		return errors.New("vehicle number should not empty")
	}
	if tripSheetReq.VehicleSize == "" {
		trp.l.Error("Error VehicleSize: vehicle size should not empty")
		return errors.New("vehicle size should not empty")
	}
	if tripSheetReq.MobileNumber == "" {
		trp.l.Error("Error MobileNumber: mobile numner should not empty")
		return errors.New("mobile numner should not empty")
	}
	if tripSheetReq.DriverName == "" {
		trp.l.Error("Error DriverName: driver name should not empty")
		return errors.New("driver name should not empty")
	}
	if tripSheetReq.DriverLicenseImage == "" {
		trp.l.Error("Error DriverLicenseImage: driver license required")
		return errors.New("driver license required")
	}

	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_ADHOC_TRIP) {

	}
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_ADHOC_TRIP) {
	}

	return nil
}

func (trp *TripSheetObj) validateUpdateTrip(tripSheetReq dtos.UpdateTripSheetHeader) error {

	//if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) {
	if tripSheetReq.TripSheetType == "" {
		trp.l.Error("Error TripSheetType: seelct trip sheet type")
		return errors.New("select trip sheet type")
	}
	if tripSheetReq.TripSheetNum == "" {
		trp.l.Error("Error TripSheetNum: trip sheet number should not empty")
		return errors.New("trip sheet number should not empty")
	}
	if tripSheetReq.CustomerID == 0 {
		trp.l.Error("Error CustomerID: select customer")
		return errors.New("select customer")
	}
	trp.l.Info("validateUpdateTrip: ", tripSheetReq.LoadingPointIDs, tripSheetReq.UnLoadingPointIDs)
	if len(tripSheetReq.LoadingPointIDs) == 0 {
		trp.l.Error("Error LoadingPointID: select from (loading point)")
		return errors.New("select from (loading point)")
	}
	if len(tripSheetReq.UnLoadingPointIDs) == 0 {
		trp.l.Error("Error UnLoadingPointID: select to (unloading point)")
		return errors.New("select to (unloading point)")
	}
	if tripSheetReq.TripType == "" {
		trp.l.Error("Error TripType: select trip type")
		return errors.New("select trip type")
	}
	if tripSheetReq.OpenTripDateTime == "" {
		trp.l.Error("Error OpenTripDateTime: open trip date should not empty")
		return errors.New("open trip date should not empty")
	}
	if tripSheetReq.VehicleCapacityTon == "" {
		trp.l.Error("Error VehicleCapacityTon: vehicle capacity ton should not empty")
		return errors.New("vehicle capacity ton should not empty")
	}
	if tripSheetReq.VendorID == 0 {
		trp.l.Error("Error VendorID: vendor required")
		return errors.New("vendor required")
	}
	if tripSheetReq.VehicleNumber == "" {
		trp.l.Error("Error VehicleNumber: vehicle number should not empty")
		return errors.New("vehicle number should not empty")
	}
	if tripSheetReq.VehicleSize == "" {
		trp.l.Error("Error VehicleSize: vehicle size should not empty")
		return errors.New("vehicle size should not empty")
	}
	if tripSheetReq.MobileNumber == "" {
		trp.l.Error("Error MobileNumber: mobile numner should not empty")
		return errors.New("mobile numner should not empty")
	}
	if tripSheetReq.DriverName == "" {
		trp.l.Error("Error DriverName: driver name should not empty")
		return errors.New("driver name should not empty")
	}
	if tripSheetReq.LRNUmber == "" {
		trp.l.Error("Error LRNUmber: LR number should not empty")
		return errors.New("LR number should not empty")
	}

	//Customer
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) {
		/*if tripSheetReq.CustomerBaseRate == 0.0 {
			trp.l.Error("Error CustomerBaseRate: customer base rate should not empty")
			return errors.New("customer base rate should not empty")
		}
		if tripSheetReq.CustomerKMCost == 0.0 {
			trp.l.Error("Error CustomerKMCost: customer km cost should not empty")
			return errors.New("customer km cost should not empty")
		}
		if tripSheetReq.CustomerToll == 0.0 {
			trp.l.Error("Error CustomerToll: customer toll should not empty")
			return errors.New("customer toll should not empty")
		}*/

		// if tripSheetReq.VendorBaseRate == 0.0 {
		// 	trp.l.Error("Error DriverName: vendor base rate should not empty")
		// 	return errors.New("VendorBaseRate base rate should not empty")
		// }
		// if tripSheetReq.VendorKMCost == 0.0 {
		// 	trp.l.Error("Error VendorKMCost: vendor km cost should not empty")
		// 	return errors.New("vendor km cost should not empty")
		// }
		// if tripSheetReq.VendorToll == 0.0 {
		// 	trp.l.Error("Error VendorKMCost: vendor toll should not empty")
		// 	return errors.New("vendor toll should not empty")
		// }
	}

	/*if strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_ADHOC_TRIP) {
		if tripSheetReq.CustomerBaseRate == 0.0 {
			trp.l.Error("Error CustomerBaseRate: customer base rate should not empty")
			return errors.New("customer base rate should not empty")
		}
	}*/

	// Vendor
	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_SCHEDULED_TRIP) ||
		strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_ADHOC_TRIP) {
		if tripSheetReq.ZonalName == "" {
			trp.l.Error("Error ZonalName: zonal name should not empty")
			return errors.New("zonal name should not empty")
		}

		/*if tripSheetReq.CustomerPerLoadHire == 0.0 {
			trp.l.Error("Error CustomerPerLoadHire: per load hire should not empty")
			return errors.New("per load hire should not empty")
		}
		if tripSheetReq.CustomerRunningKM == 0.0 {
			trp.l.Error("Error CustomerRunningKM: running km should not empty")
			return errors.New("running km should not empty")
		}
		if tripSheetReq.CustomerPerKMPrice == 0.0 {
			trp.l.Error("Error CustomerPerKMPrice: per km price should not empty")
			return errors.New("per km price should not empty")
		}
		if tripSheetReq.CustomerPlacedVehicleSize == "" {
			trp.l.Error("Error CustomerPlacedVehicleSize: placed vehicle size should not empty")
			return errors.New("placed vehicle size should not empty")
		}
		if tripSheetReq.CustomerLoadCancelled == "" {
			trp.l.Error("Error CustomerLoadCancelled: load cancelled should not empty")
			return errors.New("load cancelled should not empty")
		}*/
	}

	/*if tripSheetReq.CustomerTotalHire == 0.0 {
		trp.l.Error("Error CustomerTotalHire: customer total hire should not empty")
		return errors.New("customer total hire should not empty")
	}
	if tripSheetReq.CustomerCloseTripDateTime == "" {
		trp.l.Error("Error CustomerTotalHire: close trip date & time should not empty")
		return errors.New("close trip date & time should not empty")
	} */

	// Vendor
	// if tripSheetReq.VendorTotalHire == 0.0 {
	// 	trp.l.Error("Error VendorTotalHire: vendor toll should not empty")
	// 	return errors.New("vendor toll should not empty")
	// }
	// if tripSheetReq.VendorAdvance == 0.0 {
	// 	trp.l.Error("Error VendorAdvance: vendor advance amount not empty")
	// 	return errors.New("vendor advance amount not empty")
	// }
	// if tripSheetReq.VendorPaidDate == "" {
	// 	trp.l.Error("Error VendorPaidDate: vendor paid date not empty")
	// 	return errors.New("vendor paid date not empty")
	// }
	// if tripSheetReq.VendorBalanceAmount == 0.0 {
	// 	trp.l.Error("Error VendorBalanceAmount: vendor balance amount not empty")
	// 	return errors.New("vendor balance amount not empty")
	// }

	if strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_ADHOC_TRIP) {
		if tripSheetReq.CustomerReportedDateTimeForHaltingCalc == "" {
			// trp.l.Error("Error CustomerReportedDateTimeForHaltingCalc: Reported date time should not empty")
			// return errors.New("reported date time should not empty")
		}
	}

	if tripSheetReq.LRGateImage == "" &&
		(strings.EqualFold(tripSheetReq.TripSheetType, constant.LOCAL_SCHEDULED_TRIP) ||
			strings.EqualFold(tripSheetReq.TripSheetType, constant.LINE_HUAL_SCHEDULED_TRIP)) {
		trp.l.Error("Error LRGateImage: Upload LR/Gate Image should not be empty")
		return errors.New("upload LR/Gate image, should not be empty")
	}

	if fmt.Sprintf("%v", tripSheetReq.PodRequired) == "" {
		tripSheetReq.PodRequired = 0
	}

	return nil
}
