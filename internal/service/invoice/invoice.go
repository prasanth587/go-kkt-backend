package invoice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/prabha303-vi/log-util/log"
	"github.com/xuri/excelize/v2"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	"go-transport-hub/internal/service/commonsvc"
	"go-transport-hub/internal/service/notification"
	"go-transport-hub/internal/service/tripcompletexls"
	"go-transport-hub/utils"
)

type InvoiceObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	invoiceDao  daos.InvoiceDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *InvoiceObj {
	return &InvoiceObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		invoiceDao:  daos.NewInvoiceObj(l, dbConnMSSQL),
	}
}

func (inv *InvoiceObj) GetAllInvoices(limit, offset, searchText, status string, customerId string) (*dtos.InvoiceEntries, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	whereQuery := inv.invoiceDao.BuildWhereQuery(searchText, status, customerId)

	res, errA := inv.invoiceDao.GetAllInvoices(limitI, offsetI, whereQuery)
	if errA != nil {
		inv.l.Error("ERROR: GetLoadingPoints", errA)
		return nil, errA
	}

	for i := range *res {
		inv.ProcessInvoiceDatesAndStatus(&(*res)[i])
	}

	invoiceEntries := dtos.InvoiceEntries{}
	invoiceEntries.InvoiceObj = res
	invoiceEntries.Total = inv.invoiceDao.GetTotalCount(whereQuery)
	invoiceEntries.Limit = limitI
	invoiceEntries.OffSet = offsetI
	return &invoiceEntries, nil
}

func (inv *InvoiceObj) formatDateRange(start, end time.Time) string {
	return fmt.Sprintf("%s %d - %s %d %d",
		start.Format("Jan"), start.Day(),
		end.Format("Jan"), end.Day(),
		end.Year())
}

func (trp *InvoiceObj) CancelDraftInvoice(invoiceId int64) (*dtos.Messge, error) {

	trp.l.Info("CancelDraftInvoice", "InIt")
	errU := trp.invoiceDao.UpdateInvoiceStatus(invoiceId, constant.STATUS_CANCELLED)
	if errU != nil {
		errS := errU.Error()
		trp.l.Error("ERROR: CancelDraftInvoice", invoiceId, errU, errS)
		return nil, errU
	}

	// Get invoice info for notification
	invoiceRes, errU := trp.invoiceDao.GetInvoiceInfo(invoiceId)
	if errU == nil && invoiceRes != nil {
		// Get orgId from customer
		if invoiceRes.CustomerId > 0 {
			customerDao := daos.NewCustomerObj(trp.l, trp.dbConnMSSQL)
			customer, errC := customerDao.GetCustomerV1(invoiceRes.CustomerId)
			if errC == nil && customer != nil {
			// Send notification for invoice cancellation
			// Temporarily disabled
			// notificationSvc := notification.New(trp.l, trp.dbConnMSSQL)
			// invoiceNo := invoiceRes.InvoiceNumber
			// if invoiceNo == "" {
			// 	invoiceNo = fmt.Sprintf("Draft Invoice #%d", invoiceId)
			// }
			// if err := notificationSvc.NotifyInvoiceCancelled(int64(customer.OrgId), invoiceId, invoiceNo); err != nil {
			// 	trp.l.Error("ERROR: Failed to send invoice cancellation notification: ", err)
					// Don't fail the request if notification fails
				}
			}
		}
	}

	trp.l.Info("DraftInvoice cancelled successfully... ", invoiceId)
	roleResponse := dtos.Messge{}
	roleResponse.Message = "Draft Invoice cancelled successfully!"
	return &roleResponse, nil
}

func (trp *InvoiceObj) UpdateInvoiceNumber(invoiceId int64, invoiceNumber string) (*dtos.Messge, error) {

	if invoiceNumber == "" {
		trp.l.Error("Error: invoice number name should not empty")
		return nil, errors.New("invoice number name should not empty")
	}

	invoiceRes, errU := trp.invoiceDao.GetInvoiceInfo(invoiceId)
	if errU != nil {
		errS := errU.Error()
		trp.l.Error("ERROR: UpdateInvoiceNumber", invoiceId, errU, errS)
		return nil, errU
	}

	if invoiceRes != nil {
		invoiceDate := utils.GetCurrentDateStr()
		errU := trp.invoiceDao.UpdateInvoiceNumber(invoiceId, invoiceNumber, constant.STATUS_INVOICE_RAISED, invoiceDate)
		if errU != nil {
			trp.l.Error("ERROR: UpdateInvoiceNumber", invoiceId, invoiceNumber, errU.Error())
			return nil, errU
		}
		tripsFound, errF := trp.invoiceDao.InvoiceNumberIsExistsAtTripSheet(invoiceId)
		if errF != nil {
			trp.l.Error("ERROR: InvoiceNumberIsExistsAtTripSheet", invoiceId, invoiceNumber, errF.Error())
			return nil, errF
		}
		if tripsFound > 0 {
			errT := trp.invoiceDao.UpdateInvoiceNumberAtTripSheet(invoiceId, invoiceNumber, constant.STATUS_INVOICE_RAISED)
			if errT != nil {
				trp.l.Error("ERROR: UpdateInvoiceNumberAtTripSheet", invoiceId, invoiceNumber, errT.Error())
				return nil, errT
			}
		} else {
			trp.l.Error("ERROR: Trips are not found for the invoice", invoiceId, invoiceNumber)
			return nil, errors.New("trips are not found for the invoice")
		}
	}

	// Send notification for invoice number update
	if invoiceRes.CustomerId > 0 {
		customerDao := daos.NewCustomerObj(trp.l, trp.dbConnMSSQL)
		customer, errC := customerDao.GetCustomerV1(invoiceRes.CustomerId)
		if errC == nil && customer != nil {
		// Temporarily disabled
		// notificationSvc := notification.New(trp.l, trp.dbConnMSSQL)
		// if err := notificationSvc.NotifyInvoiceNumberUpdated(int64(customer.OrgId), invoiceId, invoiceNumber); err != nil {
				trp.l.Error("ERROR: Failed to send invoice number update notification: ", err)
				// Don't fail the request if notification fails
			}
		}
	}

	trp.l.Info("Invoice updated successfully... ", invoiceId, invoiceNumber)
	roleResponse := dtos.Messge{}
	roleResponse.Message = "Invoice number updated  successfully! " + invoiceNumber
	return &roleResponse, nil
}

func (trp *InvoiceObj) UpdateInvoicePaid(invoiceId int64, invoicePaidDate, transactionId string) (*dtos.Messge, error) {

	if invoiceId == 0 {
		trp.l.Error("Error: is not a valid invoiceId")
		return nil, errors.New("is not a valid invoiceId")
	}
	if invoicePaidDate == "" {
		trp.l.Error("Error: invoice paid date should not empty")
		return nil, errors.New("invoice paid date should not empty")
	}
	if transactionId == "" {
		trp.l.Error("Error: transaction id should not empty")
		return nil, errors.New("transaction id should not empty")
	}

	invoiceRes, errU := trp.invoiceDao.GetInvoiceInfo(invoiceId)
	if errU != nil {
		errS := errU.Error()
		trp.l.Error("ERROR: GetInvoiceInfo", invoiceId, errU, errS)
		return nil, errU
	}

	if invoiceRes != nil {
		errU := trp.invoiceDao.UpdateInvoicePaid(invoiceId, constant.STATUS_PAID, invoicePaidDate, transactionId)
		if errU != nil {
			trp.l.Error("ERROR: UpdateInvoicePaid", invoiceId, errU.Error())
			return nil, errU
		}
		tripsFound, errF := trp.invoiceDao.InvoiceNumberIsExistsAtTripSheet(invoiceId)
		if errF != nil {
			trp.l.Error("ERROR: InvoiceNumberIsExistsAtTripSheet", invoiceId, errF.Error())
			return nil, errF
		}
		if tripsFound > 0 {
			errT := trp.invoiceDao.UpdateInvoiceNumberAtTripSheetPaid(invoiceId, constant.STATUS_PAID)
			if errT != nil {
				trp.l.Error("ERROR: UpdateInvoiceNumberAtTripSheetPaid", invoiceId, invoiceRes.InvoiceNumber, errT.Error())
				return nil, errT
			}
		} else {
			trp.l.Error("ERROR: Trips are not found for the invoice", invoiceId, invoiceRes.InvoiceNumber)
			return nil, errors.New("trips are not found for the invoice")
		}
	}

	// Send notification for invoice paid
	if invoiceRes.CustomerId > 0 {
		customerDao := daos.NewCustomerObj(trp.l, trp.dbConnMSSQL)
		customer, errC := customerDao.GetCustomerV1(invoiceRes.CustomerId)
		if errC == nil && customer != nil {
		// Temporarily disabled
		// notificationSvc := notification.New(trp.l, trp.dbConnMSSQL)
		// if err := notificationSvc.NotifyInvoicePaid(int64(customer.OrgId), invoiceId, invoiceRes.InvoiceNumber, invoiceRes.InvoiceAmount); err != nil {
				trp.l.Error("ERROR: Failed to send invoice paid notification: ", err)
				// Don't fail the request if notification fails
			}
		}
	}

	trp.l.Info("Invoice Paid updated successfully... ", invoiceId, invoiceRes.InvoiceNumber)
	roleResponse := dtos.Messge{}
	roleResponse.Message = "Invoice number updated  successfully! " + invoiceRes.InvoiceNumber
	return &roleResponse, nil
}

func (inv *InvoiceObj) DownlaodInvoiceXls(invoiceId int64) ([]byte, error) {

	if invoiceId == 0 {
		return nil, errors.New("invoiceId is not empty")
	}

	invoiceRes, tripSheetIdPesent, errA := inv.invoiceDao.GetTripSheetByInvoiceId(invoiceId)
	if errA != nil {
		inv.l.Error("ERROR: GetTripSheetByInvoiceId", errA)
		return nil, errA
	}

	f := excelize.NewFile()
	sheetName := "TripSheets"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}
	f.SetActiveSheet(index)
	styleCell, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			//Color:  "FF0000",
			Family: "Arial",
			Size:   11,
		},

		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#cfe2f3"}, // light blue low background
			Pattern: 1,                   // Solid fill
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
	})

	errSp := f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,    // No columns frozen
		YSplit:      1,    // Freeze rows above row 2 (so row 1 is frozen)
		TopLeftCell: "A2", // The first cell below the frozen pane
		ActivePane:  "bottomLeft",
	})
	if errSp != nil {
		inv.l.Error("ERROR: setting panes", err)
		return nil, errSp
	}

	for colIdx, header := range tripcompletexls.DRAFT_TRIP_SHEET_HEADER {
		// 1. Get the cell for the header (row 1, column colIdx+1)
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)

		// 2. Set the style for the header cell
		err := f.SetCellStyle(sheetName, cell, cell, styleCell)
		if err != nil {
			inv.l.Error("ERROR: setting style", err)
		}

		// 3. Set the header value
		f.SetCellValue(sheetName, cell, header)

		// 4. Set the column width based on header length
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		err = f.SetColWidth(sheetName, colName, colName, float64(len(header)+2))
		if err != nil {
			inv.l.Error("ERROR: setting column width", err)
		}
	}

	//tripSheetIdPesent
	if len(tripSheetIdPesent) != 0 {
		idsCSV := strings.Join(tripSheetIdPesent, ",")
		inv.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := daos.NewTripSheetXls(inv.l, inv.dbConnMSSQL).GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			inv.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		tpmgt := commonsvc.New(inv.l, inv.dbConnMSSQL)
		loadingMap, unloadingMap := tpmgt.BuildCityMaps(loadingPoints)

		leftAlignStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
			},
		})

		for rowIdx, sheet := range *invoiceRes {
			row := rowIdx + 2 // 1-based, and header is row 1
			//podReceived := "No"
			podRequied := "No"
			// if sheet.PodReceived == 1 {
			// 	podReceived = "Yes"
			// }
			if sheet.PodRequired == 1 {
				podRequied = "Yes"
			}
			fromLoc := ""
			toLoc := ""
			if val, ok := loadingMap[sheet.TripSheetID]; ok {
				for _, locations := range val {
					if fromLoc == "" {
						fromLoc = locations.CityCode
					} else {
						fromLoc = fromLoc + " - " + locations.CityCode
					}
				}
			}
			if val, ok := unloadingMap[sheet.TripSheetID]; ok {
				for _, locations := range val {
					if toLoc == "" {
						toLoc = locations.CityCode
					} else {
						toLoc = toLoc + " - " + locations.CityCode
					}
				}
			}

			inv.l.Debug("fromLoc", sheet.TripSheetID, fromLoc, " toLoc: ", toLoc)

			cols := []interface{}{
				sheet.TripSheetNum,
				sheet.TripType,
				sheet.TripSheetType,
				sheet.LoadHoursType,
				sheet.OpenTripDateTime,
				sheet.CustomerName,
				sheet.CustomerCode,
				fromLoc,
				toLoc,
				sheet.VehicleNumber,
				sheet.VehicleSize,
				sheet.MobileNumber,
				sheet.DriverName,
				sheet.LRNumber,
				podRequied,
				sheet.ZonalName,
				sheet.LoadStatus,
				sheet.CustomerPerLoadHire,
				sheet.CustomerRunningKM,
				sheet.CustomerBaseRate,
				sheet.CustomerPerKMPrice,
				sheet.CustomerKMCost,
				sheet.CustomerToll,
				sheet.CustomerExtraHours,
				sheet.CustomerExtraKM,
				sheet.CustomerTotalHire,
				sheet.CustomerCloseTripDateTime,
				sheet.CustomerPaymentReceivedDate,
				sheet.CustomerDebitAmount,
				sheet.CustomerBillingRaisedDate,
				sheet.CustomerLoadCancelled,
				sheet.CustomerReportedDateTimeForHaltingCalc,
				sheet.CustomerRemark,
			}
			for colIdx, val := range cols {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
				f.SetCellValue(sheetName, cell, val)
				// Apply left-align style to the cell
				err := f.SetCellStyle(sheetName, cell, cell, leftAlignStyle)
				if err != nil {
					inv.l.Error("ERROR: setting left-align style", err)
				}
			}
		}
	}

	// Write data rows

	// Delete the default Sheet1
	if err := f.DeleteSheet("Sheet1"); err != nil {
		inv.l.Warn("Could not delete Sheet1:", err)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		inv.l.Error("ERROR: WriteToBuffer", errA)
		return nil, err
	}
	return buf.Bytes(), nil

}

func (inv *InvoiceObj) GetInvoiceInfo(invoiceId int64) (*dtos.InvoiceInfo, error) {
	invoiceInfoResponse := dtos.InvoiceInfo{}
	if invoiceId == 0 {
		return nil, errors.New("invoiceId is not empty")
	}

	invoiceRes, errU := inv.invoiceDao.GetInvoiceInfo(invoiceId)
	if errU != nil {
		errS := errU.Error()
		inv.l.Error("ERROR: GetInvoiceInfo", invoiceId, errU, errS)
		return nil, errU
	}

	inv.ProcessInvoiceDatesAndStatus(invoiceRes)

	copier.Copy(&invoiceInfoResponse, &invoiceRes)

	//TotalTrips
	tripRes, tripSheetIdPesent, errA := inv.invoiceDao.GetTripSheetByInvoiceId(invoiceId)
	if errA != nil {
		inv.l.Error("ERROR: GetTripSheetByInvoiceId", errA)
		return nil, errA
	}
	var tripSheetsForInvoice []dtos.TripSheetsForInvoice

	if len(tripSheetIdPesent) != 0 {
		idsCSV := strings.Join(tripSheetIdPesent, ",")
		inv.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := daos.NewTripSheetXls(inv.l, inv.dbConnMSSQL).GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			inv.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		tpmgt := commonsvc.New(inv.l, inv.dbConnMSSQL)
		loadingMap, unloadingMap := tpmgt.BuildCityMaps(loadingPoints)

		for _, sheet := range *tripRes {

			var tsi dtos.TripSheetsForInvoice
			// Copy matching fields automatically
			copier.Copy(&tsi, &sheet)

			fromLoc := ""
			toLoc := ""
			if val, ok := loadingMap[sheet.TripSheetID]; ok {
				for _, locations := range val {
					if fromLoc == "" {
						fromLoc = locations.CityCode
					} else {
						fromLoc = fromLoc + " - " + locations.CityCode
					}
				}
			}
			if val, ok := unloadingMap[sheet.TripSheetID]; ok {
				for _, locations := range val {
					if toLoc == "" {
						toLoc = locations.CityCode
					} else {
						toLoc = toLoc + " - " + locations.CityCode
					}
				}
			}

			tsi.FromLocaion = fromLoc
			tsi.ToLocation = toLoc
			tripSheetsForInvoice = append(tripSheetsForInvoice, tsi)
		}
	}

	invoiceInfoResponse.TripSheetsForInvoice = tripSheetsForInvoice
	return &invoiceInfoResponse, nil
}

func (inv *InvoiceObj) ProcessInvoiceDatesAndStatus(invoiceRes *dtos.InvoiceObj) {
	switch invoiceRes.InvoiceStatus {
	case constant.STATUS_INVOICE_DRAFT:
		invoiceRes.InvoiceStatusDesc = constant.STATUS_INVOICE_DRAFT_DESC
	case constant.STATUS_INVOICE_RAISED:
		invoiceRes.InvoiceStatusDesc = constant.STATUS_INVOICE_RAISED_DESC
		isBreach, errN := inv.IsBreached(invoiceRes.InvoiceDate, invoiceRes.PaymentTerms)
		if errN != nil {
			inv.l.Error("ERROR: IsBreached", errN)
		}
		invoiceRes.IsOverDue = isBreach

	case constant.STATUS_CANCELLED:
		invoiceRes.InvoiceStatusDesc = constant.STATUS_INVOICE_CANCELLED_DESC
	case constant.STATUS_PAID:
		invoiceRes.InvoiceStatusDesc = constant.STATUS_PAID_DESC
	}

	if invoiceRes.WorkStartDate != "" && invoiceRes.WorkEndDate != "" {
		workStartDate, _ := time.Parse(utils.DFyyyyMMdd, invoiceRes.WorkStartDate)
		workEndDate, _ := time.Parse(utils.DFyyyyMMdd, invoiceRes.WorkEndDate)
		invoiceRes.WorkPeriod = inv.formatDateRange(workStartDate, workEndDate)
	}

	if invoiceRes.DocumentDate != "" {
		documentDate, _ := time.Parse(utils.DFyyyyMMdd, invoiceRes.DocumentDate)
		invoiceRes.DocumentDateStr = documentDate.Format(utils.DFMMMDDyyyy)
	}

}

func (inv *InvoiceObj) IsBreached(invoiceDate string, paymentTerms string) (bool, error) {
	if invoiceDate == "" {
		return false, nil
	}

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

func (inv *InvoiceObj) GetInvoicePDFInfo(invoiceID int64) (*dtos.InvoicePDFInfo, error) {

	res := dtos.InvoicePDFInfo{}

	res.Header = "BILL OF SUPPLY"
	res.BilledFromText = "Billed From : "
	res.BilledFromOfficeName = "KK TRANSPORT"
	res.BilledFromOfficeAddress = "Plot No 6, Shanmuga Nagar, Attanthangal Village, Redhills, Tiruvallur, Chennai- 600052, Tamil Nadu. "
	res.BilledFromOfficeEmail = "kktransportchennai@gmail.com"
	res.BilledFromOfficeGSTIN = "33GJFPS5370M1ZM"
	res.BilledFromOfficePan = "GJFPS5370M"
	res.BilledToText = "Billed To"
	res.BilledToOffice = "MCNEF co.in"
	res.BilledToAddress = "Attanthangal Village XXXXXXXXXXXXXXXX Bangalore"
	res.BilledToState = "Karnataka"
	res.BilledToGSTINUIN = "33AABCF8078M1Z8"

	res.PaymentTermsText = "Payment Terms :"
	res.InFavourOf = "Payment Terms :"

	res.BankBranch = "CANARA Chenani"
	res.AccountNo = "10033AABCF8078M1Z8"
	res.IFSCCode = "DNE233D4"

	invoiceDetails := dtos.InvoiceDetails{}
	invoiceDetails.Amount = "23000"
	invoiceDetails.DescriptionOfServices = "Freight Charges For The Transportation Of Your Products by Road"
	invoiceDetails.HSNSACCode = "996791"
	invoiceDetails.SLNO = "1"
	invoiceDetails.TotalAmount = "23000.00"
	invoiceDetails.AmountInWords = "Will come later, WIP"
	res.InvoiceDetails = invoiceDetails

	res.TermsAndConditions = "Transportation Of Your Products by Road"
	res.ForCompany = "KK TRANSPORT"

	billSupplyDetails := dtos.BillSupplyDetails{}
	billSupplyDetails.BillSupplyDate = "2025-11-19"
	billSupplyDetails.BillSupplyNo = "20251119"
	billSupplyDetails.BillSupplyTypeOfEnterprise = "MSME type will come"

	res.BillSupplyDetails = billSupplyDetails
	return &res, nil
}
