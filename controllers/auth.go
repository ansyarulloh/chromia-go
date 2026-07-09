package controllers

import (
	"os"
	"time"

	"chromia-api/database"
	"chromia-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struktur data untuk Register
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struktur data untuk Login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// --- FUNGSI REGISTER ---
func Register(c *fiber.Ctx) error {
	var input RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak sesuai",
		})
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email ini sudah digunakan",
		})
	}

	// Enkripsi password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengamankan password",
		})
	}

	// Buat user baru (Password diamankan)
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mendaftarkan user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Akun berhasil dibuat",
	})
}

// --- FUNGSI LOGIN ---
func Login(c *fiber.Ctx) error {
	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format data tidak sesuai",
		})
	}

	// Cari user berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Email tidak terdaftar",
		})
	}

	// Cek kecocokan password yang diketik dengan yang ada di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Password salah bro!",
		})
	}

	// Racik token JWT (Mencetak Gelang VIP)
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Berlaku 3 hari
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat token keamanan",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  t,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}