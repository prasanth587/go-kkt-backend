package notification

import (
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
)

type NotificationService struct {
	l               *log.Logger
	dbConnMSSQL     *mssqlcon.DBConn
	notificationDao daos.NotificationDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *NotificationService {
	return &NotificationService{
		l:               l,
		dbConnMSSQL:     dbConnMSSQL,
		notificationDao: daos.NewNotificationDao(l, dbConnMSSQL),
	}
}

// CreateNotificationForAll creates a notification for all users in an organization
func (ns *NotificationService) CreateNotificationForAll(orgID int64, notificationType, title, message, relatedEntityType string, relatedEntityID *int64) error {
	req := dtos.NotificationRequest{
		OrgID:             orgID,
		UserID:            nil, // nil means all users
		NotificationType:  notificationType,
		Title:             title,
		Message:           message,
		RelatedEntityType: relatedEntityType,
		RelatedEntityID:   relatedEntityID,
	}

	_, err := ns.notificationDao.CreateNotification(req)
	return err
}

// CreateNotificationForUser creates a notification for a specific user
func (ns *NotificationService) CreateNotificationForUser(orgID int64, userID int64, notificationType, title, message, relatedEntityType string, relatedEntityID *int64) error {
	req := dtos.NotificationRequest{
		OrgID:             orgID,
		UserID:            &userID,
		NotificationType:  notificationType,
		Title:             title,
		Message:           message,
		RelatedEntityType: relatedEntityType,
		RelatedEntityID:   relatedEntityID,
	}

	_, err := ns.notificationDao.CreateNotification(req)
	return err
}

// Trip Sheet Notifications
func (ns *NotificationService) NotifyTripSheetCreated(orgID int64, tripSheetID int64, tripSheetNum string, createdBy int64) error {
	title := "New Trip Sheet Created"
	message := fmt.Sprintf("Trip Sheet %s has been created", tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "trip_created", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyTripSheetStatusChanged(orgID int64, tripSheetID int64, tripSheetNum string, oldStatus, newStatus string) error {
	title := "Trip Sheet Status Updated"
	message := fmt.Sprintf("Trip Sheet %s status changed from %s to %s", tripSheetNum, oldStatus, newStatus)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "trip_status_changed", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyTripSheetClosed(orgID int64, tripSheetID int64, tripSheetNum string) error {
	title := "Trip Sheet Closed"
	message := fmt.Sprintf("Trip Sheet %s has been closed", tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "trip_closed", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyTripSheetInTransit(orgID int64, tripSheetID int64, tripSheetNum string) error {
	title := "Trip Sheet In Transit"
	message := fmt.Sprintf("Trip Sheet %s is now in transit", tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "trip_in_transit", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyTripSheetDelivered(orgID int64, tripSheetID int64, tripSheetNum string) error {
	title := "Trip Sheet Delivered"
	message := fmt.Sprintf("Trip Sheet %s has been delivered", tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "trip_delivered", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyTripSheetCompleted(orgID int64, tripSheetID int64, tripSheetNum string) error {
	title := "Trip Sheet Completed"
	message := fmt.Sprintf("Trip Sheet %s has been completed", tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "trip_completed", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyPODSubmitted(orgID int64, tripSheetID int64, tripSheetNum string) error {
	title := "POD Submitted"
	message := fmt.Sprintf("POD has been submitted for Trip Sheet %s", tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "pod_submitted", title, message, "trip_sheet", tripSheetIDPtr)
}

// Vendor Notifications
func (ns *NotificationService) NotifyVendorCreated(orgID int64, vendorID int64, vendorName string) error {
	title := "New Vendor Added"
	message := fmt.Sprintf("Vendor %s has been added", vendorName)
	vendorIDPtr := &vendorID
	return ns.CreateNotificationForAll(orgID, "vendor_created", title, message, "vendor", vendorIDPtr)
}

func (ns *NotificationService) NotifyVendorUpdated(orgID int64, vendorID int64, vendorName string) error {
	title := "Vendor Updated"
	message := fmt.Sprintf("Vendor %s has been updated", vendorName)
	vendorIDPtr := &vendorID
	return ns.CreateNotificationForAll(orgID, "vendor_updated", title, message, "vendor", vendorIDPtr)
}

// Customer Notifications
func (ns *NotificationService) NotifyCustomerCreated(orgID int64, customerID int64, customerName string) error {
	title := "New Customer Added"
	message := fmt.Sprintf("Customer %s has been added", customerName)
	customerIDPtr := &customerID
	return ns.CreateNotificationForAll(orgID, "customer_created", title, message, "customer", customerIDPtr)
}

func (ns *NotificationService) NotifyCustomerUpdated(orgID int64, customerID int64, customerName string) error {
	title := "Customer Updated"
	message := fmt.Sprintf("Customer %s has been updated", customerName)
	customerIDPtr := &customerID
	return ns.CreateNotificationForAll(orgID, "customer_updated", title, message, "customer", customerIDPtr)
}

// Payment Notifications
func (ns *NotificationService) NotifyPaymentReceived(orgID int64, tripSheetID int64, tripSheetNum string, amount float64) error {
	title := "Payment Received"
	message := fmt.Sprintf("Payment of ₹%.2f received for Trip Sheet %s", amount, tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "payment_received", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyVendorPaymentProcessed(orgID int64, tripSheetID int64, tripSheetNum string, amount float64) error {
	title := "Vendor Payment Processed"
	message := fmt.Sprintf("Payment of ₹%.2f processed for Trip Sheet %s", amount, tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "vendor_payment_processed", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyInvoiceGenerated(orgID int64, tripSheetID int64, tripSheetNum string, invoiceNo string) error {
	title := "Invoice Generated"
	message := fmt.Sprintf("Invoice %s generated for Trip Sheet %s", invoiceNo, tripSheetNum)
	tripSheetIDPtr := &tripSheetID
	return ns.CreateNotificationForAll(orgID, "invoice_generated", title, message, "trip_sheet", tripSheetIDPtr)
}

func (ns *NotificationService) NotifyInvoiceNumberUpdated(orgID int64, invoiceID int64, invoiceNo string) error {
	title := "Invoice Number Updated"
	message := fmt.Sprintf("Invoice number %s has been updated", invoiceNo)
	invoiceIDPtr := &invoiceID
	return ns.CreateNotificationForAll(orgID, "invoice_number_updated", title, message, "invoice", invoiceIDPtr)
}

func (ns *NotificationService) NotifyInvoicePaid(orgID int64, invoiceID int64, invoiceNo string, amount float64) error {
	title := "Invoice Paid"
	message := fmt.Sprintf("Invoice %s has been marked as paid (Amount: ₹%.2f)", invoiceNo, amount)
	invoiceIDPtr := &invoiceID
	return ns.CreateNotificationForAll(orgID, "invoice_paid", title, message, "invoice", invoiceIDPtr)
}

func (ns *NotificationService) NotifyInvoiceCancelled(orgID int64, invoiceID int64, invoiceNo string) error {
	title := "Invoice Cancelled"
	message := fmt.Sprintf("Invoice %s has been cancelled", invoiceNo)
	invoiceIDPtr := &invoiceID
	return ns.CreateNotificationForAll(orgID, "invoice_cancelled", title, message, "invoice", invoiceIDPtr)
}

// Branch Notifications
func (ns *NotificationService) NotifyBranchCreated(orgID int64, branchID int64, branchName string) error {
	title := "New Branch Added"
	message := fmt.Sprintf("Branch %s has been added", branchName)
	branchIDPtr := &branchID
	return ns.CreateNotificationForAll(orgID, "branch_created", title, message, "branch", branchIDPtr)
}

func (ns *NotificationService) NotifyBranchUpdated(orgID int64, branchID int64, branchName string) error {
	title := "Branch Updated"
	message := fmt.Sprintf("Branch %s has been updated", branchName)
	branchIDPtr := &branchID
	return ns.CreateNotificationForAll(orgID, "branch_updated", title, message, "branch", branchIDPtr)
}

// Loading Point Notifications
func (ns *NotificationService) NotifyLoadingPointCreated(orgID int64, loadingPointID int64, cityName string) error {
	title := "New Loading Point Added"
	message := fmt.Sprintf("Loading/Unloading point %s has been added", cityName)
	loadingPointIDPtr := &loadingPointID
	return ns.CreateNotificationForAll(orgID, "loading_point_created", title, message, "loading_point", loadingPointIDPtr)
}

func (ns *NotificationService) NotifyLoadingPointUpdated(orgID int64, loadingPointID int64, cityName string) error {
	title := "Loading Point Updated"
	message := fmt.Sprintf("Loading/Unloading point %s has been updated", cityName)
	loadingPointIDPtr := &loadingPointID
	return ns.CreateNotificationForAll(orgID, "loading_point_updated", title, message, "loading_point", loadingPointIDPtr)
}

// Employee Notifications
func (ns *NotificationService) NotifyEmployeeCreated(orgID int64, employeeID int64, employeeName string) error {
	title := "New Employee Added"
	message := fmt.Sprintf("Employee %s has been added", employeeName)
	employeeIDPtr := &employeeID
	return ns.CreateNotificationForAll(orgID, "employee_created", title, message, "employee", employeeIDPtr)
}

func (ns *NotificationService) NotifyEmployeeUpdated(orgID int64, employeeID int64, employeeName string) error {
	title := "Employee Updated"
	message := fmt.Sprintf("Employee %s has been updated", employeeName)
	employeeIDPtr := &employeeID
	return ns.CreateNotificationForAll(orgID, "employee_updated", title, message, "employee", employeeIDPtr)
}

// Role Notifications
func (ns *NotificationService) NotifyRoleCreated(orgID int64, roleID int64, roleName string) error {
	title := "New Role Created"
	message := fmt.Sprintf("Role %s has been created", roleName)
	roleIDPtr := &roleID
	return ns.CreateNotificationForAll(orgID, "role_created", title, message, "role", roleIDPtr)
}

func (ns *NotificationService) NotifyRoleUpdated(orgID int64, roleID int64, roleName string) error {
	title := "Role Updated"
	message := fmt.Sprintf("Role %s has been updated", roleName)
	roleIDPtr := &roleID
	return ns.CreateNotificationForAll(orgID, "role_updated", title, message, "role", roleIDPtr)
}

// LR Receipt Notifications
func (ns *NotificationService) NotifyLRReceiptCreated(orgID int64, lrID int64, lrNumber string) error {
	title := "LR Receipt Created"
	message := fmt.Sprintf("LR Receipt %s has been created", lrNumber)
	lrIDPtr := &lrID
	return ns.CreateNotificationForAll(orgID, "lr_receipt_created", title, message, "lr_receipt", lrIDPtr)
}

func (ns *NotificationService) NotifyLRReceiptUpdated(orgID int64, lrID int64, lrNumber string) error {
	title := "LR Receipt Updated"
	message := fmt.Sprintf("LR Receipt %s has been updated", lrNumber)
	lrIDPtr := &lrID
	return ns.CreateNotificationForAll(orgID, "lr_receipt_updated", title, message, "lr_receipt", lrIDPtr)
}

// Vehicle Size Notifications
func (ns *NotificationService) NotifyVehicleSizeCreated(orgID int64, vehicleSizeID int64, vehicleSize string) error {
	title := "New Vehicle Size Added"
	message := fmt.Sprintf("Vehicle size %s has been added", vehicleSize)
	vehicleSizeIDPtr := &vehicleSizeID
	return ns.CreateNotificationForAll(orgID, "vehicle_size_created", title, message, "vehicle_size", vehicleSizeIDPtr)
}

func (ns *NotificationService) NotifyVehicleSizeUpdated(orgID int64, vehicleSizeID int64, vehicleSize string) error {
	title := "Vehicle Size Updated"
	message := fmt.Sprintf("Vehicle size %s has been updated", vehicleSize)
	vehicleSizeIDPtr := &vehicleSizeID
	return ns.CreateNotificationForAll(orgID, "vehicle_size_updated", title, message, "vehicle_size", vehicleSizeIDPtr)
}

// GetNotifications retrieves notifications for a user/organization
func (ns *NotificationService) GetNotifications(orgID int64, userID *int64, limit, offset int) (*[]dtos.Notification, int64, error) {
	return ns.notificationDao.GetNotifications(orgID, userID, limit, offset)
}

// GetUnreadCount gets the count of unread notifications
func (ns *NotificationService) GetUnreadCount(orgID int64, userID *int64) (int64, error) {
	return ns.notificationDao.GetUnreadCount(orgID, userID)
}

// MarkAsRead marks specific notifications as read
func (ns *NotificationService) MarkAsRead(notificationIDs []int64, orgID int64, userID *int64) error {
	return ns.notificationDao.MarkAsRead(notificationIDs, orgID, userID)
}

// MarkAllAsRead marks all notifications as read for a user/organization
func (ns *NotificationService) MarkAllAsRead(orgID int64, userID *int64) error {
	return ns.notificationDao.MarkAllAsRead(orgID, userID)
}
