package db

import (
	CommentModels "fiber-crud/internal/domain/comment"
	ProductModels "fiber-crud/internal/domain/product"
	userModels "fiber-crud/internal/domain/user"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	dsn := "host=localhost user=postgres password=maulana dbname=jajal port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database: %v", err)
	}

	db = db.Debug()

	if err := db.AutoMigrate(
		&userModels.User{},
		&ProductModels.Product{},
		&CommentModels.Comment{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
