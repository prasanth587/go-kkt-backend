package attendance

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
)

type AttendanceObj struct {
	l             *log.Logger
	dbConnMSSQL   *mssqlcon.DBConn
	attendanceDao daos.AttendanceDao
}

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *AttendanceObj {
	return &AttendanceObj{
		l:             l,
		dbConnMSSQL:   dbConnMSSQL,
		attendanceDao: daos.NewAttendanceObj(l, dbConnMSSQL),
	}
}

func (att *AttendanceObj) CreateInAttendance(attendance dtos.CreateInAttendance, employeeId int64) (*dtos.Messge, error) {

	if attendance.CheckInDateStr == "" {
		att.l.Error("Error CheckInDateStr should not empty")
		return nil, errors.New("CheckInDateStr should not empty")
	}
	if attendance.InTime == "" {
		att.l.Error("Error InTime: InTime should not empty")
		return nil, errors.New("InTime should not empty")
	}
	attendance.Status = dtos.AttStatusPresent
	attendance.EmployeeID = employeeId

	attendanceId := att.attendanceDao.CheckAttendanceForEmployee(employeeId, attendance.CheckInDateStr)

	fullDateTime := fmt.Sprintf("%s %s", attendance.CheckInDateStr, attendance.InTime)

	// Layout must exactly match the format of fullDateTime
	layout := "2006-01-02 15:04"

	// Parse full datetime string into time.Time
	checkInDate, err := time.Parse(layout, fullDateTime)
	if err != nil {
		att.l.Error("attendance Error parsing date/time: ", attendance.EmployeeID, err)
	}

	attendance.CheckInDate = checkInDate.Format("2006-01-02 15:04:05")

	if attendanceId == 0 {
		attendanceId, err = att.attendanceDao.CreateInAttendance(attendance)
		if err != nil {
			att.l.Error("attendance not saved ", attendance.EmployeeID, err)
			return nil, err
		}
	}
	att.l.Info("Attendance created successfully! : ", attendanceId, attendance.EmployeeID)
	entry := dtos.AttendanceEntry{}
	entry.EntryInDate = attendance.CheckInDate
	entry.EntryInDateStr = attendance.CheckInDateStr
	entry.EntryTime = attendance.InTime
	entry.EmployeeID = attendance.EmployeeID
	entry.AttendanceId = attendanceId
	err = att.attendanceDao.CreateAttendanceQuery(entry)
	if err != nil {
		att.l.Error("CreateAttendanceQuery not saved ", attendance.EmployeeID, err)
		//return nil, err
	}

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("Attendance saved successfully: time: %v", attendance.InTime)
	return &response, nil
}

func (att *AttendanceObj) CreateOutAttendance(attendance dtos.CreateOutAttendance, employeeId int64) (*dtos.Messge, error) {

	if attendance.CheckOutDateStr == "" {
		att.l.Error("Error CheckOutDateStr should not empty")
		return nil, errors.New("CheckOutDateStr should not empty")
	}
	if attendance.OutTime == "" {
		att.l.Error("Error OutTime: OutTime should not empty")
		return nil, errors.New("OutTime should not empty")
	}
	attendance.Status = dtos.AttStatusPresent
	attendance.EmployeeID = employeeId

	attendanceId := att.attendanceDao.CheckAttendanceForEmployee(employeeId, attendance.CheckOutDateStr)

	fullDateTime := fmt.Sprintf("%s %s", attendance.CheckOutDateStr, attendance.OutTime)

	// Layout must exactly match the format of fullDateTime
	layout := "2006-01-02 15:04"

	// Parse full datetime string into time.Time
	checkOutDate, err := time.Parse(layout, fullDateTime)
	if err != nil {
		att.l.Error("attendance Error parsing date/time: ", attendance.EmployeeID, err)
	}

	attendance.CheckOutDate = checkOutDate.Format("2006-01-02 15:04:05")

	if attendanceId == 0 {
		inAttendance := dtos.CreateInAttendance{}
		inAttendance.CheckInDateStr = attendance.CheckOutDateStr
		inAttendance.CheckInDate = attendance.CheckOutDate
		inAttendance.InTime = attendance.OutTime
		inAttendance.CheckInByID = attendance.CheckOutByID
		inAttendance.CheckInLatitude = attendance.CheckOutLatitude
		inAttendance.CheckInLongitude = attendance.CheckOutLongitude
		attendance.Status = dtos.AttStatusPresent
		attendance.EmployeeID = employeeId
		attendanceId, err = att.attendanceDao.CreateInAttendance(inAttendance)
		if err != nil {
			att.l.Error("attendance not saved ", attendance.EmployeeID, err)
			return nil, err
		}
	}

	if attendanceId != 0 {
		err = att.attendanceDao.CreateOutAttendance(attendance, attendanceId)
		if err != nil {
			att.l.Error("attendance not saved ", attendance.EmployeeID, err)
			return nil, err
		}
	}

	att.l.Info("Attendance updated successfully! : ", attendanceId, attendance.EmployeeID)
	entry := dtos.AttendanceEntry{}
	entry.EntryInDate = attendance.CheckOutDate
	entry.EntryInDateStr = attendance.CheckOutDateStr
	entry.EntryTime = attendance.OutTime
	entry.EmployeeID = attendance.EmployeeID
	entry.AttendanceId = attendanceId
	err = att.attendanceDao.CreateAttendanceQuery(entry)
	if err != nil {
		att.l.Error("CreateAttendanceQuery not saved ", attendance.EmployeeID, err)
		//return nil, err
	}

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("Attendance saved successfully: time: %v", attendance.OutTime)
	return &response, nil
}

func (br *AttendanceObj) GetEmployeeAttendances(limit, offset, date string) (*dtos.EmployeeAttendances, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	res, errA := br.attendanceDao.GetEmployeeAttendances(limitI, offsetI, date)
	if errA != nil {
		br.l.Error("ERROR: GetBranchs", errA)
		return nil, errA
	}

	employeeAttendances := dtos.EmployeeAttendances{}
	employeeAttendances.EmployeeAttendance = res
	employeeAttendances.Total = br.attendanceDao.GetTotalCount()
	employeeAttendances.Limit = limitI
	employeeAttendances.OffSet = offsetI
	return &employeeAttendances, nil
}
