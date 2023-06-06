package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Order struct {
		OrderID    uint      `json:"orderId" gorm:"column:orderid;primaryKey;autoIncrement;uniqueIndex"`
		CustomerID uint      `json:"customerId" gorm:"column:customerid;"`
		ProductID  uint      `json:"productId" gorm:"column:productid;"`
		Product    string    `json:"product" gorm:"column:product;"`
		Quantity   string    `json:"quantity" gorm:"column:quantity;"`
		CreatedAt  time.Time `json:"created_at" gorm:"column:createdAt;"`
		UpdatedAt  time.Time `json:"updated_at" gorm:"column:updatedAt;"`
	}
)
type Orders []Order

func (Order) TableName() string {
	return "orders" //nama table di database
}

func (r *Order) CreateOrder(db *gorm.DB) error {
	return db.Model(Order{}).Create(r).Error
}

func (r *Order) SelectByOrderId(db *gorm.DB, ID int) error {
	return db.Model(Order{}).Where("orderId = ?", ID).First(r).Error
}

func (r *Order) UpdateOrder(db *gorm.DB) error {
	return db.Model(Order{}).Omit("orderId").Where("orderId = ?", r.OrderID).Updates(r).Error
}

func (r *Order) DeleteOrder(db *gorm.DB) error {
	return db.Model(Order{}).Where("orderId = ?", r.OrderID).Delete(r).Error
}

func (u *Orders) GetOrderListByCustomerId(db *gorm.DB, customerid int) error {
	return db.Model(Order{}).Where("customerid = ?", customerid).Find(u).Error
}

func (t *Orders) SelectAllOrder(db *gorm.DB) error {
	return db.Model(Order{}).Find(t).Error
}
