package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type LoadingPointObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewLoadingPointObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *LoadingPointObj {
	return &LoadingPointObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

type LoadingPointDao interface {
	CreateLoadingPoint(lpReq dtos.LoadingPointReq) error
	GetLoadingPoints(orgId int64, limit int64, offset int64, whereQuery string) (*[]dtos.LoadingPointEntries, error)
	GetLoadingPoint(loadingPointId int64) (*dtos.LoadingPointEntries, error)
	UpdateLoadingPoint(loadingPointId int64, lpu dtos.LoadingPointUpdate) error
	UpdateBranchActiveStatus(loadingpointId, isActive int64) error
	GetTotalCount(whereQuery string) int64

	BuildWhereQuery(orgId int64, tripSearchText, loadingPointId string) string
}

func (rl *LoadingPointObj) BuildWhereQuery(orgId int64, searchText, loadingPointId string) string {

	whereQuery := fmt.Sprintf("WHERE org_id = '%v'", orgId)

	if loadingPointId != "" {
		whereQuery = fmt.Sprintf(" %v AND loading_point_id = '%v'", whereQuery, loadingPointId)
	}

	if searchText != "" {
		whereQuery = fmt.Sprintf(" %v AND (city_code LIKE '%%%v%%' OR city_name LIKE '%%%v%%' OR address_line LIKE '%%%v%%' OR map_link LIKE '%%%v%%' OR state LIKE '%%%v%%' ) ", whereQuery, searchText, searchText, searchText, searchText, searchText)
	}

	rl.l.Info("loading_point_id whereQuery:\n ", whereQuery)

	return whereQuery
}

func (rl *LoadingPointObj) GetTotalCount(whereQuery string) int64 {
	countQuery := fmt.Sprintf(`SELECT count(*) FROM loading_point %v`, whereQuery)
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

func (br *LoadingPointObj) CreateLoadingPoint(lpReq dtos.LoadingPointReq) error {

	br.l.Info("CreateLoadingPoint : ", lpReq.CityCode)

	createLoadingPointQuery := fmt.Sprintf(`
		INSERT INTO loading_point (
			branch_id, org_id, is_active, 
			city_code, city_name, address_line, 
			map_link, state, country
		) VALUES (
			'%v', '%v', '%v', 
			'%v', '%v', '%v', 
			'%v', '%v', '%v'
		)`,
		lpReq.BranchID, lpReq.OrgID, lpReq.IsActive,
		lpReq.CityCode, lpReq.CityName, lpReq.AddressLine,
		lpReq.MapLink, lpReq.State, lpReq.Country,
	)

	br.l.Info("CreateLoadingPointQuery : ", createLoadingPointQuery)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(createLoadingPointQuery)
	if err != nil {
		br.l.Error("Error db.Exec(CreateLoadingPoint): ", err)
		return err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		br.l.Error("Error db.Exec(CreateLoadingPoint):", createdId, err)
		return err
	}
	br.l.Info("loading_point created successfully: ", createdId, lpReq.CityCode)
	return nil
}

func (br *LoadingPointObj) GetLoadingPoints(orgId int64, limit int64, offset int64, whereQuery string) (*[]dtos.LoadingPointEntries, error) {
	list := []dtos.LoadingPointEntries{}

	whereQuery = fmt.Sprintf(" %v ORDER BY updated_at DESC LIMIT %v OFFSET %v", whereQuery, limit, offset)

	loadingpointQuery := fmt.Sprintf(`SELECT loading_point_id, branch_id, is_active, 
	city_code, city_name, address_line,
	map_link, state, country FROM loading_point %v;`, whereQuery)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(loadingpointQuery)
	if err != nil {
		br.l.Error("Error LoadingPoints ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cityCode, cityName, addressLine, mapLink, state, country sql.NullString
		var loadingpointId, branchIdN, isActive sql.NullInt64

		loadingPointE := &dtos.LoadingPointEntries{}
		err := rows.Scan(&loadingpointId, &branchIdN, &isActive, &cityCode, &cityName, &addressLine, &mapLink, &state, &country)
		if err != nil {
			br.l.Error("Error GetLoadingpoints scan: ", err)
			return nil, err
		}
		loadingPointE.LoadingPointID = loadingpointId.Int64
		loadingPointE.BranchID = branchIdN.Int64
		loadingPointE.OrgID = orgId
		loadingPointE.IsActive = isActive.Int64
		loadingPointE.CityCode = cityCode.String
		loadingPointE.CityName = cityName.String
		loadingPointE.State = state.String
		loadingPointE.OrgID = orgId
		loadingPointE.Country = country.String
		loadingPointE.MapLink = mapLink.String
		loadingPointE.AddressLine = addressLine.String

		list = append(list, *loadingPointE)
	}

	return &list, nil
}

func (br *LoadingPointObj) UpdateLoadingPoint(loadingPointId int64, lpu dtos.LoadingPointUpdate) error {

	updateBranchQuery := fmt.Sprintf(`
    UPDATE loading_point SET
        branch_id = '%v', is_active = '%v', org_id = '%v',
        city_code = '%v', city_name = '%v', address_line = '%v',
        map_link = '%v', state = '%v', country = '%v'
    WHERE loading_point_id = %v;`,
		lpu.BranchID, lpu.IsActive, lpu.OrgID,
		lpu.CityCode, lpu.CityName, lpu.AddressLine,
		lpu.MapLink, lpu.State, lpu.Country,
		loadingPointId)

	br.l.Info("UpdateLoadingPoint: ", updateBranchQuery)

	_, err := br.dbConnMSSQL.GetQueryer().Exec(updateBranchQuery)
	if err != nil {
		br.l.Error("Error db.Exec(UpdateBranch): ", err)
		return err
	}
	br.l.Info("loading/unloading point updated successfully: ", loadingPointId)
	return nil
}

func (br *LoadingPointObj) GetLoadingPoint(loadingpointId int64) (*dtos.LoadingPointEntries, error) {

	loadingpointQuery := fmt.Sprintf(`SELECT loading_point_id, branch_id, is_active, city_code, city_name, address_line, map_link, state, country, org_id FROM loading_point WHERE loading_point_id = '%v';`, loadingpointId)

	br.l.Info("loadingpointQuery:\n ", loadingpointQuery)

	row := br.dbConnMSSQL.GetQueryer().QueryRow(loadingpointQuery)

	var cityCode, cityName, addressLine, mapLink, state, country sql.NullString
	var branchIdN, isActive, orgId sql.NullInt64

	loadingPointE := dtos.LoadingPointEntries{}
	err := row.Scan(&loadingpointId, &branchIdN, &isActive, &cityCode, &cityName, &addressLine, &mapLink, &state, &country, &orgId)
	if err != nil {
		br.l.Error("Error GetLoadingPoint scan: ", err)
		return nil, err
	}
	loadingPointE.LoadingPointID = loadingpointId
	loadingPointE.BranchID = branchIdN.Int64
	loadingPointE.OrgID = orgId.Int64
	loadingPointE.IsActive = isActive.Int64
	loadingPointE.CityCode = cityCode.String
	loadingPointE.CityName = cityName.String
	loadingPointE.State = state.String
	loadingPointE.Country = country.String
	loadingPointE.MapLink = mapLink.String
	loadingPointE.AddressLine = addressLine.String
	return &loadingPointE, nil
}

func (br *LoadingPointObj) UpdateBranchActiveStatus(loadingpointId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE loading_point SET is_active = '%v' WHERE loading_point_id = '%v'`, isActive, loadingpointId)

	br.l.Info("UpdateBranchActiveStatus Update query ", updateQuery)

	_, err := br.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		br.l.Error("Error db.Exec(UpdateBranchActiveStatus): ", err)
		return err
	}

	br.l.Info("loading/unloading point updated successfully ", loadingpointId)

	return nil
}
