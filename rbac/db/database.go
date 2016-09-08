package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type (
	MgoConf struct {
		Url string
	}

	DataBase struct {
		session *mgo.Session
		dbName  string
	}
)

func Init(conf *MgoConf) (base *DataBase, err error) {
	mgoAddr, dbName, err := parseMgoAddr(conf.Url)
	if err != nil {
		return
	}
	return newDatabaseWithTimeout(mgoAddr, dbName, false)
}

func parseMgoAddr(url string) (mgoAddr string, dbName string, err error) {
	DbPos := strings.LastIndex(url, "/")
	if DbPos == -1 {
		err = errors.New("mgoDns don't contain '/'")
		return
	}
	mgoAddr = url[:DbPos]
	dbName = url[DbPos+1:]
	return
}

func newDatabaseWithTimeout(mgoAddr, dbName string, useTimeout bool, timeouts ...time.Duration) (res *DataBase, err error) {
	var mgoSession *mgo.Session
	var timeout time.Duration
	if useTimeout {
		if len(timeouts) > 0 {
			timeout = timeouts[0]
		}
		mgoSession, err = mgo.DialWithTimeout(mgoAddr, timeout)
	} else {
		mgoSession, err = mgo.Dial(mgoAddr)
	}
	if err != nil {
		err = fmt.Errorf("mgo.Dial error: %s", err)
		return
	}

	res = &DataBase{
		session: mgoSession,
		dbName:  dbName,
	}
	res.session.SetSyncTimeout(0)
	return
}

func (m *DataBase) C(colName string) *mgo.Collection {
	return m.session.DB(m.dbName).C(colName)
}

func (m *DataBase) Close() {
	m.session.Close()
	m = nil
}

func (m *DataBase) Copy() *DataBase {
	return &DataBase{
		session: m.session.Copy(),
		dbName:  m.dbName,
	}
}

func (m *DataBase) Session() *mgo.Session {
	return m.session
}
