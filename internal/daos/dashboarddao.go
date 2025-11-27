package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type DashBoardObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewDashBoardObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *DashBoardObj {
	return &DashBoardObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

var (
// UPDATE_ERROR = "ERROR UpdateUserLoginAterCredentialSuccess update: %v "
// VERIFY_ERROR = "ERROR VerifyCredentials: %v "
// INVALIDE_    = "invalide credentials"
// LOGIN_FAILED = "login failed  - "
)

type DashBoardDao interface {
	GetEmpolyee() (*dtos.Employee, error)
	GetLoadingUnLoadingPoints() (*dtos.LoadUnloadPoints, error)
	GetCustomer() (*dtos.CustomerStats, error)
	GetVendor() (*dtos.VendorStats, error)
	GetVhicle() (*dtos.VehicleStats, error)
	GetTripStatusCounts(tripStatus, start, end, column string) (*dtos.TripStatusStats, error)
	GetTripGraph(startDateTime, endDateTime string) (*[]dtos.TripCountWithType, error)
	GetCountByInvoiceStatus(invoiceStatus string) int64
	GetAllInvoicesWithoutLimit(status string) (*[]dtos.InvoiceBreach, error)
}

func (rl *DashBoardObj) GetVhicle() (*dtos.VehicleStats, error) {
	countQuery := "SELECT count(*) FROM vehicles;;"
	rl.l.Info("GetVhicle: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var inActive, totalCount sql.NullInt64

	errE := row.Scan(&totalCount)
	if errE != nil {
		rl.l.Error("Error GetVhicle scan: ", errE)
		return nil, errE
	}
	model := dtos.VehicleStats{
		ActiveCount:   totalCount.Int64,
		InActiveCount: inActive.Int64,
		TotalCount:    totalCount.Int64,
	}

	return &model, nil
}

func (rl *DashBoardObj) GetCustomer() (*dtos.CustomerStats, error) {

	countQuery := "SELECT SUM(is_active = '1') AS active_count, SUM(is_active = '0') AS in_active_count, COUNT(*) AS total_count  FROM customers;"
	rl.l.Info("GetEmpolyee: ", countQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var active, inActive, totalCount sql.NullInt64

	errE := row.Scan(&active, &inActive, &totalCount)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	model := dtos.CustomerStats{
		ActiveCount:   active.Int64,
		InActiveCount: inActive.Int64,
		TotalCount:    totalCount.Int64,
	}

	return &model, nil
}

func (rl *DashBoardObj) GetVendor() (*dtos.VendorStats, error) {

	countQuery := "SELECT SUM(is_active = '1') AS active_count, SUM(is_active = '0') AS in_active_count, COUNT(*) AS total_count  FROM vendors;"
	rl.l.Info("GetEmpolyee: ", countQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var active, inActive, totalCount sql.NullInt64

	errE := row.Scan(&active, &inActive, &totalCount)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	model := dtos.VendorStats{
		ActiveCount:   active.Int64,
		InActiveCount: inActive.Int64,
		TotalCount:    totalCount.Int64,
	}

	return &model, nil
}

func (rl *DashBoardObj) GetLoadingUnLoadingPoints() (*dtos.LoadUnloadPoints, error) {

	countQuery := "SELECT SUM(is_active = '1') AS active_count, SUM(is_active = '0') AS in_active_count, COUNT(*) AS total_count  FROM loading_point;"
	rl.l.Info("GetEmpolyee: ", countQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var active, inActive, totalCount sql.NullInt64

	errE := row.Scan(&active, &inActive, &totalCount)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	load := dtos.LoadUnloadPoints{
		ActiveCount:   active.Int64,
		InActiveCount: inActive.Int64,
		TotalCount:    totalCount.Int64,
	}

	return &load, nil
}

func (rl *DashBoardObj) GetEmpolyee() (*dtos.Employee, error) {

	countQuery := "SELECT SUM(is_active = '1') AS active_count, SUM(is_active = '0') AS in_active_count, COUNT(*) AS total_count FROM employee;"
	rl.l.Info("GetEmpolyee: ", countQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var active, inActive, totalCount sql.NullInt64

	errE := row.Scan(&active, &inActive, &totalCount)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	emp := dtos.Employee{
		ActiveCount:   active.Int64,
		InActiveCount: inActive.Int64,
		TotalCount:    totalCount.Int64,
	}

	return &emp, nil
}

func (rl *DashBoardObj) GetTripStatusCounts(tripStatus, start, end, column string) (*dtos.TripStatusStats, error) {

	whereQuery := fmt.Sprintf(`SELECT 
  		SUM(CASE WHEN load_status = '%v' AND trip_sheet_type = '%v' THEN 1 ELSE 0 END) AS local_adhoc_trip_count,
  		SUM(CASE WHEN load_status = '%v' AND trip_sheet_type = '%v' THEN 1 ELSE 0 END) AS local_scheduled_trip_count,
  		SUM(CASE WHEN load_status = '%v' AND trip_sheet_type = '%v' THEN 1 ELSE 0 END) AS linehaul_scheduled_trip_count,
  		SUM(CASE WHEN load_status = '%v' AND trip_sheet_type = '%v' THEN 1 ELSE 0 END) AS linehaul_adhoc_trip_count
		FROM trip_sheet_header WHERE %v BETWEEN '%v' AND '%v';`,
		tripStatus,
		constant.LOCAL_ADHOC_TRIP,
		tripStatus,
		constant.LOCAL_SCHEDULED_TRIP,
		tripStatus,
		constant.LINE_HUAL_SCHEDULED_TRIP,
		tripStatus,
		constant.LINE_HUAL_ADHOC_TRIP,
		column,
		start,
		end)

	rl.l.Info("GetTripStatusCounts: ", whereQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(whereQuery)
	var localAdhocTripCount, localScheduledTripCount, linehaulScheduledTripCount,
		linehaulAdhocTripCount sql.NullInt64

	errE := row.Scan(&localAdhocTripCount, &localScheduledTripCount, &linehaulScheduledTripCount, &linehaulAdhocTripCount)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	tripStat := dtos.TripStatusStats{}
	tripStat.LocalAdhocTrip = localAdhocTripCount.Int64
	tripStat.LocalScheduledTrip = localScheduledTripCount.Int64
	tripStat.LineHaulScheduledTrip = linehaulScheduledTripCount.Int64
	tripStat.LineHaulAdhocTrip = linehaulAdhocTripCount.Int64
	tripStat.TotalCount = tripStat.LocalAdhocTrip + tripStat.LocalScheduledTrip + tripStat.LineHaulScheduledTrip + tripStat.LineHaulAdhocTrip
	return &tripStat, nil
}

func (rl *DashBoardObj) GetTripGraph(startDateTime, endDateTime string) (*[]dtos.TripCountWithType, error) {
	list := []dtos.TripCountWithType{}
	whereQuery := fmt.Sprintf(`
            SELECT  DATE(open_trip_date_time) AS date,
                    trip_sheet_type,
                    COUNT(*) AS tripCount
            FROM trip_sheet_header
            WHERE open_trip_date_time BETWEEN '%s' AND '%s'
            GROUP BY DATE(open_trip_date_time), trip_sheet_type
            ORDER BY date;`, startDateTime, endDateTime)

	rl.l.Info("GetTripGraph query: ", whereQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(whereQuery)
	if err != nil {
		rl.l.Error("Error Customers ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dateT, tripSheetType sql.NullString
		var tripCount sql.NullInt64

		tripC := &dtos.TripCountWithType{}
		err := rows.Scan(&dateT, &tripSheetType, &tripCount)
		if err != nil {
			rl.l.Error("Error GetTripGraph scan: ", err)
			return nil, err
		}
		tripC.TripCount = tripCount.Int64
		tripC.Date = dateT.String
		tripC.TripSheetType = tripSheetType.String

		list = append(list, *tripC)
	}

	return &list, nil
}

func (rl *DashBoardObj) GetCountByInvoiceStatus(invoiceStatus string) int64 {
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

func (br *DashBoardObj) GetAllInvoicesWithoutLimit(status string) (*[]dtos.InvoiceBreach, error) {
	list := []dtos.InvoiceBreach{}

	customerInvoiceQuery := fmt.Sprintf(`SELECT ci.id, ci.invoice_number, ci.invoice_status, ci.invoice_date, c.payment_terms from customer_invoice as ci LEFT JOIN customers c ON c.customer_id = ci.customer_id Where ci.invoice_status = '%v';`, status)

	br.l.Info("customerInvoiceQuery:\n ", customerInvoiceQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(customerInvoiceQuery)
	if err != nil {
		br.l.Error("Error LoadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invoiceNumber, invoiceStatus, invoiceDate, paymentTerms sql.NullString
		var id sql.NullInt64

		invoiceObj := &dtos.InvoiceBreach{}
		err := rows.Scan(&id, &invoiceNumber, &invoiceStatus, &invoiceDate, &paymentTerms)
		if err != nil {
			br.l.Error("Error GetAllInvoices scan: ", err)
			return nil, err
		}
		invoiceObj.ID = id.Int64
		invoiceObj.InvoiceNumber = invoiceNumber.String
		invoiceObj.InvoiceStatus = invoiceStatus.String
		invoiceObj.InvoiceDate = invoiceDate.String
		invoiceObj.InvoiceDate = invoiceDate.String
		invoiceObj.PaymentTerms = paymentTerms.String
		list = append(list, *invoiceObj)
	}

	return &list, nil
}
