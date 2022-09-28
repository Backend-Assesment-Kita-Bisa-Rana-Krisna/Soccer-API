
# Soccer API
Soccer API that can give a list of soccer players in a team

## Prerequisite
* Golang >= 1.19
* <a style="width:50%" href="https://github.com/gin-gonic/gin" target="_blank">Gin Framework</a>
* <a style="width:50%" href="https://www.mongodb.com/" target="_blank">MongoDB >= 5.0.6</a>
## Installation
Gunakan package manager [go](https://go.dev/) untuk install
```bash
go get
```
Buat file [.env] pada folder untuk konfigurasi Golang Apps
```bash
cp .env.example .env
```
lalu ubah content yg ada di **.env**
Buat database di MongoDB, nama database sesuaikan dengan yang ada di **.env**
## Usage
### Run app
```bash
go run main.go
```
akses service di http://localhost:8000

### Documentation

[https://documenter.getpostman.com/view/4027401/2s83mXN7HS](https://documenter.getpostman.com/view/4027401/2s83mXN7HS)
