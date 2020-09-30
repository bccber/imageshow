package db

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"imageshow/conf"
)

var (
	spiderDB *sql.DB
	masterDB []*sql.DB
	slaveDB  []*sql.DB
)

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()

	return db, err
}

func init() {
	var err error

	// 初始化爬虫库
	spiderDB, err = openDB(conf.Config.SpiderDB)
	if err != nil {
		panic(err)
	}

	// 初始化主库
	for _, connStr := range conf.Config.MasterDB {
		db, err := openDB(connStr)
		if err != nil {
			panic(err)
		}

		masterDB = append(masterDB, db)
	}

	// 初始化从库
	for _, connStr := range conf.Config.SlaveDB {
		db, err := openDB(connStr)
		if err != nil {
			panic(err)
		}

		slaveDB = append(slaveDB, db)
	}
}

// getMasterDB 通过id取模获取主库对象
func getMasterDB(id int64) (*sql.DB, error) {
	dbLen := int64(len(masterDB))
	if dbLen <= 0 {
		return nil, errors.New("MasterDB为空")
	}
	dbIdx := id % dbLen
	if dbIdx < 0 || dbIdx >= dbLen {
		return nil, errors.New("MasterDB dbIdx超出范围")
	}

	return masterDB[dbIdx], nil
}

// getSlaveDB 通过id取模获取从库对象
func getSlaveDB(id int64) (*sql.DB, error) {
	dbLen := int64(len(slaveDB))
	if dbLen <= 0 {
		return nil, errors.New("SlaveDB为空")
	}
	dbIdx := id % dbLen
	if dbIdx < 0 || dbIdx >= dbLen {
		return nil, errors.New("SlaveDB dbIdx超出范围")
	}

	return slaveDB[dbIdx], nil
}
