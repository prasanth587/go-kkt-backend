package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/branch"
)

func branchHub(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/create/branch", wrapHandler(recoverHandler.ThenFunc(createBranch)))
	router.GET("/v1/:orgId/branches", wrapHandler(recoverHandler.ThenFunc(getBranches)))
	router.PUT("/v1/branch/:branchId/update/status", wrapHandler(recoverHandler.ThenFunc(updateBranchActiveStatus)))
	router.POST("/v1/branch/:branchId/update", wrapHandler(recoverHandler.ThenFunc(updateBranch)))
}

func createBranch(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	bra := branch.New(rd.l, rd.dbConnMSSQL)
	var branchReq dtos.BranchReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&branchReq)
	if err != nil {
		rd.l.Error("createBranch: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := bra.CreateBranch(branchReq)
	if err != nil {
		rd.l.Error("createBranch error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func getBranches(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	orgID, isErr := GetIDFromParams(w, r, "orgId")
	rd.l.Info("getBranches", "limit", limit, "offset", offset, "orgID", orgID)
	if !isErr {
		return
	}
	bra := branch.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetBranch(orgID, limit, offset)
	if err != nil {
		rd.l.Error("GetBranch error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateBranch(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ven := branch.New(rd.l, rd.dbConnMSSQL)
	var branchObj dtos.BranchUpdate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&branchObj)
	if err != nil {
		rd.l.Error("updateBranch: request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	branchId, isErr := GetIDFromParams(w, r, "branchId")
	if !isErr {
		return
	}
	rd.l.Info("updateBranch", "branchId", branchId)

	res, err := ven.UpdateBranch(branchId, branchObj)
	if err != nil {
		rd.l.Error("updateBranch error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func updateBranchActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActive := keys.Get("isActive")
	branchId, isErr := GetIDFromParams(w, r, "branchId")
	rd.l.Info("updateBranchActiveStatus", "branchId", branchId, "isActive", isActive)
	if !isErr {
		return
	}

	ven := branch.New(rd.l, rd.dbConnMSSQL)
	res, err := ven.UpdatebranchActiveStatus(isActive, branchId)
	if err != nil {
		rd.l.Error("updateBranchActiveStatus error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
