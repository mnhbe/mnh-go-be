package model

import "time"

type Role int

var (
  AdminSystem  Role = 1
  AdminRegion  Role = 2
  AdminCompany Role = 3
)

func (r *Role) ToString() string {
  switch *r {
  case AdminSystem:
    return "AdminSystem"
  case AdminRegion:
    return "AdminRegion"
  case AdminCompany:
    return "AdminCompany"
  }
  return ""
}

type UserRole struct {
  ID        uint `gorm:"primary_key"`
  Name      string
  CreatedAt time.Time  `gorm:"column:createdAt"`
  UpdatedAt time.Time  `gorm:"column:updatedAt"`
  DeletedAt *time.Time `gorm:"column:deletedAt"`
}

func (UserRole) TableName() string { return "UserRoles" }
