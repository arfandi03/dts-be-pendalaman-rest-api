package database

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"-" `
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

func (auth *Auth) SignUp(db *gorm.DB) error {
	// select * from users where username="user"
	fmt.Println("cari user")
	if err := db.Where(&Auth{Username: auth.Username}).First(auth).Error; err != nil {
		fmt.Println("masuk")
		if err == gorm.ErrRecordNotFound { // jika tidak ditemukan
			fmt.Println("user tidak ada")
			res := db.Create(auth)
			if res.Error != nil { // membuat user baru
				return res.Error
			}
		}
	} else {
		return errors.Errorf("Duplicate username")
	}
	return nil
}

func (auth *Auth) Login(db *gorm.DB) (*Auth, error) {
	if err := db.Where(&Auth{Username: auth.Username, Password: auth.Password}).First(auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound { // jika tidak ditemukan
			return nil, errors.Errorf("Incorect email or password")
		}
	}
	return auth, nil
}

func ValidateAuth(token string, db *gorm.DB) (*Auth, error) {
	var auth Auth

	if err := db.Where(&Auth{Token: token}).First(&auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound { // jika tidak ditemukan
			return nil, errors.Errorf("Invalid token")
		}
	}
	return &auth, nil
}
