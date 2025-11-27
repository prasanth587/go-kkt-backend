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

type InvoiceObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewInvoiceObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *InvoiceObj {
	return &InvoiceObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type InvoiceDao interface {
	GetTotalCount(whereQuery string) int64
	BuildWhereQuery(searchText, status, customerId string) string
	GetAllInvoices(limit int64, offset int64, whereQuery string) (*[]dtos.InvoiceObj, error)
	UpdateInvoiceStatus(invoiceId int64, status string) error
	GetInvoiceInfo(invoiceId int64) (*dtos.InvoiceObj, error)
	UpdateInvoiceNumber(invoiceId int64, invoiceNumber string, status, invoiceDate string) error
	UpdateInvoiceNumberAtTripSheet(invoiceId int64, invoiceNumber string, status string) error
	InvoiceNumberIsExistsAtTripSheet(invoiceId int64) (int64, error)
	UpdateInvoicePaid(invoiceId int64, status, invoicePaidDate, transactionId string) error
	UpdateInvoiceNumberAtTripSheetPaid(invoiceId int64, status string) error
	GetTripSheetByInvoiceId(invoiceId int64) (*[]dtos.DownloadTripSheetXls, []string, error)
	GetCountByInvoiceStatus(invoiceStatus string) int64
}

func (rl *InvoiceObj) BuildWhereQuery(searchText, status, customerId string) string {
	var conditions []string

	if status != "" {
		arr := strings.Split(status, ",")
		for i, v := range arr {
			arr[i] = fmt.Sprintf("'%s'", v)
		}
		statuses := strings.Join(arr, ",")
		conditions = append(conditions, fmt.Sprintf("ci.invoice_status IN (%v)", statuses))
	}
	if customerId != "" {
		conditions = append(conditions, fmt.Sprintf("ci.customer_id = '%v'", customerId))
	}

	if searchText != "" {
		cond := fmt.Sprintf("(ci.invoice_number LIKE '%%%v%%' OR ci.invoice_ref LIKE '%%%v%%' OR ci.work_type LIKE '%%%v%%' OR ci.trip_ref LIKE '%%%v%%')",
			searchText, searchText, searchText, searchText)
		conditions = append(conditions, cond)
	}

	whereQuery := ""
	if len(conditions) > 0 {
		whereQuery = "WHERE " + strings.Join(conditions, " AND ")
	}

	rl.l.Info("invoice whereQuery:\n ", whereQuery)

	return whereQuery
}

func (rl *InvoiceObj) GetTotalCount(whereQuery string) int64 {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM customer_invoice as ci %v`, whereQuery)
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

func (br *InvoiceObj) GetAllInvoices(limit int64, offset int64, whereQuery string) (*[]dtos.InvoiceObj, error) {
	list := []dtos.InvoiceObj{}

	whereQuery = fmt.Sprintf(" %v ORDER BY ci.updated_at DESC LIMIT %v OFFSET %v", whereQuery, limit, offset)

	customerInvoiceQuery := fmt.Sprintf(`SELECT ci.id, ci.invoice_number, ci.invoice_ref, 
	ci.work_type, ci.work_start_date, ci.work_end_date,
	ci.document_date, ci.invoice_status, ci.invoice_amount,
	ci.invoice_date, ci.payment_date, ci.customer_id, ci.customer_name, ci.customer_code, c.payment_terms FROM customer_invoice ci
    LEFT JOIN customers c ON c.customer_id = ci.customer_id %v;`, whereQuery)

	br.l.Info("customerInvoiceQuery:\n ", customerInvoiceQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(customerInvoiceQuery)
	if err != nil {
		br.l.Error("Error LoadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invoiceNumber, invoiceRef, workType, workStartDate, workEndDate, documentDate sql.NullString
		var invoiceStatus, invoiceDate, paymentDate, customerName, customerCode, paymentTerms sql.NullString
		var id, customerId sql.NullInt64
		var invoiceAmount sql.NullFloat64

		invoiceObj := &dtos.InvoiceObj{}
		err := rows.Scan(&id, &invoiceNumber, &invoiceRef, &workType, &workStartDate, &workEndDate,
			&documentDate, &invoiceStatus, &invoiceAmount, &invoiceDate, &paymentDate, &customerId, &customerName, &customerCode, &paymentTerms)
		if err != nil {
			br.l.Error("Error GetAllInvoices scan: ", err)
			return nil, err
		}
		invoiceObj.ID = id.Int64
		invoiceObj.InvoiceNumber = invoiceNumber.String
		invoiceObj.InvoiceRef = invoiceRef.String
		//invoiceObj.WorkType = workType.String
		invoiceObj.WorkStartDate = workStartDate.String
		invoiceObj.WorkEndDate = workEndDate.String
		invoiceObj.DocumentDate = documentDate.String
		invoiceObj.InvoiceStatus = invoiceStatus.String
		invoiceObj.InvoiceAmount = invoiceAmount.Float64
		invoiceObj.InvoiceDate = invoiceDate.String
		invoiceObj.PaymentDate = paymentDate.String
		invoiceObj.CustomerCode = customerCode.String
		invoiceObj.CustomerName = customerName.String
		invoiceObj.CustomerId = customerId.Int64
		invoiceObj.PaymentTerms = paymentTerms.String
		list = append(list, *invoiceObj)
	}

	return &list, nil
}

func (inv *InvoiceObj) UpdateInvoiceStatus(invoiceId int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE customer_invoice SET invoice_status = '%v' WHERE id = '%v'`, status, invoiceId)

	inv.l.Info("UpdateInvoiceStatus Update query: ", updateQuery)

	res, err := inv.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		inv.l.Error("Error db.Exec(UpdateInvoiceStatus): ", err)
		return err
	}
	inv.l.Info("UpdateInvoiceStatus Update query: ", res)
	return nil
}

func (br *InvoiceObj) GetInvoiceInfo(invoiceId int64) (*dtos.InvoiceObj, error) {

	customerInvoiceQuery := fmt.Sprintf(`SELECT id, invoice_number, invoice_ref, 
	work_type, work_start_date, work_end_date,
	document_date, invoice_status, invoice_amount,
	invoice_date,payment_date FROM customer_invoice where id = %v;`, invoiceId)

	br.l.Info("customerInvoiceQuery:\n ", customerInvoiceQuery)

	row := br.dbConnMSSQL.GetQueryer().QueryRow(customerInvoiceQuery)

	var invoiceNumber, invoiceRef, workType, workStartDate, workEndDate, documentDate sql.NullString
	var invoiceStatus, invoiceDate, paymentDate sql.NullString
	var id sql.NullInt64
	var invoiceAmount sql.NullFloat64

	invoiceObj := &dtos.InvoiceObj{}
	err := row.Scan(&id, &invoiceNumber, &invoiceRef, &workType, &workStartDate, &workEndDate,
		&documentDate, &invoiceStatus, &invoiceAmount, &invoiceDate, &paymentDate)
	if err != nil {
		br.l.Error("Error GetAllInvoice scan: ", err)
		return nil, err
	}
	invoiceObj.ID = id.Int64
	invoiceObj.InvoiceNumber = invoiceNumber.String
	invoiceObj.InvoiceRef = invoiceRef.String
	//invoiceObj.WorkType = workType.String
	invoiceObj.WorkStartDate = workStartDate.String
	invoiceObj.WorkEndDate = workEndDate.String
	invoiceObj.DocumentDate = documentDate.String
	invoiceObj.InvoiceStatus = invoiceStatus.String
	invoiceObj.InvoiceAmount = invoiceAmount.Float64
	invoiceObj.InvoiceDate = invoiceDate.String
	invoiceObj.PaymentDate = paymentDate.String

	return invoiceObj, nil
}

func (inv *InvoiceObj) UpdateInvoiceNumber(invoiceId int64, invoiceNumber, status, invoiceDate string) error {

	updateQuery := fmt.Sprintf(`UPDATE customer_invoice SET 
	 invoice_number = '%v',
	 invoice_status = '%v',
	 invoice_date = '%v'
	 WHERE id = '%v'`, invoiceNumber, status, invoiceDate, invoiceId)

	inv.l.Info("UpdateInvoiceNumber Update query: ", updateQuery)

	_, err := inv.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		inv.l.Error("Error db.Exec(UpdateInvoiceNumber): ", err)
		return err
	}

	inv.l.Info("UpdateInvoiceNumber Success: ")
	return nil
}

func (ts *InvoiceObj) InvoiceNumberIsExistsAtTripSheet(invoiceId int64) (int64, error) {

	query := fmt.Sprintf("SELECT COUNT(*) FROM trip_sheet_header WHERE customer_invoice_id = '%v'", invoiceId)
	ts.l.Info("InvoiceNumberIsExistsAtTripSheet query: ", query)

	var tripCount sql.NullInt64
	err := ts.dbConnMSSQL.GetQueryer().QueryRow(query).Scan(&tripCount)
	if err != nil {
		ts.l.Error("Error db.Exec(IsTripExists): ", err)
		return 0, err
	}
	ts.l.Info("found total Invoice ID at trip_sheet_header ", tripCount.Int64)
	return tripCount.Int64, nil
}

func (inv *InvoiceObj) UpdateInvoiceNumberAtTripSheet(invoiceId int64, invoiceNumber string, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET 
	 customer_invoice_no = '%v',
	 load_status = '%v'
	 WHERE customer_invoice_id = '%v'`, invoiceNumber, status, invoiceId)

	inv.l.Info("UpdateInvoiceNumberAtTripSheet Update query: ", updateQuery)

	_, err := inv.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		inv.l.Error("Error db.Exec(UpdateInvoiceNumberAtTripSheet): ", err)
		return err
	}

	return nil
}

func (inv *InvoiceObj) UpdateInvoicePaid(invoiceId int64, status, invoicePaidDate, transactionId string) error {

	updateQuery := fmt.Sprintf(`UPDATE customer_invoice SET invoice_status = '%v', payment_date = '%v', transaction_id = '%v' WHERE id = '%v'`, status, invoicePaidDate, transactionId, invoiceId)

	inv.l.Info("UpdateInvoicePaid Update query: ", updateQuery)

	_, err := inv.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		inv.l.Error("Error db.Exec(UpdateInvoicePaid): ", err)
		return err
	}

	inv.l.Info("UpdateInvoicePaid Success: ")
	return nil
}

func (inv *InvoiceObj) UpdateInvoiceNumberAtTripSheetPaid(invoiceId int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET load_status = '%v' WHERE customer_invoice_id = '%v'`, status, invoiceId)

	inv.l.Info("UpdateInvoiceNumberAtTripSheetPaid Update query: ", updateQuery)

	_, err := inv.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		inv.l.Error("Error db.Exec(UpdateInvoiceNumberAtTripSheetPaid): ", err)
		return err
	}

	return nil
}

func (ts *InvoiceObj) GetTripSheetByInvoiceId(invoiceId int64) (*[]dtos.DownloadTripSheetXls, []string, error) {
	list := []dtos.DownloadTripSheetXls{}
	var tripSheetIdPesent []string
	whereQuery := fmt.Sprintf(" WHERE t.customer_invoice_id = '%v' ORDER BY t.trip_sheet_id ASC;", invoiceId)

	query := fmt.Sprintf(`
    SELECT 
        t.trip_sheet_id, t.trip_sheet_num, t.trip_type, t.trip_sheet_type, 
        t.load_hours_type, t.open_trip_date_time, t.branch_id, t.customer_id, 
        c.customer_name, c.customer_code,
        t.vendor_id, 
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
	LEFT JOIN vehicle_size_type vs ON t.vehicle_size = vs.vehicle_size_id 
    %v`, whereQuery)

	ts.l.Info("GetTripSheetByInvoiceId query:\n ", query)
	rows, err := ts.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		ts.l.Error("Error GetTripSheetByInvoiceId ", err)
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

func (rl *InvoiceObj) GetCountByInvoiceStatus(invoiceStatus string) int64 {
	countQuery := fmt.Sprintf(`Select count(*) from customer_invoice where invoice_status =  '%v'`, invoiceStatus)
	rl.l.Info(" GetCountByInvoiceStatus select query: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCountByInvoiceStatus scan: ", errE)
		return 0
	}

	return count.Int64
}

func (rl *InvoiceObj) GetInvoices(invoiceStatus string) int64 {
	countQuery := fmt.Sprintf(`Select count(*) from customer_invoice where invoice_status =  '%v'`, invoiceStatus)
	rl.l.Info(" GetCountByInvoiceStatus select query: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCountByInvoiceStatus scan: ", errE)
		return 0
	}

	return count.Int64
}
