package daos

import (
	"database/sql"
	"fmt"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
)

type BranchObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewBranchObj(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *BranchObj {
	return &BranchObj{
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

type BranchDao interface {
	CreateBranch(branchReq dtos.BranchReq) error
	GetBranches(orgId int64, limit int64, offset int64) (*[]dtos.Branch, error)
	UpdateBranch(branchId int64, veh dtos.BranchUpdate) error
	GetBranch(branchId int64) (*dtos.Branch, error)
	UpdateBranchActiveStatus(branchId, isActive int64) error
	// UpdateBranchImagePath(updateQuery string) error
	GetTotalCount(orgId int64) int64
}

func (br *BranchObj) GetTotalCount(orgId int64) int64 {

	countQuery := fmt.Sprintf(`SELECT count(*) FROM branch WHERE org_id = '%v'`, orgId)
	br.l.Info(" GetTotalCount select query: ", countQuery)
	row := br.dbConnMSSQL.GetQueryer().QueryRow(countQuery)
	var count sql.NullInt64

	errE := row.Scan(&count)
	if errE != nil {
		br.l.Error("Error GetCount scan: ", errE)
		return 0
	}

	return count.Int64
}

func (br *BranchObj) CreateBranch(branch dtos.BranchReq) error {

	br.l.Info("CreateBranch : ", branch.BranchName)

	createBranchQuery := fmt.Sprintf(`INSERT INTO branch (
		branch_name, branch_code, org_id, 
		is_active, address_line1, address_line2, 
		city, state, country)
		VALUES
		('%v', '%v', %v, 
		%v, '%v', '%v', 
		'%v', '%v', '%v')`,
		branch.BranchName, branch.BranchCode, branch.OrgID,
		branch.IsActive, branch.AddressLine1, branch.AddressLine2,
		branch.City, branch.State, branch.Country)

	br.l.Info("createBranchQuery : ", createBranchQuery)

	roleResult, err := br.dbConnMSSQL.GetQueryer().Exec(createBranchQuery)
	if err != nil {
		br.l.Error("Error db.Exec(CreateBranch): ", err)
		return err
	}
	createdId, err := roleResult.LastInsertId()
	if err != nil {
		br.l.Error("Error db.Exec(CreateBranch):", createdId, err)
		return err
	}
	br.l.Info("branch created successfully: ", createdId, branch.BranchCode)
	return nil
}

func (br *BranchObj) GetBranches(orgId int64, limit int64, offset int64) (*[]dtos.Branch, error) {
	list := []dtos.Branch{}

	branchQuery := fmt.Sprintf(`SELECT branch_id, branch_name, branch_code, address_line1, address_line2, is_active, city, state, country FROM branch WHERE org_id = '%v' ORDER BY updated_at DESC LIMIT %v OFFSET %v;`, orgId, limit, offset)

	br.l.Info("branchQuery:\n ", branchQuery)

	rows, err := br.dbConnMSSQL.GetQueryer().Query(branchQuery)
	if err != nil {
		br.l.Error("Error Branches ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branchName, branchCode, addressLine1, addressLine2, city, state, country sql.NullString
		var branchId, isActive sql.NullInt64

		branch := &dtos.Branch{}
		err := rows.Scan(&branchId, &branchName, &branchCode, &addressLine1, &addressLine2, &isActive, &city, &state,
			&country)
		if err != nil {
			br.l.Error("Error GetBranchs scan: ", err)
			return nil, err
		}
		branch.BranchId = branchId.Int64
		branch.BranchName = branchName.String
		branch.BranchCode = branchCode.String
		branch.AddressLine1 = addressLine1.String
		branch.AddressLine2 = addressLine2.String
		branch.IsActive = isActive.Int64
		branch.State = state.String
		branch.City = city.String
		branch.OrgID = orgId
		branch.Country = country.String
		list = append(list, *branch)
	}

	return &list, nil
}

func (br *BranchObj) UpdateBranch(branchId int64, branch dtos.BranchUpdate) error {

	updateBranchQuery := fmt.Sprintf(`
    UPDATE branch SET
        branch_name = '%v',
        branch_code = '%v',
        org_id = %v,
        is_active = %v,
        address_line1 = '%v',
        address_line2 = '%v',
        city = '%v',
        state = '%v',
        country = '%v'
    WHERE branch_id = %v;`,
		branch.BranchName,
		branch.BranchCode,
		branch.OrgID,
		branch.IsActive,
		branch.AddressLine1,
		branch.AddressLine2,
		branch.City,
		branch.State,
		branch.Country,
		branchId)

	br.l.Info("UpdateBranch query: ", updateBranchQuery)

	_, err := br.dbConnMSSQL.GetQueryer().Exec(updateBranchQuery)
	if err != nil {
		br.l.Error("Error db.Exec(UpdateBranch): ", err)
		return err
	}
	br.l.Info("branch updated successfully: ", branchId)
	return nil
}

func (br *BranchObj) GetBranch(branchId int64) (*dtos.Branch, error) {

	branchQuery := fmt.Sprintf(`SELECT branch_id, branch_name, branch_code, address_line1, address_line2, is_active, city, state, country, org_id FROM branch WHERE branch_id = '%v';`, branchId)

	br.l.Info("branchQuery:\n ", branchQuery)

	row := br.dbConnMSSQL.GetQueryer().QueryRow(branchQuery)

	var branchName, branchCode, addressLine1, addressLine2, city, state, country sql.NullString
	var branchID, isActive, orgId sql.NullInt64

	branch := dtos.Branch{}
	err := row.Scan(&branchID, &branchName, &branchCode, &addressLine1, &addressLine2, &isActive, &city, &state,
		&country, &orgId)
	if err != nil {
		br.l.Error("Error GetBranchs scan: ", err)
		return nil, err
	}
	branch.BranchId = branchId
	branch.BranchName = branchName.String
	branch.BranchCode = branchCode.String
	branch.AddressLine1 = addressLine1.String
	branch.AddressLine2 = addressLine2.String
	branch.IsActive = isActive.Int64
	branch.State = state.String
	branch.City = city.String
	branch.Country = country.String
	branch.OrgID = orgId.Int64

	return &branch, nil
}

func (br *BranchObj) UpdateBranchActiveStatus(branchId, isActive int64) error {

	updateQuery := fmt.Sprintf(`UPDATE branch SET is_active = '%v' WHERE branch_id = '%v'`, isActive, branchId)

	br.l.Info("UpdateBranchActiveStatus Update query ", updateQuery)

	_, err := br.dbConnMSSQL.GetQueryer().Exec(updateQuery)
	if err != nil {
		br.l.Error("Error db.Exec(UpdateBranchActiveStatus): ", err)
		return err
	}

	br.l.Info("branch status updated successfully: ", branchId)

	return nil
}
