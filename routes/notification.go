package routes

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"go-transport-hub/dtos"
	"go-transport-hub/internal/service/notification"
)

func notificationRoutes(router *httprouter.Router, recoverHandler alice.Chain) {
	router.GET("/v1/:orgId/notifications", wrapHandler(recoverHandler.ThenFunc(getNotifications)))
	router.GET("/v1/:orgId/notifications/unread-count", wrapHandler(recoverHandler.ThenFunc(getUnreadCount)))
	router.POST("/v1/:orgId/notifications/mark-read", wrapHandler(recoverHandler.ThenFunc(markAsRead)))
	router.POST("/v1/:orgId/notifications/mark-all-read", wrapHandler(recoverHandler.ThenFunc(markAllAsRead)))
}

func getNotifications(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	if !isErr {
		return
	}

	// Get user ID from query params (optional)
	userIDStr := r.URL.Query().Get("user_id")
	var userID *int64
	if userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}

	// Get limit and offset
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	notificationSvc := notification.New(rd.l, rd.dbConnMSSQL)
	notifications, total, err := notificationSvc.GetNotifications(orgID, userID, limit, offset)
	if err != nil {
		rd.l.Error("GetNotifications error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	unreadCount, err := notificationSvc.GetUnreadCount(orgID, userID)
	if err != nil {
		rd.l.Error("GetUnreadCount error: ", err)
		unreadCount = 0
	}

	response := dtos.NotificationResponse{
		Notifications: *notifications,
		UnreadCount:   unreadCount,
		Total:         total,
	}

	writeJSONStruct(response, http.StatusOK, rd)
}

func getUnreadCount(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	if !isErr {
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	var userID *int64
	if userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}

	notificationSvc := notification.New(rd.l, rd.dbConnMSSQL)
	count, err := notificationSvc.GetUnreadCount(orgID, userID)
	if err != nil {
		rd.l.Error("GetUnreadCount error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	writeJSONStruct(map[string]int64{"unread_count": count}, http.StatusOK, rd)
}

func markAsRead(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	if !isErr {
		return
	}

	var req dtos.MarkAsReadRequest
	if !parseJSON(rd.w, rd.r.Body, &req) {
		rd.l.Error("markAsRead parseJSON error")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	var userID *int64
	if userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}

	notificationSvc := notification.New(rd.l, rd.dbConnMSSQL)
	err := notificationSvc.MarkAsRead(req.NotificationIDs, orgID, userID)
	if err != nil {
		rd.l.Error("MarkAsRead error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	writeJSONMessage("Notifications marked as read", MSG, http.StatusOK, rd)
}

func markAllAsRead(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	orgID, isErr := GetIDFromParams(w, r, "orgId")
	if !isErr {
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	var userID *int64
	if userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}

	notificationSvc := notification.New(rd.l, rd.dbConnMSSQL)
	err := notificationSvc.MarkAllAsRead(orgID, userID)
	if err != nil {
		rd.l.Error("MarkAllAsRead error: ", err)
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	writeJSONMessage("All notifications marked as read", MSG, http.StatusOK, rd)
}
