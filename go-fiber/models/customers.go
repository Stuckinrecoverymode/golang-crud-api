package models

import "gorm.io/gorm"

type Customer struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	CustName    *string `json:"customer"`
	PhoneNumber *string `json:"phone"`
	Date        *string `json:"date"`
	State       *string `json:"state"`
	DeviceName  *string `json:"device"`
	Problem     *string `json:"problem"`
	Mail        *string `json:"mail"`
	Price       *string `json:"price"`
}

func MigrateCustomers(db *gorm.DB) error {
	err := db.AutoMigrate(&Customer{})
	return err
}