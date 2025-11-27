package branch

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/dtos"
	"go-transport-hub/internal/daos"
)

type BranchObj struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
	branchDao   daos.BranchDao
}

var (
	ErrUnableToPingDB = errors.New("Unable to ping database")
	USER_SUCCESS      = "User logged in successfully"
	ERROR_IN_UPDATE   = "Error; UpdateUserLoginAterCredentialSuccess  - "
	INVALIDE_         = "invalid credentials"
	LOGIN_FAILED      = "login failed  - "
)

const (
	EmployeeMobilePattern = `^((\+)?(\d{2}[-])?(\d{10}){1})?(\d{11}){0,1}?$`
)

// DatePattern
const (
	DatePattern = `^\d{1,2}\/\d{1,2}\/\d{4}$`
)

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *BranchObj {
	return &BranchObj{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		branchDao:   daos.NewBranchObj(l, dbConnMSSQL),
	}
}

func (br *BranchObj) CreateBranch(branchReq dtos.BranchReq) (*dtos.Messge, error) {

	if branchReq.BranchName == "" {
		br.l.Error("Error branch name should not empty")
		return nil, errors.New("branch name should not empty")
	}
	if branchReq.BranchCode == "" {
		br.l.Error("Error BranchCode: branch code should not empty")
		return nil, errors.New("branch code should not empty")
	}
	branchReq.BranchCode = strings.ToUpper(branchReq.BranchCode)

	err1 := br.branchDao.CreateBranch(branchReq)
	if err1 != nil {
		br.l.Error("branch not saved ", branchReq.BranchName, err1)
		return nil, err1
	}

	br.l.Info("Branch created successfully! : ", branchReq.BranchName)

	response := dtos.Messge{}
	response.Message = fmt.Sprintf("Branch saved successfully: %s", branchReq.BranchName)
	return &response, nil
}

func (br *BranchObj) GetBranch(orgId int64, limit string, offset string) (*dtos.Branches, error) {

	limitI, errInt := strconv.ParseInt(limit, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid limit")
	}
	offsetI, errInt := strconv.ParseInt(offset, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid offset")
	}

	res, errA := br.branchDao.GetBranches(orgId, limitI, offsetI)
	if errA != nil {
		br.l.Error("ERROR: GetBranchs", errA)
		return nil, errA
	}

	branchEntries := dtos.Branches{}
	branchEntries.Branch = res
	branchEntries.Total = br.branchDao.GetTotalCount(orgId)
	branchEntries.Limit = limitI
	branchEntries.OffSet = offsetI
	return &branchEntries, nil
}

func (br *BranchObj) UpdateBranch(branchId int64, branchReq dtos.BranchUpdate) (*dtos.Messge, error) {

	if branchReq.BranchName == "" {
		br.l.Error("Error branch name should not empty")
		return nil, errors.New("branch name should not empty")
	}
	if branchReq.BranchCode == "" {
		br.l.Error("Error BranchCode: branch code should not empty")
		return nil, errors.New("branch code should not empty")
	}
	branchReq.BranchCode = strings.ToUpper(branchReq.BranchCode)

	branchInfo, errV := br.branchDao.GetBranch(branchId)
	if errV != nil {
		br.l.Error("ERROR: branch not found", branchId, errV)
		return nil, errV
	}
	jsonBytes, _ := json.Marshal(branchInfo)
	br.l.Info("GetBranch: ******* ", string(jsonBytes))

	err1 := br.branchDao.UpdateBranch(branchId, branchReq)
	if err1 != nil {
		br.l.Error("branch not updated ", branchReq.BranchName, err1)
		return nil, err1
	}

	br.l.Info("branch updated successfully! : ", branchReq.BranchName)

	roleResponse := dtos.Messge{}
	roleResponse.Message = fmt.Sprintf("branch updated successfully!: %s", branchReq.BranchName)
	return &roleResponse, nil
}

func (ul *BranchObj) UpdatebranchActiveStatus(isA string, branchId int64) (*dtos.Messge, error) {

	isActive, errInt := strconv.ParseInt(isA, 10, 64)
	if errInt != nil {
		return nil, errors.New("invalid status update")
	}
	if branchId == 0 {
		ul.l.Error("unknown branch", branchId)
		return nil, errors.New("unknown branch")
	}
	statusRes := dtos.Messge{}
	branchInfo, errA := ul.branchDao.GetBranch(branchId)
	if errA != nil {
		ul.l.Error("ERROR: branch not found", errA)
		return nil, errA
	}
	jsonBytes, _ := json.Marshal(branchInfo)
	ul.l.Info("branchInfo: ", string(jsonBytes))

	errU := ul.branchDao.UpdateBranchActiveStatus(branchId, isActive)
	if errU != nil {
		ul.l.Error("ERROR: UpdatebranchActiveStatus ", errU)
		return nil, errU
	}
	statusRes.Message = "branch updated successfully : " + branchInfo.BranchName
	return &statusRes, nil
}
