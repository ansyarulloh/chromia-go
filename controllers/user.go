package controllers

import (
	"time"

	"chromia-api/database"
	"chromia-api/models"

	"github.com/gofiber/fiber/v2"
)

func GetUserStatus(c *fiber.Ctx) error {
	// 1. Ambil user_id yang dititipkan oleh si Satpam (Middleware) di Context Locals
	userID := c.Locals("user_id").(string)

	// 2. Cari data user lengkapnya di database berdasarkan ID tersebut
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User tidak ditemukan",
		})
	}

	// 3. LOGIKA RESET KUOTA (Cek apakah sudah ganti hari)
	// Kita bandingkan tanggal hari ini dengan tanggal terakhir kali user melakukan aktivitas
	formatTanggal := "2006-01-02" // Format standar YYYY-MM-DD di Go
	hariIni := time.Now().Format(formatTanggal)
	hariTerakhirScan := user.LastScanDate.Format(formatTanggal)

	if hariIni != hariTerakhirScan {
		// Kalau tanggalnya beda, berarti udah ganti hari. Reset hitungan jadi 0!
		user.DailyScanCount = 0
		user.LastScanDate = time.Now()
		database.DB.Save(&user) // Simpan perubahan terbaru ke PostgreSQL
	}

	// 4. Kalkulasi sisa kuota scan untuk dikirim ke Flutter
	limitMaksimal := 5
	sisaKuota := limitMaksimal - user.DailyScanCount
	if sisaKuota < 0 {
		sisaKuota = 0
	}

	// 5. Ambil Riwayat Scan
	var history []models.ScanHistory
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(5).Find(&history)

	// 6. Kirim respon sukses beserta data kuota ter-update dan riwayat
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":              user.ID,
			"name":            user.Name,
			"email":           user.Email,
			"avatar_url":      user.AvatarURL,
			"is_premium":      user.IsPremium,
			"remaining_scans": sisaKuota, 
			"history":         history,
		},
	})
}