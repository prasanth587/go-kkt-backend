package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/dtos/schema"
	"go-transport-hub/internal/service/userlogin"
)

func login(router *httprouter.Router, recoverHandler alice.Chain) {
	router.POST("/v1/adminLogin", wrapHandler(recoverHandler.ThenFunc(Login)))
	router.POST("/v1/user/create", wrapHandler(recoverHandler.ThenFunc(CreateUser)))
	router.POST("/user/v1/:userId/update", wrapHandler(recoverHandler.ThenFunc(UpdateUser)))
	router.PUT("/user/v1/:userId/update/status", wrapHandler(recoverHandler.ThenFunc(UpdateUserActiveStatus)))
	router.GET("/v1/:orgId/users/logins", wrapHandler(recoverHandler.ThenFunc(GetUsers)))
	router.POST("/forgot/v1/password", wrapHandler(recoverHandler.ThenFunc(UpdatePassword)))
}

func Login(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ul := userlogin.New(rd.l, rd.dbConnMSSQL)
	var reqbody dtos.AdminLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqbody)
	if err != nil {
		rd.l.Error("request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := ul.UserLogin(reqbody)
	if err != nil {
		rd.l.Error("UserLogin error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ul := userlogin.New(rd.l, rd.dbConnMSSQL)
	var reqbody schema.UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqbody)
	if err != nil {
		rd.l.Error("request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := ul.CreateUser(reqbody)
	if err != nil {
		rd.l.Error("Create UserLogin error:", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ul := userlogin.New(rd.l, rd.dbConnMSSQL)
	var reqbody schema.UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqbody)
	if err != nil {
		rd.l.Error("request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := ul.UpdateUser(reqbody)
	if err != nil {
		rd.l.Error("Update UserLogin error:", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func UpdateUserActiveStatus(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	isActiveStr := keys.Get("isActive")
	userId, isErr := GetIDFromParams(w, r, "userId")
	if !isErr {
		return
	}

	isActive := 0
	if isActiveStr == "1" {
		isActive = 1
	}

	ul := userlogin.New(rd.l, rd.dbConnMSSQL)
	res, err := ul.UpdateUserActiveStatus(userId, isActive)
	if err != nil {
		rd.l.Error("UpdateUserActiveStatus error:", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	keys := r.URL.Query()
	limit := keys.Get("limit")
	offset := keys.Get("offset")
	searchText := strings.TrimSpace(keys.Get("search"))
	userId := keys.Get("user_id")

	rd.l.Info("GetUsers", "limit:", limit, "offset:", offset, "searchText:", searchText, "user_id", userId)
	bra := userlogin.New(rd.l, rd.dbConnMSSQL)
	res, err := bra.GetUserList(limit, offset, searchText, userId)
	if err != nil {
		rd.l.Error("GetUsers error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)
	ul := userlogin.New(rd.l, rd.dbConnMSSQL)
	var reqbody dtos.UpdatePassword
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqbody)
	if err != nil {
		rd.l.Error("request error: ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	res, err := ul.UpdatePassword(reqbody)
	if err != nil {
		rd.l.Error("UpdatePassword error:", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}
	writeJSONStruct(res, http.StatusOK, rd)
}
