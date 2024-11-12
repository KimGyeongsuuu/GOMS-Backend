FROM golang:1.22 AS builder

WORKDIR /GOMS-Backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goms-backend-go

FROM alpine:latest

COPY --from=builder /GOMS-Backend/resource/app.yml /resource/app.yml

EXPOSE 8080

CMD ["/goms-backend-go"]