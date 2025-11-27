package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/loadingpoint"
)

func loadingUnloadingHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/loading/point", wrapHandler(recoverHandler.ThenFunc(createLoadingPoint)))
	router.GET("/v1/:orgId/loading/points", wrapHandler(recoverHandler.ThenFunc(getLoadingPoints)))
	router.PUT("/v1/loading/:loadingPointId/update/status", wrapHandler(recoverHandler.ThenFunc(updateLoadingActiveStatus)))
	router.POST("/v1/loading/:loadingPointId/update", wrapHandler(recoverHandler.ThenFunc(updateLoading)))
}

func createLoadingPoint(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	loadingP := loadingpoint.New(rd.l, rd.dbConnMSSQL)
	var loadingReq dtos.LoadingPointReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loadingReq)
	if err != nil {
		rd.l.Error("createLoadingPoint: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := loadingP.CreateLoadingPoint(loadingReq)
	if err != nil {
		rd.l.Error("createLoading error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getLoadingPoints(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	searchText := strings.TrimSpace(keys.Get("search"))
	loadingPointId := keys.Get("loading_point_id")

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getLoadinges", "limit:", limit, "offset:", offset, "orgID:", orgID, "searchText:", searchText, "loadingPointId:", loadingPointId)
	if !isErr {
		return
	}
	bra := loadingpoint.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetLoadingPoint(orgID, limit, offset, searchText, loadingPointId)
	if err != nil {
		rd.l.Error("GetLoadingPoints error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateLoading(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	load := loadingpoint.New(rd.l, rd.dbConnMSSQL)
	var loadingObj dtos.LoadingPointUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loadingObj)
	if err != nil {
		rd.l.Error("updateLoading: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	loadingPointId, isErr := GetIDFromParams(w, r, "loadingPointId")
	if !isErr {
		return
	}
	rd.l.Info("updateLoading", "loadingPointId", loadingPointId)

	res, err := load.UpdateLoadingPoint(loadingPointId, loadingObj)
	if err != nil {
		rd.l.Error("updateLoading error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateLoadingActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	loadingPointId, isErr := GetIDFromParams(w, r, "loadingPointId")
	rd.l.Info("updateLoadingActiveStatus", "loadingPointId", loadingPointId, "isActive", isActive)
	if !isErr {
		return
	}

	ven := loadingpoint.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.UpdateloadingpointActiveStatus(isActive, loadingPointId)
	if err != nil {
		rd.l.Error("updateLoadingActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
