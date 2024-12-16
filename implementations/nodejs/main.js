const crypto = require('node:crypto');

/**
 * generateOTP
 * 
 * - Gizli anahtar (key) kullanarak 16 byte rastgele R üretir.
 * - HMAC-SHA256(R) hesaplar.
 * - R ve MAC değerlerini hex formatında birleştirerek döndürür.
 * 
 * @param {Buffer} key - Gizli anahtar (Buffer)
 * @returns {string} Hex formatında "R + MAC"
 */
function generateOTP(key) {
  // 1. 16 byte rastgele R üret
  const r = crypto.randomBytes(16);

  // 2. MAC hesapla: MAC = HMAC-SHA256(key, R)
  const hmac = crypto.createHmac('sha256', key);
  hmac.update(r);
  const mac = hmac.digest(); // raw Buffer

  // 3. OTP olarak hex(R) + hex(MAC)
  const otp = r.toString('hex') + mac.toString('hex');
  return otp;
}

/**
 * verifyOTP
 * 
 * - Üretilmiş OTP değerini doğrular.
 * - Hex formatındaki OTP içinden R ve MAC parçalarını ayrıştırır.
 * - Aynı key ile HMAC hesaplar, karşılaştırır.
 * 
 * @param {Buffer} key - Gizli anahtar
 * @param {string} otp - Hex formatında "R + MAC"
 * @returns {boolean} true ise geçerli, false aksi halde
 */
function verifyOTP(key, otp) {
  // Toplam 48 byte => 16 (R) + 32 (MAC), hex olarak 96 karakter
  if (otp.length !== 96) {
    return false;
  }

  // R: ilk 32 hex karakter => 16 byte
  const rHex = otp.slice(0, 32);
  // MAC: kalan 64 hex karakter => 32 byte
  const macHex = otp.slice(32);

  const r = Buffer.from(rHex, 'hex');
  const mac = Buffer.from(macHex, 'hex');

  // 2. Beklenen MAC'i yeniden hesapla
  const hmac = crypto.createHmac('sha256', key);
  hmac.update(r);
  const expectedMac = hmac.digest();

  // 3. Sabit zamanlı karşılaştırma (crypto.timingSafeEqual)
  if (mac.length !== expectedMac.length) {
    return false;
  }
  return crypto.timingSafeEqual(mac, expectedMac);
}

// =========== Örnek Kullanım ===========

// Gizli anahtar (örnek). Gerçekte güvenli şekilde saklayın.
const key = Buffer.from('supersecretkey123', 'utf8'); // 16 byte

// OTP üret
const otp = generateOTP(key);
console.log("Üretilen OTP (hex):", otp);

// Doğrula
if (verifyOTP(key, otp)) {
  console.log("OTP doğrulaması BAŞARILI.");
} else {
  console.log("OTP doğrulaması BAŞARISIZ!");
}

// Bozulmuş OTP örneği
const tampered = `${otp.slice(0, -1)}0`; // Son hex karakterini değiştir
if (verifyOTP(key, tampered)) {
  console.log("Bozulmuş OTP geçerli görünüyor, hata!");
} else {
  console.log("Bozulmuş OTP reddedildi, beklenen davranış.");
}
