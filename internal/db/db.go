package db

import (
	"github.com/smailic05/TaskInfoblox/internal/model"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{db}
}

func (db *DB) Load(user model.User) ([]model.User, error) {
	users := make([]model.User, 0)
	err := db.Where("username LIKE ? AND address LIKE ? AND phone LIKE ?", user.Username, user.Address, user.Phone).Find(&users).Error
	return users, err
}

func (db *DB) LoadOne(user model.User) (model.User, error) {
	err := db.Where(user).First(&user).Error
	return user, err
}

func (db *DB) Store(user model.User) error {
	err := db.Create(&user).Error
	return err
}

func (db *DB) Update(user model.User) error {
	err := db.Save(&user).Error
	return err
}

func (db *DB) DeleteUser(user model.User) error {
	err := db.Where(user).Delete(&user).Error
	return err
}
