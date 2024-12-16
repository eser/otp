<?php

/**
 * generateOTP
 * 
 * - Gizli anahtar (key) kullanarak 16 byte rastgele R üretir.
 * - HMAC-SHA256(R) hesaplar.
 * - R ve MAC değerlerini hex formatında birleştirerek döndürür.
 * 
 * @param string $key Gizli anahtar (16 veya 32 byte olabilir)
 * @return string Hex formatında "R + MAC"
 */
function generateOTP($key) {
    // 1. Rastgele 16 byte üret
    $r = random_bytes(16);

    // 2. HMAC-SHA256 hesapla (binary output)
    $mac = hash_hmac('sha256', $r, $key, true); // 'true' => raw binary output

    // 3. R + MAC birleşimi (hex)
    $otp = bin2hex($r) . bin2hex($mac);

    return $otp;
}

/**
 * verifyOTP
 * 
 * - Üretilmiş OTP değerini doğrular.
 * - Hex formatındaki OTP içinden R ve MAC parçalarını ayrıştırır.
 * - Aynı key ile HMAC hesaplar, karşılaştırır.
 * 
 * @param string $key Gizli anahtar
 * @param string $otp Hex formatında "R + MAC"
 * @return bool Geçerliyse true, değilse false
 */
function verifyOTP($key, $otp) {
    // Beklenen toplam uzunluk = 16 byte (R) + 32 byte (MAC) = 48 byte.
    // Hex string olarak 48*2 = 96 karakter.
    if (strlen($otp) !== 96) {
        // Uzunluk beklenenden farklı
        return false;
    }

    // 1. Parçaları ayır
    $rHex = substr($otp, 0, 32);   // ilk 32 hex karakter (16 byte)
    $macHex = substr($otp, 32);   // geri kalan 64 hex karakter (32 byte)
    
    $r = hex2bin($rHex);
    $mac = hex2bin($macHex);

    if ($r === false || $mac === false) {
        return false;
    }

    // 2. Beklenen MAC’i yeniden hesapla
    $expectedMac = hash_hmac('sha256', $r, $key, true);

    // 3. Sabit zamanlı (constant time) karşılaştırma
    return hash_equals($mac, $expectedMac);
}

// =========== Örnek Kullanım ===========

// Gizli anahtar (örnek). Gerçekte güvenli şekilde saklayın.
$key = "supersecretkey123"; // 16 byte

// OTP üret
$otp = generateOTP($key);
echo "Üretilen OTP (hex): $otp\n";

// Doğrula
if (verifyOTP($key, $otp)) {
    echo "OTP doğrulaması BAŞARILI.\n";
} else {
    echo "OTP doğrulaması BAŞARISIZ!\n";
}

// Bozulmuş OTP örneği
$tampered = substr($otp, 0, -1) . '0'; // Son karakteri değiştir
if (verifyOTP($key, $tampered)) {
    echo "Bozulmuş OTP geçerli görünüyor, hata!\n";
} else {
    echo "Bozulmuş OTP reddedildi, beklenen davranış.\n";
}
