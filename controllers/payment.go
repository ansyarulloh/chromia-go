package controllers

import (
	"chromia-api/database"
	"chromia-api/models"

	"github.com/gofiber/fiber/v2"
)

func SubscribePremium(c *fiber.Ctx) error {
	// 1. Ambil ID User yang dikasih sama Satpam JWT
	userID := c.Locals("user_id").(string)

	// 2. Cari data user di database berdasarkan ID
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User tidak ditemukan",
		})
	}

	// Cek apakah user sudah premium sebelumnya
	if user.IsPremium {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Akun kamu sudah berstatus Premium, bro!",
		})
	}

	// 3. Eksekusi Upgrade: Ubah IsPremium menjadi true
	user.IsPremium = true
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memproses upgrade Premium ke database",
		})
	}

	// 4. Kirim respon sukses ke Flutter
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Selamat! Akun berhasil di-upgrade ke Premium. Scan warna sekarang sepuasnya tanpa batas! 🚀",
		"data": fiber.Map{
			"id":         user.ID,
			"name":       user.Name,
			"is_premium": user.IsPremium,
		},
	})
}