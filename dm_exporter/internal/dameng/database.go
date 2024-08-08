package dameng

import (
	"database/sql"

	_ "dm"
	"dm_exporter/global"

	"github.com/prometheus/common/log"
)

func Connect(dsn string) *sql.DB {
	log.Debugln("Launching connection: ", dsn)
	db, err := sql.Open("dm", dsn)
	if err != nil {
		log.Errorln("Error while connecting to", dsn)
		panic(err)
	}
	log.Debugln("set max idle connections to ", *global.MaxIdleConns)
	db.SetMaxIdleConns(*global.MaxIdleConns)
	log.Debugln("set max open connections to ", *global.MaxOpenConns)
	db.SetMaxOpenConns(*global.MaxOpenConns)
	log.Debugln("Successfully connected to: ", dsn)
	return db
}
