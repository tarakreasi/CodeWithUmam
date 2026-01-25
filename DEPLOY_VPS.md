# 🚀 Deploy Go API ke VPS Debian - Step by Step

Panduan lengkap deploy project Go CRUD API ke VPS Debian dari GitHub.

---

## 📋 Persiapan

Yang kamu butuhkan:
- ✅ VPS Debian (sudah running)
- ✅ SSH access ke VPS
- ✅ Project sudah di GitHub

---

## 🔧 Step 1: Login ke VPS

```bash
ssh user@ip-vps-kamu
```

Ganti:
- `user` = username VPS kamu
- `ip-vps-kamu` = IP address VPS

**Contoh:**
```bash
ssh root@192.168.1.100
```

---

## 🔧 Step 2: Update System (Opsional tapi Recommended)

```bash
sudo apt update && sudo apt upgrade -y
```

Tunggu sampai selesai.

---

## 🔧 Step 3: Install Git

Cek git udah installed belum:
```bash
git --version
```

Kalo error/belum ada, install:
```bash
sudo apt install git -y
```

---

## 🔧 Step 4: Install Go (Golang)

### Cek Go udah ada belum:
```bash
go version
```

### Kalo belum ada, install Go:

```bash
# Download Go (versi 1.21.6 - sesuaikan dengan versi terbaru)
cd ~
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# Extract ke /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# Set PATH di .bashrc atau .profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc

# Reload .bashrc
source ~/.bashrc

# Cek instalasi
go version
```

Output harusnya: `go version go1.21.6 linux/amd64`

---

## 🔧 Step 5: Clone Project dari GitHub

```bash
# Bikin folder untuk project (opsional)
mkdir -p ~/projects
cd ~/projects

# Clone repository
git clone https://github.com/tarakreasi/CodeWithUmam.git

# Masuk ke folder project
cd CodeWithUmam
```

---

## 🔧 Step 6: Cek Files & Dependencies

```bash
# Lihat isi folder
ls -la

# Cek go.mod (harusnya cuma 3 baris, no deps!)
cat go.mod
```

Karena project ini **pure standard library**, gak perlu install dependencies tambahan!

---

## 🔧 Step 7: Build Aplikasi (Opsional)

Ada 2 cara jalanin:

### Cara 1: Langsung Run (Development)
```bash
go run .
```

### Cara 2: Build Binary (Production - Recommended)
```bash
# Build binary
go build -o codeWithUmam .

# Jalankan binary
./codeWithUmam
```

**Output:**
```
Server running on http://localhost:8080
```

---

## 🔧 Step 8: Test dari VPS

Buka terminal/SSH baru, test API:

```bash
# Test GET
curl http://localhost:8080/categories

# Test POST
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "description": "Dari VPS"}'
```

---

## 🔧 Step 9: Akses dari Luar (Public Access)

Kalo mau diakses dari internet, ada beberapa opsi:

### Opsi A: Ganti Port Binding di Code

Edit `main.go`:
```go
// Dari:
log.Fatal(http.ListenAndServe(":8080", nil))

// Jadi (listen di semua interface):
log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
```

Commit & push perubahan, terus pull di VPS:
```bash
git pull origin main
go build -o codeWithUmam .
./codeWithUmam
```

### Opsi B: Pakai Reverse Proxy (Nginx) - Recommended!

Install Nginx:
```bash
sudo apt install nginx -y
```

**Akan dibahas di guide terpisah jika dibutuhkan.**

---

## 🔧 Step 10: Jalankan di Background (Production)

Agar server tetap jalan pas SSH ditutup:

### Cara 1: Pakai `nohup`
```bash
nohup ./codeWithUmam > app.log 2>&1 &
```

### Cara 2: Pakai `screen` (Recommended)
```bash
# Install screen
sudo apt install screen -y

# Bikin session baru
screen -S goapi

# Jalankan app
./codeWithUmam

# Detach: tekan Ctrl+A, terus D
# Attach lagi: screen -r goapi
```

### Cara 3: Pakai Systemd Service (Paling Pro)
```bash
# Bikin service file
sudo nano /etc/systemd/system/codewithumam.service
```

Isi file:
```ini
[Unit]
Description=CodeWithUmam Go API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/projects/CodeWithUmam
ExecStart=/root/projects/CodeWithUmam/codeWithUmam
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Aktifkan service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable codewithumam
sudo systemctl start codewithumam

# Cek status
sudo systemctl status codewithumam
```

---

## 📊 Monitoring & Troubleshooting

### Cek app jalan:
```bash
ps aux | grep codeWithUmam
```

### Cek port 8080 listen:
```bash
sudo netstat -tulpn | grep 8080
# atau
sudo ss -tulpn | grep 8080
```

### Lihat log (kalo pakai systemd):
```bash
sudo journalctl -u codewithumam -f
```

---

## 🔥 Firewall (Jika Perlu)

Kalo VPS pakai firewall (UFW):
```bash
# Cek status
sudo ufw status

# Allow port 8080
sudo ufw allow 8080/tcp

# Reload
sudo ufw reload
```

---

## 🎯 Testing dari Luar

Dari komputer lokal (bukan VPS):
```bash
# Ganti IP_VPS dengan IP VPS kamu
curl http://IP_VPS:8080/categories
```

---

## 🚀 Next Steps (Opsional)

1. **Setup Domain** → Pakai Cloudflare/Namecheap
2. **SSL Certificate** → Let's Encrypt (gratis)
3. **Nginx Reverse Proxy** → Lebih aman & flexible
4. **Auto Deploy** → GitHub Actions + SSH

---

**Selamat! API kamu sekarang live di VPS!** 🎉

Butuh bantuan lebih lanjut? Tanya aja! 😊
