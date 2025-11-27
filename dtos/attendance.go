package dtos

import "time"

type AttendanceStatus string

const (
	AttStatusPresent   AttendanceStatus = "Present"
	AttStatusLeave     AttendanceStatus = "Leave"
	AttStatusAbsent    AttendanceStatus = "Absent"
	AttStatusRequested AttendanceStatus = "Requested"
)

type Attendance struct {
	AttendanceID      uint64    `json:"attendance_id"`
	CheckInDate       time.Time `json:"check_in_date"`
	CheckInDateStr    string    `json:"check_in_date_str"`
	CheckOutDate      time.Time `json:"check_out_date"`
	CheckOutDateStr   string    `json:"check_out_date_str"`
	InTime            string    `json:"in_time"`
	OutTime           string    `json:"out_time"`
	LateTime          int       `json:"late_time"`
	OverTime          int       `json:"over_time"`
	Duration          string    `json:"duration"`
	Hours             int       `json:"hours"`
	Minutes           int       `json:"minutes"`
	CheckINLatitude   float64   `json:"check_in_latitude"`
	CheckINLongitude  float64   `json:"check_in_longitude"`
	CheckINArea       string    `json:"check_in_area"`
	CheckINCity       string    `json:"check_in_city"`
	CheckOUTLatitude  float64   `json:"check_out_latitude"`
	CheckOUTLongitude float64   `json:"check_out_longitude"`
	CheckOUTArea      string    `json:"check_out_area"`
	CheckOUTCity      string    `json:"check_out_city"`
	IsCompleted       bool      `json:"is_completed"`
	ManualCheckedOut  bool      `json:"manual_checked_out"`
	CheckOutByID      uint64    `json:"check_out_by_id"`
	ManualCheckedIn   bool      `json:"manual_checked_in"`
	CheckInByID       uint64    `json:"check_in_by_id"`
	Status            string    `json:"status"`
}

type CreateInAttendance struct {
	EmployeeID       int64            `json:"employee_id"`
	CheckInDateStr   string           `json:"check_in_date_str"`
	CheckInDate      string           `json:"check_in_date"`
	InTime           string           `json:"in_time"`
	CheckInLatitude  float64          `json:"check_in_latitude"`
	CheckInLongitude float64          `json:"check_in_longitude"`
	CheckInByID      uint64           `json:"check_in_by_id"`
	Status           AttendanceStatus `json:"status"`
}
type CreateOutAttendance struct {
	EmployeeID        int64            `json:"employee_id"`
	CheckOutDateStr   string           `json:"check_out_date_str"`
	CheckOutDate      string           `json:"check_out_date"`
	OutTime           string           `json:"out_time"`
	CheckOutLatitude  float64          `json:"check_out_latitude"`
	CheckOutLongitude float64          `json:"check_out_longitude"`
	CheckOutByID      uint64           `json:"check_out_by_id"`
	Status            AttendanceStatus `json:"status"`
}

type AttendanceEntry struct {
	AttendanceId   int64  `json:"attendance_id"`
	EmployeeID     int64  `json:"employee_id"`
	EntryInDateStr string `json:"entry_in_date_str"`
	EntryInDate    string `json:"check_in_date"`
	EntryTime      string `json:"entry_time"`
}

type EmployeeAttendance struct {
	AttendanceId   int64  `json:"attendance_id"`
	EmployeeID     int64  `json:"employee_id"`
	EmployeeIDText string `json:"emp_id_text"`
	EmployeeName   string `json:"employee_name"`
	OutTime        string `json:"out_time"`
	CheckInDateStr string `json:"check_in_date_str"`
	InTime         string `json:"in_time"`
	Status         string `json:"status"`
	IsActive       bool   `json:"is_active"`
}

type EmployeeAttendances struct {
	EmployeeAttendance *[]EmployeeAttendance `json:"employees_attendance"`
	Total              int64                 `json:"total"`
	Limit              int64                 `json:"limit"`
	OffSet             int64                 `json:"offset"`
}
