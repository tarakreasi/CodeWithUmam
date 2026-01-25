# 📚 Panduan Git - Step by Step

## 🎯 Yang Bakal Kita Lakukan

1. Cek git udah install belum
2. Setting email & nama di git
3. Bikin `.gitignore` (file yang gak perlu di-upload)
4. Initialize git repo
5. Add & commit semua file
6. Connect ke remote repository (GitHub/GitLab)
7. Push ke remote

---

## ✅ Step 1: Cek Git Installed

```bash
git --version
```

Harusnya muncul versi git. Kalo error "command not found", install dulu git.

---

## ✅ Step 2: Setting Email & Nama

Ganti `YOUR_EMAIL` dan `YOUR_NAME` dengan info kamu:

```bash
# Set email (wajib!)
git config --global user.email "your.email@example.com"

# Set nama
git config --global user.name "Nama Kamu"
```

**Cek udah ke-set:**
```bash
git config --global --list
```

---

## ✅ Step 3: Bikin .gitignore

File ini buat ngasih tau Git file mana yang **gak usah** di-track (misal: file temporary).

**Isi `.gitignore` untuk project Go:**
```
# Binary hasil compile
codeWithUmam

# Go workspace file (opsional)
go.work
```

---

## ✅ Step 4: Initialize Git Repository

```bash
# Masuk ke folder project
cd /home/twantoro/Project/Golang/codeWithUmam

# Initialize git
git init
```

Output: `Initialized empty Git repository in ...`

---

## ✅ Step 5: Add Semua File

```bash
# Add semua file ke staging
git add .

# Cek status (opsional, buat mastiin)
git status
```

Harusnya muncul file: `go.mod`, `main.go`, `handlers.go`, `models.go`, `.gitignore`

---

## ✅ Step 6: Commit Pertama

```bash
git commit -m "Initial commit: Simple Go CRUD API"
```

Output: `[main (root-commit) xxxxxx] Initial commit: Simple Go CRUD API`

---

## ✅ Step 7: Connect ke Remote Repository

Ganti `YOUR_REPO_URL` dengan URL repo kamu (dari GitHub/GitLab):

```bash
git remote add origin YOUR_REPO_URL
```

**Contoh:**
```bash
git remote add origin https://github.com/username/codeWithUmam.git
```

**Cek remote udah connect:**
```bash
git remote -v
```

---

## ✅ Step 8: Push ke Remote

```bash
# Push ke branch main
git push -u origin main
```

Kalo error **"branch main doesn't exist"**, coba:
```bash
git push -u origin master
```

**Keterangan:**
- `-u` = set upstream (biar next time tinggal `git push` aja)
- `origin` = nama remote (default)
- `main`/`master` = nama branch

---

## 🎉 Selesai!

Project kamu sekarang udah di GitHub/GitLab! 

### Next Time Mau Update

```bash
# Add perubahan
git add .

# Commit dengan pesan
git commit -m "Pesan commit kamu"

# Push
git push
```

---

## 🔧 Troubleshooting

### Error: "please tell me who you are"
→ Belum setting email/nama, balik ke Step 2

### Error: "failed to push"
→ Cek koneksi internet atau URL repo salah

### Error: "permission denied"
→ Perlu setup SSH key atau pakai Personal Access Token di GitHub

---

**Butuh bantuan?** Kasih tau error message nya! 😊
