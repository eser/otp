package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// generateOTP, gizli bir key kullanarak rastgele bir R üretir
// ve HMAC-SHA256 hesaplayarak R || MAC (HEX formatında) şeklinde bir çıktı oluşturur.
//
// - key: Paylaşılan gizli anahtar (ör. 16 ya da 32 byte)
// - returns: HEX string formatında "R + MAC".
func generateOTP(key []byte) (string, error) {
	// 1. Rastgele R üret (16 byte örnek; ihtiyaç durumuna göre 8, 32 vb. ayarlanabilir).
	rBytes := make([]byte, 16)
	_, err := rand.Read(rBytes)
	if err != nil {
		return "", fmt.Errorf("rastgele R üretilemedi: %v", err)
	}

	// 2. HMAC-SHA256 ile MAC hesapla: MAC = HMAC(key, R)
	mac := hmacSHA256(key, rBytes)

	// 3. Nihai OTP kodu = R || MAC.
	//   Dışarıya hex encoding ile verelim: hex(R) + hex(MAC)
	otp := hex.EncodeToString(rBytes) + hex.EncodeToString(mac)

	return otp, nil
}

// verifyOTP, generateOTP ile üretilmiş bir OTP kodunu doğrular.
//
// - key: Paylaşılan gizli anahtar
// - otpHex: HEX formatında "R + MAC" birleşimi (generateOTP fonksiyonundan dönen string)
// - returns: bool (OTP geçerliyse true, yoksa false)
func verifyOTP(key []byte, otpHex string) bool {
	// R(16 byte) + MAC(32 byte) = toplam 48 byte.
	// Dolayısıyla hex olarak 96 karakter (16*2 + 32*2) bekleriz.
	if len(otpHex) != (16+32)*2 {
		fmt.Println("OTP uzunluğu beklenen formatla uyuşmuyor.")
		return false
	}

	// 1. OTP string'ini R ve MAC parçalarına ayır
	rPartHex := otpHex[:32]   // 16 byte = 32 hex karakter
	macPartHex := otpHex[32:] // kalan 64 hex karakter = 32 byte

	rBytes, err := hex.DecodeString(rPartHex)
	if err != nil {
		fmt.Println("R decode edilemedi:", err)
		return false
	}

	macBytes, err := hex.DecodeString(macPartHex)
	if err != nil {
		fmt.Println("MAC decode edilemedi:", err)
		return false
	}

	// 2. Beklenen MAC'i yeniden hesapla
	expectedMac := hmacSHA256(key, rBytes)

	// 3. MAC karşılaştırması: Sabit zamanlı (constant time) karşılaştırma önerilir.
	// Go'nun hmac.Equal() fonksiyonunu kullanabiliriz.
	if hmac.Equal(macBytes, expectedMac) {
		return true
	}

	return false
}

// hmacSHA256, key ve data'yı kullanarak HMAC-SHA256 çıktısı döndürür.
func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil) // 32 byte (256 bit) döner
}

func main() {
	// Örnek key (16 byte). Gerçek senaryoda güvenli bir şekilde saklanmalı.
	key := []byte("supersecretkey123")

	// OTP üret
	otp, err := generateOTP(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("Üretilen OTP kodu (HEX):", otp)

	// Doğrula
	isValid := verifyOTP(key, otp)
	if isValid {
		fmt.Println("OTP doğrulaması BAŞARILI.")
	} else {
		fmt.Println("OTP doğrulaması BAŞARISIZ!")
	}

	// Geçersiz senaryo örneği: OTP’nin bir karakterini değiştirip deneyelim
	tampered := otp[:len(otp)-1] + "0" // Son hex karakterini değiştir
	isValid = verifyOTP(key, tampered)
	if isValid {
		fmt.Println("Bozulmuş OTP geçerli görünüyor, hata!")
	} else {
		fmt.Println("Bozulmuş OTP reddedildi, beklenen davranış.")
	}
}
