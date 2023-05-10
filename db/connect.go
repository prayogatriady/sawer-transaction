package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConnectDB struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func NewConnectDB(db_user, db_password, db_host, db_port, db_name string) *ConnectDB {
	return &ConnectDB{
		DB_USER:     db_user,
		DB_PASSWORD: db_password,
		DB_HOST:     db_host,
		DB_PORT:     db_port,
		DB_NAME:     db_name,
	}
}

func (c *ConnectDB) InitMySQL() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DB_USER, c.DB_PASSWORD, c.DB_HOST, c.DB_PORT, c.DB_NAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = c.CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func (c *ConnectDB) CreateTables(db *gorm.DB) error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
		id BIGINT NOT NULL AUTO_INCREMENT,
		user_id BIGINT NOT NULL,
		amount INT DEFAULT 0,
		sawer_user_id BIGINT NOT NULL,
		transaction_status VARCHAR(30) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL,
		INDEX (id),
		PRIMARY KEY (id)
	)`

	if err := db.Exec(query).Error; err != nil {
		return err
	}
	return nil
}
