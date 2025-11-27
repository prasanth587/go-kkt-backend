package daos

import (
	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
)

type PingCheckStruct struct {
	Count int32 `json:"count"`
}

type Ping struct {
	l           *log.Logger
	dbConnMSSQL *mssqlcon.DBConn
}

func NewPing(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *Ping {
	return &Ping{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
	}
}

func (p *Ping) Ping() (bool, error) {
	pingModel := PingCheckStruct{}
	errMs := p.dbConnMSSQL.GetQueryer().QueryRow("SELECT 1").Scan(&pingModel.Count)
	if errMs != nil {
		p.l.Error("Ping error - MSSQL-", errMs)
		return false, errMs
	}
	return (pingModel.Count == 1), errMs
}
