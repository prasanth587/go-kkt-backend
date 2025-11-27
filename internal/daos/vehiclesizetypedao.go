package daos

import (
	"database/sql"
	"fmt"

	"go-transport-hub/dtos"
)

func (rl *VehicleObj) VehicleSizeTypesGetTotalCount() int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM vehicle_size_type WHERE is_active = '%v'`, 1)
	rl.l.Info(" VehicleSizeTypesGetTotalCount select query: ", countQuery)
	row := rl.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		rl.l.Error("Error VehicleSizeTypesGetTotalCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (rl *VehicleObj) CreateVehicleSizeTypes(vehicleReg dtos.VehicleSizeType) error {

	rl.l.Info("CreateDriver : ", vehicleReg.VehicleSize, vehicleReg.VehicleType)

	createVehicleSizeTypesQuery := fmt.Sprintf(`INSERT INTO vehicle_size_type (
		vehicle_size, vehicle_type, status, is_active)
		VALUES
		( '%v', '%v', '%v', '%v')`,
		vehicleReg.VehicleSize, vehicleReg.VehicleType, vehicleReg.Status, vehicleReg.IsActive)

	rl.l.Info("createVehicleSizeTypesQuery : ", createVehicleSizeTypesQuery)

	roleResult, err := rl.dbConnMSSQL.GetQueryer().Exec(createVehicleSizeTypesQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(createVehicleSizeTypesQuery): ", err)
		return err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		rl.l.Error("Error db.Exec(createVehicleSizeTypesQuery):", createdId, err)
		return err
	}
	rl.l.Info("Vehicle size type created successfully: ", createdId, vehicleReg.VehicleSize, vehicleReg.VehicleType)
	return nil
}

func (rl *VehicleObj) GetVehicleSizeTypes(limit int64, offset int64) (*[]dtos.VehicleSizeTypeObj, error) {
	list := []dtos.VehicleSizeTypeObj{}

	vehicleSizeTypesQuery := "SELECT vehicle_size_id, vehicle_size, vehicle_type, status, is_active  FROM vehicle_size_type ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	rl.l.Info("vehicleSizeTypesQuery: ", vehicleSizeTypesQuery)

	rows, err := rl.dbConnMSSQL.GetQueryer().Query(vehicleSizeTypesQuery, limit, offset)
	if err != nil {
		rl.l.Error("Error Vehicles ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicleType, vehicleSize, status sql.NullString

		var vehicleSizeId, isActive sql.NullInt64

		vehicleRes := &dtos.VehicleSizeTypeObj{}
		err := rows.Scan(&vehicleSizeId, &vehicleSize, &vehicleType, &status, &isActive)
		if err != nil {
			rl.l.Error("Error GetVehicles scan: ", err)
			return nil, err
		}
		vehicleRes.VehicleSizeId = vehicleSizeId.Int64
		vehicleRes.VehicleType = vehicleType.String
		vehicleRes.VehicleSize = vehicleSize.String
		vehicleRes.Status = status.String
		vehicleRes.IsActive = isActive.Int64

		list = append(list, *vehicleRes)
	}

	return &list, nil
}

//

func (rl *VehicleObj) GetVehicleSizeType(vehicleSizeID int64) (*dtos.VehicleSizeTypeObj, error) {

	vehicleQuery := "SELECT vehicle_size_id, vehicle_size, vehicle_type, status, is_active  FROM vehicle_size_type WHERE vehicle_size_id = ? "
	rl.l.Info("GetVehicleSizeType: ", vehicleQuery)

	row := rl.dbConnMSSQL.GetQueryer().QueryRow(vehicleQuery, vehicleSizeID)

	var vehicleType, vehicleSize, status sql.NullString

	var vehicleSizeId, isActive sql.NullInt64

	vehicleRes := &dtos.VehicleSizeTypeObj{}
	err := row.Scan(&vehicleSizeId, &vehicleSize, &vehicleType, &status, &isActive)
	if err != nil {
		rl.l.Error("Error GetVehicles scan: ", err)
		return nil, err
	}

	vehicleRes.VehicleSizeId = vehicleSizeId.Int64
	vehicleRes.VehicleType = vehicleType.String
	vehicleRes.VehicleSize = vehicleSize.String
	vehicleRes.Status = status.String
	vehicleRes.IsActive = isActive.Int64

	return vehicleRes, nil
}

func (rl *VehicleObj) UpdateVehicleSizeType(vehicleSizeId int64, veh dtos.VehicleSizeTypeUpdate) error {

	updateVehicleSizeTypeQuery := fmt.Sprintf(`
    UPDATE vehicle_size_type SET 
        vehicle_size = '%v',
        vehicle_type = '%v',
        status = '%v',
        is_active = '%v'
    WHERE vehicle_size_id = '%v'`,
		veh.VehicleSize, veh.VehicleType, veh.Status, veh.IsActive, vehicleSizeId)

	rl.l.Info("updateVehicleSizeType query: ", updateVehicleSizeTypeQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateVehicleSizeTypeQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(updateVehicleSizeTypeQuery): ", err)
		return err
	}
	rl.l.Info("Update Vehicle size type updated successfully: ", vehicleSizeId, veh.VehicleSize)
	return nil
}

func (rl *VehicleObj) UpdateVehicleSizeTypeActiveStatus(vehicleSizeId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE vehicle_size_type SET is_active = '%v' WHERE vehicle_size_id = '%v'`, isActive, vehicleSizeId)

	rl.l.Info("UpdateVehicleSizeTypeActiveStatus Update query ", updateQuery)

	_, err := rl.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		rl.l.Error("Error db.Exec(UpdateVehicleSizeTypeActiveStatus): ", err)
		return err
	}

	rl.l.Info("Update Vehicle size type updated successfully: ", vehicleSizeId)

	return nil
}
