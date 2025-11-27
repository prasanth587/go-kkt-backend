package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/lrreceipt"
)

func lrReceipt(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/lr/receipt", wrapHandler(recoverHandler.ThenFunc(createLRReceipt)))
	router.POST("/v1/lr/receipt/:lrId/update", wrapHandler(recoverHandler.ThenFunc(updateLR)))
	router.GET("/v1/:orgId/lr/records", wrapHandler(recoverHandler.ThenFunc(getLRRecords)))

	// router.PUT("/v1/manage/pod/:podId/update", wrapHandler(recoverHandler.ThenFunc(updatePod)))
	//router.PUT("/v1/pod/update/:podId/status", wrapHandler(recoverHandler.ThenFunc(updatePodStatus)))
	// router.POST("/v1/pod/:tripSheetId/upload", wrapHandler(recoverHandler.ThenFunc(uploadPodImages)))
}

func createLRReceipt(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	lrm := lrreceipt.New(rd.l, rd.dbConnMSSQL)
	var lrReq dtos.LRReceiptReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&lrReq)
	if err != nil {
		rd.l.Error("createManagePod: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := lrm.CreateLRReceipt(lrReq)
	if err != nil {
		rd.l.Error("createManagePod error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getLRRecords(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getLRRecords", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}
	lrId := keys.Get("lr_id")
	tripSheetNum := keys.Get("trip_sheet_num")
	lrNumber := keys.Get("lr_number")
	tripDate := keys.Get("trip_date")
	searchText := strings.TrimSpace(keys.Get("search"))
	lrG := lrreceipt.New(rd.l, rd.dbConnMSSQL)
	res, err := lrG.GetLRRecords(orgID, limit, offset, lrId, tripSheetNum, lrNumber, tripDate, searchText)
	if err != nil {
		rd.l.Error("GetLRRecords error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateLR(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	lrR := lrreceipt.New(rd.l, rd.dbConnMSSQL)
	var lrObj dtos.LRReceiptUpdateReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&lrObj)
	if err != nil {
		rd.l.Error("updateLR: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	lrId, isErr := GetIDFromParams(w, r, "lrId")
	if !isErr {
		return
	}
	rd.l.Info("updateLR", "lrId", lrId)

	res, err := lrR.UpdateLR(lrId, lrObj)
	if err != nil {
		rd.l.Error("updatePod error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

// func updatePodStatus(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	keys := r.URL.Query()
// 	status := keys.Get("status")
// 	podId, isErr := GetIDFromParams(w, r, "podId")
// 	rd.l.Info("updatePodStatus", "podId", podId, "status:", status)
// 	if !isErr {
// 		return
// 	}

// 	podM := managepod.New(rd.l, rd.dbConnMSSQL)
// 	res, err := podM.UpdatePODStatus(status, podId)
// 	if err != nil {
// 		rd.l.Error("UpdatePODStatus error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }

// func uploadPodImages(w http.ResponseWriter, r *http.Request) {
// 	rd := logAndGetContext(w, r)
// 	podM := managepod.New(rd.l, rd.dbConnMSSQL)
// 	tripSheetId, isErr := GetIDFromParams(w, r, "tripSheetId")
// 	if !isErr {
// 		return
// 	}
// 	keys := r.URL.Query()
// 	imageFor := keys.Get("imageFor")
// 	// Get the file from the request
// 	file, header, err := r.FormFile("image")
// 	if err != nil {
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	defer file.Close()

// 	res, err := podM.UploadPodDoc(tripSheetId, imageFor, file, header)
// 	if err != nil {
// 		rd.l.Error("UploadPodDoc error: ", err)
// 		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
// 		return
// 	}
// 	writeJSONStruct(res, http.StatusOK, rd)
// }
