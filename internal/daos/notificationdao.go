package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type NotificationDao interface {
	CreateNotification(notification dtos.NotificationRequest) (int64, error)
	GetNotifications(orgID int64, userID *int64, limit, offset int) (*[]dtos.Notification, int64, error)
	GetUnreadCount(orgID int64, userID *int64) (int64, error)
	MarkAsRead(notificationIDs []int64, orgID int64, userID *int64) error
	MarkAllAsRead(orgID int64, userID *int64) error
}

type NotificationObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewNotificationDao(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) NotificationDao {
	return &NotificationObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

func (n *NotificationObj) CreateNotification(notification dtos.NotificationRequest) (int64, error) {
	userIDStr := "NULL"
	if notification.UserID != nil {
		userIDStr = fmt.Sprintf("%d", *notification.UserID)
	}

	relatedEntityIDStr := "NULL"
	if notification.RelatedEntityID != nil {
		relatedEntityIDStr = fmt.Sprintf("%d", *notification.RelatedEntityID)
	}

	query := fmt.Sprintf(`INSERT INTO notifications 
		(org_id, user_id, notification_type, title, message, related_entity_type, related_entity_id) 
		VALUES (%d, %s, '%s', '%s', '%s', '%s', %s)`,
		notification.OrgID,
		userIDStr,
		notification.NotificationType,
		notification.Title,
		notification.Message,
		notification.RelatedEntityType,
		relatedEntityIDStr,
	)

	n.l.Info("CreateNotification query: ", query)

	result, err := n.dbConnMSSQL.GetQueryer().Exec(query)
	if err != nil {
		n.l.Error("Error creating notification: ", err)
		return 0, err
	}

	notificationID, err := result.LastInsertId()
	if err != nil {
		n.l.Error("Error getting last insert ID: ", err)
		return 0, err
	}

	return notificationID, nil
}

func (n *NotificationObj) GetNotifications(orgID int64, userID *int64, limit, offset int) (*[]dtos.Notification, int64, error) {
	userFilter := ""
	if userID != nil {
		userFilter = fmt.Sprintf("AND (user_id = %d OR user_id IS NULL)", *userID)
	} else {
		userFilter = "AND user_id IS NULL"
	}

	query := fmt.Sprintf(`SELECT 
		notification_id, org_id, user_id, notification_type, title, message, 
		related_entity_type, related_entity_id, is_read, created_at
		FROM notifications 
		WHERE org_id = %d %s
		ORDER BY created_at DESC
		LIMIT %d OFFSET %d`,
		orgID, userFilter, limit, offset)

	n.l.Info("GetNotifications query: ", query)

	rows, err := n.dbConnMSSQL.GetQueryer().Query(query)
	if err != nil {
		n.l.Error("Error fetching notifications: ", err)
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []dtos.Notification
	for rows.Next() {
		var notification dtos.Notification
		var userID, relatedEntityID sql.NullInt64

		err := rows.Scan(
			&notification.NotificationID,
			&notification.OrgID,
			&userID,
			&notification.NotificationType,
			&notification.Title,
			&notification.Message,
			&notification.RelatedEntityType,
			&relatedEntityID,
			&notification.IsRead,
			&notification.CreatedAt,
		)
		if err != nil {
			n.l.Error("Error scanning notification: ", err)
			continue
		}

		if userID.Valid {
			notification.UserID = &userID.Int64
		}
		if relatedEntityID.Valid {
			notification.RelatedEntityID = &relatedEntityID.Int64
		}

		notifications = append(notifications, notification)
	}

	// Get total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM notifications WHERE org_id = %d %s`, orgID, userFilter)
	var total int64
	err = n.dbConnMSSQL.GetQueryer().QueryRow(countQuery).Scan(&total)
	if err != nil {
		n.l.Error("Error getting notification count: ", err)
		total = int64(len(notifications))
	}

	return &notifications, total, nil
}

func (n *NotificationObj) GetUnreadCount(orgID int64, userID *int64) (int64, error) {
	userFilter := ""
	if userID != nil {
		userFilter = fmt.Sprintf("AND (user_id = %d OR user_id IS NULL)", *userID)
	} else {
		userFilter = "AND user_id IS NULL"
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM notifications 
		WHERE org_id = %d %s AND is_read = FALSE`,
		orgID, userFilter)

	var count int64
	err := n.dbConnMSSQL.GetQueryer().QueryRow(query).Scan(&count)
	if err != nil {
		n.l.Error("Error getting unread count: ", err)
		return 0, err
	}

	return count, nil
}

func (n *NotificationObj) MarkAsRead(notificationIDs []int64, orgID int64, userID *int64) error {
	if len(notificationIDs) == 0 {
		return nil
	}

	idsStr := ""
	for i, id := range notificationIDs {
		if i > 0 {
			idsStr += ","
		}
		idsStr += fmt.Sprintf("%d", id)
	}

	userFilter := ""
	if userID != nil {
		userFilter = fmt.Sprintf("AND (user_id = %d OR user_id IS NULL)", *userID)
	}

	query := fmt.Sprintf(`UPDATE notifications 
		SET is_read = TRUE 
		WHERE notification_id IN (%s) AND org_id = %d %s`,
		idsStr, orgID, userFilter)

	n.l.Info("MarkAsRead query: ", query)

	_, err := n.dbConnMSSQL.GetQueryer().Exec(query)
	if err != nil {
		n.l.Error("Error marking notifications as read: ", err)
		return err
	}

	return nil
}

func (n *NotificationObj) MarkAllAsRead(orgID int64, userID *int64) error {
	userFilter := ""
	if userID != nil {
		userFilter = fmt.Sprintf("AND (user_id = %d OR user_id IS NULL)", *userID)
	}

	query := fmt.Sprintf(`UPDATE notifications 
		SET is_read = TRUE 
		WHERE org_id = %d %s AND is_read = FALSE`,
		orgID, userFilter)

	n.l.Info("MarkAllAsRead query: ", query)

	_, err := n.dbConnMSSQL.GetQueryer().Exec(query)
	if err != nil {
		n.l.Error("Error marking all notifications as read: ", err)
		return err
	}

	return nil
}
