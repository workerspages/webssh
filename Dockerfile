# Stage 1: Build Frontend
FROM node:18-alpine as frontend-builder
WORKDIR /webssh/frontend
COPY frontend/package.json frontend/package-lock.json ./
# 忽略一些可能的 npm 错误，确保依赖安装
RUN apk add --no-cache python3 make g++
ENV NODE_OPTIONS=--openssl-legacy-provider
RUN npm install --legacy-peer-deps
COPY frontend .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine as backend-builder
WORKDIR /webssh
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 将构建好的前端资源覆盖到 public 目录 (main.go 会 embed 这个目录)
COPY --from=frontend-builder /webssh/public ./public
RUN CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o webssh main.go

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app
# 安装基础依赖
RUN apk --no-cache add ca-certificates tzdata
COPY --from=backend-builder /webssh/webssh .
EXPOSE 8888
CMD ["./webssh"]
