package ping

import (
	"errors"

	"github.com/prabha303-vi/log-util/log"

	"go-transport-hub/dbconn/mssqlcon"
	"go-transport-hub/internal/daos"
)

type Ping struct {
	l           *log.Logger
	ping        *daos.Ping
	dbConnMSSQL *mssqlcon.DBConn
}

var (
	ErrUnableToPingDB = errors.New("unable to ping database")
)

func New(l *log.Logger, dbConnMSSQL *mssqlcon.DBConn) *Ping {
	return &Ping{
		l:           l,
		dbConnMSSQL: dbConnMSSQL,
		ping:        daos.NewPing(l, dbConnMSSQL),
	}
}

func (p *Ping) Ping() error {
	ok, err := p.ping.Ping()
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnableToPingDB
	}
	return nil
}
