package iris_lib

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func InitDb(cfg *MysqlDB, pool *MysqlPool, log *LogConf) *gorm.DB {
	gormConfig := &gorm.Config{}
	if cfg.EnableSqlLog {
		gormConfig = &gorm.Config{Logger: NewGormLog(log.FilePath, log.FileName)}
	}
	return createSqlConnection(cfg.Conn, pool.MaxIdleConn, pool.MaxOpenConn, pool.ConnMaxLifetime, gormConfig)
}

func createSqlConnection(dsn string, maxIdleConn, maxOpenConn int, connMaxLifetime time.Duration, gormConfig *gorm.Config) *gorm.DB {
	if db, err := gorm.Open(mysql.Open(dsn), gormConfig); err != nil {
		SystemLogger.Error("mysql gorm v2 init failed", zap.String("dsn", dsn), zap.Error(err))
		panic(err)
	} else {
		Mysqldb, _ := db.DB()
		Mysqldb.SetMaxIdleConns(maxIdleConn)
		Mysqldb.SetMaxOpenConns(maxOpenConn)
		Mysqldb.SetConnMaxLifetime(time.Minute * connMaxLifetime)
		SystemLogger.Info("mysql gorm v2 init success")
		return db
	}

}

func NewGormLog(logPath string, applicationName string) logger.Interface {
	infoPath := logPath + "/" + applicationName + "-" + "sql.log"
	return logger.New(
		log.New(NewWriter(infoPath), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)
}
