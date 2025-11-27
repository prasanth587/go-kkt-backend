package daos

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type LRReceipt struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewLRReceiptObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *LRReceipt {
	return &LRReceipt{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type LRReceiptDao interface {
	CreateLRReceipt(lrReq dtos.LRReceiptReq) error
	UpdateTripSheetHeader(tripSheetID int64, isLRGenerated int64) error
	BuildWhereQuery(orgId int64, lrId, tripSheetNum, lrNumber, tripDate, searchText string) string
	GetLRRecord(lrId int64) (*dtos.LRReceipt, error)
	UpdateLR(lrId int64, lrReq dtos.LRReceiptUpdateReq) error
	GetLRRecords(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.LRReceipt, []string, error)
	GetTotalCount(whereQuery string) int64

	GetLoadingUnloadingPoints(tripSheetIds string) (*[]dtos.PodLoadUnLoad, error)
	GetLocationNameById(loadingPointId int64) (*dtos.LoadUnLoadLoc, error)
}

func (mp *LRReceipt) BuildWhereQuery(orgId int64, lrId, tripSheetNum, lrNumber, tripDate, searchText string) string {

	whereQuery := fmt.Sprintf("WHERE org_id = '%v'", orgId)

	if tripSheetNum != "" {
		whereQuery = fmt.Sprintf(" %v AND trip_sheet_num = '%v'", whereQuery, tripSheetNum)
	}
	if lrNumber != "" {
		whereQuery = fmt.Sprintf(" %v AND lr_number = '%v'", whereQuery, lrNumber)
	}
	if tripDate != "" {
		whereQuery = fmt.Sprintf(" %v AND trip_date = '%v'", whereQuery, tripDate)
	}

	if searchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (trip_sheet_num LIKE '%%%v%%' OR lr_number LIKE '%%%v%%' OR vehicle_number LIKE '%%%v%%' OR invoice_number LIKE '%%%v%%' OR consignor_name LIKE '%%%v%%' OR consignor_address LIKE '%%%v%%' OR consignee_name LIKE '%%%v%%' OR consignee_address LIKE '%%%v%%' OR goods_type LIKE '%%%v%%' OR remark LIKE '%%%v%%' ) ", whereQuery, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText, searchText)
	}

	mp.l.Info("lr whereQuery:\n ", whereQuery)
	return whereQuery
}

func (rl *LRReceipt) GetTotalCount(whereQuery string) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM lr_receipt_header %v`, whereQuery)
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

func (mp *LRReceipt) CreateLRReceipt(lrReq dtos.LRReceiptReq) error {

	mp.l.Info("CreateLRReceipt : ", lrReq.LRNumber)

	insertQuery := fmt.Sprintf(`INSERT INTO lr_receipt_header (
    trip_sheet_id, trip_sheet_num, lr_number, trip_date, vehicle_number, vehicle_size,
    invoice_number, invoice_value, consignor_name, consignor_address, consignor_gst,
    consignee_name, consignee_address, consignee_gst, goods_type, goods_weight,
    quantity_in_pieces, remark, org_id, driver_name, driver_mobile_number
		) VALUES (
			'%v', '%v', '%v', '%v', '%v', '%v',
			'%v', '%v', '%v', '%v', '%v',
			'%v', '%v', '%v', '%v', '%v',
			'%v', '%v', '%v','%v', '%v'
		)`,
		lrReq.TripSheetID, lrReq.TripSheetNum, lrReq.LRNumber, lrReq.TripDate, lrReq.VehicleNumber, lrReq.VehicleSize,
		lrReq.InvoiceNumber, lrReq.InvoiceValue, lrReq.ConsignorName, lrReq.ConsignorAddress, lrReq.ConsignorGST,
		lrReq.ConsigneeName, lrReq.ConsigneeAddress, lrReq.ConsigneeGST, lrReq.GoodsType, lrReq.GoodsWeight,
		lrReq.QuantityInPieces, lrReq.Remark, lrReq.OrgID, lrReq.DriverName, lrReq.DriverMobileNumber,
	)

	mp.l.Info("CreateLRReceipt : ", insertQuery)

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
	mp.l.Info("Pod created successfully ID: ", lrReq.TripSheetNum, insertedID)
	return nil
}

func (lr *LRReceipt) UpdateTripSheetHeader(tripSheetID int64, isLRGenerated int64) error {

	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET is_lr_generated = '%v' WHERE trip_sheet_id = '%v'`, isLRGenerated, tripSheetID)

	lr.l.Info("LR UpdateTripSheetHeader query: ", updateQuery)

	_, err := lr.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		lr.l.Error("Error db.Exec(UpdateTripSheetHeader): ", err)
		return err
	}
	lr.l.Info("lr updated successfully ", tripSheetID)
	return nil
}

func (lr *LRReceipt) GetLRRecord(lrId int64) (*dtos.LRReceipt, error) {

	lrQuery := fmt.Sprintf(`SELECT lr_id, trip_sheet_id, trip_sheet_num, lr_number, trip_date, vehicle_number, 
	vehicle_size, invoice_number, invoice_value, consignor_name, consignor_address, consignor_gst, 
	consignee_name, consignee_address, consignee_gst, goods_type, goods_weight, 
	quantity_in_pieces, remark, org_id, driver_name, driver_mobile_number FROM lr_receipt_header WHERE lr_id = '%v'`, lrId)

	lr.l.Info("lrQuery:\n ", lrQuery)

	row := lr.dbConnMSSQL.GetQueryer().QueryRow(lrQuery)

	var (
		lrID, tripSheetID, orgID                                                                     sql.NullInt64
		tripSheetNum, lrNumber, tripDate, vehicleNumber, vehicleSize                                 sql.NullString
		invoiceNumber, consignorName, consignorAddress, consignorGST, driverName, driverMobileMumber sql.NullString
		consigneeName, consigneeAddress, consigneeGST, goodsType                                     sql.NullString
		goodsWeight, quantityInPieces, remark                                                        sql.NullString
		invoiceValue                                                                                 sql.NullFloat64
	)

	lrR := &dtos.LRReceipt{}
	err := row.Scan(
		&lrID, &tripSheetID, &tripSheetNum, &lrNumber, &tripDate, &vehicleNumber,
		&vehicleSize, &invoiceNumber, &invoiceValue, &consignorName, &consignorAddress,
		&consignorGST, &consigneeName, &consigneeAddress, &consigneeGST, &goodsType,
		&goodsWeight, &quantityInPieces, &remark, &orgID, &driverName, &driverMobileMumber,
	)
	if err != nil {
		lr.l.Error("Error LRReceipt scan: ", err)
		return nil, err
	}
	lrR.LRId = lrId
	lrR.TripSheetID = tripSheetID.Int64
	lrR.TripSheetNum = tripSheetNum.String
	lrR.LRNumber = lrNumber.String
	lrR.TripDate = tripDate.String
	lrR.VehicleNumber = vehicleNumber.String
	lrR.VehicleSize = vehicleSize.String
	lrR.InvoiceNumber = invoiceNumber.String
	lrR.InvoiceValue = invoiceValue.Float64
	lrR.ConsigneeName = consigneeName.String
	lrR.ConsigneeAddress = consigneeAddress.String
	lrR.ConsigneeGST = consigneeGST.String
	lrR.ConsignorName = consignorName.String
	lrR.ConsignorAddress = consignorAddress.String
	lrR.ConsignorGST = consignorGST.String
	lrR.GoodsType = goodsType.String
	lrR.GoodsWeight = goodsWeight.String
	lrR.QuantityInPieces = quantityInPieces.String
	lrR.Remark = remark.String
	lrR.OrgID = orgID.Int64
	lrR.DriverName = driverName.String
	lrR.DriverMobileNumber = driverMobileMumber.String

	return lrR, nil
}

func (mp *LRReceipt) UpdateLR(lrId int64, lrReq dtos.LRReceiptUpdateReq) error {

	updateLRQuery := fmt.Sprintf(`
    UPDATE lr_receipt_header SET
        trip_sheet_id = '%v',
        trip_sheet_num = '%v',
        lr_number = '%v',
        trip_date = '%v',
        vehicle_number = '%v',
        vehicle_size = '%v',
        invoice_number = '%v',
        invoice_value = '%v',
        consignor_name = '%v',
        consignor_address = '%v',
        consignor_gst = '%v',
        consignee_name = '%v',
        consignee_address = '%v',
        consignee_gst = '%v',
        goods_type = '%v',
        goods_weight = '%v',
        quantity_in_pieces = '%v',
		remark = '%v',
		driver_name = '%v',
		driver_mobile_number = '%v'
    WHERE lr_id = %v;`,
		lrReq.TripSheetID,
		lrReq.TripSheetNum,
		lrReq.LRNumber,
		lrReq.TripDate,
		lrReq.VehicleNumber,
		lrReq.VehicleSize,
		lrReq.InvoiceNumber,
		lrReq.InvoiceValue,
		lrReq.ConsignorName,
		lrReq.ConsignorAddress,
		lrReq.ConsignorGST,
		lrReq.ConsigneeName,
		lrReq.ConsigneeAddress,
		lrReq.ConsigneeGST,
		lrReq.GoodsType,
		lrReq.GoodsWeight,
		lrReq.QuantityInPieces,
		lrReq.Remark,
		lrReq.DriverName,
		lrReq.DriverMobileNumber,
		lrId,
	)

	mp.l.Info("updateLRQuery query: ", updateLRQuery)

	_, err := mp.dbConnMSSQL.GetQueryer().Exec(updateLRQuery)
	if err != nil {
		mp.l.Error("Error db.Exec(updateLRQuery): ", err)
		return err
	}
	mp.l.Info("LR updated successfully: ", lrId, lrReq.LRNumber)
	return nil
}

func (mp *LRReceipt) GetLRRecords(orgId int64, whereQuery string, limit int64, offset int64) (*[]dtos.LRReceipt, []string, error) {
	list := []dtos.LRReceipt{}
	var tripSheetIds []string
	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v;", whereQuery, limit, offset)

	mp.l.Info("lr whereQuery:\n ", whereQuery)

	lrQuery := fmt.Sprintf(`SELECT lr_id, trip_sheet_id, trip_sheet_num, lr_number, trip_date, vehicle_number, 
	vehicle_size, invoice_number, invoice_value, consignor_name, consignor_address, consignor_gst, 
	consignee_name, consignee_address, consignee_gst, goods_type, goods_weight, 
	quantity_in_pieces, remark, org_id, driver_name, driver_mobile_number FROM lr_receipt_header %v `, whereQuery)

	rows, err := mp.dbConnMSSQL.GetQueryer().Query(lrQuery)
	if err != nil {
		mp.l.Error("Error GetManagePods ", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var (
			lrID, tripSheetID, orgID                                                                     sql.NullInt64
			tripSheetNum, lrNumber, tripDate, vehicleNumber, vehicleSize                                 sql.NullString
			invoiceNumber, consignorName, consignorAddress, consignorGST, driverName, driverMobileMumber sql.NullString
			consigneeName, consigneeAddress, consigneeGST, goodsType                                     sql.NullString
			goodsWeight, quantityInPieces, remark                                                        sql.NullString
			invoiceValue                                                                                 sql.NullFloat64
		)
		lrR := &dtos.LRReceipt{}
		err := rows.Scan(
			&lrID, &tripSheetID, &tripSheetNum, &lrNumber, &tripDate, &vehicleNumber,
			&vehicleSize, &invoiceNumber, &invoiceValue, &consignorName, &consignorAddress,
			&consignorGST, &consigneeName, &consigneeAddress, &consigneeGST, &goodsType,
			&goodsWeight, &quantityInPieces, &remark, &orgID, &driverName, &driverMobileMumber,
		)
		if err != nil {
			mp.l.Error("Error GetLRRecords scan: ", err)
			return nil, nil, err
		}

		lrR.LRId = lrID.Int64
		lrR.TripSheetID = tripSheetID.Int64
		lrR.TripSheetNum = tripSheetNum.String
		lrR.LRNumber = lrNumber.String
		lrR.TripDate = tripDate.String
		lrR.VehicleNumber = vehicleNumber.String
		lrR.VehicleSize = vehicleSize.String
		lrR.InvoiceNumber = invoiceNumber.String
		lrR.InvoiceValue = invoiceValue.Float64
		lrR.ConsigneeName = consigneeName.String
		lrR.ConsigneeAddress = consigneeAddress.String
		lrR.ConsigneeGST = consigneeGST.String
		lrR.ConsignorName = consignorName.String
		lrR.ConsignorAddress = consignorAddress.String
		lrR.ConsignorGST = consignorGST.String
		lrR.GoodsType = goodsType.String
		lrR.GoodsWeight = goodsWeight.String
		lrR.QuantityInPieces = quantityInPieces.String
		lrR.Remark = remark.String
		lrR.OrgID = orgId
		lrR.DriverName = driverName.String
		lrR.DriverMobileNumber = driverMobileMumber.String

		if tripSheetID.Valid {
			tripSheetIdStr := strconv.FormatInt(tripSheetID.Int64, 10)
			tripSheetIds = append(tripSheetIds, tripSheetIdStr)
		}

		list = append(list, *lrR)
	}

	return &list, tripSheetIds, nil
}

func (mp *LRReceipt) GetLoadingUnloadingPoints(tripSheetIds string) (*[]dtos.PodLoadUnLoad, error) {
	list := []dtos.PodLoadUnLoad{}
	loadUnLoadQuery := fmt.Sprintf(`SELECT trip_sheet_id, type, loading_point_id FROM trip_sheet_header_load_unload_points WHERE trip_sheet_id IN (%v)  ORDER BY load_unload_point_id ASC; `, tripSheetIds)

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

func (mp *LRReceipt) GetLocationNameById(loadingPointId int64) (*dtos.LoadUnLoadLoc, error) {

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
