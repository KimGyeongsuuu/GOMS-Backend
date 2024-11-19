package mysql

import (
	"fmt"

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
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		return db, nil
	}

	return nil, fmt.Errorf("failed to connect to the database: %w", err)
}
