package database

import (
	"database/sql"
	"fmt"
	"log"
	"working-day-api/config"
	"working-day-api/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Gorm struct {
	*gorm.DB
}

var DB *gorm.DB

func Connection(config *config.AppVars) {

	dsn_write := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		config.DB.DBParams.Username,
		config.DB.DBParams.Password,
		config.DB.DBParams.Host,
		3306,
		config.DB.DBParams.Name,
	)

	sqlDB, err := sql.Open("mysql", dsn_write)
	if err != nil {
		log.Println("An error occurred, please try again later. (2)")
		return
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(config.DB.DBParams.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.DB.DBParams.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(config.DB.DBParams.MaxLifetime)

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:  dsn_write,
		Conn: sqlDB,
	}), &gorm.Config{})

	err = DB.AutoMigrate(&domain.Role{}, &domain.User{}, &domain.Task{})

	if err != nil {
		log.Println("An error occurred, please try again later. (3)", err)
		return
	}

	SeedRoles()
}
