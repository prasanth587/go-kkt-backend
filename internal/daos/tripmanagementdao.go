package daos

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/dtos/schema"
)

type TripSheetObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewTripSheetObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *TripSheetObj {
	return &TripSheetObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type TripSheetDao interface {
	GetTotalCount(whereQuery string) int64
	CreateTripSheet(tripSheetNumber, createTripSheetQuery string) (int64, error)
	CreateTripSheetBuildQuery(tripSheetReq dtos.CreateTripSheetHeader) string
	SaveTripSheetLoadingPoint(tripSheetID, loadingPointID uint64, lpType string) error
	GetTripSheets(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.TripSheet, error)
	UpdateTripSheetImagePath(updateQuery string) error
	GetTripSheet(tripSheetId int64) (*dtos.TripSheet, error)
	BuildWhereQuery(orgId int64, tripSheetid, tripStatus, tripSearchText, fromDate, toDate, podRequired, podReceived string) string
	UpdateTripSheetHeader(tripSheetId int64, tripSheetUpdateReq dtos.UpdateTripSheetHeader) error
	GetTripSheetLoadUnLoadPoints(tripSheetID int64) (*[]dtos.TripSheetLoadUnLoadPoints, error)
	DeleteTripSheetLoadUnLoadPoints(tripSheetID int64, lpType string) error
	CancelTripSheet(tripSheetID int64, status string) error
	CancelTripSheetUpdateToPOD(tripSheetID int64, status string) error
	GetTripStats(whereQuery string) (*[]dtos.TripStats, error)

	// vehicle size type
	GetVehicleSizeType(vehicleSizeID int64) (*dtos.VehicleSizeTypeObj, error)
	GetUserLoginDetails(loginId string) (*schema.UserLogin, error)
}

func (rl *TripSheetObj) BuildWhereQuery(orgId int64, tripSheetid, tripStatus, tripSearchText, fromDate, toDate, podRequired, podReceived string) string {

	whereQuery := fmt.Sprintf("WHERE org_id = '%v'", orgId)

	if tripSheetid != "" {
		whereQuery = fmt.Sprintf(" %v AND trip_sheet_id = '%v'", whereQuery, tripSheetid)
	}
	if podRequired != "" {
		whereQuery = fmt.Sprintf(" %v AND pod_required = '%v'", whereQuery, podRequired)
	}
	if podReceived != "" {
		whereQuery = fmt.Sprintf(" %v AND pod_received = '%v'", whereQuery, podReceived)
	}
	// if tripStatus != "" {
	// 	whereQuery = fmt.Sprintf(" %v AND load_status = '%v'", whereQuery, tripStatus)
	// }
	if len(tripStatus) > 0 {
		var quoted []string
		tripStatusSlice := strings.Split(tripStatus, ",")
		for _, s := range tripStatusSlice {
			quoted = append(quoted, fmt.Sprintf("'%s'", s))
		}
		whereQuery = fmt.Sprintf(" %v AND load_status IN (%s)", whereQuery, strings.Join(quoted, ", "))
	}

	if fromDate != "" && toDate != "" {
		whereQuery = fmt.Sprintf(" %v AND (open_trip_date_time >= '%v' AND open_trip_date_time <= '%v' )", whereQuery, fromDate, toDate)
	}

	if tripSearchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (trip_sheet_num LIKE '%%%v%%' OR trip_sheet_type LIKE '%%%v%%' OR load_hours_type LIKE '%%%v%%' OR open_trip_date_time LIKE '%%%v%%' OR zonal_name LIKE '%%%v%%' OR vehicle_number LIKE '%%%v%%' OR vehicle_size LIKE '%%%v%%' OR mobile_number LIKE '%%%v%%' OR load_status LIKE '%%%v%%') ", whereQuery, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText)
	}

	rl.l.Info("tripSheet whereQuery:\n ", whereQuery)

	return whereQuery
}

func (rl *TripSheetObj) UpdateTripSheetHeader(tripSheetId int64, uts dtos.UpdateTripSheetHeader) error {

	updateTripSheetHeader := fmt.Sprintf(`UPDATE trip_sheet_header SET 
        trip_sheet_type = '%v',
        trip_type = '%v',
        vehicle_number = '%v',
        vehicle_size = '%v',
        open_trip_date_time = '%v',
        lr_number = '%v',
        load_hours_type = '%v',
        customer_id = '%v',
        customer_base_rate = '%v',
        customer_km_cost = '%v',
        customer_toll = '%v',
        customer_extra_hours = '%v',
        customer_extra_km = '%v',
        customer_total_hire = '%v',
        customer_close_trip_date_time = '%v',
        customer_invoice_no = '%v',
        customer_payment_received_date = '%v',
        customer_debit_amount = '%v',
        customer_remark = '%v',
        customer_billing_raised_date = '%v',
        pod_required = '%v',
        vendor_id = '%v',
        vendor_base_rate = '%v',
        vendor_km_cost = '%v',
        vendor_toll = '%v',
        vendor_total_hire = '%v',
        vendor_advance = '%v',
        vendor_paid_date = '%v',
        vendor_cross_dock = '%v',
        vendor_remark = '%v',
        vendor_debit_amount = '%v',
        vendor_balance_amount = '%v',
		vehicle_capacity_ton = '%v',
		vendor_monul = '%v',
		vendor_total_amount = '%v',
        mobile_number = '%v',
        driver_name = '%v',
		customer_per_load_hire = '%v',
		customer_running_km = '%v',
		customer_per_km_price = '%v',
		customer_placed_vehicle_size = '%v',
		customer_load_cancelled = '%v',
		vendor_break_down = '%v',
		load_status = '%v',
		customer_reported_date_time_for_halting_calc = '%v',
		vendor_paid_by = '%v',
		vendor_halting_days = '%v',
		vendor_halting_paid = '%v',
		vendor_extra_delivery = '%v',
		vendor_load_unload_amount = '%v',
		pod_received = '%v',
		zonal_name = '%v',
		trip_submitted_date = '%v',
		vehicle_size_id = '%v',
		vendor_commission = '%v',
		managed_by = '%v'
    WHERE trip_sheet_id = '%v'`,
		uts.TripSheetType, uts.TripType, uts.VehicleNumber, uts.VehicleSize, uts.OpenTripDateTime, uts.LRNUmber, uts.LoadHoursType, uts.CustomerID,
		uts.CustomerBaseRate, uts.CustomerKMCost, uts.CustomerToll, uts.CustomerExtraHours, uts.CustomerExtraKM, uts.CustomerTotalHire, uts.CustomerCloseTripDateTime,
		uts.CustomerInvoiceNo, uts.CustomerPaymentReceivedDate, uts.CustomerDebitAmount, uts.CustomerRemark, uts.CustomerBillingRaisedDate, uts.PodRequired,
		uts.VendorID, uts.VendorBaseRate, uts.VendorKMCost, uts.VendorToll, uts.VendorTotalHire, uts.VendorAdvance, uts.VendorPaidDate, uts.VendorCrossDock,
		uts.VendorRemark, uts.VendorDebitAmount, uts.VendorBalanceAmount, uts.VehicleCapacityTon, uts.VendorMonul, uts.VendorTotalAmount, uts.MobileNumber, uts.DriverName,
		uts.CustomerPerLoadHire, uts.CustomerRunningKM, uts.CustomerPerKMPrice, uts.CustomerPlacedVehicleSize, uts.CustomerLoadCancelled, uts.VendorBreakDown,
		uts.LoadStatus, uts.CustomerReportedDateTimeForHaltingCalc, uts.VendorPaidBy, uts.VendorHaltingDays,
		uts.VendorHaltingPaid, uts.VendorExtraDelivery, uts.VendorLoadUnLoadAmount, uts.PodReceived, uts.ZonalName,
		uts.TripSubmittedDate, uts.VehicleSizeID, uts.VendorCommission, uts.UserLoginId, tripSheetId)

	rl.l.Info("UpdateTripSheetHeader query: ", updateTripSheetHeader)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateTripSheetHeader)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVehicle): ", err)
		return err
	}
	rl.l.Info("Vehicle updated successfully: ", uts.TripSheetNum)

	return nil
}

func (rl *TripSheetObj) GetTotalCount(whereQuery string) int64 {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM trip_sheet_header %v`, whereQuery)
	rl.l.Info(" GetTotalCount select query: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (br *TripSheetObj) CreateTripSheet(tripSheetNumber, createTripSheetQuery string) (int64, error) {

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(createTripSheetQuery)
	if err != nil {
		br.l.Error("Error db.Exec(CreateTripSheet): ", err)
		return 0, err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		br.l.Error("Error db.Exec(CreateTripSheet):", createdId, err)
		return 0, err
	}

	br.l.Info("TripSheet created successfully: ", createdId, tripSheetNumber)
	return createdId, err
}

func (trp *TripSheetObj) CreateTripSheetBuildQuery(tripSheetReq dtos.CreateTripSheetHeader) string {

	createTripSheetQuery := fmt.Sprintf(`INSERT INTO trip_sheet_header (
		trip_sheet_num, trip_type, trip_sheet_type, 
		load_hours_type, open_trip_date_time, branch_id, 
		customer_id, 
		vendor_id, vehicle_capacity_ton, vehicle_number, 
		vehicle_size, vehicle_size_id, mobile_number, driver_name, 
		driver_license_image, org_id, zonal_name, load_status, created_by, managed_by)
	VALUES
		('%v', '%v', '%v', 
		'%v', '%v', '%v', 
		'%v', 
		'%v', '%v', '%v', 
		'%v', '%v', '%v', '%v', 
		'%v','%v','%v', '%v','%v','%v')`,
		tripSheetReq.TripSheetNum,
		tripSheetReq.TripType,
		tripSheetReq.TripSheetType,
		tripSheetReq.LoadHoursType,
		tripSheetReq.OpenTripDateTime,
		tripSheetReq.BranchID,
		tripSheetReq.CustomerID,
		tripSheetReq.VendorID,
		tripSheetReq.VehicleCapacityTon,
		tripSheetReq.VehicleNumber,
		tripSheetReq.VehicleSize,
		tripSheetReq.VehicleSizeID,
		tripSheetReq.MobileNumber,
		tripSheetReq.DriverName,
		tripSheetReq.DriverLicenseImage, tripSheetReq.OrgId, tripSheetReq.ZonalName, tripSheetReq.LoadStatus, tripSheetReq.UserLoginId, tripSheetReq.UserLoginId)

	trp.l.Info("CreateTripSheetBuildQuery: ", tripSheetReq.OpenTripDateTime, createTripSheetQuery)
	return createTripSheetQuery
}

func (br *TripSheetObj) SaveTripSheetLoadingPoint(tripSheetID, loadingPointID uint64, lpType string) error {

	insertTripSheetLoadUnloadPointsQuery := fmt.Sprintf(`INSERT INTO trip_sheet_header_load_unload_points 
	(trip_sheet_id, loading_point_id, type)
	VALUES ('%v', '%v', '%v')`, tripSheetID, loadingPointID, lpType)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(insertTripSheetLoadUnloadPointsQuery)
	if err != nil {
		br.l.Error("Error db.Exec(SaveTripSheetLoadingPoint): ", err)
		return err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		br.l.Error("Error db.Exec(SaveTripSheetLoadingPoint):", createdId, err)
		return err
	}

	br.l.Info("loading_point created successfully: ", createdId, loadingPointID)

	return nil
}

func (ts *TripSheetObj) GetTripSheets(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.TripSheet, error) {
	list := []dtos.TripSheet{}

	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)
	// SQL query to fetch trip sheets
	query := fmt.Sprintf(`
        SELECT 
            trip_sheet_id, trip_sheet_num, trip_type, trip_sheet_type, 
            load_hours_type, open_trip_date_time, branch_id, customer_id, 
            vendor_id, 
            vehicle_capacity_ton, vehicle_number, vehicle_size, 
            mobile_number, driver_name, driver_license_image, lr_gate_image, lr_number,
            customer_base_rate, customer_km_cost, customer_toll, customer_extra_hours,
            customer_extra_km, customer_total_hire, customer_close_trip_date_time,
            customer_invoice_no, customer_payment_received_date, customer_debit_amount,
            customer_remark, customer_billing_raised_date, pod_required,
            vendor_base_rate, vendor_km_cost, vendor_toll, vendor_total_hire,
            vendor_advance, vendor_paid_date, vendor_cross_dock, vendor_remark,
            vendor_debit_amount, vendor_balance_amount,pod_received,
			customer_per_load_hire, customer_running_km, customer_per_km_price, customer_placed_vehicle_size,
			customer_load_cancelled, vendor_paid_by, vendor_load_unload_amount, vendor_halting_days,
			vendor_halting_paid, vendor_extra_delivery, load_status, customer_reported_date_time_for_halting_calc, vendor_break_down, zonal_name,
			trip_submitted_date, trip_closed_date, trip_delivered_date, trip_completed_date, vendor_monul, vendor_total_amount, vehicle_size_id, vendor_commission
        FROM trip_sheet_header %v`, whereQuery)

	ts.l.Info("TripSheet Query:\n", query)

	// Execute the query
	rows, err := ts.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		ts.l.Error("Error fetching trip sheets: ", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results
	for rows.Next() {
		var (
			// Core trip details
			tripSheetNum, tripType, tripSheetType, loadHoursType, openTripDateTime, vehicleCapacityTon, vehicleNumber, vehicleSize, mobileNumber, driverName,
			driverLicenseImage, lrGateImage, lrNumber, customerCloseTripDT, customerInvoiceNo,
			customerPaymentDate, customerRemark, customerBillingDate, vendorPaidDate, vendorCrossDock, vendorRemark,
			customerPlacedVehicleSize, customerLoadCancelled, vendorPaidBy, loadStatus,
			customerReportedDateTimeForHaltingCalc, zonalName, tripSubmittedDate, tripClosedDate, tripDeliveredDate, tripCompletedDate, vendorBreakDown sql.NullString

			tripSheetID, branchID, customerID, vendorID, podRequired, podReceived, vehicleSizeID sql.NullInt64
			customerBaseRate, customerKmCost, customerToll, customerExtraHours, customerExtraKm, customerTotalHire, customerDebitAmount,
			vendorBaseRate, vendorKmCost, vendorToll, vendorTotalHire, vendorAdvance, vendorDebitAmount, vendorBalanceAmount,
			customerPerLoadHire, customerRunningKM, customerPerKMPrice, vendorLoadUnLoadAmount, vendorHaltingDays,
			vendorHaltingPaid, vendorExtraDelivery, vendorMonul, vendorTotalAmount, vendorCommission sql.NullFloat64
		)

		// Scan the row into variables
		err := rows.Scan(
			&tripSheetID,
			&tripSheetNum,
			&tripType,
			&tripSheetType,
			&loadHoursType,
			&openTripDateTime,
			&branchID,
			&customerID,
			&vendorID,
			&vehicleCapacityTon,
			&vehicleNumber,
			&vehicleSize,
			&mobileNumber,
			&driverName,
			&driverLicenseImage,
			&lrGateImage,
			&lrNumber,
			&customerBaseRate,
			&customerKmCost,
			&customerToll,
			&customerExtraHours,
			&customerExtraKm,
			&customerTotalHire,
			&customerCloseTripDT,
			&customerInvoiceNo,
			&customerPaymentDate,
			&customerDebitAmount,
			&customerRemark,
			&customerBillingDate,
			&podRequired,
			&vendorBaseRate,
			&vendorKmCost,
			&vendorToll,
			&vendorTotalHire,
			&vendorAdvance,
			&vendorPaidDate,
			&vendorCrossDock,
			&vendorRemark,
			&vendorDebitAmount,
			&vendorBalanceAmount,
			&podReceived,
			&customerPerLoadHire,
			&customerRunningKM,
			&customerPerKMPrice,
			&customerPlacedVehicleSize,
			&customerLoadCancelled,
			&vendorPaidBy,
			&vendorLoadUnLoadAmount,
			&vendorHaltingDays,
			&vendorHaltingPaid,
			&vendorExtraDelivery,
			&loadStatus,
			&customerReportedDateTimeForHaltingCalc,
			&vendorBreakDown,
			&zonalName,
			&tripSubmittedDate,
			&tripClosedDate,
			&tripDeliveredDate,
			&tripCompletedDate,
			&vendorMonul,
			&vendorTotalAmount,
			&vehicleSizeID,
			&vendorCommission,
		)

		if err != nil {
			ts.l.Error("Error scanning trip sheet row: ", err)
			return nil, err
		}

		// Create a TripSheet instance and populate it with data from the row.
		trip := dtos.TripSheet{
			TripSheetID:                            tripSheetID.Int64,
			TripSheetNum:                           tripSheetNum.String,
			TripType:                               tripType.String,
			TripSheetType:                          tripSheetType.String,
			LoadHoursType:                          loadHoursType.String,
			OpenTripDateTime:                       openTripDateTime.String,
			BranchID:                               branchID.Int64,
			VendorID:                               vendorID.Int64,
			VehicleCapacityTon:                     vehicleCapacityTon.String,
			VehicleNumber:                          vehicleNumber.String,
			VehicleSize:                            vehicleSize.String,
			MobileNumber:                           mobileNumber.String,
			DriverName:                             driverName.String,
			DriverLicenseImage:                     driverLicenseImage.String,
			LRGateImage:                            lrGateImage.String,
			LRNumber:                               lrNumber.String,
			CustomerBaseRate:                       customerBaseRate.Float64,
			CustomerKMCost:                         customerKmCost.Float64,
			CustomerToll:                           customerToll.Float64,
			CustomerExtraHours:                     customerExtraHours.Float64,
			CustomerExtraKM:                        customerExtraKm.Float64,
			CustomerTotalHire:                      customerTotalHire.Float64,
			CustomerCloseTripDateTime:              customerCloseTripDT.String,
			CustomerInvoiceNo:                      customerInvoiceNo.String,
			CustomerPaymentReceivedDate:            customerPaymentDate.String,
			CustomerDebitAmount:                    customerDebitAmount.Float64,
			CustomerRemark:                         customerRemark.String,
			CustomerBillingRaisedDate:              customerBillingDate.String,
			PodRequired:                            podRequired.Int64,
			VendorBaseRate:                         vendorBaseRate.Float64,
			VendorKMCost:                           vendorKmCost.Float64,
			VendorToll:                             vendorToll.Float64,
			VendorTotalHire:                        vendorTotalHire.Float64,
			VendorAdvance:                          vendorAdvance.Float64,
			VendorPaidDate:                         vendorPaidDate.String,
			VendorCrossDock:                        vendorCrossDock.String,
			VendorRemark:                           vendorRemark.String,
			VendorDebitAmount:                      vendorDebitAmount.Float64,
			VendorBalanceAmount:                    vendorBalanceAmount.Float64,
			VendorMonul:                            vendorMonul.Float64,
			VendorTotalAmount:                      vendorTotalAmount.Float64,
			PodReceived:                            podReceived.Int64,
			CustomerPerLoadHire:                    customerPerLoadHire.Float64,
			CustomerRunningKM:                      customerRunningKM.Float64,
			CustomerPerKMPrice:                     customerPerKMPrice.Float64,
			CustomerPlacedVehicleSize:              customerPlacedVehicleSize.String,
			CustomerLoadCancelled:                  customerLoadCancelled.String,
			VendorPaidBy:                           vendorPaidBy.String,
			VendorLoadUnLoadAmount:                 vendorLoadUnLoadAmount.Float64,
			VendorHaltingDays:                      vendorHaltingDays.Float64,
			VendorHaltingPaid:                      vendorHaltingPaid.Float64,
			VendorExtraDelivery:                    vendorExtraDelivery.Float64,
			LoadStatus:                             loadStatus.String,
			CustomerReportedDateTimeForHaltingCalc: customerReportedDateTimeForHaltingCalc.String,
			VendorBreakDown:                        vendorBreakDown.String,
			ZonalName:                              zonalName.String,
			TripSubmittedDate:                      tripSubmittedDate.String,
			TripClosedDate:                         tripClosedDate.String,
			TripDeliveredDate:                      tripDeliveredDate.String,
			TripCompletedDate:                      tripCompletedDate.String,
			VehicleSizeID:                          vehicleSizeID.Int64,
			VendorCommission:                       vendorCommission.Float64,
		}
		loadingPs, err := ts.GetTripSheetLoadingPoint(tripSheetID.Int64, constant.LOADING_POINT)
		if err != nil {
			ts.l.Error("Error fetching trip sheets: ", err)
			return nil, err
		}

		unLoadingPs, err := ts.GetTripSheetUnLoadingPoint(tripSheetID.Int64, constant.UN_LOADING_POINT)
		if err != nil {
			ts.l.Error("Error fetching trip sheets: ", err)
			return nil, err
		}

		customerObj, errC := ts.GetTripSheetCustomer(customerID.Int64)
		if errC != nil {
			ts.l.Error("Error GetTripSheetCustomer trip sheets: ", errC)
			return nil, errC
		}
		trip.Customer = customerObj

		// Vendor details

		vendorObj, errV := ts.GetTripSheetVendor(vendorID.Int64)
		if errV != nil {
			ts.l.Error("Error GetTripSheetVendor trip sheets: ", errV)
			return nil, errC
		}
		trip.Vendor = vendorObj

		trip.LoadingPointIDs = loadingPs
		trip.UnLoadingPointIDs = unLoadingPs
		// Append the populated TripSheet struct to the list.
		list = append(list, trip)
	}

	return &list, nil // Return the list of TripSheets.
}

func (br *TripSheetObj) GetTripSheetLoadingPoint(tripSheetID int64, lp_up_type string) (*[]dtos.TripSheetLoading, error) {
	list := []dtos.TripSheetLoading{}

	loadingpointQuery := fmt.Sprintf(`SELECT tshlup.loading_point_id, lp.city_code FROM  trip_sheet_header_load_unload_points tshlup LEFT JOIN loading_point lp ON tshlup.loading_point_id = lp.loading_point_id WHERE tshlup.trip_sheet_id = '%v' AND type = '%v' ORDER BY tshlup.load_unload_point_id ASC;`, tripSheetID, lp_up_type)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error GetTripSheetLoadingPoint ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cityCode sql.NullString
		var loadingpointId sql.NullInt64

		loadingPointE := &dtos.TripSheetLoading{}
		err := rows.Scan(&loadingpointId, &cityCode)
		if err != nil {
			br.l.Error("Error GetTripSheetLoadingPoint scan: ", err)
			return nil, err
		}
		loadingPointE.LoadingPointId = loadingpointId.Int64
		loadingPointE.CityCode = cityCode.String
		list = append(list, *loadingPointE)
	}

	return &list, nil
}

func (br *TripSheetObj) GetTripSheetCustomer(customerId int64) (*dtos.Customers, error) {
	customerQuery := fmt.Sprintf(`SELECT customer_name, customer_code FROM customers WHERE customer_id = '%v';`, customerId)
	br.l.Info("GetTripSheetCustomer:\n ", customerQuery)
	row := br.dbConnMSSQL.GetQueryer().QueryRow(customerQuery)
	var customerName, customerCode sql.NullString
	err := row.Scan(&customerName, &customerCode)
	if err != nil {
		br.l.Error("Error GetTripSheetCustomer scan: ", err)
		return nil, err
	}

	customerObj := dtos.Customers{}
	customerObj.CustomerId = customerId
	customerObj.CustomerName = fmt.Sprintf("%v - %v", customerCode.String, customerName.String)
	return &customerObj, nil
}

func (br *TripSheetObj) GetTripSheetVendor(vendorId int64) (*dtos.VendorT, error) {
	vendorQuery := fmt.Sprintf(`SELECT vendor_name, vendor_code FROM vendors WHERE vendor_id = '%v';`, vendorId)
	br.l.Info("GetTripSheetVendor:\n ", vendorQuery)
	row := br.dbConnMSSQL.GetQueryer().QueryRow(vendorQuery)
	var vendorName, vendorCode sql.NullString
	err := row.Scan(&vendorName, &vendorCode)
	if err != nil {
		br.l.Error("Error GetTripSheetVendor scan: ", err)
		return nil, err
	}

	VendorObj := dtos.VendorT{}
	VendorObj.VendorId = vendorId
	VendorObj.VendorName = fmt.Sprintf("%v - %v", vendorCode.String, vendorName.String)
	return &VendorObj, nil
}

func (br *TripSheetObj) GetTripSheetUnLoadingPoint(tripSheetID int64, lp_up_type string) (*[]dtos.TripSheetUnLoading, error) {
	list := []dtos.TripSheetUnLoading{}

	loadingpointQuery := fmt.Sprintf(`SELECT tshlup.loading_point_id, lp.city_code FROM  trip_sheet_header_load_unload_points tshlup LEFT JOIN loading_point lp ON tshlup.loading_point_id = lp.loading_point_id WHERE tshlup.trip_sheet_id = '%v' AND type = '%v' ORDER BY tshlup.load_unload_point_id ASC;`, tripSheetID, lp_up_type)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error GetTripSheetUnLoadingPoint ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cityCode sql.NullString
		var loadingpointId sql.NullInt64

		loadingPointE := &dtos.TripSheetUnLoading{}
		err := rows.Scan(&loadingpointId, &cityCode)
		if err != nil {
			br.l.Error("Error GetTripSheetUnLoadingPoint scan: ", err)
			return nil, err
		}
		loadingPointE.UnLoadingPointId = loadingpointId.Int64
		loadingPointE.CityCode = cityCode.String
		list = append(list, *loadingPointE)
	}

	return &list, nil
}

func (ts *TripSheetObj) GetTripSheet(tripSheetId int64) (*dtos.TripSheet, error) {

	// SQL query to fetch trip sheets
	query := fmt.Sprintf(`
        SELECT 
            trip_sheet_num, trip_type, trip_sheet_type, 
            load_hours_type, open_trip_date_time, branch_id, customer_id, 
            vendor_id, vehicle_capacity_ton, vehicle_number, vehicle_size, 
            mobile_number, driver_name, driver_license_image, lr_gate_image, lr_number,
            customer_base_rate, customer_km_cost, customer_toll, customer_extra_hours,
            customer_extra_km, customer_total_hire, customer_close_trip_date_time,
            customer_invoice_no, customer_payment_received_date, customer_debit_amount,
            customer_remark, customer_billing_raised_date, pod_required,
            vendor_base_rate, vendor_km_cost, vendor_toll, vendor_total_hire,
            vendor_advance, vendor_paid_date, vendor_cross_dock, vendor_remark,
            vendor_debit_amount, vendor_balance_amount, pod_received,customer_per_load_hire, customer_running_km, customer_per_km_price, customer_placed_vehicle_size,
			customer_load_cancelled, vendor_paid_by, vendor_load_unload_amount, vendor_halting_days,
			vendor_halting_paid, vendor_extra_delivery, load_status, customer_reported_date_time_for_halting_calc, vendor_break_down,
			trip_submitted_date, trip_closed_date, trip_delivered_date, trip_completed_date, vendor_monul,vendor_commission, vehicle_size_id, org_id
        FROM trip_sheet_header 
        WHERE trip_sheet_id = '%d';`, tripSheetId)

	ts.l.Info("TripSheet Query:\n", query)

	// Execute the query
	row := ts.dbConnMSSQL.GetQueryer().QueryRow(query)

	var (
		// Core trip details
		tripSheetNum, tripType, tripSheetType, loadHoursType, openTripDateTime, vehicleCapacityTon, vehicleNumber, vehicleSize, mobileNumber, driverName,
		driverLicenseImage, lrGateImage, lrNumber, customerCloseTripDT, customerInvoiceNo,
		customerPaymentDate, customerRemark, customerBillingDate, vendorPaidDate, vendorCrossDock, vendorRemark,
		customerPlacedVehicleSize, customerLoadCancelled, vendorPaidBy, loadStatus,
		customerReportedDateTimeForHaltingCalc, tripSubmittedDate, tripClosedDate, tripDeliveredDate, tripCompletedDate, vendorBreakDown sql.NullString

		branchID, customerID, vendorID, podRequired, podReceived, vehicleSizeId, orgID sql.NullInt64
		customerBaseRate, customerKmCost, customerToll, customerExtraHours, customerExtraKm, customerTotalHire, customerDebitAmount,
		vendorBaseRate, vendorKmCost, vendorToll, vendorTotalHire, vendorAdvance, vendorDebitAmount, vendorBalanceAmount,
		customerPerLoadHire, customerRunningKM, customerPerKMPrice, vendorLoadUnLoadAmount, vendorHaltingDays,
		vendorHaltingPaid, vendorExtraDelivery, vendorMonul, vendorCommission sql.NullFloat64
	)
	err := row.Scan(
		&tripSheetNum,
		&tripType,
		&tripSheetType,
		&loadHoursType,
		&openTripDateTime,
		&branchID,
		&customerID,
		&vendorID,
		&vehicleCapacityTon,
		&vehicleNumber,
		&vehicleSize,
		&mobileNumber,
		&driverName,
		&driverLicenseImage,
		&lrGateImage,
		&lrNumber,
		&customerBaseRate,
		&customerKmCost,
		&customerToll,
		&customerExtraHours,
		&customerExtraKm,
		&customerTotalHire,
		&customerCloseTripDT,
		&customerInvoiceNo,
		&customerPaymentDate,
		&customerDebitAmount,
		&customerRemark,
		&customerBillingDate,
		&podRequired,
		&vendorBaseRate,
		&vendorKmCost,
		&vendorToll,
		&vendorTotalHire,
		&vendorAdvance,
		&vendorPaidDate,
		&vendorCrossDock,
		&vendorRemark,
		&vendorDebitAmount,
		&vendorBalanceAmount,
		&podReceived,
		&customerPerLoadHire,
		&customerRunningKM,
		&customerPerKMPrice,
		&customerPlacedVehicleSize,
		&customerLoadCancelled,
		&vendorPaidBy,
		&vendorLoadUnLoadAmount,
		&vendorHaltingDays,
		&vendorHaltingPaid,
		&vendorExtraDelivery,
		&loadStatus,
		&customerReportedDateTimeForHaltingCalc,
		&vendorBreakDown,
		&tripSubmittedDate,
		&tripClosedDate,
		&tripDeliveredDate,
		&tripCompletedDate,
		&vendorMonul,
		&vendorCommission,
		&vehicleSizeId,
		&orgID,
	)

	if err != nil {
		ts.l.Error("Error scanning trip sheet row: ", err)
		return nil, err
	}

	trip := dtos.TripSheet{
		TripSheetID:                            tripSheetId,
		TripSheetNum:                           tripSheetNum.String,
		TripType:                               tripType.String,
		TripSheetType:                          tripSheetType.String,
		LoadHoursType:                          loadHoursType.String,
		OpenTripDateTime:                       openTripDateTime.String,
		BranchID:                               branchID.Int64,
		VendorID:                               vendorID.Int64,
		VehicleCapacityTon:                     vehicleCapacityTon.String,
		VehicleNumber:                          vehicleNumber.String,
		VehicleSize:                            vehicleSize.String,
		MobileNumber:                           mobileNumber.String,
		DriverName:                             driverName.String,
		DriverLicenseImage:                     driverLicenseImage.String,
		LRGateImage:                            lrGateImage.String,
		LRNumber:                               lrNumber.String,
		CustomerBaseRate:                       customerBaseRate.Float64,
		CustomerKMCost:                         customerKmCost.Float64,
		CustomerToll:                           customerToll.Float64,
		CustomerExtraHours:                     customerExtraHours.Float64,
		CustomerExtraKM:                        customerExtraKm.Float64,
		CustomerTotalHire:                      customerTotalHire.Float64,
		CustomerCloseTripDateTime:              customerCloseTripDT.String,
		CustomerInvoiceNo:                      customerInvoiceNo.String,
		CustomerPaymentReceivedDate:            customerPaymentDate.String,
		CustomerDebitAmount:                    customerDebitAmount.Float64,
		CustomerRemark:                         customerRemark.String,
		CustomerBillingRaisedDate:              customerBillingDate.String,
		PodRequired:                            podRequired.Int64,
		VendorBaseRate:                         vendorBaseRate.Float64,
		VendorKMCost:                           vendorKmCost.Float64,
		VendorToll:                             vendorToll.Float64,
		VendorTotalHire:                        vendorTotalHire.Float64,
		VendorAdvance:                          vendorAdvance.Float64,
		VendorPaidDate:                         vendorPaidDate.String,
		VendorCrossDock:                        vendorCrossDock.String,
		VendorRemark:                           vendorRemark.String,
		VendorDebitAmount:                      vendorDebitAmount.Float64,
		VendorBalanceAmount:                    vendorBalanceAmount.Float64,
		PodReceived:                            podReceived.Int64,
		CustomerReportedDateTimeForHaltingCalc: customerReportedDateTimeForHaltingCalc.String,
		VendorBreakDown:                        vendorBreakDown.String,
		LoadStatus:                             loadStatus.String,
		TripSubmittedDate:                      tripSubmittedDate.String,
		TripClosedDate:                         tripClosedDate.String,
		TripDeliveredDate:                      tripDeliveredDate.String,
		TripCompletedDate:                      tripCompletedDate.String,
		VendorMonul:                            vendorMonul.Float64,
		VendorCommission:                       vendorCommission.Float64,
		VendorLoadUnLoadAmount:                 vendorLoadUnLoadAmount.Float64,
		VendorHaltingPaid:                      vendorHaltingPaid.Float64,
		VehicleSizeID:                          vehicleSizeId.Int64,
	}
	loadingPs, err := ts.GetTripSheetLoadingPoint(tripSheetId, constant.LOADING_POINT)
	if err != nil {
		ts.l.Error("Error fetching trip sheets: ", err)
		return nil, err
	}

	unLoadingPs, err := ts.GetTripSheetUnLoadingPoint(tripSheetId, constant.UN_LOADING_POINT)
	if err != nil {
		ts.l.Error("Error fetching trip sheets: ", err)
		return nil, err
	}
	customerObj, errC := ts.GetTripSheetCustomer(customerID.Int64)
	if errC != nil {
		ts.l.Error("Error GetTripSheetCustomer trip sheets: ", errC)
		return nil, errC
	}

	vendorObj, errV := ts.GetTripSheetVendor(vendorID.Int64)
	if errV != nil {
		ts.l.Error("Error GetTripSheetVendor trip sheets: ", errV)
		return nil, errC
	}
	trip.Vendor = vendorObj

	trip.Customer = customerObj
	trip.LoadingPointIDs = loadingPs
	trip.UnLoadingPointIDs = unLoadingPs
	return &trip, nil
}

// func (br *TripSheetObj) GetCustomers(orgId int64) (*[]dtos.Customers, error) {
// 	list := []dtos.Customers{}

// 	loadingpointQuery := fmt.Sprintf(`SELECT customer_id, customer_name, customer_code FROM customers WHERE org_id = '%v' AND is_active = '1' ORDER BY customer_code ASC;`, orgId)

// 	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

// 	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
// 	if err != nil {
// 		br.l.Error("Error GetCustomers ", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var customerName, customerCode sql.NullString

// 		var customerId sql.NullInt64

// 		customerE := &dtos.Customers{}
// 		err := rows.Scan(&customerId, &customerName, &customerCode)
// 		if err != nil {
// 			br.l.Error("Error GetCustomers scan: ", err)
// 			return nil, err
// 		}
// 		customerE.CustomerId = customerId.Int64
// 		customerE.CustomerName = fmt.Sprintf("%v - %v", customerCode.String, customerName.String)
// 		list = append(list, *customerE)
// 	}

// 	return &list, nil
// }

// func (br *TripSheetObj) GetBranches(orgId int64) (*[]dtos.BranchT, error) {
// 	list := []dtos.BranchT{}

// 	branchQuery := fmt.Sprintf(`SELECT branch_id, branch_name, branch_code FROM branch WHERE org_id = '%v' AND is_active='1' ORDER BY branch_code ASC;`, orgId)

// 	br.l.Info("branchQuery:\n ", branchQuery)

// 	rows, err := br.dbConnMSSQL.GetQueryer().Query(branchQuery)
// 	if err != nil {
// 		br.l.Error("Error Branches ", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var branchName, branchCode sql.NullString
// 		var branchId sql.NullInt64

// 		branch := &dtos.BranchT{}
// 		err := rows.Scan(&branchId, &branchName, &branchCode)
// 		if err != nil {
// 			br.l.Error("Error GetBranchs scan: ", err)
// 			return nil, err
// 		}
// 		branch.BranchId = branchId.Int64
// 		branch.BranchName = fmt.Sprintf("%v - %v", branchCode.String, branchName.String)
// 		list = append(list, *branch)
// 	}

// 	return &list, nil
// }

// func (br *TripSheetObj) GetVendors(orgId int64) (*[]dtos.VendorT, error) {
// 	list := []dtos.VendorT{}

// 	vendorQuery := fmt.Sprintf(`SELECT vendor_id, vendor_name, vendor_code FROM vendors WHERE org_id = '%v' AND is_active='1' ORDER BY vendor_code ASC;`, orgId)

// 	br.l.Info("vendorQuery:\n ", vendorQuery)

// 	rows, err := br.dbConnMSSQL.GetQueryer().Query(vendorQuery)
// 	if err != nil {
// 		br.l.Error("Error GetVendors ", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var vendorName, vendorCode sql.NullString
// 		var vendorID sql.NullInt64

// 		vendor := &dtos.VendorT{}
// 		err := rows.Scan(&vendorID, &vendorName, &vendorCode)
// 		if err != nil {
// 			br.l.Error("Error GetVendors scan: ", err)
// 			return nil, err
// 		}
// 		vendor.VendorId = vendorID.Int64
// 		vendor.VendorName = fmt.Sprintf("%v - %v", vendorCode.String, vendorName.String)
// 		list = append(list, *vendor)
// 	}
// 	return &list, nil
// }

// func (br *TripSheetObj) InsertAndGetTripNumber(tripSheeetNumber string) (int64, error) {

// 	tripSheetNum := fmt.Sprintf(`INSERT INTO trip_sheet_num (year) VALUES ( '%v' )`, tripSheeetNumber)

// 	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(tripSheetNum)
// 	if err != nil {
// 		br.l.Error("Error db.Exec(InsertAndGetTripNumber): ", err)
// 		return 0, err
// 	}
// 	createdId, _ := roleResult.LastInsertId()

// 	return createdId, nil

// }

func (rl *TripSheetObj) UpdateTripSheetImagePath(updateQuery string) error {

	rl.l.Info("UpdateTripSheetImagePath Update query: ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateTripSheetImagePath): ", err)
		return err
	}

	return nil
}

func (br *TripSheetObj) GetTripSheetLoadUnLoadPoints(tripSheetID int64) (*[]dtos.TripSheetLoadUnLoadPoints, error) {
	list := []dtos.TripSheetLoadUnLoadPoints{}

	loadingpointQuery := fmt.Sprintf(`SELECT loading_point_id, type FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id = '%v' ORDER BY loading_point_id ASC;`, tripSheetID)

	br.l.Info("GetTripSheetLoadUnLoadPoints:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error GetTripSheetLoadingPoint ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var typeL sql.NullString
		var loadingpointId sql.NullInt64

		loadingPointE := &dtos.TripSheetLoadUnLoadPoints{}
		err := rows.Scan(&loadingpointId, &typeL)
		if err != nil {
			br.l.Error("Error GetTripSheetLoadingPoint scan: ", err)
			return nil, err
		}
		loadingPointE.LoadingPointID = loadingpointId.Int64
		loadingPointE.Type = typeL.String
		list = append(list, *loadingPointE)
	}

	return &list, nil
}

func (trp *TripSheetObj) CancelTripSheet(tripSheetID int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET load_status = '%v' WHERE trip_sheet_id = '%v'`, status, tripSheetID)

	trp.l.Info("CancelTripSheet Update query: ", updateQuery)

	res, err := trp.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		trp.l.Error("Error db.Exec(CancelTripSheet): ", err)
		return err
	}
	trp.l.Info("CancelTripSheet Update successfully: ", res)
	return nil
}

func (trp *TripSheetObj) CancelTripSheetUpdateToPOD(tripSheetID int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE manage_pod SET pod_status = '%v' WHERE trip_sheet_id = '%v'`, status, tripSheetID)

	trp.l.Info("CancelTripSheetUpdateToPOD Update query: ", updateQuery)

	res, err := trp.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		trp.l.Error("Error db.Exec(CancelTripSheetUpdateToPOD): ", err)
		return err
	}
	trp.l.Info("CancelTripSheetUpdateToPOD Update query: ", res)
	return nil
}

func (br *TripSheetObj) GetTripStats(whereQuery string) (*[]dtos.TripStats, error) {
	list := []dtos.TripStats{}

	query := fmt.Sprintf(`
        SELECT trip_sheet_id, trip_sheet_num, trip_type, trip_sheet_type, pod_required, pod_received, load_status FROM trip_sheet_header %v`, whereQuery)

	br.l.Info("GetTripStats query:\n ", query)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		br.l.Error("Error GetTripSheetUnLoadingPoint ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			tripSheetNum, tripType, tripSheetType, loadStatus sql.NullString
			tripSheetID, podRequired, podReceived             sql.NullInt64
		)

		tripStats := &dtos.TripStats{}
		err := rows.Scan(&tripSheetID, &tripSheetNum, &tripType, &tripSheetType, &podRequired, &podReceived, &loadStatus)
		if err != nil {
			br.l.Error("Error GetTripSheetUnLoadingPoint scan: ", err)
			return nil, err
		}
		tripStats.TripSheetID = tripSheetID.Int64
		tripStats.TripSheetNum = tripSheetNum.String
		tripStats.TripSheetType = tripSheetType.String
		tripStats.TripType = tripType.String
		tripStats.LoadStatus = loadStatus.String
		tripStats.PodRequired = podRequired.Int64
		tripStats.PodReceived = podReceived.Int64
		list = append(list, *tripStats)
	}
	return &list, nil
}

func (rl *TripSheetObj) GetVehicleSizeType(vehicleSizeID int64) (*dtos.VehicleSizeTypeObj, error) {

	vehicleQuery := "SELECT vehicle_size_id, vehicle_size, vehicle_type, status, is_active  FROM vehicle_size_type WHERE vehicle_size_id = ? "
	rl.l.Info("GetVehicleSizeType: ", vehicleQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(vehicleQuery, vehicleSizeID)

	var vehicleType, vehicleSize, status sql.NullString

	var vehicleSizeId, isActive sql.NullInt64

	vehicleRes := &dtos.VehicleSizeTypeObj{}
	err := row.Scan(&vehicleSizeId, &vehicleSize, &vehicleType, &status, &isActive)
	if err != nil {
		rl.l.Error("Error GetVehicles scan: ", vehicleSizeID, err)
		return nil, err
	}

	vehicleRes.VehicleSizeId = vehicleSizeId.Int64
	vehicleRes.VehicleType = vehicleType.String
	vehicleRes.VehicleSize = vehicleSize.String
	vehicleRes.Status = status.String
	vehicleRes.IsActive = isActive.Int64

	return vehicleRes, nil
}

func (ul *TripSheetObj) GetUserLoginDetails(loginId string) (*schema.UserLogin, error) {
	userLogin := schema.UserLogin{}

	var firstName, accessToken, loginType, mobileNo sql.NullString
	var employeeId, roleID, orgId, version, isSuperAdmin, isAdmin sql.NullInt64

	query := `SELECT id,first_name,access_token,role_id,employee_id,is_super_admin,is_admin,org_id,login_type, version, mobile_no FROM user_login WHERE id = ?`
	row := ul.dbConnMSSQL.GetQueryer().QueryRow(query, loginId)
	errS := row.Scan(&userLogin.ID, &firstName, &accessToken, &roleID, &employeeId, &isSuperAdmin, &isAdmin, &orgId, &loginType, &version, &mobileNo)
	if errS != nil {
		ul.l.Error(VERIFY_ERROR, loginId, errS)
		return nil, errS
	}
	userLogin.FirstName = firstName.String
	userLogin.AccessToken = accessToken.String
	userLogin.LoginType = loginType.String
	userLogin.EmployeeId = employeeId.Int64
	userLogin.RoleID = roleID.Int64
	userLogin.OrgId = orgId.Int64
	userLogin.Version = version.Int64
	userLogin.IsSuperAdmin = int(isSuperAdmin.Int64)
	userLogin.IsAdmin = int(isAdmin.Int64)
	userLogin.MobileNo = mobileNo.String
	return &userLogin, nil
}
