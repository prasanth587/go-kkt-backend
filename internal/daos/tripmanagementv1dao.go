package daos

import "fmt"

func (br *TripSheetObj) DeleteTripSheetLoadUnLoadPoints(tripSheetID int64, lpType string) error {

	deleteQuery := fmt.Sprintf(`DELETE FROM trip_sheet_header_load_unload_points
WHERE trip_sheet_id = '%v' AND type = '%v';;
`, tripSheetID, lpType)

	br.l.Info("DeleteTripSheetLoadUnLoadPoints query:", deleteQuery)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(deleteQuery)
	if err != nil {
		br.l.Error("Error db.Exec(DeleteTripSheetLoadUnLoadPoints): ", err)
		return err
	}

	br.l.Info("deleted successfully: ", lpType, roleResult)

	return nil
}
