package controllers

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"chromia-api/database"
	"chromia-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ScanColor(c *fiber.Ctx) error {
	// 1. Ambil ID User dari Satpam (Middleware)
	userID := c.Locals("user_id").(string)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User tidak ditemukan"})
	}

	// 2. Cek Kuota (Maksimal 5 kali untuk akun Free)
	if !user.IsPremium && user.DailyScanCount >= 5 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Limit harian habis. Silakan upgrade ke Premium untuk scan sepuasnya!",
		})
	}

	log.Println("1. Menerima request dari user:", userID)
	// 3. Tangkap file gambar dari request
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Error: Gambar tidak ditemukan")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Gambar tidak ditemukan, pastikan mengirim field 'image'"})
	}
	log.Println("2. File gambar diterima, ukuran:", file.Size)

	fileData, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal membuka gambar"})
	}
	defer fileData.Close()

	imgBytes, err := io.ReadAll(fileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal membaca data gambar"})
	}

	// Ambil format ekstensi gambar dan buang tulisan "image/"
	mimeType := file.Header.Get("Content-Type") 
	formatGambar := "jpeg" // Default
	if strings.Contains(mimeType, "/") {
		formatGambar = strings.Split(mimeType, "/")[1] 
	}

	// 4. Inisialisasi Klien Google Gemini AI
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyambung ke server AI"})
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.1-flash-lite")

	// 5. Siapkan Prompt / Perintah untuk AI
	promptText := "Kamu adalah asisten pendeteksi warna. Fokus HANYA pada benda yang berada tepat di TENGAH gambar. Apa warna dari benda di tengah tersebut? Jawab HANYA dengan format persis seperti ini tanpa kata-kata lain: Nama Warna | #KodeHex (contoh: Merah Gelap | #8B0000)."
	prompt := genai.Text(promptText)
	imgPart := genai.ImageData(formatGambar, imgBytes) 

	log.Println("3. Mulai mengirim gambar ke server Gemini AI... (Format:", formatGambar, ")")

	// 6. Tembak Permintaan ke Server Google
	resp, err := model.GenerateContent(ctx, prompt, imgPart)
	log.Println("4. Selesai menerima balasan dari Gemini AI")
	if err != nil {
		// Kita print error aslinya ke terminal biar keliatan jelas
		log.Println("Penyebab Google nolak:", err) 
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", 
			"message": "AI gagal memproses warna dari gambar. Cek terminal Go!", 
		})
	}

	// 7. Ambil Teks Jawaban dari Gemini dan Parse Formatnya
	var hasilWarna string
	var kodeHex string
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			hasilKotor := strings.TrimSpace(string(txt))
			// Hilangkan markdown bold/italic jika ada (Gemini suka bandel)
			hasilKotor = strings.ReplaceAll(hasilKotor, "*", "")
			
			pecahan := strings.Split(hasilKotor, "|")
			if len(pecahan) >= 2 {
				hasilWarna = strings.TrimSpace(pecahan[0])
				kodeHex = strings.TrimSpace(pecahan[1])
			} else {
				hasilWarna = hasilKotor
				kodeHex = "#808080" // Default Abu-abu
			}
		}
	}

	// 8. Potong Kuota Scan (Hanya untuk user Free)
	if !user.IsPremium {
		user.DailyScanCount += 1
		database.DB.Save(&user)
	}

	// Simpan ke Riwayat Scan
	riwayat := models.ScanHistory{
		UserID:    userID,
		ColorName: hasilWarna,
		HexCode:   kodeHex,
	}
	database.DB.Create(&riwayat)

	// Hitung sisa kuota untuk ditampilkan di UI
	sisaKuota := 5 - user.DailyScanCount
	if user.IsPremium {
		sisaKuota = 9999 // Unlimited
	} else if sisaKuota < 0 {
		sisaKuota = 0
	}

	// 9. Kembalikan Jawaban ke Klien (Flutter)
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"color":           hasilWarna,
			"hex":             kodeHex,
			"remaining_scans": sisaKuota,
			"is_premium":      user.IsPremium,
		},
	})
}