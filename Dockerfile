# Gunakan gambar dasar Golang
FROM golang:1.22-alpine

# Setel variabel lingkungan untuk root direktori
ENV APP_HOME /app
WORKDIR $APP_HOME

# Salin go.mod dan go.sum ke direktori kerja
COPY go.mod go.sum ./

# Unduh dependencies Go
RUN go mod download

# Salin seluruh kode sumber ke dalam container
COPY . .

# Kompilasi aplikasi Go
RUN go build -ldflags "-s -w" -o main .

# Eksekusi aplikasi Go
#CMD ["./main"]
ENTRYPOINT ["./main"]
