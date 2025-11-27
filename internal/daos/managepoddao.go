package daos

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/constant"
	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type ManagePod struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewManagePodObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *ManagePod {
	return &ManagePod{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type ManagePodDao interface {
	CreateManagePod(podReq dtos.ManagePodReq) error
	BuildWhereQuery(orgId int64, podId, podStatus, tripType, searchText string) string
	GetManagePods(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.ManagePod, []string, error)
	GetManagePod(podId int64) (*dtos.ManagePod, error)
	UpdatePod(podId int64, podReq dtos.UpdateManagePodReq) error
	UpdatePODStatus(podId int64, status string) error
	UpdatePODStatusDelete(podId int64, lrNumber, status string) error
	GetLoadingUnloadingPoints(tripSheetIds string) (*[]dtos.PodLoadUnLoad, error)
	GetTotalCount(whereQuery string) int64
	GetLocationNameById(loadingPointId int64) (*dtos.LoadUnLoadLoc, error)
	UpdateTripSheetStatus(tripSheetId int64, status, statusDateUpdateQuery string) error
	UpdatePodImagePath(updateQuery string) error
}

func (mp *ManagePod) BuildWhereQuery(orgId int64, podId, podStatus, tripType, searchText string) string {

	whereQuery := fmt.Sprintf("WHERE org_id = '%v'", orgId)

	if podId != "" {
		whereQuery = fmt.Sprintf(" %v AND pod_id = '%v'", whereQuery, podId)
	}
	if tripType != "" {
		whereQuery = fmt.Sprintf(" %v AND LOWER(trip_type) = LOWER('%v')", whereQuery, tripType)
	}

	whereQuery = fmt.Sprintf(" %v AND pod_status != '%v'", whereQuery, constant.STATUS_DELETED)
	if podStatus != "" {
		whereQuery = fmt.Sprintf(" %v AND LOWER(pod_status) = LOWER('%v')", whereQuery, podStatus)
	}

	if searchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (trip_sheet_num LIKE '%%%v%%' OR lr_number LIKE '%%%v%%' OR customer_name LIKE '%%%v%%' OR pod_status LIKE '%%%v%%' OR pod_remark LIKE '%%%v%%' OR paid_by LIKE '%%%v%%' ) ", whereQuery, searchText, searchText, searchText, searchText, searchText, searchText)
	}

	mp.l.Info("pod whereQuery:\n ", whereQuery)

	return whereQuery
}

func (rl *ManagePod) GetTotalCount(whereQuery string) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM manage_pod %v`, whereQuery)
	rl.l.Info(" GetTotalCount select query: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (mp *ManagePod) CreateManagePod(podReq dtos.ManagePodReq) error {

	mp.l.Info("CreateManagePod : ", podReq.LRNumber)

	insertQuery := fmt.Sprintf(`INSERT INTO manage_pod (
        trip_sheet_id, trip_sheet_num, lr_number, customer_name, customer_id, 
		unloading_charges, unloading_date, paid_by, pod_submited_date, late_submission_debit,
        pod_remark, pod_status, org_id, pod_doc, halting_amount, 
		trip_type, kilometers_covered, halting_days
    ) VALUES (
        '%v', '%v', '%v', '%v','%v',
        '%v', '%v', '%v', '%v','%v',
        '%v', '%v', '%v','%v', '%v',
		'%v', '%v', '%v'
    )`,
		podReq.TripSheetID, podReq.TripSheetNum, podReq.LRNumber, podReq.CustomerName, podReq.CustomerID,
		podReq.UnPoadingCharges, podReq.UnPoadingDate, podReq.PaidBy, podReq.PODSubmitedDate, podReq.LateSubmissionDebit,
		podReq.PodRemark, podReq.PodStatus, podReq.OrgId, podReq.PodUpload, podReq.HaltingAmount, podReq.TripType,
		podReq.RunningKM, podReq.HaltingDays,
	)

	mp.l.Info("CreateManagePod : ", insertQuery)

	result, err := mp.dbConnMSSQL.GetQueryer().Exec(insertQuery)
	if err != nil {
		mp.l.Error("Error creating manage_pod: ", err)
		return err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		mp.l.Error("Error getting last insert ID: ", err)
		return err
	}
	mp.l.Info("Pod created successfully ID: ", podReq.TripSheetNum, insertedID)
	return nil
}

func (mp *ManagePod) GetManagePods(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.ManagePod, []string, error) {
	list := []dtos.ManagePod{}
	var tripSheetIds []string

	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)

	mp.l.Info("pods whereQuery:\n ", whereQuery)

	podQuery := fmt.Sprintf(`SELECT pod_id, trip_sheet_id, trip_sheet_num, lr_number, customer_name,customer_id, send_by, pod_status,
        pod_remark, late_submission_debit, pod_submited_date,halting_amount,paid_by,unloading_date, unloading_charges, 
		pod_doc, trip_type,kilometers_covered, halting_days FROM manage_pod %v `, whereQuery)

	rows, err := mp.dbConnMSSQL.GetQueryer().Query(podQuery)
	if err != nil {
		mp.l.Error("Error GetManagePods ", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tripSheetNum, lrNumber, submitedDate, customerName, sendBy,
			podRemark, podStatus,
			podSubmitedDate, paidBy, unloadingDate, podDoc, tripType sql.NullString
		var podId, tripSheetId, customerId sql.NullInt64
		var haltingAmount, unloadingCharges, kilometersCovered, haltingDays, lateSubmissionDebit sql.NullFloat64

		pod := &dtos.ManagePod{}
		err := rows.Scan(&podId, &tripSheetId, &tripSheetNum, &lrNumber, &customerName, &customerId, &sendBy,
			&podStatus, &podRemark, &lateSubmissionDebit, &podSubmitedDate, &haltingAmount, &paidBy, &unloadingDate,
			&unloadingCharges, &podDoc, &tripType, &kilometersCovered, &haltingDays)
		if err != nil {
			mp.l.Error("Error GetManagePods scan: ", err)
			return nil, nil, err
		}
		pod.PodId = podId.Int64
		pod.TripSheetID = tripSheetId.Int64
		pod.TripSheetNum = tripSheetNum.String
		pod.LRNumber = lrNumber.String
		pod.SendBy = sendBy.String
		pod.SubmitedDate = submitedDate.String
		pod.CustomerName = customerName.String
		pod.CustomerID = customerId.Int64
		pod.PodStatus = podStatus.String
		pod.PodRemark = podRemark.String
		pod.LateSubmissionDebit = lateSubmissionDebit.Float64
		pod.PODSubmitedDate = podSubmitedDate.String
		pod.HaltingAmount = haltingAmount.Float64
		pod.PaidBy = paidBy.String
		pod.UnPoadingDate = unloadingDate.String
		pod.UnPoadingCharges = unloadingCharges.Float64
		pod.PodUpload = podDoc.String
		pod.TripType = tripType.String
		pod.KilometersCovered = kilometersCovered.Float64
		pod.HaltingDays = haltingDays.Float64
		list = append(list, *pod)

		var tripSheetIdStr string
		if tripSheetId.Valid {
			tripSheetIdStr = strconv.FormatInt(tripSheetId.Int64, 10)
			tripSheetIds = append(tripSheetIds, tripSheetIdStr)
		}

	}

	return &list, tripSheetIds, nil
}

func (mp *ManagePod) GetLoadingUnloadingPoints(tripSheetIds string) (*[]dtos.PodLoadUnLoad, error) {
	list := []dtos.PodLoadUnLoad{}
	loadUnLoadQuery := fmt.Sprintf(`SELECT trip_sheet_id, type, loading_point_id FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id IN (%v)  ORDER BY trip_sheet_id `, tripSheetIds)

	mp.l.Info("GetLoadingUnloadingPoints whereQuery:\n ", loadUnLoadQuery)

	rows, err := mp.dbConnMSSQL.GetQueryer().Query(loadUnLoadQuery)
	if err != nil {
		mp.l.Error("Error GetLoadingUnloadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		loadUnLoad := &dtos.PodLoadUnLoad{}

		var tripSheetId, loadUnloadPointId sql.NullInt64
		var typeT sql.NullString
		err := rows.Scan(&tripSheetId, &typeT, &loadUnloadPointId)
		if err != nil {
			mp.l.Error("Error GetManagePods scan: ", err)
			return nil, err
		}

		loadUnLoad.TripSheetId = tripSheetId.Int64
		loadUnLoad.Type = typeT.String
		loadUnLoad.LoadingPointId = loadUnloadPointId.Int64

		list = append(list, *loadUnLoad)
	}
	return &list, nil
}

func (mp *ManagePod) UpdatePod(podId int64, podReq dtos.UpdateManagePodReq) error {

	updatePodQuery := fmt.Sprintf(`
    UPDATE manage_pod SET
        trip_sheet_id = '%v',
        trip_sheet_num = '%v',
        lr_number = '%v',
		customer_id = '%v',
		customer_name = '%v',
		pod_doc = '%v',
		unloading_charges = '%v',
		unloading_date = '%v',
		paid_by = '%v',
		pod_submited_date = '%v',
		late_submission_debit = '%v',
		pod_remark = '%v',
		trip_type = '%v',
		halting_amount = '%v',
		pod_status = '%v',
		kilometers_covered = '%v',
		halting_days = '%v'
    WHERE pod_id = %v;`,
		podReq.TripSheetID,
		podReq.TripSheetNum,
		podReq.LRNumber,
		podReq.CustomerID,
		podReq.CustomerName,
		podReq.PodUpload,
		podReq.UnPoadingCharges,
		podReq.UnPoadingDate,
		podReq.PaidBy, podReq.PODSubmitedDate,
		podReq.LateSubmissionDebit,
		podReq.PodRemark,
		podReq.TripType,
		podReq.HaltingAmount,
		podReq.PodStatus, podReq.RunningKM, podReq.HaltingDays,
		podId)

	mp.l.Info("updatePodQuery query: ", updatePodQuery)

	_, err := mp.dbConnMSSQL.GetQueryer().Exec(updatePodQuery)
	if err != nil {
		mp.l.Error("Error db.Exec(updatePodQuery): ", err)
		return err
	}
	mp.l.Info("pod updated successfully: ", podId, podReq.LRNumber)
	return nil
}

func (mp *ManagePod) GetManagePod(podId int64) (*dtos.ManagePod, error) {

	podQuery := fmt.Sprintf(`SELECT trip_sheet_id, trip_sheet_num, lr_number, customer_name,customer_id, send_by, pod_status,
        pod_remark, late_submission_debit, pod_submited_date,halting_amount,paid_by,unloading_date, unloading_charges, 
		pod_doc, trip_type,kilometers_covered, halting_days FROM manage_pod where pod_id = '%v'`, podId)

	mp.l.Info("podQuery:\n ", podQuery)

	row := mp.dbConnMSSQL.GetQueryer().QueryRow(podQuery)

	var tripSheetNum, lrNumber, submitedDate, customerName, sendBy,
		podRemark, podStatus,
		podSubmitedDate, paidBy, unloadingDate, podDoc, tripType sql.NullString
	var tripSheetId, customerId sql.NullInt64
	var haltingAmount, unloadingCharges, kilometersCovered, haltingDays, lateSubmissionDebit sql.NullFloat64

	pod := &dtos.ManagePod{}
	err := row.Scan(&tripSheetId, &tripSheetNum, &lrNumber, &customerName, &customerId, &sendBy,
		&podStatus, &podRemark, &lateSubmissionDebit, &podSubmitedDate, &haltingAmount, &paidBy, &unloadingDate,
		&unloadingCharges, &podDoc, &tripType, &kilometersCovered, &haltingDays)
	if err != nil {
		mp.l.Error("Error GetManagePods scan: ", err)
		return nil, err
	}
	pod.PodId = podId
	pod.TripSheetID = tripSheetId.Int64
	pod.TripSheetNum = tripSheetNum.String
	pod.LRNumber = lrNumber.String
	pod.SendBy = sendBy.String
	pod.SubmitedDate = submitedDate.String
	pod.CustomerName = customerName.String
	pod.CustomerID = customerId.Int64
	pod.PodStatus = podStatus.String
	pod.PodRemark = podRemark.String
	pod.LateSubmissionDebit = lateSubmissionDebit.Float64
	pod.PODSubmitedDate = podSubmitedDate.String
	pod.HaltingAmount = haltingAmount.Float64
	pod.PaidBy = paidBy.String
	pod.UnPoadingDate = unloadingDate.String
	pod.UnPoadingCharges = unloadingCharges.Float64
	pod.PodUpload = podDoc.String
	pod.TripType = tripType.String
	pod.KilometersCovered = kilometersCovered.Float64
	pod.HaltingDays = haltingDays.Float64
	return pod, nil
}

func (mp *ManagePod) UpdatePODStatus(podId int64, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE manage_pod SET pod_status = '%v' WHERE pod_id = '%v'`, status, podId)

	mp.l.Info("UpdatePODStatus Update query ", updateQuery)

	_, err := mp.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		mp.l.Error("Error db.Exec(UpdatePODStatus): ", err)
		return err
	}

	mp.l.Info("pod status updated successfully: ", podId)

	return nil
}

func (mp *ManagePod) UpdatePODStatusDelete(podId int64, lrNumber, status string) error {

	updateQuery := fmt.Sprintf(`UPDATE manage_pod SET 
	pod_status = '%v', 
	lr_number = '%v'
	WHERE pod_id = '%v'`, status, lrNumber, podId)

	mp.l.Info("UpdatePODStatusDelete Update query ", updateQuery)

	_, err := mp.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		mp.l.Error("Error db.Exec(UpdatePODStatusDelete): ", err)
		return err
	}

	mp.l.Info("pod status delete updated successfully: ", podId)

	return nil
}

func (mp *ManagePod) UpdateTripSheetStatus(tripSheetId int64, status, statusDateUpdateQuery string) error {

	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET load_status = '%v' %v WHERE trip_sheet_id = '%v'`, status, statusDateUpdateQuery, tripSheetId)

	mp.l.Info("UpdateTripSheetStatus Update query ", updateQuery)

	_, err := mp.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		mp.l.Error("Error db.Exec(UpdateTripSheetStatus): ", err)
		return err
	}
	mp.l.Info("trip_sheet_header status updated successfully: ", tripSheetId, status)

	return nil
}

func (mp *ManagePod) GetLocationNameById(loadingPointId int64) (*dtos.LoadUnLoadLoc, error) {

	locQuery := fmt.Sprintf(`SELECT city_code, city_name FROM loading_point WHERE loading_point_id = '%v' `, loadingPointId)
	mp.l.Info("GetLocationNameById whereQuery:\n ", locQuery)

	row := mp.dbConnMSSQL.GetQueryer().QueryRow(locQuery)
	var cityCode, cityName sql.NullString
	errE := row.Scan(&cityCode, &cityName)
	if errE != nil {
		mp.l.Error("Error GetCount scan: ", errE)
		return nil, errE
	}
	loadUnLoad := &dtos.LoadUnLoadLoc{}
	loadUnLoad.CityCode = cityCode.String
	loadUnLoad.CityName = cityName.String
	loadUnLoad.LoadingPointId = loadingPointId

	return loadUnLoad, nil
}

func (rl *ManagePod) UpdatePodImagePath(updateQuery string) error {

	rl.l.Info("UpdatePodImagePath Update query: ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdatePodImagePath) pod: ", err)
		return err
	}

	return nil
}
