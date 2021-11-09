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
	result := db.Where("username LIKE ? AND address LIKE ? AND phone LIKE ?", user.Username, user.Address, user.Phone).Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

func (db *DB) LoadOne(user model.User) (model.User, error) {
	result := db.Where(user).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (db *DB) Store(user model.User) error {
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *DB) Update(user model.User) error {
	err := db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteUser(user model.User) error {
	err := db.Where(user).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
