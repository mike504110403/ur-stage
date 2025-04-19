FROM golang:alpine AS builder

RUN apk add --no-cache git
## RUN GOCACHE=OFF
RUN go env -w GOPRIVATE=gitlab.com

# modules token
ARG MODULES_TOKEN
RUN git config --global url."https://:${MODULES_TOKEN}@gitlab.com".insteadOf "https://gitlab.com"


# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 文件並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製其餘的應用程序文件
COPY . .

# 構建應用程序
RUN go build -o main .

# 使用輕量級的Alpine作為運行鏡像
FROM alpine:3.18

# 安裝 tzdata 以設置時區
RUN apk add --no-cache tzdata

# 設置時區為 Asia/Shanghai（UTC+8）
ENV TZ=Asia/Shanghai

WORKDIR /

# 複製編譯好的應用程序到運行鏡像
COPY --from=builder /app/main .

# 確保可執行文件具有執行權限
RUN chmod +x ./main
# RUN ls -la /

# 設置容器啟動命令
CMD ["./main"]