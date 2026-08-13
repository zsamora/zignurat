package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	Username  string         `gorm:"column:username"`
	Password  string         `gorm:"column:password"`
	Module    uint8          `gorm:"column:module"`
	AccRole   uint8          `gorm:"column:acc_role"`
	OwnerUUID *uuid.UUID     `gorm:"column:owner_uuid;type:uuid"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	JWTToken  string         `gorm:"column:jwt_token"`
	RefrToken string         `gorm:"column:refr_token"`
}

func (Account) TableName() string {
	return "accounts"
}
