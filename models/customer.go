package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Customer struct {
		CustomerID  uint      `json:"customerId" gorm:"column:customerid;primaryKey;autoIncrement;"`
		FullName    string    `json:"fullName" gorm:"column:fullName;size:50;" validate:"required"`
		Username    string    `json:"username" gorm:"size:255;not null;unique;"`
		Password    string    `json:"password" gorm:"size:255;not null;"`
		Email       string    `json:"email" gorm:"column:email;size:50;" validate:"required,email"`
		PhoneNumber string    `json:"phoneNumber" gorm:"column:phoneNumber;type:varchar(20);" validate:"required"`
		CreatedAt   time.Time `json:"created_at" gorm:"column:createdAt;"`
		UpdatedAt   time.Time `json:"updated_at" gorm:"column:updatedAt;"`
	}
)
type Customers []Customer

func (Customer) TableName() string {
	return "customers" //nama table di database
}

func (t *Customer) InsertCustomer(db *gorm.DB) error {
	return db.Model(Customer{}).Create(t).Error
}

func (r *Customer) SelectById(db *gorm.DB, ID int) error {
	return db.Model(Customer{}).Where("customerId = ?", ID).First(r).Error
}

func (r *Customer) UpdateCustomer(db *gorm.DB) error {
	return db.Model(Customer{}).Omit("customerid").Where("customerid = ?", r.CustomerID).Updates(r).Error
}

func (r *Customer) DeleteCustomer(db *gorm.DB) error {
	return db.Model(Customer{}).Where("customerid = ?", r.CustomerID).Delete(r).Error
}

func (r *Customers) SelectAll(db *gorm.DB) error {
	return db.Model(Customer{}).Find(r).Error
}
