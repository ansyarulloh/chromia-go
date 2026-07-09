package database

import (
	"fmt"
	"log"
	"os"

	// Sesuaikan "chromia-api" dengan nama module di go.mod kamu
	"chromia-api/models" 

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Yah, gagal konek ke database! \n", err)
	}

	fmt.Println("Mantap! Sukses konek ke PostgreSQL 🎉")
	
	// --- TAMBAHKAN BARIS INI ---
	fmt.Println("Sedang menjalankan Auto Migration...")
	err = database.AutoMigrate(&models.User{}, &models.ScanHistory{})
	if err != nil {
		log.Println("Gagal melakukan migrasi tabel:", err)
	} else {
		fmt.Println("Database Migration sukses! Tabel 'users' siap digunakan 🚀")
	}

	DB = database
}