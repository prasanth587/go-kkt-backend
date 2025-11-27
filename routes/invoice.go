package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/internal/service/invoice"
)

func invoiceHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/all/v1/invoices", wrapHandler(recoverHandler.ThenFunc(getAllInvoices)))
	router.PUT("/v1/invoice/:invoiceId/cancel", wrapHandler(recoverHandler.ThenFunc(cancelDraftInvoiceStatus)))
	router.PUT("/v1/invoice/:invoiceId/update/number", wrapHandler(recoverHandler.ThenFunc(updateInvoiceNumber)))
	router.PUT("/v1/invoice/:invoiceId/update/paid", wrapHandler(recoverHandler.ThenFunc(updatePaid)))
	router.GET("/invoice/v1/:invoiceId/downlaod/xls", wrapHandler(recoverHandler.ThenFunc(downloadInvoiceXls)))
	router.GET("/invoice/v1/:invoiceId/downlaod/pdf", wrapHandler(recoverHandler.ThenFunc(downloadInvoicePDF)))
	router.GET("/invoice/v1/:invoiceId/info", wrapHandler(recoverHandler.ThenFunc(invoiceInfo)))
}

func getAllInvoices(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	status := strings.TrimSpace(keys.Get("status"))
	customerId := strings.TrimSpace(keys.Get("customer_id"))
	searchText := strings.TrimSpace(keys.Get("search"))

	rd.l.Info("getAllInvoices", "limit:", limit, "offset:", offset, "searchText:", searchText)

	inv := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := inv.GetAllInvoices(limit, offset, searchText, status, customerId)
	if err != nil {
		rd.l.Error("GetLoadingPoints error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

// func updateLoading(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	load := loadingpoint.New(rd.l, rd.dbConnMSSQL)
// 	var loadingObj dtos.LoadingPointUpdate
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&loadingObj)
// 	if err != nil {
// 		rd.l.Error("updateLoading: request error: ", err.Error())
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	loadingPointId, isErr := GetIDFromParams(w, r, "loadingPointId")
// 	if !isErr {
// 		return
// 	}
// 	rd.l.Info("updateLoading", "loadingPointId", loadingPointId)

// 	res, err := load.UpdateLoadingPoint(loadingPointId, loadingObj)
// 	if err != nil {
// 		rd.l.Error("updateLoading error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

func cancelDraftInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	invoiceId, isErr := GetIDFromParams(w, r, "invoiceId")
	rd.l.Info("cancelDraftInvoiceStatus", "invoiceId", invoiceId)
	if !isErr {
		return
	}

	ven := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.CancelDraftInvoice(invoiceId)
	if err != nil {
		rd.l.Error("cancelDraftInvoiceStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateInvoiceNumber(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	invoiceNumber := keys.Get("invoice_number")
	invoiceId, isErr := GetIDFromParams(w, r, "invoiceId")
	rd.l.Info("updateInvoiceNumber", "invoiceId", invoiceId, "invoiceNumber", invoiceNumber)
	if !isErr {
		return
	}

	ven := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.UpdateInvoiceNumber(invoiceId, invoiceNumber)
	if err != nil {
		rd.l.Error("UpdateInvoiceNumber error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updatePaid(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	invoicePaidDate := keys.Get("invoice_paid_date")
	transactionId := keys.Get("transaction_id")
	invoiceId, isErr := GetIDFromParams(w, r, "invoiceId")
	rd.l.Info("UpdateInvoicePaid", "invoiceId", invoiceId, "invoicePaidDate", invoicePaidDate, "transactionId", transactionId)
	if !isErr {
		return
	}

	ven := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.UpdateInvoicePaid(invoiceId, invoicePaidDate, transactionId)
	if err != nil {
		rd.l.Error("UpdateInvoicePaid error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func downloadInvoiceXls(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	//keys := r.URL.Query()

	invoiceId, isErr := GetIDFromParams(w, r, "invoiceId")
	rd.l.Info("downlaod invoice", invoiceId)
	if !isErr {
		return
	}
	dn := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := dn.DownlaodInvoiceXls(invoiceId)
	if err != nil {
		rd.l.Error("DownlaodInvoiceXls error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	fileName, _ := extractDate()
	fileName = fmt.Sprintf("attachment; filename=KKT_Invoice_%s.xlsx", fileName)

	// Set the correct headers for XLSX file download
	w.Header().Set("Content-Disposition", fileName)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write the file bytes to the response
	w.Write(res)
}

func downloadInvoicePDF(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	//keys := r.URL.Query()

	invoiceId, isErr := GetIDFromParams(w, r, "invoiceId")
	rd.l.Info("downlaod invoice", invoiceId)
	if !isErr {
		return
	}
	inv := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := inv.GetInvoicePDFInfo(invoiceId)
	if err != nil {
		rd.l.Error("GetLoadingPoints error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func invoiceInfo(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	invoiceId, isErr := GetIDFromParams(w, r, "invoiceId")

	rd.l.Info("invoiceInfo", "invoiceId", invoiceId)
	if !isErr {
		return
	}

	inv := invoice.New(rd.l, rd.dbConnMSSQL)
	res, err := inv.GetInvoiceInfo(invoiceId)
	if err != nil {
		rd.l.Error("GetLoadingPoints error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
