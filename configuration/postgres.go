package configuration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DB struct {
	DbPostgres *gorm.DB
}

var (
	onceSinglePostgres sync.Once
	instanceDB         *DB
)

// InitSingleDB ...
func InitSingleDB(db string, logMode bool) *gorm.DB {
	onceSinglePostgres.Do(func() {
		host := strings.Split(strings.Split(db, " port=")[0], "host=")[1]
		logs := fmt.Sprintf("[INFO] Connected to POSTGRES TYPE = %s | LogMode = %+v", host, logMode)

		sqlDB, err := sql.Open("postgres", db)
		if err != nil {
			logs = fmt.Sprintf("[ERROR] Failed to connect to POSTGRES with err %s. Config=%s", err.Error(), host)
			log.Fatalln(logs)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		dialect := postgres.New(postgres.Config{Conn: sqlDB})
		loggerLevel := logger.Error
		if logMode {
			loggerLevel = logger.Info
		}

		dbConnection, err := gorm.Open(dialect, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: time.Second,
					LogLevel:      loggerLevel,
				},
			),
		})
		if err != nil {
			logs = fmt.Sprintf("[ERROR] Failed to connect to POSTGRES with err %s. Config=%s", err.Error(), host)
			log.Fatalln(logs)
		}
		log.Println(logs)
		instanceDB = &DB{DbPostgres: dbConnection}
	})
	return instanceDB.DbPostgres
}
