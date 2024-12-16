# otp

Buradaki yaklaşım, “OTP gibi” tek seferlik veya kısa ömürlü rastgele değerler üretirken, aynı zamanda her üretilen değerin (örneğin bir sayı dizisinin) bit seviyesinde “doğru” (otantik) olup olmadığını test etmeye yarar. Böylece her değer, içinde gizli bir imza veya pattern taşıyabilir ve bunu bitwise olarak doğrulayabilirsiniz.

## Yaklaşım Özeti

1. Rastgele bitler (R): Önce yeterince uzun ve kriptografik olarak güvenli rastgele bir bit dizisi üretin (örneğin 128 bit).

2. MAC/HMAC veya imza (S): Bu rastgele bitler için gizli anahtarınızı kullanarak bir bütünlük kontrol değeri (MAC – Message Authentication Code) veya bir dijital imza oluşturun.

3. Birleştirme: Nihai OTP benzeri değeri üretmek için rastgele bitlerle (R) hesaplanan MAC/imzanın bir kısmını (S) birleştirin. Böylece ortaya çıkan kod, hem rastgele hem de bitwise doğrulanabilir bir yapı taşır.

4. Doğrulama: Alıcı taraf, aynı gizli anahtarla (veya uygun doğrulama mekanizmasıyla) üretilen MAC/imza kısmını kontrol eder. Böylece kodun gerçekten o anahtarla ve belirlenmiş pattern ile üretildiği bit düzeyinde tespit edilmiş olur.

Bu yöntemi hem sayaç tabanlı (HOTP) hem de zaman tabanlı (TOTP) OTP tasarımlarına uygulayabilir veya sadece “tek seferlik rastgele sayı + gizli anahtarla imza” şeklinde düzenleyebilirsiniz.

## Adım Adım Uygulama

### 1. Gizli Anahtar Oluşturma
- Öncelikle, sistemde herkesin erişemeyeceği güvenli bir gizli anahtar (secret key) oluşturun.

- Anahtar yeterince uzun olmalı. Mesela 128 bit veya 256 bit’lik bir key yaygın olarak kullanılır.

Bu anahtar ileride hem rastgele sayıları üretirken hem de imza (veya MAC) oluştururken kullanılacaktır.

### 2. Rastgele Bit Dizisi (R) Oluşturma
- Bir kriptografik rastgele sayı üreteci (CSPRNG: Cryptographically Secure PseudoRandom Number Generator) kullanarak 128 bit’lik (veya ihtiyaca göre daha uzun/kısa) bir rastgele bit dizisi üretin.

- Bu dizi “OTP” görevini görebilecek temel verilerdir.

Örnek: 128 bit’lik R dizimiz şu şekilde olsun:

R = 11001010 ... (toplam 128 bit) ...

### 3. MAC veya İmza Hesaplama (S)
- Rastgele bit dizisinin (R) bütünlüğünü ve kaynaktan (gizli anahtardan) geldiğini doğrulamak için bir MAC algoritması kullanabilirsiniz. (HMAC-SHA256 gibi.)

- Alternatif olarak asimetrik kriptografi kullanılarak dijital imza da üretilebilir (ör. ECDSA, RSA gibi), ancak bu daha maliyetli ve gereksiz olabilir. HMAC genelde OTP sistemlerinde yeterlidir.

Örnek: HMAC-SHA256 kullandığınızı varsayalım:

S = HMAC-SHA256(Key, R)

Bu işlem 256 bit’lik bir çıktı verecektir.

### 4. Nihai Değerin Oluşturulması
- Üretilen 128 bit’lik rastgele dizi R ve 256 bit’lik MAC çıktısı S, tek bir “OTP kodu” gibi düşünülebilir. Fakat son kullanıcıya çok uzun bir bit dizisi iletmek pratik olmayabilir.

- Bunun yerine “R + S’in bir bölümü” birleştirilerek daha kısa, ama doğrulaması bit bazında yapılabilen bir kod üretilir.

Örneğin:

1. Rastgele 128 bit “R”

2. MAC (HMAC) sonucu 256 bit “S”

3. Nihai kod olarak R || (S’nin ilk N biti) kullanabilirsiniz. Mesela S’nin ilk 32 bit’ini alıp eklemek pratik bir yöntemdir:

OTP = R (128 bit) || S[0:32]  (32 bit)

Böylece toplam 160 bit’lik bir yapı elde etmiş olursunuz.

#### Bitwise Pattern
Kodun “içinde” özel bir bit düzeni olduğunu söylemek için, MAC’in belli bitlerinin veya R’nin içindeki belli bitlerin belirli bir pattern’a uygun olması sağlanabilir. Ancak genel yaklaşım şu olur:

- Siz rastgele R’yi ürettikten sonra MAC hesaplıyorsunuz. MAC de key’e ve R’ye bağlı olduğu için bit düzeyinde bir “imza/pattern” taşıyor.

- Doğrulama aşamasında MAC kısmını R’yi de kullanarak yeniden hesapladığınızda bit düzeyinde tamamen eşleşip eşleşmediğini kontrol ediyorsunuz.

### 5. Doğrulama

Kod (OTP) tarafınıza ulaştığında doğrulamak için şu işlemleri yaparsınız:

1. OTP içinden R ve eklenmiş MAC parçasını ayırın.

2. Aynı gizli anahtarı (Key) ve R’yi kullanarak S' = HMAC-SHA256(Key, R) hesaplayın.

3. S’ çıktısının aynı bit dilimlerini alın (mesela ilk 32 bit). Bunu OTP içindeki MAC dilimiyle karşılaştırın.

4. Eğer birebir eşleşiyorsa kod bit düzeyinde doğrulanır; aksi takdirde sahte veya bozunmuş kabul edilir.

## Ek Güvenlik ve Tasarım İpuçları

### 1. Sayaç veya Zaman Tabanlı Tekrar Engeli
OTP’nin yeniden kullanılmasını engellemek istiyorsanız, R’yi doğrudan rastgele yerine, bir sayaç (counter) ya da zaman damgası ile birleştirebilirsiniz (ör. TOTP gibi).

- TOTP’de: R = HMAC-SHA1(Key, ZamanBlok) gibi bir şey üretilir. Her 30 saniyede bir değişir. Ardından tekrar MAC alarak imza eklenir.

- HOTP’de: Sayaç her kullanımda artar.

### 2. Uzunluk Ayarı
Kodun pratikte fazla uzun olmaması için “R” ve “S” değerlerinin bit sayısı kısaltılabilir. Ama unutmayın: Ne kadar çok kısaltırsanız, brute force veya çakışma riski de o kadar artar.

### 3. Asimetrik İmza Alternatifi
Eğer aynı gizli anahtarı paylaşmak istemiyorsanız, asimetrik kriptografi (ör. RSA, ECDSA) kullanarak R’ye dijital imza da ekleyebilirsiniz. Karşı taraf, sizin public key’inizle imzayı doğrulayabilir. Böylece bitwise olarak yine “bu sayı belli bir pattern/imza taşıyor mu?” diye kontrol edilebilir. Fakat OTP senaryolarında HMAC daha yaygındır.

### 4. Hash ve Pattern Gömme
“Bitwise pattern”i tam olarak manuel şekilde gömmek istiyorsanız, R üretildikten sonra R’nin belli bitlerini bir pattern’e zorlayacak şekilde seçebilirsiniz (örneğin R’nin son 8 bitini sabit bir mask ile birleştirmek). Ardından HMAC alırsınız. Böylece hem rastgelelik korunur hem de kod içinde “göze çarpan” bir bit pattern yer alır. Doğrulamada, önce pattern kontrol edilir, sonra HMAC eşleşmesi yapılır.

## Sonuç

Bu yaklaşım sayesinde:

- Her üretilen sayı (OTP) bit düzeyinde rastgelelik içerir.

- Bitwise doğrulama yapabilmenizi sağlayacak bir MAC/imza parçası da bu sayıların içinde gömülü durur.

- Kod size ulaştığında, aynı gizli anahtar ve algoritmayla tekrar MAC hesaplayıp bit bit kıyaslayarak kodun otantik olup olmadığını çok hızlı tespit edebilirsiniz.

Gerçek uygulamada HMAC tabanlı yaklaşım, hem basitliği hem de güçlü güvenlik özellikleri nedeniyle yaygın olarak kullanılır. Eğer istemci ve sunucu aynı gizli anahtarı paylaşıyorsa (paylaşımlı anahtar modeli), her üretilen rastgele sayının yanında kendi MAC değerini de saklamak OTP algoritması için yeterlidir.

Bu şekilde, bitwise olarak her sayının (OTP’nin) gerçekten sizin belirlediğiniz pattern/gizli anahtar çerçevesinde üretildiğini doğrulayabilirsiniz.
