package db

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func Migrate(st interface{}) error {
	return db.AutoMigrate(st)
}

func Create(ctx context.Context, data interface{}) error {
	return db.WithContext(ctx).Create(data).Error
}

func First(ctx context.Context, data interface{}, args ...interface{}) error {
	return db.WithContext(ctx).First(data, args).Error
}

func Save(ctx context.Context, data interface{}) error {
	return db.WithContext(ctx).Save(data).Error
}

func FindByIds(ctx context.Context, model interface{}, ids []int) error {
	return db.WithContext(ctx).Find(model, ids).Error
}
