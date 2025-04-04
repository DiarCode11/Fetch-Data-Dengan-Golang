package config

import (
	"fmt"
	"log"
	"os"
	"context"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil Konfigurasi di .env
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		fmt.Println("DATABASE URL TIDAK DITEMUKAN")
	}

	// Buat koneksi Database
	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		fmt.Println("Gagal terhubung ke database: ", err)
	}

	DB = dbPool
	fmt.Println("Berhasil terhubung ke database")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Koneksi database terputus")
	}
}