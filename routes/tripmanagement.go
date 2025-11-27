package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/tripmanagement"
)

func tripManagementHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/trip/sheet", wrapHandler(recoverHandler.ThenFunc(createTripSheet)))
	router.GET("/v1/:orgId/trip/sheets", wrapHandler(recoverHandler.ThenFunc(getTripSheets)))
	router.POST("/v1/trip/sheet/update/:tripSheetId", wrapHandler(recoverHandler.ThenFunc(updateTripSheet)))
	router.POST("/v1/trip/sheet/doc/upload", wrapHandler(recoverHandler.ThenFunc(uploadTripSheetImages)))
	router.PUT("/v1/trip/sheet/:tripSheetId/cancel", wrapHandler(recoverHandler.ThenFunc(cancelTripSheet)))
	router.GET("/v1/:orgId/trip/sheets/stats", wrapHandler(recoverHandler.ThenFunc(getTripStats)))

	router.GET("/trip/sheet/v1/:tripSheetId/challan/info", wrapHandler(recoverHandler.ThenFunc(tripSheetChallanInfo)))
}

func cancelTripSheet(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	tripSheetId, isErr := GetIDFromParams(w, r, "tripSheetId")
	rd.l.Info("cancelTripSheet", "trip to candel. ID:", tripSheetId)
	if !isErr {
		return
	}

	trp := tripmanagement.New(rd.l, rd.dbConnMSSQL)
	res, err := trp.CancelTripSheet(tripSheetId)
	if err != nil {
		rd.l.Error("cancelTripSheet error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func createTripSheet(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	cus := tripmanagement.New(rd.l, rd.dbConnMSSQL)
	var tripReq dtos.CreateTripSheetHeader
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tripReq)
	if err != nil {
		rd.l.Error("createcustomer: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := cus.CreateTripSheetHeader(tripReq)
	if err != nil {
		rd.l.Error("CreateCustomer error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getTripSheets(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	status := keys.Get("status")
	tripSearchText := strings.TrimSpace(keys.Get("search"))
	fromDate := keys.Get("fromDate")
	toDate := keys.Get("toDate")
	podRequired := keys.Get("pod_required")
	podReceived := keys.Get("pod_received")
	tripSheetId := keys.Get("trip_sheet_id")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getTripSheetNumber", "fromDate: ", fromDate, "toDate: ", toDate, "orgID: ",
		orgID, "status: ", status, "podRequired", podRequired, "podReceived", podReceived)
	if !isErr {
		return
	}
	bra := tripmanagement.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetTripSheets(orgID, tripSheetId, limit, offset, status, tripSearchText, fromDate, toDate, podRequired, podReceived)
	if err != nil {
		rd.l.Error("GetTripSheets error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadTripSheetImages(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	tripsheet := tripmanagement.New(rd.l, rd.dbConnMSSQL)

	keys := r.URL.Query()
	imageFor := keys.Get("imageFor")
	tripSheetNumber := keys.Get("tripSheetNumber")

	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		rd.l.Error("uploadTripSheetImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := tripsheet.UploadTripSheetImages(imageFor, tripSheetNumber, file, header)
	if err != nil {
		rd.l.Error("UploadCustomerImages error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateTripSheet(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	load := tripmanagement.New(rd.l, rd.dbConnMSSQL)
	var updateTripSheet dtos.UpdateTripSheetHeader
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updateTripSheet)
	if err != nil {
		rd.l.Error("UpdateTripSheetHeader: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	tripSheetId, isErr := GetIDFromParams(w, r, "tripSheetId")
	if !isErr {
		return
	}
	rd.l.Info("UpdateTripSheetHeader", "tripSheetId", tripSheetId)

	res, err := load.UpdateTripSheetHeader(tripSheetId, updateTripSheet)
	if err != nil {
		rd.l.Error("UpdateTripSheetHeader error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getTripStats(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	fromDate := keys.Get("from_date")
	toDate := keys.Get("to_date")
	podRequired := keys.Get("pod_required")
	podReceived := keys.Get("pod_received")
	status := keys.Get("status")
	tripSheetId := keys.Get("trip_sheet_id")
	tripSearchText := strings.TrimSpace(keys.Get("search"))
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getTripSheetNumber", "fromDate: ", fromDate, "toDate: ", toDate, "orgID: ",
		orgID, "status: ", status, "podRequired", podRequired, "podReceived", podReceived)
	if !isErr {
		return
	}
	bra := tripmanagement.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetTripStats(orgID, tripSheetId, fromDate, toDate, status, tripSearchText, podRequired, podReceived)
	if err != nil {
		rd.l.Error("GetTripSheets error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func tripSheetChallanInfo(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	tripSheetId, isErr := GetIDFromParams(w, r, "tripSheetId")
	rd.l.Info("tripSheetChallanInfo", "tripSheetId: ", tripSheetId)
	if !isErr {
		return
	}
	keys := r.URL.Query()
	loginId := keys.Get("login_id")
	bra := tripmanagement.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetTripSheetChallanInfo(tripSheetId, loginId)
	if err != nil {
		rd.l.Error("GetTripSheetChallanInfo error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
