# ---- build stage ----
FROM golang:1.24.5-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./main.go

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /app/app /app/app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
