package schema

import (
	"time"
)

type Organisation struct {
	OrgId        int64     `json:"org_id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	DomainName   string    `json:"domain_name"`
	EmailId      string    `json:"email_id"`
	ContactName  string    `json:"contact_name"`
	ContactNo    string    `json:"contact_no"`
	IsActive     bool      `json:"is_active"`
	LogoPath     string    `json:"logo_path"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int64     `json:"version"`
}

type UserLogin struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	MobileNo     string    `json:"mobile_no"`
	EmailId      string    `json:"email_id"`
	IsActive     bool      `json:"is_active"`
	ContactName  string    `json:"contact_name"`
	ContactNo    string    `json:"contact_no"`
	Password     string    `json:"password"`
	PasswordStr  string    `json:"password_str"`
	AccessToken  string    `json:"access_token"`
	RoleID       int64     `json:"role_id"`
	RoleName     string    `json:"role_name"`
	OrgId        int64     `json:"org_id"`
	EmployeeId   int64     `json:"employee_id"`
	IsSuperAdmin int       `json:"is_super_admin"`
	IsAdmin      int       `json:"is_admin"`
	LoginType    string    `json:"login_type"`
	LastLogin    time.Time `json:"last_login"`
	FacilityID   int64     `json:"facility_id"`
	Version      int64     `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	RoleId      int64     `json:"role_id"`     // Unique identifier for the role
	RoleCode    string    `json:"role_code"`   // Code for the role
	RoleName    string    `json:"role_name"`   // Name of the role
	Description string    `json:"description"` // Description of the role
	OrgId       int64     `json:"org_id"`      // Organization ID associated with the role
	IsActive    int       `json:"is_active"`   // Indicates if the role is active
	Version     int64     `json:"version"`     // Version of the role record
	CreatedAt   time.Time `json:"created_at"`  // Timestamp when the role was created
	UpdatedAt   time.Time `json:"updated_at"`  // Timestamp when the role was last updated
}

type Employee struct {
	EmpID         int64     `json:"emp_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	EmployeeCode  string    `json:"employee_code"`
	MobileNo      string    `json:"mobile_no"`
	EmailID       string    `json:"email_id"`
	RoleID        int64     `json:"role_id"`
	DOB           time.Time `json:"dob" db:"dob"`
	Gender        string    `json:"gender"`
	AadharNo      string    `json:"aadhar_no"`
	AccessNo      string    `json:"access_no"`
	IsActive      int64     `json:"is_active"`
	AccessToken   string    `json:"access_token"`
	JoiningDate   time.Time `json:"joining_date"`
	RelievingDate time.Time `json:"relieving_date"`
	AddressLine1  string    `json:"address_line1"`
	AddressLine2  string    `json:"address_line2"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	IsSuperAdmin  int64     `json:"is_super_admin" `
	IsAdmin       int64     `json:"is_admin"`
	OrgID         int64     `json:"org_id"`
	LoginType     string    `json:"login_type" `
	Image         string    `json:"image"`
	FacilityID    int64     `json:"facility_id"`
	Version       int64     `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TruckTypes struct {
	TruckTypeId int64  `json:"truck_type_id"`
	Category    string `json:"category"`
	TruckType   string `json:"truck_type"`
}

type PermissionLabel struct {
	PermisstionLabelID int64  `json:"permisstion_label_id"`
	PermisstionLabel   string `json:"permisstion_label"`
	Description        string `json:"description"`
}

type WebsiteScreens struct {
	WebsiteScreenID int64  `json:"website_screen_id"`
	ScreenName      string `json:"screen_name"`
	Description     string `json:"description"`
}

type AppConfig struct {
	AppConfigID int64  `json:"app_config_id"`
	ConfigCode  string `json:"config_code"`
	ConfigName  string `json:"config_name"`
	ConfigValue string `json:"value"`
}

type EmpAttendance struct {
	CheckINDate string `json:"check_in_date_str"`
	InTime      string `json:"in_time"`
	OutTime     string `json:"out_time"`
}
