FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum* ./
RUN go mod tidy

COPY backend/ ./
RUN go build -o server ./cmd/KitaF_Q-Tracker/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=frontend-builder /app/frontend/out ./public
COPY --from=backend-builder /app/backend/server ./server

EXPOSE 8080
CMD ["./server"]
