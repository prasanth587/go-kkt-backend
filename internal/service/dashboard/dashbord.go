package dashboard

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
)

type DashBoardObj struct {
	l            *log.Logger
	dbConnMSSQL  *mssqlcon.DBConn
	dashBoardDao daos.DashBoardDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *DashBoardObj {
	return &DashBoardObj{
		l:            l,
		dbConnMSSQL:  dbConnMSSQL,
		dashBoardDao: daos.NewDashBoardObj(l, dbConnMSSQL),
	}
}

func (dab *DashBoardObj) GetStats(startDate, endDate, employee, vendor, vehicle, customer, loadUnloadPoints,
	tripCreated, tripSubmitted, tripDelivered, tripClosed, tripGraph, tripCompleted, inventory string) (*dtos.DashBoardRes, error) {
	res := dtos.DashBoardRes{}

	if startDate == "" {
		return nil, errors.New("start date should not empty")
	}
	if endDate == "" {
		return nil, errors.New("end date should not empty")
	}

	if employee == "true" {
		employee, errA := dab.dashBoardDao.GetEmpolyee()
		if errA != nil {
			dab.l.Error("ERROR: GetEmpolyee", errA)
			return nil, errA
		}
		res.Employee = employee
	}
	if loadUnloadPoints == "true" {
		loadRes, errA := dab.dashBoardDao.GetLoadingUnLoadingPoints()
		if errA != nil {
			dab.l.Error("ERROR: GetLoadingUnLoadingPoints", errA)
			return nil, errA
		}
		res.LoadUnloadPoints = loadRes
	}

	if customer == "true" {
		loadRes, errA := dab.dashBoardDao.GetCustomer()
		if errA != nil {
			dab.l.Error("ERROR: GetCustomer", errA)
			return nil, errA
		}
		res.Customer = loadRes
	}
	if vendor == "true" {
		loadRes, errA := dab.dashBoardDao.GetVendor()
		if errA != nil {
			dab.l.Error("ERROR: GetVendor", errA)
			return nil, errA
		}
		res.Vendor = loadRes
	}
	if vehicle == "true" {
		loadRes, errA := dab.dashBoardDao.GetVhicle()
		if errA != nil {
			dab.l.Error("ERROR: GetVendor", errA)
			return nil, errA
		}
		res.Vehicle = loadRes
	}
	if tripCreated == "true" {
		column := "open_trip_date_time"
		stats, errA := dab.dashBoardDao.GetTripStatusCounts(constant.STATUS_CREATED, startDate, endDate, column)
		if errA != nil {
			dab.l.Error("ERROR: GetTripStatusCounts STATUS_CREATED", errA)
			return nil, errA
		}
		res.TripCreated = stats
	}
	if tripSubmitted == "true" {
		column := "trip_submitted_date"
		stats, errA := dab.dashBoardDao.GetTripStatusCounts(constant.STATUS_SUBMITTED, startDate, endDate, column)
		if errA != nil {
			dab.l.Error("ERROR: GetTripStatusCounts STATUS_SUBMITTED ", errA)
			return nil, errA
		}
		res.TripSubmitted = stats
	}
	if tripDelivered == "true" {
		column := "trip_delivered_date"
		stats, errA := dab.dashBoardDao.GetTripStatusCounts(constant.STATUS_DELIVERED, startDate, endDate, column)
		if errA != nil {
			dab.l.Error("ERROR: GetTripStatusCounts STATUS_DELIVERED", errA)
			return nil, errA
		}
		res.TripDelivered = stats
	}
	if tripClosed == "true" {
		column := "trip_closed_date"
		stats, errA := dab.dashBoardDao.GetTripStatusCounts(constant.STATUS_CLOSED, startDate, endDate, column)
		if errA != nil {
			dab.l.Error("ERROR: GetTripStatusCounts STATUS_CLOSED", errA)
			return nil, errA
		}
		res.TripClosed = stats
	}
	if tripCompleted == "true" {
		column := "trip_completed_date"
		stats, errA := dab.dashBoardDao.GetTripStatusCounts(constant.STATUS_COMPLETED, startDate, endDate, column)
		if errA != nil {
			dab.l.Error("ERROR: GetTripStatusCounts STATUS_COMPLETED", errA)
			return nil, errA
		}
		res.TripCompleted = stats
	}

	if tripGraph == "true" {
		statsList, errA := dab.dashBoardDao.GetTripGraph(startDate, endDate)
		if errA != nil {
			dab.l.Error("ERROR: GetTripStatusCounts STATUS_CLOSED", errA)
			return nil, errA
		}

		if len(*statsList) > 0 {
			res.TripGraphStats = aggregateTripCounts(statsList)
		}
	}

	if inventory == "true" {
		inv := dtos.InvStatusResponse{}
		draftCount := dab.dashBoardDao.GetCountByInvoiceStatus(constant.STATUS_INVOICE_DRAFT)

		resN, errN := dab.dashBoardDao.GetAllInvoicesWithoutLimit(constant.STATUS_INVOICE_RAISED)
		if errN != nil {
			dab.l.Error("ERROR: GetAllInvoicesWithoutLimit", errN)
			//return nil, errN
		}
		//
		breachCount := 0
		for _, invoiceBreach := range *resN {
			isBreach, errN := dab.IsBreached(invoiceBreach.InvoiceDate, invoiceBreach.PaymentTerms)
			if errN != nil {
				dab.l.Error("ERROR: IsBreached", errN)
				return nil, errN
			}
			if isBreach {
				breachCount++
			}
		}

		inv.DraftCount = draftCount
		inv.InvcoiceRaised = len(*resN)
		inv.InvcoiceOverDue = breachCount
		res.InvoiceStats = inv
	}

	return &res, nil
}

func aggregateTripCounts(input *[]dtos.TripCountWithType) *[]dtos.TripSummary {
	// Map date -> TripSummary
	summaryMap := make(map[string]*dtos.TripSummary)

	for _, record := range *input {
		s, exists := summaryMap[record.Date]
		if !exists {
			s = &dtos.TripSummary{
				Date:   record.Date,
				ByType: make(map[string]int64),
			}
			summaryMap[record.Date] = s
		}
		s.ByType[record.TripSheetType] += record.TripCount
		s.TotalTrips += record.TripCount
	}

	// Convert map to slice
	result := make([]dtos.TripSummary, 0, len(summaryMap))
	for _, v := range summaryMap {
		result = append(result, *v)
	}

	// Optional: sort by date ascending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return &result
}

func (inv *DashBoardObj) IsBreached(invoiceDate string, paymentTerms string) (bool, error) {
	if invoiceDate == "" {
		return false, nil
	}

	paymentTerms = "2"
	if paymentTerms == "" {
		paymentTerms = "30"
	}

	// Parse the invoice date
	dateLayout := "2006-01-02"
	invoiceD, err := time.Parse(dateLayout, invoiceDate)
	if err != nil {
		return false, fmt.Errorf("invalid date: %v", err)
	}

	// Parse the term days
	//terms := strings.Fields(paymentTerms) // ["60", "days"]
	days, err := strconv.Atoi(paymentTerms)
	if err != nil {
		//return false, fmt.Errorf("invalid payment terms: %v", err)
		inv.l.Error("ERROR: invalid payment terms: %v", err)
		days = 30
	}

	// Calculate due date
	dueDate := invoiceD.AddDate(0, 0, days)

	// Compare
	breached := time.Now().After(dueDate)
	return breached, nil
}
