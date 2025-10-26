package model

import (
  "strings"
  "time"

  "github.com/jinzhu/gorm"
)

type User struct {
  ID        uint    `gorm:"primary_key,column:id"`
  FirstName string  `gorm:"column:firstName"`
  LastName  string  `gorm:"column:lastName"`
  Phone     *string `gorm:"column:phone"`
  Email     *string `gorm:"column:email"`
  Password  *string `gorm:"column:password"`
  AvatarUrl string  `gorm:"column:avatarUrl"`
  IsLock    bool
  LastLogin *time.Time `gorm:"column:lastLogin"`

  UserRole   UserRole `gorm:"foreignkey:UserRoleId"`
  UserRoleID uint     `gorm:"column:UserRoleId"`

  CreatedAt time.Time `gorm:"column:createdAt"`
  UpdatedAt time.Time `gorm:"column:updatedAt"`
  DeletedAt time.Time `gorm:"column:deletedAt"`
}

func (u User) TableName() string { return "Users" }

func (u User) AutoMigrate(db *gorm.DB) {
  db.AutoMigrate(&User{})
}

func (u User) DisplayName() string {
  return *u.Email
}

func (u User) GetFullName() string {
  s := strings.TrimSpace(strings.TrimSpace(u.FirstName) + " " + strings.TrimSpace(u.LastName))
  return s
}
