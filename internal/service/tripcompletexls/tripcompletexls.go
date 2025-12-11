package tripcompletexls

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prabha303-vi/log-util/log"
	"github.com/xuri/excelize/v2"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
	"go-transport-hub/internal/service/commonsvc"
	"go-transport-hub/utils"
)

type TripSheetXls struct {
	l               *log.Logger
	dbConnMSSQL     *mssqlcon.DBConn
	tripSheetXlsDao daos.TripSheetXlsDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *TripSheetXls {
	return &TripSheetXls{
		l:               l,
		dbConnMSSQL:     dbConnMSSQL,
		tripSheetXlsDao: daos.NewTripSheetXls(l, dbConnMSSQL),
	}
}

func (ts *TripSheetXls) GetTripSheets(orgId int64, tripSheetId string, startTime string, endTime, tripStatus,
	tripSearchTex, fromDate, toDate, podRequired, podReceived, limit, offset, customers, vendors string) (*dtos.TripSheetXlsRes, error) {

	if startTime == "" {
		return nil, errors.New("start date should not empty")
	}
	if endTime == "" {
		return nil, errors.New("end date should not empty")
	}
	limitI := 100
	offsetI := 0
	if limit != "" {
		limitS, errInt := strconv.ParseInt(limit, 10, 64)
		if errInt != nil {
			return nil, errors.New("invalid limit")
		}
		limitI = int(limitS)
	}
	if offset != "" {
		offsetIS, errInt := strconv.ParseInt(offset, 10, 64)
		if errInt != nil {
			return nil, errors.New("invalid offset")
		}
		offsetI = int(offsetIS)
	}

	whereQuery := ts.tripSheetXlsDao.BuildWhereQuery(orgId, tripSheetId, tripStatus, tripSearchTex, fromDate, toDate, podRequired, podReceived, customers, vendors)
	res, tripSheetIds, errA := ts.tripSheetXlsDao.GetTripSheetsV1Sort(orgId, whereQuery, limitI, offsetI)
	if errA != nil {
		ts.l.Error("ERROR: GetTripSheets", errA)
		return nil, errA
	}

	if len(tripSheetIds) != 0 {
		idsCSV := strings.Join(tripSheetIds, ",")
		ts.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := ts.tripSheetXlsDao.GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			ts.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		//ts.l.Info("loadingPoints: ", utils.MustMarshal(loadingPoints, "loadingPoints"))
		tpmgt := commonsvc.New(ts.l, ts.dbConnMSSQL)
		loadingMap, unloadingMap := tpmgt.BuildCityMaps(loadingPoints)

		//ts.l.Info("loadingMap: ******* ", utils.MustMarshal(loadingMap, "loadingMap"))
		//ts.l.Info("unloadingMap: ******* ", utils.MustMarshal(unloadingMap, "unloadingMap"))

		for i := range *res {
			lrRes := &(*res)[i]
			tripSheetID := lrRes.TripSheetID
			fromLoc := ""
			if val, ok := loadingMap[tripSheetID]; ok {
				for _, locations := range val {
					if fromLoc == "" {
						fromLoc = locations.CityCode
					} else {
						fromLoc = fromLoc + " - " + locations.CityCode
					}
				}
			}
			lrRes.FromLocations = fromLoc

			toLoc := ""
			if val, ok := unloadingMap[tripSheetID]; ok {
				for _, locations := range val {
					if toLoc == "" {
						toLoc = locations.CityCode
					} else {
						toLoc = toLoc + " - " + locations.CityCode
					}
				}
			}
			lrRes.ToLocations = toLoc
		}
	}

	loadingpointEntries := dtos.TripSheetXlsRes{}
	loadingpointEntries.TripSheet = res
	loadingpointEntries.Total = ts.tripSheetXlsDao.GetTotalCount(whereQuery)
	loadingpointEntries.Limit = int64(limitI)
	loadingpointEntries.OffSet = int64(offsetI)
	return &loadingpointEntries, nil
}

func (ts *TripSheetXls) GetTripsByIds(orgId int64, tripSheetIds string) ([]byte, error) {

	if tripSheetIds == "" {
		return nil, errors.New("tripSheetIds is not empty")
	}
	tripSheetIds = strings.ReplaceAll(tripSheetIds, " ", "")
	matched, _ := regexp.MatchString(`^\d+(,\d+)*$`, tripSheetIds)
	ts.l.Info("tripSheetIds is expected format: ", matched)

	if !ts.isNotValideTripSheetIDs(tripSheetIds) {
		ts.l.Error("tripSheetIds is not a expected format. ", tripSheetIds)
		return nil, errors.New("tripSheetIds is not a expected format. ")
	}

	res, tripSheetIdPesent, errA := ts.tripSheetXlsDao.GetTripSheetByIds(orgId, tripSheetIds)
	if errA != nil {
		ts.l.Error("ERROR: GetTripSheets", errA)
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
		ts.l.Error("ERROR: setting panes", err)
		return nil, errSp
	}

	for colIdx, header := range TRIP_SHEET_HEADER {
		// 1. Get the cell for the header (row 1, column colIdx+1)
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)

		// 2. Set the style for the header cell
		err := f.SetCellStyle(sheetName, cell, cell, styleCell)
		if err != nil {
			ts.l.Error("ERROR: setting style", err)
		}

		// 3. Set the header value
		f.SetCellValue(sheetName, cell, header)

		// 4. Set the column width based on header length
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		err = f.SetColWidth(sheetName, colName, colName, float64(len(header)+2))
		if err != nil {
			ts.l.Error("ERROR: setting column width", err)
		}
	}

	//tripSheetIdPesent
	if len(tripSheetIdPesent) != 0 {
		idsCSV := strings.Join(tripSheetIdPesent, ",")
		ts.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := ts.tripSheetXlsDao.GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			ts.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		tpmgt := commonsvc.New(ts.l, ts.dbConnMSSQL)
		loadingMap, unloadingMap := tpmgt.BuildCityMaps(loadingPoints)

		leftAlignStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
			},
		})

		for rowIdx, sheet := range *res {
			row := rowIdx + 2 // 1-based, and header is row 1
			podReceived := "No"
			podRequied := "No"
			if sheet.PodReceived == 1 {
				podReceived = "Yes"
			}
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

			ts.l.Debug("fromLoc", sheet.TripSheetID, fromLoc, " toLoc: ", toLoc)

			cols := []interface{}{
				sheet.TripSheetID,
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
				sheet.VendorBaseRate,
				sheet.VendorKMCost,
				sheet.VendorToll,
				sheet.VendorTotalHire,
				sheet.VendorAdvance,
				sheet.VendorPaidDate,
				sheet.VendorCrossDock,
				sheet.VendorRemark,
				sheet.VendorDebitAmount,
				sheet.VendorBalanceAmount,
				sheet.VendorBreakDown,
				sheet.VendorPaidBy,
				sheet.VendorLoadUnLoadAmount,
				sheet.VendorHaltingDays,
				sheet.VendorHaltingPaid,
				sheet.VendorExtraDelivery,
				podReceived,
				sheet.LoadStatus,
				sheet.ZonalName,
				sheet.CustomerInvoiceNo,
				sheet.CustomerBaseRate,
				sheet.CustomerKMCost,
				sheet.CustomerToll,
				sheet.CustomerExtraHours,
				sheet.CustomerExtraKM,
				sheet.CustomerTotalHire,
				sheet.CustomerCloseTripDateTime,
				sheet.CustomerPaymentReceivedDate,
				sheet.CustomerDebitAmount,
				sheet.CustomerBillingRaisedDate,
				sheet.CustomerPerLoadHire,
				sheet.CustomerRunningKM,
				sheet.CustomerPerKMPrice,
				sheet.CustomerPlacedVehicleSize,
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
					ts.l.Error("ERROR: setting left-align style", err)
				}
			}
		}
	}

	// Write data rows

	// Delete the default Sheet1
	if err := f.DeleteSheet("Sheet1"); err != nil {
		ts.l.Warn("Could not delete Sheet1:", err)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		ts.l.Error("ERROR: WriteToBuffer", errA)
		return nil, err
	}
	return buf.Bytes(), nil

}

func (ts *TripSheetXls) GetTripsByIdsTally(orgId int64, tripSheetIds string) ([]byte, error) {
	if tripSheetIds == "" {
		return nil, errors.New("tripSheetIds should not be empty")
	}
	tripSheetIds = strings.ReplaceAll(tripSheetIds, " ", "")

	if !ts.isNotValideTripSheetIDs(tripSheetIds) {
		ts.l.Error("tripSheetIds is not a expected format. ", tripSheetIds)
		return nil, errors.New("tripSheetIds is not a expected format")
	}

	res, tripSheetIdPesent, errA := ts.tripSheetXlsDao.GetTripSheetByIds(orgId, tripSheetIds)
	if errA != nil {
		ts.l.Error("ERROR: GetTripSheetByIds", errA)
		return nil, errA
	}

	if len(*res) == 0 {
		return nil, errors.New("no trip sheets found")
	}

	// Generate Tally XML
	xml := ts.generateTallyXML(res, tripSheetIdPesent)
	return []byte(xml), nil
}

func (ts *TripSheetXls) generateTallyXML(trips *[]dtos.DownloadTripSheetXls, tripSheetIds []string) string {
	var xml strings.Builder
	xml.WriteString(`<?xml version="1.0"?>`)
	xml.WriteString("\n<ENVELOPE>")
	xml.WriteString("\n  <HEADER>")
	xml.WriteString("\n    <VERSION>1</VERSION>")
	xml.WriteString("\n    <TALLYREQUEST>Import</TALLYREQUEST>")
	xml.WriteString("\n    <TYPE>Data</TYPE>")
	xml.WriteString("\n    <ID>Vouchers</ID>")
	xml.WriteString("\n  </HEADER>")
	xml.WriteString("\n  <BODY>")
	xml.WriteString("\n    <DESC>")
	xml.WriteString("\n      <STATICVARIABLES>")
	xml.WriteString("\n        <IMPORTDUPS>@@DUPCOMBINE</IMPORTDUPS>")
	xml.WriteString("\n      </STATICVARIABLES>")
	xml.WriteString("\n    </DESC>")
	xml.WriteString("\n    <DATA>")
	xml.WriteString("\n      <TALLYMESSAGE>")

	// Collect unique customer names and vendor names
	customerMap := make(map[string]bool)
	vendorMap := make(map[string]bool)

	for _, trip := range *trips {
		// Collect customers
		if trip.CustomerName != "" && trip.CustomerTotalHire > 0 && trip.CustomerInvoiceNo != "" {
			customerMap[trip.CustomerName] = true
		}

		// Collect vendors
		vendorTotal := trip.VendorTotalHire
		if trip.VendorLoadUnLoadAmount > 0 {
			vendorTotal += trip.VendorLoadUnLoadAmount
		}
		if trip.VendorHaltingPaid > 0 {
			vendorTotal += trip.VendorHaltingPaid
		}
		if trip.VendorExtraDelivery > 0 {
			vendorTotal += trip.VendorExtraDelivery
		}
		if trip.VendorAdvance > 0 {
			vendorTotal -= trip.VendorAdvance
		}

		if vendorTotal > 0 && trip.VendorPaidDate != "" {
			vendorName := ts.formatVendorName(trip.VendorName, trip.VendorCode, trip.VendorID)
			if vendorName != "" {
				vendorMap[vendorName] = true
				ts.l.Info("Collecting vendor for ledger creation: ", vendorName)
			}
		}
	}

	// Generate "Sales - Transport" ledger (required for sales vouchers)
	xml.WriteString("\n        <LEDGER NAME=\"Sales - Transport\" ACTION=\"Create\">")
	xml.WriteString("\n          <NAME>Sales - Transport</NAME>")
	xml.WriteString("\n          <PARENT>Sales Accounts</PARENT>")
	xml.WriteString("\n        </LEDGER>")

	// Ensure Cash ledger exists (required for payment vouchers)
	xml.WriteString("\n        <LEDGER NAME=\"Cash\" ACTION=\"Create\">")
	xml.WriteString("\n          <NAME>Cash</NAME>")
	xml.WriteString("\n          <PARENT>Cash-in-Hand</PARENT>")
	xml.WriteString("\n        </LEDGER>")

	// Generate Customer Ledgers (Sundry Debtors)
	for customerName := range customerMap {
		if customerName != "" {
			xml.WriteString("\n        <LEDGER NAME=\"" + ts.escapeXML(customerName) + "\" ACTION=\"Create\">")
			xml.WriteString("\n          <NAME>" + ts.escapeXML(customerName) + "</NAME>")
			xml.WriteString("\n          <PARENT>Sundry Debtors</PARENT>")
			xml.WriteString("\n        </LEDGER>")
		}
	}

	// Generate Vendor Ledgers (Sundry Creditors)
	for vendorName := range vendorMap {
		if vendorName != "" {
			escapedName := ts.escapeXML(vendorName)
			ts.l.Info("Creating vendor ledger: ", escapedName)
			xml.WriteString("\n        <LEDGER NAME=\"" + escapedName + "\" ACTION=\"Create\">")
			xml.WriteString("\n          <NAME>" + escapedName + "</NAME>")
			xml.WriteString("\n          <PARENT>Sundry Creditors</PARENT>")
			xml.WriteString("\n        </LEDGER>")
		}
	}

	// Close the first TALLYMESSAGE for ledgers and start vouchers in same envelope
	xml.WriteString("\n      </TALLYMESSAGE>")
	xml.WriteString("\n      <TALLYMESSAGE>")

	voucherCount := 0
	for _, trip := range *trips {
		// Generate Sales Voucher for Customer (if customer has billing)
		if trip.CustomerTotalHire > 0 && trip.CustomerInvoiceNo != "" {
			voucherCount++
			xml.WriteString("\n        <VOUCHER REMOTEID=\"\" VCHKEY=\"\" VCHTYPE=\"Sales\" ACTION=\"Create\">")
			xml.WriteString("\n          <DATE>" + ts.formatDateForTally(trip.CustomerBillingRaisedDate, trip.OpenTripDateTime) + "</DATE>")
			xml.WriteString("\n          <NARRATION>Transport Service - " + trip.TripSheetNum)
			if trip.LRNumber != "" {
				xml.WriteString(" - LR: " + trip.LRNumber)
			}
			xml.WriteString("</NARRATION>")
			xml.WriteString("\n          <VOUCHERTYPE>Sales</VOUCHERTYPE>")
			xml.WriteString("\n          <VOUCHERNUMBER>" + trip.CustomerInvoiceNo + "</VOUCHERNUMBER>")
			xml.WriteString("\n          <PARTYNAME>" + ts.escapeXML(trip.CustomerName) + "</PARTYNAME>")

			xml.WriteString("\n          <ALLLEDGERENTRIES.LIST>")
			xml.WriteString("\n            <LEDGERNAME>" + ts.escapeXML(trip.CustomerName) + "</LEDGERNAME>")
			xml.WriteString("\n            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>")
			xml.WriteString("\n            <AMOUNT>" + ts.formatAmount(trip.CustomerTotalHire) + "</AMOUNT>")
			xml.WriteString("\n          </ALLLEDGERENTRIES.LIST>")

			xml.WriteString("\n          <ALLLEDGERENTRIES.LIST>")
			xml.WriteString("\n            <LEDGERNAME>Sales - Transport</LEDGERNAME>")
			xml.WriteString("\n            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>")
			xml.WriteString("\n            <AMOUNT>" + ts.formatAmount(trip.CustomerTotalHire) + "</AMOUNT>")
			xml.WriteString("\n          </ALLLEDGERENTRIES.LIST>")

			xml.WriteString("\n        </VOUCHER>")
		}

		// Generate Payment Voucher for Vendor (if vendor payment exists)
		// Calculate vendor total: VendorTotalHire + VendorLoadUnLoadAmount + VendorHaltingPaid + VendorExtraDelivery
		vendorTotal := trip.VendorTotalHire
		if trip.VendorLoadUnLoadAmount > 0 {
			vendorTotal += trip.VendorLoadUnLoadAmount
		}
		if trip.VendorHaltingPaid > 0 {
			vendorTotal += trip.VendorHaltingPaid
		}
		if trip.VendorExtraDelivery > 0 {
			vendorTotal += trip.VendorExtraDelivery
		}
		if trip.VendorAdvance > 0 {
			vendorTotal -= trip.VendorAdvance // Advance is deducted
		}

		if vendorTotal > 0 && trip.VendorPaidDate != "" {
			voucherCount++
			xml.WriteString("\n        <VOUCHER REMOTEID=\"\" VCHKEY=\"\" VCHTYPE=\"Payment\" ACTION=\"Create\">")
			xml.WriteString("\n          <DATE>" + ts.formatDateForTally(trip.VendorPaidDate, trip.OpenTripDateTime) + "</DATE>")
			xml.WriteString("\n          <NARRATION>Vendor Payment - " + trip.TripSheetNum)
			if trip.LRNumber != "" {
				xml.WriteString(" - LR: " + trip.LRNumber)
			}
			xml.WriteString("</NARRATION>")
			xml.WriteString("\n          <VOUCHERTYPE>Payment</VOUCHERTYPE>")

			// Use consistent vendor name formatting (same as ledger creation)
			vendorName := ts.formatVendorName(trip.VendorName, trip.VendorCode, trip.VendorID)
			escapedVendorName := ts.escapeXML(vendorName)
			ts.l.Info("Using vendor name in voucher: ", escapedVendorName, " (original: ", vendorName, ")")

			xml.WriteString("\n          <ALLLEDGERENTRIES.LIST>")
			xml.WriteString("\n            <LEDGERNAME>Cash</LEDGERNAME>")
			xml.WriteString("\n            <ISDEEMEDPOSITIVE>Yes</ISDEEMEDPOSITIVE>")
			xml.WriteString("\n            <AMOUNT>" + ts.formatAmount(vendorTotal) + "</AMOUNT>")
			xml.WriteString("\n          </ALLLEDGERENTRIES.LIST>")

			xml.WriteString("\n          <ALLLEDGERENTRIES.LIST>")
			xml.WriteString("\n            <LEDGERNAME>" + escapedVendorName + "</LEDGERNAME>")
			xml.WriteString("\n            <ISDEEMEDPOSITIVE>No</ISDEEMEDPOSITIVE>")
			xml.WriteString("\n            <AMOUNT>" + ts.formatAmount(vendorTotal) + "</AMOUNT>")
			xml.WriteString("\n          </ALLLEDGERENTRIES.LIST>")

			xml.WriteString("\n        </VOUCHER>")
		}
	}

	xml.WriteString("\n      </TALLYMESSAGE>")
	xml.WriteString("\n    </DATA>")
	xml.WriteString("\n  </BODY>")
	xml.WriteString("\n</ENVELOPE>")

	ledgerCount := len(customerMap) + len(vendorMap) + 1 // +1 for Sales - Transport
	ts.l.Info("Generated Tally XML with ", ledgerCount, " ledgers and ", voucherCount, " vouchers")
	return xml.String()
}

func (ts *TripSheetXls) formatDateForTally(dateStr, fallbackDate string) string {
	if dateStr == "" {
		dateStr = fallbackDate
	}
	// Tally expects YYYYMMDD format
	// Try parsing common formats
	layouts := []string{"2006-01-02", "2006-01-02 15:04:05", "02-01-2006", "02/01/2006", "2006-01-02T15:04:05"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("20060102")
		}
	}
	// If parsing fails, use current date
	return time.Now().Format("20060102")
}

func (ts *TripSheetXls) formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func (ts *TripSheetXls) escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	// Remove invalid XML control characters (0x00-0x1F except tab, newline, carriage return)
	var result strings.Builder
	for _, r := range s {
		if r >= 0x20 || r == 0x09 || r == 0x0A || r == 0x0D {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// formatVendorName formats vendor name consistently for both ledger creation and vouchers
func (ts *TripSheetXls) formatVendorName(vendorName, vendorCode string, vendorID int64) string {
	if vendorName == "" {
		if vendorCode != "" {
			return fmt.Sprintf("%s - Vendor_%d", vendorCode, vendorID)
		}
		return fmt.Sprintf("Vendor_%d", vendorID)
	}
	if vendorCode != "" {
		return fmt.Sprintf("%s - %s", vendorCode, vendorName)
	}
	return vendorName
}

func (ts *TripSheetXls) isNotValideTripSheetIDs(tripSheetIDs string) bool {
	pattern := `^\d+(,\d+)*$`
	matched, errA := regexp.MatchString(pattern, tripSheetIDs)
	if errA != nil {
		ts.l.Error("ERROR: notNotValideTripSheetIDs", errA)
		return false
	}
	return matched
}

func (ts *TripSheetXls) UpdateTripSheetXls(xlFile *excelize.File) (*dtos.XlsUpdateMessge, error) {

	sheetName := xlFile.GetSheetName(0)
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		ts.l.Error("failed to get rows: %w", err)
		return nil, err
	}
	//
	headerMap := make(map[string]int)
	for i, cell := range rows[0] {
		headerMap[cell] = i
	}

	maxCount := 0
	errorRows := []dtos.Messge{}
	successRows := []dtos.Messge{}

	// Process each data row (skip header row)
	for rowIdx, row := range rows[1:] {
		if len(row) < 1 {
			continue // Skip empty rows
		}
		errRow := dtos.Messge{}
		if maxCount > constant.MAX_ROW_UPDATE_TRIPSHEET {
			ts.l.Error("Max row reached: ", row[0], row[1], rowIdx)
			errRow.Message = fmt.Sprintf("Maximum row reached. found more than %v continuing", rowIdx)
			errorRows = append(errorRows, errRow)
			break
		}

		maxCount++

		// Get trip_sheet_id (first column)
		tripSheetID, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			ts.l.Error("failed to get rows: ", row[0], row[1], rowIdx, err)
			errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
			errorRows = append(errorRows, errRow)
			continue
		}

		isExists := ts.tripSheetXlsDao.IsTripExists(tripSheetID)
		if !isExists {
			ts.l.Error("row not found in the table: ", row[0], row[1], tripSheetID)
			errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
			errorRows = append(errorRows, errRow)
			continue
		}

		// Prepare update data
		updateData := dtos.TripSheetUpdateData{
			CustomerInvoiceNo:                      getStringValue(row, headerMap, "Customer Invoice No"),
			CustomerCloseTripDateTime:              getStringValue(row, headerMap, "Customer Close Trip Date Time"),
			CustomerPaymentReceivedDate:            getStringValue(row, headerMap, "Customer Payment Received Date"),
			CustomerRemark:                         getStringValue(row, headerMap, "Customer Remark"),
			CustomerBillingRaisedDate:              getStringValue(row, headerMap, "Customer Billing Raised Date"),
			CustomerBaseRate:                       getFloatValue(row, headerMap, "Customer Base Rate"),
			CustomerKMCost:                         getFloatValue(row, headerMap, "Customer KM Cost"),
			CustomerToll:                           getFloatValue(row, headerMap, "Customer Toll"),
			CustomerExtraHours:                     getFloatValue(row, headerMap, "Customer Extra Hours"),
			CustomerExtraKM:                        getFloatValue(row, headerMap, "Customer Extra KM"),
			CustomerTotalHire:                      getFloatValue(row, headerMap, "Customer Total Hire"),
			CustomerDebitAmount:                    getFloatValue(row, headerMap, "Customer Debit Amount"),
			CustomerPerLoadHire:                    getFloatValue(row, headerMap, "Customer Per Load Hire"),
			CustomerRunningKM:                      getFloatValue(row, headerMap, "Customer Running KM"),
			CustomerPerKMPrice:                     getFloatValue(row, headerMap, "Customer Per KM Price"),
			CustomerPlacedVehicleSize:              getStringValue(row, headerMap, "Customer Placed Vehicle Size"),
			CustomerLoadCancelled:                  getStringValue(row, headerMap, "Customer Load Cancelled"),
			CustomerReportedDateTimeForHaltingCalc: getStringValue(row, headerMap, "Customer Reported For Halting"),
		}
		updateData.LoadStatus = constant.STATUS_COMPLETED
		ts.l.Info("updateData.CustomerInvoiceNo", updateData.CustomerInvoiceNo)
		if updateData.CustomerInvoiceNo == "" {
			ts.l.Error("ERROR: Invoice number should not empty")
			errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], "Customer Invoice number empty")
			errorRows = append(errorRows, errRow)
			continue
		}

		if updateData.CustomerCloseTripDateTime != "" {
			if strings.Contains(updateData.CustomerCloseTripDateTime, "/") {
				date, err := utils.HandleZeroPaddingWithTime(strings.ReplaceAll(updateData.CustomerCloseTripDateTime, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerCloseTripDateTime): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerCloseTripDateTime = date
				ts.l.Info("Modified CustomerCloseTripDateTime string: ", row[0], row[1], updateData.CustomerCloseTripDateTime)
			}
		}
		if updateData.CustomerPaymentReceivedDate != "" {
			if strings.Contains(updateData.CustomerPaymentReceivedDate, "/") {
				date, err := utils.HandleZeroPadding(strings.ReplaceAll(updateData.CustomerPaymentReceivedDate, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerPaymentReceivedDate): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerPaymentReceivedDate = date
				ts.l.Info("Modified CustomerPaymentReceivedDate string: ", row[0], row[1], updateData.CustomerPaymentReceivedDate)
			}
		}

		if updateData.CustomerBillingRaisedDate != "" {
			if strings.Contains(updateData.CustomerBillingRaisedDate, "/") {
				date, err := utils.HandleZeroPadding(strings.ReplaceAll(updateData.CustomerBillingRaisedDate, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerBillingRaisedDate): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerBillingRaisedDate = date
				ts.l.Info("Modified CustomerBillingRaisedDate string: ", row[0], row[1], updateData.CustomerBillingRaisedDate)
			}
		}
		if updateData.CustomerReportedDateTimeForHaltingCalc != "" {
			if strings.Contains(updateData.CustomerReportedDateTimeForHaltingCalc, "/") {
				date, err := utils.HandleZeroPadding(strings.ReplaceAll(updateData.CustomerReportedDateTimeForHaltingCalc, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerReportedDateTimeForHaltingCalc): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerReportedDateTimeForHaltingCalc = date
				ts.l.Info("Modified CustomerReportedDateTimeForHaltingCalc string: ", row[0], row[1], updateData.CustomerReportedDateTimeForHaltingCalc)
			}
		}

		errXls := ts.tripSheetXlsDao.XLSUpdateTripSheetHeader(tripSheetID, updateData)
		if errXls != nil {
			ts.l.Error("ERROR: XLSUpdateTripSheetHeader", errXls)
			errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], errXls)
			errorRows = append(errorRows, errRow)
			continue
		}

		successRow := dtos.Messge{}
		successRow.Message = fmt.Sprintf("Success: %v", row[1])
		successRows = append(successRows, successRow)
	}

	roleResponse := dtos.XlsUpdateMessge{}
	roleResponse.ErrorRows = &errorRows
	roleResponse.Success = &successRows
	roleResponse.ErrorTripsCount = len(errorRows)
	roleResponse.SuccessTripsCount = len(successRows)
	return &roleResponse, nil
}

func getStringValue(row []string, headerMap map[string]int, header string) string {
	if idx, ok := headerMap[header]; ok && idx < len(row) {
		return row[idx]
	}
	return ""
}

func getFloatValue(row []string, headerMap map[string]int, header string) float64 {
	if idx, ok := headerMap[header]; ok && idx < len(row) && row[idx] != "" {
		val, err := strconv.ParseFloat(row[idx], 64)
		if err == nil {
			return val
		}
	}
	return 0
}

func (ts *TripSheetXls) GetTripSheetsInvoice(orgId int64, tripSheetId string, startTime string, endTime, tripStatus,
	tripSearchTex, fromDate, toDate, podRequired, podReceived, limit, offset, customers, vendors string) (*dtos.TripSheetXlsRes, error) {

	if startTime == "" && endTime == "" {
		// Both empty is allowed; do nothing
	} else if startTime == "" {
		return nil, errors.New("start/end date should not be empty")
	} else if endTime == "" {
		return nil, errors.New("start/end date should not be empty")
	}

	limitI := 50
	offsetI := 0
	if limit != "" {
		limitS, errInt := strconv.ParseInt(limit, 10, 64)
		if errInt != nil {
			return nil, errors.New("invalid limit")
		}
		limitI = int(limitS)
	}
	if offset != "" {
		offsetIS, errInt := strconv.ParseInt(offset, 10, 64)
		if errInt != nil {
			return nil, errors.New("invalid offset")
		}
		offsetI = int(offsetIS)
	}

	whereQuery := ts.tripSheetXlsDao.BuildWhereQuery(orgId, tripSheetId, tripStatus, tripSearchTex, fromDate, toDate, podRequired, podReceived, customers, vendors)
	res, tripSheetIds, errA := ts.tripSheetXlsDao.GetTripSheetsV1Sort(orgId, whereQuery, limitI, offsetI)
	if errA != nil {
		ts.l.Error("ERROR: GetTripSheets", errA)
		return nil, errA
	}

	if len(tripSheetIds) != 0 {
		idsCSV := strings.Join(tripSheetIds, ",")
		ts.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := ts.tripSheetXlsDao.GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			ts.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		//ts.l.Info("loadingPoints: ", utils.MustMarshal(loadingPoints, "loadingPoints"))
		tpmgt := commonsvc.New(ts.l, ts.dbConnMSSQL)
		loadingMap, unloadingMap := tpmgt.BuildCityMaps(loadingPoints)

		//ts.l.Info("loadingMap: ******* ", utils.MustMarshal(loadingMap, "loadingMap"))
		//ts.l.Info("unloadingMap: ******* ", utils.MustMarshal(unloadingMap, "unloadingMap"))

		for i := range *res {
			lrRes := &(*res)[i]
			tripSheetID := lrRes.TripSheetID
			fromLoc := ""
			if val, ok := loadingMap[tripSheetID]; ok {
				for _, locations := range val {
					if fromLoc == "" {
						fromLoc = locations.CityCode
					} else {
						fromLoc = fromLoc + " - " + locations.CityCode
					}
				}
			}
			lrRes.FromLocations = fromLoc

			toLoc := ""
			if val, ok := unloadingMap[tripSheetID]; ok {
				for _, locations := range val {
					if toLoc == "" {
						toLoc = locations.CityCode
					} else {
						toLoc = toLoc + " - " + locations.CityCode
					}
				}
			}
			lrRes.ToLocations = toLoc
		}
	}

	loadingpointEntries := dtos.TripSheetXlsRes{}
	loadingpointEntries.TripSheet = res
	loadingpointEntries.Total = ts.tripSheetXlsDao.GetTotalCount(whereQuery)
	loadingpointEntries.Limit = int64(limitI)
	loadingpointEntries.OffSet = int64(offsetI)
	return &loadingpointEntries, nil
}

func (ts *TripSheetXls) GetTripsForDraftInvoice(orgId int64, tripSheetIds string) ([]byte, error) {

	if tripSheetIds == "" {
		return nil, errors.New("tripSheetIds is not empty")
	}

	tripSheetIds = strings.ReplaceAll(tripSheetIds, " ", "")
	matched, _ := regexp.MatchString(`^\d+(,\d+)*$`, tripSheetIds)
	ts.l.Info("tripSheetIds is expected format: ", matched)

	if !ts.isNotValideTripSheetIDs(tripSheetIds) {
		ts.l.Error("tripSheetIds is not a expected format. ", tripSheetIds)
		return nil, errors.New("tripSheetIds is not a expected format. ")
	}

	customerIs, errC := ts.tripSheetXlsDao.CheckSelectedTripsAreSingleCustomers(tripSheetIds)
	if errC != nil {
		ts.l.Error("ERROR: CheckSelectedTripsAreSingleCustomers", errC)
		return nil, errC
	}

	isSameCustomer := ts.CheckAllSame(customerIs)
	if !isSameCustomer {
		return nil, errors.New("selected trips are found different customers, expecting single customers")
	}

	res, tripSheetIdPesent, errA := ts.tripSheetXlsDao.GetTripSheetByIds(orgId, tripSheetIds)
	if errA != nil {
		ts.l.Error("ERROR: GetTripSheets", errA)
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

	styleCellAfter16, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			//Color:  "FF0000",
			Family: "Arial",
			Size:   11,
		},

		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#c5f3f4"}, // light blue low background
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
		ts.l.Error("ERROR: setting panes", err)
		return nil, errSp
	}

	for colIdx, header := range DRAFT_TRIP_SHEET_HEADER {
		// 1. Get the cell for the header (row 1, column colIdx+1)
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)

		var styleCellV int
		if colIdx > 16 {
			//ts.l.Info("styleCellV", colIdx, header)
			styleCellV = styleCellAfter16
		} else {
			styleCellV = styleCell
		}
		// 2. Set the style for the header cell
		err := f.SetCellStyle(sheetName, cell, cell, styleCellV)
		if err != nil {
			ts.l.Error("ERROR: setting style", err)
		}

		// 3. Set the header value
		f.SetCellValue(sheetName, cell, header)

		// 4. Set the column width based on header length
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		err = f.SetColWidth(sheetName, colName, colName, float64(len(header)+2))
		if err != nil {
			ts.l.Error("ERROR: setting column width", err)
		}
	}

	//tripSheetIdPesent
	if len(tripSheetIds) != 0 {
		idsCSV := strings.Join(tripSheetIdPesent, ",")
		ts.l.Info("idsCSV commas: ", idsCSV)
		loadingPoints, err := ts.tripSheetXlsDao.GetLoadingUnloadingPoints(idsCSV)
		if err != nil {
			ts.l.Error("ERROR: loadingPoints", err)
			return nil, err
		}
		tpmgt := commonsvc.New(ts.l, ts.dbConnMSSQL)
		loadingMap, unloadingMap := tpmgt.BuildCityMaps(loadingPoints)

		leftAlignStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
			},
		})

		for rowIdx, sheet := range *res {
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

			ts.l.Debug("fromLoc", sheet.TripSheetID, fromLoc, " toLoc: ", toLoc)

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
					ts.l.Error("ERROR: setting left-align style", err)
				}
			}
		}
	}

	// Write data rows

	// Delete the default Sheet1
	if err := f.DeleteSheet("Sheet1"); err != nil {
		ts.l.Warn("Could not delete Sheet1:", err)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		ts.l.Error("ERROR: WriteToBuffer", errA)
		return nil, err
	}
	return buf.Bytes(), nil

}

func (ts *TripSheetXls) CheckAllSame(numbers []int64) bool {
	allSame := true

	if len(numbers) == 0 {
		fmt.Println("Array is empty")
		return allSame
	}

	first := numbers[0]
	for i, num := range numbers {
		if num != first {
			allSame = false
			ts.l.Errorf("Difference found at index %d: value %d (expected %d)\n", i, num, first)
		}
	}
	return allSame
}

func (ts *TripSheetXls) GenerateDraftInvoiceByTripxls(xlFile *excelize.File) (*dtos.XlsDraftInvoiceMessge, error) {

	sheetName := xlFile.GetSheetName(0)
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		ts.l.Error("failed to get rows: %w", err)
		return nil, err
	}
	//
	headerMap := make(map[string]int)
	for i, cell := range rows[0] {
		headerMap[cell] = i
	}

	maxCount := 0
	errorRows := []dtos.Messge{}
	successRows := []dtos.Messge{}
	tripSheetNumbers := []string{}
	tripSheetNumbersArray := []string{}
	tripOpenTripDateArray := []string{}
	tripSheetIds := []int64{}
	invoiceRef := ""

	// Process each data row (skip header row)
	for rowIdx, row := range rows[1:] {
		if len(row) < 1 {
			continue // Skip empty rows
		}
		errRow := dtos.Messge{}
		if maxCount > constant.MAX_ROW_UPDATE_TRIPSHEET {
			ts.l.Error("Max row reached: ", row[0], row[1], rowIdx)
			errRow.Message = fmt.Sprintf("Maximum row reached. found more than %v continuing", rowIdx)
			errorRows = append(errorRows, errRow)
			break
		}

		maxCount++

		// Get trip_sheet_id (first column)
		tripSheetNumber := row[0]
		if tripSheetNumber == "" {
			ts.l.Error("failed to get rows: ", row[0], row[1], rowIdx)
			errRow.Message = fmt.Sprintf("%v, Error No tripsheet number: %v", rowIdx, row[0])
			errorRows = append(errorRows, errRow)
			continue
		}

		isExists := ts.tripSheetXlsDao.IsTripSheetNumberExists(tripSheetNumber)
		if !isExists {
			ts.l.Error("row not found in the table: ", row[0], row[1], tripSheetNumber)
			errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
			errorRows = append(errorRows, errRow)
			continue
		}

		openTripDateTime := row[4]
		if openTripDateTime != "" {
			tripOpenTripDateArray = append(tripOpenTripDateArray, openTripDateTime)
		}

		// Prepare update data
		updateData := dtos.UpdateDraftInvoiceData{
			CustomerPerLoadHire:                    getFloatValue(row, headerMap, CustomerPerLoadHire),
			CustomerRunningKM:                      getFloatValue(row, headerMap, CustomerRunningKM),
			CustomerBaseRate:                       getFloatValue(row, headerMap, CustomerBaseRate),
			CustomerPerKMPrice:                     getFloatValue(row, headerMap, CustomerPerKMPrice),
			CustomerKMCost:                         getFloatValue(row, headerMap, CustomerKMCost),
			CustomerToll:                           getFloatValue(row, headerMap, CustomerToll),
			CustomerExtraHours:                     getFloatValue(row, headerMap, CustomerExtraHours),
			CustomerExtraKM:                        getFloatValue(row, headerMap, CustomerExtraKM),
			CustomerTotalHire:                      getFloatValue(row, headerMap, CustomerTotalHire),
			CustomerCloseTripDateTime:              getStringValue(row, headerMap, CustomerCloseTripDateTime),
			CustomerPaymentReceivedDate:            getStringValue(row, headerMap, CustomerPaymentReceivedDate),
			CustomerDebitAmount:                    getFloatValue(row, headerMap, CustomerDebitAmount),
			CustomerBillingRaisedDate:              getStringValue(row, headerMap, CustomerBillingRaisedDate),
			CustomerLoadCancelled:                  getStringValue(row, headerMap, CustomerLoadCancelled),
			CustomerReportedDateTimeForHaltingCalc: getStringValue(row, headerMap, CustomerReportedForHalting),
			CustomerRemark:                         getStringValue(row, headerMap, CustomerRemark),
		}

		//updateData.LoadStatus = constant.STATUS_INVOICE_RAISED
		ts.l.Info("updateData.tripSheetNumber", tripSheetNumber)
		// if updateData.CustomerInvoiceNo == "" {
		// 	ts.l.Error("ERROR: Invoice number should not empty")
		// 	errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], "Customer Invoice number empty")
		// 	errorRows = append(errorRows, errRow)
		// 	continue
		// }

		if updateData.CustomerCloseTripDateTime != "" {
			if strings.Contains(updateData.CustomerCloseTripDateTime, "/") {
				date, err := utils.HandleZeroPaddingWithTime(strings.ReplaceAll(updateData.CustomerCloseTripDateTime, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerCloseTripDateTime): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerCloseTripDateTime = date
				ts.l.Info("Modified CustomerCloseTripDateTime string: ", row[0], row[1], updateData.CustomerCloseTripDateTime)
			}
		}
		if updateData.CustomerPaymentReceivedDate != "" {
			if strings.Contains(updateData.CustomerPaymentReceivedDate, "/") {
				date, err := utils.HandleZeroPadding(strings.ReplaceAll(updateData.CustomerPaymentReceivedDate, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerPaymentReceivedDate): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerPaymentReceivedDate = date
				ts.l.Info("Modified CustomerPaymentReceivedDate string: ", row[0], row[1], updateData.CustomerPaymentReceivedDate)
			}
		}

		if updateData.CustomerBillingRaisedDate != "" {
			if strings.Contains(updateData.CustomerBillingRaisedDate, "/") {
				date, err := utils.HandleZeroPadding(strings.ReplaceAll(updateData.CustomerBillingRaisedDate, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerBillingRaisedDate): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerBillingRaisedDate = date
				ts.l.Info("Modified CustomerBillingRaisedDate string: ", row[0], row[1], updateData.CustomerBillingRaisedDate)
			}
		}
		if updateData.CustomerReportedDateTimeForHaltingCalc != "" {
			if strings.Contains(updateData.CustomerReportedDateTimeForHaltingCalc, "/") {
				date, err := utils.HandleZeroPadding(strings.ReplaceAll(updateData.CustomerReportedDateTimeForHaltingCalc, "/", "-"))
				if err != nil {
					ts.l.Error("date format issue (CustomerReportedDateTimeForHaltingCalc): ", row[0], row[1], rowIdx, err)
					errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], err)
					errorRows = append(errorRows, errRow)
					continue
				}
				updateData.CustomerReportedDateTimeForHaltingCalc = date
				ts.l.Info("Modified CustomerReportedDateTimeForHaltingCalc string: ", row[0], row[1], updateData.CustomerReportedDateTimeForHaltingCalc)
			}
		}

		errXls := ts.tripSheetXlsDao.XLSUpdateTripSheetHeaderForDraftInvoice(tripSheetNumber, updateData)
		if errXls != nil {
			ts.l.Error("ERROR: XLSUpdateTripSheetHeader", errXls)
			errRow.Message = fmt.Sprintf("%v, Error Msg: %v, %v", rowIdx, row[1], errXls)
			errorRows = append(errorRows, errRow)
			continue
		}

		successRow := dtos.Messge{}
		successRow.Message = fmt.Sprintf("Success: %v", row[0])
		successRows = append(successRows, successRow)

		tripSheetNumbers = append(tripSheetNumbers, tripSheetNumber)
		tripSheetNumbersArray = append(tripSheetNumbersArray, tripSheetNumber)
	}

	ts.l.Info("Trips updated successfully", len(tripSheetNumbers), tripSheetNumbers)
	for i, v := range tripSheetNumbers {
		tripSheetNumbers[i] = "'" + v + "'"
	}
	inClauseTripNum := "(" + strings.Join(tripSheetNumbers, ",") + ")"
	draftInvoiceTrips, errT := ts.tripSheetXlsDao.GetTripsByTripSheetNumbers(inClauseTripNum)
	if errT != nil {
		ts.l.Error("ERROR: GetTripsByTripSheetNumbers", errT)
		return nil, errT
	}

	draftInvoiceTotalAmount := 0.0
	if draftInvoiceTrips != nil {

		// Calculating total invoice amount
		for _, trips := range *draftInvoiceTrips {
			ts.l.Info("draftInvoiceTrips", trips.TripSheetID, trips.TripSheetNum, "LoadHire", trips.CustomerPerLoadHire, "Running KM", trips.CustomerRunningKM, "KM Price", trips.CustomerPerKMPrice)

			if trips.TripType == constant.LOCAL_SCHEDULED_TRIP {
				draftInvoiceTotalAmount += (trips.CustomerBaseRate + trips.CustomerKMCost + trips.CustomerToll)
			} else {
				if trips.CustomerPerLoadHire > 0 {
					draftInvoiceTotalAmount += trips.CustomerPerLoadHire
				} else {
					draftInvoiceTotalAmount += (trips.CustomerRunningKM * trips.CustomerPerKMPrice)
				}
			}
			tripSheetIds = append(tripSheetIds, trips.TripSheetID)

			//Update Total Hire Note: Right no need to update the Total Hire Field
		}
		var workStartDate, workEndDate string
		if len(tripOpenTripDateArray) > 0 {
			workStartDate, workEndDate = ts.getMinMaxDates(tripOpenTripDateArray)
		}

		if len(tripSheetIds) != 0 {

			customerName, customerCode, customerId := ts.tripSheetXlsDao.GetCustomerIDAndName(tripSheetIds[0])

			tripSheetNumbers := strings.Join(tripSheetNumbersArray, ",")
			invoiceRef = fmt.Sprintf("%s-%s-%s", time.Now().Format("20060102"), randomLetters(4), randomLetters(4))
			ts.l.Info("draftInvoiceTotalAmount", draftInvoiceTotalAmount, invoiceRef, tripSheetNumbers)

			// Generate Trip Sheet
			draftInvoice := dtos.CreateDraftInvoice{}
			draftInvoice.InvoiceRefID = invoiceRef
			draftInvoice.WorkType = "Spot"
			draftInvoice.WorkStartDate = workStartDate
			draftInvoice.WorkEndDate = workEndDate
			draftInvoice.DocumentDate = utils.GetCurrentDateStr()
			draftInvoice.InvoiceStatus = constant.STATUS_INVOICE_DRAFT
			draftInvoice.TripRef = tripSheetNumbers
			draftInvoice.InvoiceAmount = draftInvoiceTotalAmount
			draftInvoice.CustomerCode = customerCode
			draftInvoice.CustomerId = customerId
			draftInvoice.CustomerName = customerName

			invoiceRowId, errI := ts.tripSheetXlsDao.CreateInvoiceDraft(draftInvoice)
			if errI != nil {
				ts.l.Error("ERROR: CreateInvoinceDraft", errI)
				return nil, errI
			}

			inClauseTripID := strings.Join(int64SliceToStringSlice(tripSheetIds), ",")
			errIn := ts.tripSheetXlsDao.UpdateInvoiceIDToTripHeader(invoiceRowId, inClauseTripID)
			if errIn != nil {
				ts.l.Error("ERROR: UpdateInvoiceIDToTripHeader", errIn)
				return nil, errIn
			}
		}

	}

	roleResponse := dtos.XlsDraftInvoiceMessge{}
	roleResponse.ErrorRows = &errorRows
	roleResponse.Success = &successRows
	roleResponse.ErrorTripsCount = len(errorRows)
	roleResponse.SuccessTripsCount = len(successRows)
	roleResponse.Msg = fmt.Sprintf("Draft Invoice created. trips: %v, Invoice Value: %.2f Draft: %s", len(tripSheetNumbersArray), draftInvoiceTotalAmount, invoiceRef)
	return &roleResponse, nil
}

func getEndDate() time.Time {
	currentDate := time.Now()
	endDate := currentDate.AddDate(0, 0, 2) // Add 2 days
	return endDate
}

func randomLetters(length int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func int64SliceToStringSlice(ints []int64) []string {
	strSlice := make([]string, len(ints))
	for i, v := range ints {
		strSlice[i] = strconv.FormatInt(v, 10) // convert int64 to decimal string
	}
	return strSlice
}

func (ts *TripSheetXls) getMinMaxDates(dates []string) (string, string) {
	layout := "2006-01-02 15:04"
	if len(dates) == 0 {
		return utils.GetCurrentDateStr(), utils.GetCurrentDateStr()
	}

	minDate, err := time.Parse(layout, dates[0])
	if err != nil {
		return utils.GetCurrentDateStr(), utils.GetCurrentDateStr()
	}
	maxDate := minDate

	for _, d := range dates[1:] {
		parsedDate, err := time.Parse(layout, d)
		if err != nil {
			return utils.GetCurrentDateStr(), utils.GetCurrentDateStr()
		}
		if parsedDate.Before(minDate) {
			minDate = parsedDate
		}
		if parsedDate.After(maxDate) {
			maxDate = parsedDate
		}
	}

	return utils.ConvertDateToStringDFyyyyMMdd(minDate), utils.ConvertDateToStringDFyyyyMMdd(maxDate)
}
