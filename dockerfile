# 使用官方 Go 映像作為基礎，選擇輕量化的 alpine 版本
FROM golang:1.24.5 AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum，先下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製專案程式碼
COPY . .

# 編譯 Go 應用程式，關閉 CGO 並生成靜態二進位檔案
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用輕量化的 alpine 映像作為最終映像
FROM alpine:latest

# 安裝必要的證書（如果需要 HTTPS 或外部連線）
RUN apk --no-cache add ca-certificates

# 設置工作目錄
WORKDIR /app

# 從 builder 階段複製編譯好的二進位檔案
COPY --from=builder /app/main .

# 暴露應用程式埠（假設 Gin 使用 8080）
EXPOSE 8080

# 啟動應用程式
CMD ["./main"]