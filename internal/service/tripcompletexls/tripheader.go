package tripcompletexls

const (
	TSID                        = "TSID"
	TripSheetNumber             = "TripSheet Number"
	TripType                    = "Trip Type"
	TripSheetType               = "TripSheet Type"
	LoadHoursType               = "Load Hours Type"
	OpenTripDateTime            = "Open Trip Date Time"
	CustomerName                = "Customer Name"
	CustomerCode                = "Customer Code"
	From                        = "From Locations"
	To                          = "To Locations"
	VehicleNumber               = "Vehicle Number"
	VehicleSize                 = "Vehicle Size"
	MobileNumber                = "Mobile Number"
	DriverName                  = "Driver Name"
	LRNumber                    = "LR Number"
	PODRequired                 = "POD Required"
	VendorBaseRate              = "Vendor Base Rate"
	VendorKMCost                = "Vendor KM Cost"
	VendorToll                  = "Vendor Toll"
	VendorTotalHire             = "Vendor Total Hire"
	VendorAdvance               = "Vendor Advance"
	VendorPaidDate              = "Vendor Paid Date"
	VendorCrossDock             = "Vendor Cross Dock"
	VendorRemark                = "Vendor Remark"
	VendorDebitAmount           = "Vendor Debit Amount"
	VendorBalanceAmount         = "Vendor Balance Amount"
	VendorBreakDown             = "Vendor Break Down"
	VendorPaidBy                = "Vendor Paid By"
	VendorLoadUnLoadAmount      = "Vendor Load UnLoad Amount"
	VendorHaltingDays           = "Vendor Halting Days"
	VendorHaltingPaid           = "Vendor Halting Paid"
	VendorExtraDelivery         = "Vendor Extra Delivery"
	PODReceived                 = "POD Received"
	LoadStatus                  = "Load Status"
	ZonalName                   = "Zonal Name"
	CustomerInvoiceNo           = "Customer Invoice No"
	CustomerBaseRate            = "Customer Base Rate"
	CustomerKMCost              = "Customer KM Cost"
	CustomerToll                = "Customer Toll"
	CustomerExtraHours          = "Customer Extra Hours"
	CustomerExtraKM             = "Customer Extra KM"
	CustomerTotalHire           = "Customer Total Hire"
	CustomerCloseTripDateTime   = "Customer Close Trip Date Time"
	CustomerPaymentReceivedDate = "Customer Payment Received Date"
	CustomerDebitAmount         = "Customer Debit Amount"
	CustomerBillingRaisedDate   = "Customer Billing Raised Date"
	CustomerPerLoadHire         = "Customer Per Load Hire"
	CustomerRunningKM           = "Customer Running KM"
	CustomerPerKMPrice          = "Customer Per KM Price"
	CustomerPlacedVehicleSize   = "Customer Placed Vehicle Size"
	CustomerLoadCancelled       = "Customer Load Cancelled"
	CustomerReportedForHalting  = "Customer Reported For Halting (Date)"
	CustomerRemark              = "Customer Remark"
)

var DRAFT_TRIP_SHEET_HEADER = []string{
	TripSheetNumber,
	TripType,
	TripSheetType,
	LoadHoursType,
	OpenTripDateTime,
	CustomerName,
	CustomerCode,
	From,
	To,
	VehicleNumber,
	VehicleSize,
	MobileNumber,
	DriverName,
	LRNumber,
	PODRequired,
	ZonalName,
	LoadStatus,
	CustomerPerLoadHire,
	CustomerRunningKM,
	CustomerBaseRate,
	CustomerPerKMPrice,
	CustomerKMCost,
	CustomerToll,
	CustomerExtraHours,
	CustomerExtraKM,
	CustomerTotalHire,
	CustomerCloseTripDateTime,
	CustomerPaymentReceivedDate,
	CustomerDebitAmount,
	CustomerBillingRaisedDate,
	CustomerLoadCancelled,
	CustomerReportedForHalting,
	CustomerRemark,
}

var TRIP_SHEET_HEADER = []string{
	TSID,
	TripSheetNumber,
	TripType,
	TripSheetType,
	LoadHoursType,
	OpenTripDateTime,
	CustomerName,
	CustomerCode,
	From,
	To,
	VehicleNumber,
	VehicleSize,
	MobileNumber,
	DriverName,
	LRNumber,
	PODRequired,
	VendorBaseRate,
	VendorKMCost,
	VendorToll,
	VendorTotalHire,
	VendorAdvance,
	VendorPaidDate,
	VendorCrossDock,
	VendorRemark,
	VendorDebitAmount,
	VendorBalanceAmount,
	VendorBreakDown,
	VendorPaidBy,
	VendorLoadUnLoadAmount,
	VendorHaltingDays,
	VendorHaltingPaid,
	VendorExtraDelivery,
	PODReceived,
	LoadStatus,
	ZonalName,
	CustomerInvoiceNo,
	CustomerBaseRate,
	CustomerKMCost,
	CustomerToll,
	CustomerExtraHours,
	CustomerExtraKM,
	CustomerTotalHire,
	CustomerCloseTripDateTime,
	CustomerPaymentReceivedDate,
	CustomerDebitAmount,
	CustomerBillingRaisedDate,
	CustomerPerLoadHire,
	CustomerRunningKM,
	CustomerPerKMPrice,
	CustomerPlacedVehicleSize,
	CustomerLoadCancelled,
	CustomerReportedForHalting,
	CustomerRemark,
}
