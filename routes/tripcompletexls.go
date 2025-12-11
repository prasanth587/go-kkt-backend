package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/xuri/excelize/v2"

	"go-transport-hub/internal/service/tripcompletexls"
	"go-transport-hub/utils"
)

func tripcompleteXLS(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/v1/:orgId/trip/sheets/xls", wrapHandler(recoverHandler.ThenFunc(getTripSheetsList)))
	router.GET("/v1/:orgId/downlaod/trip/sheet", wrapHandler(recoverHandler.ThenFunc(downloadTripsByIds)))
	// Tally export route removed
	// router.GET("/v1/:orgId/download/trip/sheet/tally", wrapHandler(recoverHandler.ThenFunc(downloadTripsByIdsTally)))
	router.POST("/v1/update/customer/invoice/xls", wrapHandler(recoverHandler.ThenFunc(updateCustomerPaymantInfoXls)))
}

func paymentAndInvoice(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/v1/:orgId/trip/sheets/xls/invoice", wrapHandler(recoverHandler.ThenFunc(getTripSheetsInvoiceList)))
	router.GET("/v1/:orgId/downlaod/trips/draft/invoice", wrapHandler(recoverHandler.ThenFunc(downloadTripsByIdsForDraftInvoice)))
	router.POST("/v1/upload/draft/invoice/trips/xls", wrapHandler(recoverHandler.ThenFunc(generateDraftInvoiceByTripxls)))
}
func getTripSheetsList(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	status := keys.Get("status")
	tripSearchText := strings.TrimSpace(keys.Get("search"))
	fromDate := keys.Get("start_time")
	toDate := keys.Get("end_time")
	podRequired := keys.Get("pod_required")
	podReceived := keys.Get("pod_received")
	tripSheetId := keys.Get("trip_sheet_id")
	customers := keys.Get("customers")
	vendors := keys.Get("vnedors")
	limit := keys.Get("limit")
	offset := keys.Get("offset")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getTripSheetsList", "fromDate: ", fromDate, "toDate: ", toDate, "orgID: ",
		orgID, "status: ", status, "podRequired", podRequired, "podReceived", podReceived)
	if !isErr {
		return
	}
	bra := tripcompletexls.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetTripSheets(orgID, tripSheetId, fromDate, toDate, status, tripSearchText, fromDate, toDate, podRequired, podReceived, limit, offset, customers, vendors)
	if err != nil {
		rd.l.Error("GetTripSheets error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getTripSheetsInvoiceList(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	status := keys.Get("status")
	tripSearchText := strings.TrimSpace(keys.Get("search"))
	fromDate := keys.Get("start_time")
	toDate := keys.Get("end_time")
	podRequired := keys.Get("pod_required")
	podReceived := keys.Get("pod_received")
	tripSheetId := keys.Get("trip_sheet_id")
	customers := keys.Get("customers")
	vendors := keys.Get("vendors")
	limit := keys.Get("limit")
	offset := keys.Get("offset")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getTripSheetsInvoiceList", "fromDate: ", fromDate, "toDate: ", toDate, "orgID: ",
		orgID, "status: ", status, "podRequired", podRequired, "podReceived", podReceived)
	if !isErr {
		return
	}
	bra := tripcompletexls.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetTripSheetsInvoice(orgID, tripSheetId, fromDate, toDate, status, tripSearchText, fromDate, toDate,
		podRequired, podReceived, limit, offset, customers, vendors)
	if err != nil {
		rd.l.Error("GetTripSheets error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func downloadTripsByIds(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	tripSheetids := keys.Get("trip_sheet_ids")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("orgID", orgID, "tripSheetids", tripSheetids)
	if !isErr {
		return
	}
	dn := tripcompletexls.New(rd.l, rd.dbConnMSSQL)
	res, err := dn.GetTripsByIds(orgID, tripSheetids)
	if err != nil {
		rd.l.Error("GetTripSheets error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	fileName, _ := extractDate()
	fileName = fmt.Sprintf("attachment; filename=KKT_TripSheet_%s.xlsx", fileName)

	// Set the correct headers for XLSX file download
	w.Header().Set("Content-Disposition", fileName)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write the file bytes to the response
	w.Write(res)
}

func extractDate() (string, error) {
	layout := "2006-01-02T15:04"
	dateStr := time.Now().In(utils.TimeLoc()).Format("2006-01-02T15:04")
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return "TripSheets", err
	}
	return t.Format("2006-01-02"), nil
}

func updateCustomerPaymantInfoXls(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	err := r.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		rd.l.Error("updateCustomerPaymantInfoXls error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	tripsheet := tripcompletexls.New(rd.l, rd.dbConnMSSQL)

	file, header, err := r.FormFile("xlsFile")
	if err != nil {
		rd.l.Error("updateCustomerPaymantInfoXls error 1111: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	if !isValidExcelExtension(header.Filename) {
		rd.l.Error("is not a valide file to upload: ", header.Filename)
		writeJSONMessage("is not a valide file to upload. file: "+header.Filename, ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	// Option 1: Use the file directly (if you don't need to read it multiple times)
	xlFile, err := excelize.OpenReader(file)
	if err != nil {
		rd.l.Error("*********** failed to read Excel file:", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusInternalServerError, rd)
		return
	}
	defer xlFile.Close()

	rd.l.Info("no of sheet found the file ", xlFile.SheetCount, header.Filename)

	res, err := tripsheet.UpdateTripSheetXls(xlFile)
	if err != nil {
		rd.l.Error("updateCustomerPaymantInfoXls 222 error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func isValidExcelExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".xlsx" || ext == ".xls"
}

func downloadTripsByIdsForDraftInvoice(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	tripSheetids := keys.Get("trip_sheet_ids")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("orgID", orgID, "tripSheetids", tripSheetids)
	if !isErr {
		return
	}
	dn := tripcompletexls.New(rd.l, rd.dbConnMSSQL)
	res, err := dn.GetTripsForDraftInvoice(orgID, tripSheetids)
	if err != nil {
		rd.l.Error("GetTripSheets error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	fileName, _ := extractDate()
	fileName = fmt.Sprintf("attachment; filename=KKT_CustomerPaymnetTripSheet_%s.xlsx", fileName)

	// Set the correct headers for XLSX file download
	w.Header().Set("Content-Disposition", fileName)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write the file bytes to the response
	w.Write(res)
	//writeJSONStruct(res, http.StatusOK, rd)

}

func downloadTripsByIdsTally(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	tripSheetids := keys.Get("trip_sheet_ids")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("orgID", orgID, "tripSheetids", tripSheetids)
	if !isErr {
		return
	}
	dn := tripcompletexls.New(rd.l, rd.dbConnMSSQL)
	res, err := dn.GetTripsByIdsTally(orgID, tripSheetids)
	if err != nil {
		rd.l.Error("GetTripsByIdsTally error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	fileName, _ := extractDate()
	fileName = fmt.Sprintf("attachment; filename=KKT_TripSheet_Tally_%s.xml", fileName)

	// Set the correct headers for XML file download
	w.Header().Set("Content-Disposition", fileName)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	// Write the XML bytes to the response
	w.Write(res)
}

func generateDraftInvoiceByTripxls(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	err := r.ParseMultipartForm(100 << 20) // 100 MB
	if err != nil {
		rd.l.Error("generateDraftInvoiceByTripxls error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	tripsheet := tripcompletexls.New(rd.l, rd.dbConnMSSQL)

	file, header, err := r.FormFile("xlsFile")
	if err != nil {
		rd.l.Error("generateDraftInvoiceByTripxls error 1111: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	if !isValidExcelExtension(header.Filename) {
		rd.l.Error("is not a valide file to upload: ", header.Filename)
		writeJSONMessage("is not a valide file to upload. file: "+header.Filename, ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	// Option 1: Use the file directly (if you don't need to read it multiple times)
	xlFile, err := excelize.OpenReader(file)
	if err != nil {
		rd.l.Error("*********** failed to read Excel file:", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusInternalServerError, rd)
		return
	}
	defer xlFile.Close()

	rd.l.Info("no of sheet found the file ", xlFile.SheetCount, header.Filename)

	res, err := tripsheet.GenerateDraftInvoiceByTripxls(xlFile)
	if err != nil {
		rd.l.Error("updateCustomerPaymantInfoXls 222 error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
