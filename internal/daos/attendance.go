package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type AttendanceObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewAttendanceObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *AttendanceObj {
	return &AttendanceObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

var (
// UPDATE_ERROR = "ERROR UpdateUserLoginAterCredentialSuccess update: %v "
// VERIFY_ERROR = "ERROR VerifyCredentials: %v "
// INVALIDE_    = "invalide credentials"
// LOGIN_FAILED = "login failed  - "
)

type AttendanceDao interface {
	CreateInAttendance(attendance dtos.CreateInAttendance) (int64, error)
	CreateOutAttendance(att dtos.CreateOutAttendance, attendanceId int64) error
	CheckAttendanceForEmployee(employeeId int64, checkInDateStr string) int64
	CreateAttendanceQuery(entry dtos.AttendanceEntry) error
	GetEmployeeAttendances(limit int64, offset int64, date string) (*[]dtos.EmployeeAttendance, error)
	GetTotalCount() int64
}

func (rl *AttendanceObj) GetTotalCount() int64 {

	empSelectQuery := "SELECT count(*) FROM employee"
	rl.l.Info(" GetTotalCount select query: ", empSelectQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(empSelectQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (br *AttendanceObj) CreateInAttendance(att dtos.CreateInAttendance) (int64, error) {

	br.l.Info("CreateAttendance : ", att.EmployeeID)

	createAttendanceQuery := fmt.Sprintf(`INSERT INTO attendance (
        employee_id, 
        check_in_date_str,
		check_in_date,
        in_time,
        check_in_latitude,
        check_in_longitude,
        check_in_by_id,
        status
        ) VALUES (
        %v, '%v', '%v', '%v', '%f', '%f', %v, '%v'
    );`,
		att.EmployeeID,
		att.CheckInDateStr,
		att.CheckInDate,
		att.InTime,
		att.CheckInLatitude,
		att.CheckInLongitude,
		att.CheckInByID,
		att.Status,
	)

	br.l.Info("createAttendanceQuery : ", createAttendanceQuery)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(createAttendanceQuery)
	if err != nil {
		br.l.Error("Error db.Exec(CreateAttendance): ", err)
		return 0, err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		br.l.Error("Error db.Exec(CreateAttendance):", createdId, err)
		return 0, err
	}
	br.l.Info("attendance created successfully", createdId, att.EmployeeID, att.InTime)
	return createdId, nil
}

func (br *AttendanceObj) CheckAttendanceForEmployee(employeeId int64, checkInDateStr string) int64 {

	attQuery := fmt.Sprintf(`SELECT attendance_id from attendance Where employee_id = '%v' AND check_in_date_str = '%v';`, employeeId, checkInDateStr)
	br.l.Info("attQuery:\n ", attQuery)
	row := br.dbConnMSSQL.GetQueryer().QueryRow(attQuery)
	var attendanceId sql.NullInt64
	err := row.Scan(&attendanceId)
	if err != nil {
		br.l.Error("Error CheckAttendanceForEmployee scan: ", err)
	}
	return attendanceId.Int64
}

// CreateAttendanceQuery(entry dtos.AttendanceEntry) error
func (br *AttendanceObj) CreateAttendanceQuery(entry dtos.AttendanceEntry) error {
	br.l.Info("CreateAttendanceEntry : ", entry.AttendanceId)

	// Prepare insert query
	createAttendanceLogQuery := fmt.Sprintf(`INSERT INTO attendance_entry_log (
        attendance_id,
        employee_id,
        entry_in_date,
        entry_in_date_str,
        entry_time
    ) VALUES (
        %v, %v, '%v', '%v', '%v'
    );`,
		entry.AttendanceId,
		entry.EmployeeID,
		entry.EntryInDate,
		entry.EntryInDateStr,
		entry.EntryTime,
	)

	br.l.Info("createAttendanceLogQuery : ", createAttendanceLogQuery)

	// Execute the query
	result, err := br.dbConnMSSQL.GetQueryer().Exec(createAttendanceLogQuery)
	if err != nil {
		br.l.Error("Error db.Exec(CreateAttendanceQuery): ", err)
		return err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		br.l.Error("Error fetching inserted ID in CreateAttendanceQuery:", err)
		return err
	}

	br.l.Info("attendance entry log created successfully", insertedID, entry.EmployeeID, entry.EntryTime)
	return nil
}

func (br *AttendanceObj) CreateOutAttendance(att dtos.CreateOutAttendance, attendanceId int64) error {

	br.l.Info("CreateAttendance : ", att.EmployeeID)

	updateAttendanceQuery := fmt.Sprintf(`UPDATE attendance SET
        check_out_date_str = '%v',
        check_out_date = '%v',
        out_time = '%v',
        check_out_latitude = %f,
        check_out_longitude = %f,
        check_out_by_id = %v,
        status = '%v'
    WHERE attendance_id = %v;`,
		att.CheckOutDateStr,
		att.CheckOutDate,
		att.OutTime,
		att.CheckOutLatitude,
		att.CheckOutLongitude,
		att.CheckOutByID,
		att.Status,
		attendanceId, // make sure your struct includes this field
	)

	br.l.Info("updateAttendanceQuery : ", updateAttendanceQuery)

	_, err := br.dbConnMSSQL.GetQueryer().Exec(updateAttendanceQuery)
	if err != nil {
		br.l.Error("Error db.Exec(CreateAttendance): ", err)
		return err
	}

	br.l.Info("attendance updated successfully", att.EmployeeID, att.OutTime)
	return nil
}

// SELECT e.emp_id,e.first_name,e.last_name, a.check_in_date_str, a.in_time, a.out_time, a.status FROM employee as e
// LEFT JOIN attendance as a ON a.employee_id = e.emp_id order by a.in_time DESC LIMIT 10 OFFSET 0

func (br *AttendanceObj) GetEmployeeAttendances(limit int64, offset int64, date string) (*[]dtos.EmployeeAttendance, error) {
	list := []dtos.EmployeeAttendance{}

	employeeAttendancesQuery := fmt.Sprintf(`SELECT e.emp_id,e.first_name,e.last_name, a.check_in_date_str, a.in_time, a.out_time, a.status, a.attendance_id, e.is_active FROM employee as e 
	LEFT JOIN attendance as a ON a.employee_id = e.emp_id order by a.in_time DESC LIMIT %v OFFSET %v;`, limit, offset)

	if date != "" {
		employeeAttendancesQuery = fmt.Sprintf(`SELECT e.emp_id,e.first_name,e.last_name, a.check_in_date_str, a.in_time, a.out_time, a.status, a.attendance_id, e.is_active FROM employee as e 
	LEFT JOIN attendance as a ON (a.employee_id = e.emp_id AND a.check_in_date_str = '%v' ) order by a.in_time DESC LIMIT %v OFFSET %v;`, date, limit, offset)
	}

	br.l.Info("employeeAttendancesQuery:\n ", employeeAttendancesQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(employeeAttendancesQuery)
	if err != nil {
		br.l.Error("Error GetEmployeeAttendances ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var firstName, lastname, checkInDateStr, inIime, outTime, status sql.NullString
		var empId, attendanceId, isActive sql.NullInt64

		empAtt := &dtos.EmployeeAttendance{}
		err := rows.Scan(&empId, &firstName, &lastname, &checkInDateStr, &inIime, &outTime, &status, &attendanceId, &isActive)
		if err != nil {
			br.l.Error("Error GetEmployeeAttendances scan: ", err)
			return nil, err
		}
		empAtt.EmployeeID = empId.Int64
		empAtt.EmployeeIDText = fmt.Sprintf("EMP_%v", empId.Int64)
		empAtt.AttendanceId = attendanceId.Int64
		empAtt.EmployeeName = fmt.Sprintf("%s %s", firstName.String, lastname.String)
		empAtt.CheckInDateStr = checkInDateStr.String
		empAtt.InTime = inIime.String
		empAtt.OutTime = outTime.String
		empAtt.Status = status.String
		empAtt.IsActive = isActive.Int64 == 1
		if checkInDateStr.String == "" || empAtt.Status == "" {
			empAtt.Status = "Absent"
		}

		list = append(list, *empAtt)
	}

	return &list, nil
}
