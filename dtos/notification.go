package dtos

type Notification struct {
	NotificationID    int64  `json:"notification_id"`
	OrgID             int64  `json:"org_id"`
	UserID            *int64 `json:"user_id"`           // null means notification for all users in org
	NotificationType  string `json:"notification_type"` // trip_created, trip_status_changed, vendor_created, etc.
	Title             string `json:"title"`
	Message           string `json:"message"`
	RelatedEntityType string `json:"related_entity_type"` // trip_sheet, vendor, customer, etc.
	RelatedEntityID   *int64 `json:"related_entity_id"`
	IsRead            bool   `json:"is_read"`
	CreatedAt         string `json:"created_at"`
}

type NotificationRequest struct {
	OrgID             int64  `json:"org_id"`
	UserID            *int64 `json:"user_id"` // null for all users
	NotificationType  string `json:"notification_type"`
	Title             string `json:"title"`
	Message           string `json:"message"`
	RelatedEntityType string `json:"related_entity_type"`
	RelatedEntityID   *int64 `json:"related_entity_id"`
}

type NotificationResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int64          `json:"unread_count"`
	Total         int64          `json:"total"`
}

type MarkAsReadRequest struct {
	NotificationIDs []int64 `json:"notification_ids"`
}
