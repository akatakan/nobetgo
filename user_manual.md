# NöbetGo Kullanım Kılavuzu

## İçindekiler
1. [Giriş](#1-giriş)
2. [Kurulum ve Başlangıç](#2-kurulum-ve-başlangıç)
3. [Tanımlamalar](#3-tanımlamalar)
    - [Bölüm Yönetimi](#31-bölüm-yönetimi)
    - [Nöbet Tipleri](#32-nöbet-tipleri)
    - [Personel Yönetimi](#33-personel-yönetimi)
4. [Nöbet Planlama](#4-nöbet-planlama)
    - [Otomatik Planlama Sihirbazı](#41-otomatik-planlama-sihirbazı)
    - [Takvim ve Manuel Düzenleme](#42-takvim-ve-manuel-düzenleme)
5. [Puantaj ve Devam Takibi](#5-puantaj-ve-devam-takibi)
6. [Raporlar](#6-raporlar)

---

## 1. Giriş
**NöbetGo**, sağlık personeli ve vardiyalı çalışan ekipler için geliştirilmiş, adil ve optimize edilmiş nöbet çizelgeleri oluşturan modern bir web uygulamasıdır.

## 2. Kurulum ve Başlangıç
Uygulama yerel ağınızda veya sunucunuzda çalışır. Tarayıcınızdan (Chrome, Edge vb.) belirtilen adrese (örn. `http://localhost:5173`) giderek erişebilirsiniz.

### 2.1 Şifre Sıfırlama
Şifrenizi unutursanız:
1. Giriş sayfasındaki **"Şifremi Unuttum"** bağlantısına tıklayın.
2. E-posta adresinizi girin.
3. Sistem yöneticiniz şifre sıfırlama bağlantısını terminal kayıtlarından (Log) kontrol edebilir ve size iletebilir.
4. Bağlantıya tıklayarak yeni şifrenizi belirleyebilirsiniz (En az 6 karakter).

### 2.2 Güvenlik
Sistem, kötü niyetli denemeleri engellemek için **IP tabanlı hız sınırlaması** uygular. Arka arkaya çok fazla hatalı giriş denemesi yaparsanız erişiminiz kısa süreliğine kısıtlanabilir.

## 3. Tanımlamalar
Sistemi kullanmaya başlamadan önce temel tanımları yapmanız gerekir.

### 3.1 Bölüm Yönetimi
Sol menüden **Bölümler** sayfasına gidin.
- **Yeni Bölüm Ekle**: Sağ üstteki "Yeni Bölüm" butonuna tıklayın.
    - **Bölüm Adı**: Örn. "Cerrahi", "Acil".
    - **Kat**: Kat numarasını girin. **0 (Zemin)** ve **negatif değerler (-1, -2 vb.)** desteklenmektedir.
- **Düzenle/Sil**: Mevcut kartların üzerindeki kalem veya çöp kutusu ikonlarını kullanın.

### 3.2 Nöbet Tipleri
Sol menüden **Nöbet Tipleri** sayfasına gidin.
- **Yeni Nöbet Tipi**: "Yeni Nöbet Tipi" butonuna tıklayın.
    - **Ad**: Örn. "Gündüz", "Gece", "24 Saat".
    - **Renk**: Takvimde görünecek rengi seçin.
    - **Saatler**: Başlangıç ve bitiş saatlerini girin.
- **Önemli**: Gece nöbetleri için bitiş saati ertesi güne sarkıyorsa sistem bunu otomatik algılar ve süreyi doğru hesaplar.

### 3.3 Personel Yönetimi
Sol menüden **Personel** sayfasına gidin.
- **Yeni Personel**: Ad, Soyad, Ünvan ve Bölüm bilgilerini girerek personel ekleyin.
- **Excel İçe Aktar**: Toplu personel yüklemek için kullanılır.
    - **Excel Dosya Formatı**:
        - 1. Satır: Başlıklar (Sırayla: Ad, Soyad, E-posta, Telefon, Ünvan, Bölüm, Saatlik Ücret)
        - 2. Satırdan itibaren veriler girilmelidir.
        - **Ünvan** ve **Bölüm** isimleri sistemdekilerle birebir aynı olmalıdır. Eşleşmezse kayıt atlanabilir.
- **Aktif/Pasif**: Personeli geçici olarak pasife alarak nöbet listelerine dahil edilmemesini sağlayabilirsiniz.

## 4. Nöbet Planlama

### 4.1 Otomatik Planlama Sihirbazı
En güçlü özellik burasıdır.
1. Sol menüden **Otomatik Planla** sayfasına gidin.
2. "Başlayalım" butonuna basın.
3. **Parametreler**:
    - **Bölüm**: Planlama yapılacak bölümü seçin.
    - **Ay/Yıl**: Hangi ay için planlama yapacağınızı seçin.
    - **Nöbet Tipleri**: Bu çizelgede kullanılacak vardiyaları işaretleyin.
    - **Personel**: Listeye dahil edilecek personelleri seçin. (İzinli personellerin seçimini kaldırabilirsiniz).
    - **Kurallar**: Ek mesai eşiği ve çarpanlarını belirleyin.
4. "Sihirbazı Çalıştır" butonuna tıklayın. Sistem adil bir dağılım yaparak takvimi oluşturur.

### 4.2 Takvim ve Manuel Düzenleme
Oluşturulan çizelgeyi **Nöbet Takvimi** sayfasından inceleyebilirsiniz.
- **Yeni Nöbet Ekle**: İstediğiniz günün kutusundaki **+** butonuna tıklayın.
- **Düzenle**: Mevcut bir nöbet kutusuna tıklayın. Personeli veya vardiya tipini değiştirebilirsiniz.
- **Sil**: Nöbet detay penceresindeki **Sil** butonunu kullanın.
- **Sürükle-Bırak**: Bir nöbeti tutup başka bir güne sürükleyerek tarihini değiştirebilirsiniz.

### 4.3 Boş Kalan (Atanmamış) Nöbetler
Otomatik planlama sihirbazı, kuralları (zorunlu dinlenme süreleri, yıllık izinler vb.) çiğnememek adına bazı nöbetlere personel atayamayabilir.
- **Görünüm**: Bu nöbetler takvimde **kırmızı kesik çizgili** bir çerçeve ile ve **"BOŞ NÖBET"** uyarısı ile görünür.
- **Çözüm**: Bu kutulara tıklayarak personeli manuel olarak atayabilir veya kuralları esneterek sihirbazı tekrar çalıştırabilirsiniz.

## 5. Puantaj ve Devam Takibi
Gerçekleşen nöbet saatlerini takip etmek için **Puantaj** sayfasını kullanın.
- İlgili Ay ve Bölümü seçin.
- Listede personellerin planlanan nöbetleri görünür.
- **Kaydet/Düzenle**:
    - Personelin gerçekleşen giriş-çıkış saatlerini girin.
    - Eğer nöbet gece başlayıp sabah bitiyorsa (örn. 20:00 - 08:00), sistem bunu otomatik olarak **ertesi gün** olarak algılar ve süreyi (12 saat) doğru hesaplar.
    - Fazla mesai varsa sistem, "Planlanan Süre" ile "Gerçekleşen Süre" arasındaki farkı hesaplayarak **Ek Mesai** olarak kaydeder.

### 5.1 Manuel Puantaj Ekleme (Nöbetsiz Personel İçin)
Planlı nöbeti olmayan veya normal mesai çalışan personeller için "Manuel Ekle" özelliğini kullanabilirsiniz.
1. **Puantaj** sayfasında üstteki **Manuel Ekle** butonuna tıklayın.
2. **Personel**, **Tarih** ve referans alınacak **Vardiya Tipi**ni (örn. Gündüz 08-17) seçin.
3. Giriş ve Çıkış saatlerini girin.
4. **Kaydet** butonuna basın.
Sistem, seçilen vardiya süresini baz alarak (örn. 8 saat) üzerindeki çalışmayı **Ek Mesai** olarak kaydedecektir.

## 6. Raporlar
Sol menüden **Ek Mesai** raporlarına ulaşabilirsiniz.
- **Veri İçeriği**: Seçilen ay için personellerin toplam ek mesai saatlerini (Hafta Tatili, Resmi Tatil, Fazla Mesai) ve bunların **parasal karşılıklarını** görüntüler.
- **Hesaplama**: Ücretler, personelin "Saatlik Ücret"i ve sistemdeki "Mesai Çarpanları" (1.5, 2.0 vb.) kullanılarak otomatik hesaplanır.
- **İndirme**: Listeyi yazdırabilir veya Excel'e aktarabilirsiniz.

---

*Not: Sistemdeki Personel, İzin ve Puantaj listeleri gibi yoğun veri içeren tablolar **sayfalamalı (pagination)** bir yapıya sahiptir. Listenin altındaki navigasyon butonlarını kullanarak diğer sayfalara geçebilirsiniz.*
