package mysql

import (
	"fmt"
	"log"
	"time"

	"GOMS-BACKEND-GO/global/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLConnection() (*gorm.DB, error) {
	user := config.Data().Mysql.User
	password := config.Data().Mysql.Pass
	host := config.Data().Mysql.Host
	port := config.Data().Mysql.Port
	database := config.Data().Mysql.Db

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", user, password, host, port, database)
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("Retrying MySQL connection (%d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to the database: %w", err)
}
