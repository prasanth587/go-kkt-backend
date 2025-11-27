package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/managepod"
)

func managePod(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/manage/pod", wrapHandler(recoverHandler.ThenFunc(createManagePod)))
	router.GET("/v1/:orgId/manage/pods", wrapHandler(recoverHandler.ThenFunc(getPods)))
	router.POST("/v1/manage/pod/:podId/update", wrapHandler(recoverHandler.ThenFunc(updatePod)))
	router.PUT("/v1/pod/update/:podId/status", wrapHandler(recoverHandler.ThenFunc(updatePodStatus)))
	router.POST("/v1/pod/:tripSheetId/upload", wrapHandler(recoverHandler.ThenFunc(uploadPodImages)))
}

func createManagePod(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	podm := managepod.New(rd.l, rd.dbConnMSSQL)
	var podReq dtos.ManagePodReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&podReq)
	if err != nil {
		rd.l.Error("createManagePod: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := podm.CreateManagePod(podReq)
	if err != nil {
		rd.l.Error("createManagePod error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getPods(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	searchText := strings.TrimSpace(keys.Get("search"))
	podId := keys.Get("pod_id")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getBranches", "limit", limit, "offset", offset, "orgID", orgID, "podId", podId, "searchText", searchText)
	if !isErr {
		return
	}
	podStatus := keys.Get("pod_status")
	tripType := keys.Get("trip_type")
	mp := managepod.New(rd.l, rd.dbConnMSSQL)
	res, err := mp.GetPods(orgID, limit, offset, podId, podStatus, tripType, searchText)
	if err != nil {
		rd.l.Error("GetBranch error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updatePod(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	podM := managepod.New(rd.l, rd.dbConnMSSQL)
	var podObj dtos.UpdateManagePodReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&podObj)
	if err != nil {
		rd.l.Error("updatePod: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	podId, isErr := GetIDFromParams(w, r, "podId")
	if !isErr {
		return
	}
	rd.l.Info("updatePod", "podId", podId)

	res, err := podM.UpdatePod(podId, podObj)
	if err != nil {
		rd.l.Error("updatePod error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updatePodStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	status := keys.Get("status")
	podId, isErr := GetIDFromParams(w, r, "podId")
	rd.l.Info("updatePodStatus", "podId", podId, "status:", status)
	if !isErr {
		return
	}

	podM := managepod.New(rd.l, rd.dbConnMSSQL)
	res, err := podM.UpdatePODStatus(status, podId)
	if err != nil {
		rd.l.Error("UpdatePODStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func uploadPodImages(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	podM := managepod.New(rd.l, rd.dbConnMSSQL)
	tripSheetId, isErr := GetIDFromParams(w, r, "tripSheetId")
	if !isErr {
		return
	}
	keys := r.URL.Query()
	imageFor := keys.Get("imageFor")
	// Get the file from the request
	file, header, err := r.FormFile("image")
	if err != nil {
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	res, err := podM.UploadPodDoc(tripSheetId, imageFor, file, header)
	if err != nil {
		rd.l.Error("UploadPodDoc error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
