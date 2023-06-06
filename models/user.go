package models

import (
	"errors"
	"html"
	"strings"
	"tes_dbo/utils/token"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User struct {
		UserId      uint      `json:"userId" gorm:"column:userId;primaryKey;autoIncrement;"`
		FullName    string    `json:"fullName" gorm:"column:fullName;size:50;" validate:"required"`
		Username    string    `json:"username" gorm:"size:255;not null;unique;"`
		Password    string    `json:"password" gorm:"size:255;not null;"`
		Email       string    `json:"email" gorm:"column:email;size:50;" validate:"required,email"`
		PhoneNumber string    `json:"phoneNumber" gorm:"column:phoneNumber;type:varchar(20);" validate:"required"`
		Role        string    `json:"role" gorm:"column:role;" validate:"required"`
		CreatedAt   time.Time `json:"created_at" gorm:"column:createdAt;"`
		UpdatedAt   time.Time `json:"updated_at" gorm:"column:updatedAt;"`
	}
)
type Users []User

func (User) TableName() string {
	return "users" //nama table di database
}
func (u *User) InsertUser(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return db.Model(User{}).Create(u).Error
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(db *gorm.DB, username string, password string) (string, error) {

	var err error

	u := User{}

	err = db.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.UserId, u.Username, u.Role)

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserByID(db *gorm.DB, uid uint) (User, error) {

	var u User

	if err := db.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}
