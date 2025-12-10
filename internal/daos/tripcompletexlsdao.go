package daos

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type TripSheetXls struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewTripSheetXls(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *TripSheetXls {
	return &TripSheetXls{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type TripSheetXlsDao interface {
	BuildWhereQuery(orgId int64, tripSheetid, tripStatus, tripSearchText, fromDate, toDate, podRequired, podReceived, customers, vendors string) string
	GetTripSheets(orgId int64, whereQuery string, limit, offset int) (*[]dtos.TripSheetXls, []string, error)
	GetTotalCount(whereQuery string) int64
	GetTripSheetLoadingPoint(tripSheetID int64, lp_up_type string) (*[]dtos.TripSheetLoading, error)
	GetTripSheetUnLoadingPoint(tripSheetID int64, lp_up_type string) (*[]dtos.TripSheetUnLoading, error)
	GetLoadingUnloadingPoints(tripSheetIds string) (*[]dtos.LoadUnLoadObj, error)

	GetTripSheetByIds(orgId int64, tripSheetIds string) (*[]dtos.DownloadTripSheetXls, []string, error)
	XLSUpdateTripSheetHeader(tripSheetID int64, tripSheetUpdateReq dtos.TripSheetUpdateData) error
	IsTripExists(tripSheetID int64) bool
	ISTripExistsInPODManager(tripSheetID int64) bool
	GetTripSheetsV1Sort(orgId int64, whereQuery string, limit, offset int) (*[]dtos.TripSheetXlsV1, []string, error)
	CheckSelectedTripsAreSingleCustomers(tripSheetIds string) ([]int64, error)
	IsTripSheetNumberExists(tripSheetNum string) bool
	XLSUpdateTripSheetHeaderForDraftInvoice(tripSheetNum string, tripSheetUpdateReq dtos.UpdateDraftInvoiceData) error
	GetTripsByTripSheetNumbers(tripSheetNumbers string) (*[]dtos.TripSheetDraftPull, error)
	CreateInvoiceDraft(draft dtos.CreateDraftInvoice) (int64, error)
	UpdateInvoiceIDToTripHeader(invoiceRowId int64, inClauseTripID string) error
	GetCustomerIDAndName(tripSheetID int64) (string, string, int64)
}

func (rl *TripSheetXls) BuildWhereQuery(orgId int64, tripSheetid, tripStatus, tripSearchText, fromDate, toDate, podRequired, podReceived, customers, vendors string) string {

	whereQuery := fmt.Sprintf("WHERE t.org_id = '%v'", orgId)

	if tripSheetid != "" {
		whereQuery = fmt.Sprintf(" %v AND t.trip_sheet_id = '%v'", whereQuery, tripSheetid)
	}
	if podRequired != "" {
		whereQuery = fmt.Sprintf(" %v AND t.pod_required = '%v'", whereQuery, podRequired)
	}
	if podReceived != "" {
		whereQuery = fmt.Sprintf(" %v AND t.pod_received = '%v'", whereQuery, podReceived)
	}

	if customers != "" {
		whereQuery = fmt.Sprintf(" %v AND t.customer_id IN (%s)", whereQuery, customers)
	}
	if vendors != "" {
		whereQuery = fmt.Sprintf(" %v AND t.vendor_id IN (%s)", whereQuery, vendors)
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
		whereQuery = fmt.Sprintf(" %v AND t.load_status IN (%s)", whereQuery, strings.Join(quoted, ", "))
	}

	// if fromDate != "" && toDate != "" {
	// 	whereQuery = fmt.Sprintf(" %v AND (t.open_trip_date_time BETWEEN '%v' AND '%v' )", whereQuery, fromDate, toDate)
	// }

	if fromDate != "" && toDate != "" {
		whereQuery = fmt.Sprintf(" %v AND ((t.open_trip_date_time BETWEEN '%v' AND '%v') OR (t.updated_at BETWEEN '%v' AND '%v'))", whereQuery, fromDate, toDate, fromDate, toDate)
	}

	if tripSearchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (t.trip_sheet_num LIKE '%%%v%%' OR t.trip_sheet_type LIKE '%%%v%%' OR t.load_hours_type LIKE '%%%v%%' OR t.customer_invoice_no LIKE '%%%v%%' OR t.zonal_name LIKE '%%%v%%' OR t.vehicle_number LIKE '%%%v%%' OR t.vehicle_size LIKE '%%%v%%' OR t.mobile_number LIKE '%%%v%%' OR t.load_status LIKE '%%%v%%') ", whereQuery, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText, tripSearchText)
	}

	rl.l.Info("tripSheet whereQuery:\n ", whereQuery)

	return whereQuery
}

func (rl *TripSheetXls) GetTotalCount(whereQuery string) int64 {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM trip_sheet_header as t %v`, whereQuery)
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

func (ts *TripSheetXls) GetTripSheets(orgId int64, whereQuery string, limit, offset int) (*[]dtos.TripSheetXls, []string, error) {
	list := []dtos.TripSheetXls{}
	var tripSheetIds []string
	whereQuery = fmt.Sprintf(" %v ORDER BY t.updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)
	// SQL query to fetch trip sheets
	query := fmt.Sprintf(`
    SELECT 
        t.trip_sheet_id, t.trip_sheet_num, t.trip_type, t.trip_sheet_type, 
        t.load_hours_type, t.open_trip_date_time, t.branch_id, t.customer_id, 
        c.customer_name, c.customer_code,
        t.vendor_id, 
        t.vehicle_capacity_ton, t.vehicle_number, t.vehicle_size, 
        t.mobile_number, t.driver_name, t.driver_license_image, t.lr_gate_image, t.lr_number,
        t.customer_base_rate, t.customer_km_cost, t.customer_toll, t.customer_extra_hours,
        t.customer_extra_km, t.customer_total_hire, t.customer_close_trip_date_time,
        t.customer_invoice_no, t.customer_payment_received_date, t.customer_debit_amount,
        t.customer_remark, t.customer_billing_raised_date, t.pod_required,
        t.vendor_base_rate, t.vendor_km_cost, t.vendor_toll, t.vendor_total_hire,
        t.vendor_advance, t.vendor_paid_date, t.vendor_cross_dock, t.vendor_remark,
        t.vendor_debit_amount, t.vendor_balance_amount, t.pod_received,
        t.customer_per_load_hire, t.customer_running_km, t.customer_per_km_price, t.customer_placed_vehicle_size,
        t.customer_load_cancelled, t.vendor_paid_by, t.vendor_load_unload_amount, t.vendor_halting_days,
        t.vendor_halting_paid, t.vendor_extra_delivery, t.load_status, t.customer_reported_date_time_for_halting_calc, t.vendor_break_down, t.zonal_name
    FROM trip_sheet_header t
    LEFT JOIN customers c ON t.customer_id = c.customer_id
    %v`, whereQuery)

	ts.l.Info("TripSheet Query:\n", query)

	// Execute the query
	rows, err := ts.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		ts.l.Error("Error fetching trip sheets: ", err)
		return nil, tripSheetIds, err
	}
	defer rows.Close()

	// Iterate through the results
	for rows.Next() {
		var (
			// Core trip details
			tripSheetNum, tripType, tripSheetType, loadHoursType, openTripDateTime, vehicleCapacityTon, vehicleNumber, vehicleSize, mobileNumber, driverName,
			driverLicenseImage, lrGateImage, lrNumber, customerCloseTripDT, customerInvoiceNo,
			customerPaymentDate, customerRemark, customerBillingDate, vendorPaidDate, vendorCrossDock, vendorRemark,
			customerPlacedVehicleSize, customerLoadCancelled, vendorPaidBy, loadStatus, customerCode, customerName,
			customerReportedDateTimeForHaltingCalc, zonalName, vendorBreakDown sql.NullString

			tripSheetID, branchID, customerID, vendorID, podRequired, podReceived sql.NullInt64
			customerBaseRate, customerKmCost, customerToll, customerExtraHours, customerExtraKm, customerTotalHire, customerDebitAmount,
			vendorBaseRate, vendorKmCost, vendorToll, vendorTotalHire, vendorAdvance, vendorDebitAmount, vendorBalanceAmount,
			customerPerLoadHire, customerRunningKM, customerPerKMPrice, vendorLoadUnLoadAmount, vendorHaltingDays,
			vendorHaltingPaid, vendorExtraDelivery sql.NullFloat64
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
			&customerName,
			&customerCode,
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
		)

		if err != nil {
			ts.l.Error("Error scanning trip sheet row: ", err)
			return nil, tripSheetIds, err
		}

		// Create a TripSheet instance and populate it with data from the row.
		trip := dtos.TripSheetXls{
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

			CustomerId:   customerID.Int64,
			CustomerName: customerName.String,
			CustomerCode: customerCode.String,
		}

		list = append(list, trip)

		if tripSheetID.Valid {
			tripSheetIdStr := strconv.FormatInt(tripSheetID.Int64, 10)
			tripSheetIds = append(tripSheetIds, tripSheetIdStr)
		}
	}

	return &list, tripSheetIds, nil // Return the list of TripSheets.
}

func (ts *TripSheetXls) GetTripSheetsV1Sort(orgId int64, whereQuery string, limit, offset int) (*[]dtos.TripSheetXlsV1, []string, error) {
	list := []dtos.TripSheetXlsV1{}
	var tripSheetIds []string
	whereQuery = fmt.Sprintf(" %v ORDER BY t.updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)
	// SQL query to fetch trip sheets
	query := fmt.Sprintf(`
    SELECT 
        t.trip_sheet_id, t.trip_sheet_num, t.trip_type, t.trip_sheet_type, 
        t.load_hours_type, t.open_trip_date_time, t.branch_id, t.customer_id, 
        c.customer_name, c.customer_code, t.vendor_id, t.lr_number, t.customer_invoice_no, t.load_status, t.vehicle_number, v.vendor_name, v.vendor_code
    FROM trip_sheet_header t
    LEFT JOIN customers c ON t.customer_id = c.customer_id
	LEFT JOIN vendors v ON t.vendor_id = v.vendor_id
    %v`, whereQuery)

	ts.l.Info("TripSheet Query:\n", query)

	// Execute the query
	rows, err := ts.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		ts.l.Error("Error fetching trip sheets: ", err)
		return nil, tripSheetIds, err
	}
	defer rows.Close()

	// Iterate through the results
	for rows.Next() {
		var (
			// Core trip details
			tripSheetNum, tripType, tripSheetType, loadHoursType, openTripDateTime,
			customerName, customerCode, lrNumber, customerInvoiceNo, loadStatus, vehicleNumber, vendorName, vendorCode sql.NullString
			tripSheetID, branchID, customerID, vendorID sql.NullInt64
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
			&customerName,
			&customerCode,
			&vendorID,
			&lrNumber,
			&customerInvoiceNo,
			&loadStatus,
			&vehicleNumber,
			&vendorName,
			&vendorCode,
		)

		if err != nil {
			ts.l.Error("Error scanning trip sheet row: ", err)
			return nil, tripSheetIds, err
		}

		// Create a TripSheet instance and populate it with data from the row.
		trip := dtos.TripSheetXlsV1{
			TripSheetID:       tripSheetID.Int64,
			TripSheetNum:      tripSheetNum.String,
			TripType:          tripType.String,
			TripSheetType:     tripSheetType.String,
			LoadHoursType:     loadHoursType.String,
			OpenTripDateTime:  openTripDateTime.String,
			BranchID:          branchID.Int64,
			VendorID:          vendorID.Int64,
			LRNumber:          lrNumber.String,
			CustomerInvoiceNo: customerInvoiceNo.String,
			CustomerId:        customerID.Int64,
			CustomerName:      customerName.String,
			CustomerCode:      customerCode.String,
			LoadStatus:        loadStatus.String,
			VehicleNumber:     vehicleNumber.String,
			VendorName:        vendorName.String + " - " + vendorCode.String,
		}

		list = append(list, trip)

		if tripSheetID.Valid {
			tripSheetIdStr := strconv.FormatInt(tripSheetID.Int64, 10)
			tripSheetIds = append(tripSheetIds, tripSheetIdStr)
		}
	}

	return &list, tripSheetIds, nil // Return the list of TripSheets.
}

func (br *TripSheetXls) GetTripSheetLoadingPoint(tripSheetID int64, lp_up_type string) (*[]dtos.TripSheetLoading, error) {
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

func (br *TripSheetXls) GetTripSheetUnLoadingPoint(tripSheetID int64, lp_up_type string) (*[]dtos.TripSheetUnLoading, error) {
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

func (mp *TripSheetXls) GetLoadingUnloadingPoints(tripSheetIds string) (*[]dtos.LoadUnLoadObj, error) {
	list := []dtos.LoadUnLoadObj{}
	loadUnLoadQuery := fmt.Sprintf(`SELECT trip_sheet_id, type, loading_point_id FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id IN (%v)  ORDER BY load_unload_point_id ASC; `, tripSheetIds)

	mp.l.Info("GetLoadingUnloadingPoints whereQuery:\n ", loadUnLoadQuery)

	rows, err := mp.dbConnMSSQL.GetQueryer().Query(loadUnLoadQuery)
	if err != nil {
		mp.l.Error("Error GetLoadingUnloadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		loadUnLoad := &dtos.LoadUnLoadObj{}

		var tripSheetId, loadUnloadPointId sql.NullInt64
		var typeT sql.NullString
		err := rows.Scan(&tripSheetId, &typeT, &loadUnloadPointId)
		if err != nil {
			mp.l.Error("Error GetManagePods scan: ", err)
			return nil, err
		}

		loadUnLoad.TripSheetId = tripSheetId.Int64
		loadUnLoad.Type = typeT.String
		loadUnLoad.LoadingPointId = loadUnloadPointId.Int64

		list = append(list, *loadUnLoad)
	}
	return &list, nil
}

func (ts *TripSheetXls) GetTripSheetByIds(orgId int64, tripSheetIds string) (*[]dtos.DownloadTripSheetXls, []string, error) {
	list := []dtos.DownloadTripSheetXls{}
	var tripSheetIdPesent []string
	whereQuery := fmt.Sprintf(" WHERE t.trip_sheet_id IN (%v) ORDER BY t.trip_sheet_id ASC;", tripSheetIds)

	query := fmt.Sprintf(`
    SELECT 
        t.trip_sheet_id, t.trip_sheet_num, t.trip_type, t.trip_sheet_type, 
        t.load_hours_type, t.open_trip_date_time, t.branch_id, t.customer_id, 
        c.customer_name, c.customer_code,
        t.vendor_id, v.vendor_name, v.vendor_code,
        t.vehicle_capacity_ton, t.vehicle_number, CONCAT(vs.vehicle_size, ' (', vs.vehicle_type, ')') AS vehicle_size_type, 
        t.mobile_number, t.driver_name, t.driver_license_image, t.lr_gate_image, t.lr_number,
        t.customer_base_rate, t.customer_km_cost, t.customer_toll, t.customer_extra_hours,
        t.customer_extra_km, t.customer_total_hire, t.customer_close_trip_date_time,
        t.customer_invoice_no, t.customer_payment_received_date, t.customer_debit_amount,
        t.customer_remark, t.customer_billing_raised_date, t.pod_required,
        t.vendor_base_rate, t.vendor_km_cost, t.vendor_toll, t.vendor_total_hire,
        t.vendor_advance, t.vendor_paid_date, t.vendor_cross_dock, t.vendor_remark,
        t.vendor_debit_amount, t.vendor_balance_amount, t.pod_received,
        t.customer_per_load_hire, t.customer_running_km, t.customer_per_km_price, t.customer_placed_vehicle_size,
        t.customer_load_cancelled, t.vendor_paid_by, t.vendor_load_unload_amount, t.vendor_halting_days,
        t.vendor_halting_paid, t.vendor_extra_delivery, t.load_status, t.customer_reported_date_time_for_halting_calc, t.vendor_break_down, t.zonal_name
    FROM trip_sheet_header t
    LEFT JOIN customers c ON t.customer_id = c.customer_id
    LEFT JOIN vendors v ON t.vendor_id = v.vendor_id
	LEFT JOIN vehicle_size_type vs ON t.vehicle_size = vs.vehicle_size_id 
    %v`, whereQuery)

	ts.l.Info("GetTripSheetByIds query:\n ", query)
	rows, err := ts.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		ts.l.Error("Error GetTripSheetByIds ", err)
		return nil, tripSheetIdPesent, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			// Core trip details
			tripSheetNum, tripType, tripSheetType, loadHoursType, openTripDateTime, vehicleCapacityTon, vehicleNumber, vehicleSize, mobileNumber, driverName,
			driverLicenseImage, lrGateImage, lrNumber, customerCloseTripDT, customerInvoiceNo,
			customerPaymentDate, customerRemark, customerBillingDate, vendorPaidDate, vendorCrossDock, vendorRemark,
			customerPlacedVehicleSize, customerLoadCancelled, vendorPaidBy, loadStatus, customerCode, customerName,
			customerReportedDateTimeForHaltingCalc, zonalName, vendorBreakDown, vendorName, vendorCode sql.NullString

			tripSheetID, branchID, customerID, vendorID, podRequired, podReceived sql.NullInt64
			customerBaseRate, customerKmCost, customerToll, customerExtraHours, customerExtraKm, customerTotalHire, customerDebitAmount,
			vendorBaseRate, vendorKmCost, vendorToll, vendorTotalHire, vendorAdvance, vendorDebitAmount, vendorBalanceAmount,
			customerPerLoadHire, customerRunningKM, customerPerKMPrice, vendorLoadUnLoadAmount, vendorHaltingDays,
			vendorHaltingPaid, vendorExtraDelivery sql.NullFloat64
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
			&customerName,
			&customerCode,
			&vendorID,
			&vendorName,
			&vendorCode,
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
		)

		if err != nil {
			ts.l.Error("Error scanning trip sheet row: ", err)
			return nil, tripSheetIdPesent, err
		}

		vendorBreakDownF := 0.0
		if vendorBreakDown.String != "" {
			vendorBreakDownF, _ = strconv.ParseFloat(vendorBreakDown.String, 64)
		}

		// Create a TripSheet instance and populate it with data from the row.
		trip := dtos.DownloadTripSheetXls{
			TripSheetID:                            tripSheetID.Int64,
			TripSheetNum:                           tripSheetNum.String,
			TripType:                               tripType.String,
			TripSheetType:                          tripSheetType.String,
			LoadHoursType:                          loadHoursType.String,
			OpenTripDateTime:                       openTripDateTime.String,
			VendorID:                               vendorID.Int64,
			VendorName:                             vendorName.String,
			VendorCode:                             vendorCode.String,
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
			VendorBreakDown:                        vendorBreakDownF,
			ZonalName:                              zonalName.String,

			CustomerName: customerName.String,
			CustomerCode: customerCode.String,
		}

		list = append(list, trip)

		if tripSheetID.Valid {
			tripSheetIdStr := strconv.FormatInt(tripSheetID.Int64, 10)
			tripSheetIdPesent = append(tripSheetIdPesent, tripSheetIdStr)
		}
	}
	return &list, tripSheetIdPesent, nil
}

func (ts *TripSheetXls) XLSUpdateTripSheetHeader(tripSheetID int64, tripSheetUpdateReq dtos.TripSheetUpdateData) error {

	updateQuery := fmt.Sprintf(
		`UPDATE trip_sheet_header SET
        customer_invoice_no = '%v',
        customer_close_trip_date_time = '%v',
        customer_payment_received_date = '%v',
        customer_remark = '%v',
        customer_billing_raised_date = '%v',
        customer_base_rate = %v,
        customer_km_cost = %v,
        customer_toll = %v,
        customer_extra_hours = %v,
        customer_extra_km = %v,
        customer_total_hire = %v,
        customer_debit_amount = %v,
        customer_per_load_hire = %v,
        customer_running_km = %v,
        customer_per_km_price = %v,
        customer_placed_vehicle_size = '%v',
        customer_load_cancelled = '%v',
        customer_reported_date_time_for_halting_calc = '%v',
		load_status = '%v'
    WHERE trip_sheet_id = '%v'`,
		tripSheetUpdateReq.CustomerInvoiceNo,
		tripSheetUpdateReq.CustomerCloseTripDateTime,
		tripSheetUpdateReq.CustomerPaymentReceivedDate,
		tripSheetUpdateReq.CustomerRemark,
		tripSheetUpdateReq.CustomerBillingRaisedDate,
		tripSheetUpdateReq.CustomerBaseRate,
		tripSheetUpdateReq.CustomerKMCost,
		tripSheetUpdateReq.CustomerToll,
		tripSheetUpdateReq.CustomerExtraHours,
		tripSheetUpdateReq.CustomerExtraKM,
		tripSheetUpdateReq.CustomerTotalHire,
		tripSheetUpdateReq.CustomerDebitAmount,
		tripSheetUpdateReq.CustomerPerLoadHire,
		tripSheetUpdateReq.CustomerRunningKM,
		tripSheetUpdateReq.CustomerPerKMPrice,
		tripSheetUpdateReq.CustomerPlacedVehicleSize,
		tripSheetUpdateReq.CustomerLoadCancelled,
		tripSheetUpdateReq.CustomerReportedDateTimeForHaltingCalc,
		tripSheetUpdateReq.LoadStatus,
		tripSheetID,
	)

	ts.l.Info("XLSUpdateTripSheetHeader Update query: ", updateQuery)

	_, err := ts.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		ts.l.Error("Error db.Exec(XLSUpdateTripSheetHeader): ", err)
		return err
	}
	ts.l.Info("trip sheet updated successfully: ", tripSheetID, tripSheetUpdateReq.CustomerInvoiceNo)

	return nil
}

func (ts *TripSheetXls) IsTripExists(tripSheetID int64) bool {

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM trip_sheet_header WHERE trip_sheet_id = '%v')", tripSheetID)
	ts.l.Debug("IsTripExists query: ", query)

	var exists bool
	err := ts.dbConnMSSQL.GetQueryer().QueryRow(query).Scan(&exists)
	if err != nil {
		ts.l.Error("Error db.Exec(IsTripExists): ", err)
		return false
	}
	if !exists {
		ts.l.Error("IsTripExists Not found ", tripSheetID)
	}
	return exists
}

func (ts *TripSheetXls) GetCustomerIDAndName(tripSheetID int64) (string, string, int64) {

	query := fmt.Sprintf(`select t.customer_id, c.customer_name, c.customer_code from trip_sheet_header as t LEFT JOIN customers c ON t.customer_id = c.customer_id where trip_sheet_id = '%v'`, tripSheetID)
	ts.l.Debug("IsTripExists query: ", query)

	var customerName, customerCode sql.NullString
	var customerId sql.NullInt64
	err := ts.dbConnMSSQL.GetQueryer().QueryRow(query).Scan(&customerId, &customerName, &customerCode)
	if err != nil {
		ts.l.Error("Error db.Exec(GetCustomerIDAndName): ", err)
		return "", "", 0
	}

	return customerName.String, customerCode.String, customerId.Int64
}

func (ts *TripSheetXls) ISTripExistsInPODManager(tripSheetID int64) bool {

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM manage_pod WHERE trip_sheet_id = '%v')", tripSheetID)
	ts.l.Debug("ISTripExistsInPODManager query: ", query)

	var exists bool
	err := ts.dbConnMSSQL.GetQueryer().QueryRow(query).Scan(&exists)
	if err != nil {
		ts.l.Error("Error db.Exec(ISTripExistsInPODManager): ", err)
		return false
	}
	if !exists {
		ts.l.Error("ISTripExistsInPODManager Not found ", tripSheetID)
	}
	return exists
}

func (ts *TripSheetXls) CheckSelectedTripsAreSingleCustomers(tripSheetIds string) ([]int64, error) {
	var customerIds []int64
	query := fmt.Sprintf("SELECT customer_id from trip_sheet_header where trip_sheet_id IN (%s)", tripSheetIds)

	ts.l.Info("CheckSelectedTripsAreSingleCustomers: query:\n ", query)

	rows, err := ts.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		ts.l.Error("Error CheckSelectedTripsAreSingleCustomers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customerId sql.NullInt64
		err := rows.Scan(&customerId)
		if err != nil {
			ts.l.Error("Error CheckSelectedTripsAreSingleCustomers scan: ", err)
			return nil, err
		}
		customerIds = append(customerIds, customerId.Int64)
	}
	return customerIds, nil
}

func (ts *TripSheetXls) IsTripSheetNumberExists(tripSheetNum string) bool {

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM trip_sheet_header WHERE trip_sheet_num = '%v')", tripSheetNum)
	ts.l.Debug("IsTripExists query: ", query)

	var exists bool
	err := ts.dbConnMSSQL.GetQueryer().QueryRow(query).Scan(&exists)
	if err != nil {
		ts.l.Error("Error db.Exec(IsTripExists): ", err)
		return false
	}
	if !exists {
		ts.l.Error("IsTripExists Not found ", tripSheetNum)
	}
	return exists
}

func (ts *TripSheetXls) XLSUpdateTripSheetHeaderForDraftInvoice(tripSheetNum string, tripSheetUpdateReq dtos.UpdateDraftInvoiceData) error {

	updateQuery := fmt.Sprintf(
		`UPDATE trip_sheet_header SET
        customer_per_load_hire = %v,
		customer_running_km = %v,
		customer_base_rate = %v,
		customer_per_km_price = %v,
		customer_km_cost = %v,
		customer_toll = %v,
		customer_extra_hours = %v,
        customer_extra_km = %v,
        customer_total_hire = %v,
        customer_close_trip_date_time = '%v',
        customer_payment_received_date = '%v',
		customer_debit_amount = %v,
		customer_billing_raised_date = '%v',
		customer_load_cancelled = '%v',
		customer_reported_date_time_for_halting_calc = '%v',
        customer_remark = '%v'
    WHERE trip_sheet_num = '%v'`,

		tripSheetUpdateReq.CustomerPerLoadHire,
		tripSheetUpdateReq.CustomerRunningKM,
		tripSheetUpdateReq.CustomerBaseRate,
		tripSheetUpdateReq.CustomerPerKMPrice,
		tripSheetUpdateReq.CustomerKMCost,
		tripSheetUpdateReq.CustomerToll,
		tripSheetUpdateReq.CustomerExtraHours,
		tripSheetUpdateReq.CustomerExtraKM,
		tripSheetUpdateReq.CustomerTotalHire,
		tripSheetUpdateReq.CustomerCloseTripDateTime,
		tripSheetUpdateReq.CustomerPaymentReceivedDate,
		tripSheetUpdateReq.CustomerDebitAmount,
		tripSheetUpdateReq.CustomerBillingRaisedDate,
		tripSheetUpdateReq.CustomerLoadCancelled,
		tripSheetUpdateReq.CustomerReportedDateTimeForHaltingCalc,
		tripSheetUpdateReq.CustomerRemark,
		tripSheetNum,
	)

	ts.l.Info("XLSUpdateTripSheetHeaderForDraftInvoice Update query: ", updateQuery)

	_, err := ts.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		ts.l.Error("Error db.Exec(XLSUpdateTripSheetHeader): ", err)
		return err
	}
	ts.l.Info("trip sheet updated successfully (DraftInvoice): ", tripSheetNum)
	return nil
}

func (ts *TripSheetXls) GetTripsByTripSheetNumbers(tripSheetNumbers string) (*[]dtos.TripSheetDraftPull, error) {

	list := []dtos.TripSheetDraftPull{}
	loadUnLoadQuery := fmt.Sprintf(`SELECT trip_sheet_id, trip_sheet_num, customer_per_load_hire, customer_running_km, 
    customer_per_km_price,customer_base_rate, trip_type, customer_km_cost, customer_toll  from  trip_sheet_header WHERE trip_sheet_num IN %s; `, tripSheetNumbers)

	ts.l.Info("GetTripsByTripSheetNumbers query:\n ", loadUnLoadQuery)

	rows, err := ts.dbConnMSSQL.GetQueryer().Query(loadUnLoadQuery)
	if err != nil {
		ts.l.Error("Error GetTripsByTripSheetNumbers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tripSheet := &dtos.TripSheetDraftPull{}

		var tripSheetId sql.NullInt64
		var tripSheetNum, tripType sql.NullString
		var customerPerLoadHire, customerRunningKM, customerPerKMPrice, customerBaseRate, customerKmCost, customerToll sql.NullFloat64

		err := rows.Scan(&tripSheetId, &tripSheetNum, &customerPerLoadHire, &customerRunningKM, &customerPerKMPrice, &customerBaseRate, &tripType, &customerKmCost, &customerToll)
		if err != nil {
			ts.l.Error("Error GetTripsByTripSheetNumbers scan: ", err)
			return nil, err
		}

		tripSheet.TripSheetID = tripSheetId.Int64
		tripSheet.TripSheetNum = tripSheetNum.String
		tripSheet.CustomerPerLoadHire = customerPerLoadHire.Float64
		tripSheet.CustomerRunningKM = customerRunningKM.Float64
		tripSheet.CustomerPerKMPrice = customerPerKMPrice.Float64
		tripSheet.CustomerBaseRate = customerBaseRate.Float64
		tripSheet.TripType = tripType.String
		tripSheet.CustomerKMCost = customerKmCost.Float64
		tripSheet.CustomerToll = customerToll.Float64

		list = append(list, *tripSheet)
	}
	return &list, nil

}

func (ts *TripSheetXls) CreateInvoiceDraft(draft dtos.CreateDraftInvoice) (int64, error) {

	ts.l.Info("CreateInvoinceDraft : ", draft.InvoiceAmount)
	//customer_id, c.customer_name, c.customer_code
	invoiceQuery := fmt.Sprintf(`
		INSERT INTO customer_invoice (
			invoice_ref, work_type, work_start_date, 
			work_end_date, document_date, invoice_amount, 
			trip_ref, invoice_status, customer_id, customer_name, customer_code
		) VALUES (
			'%v', '%v', '%v',
			'%v', '%v', '%v',
			'%v', '%v','%v', '%v','%v'
		)`,
		draft.InvoiceRefID, draft.WorkType, draft.WorkStartDate,
		draft.WorkEndDate, draft.DocumentDate, draft.InvoiceAmount,
		draft.TripRef, draft.InvoiceStatus, draft.CustomerId, draft.CustomerName, draft.CustomerCode,
	)

	ts.l.Info("invoiceQuery : ", invoiceQuery)

	roleResult, err := ts.dbConnMSSQL.GetQueryer().Exec(invoiceQuery)
	if err != nil {
		ts.l.Error("Error db.Exec(CreateInvoinceDraft): ", err)
		return 0, err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		ts.l.Error("Error db.Exec(CreateInvoinceDraft):", createdId, err)
		return 0, err
	}
	ts.l.Info("Invoince Draft created successfully", createdId, draft.InvoiceAmount)
	return createdId, nil
}

//

func (ts *TripSheetXls) UpdateInvoiceIDToTripHeader(invoiceRowId int64, inClauseTripID string) error {

	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET customer_invoice_id = '%v' WHERE trip_sheet_id IN (%s)`, invoiceRowId, inClauseTripID)

	ts.l.Info("UpdateInvoiceIDToTripHeader Update query ", updateQuery)

	_, err := ts.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		ts.l.Error("Error db.Exec(UpdateInvoiceIDToTripHeader): ", err)
		return err
	}

	ts.l.Info("trips updated successfully: ", invoiceRowId)

	return nil
}
