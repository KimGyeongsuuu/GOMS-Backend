FROM golang:1.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goms-backend-go ./app

FROM alpine:latest

COPY --from=builder /build/goms-backend-go /goms-backend-go
COPY --from=builder /build/.env /

EXPOSE 8080

CMD ["/goms-backend-go"]

RUN ls -l /goms-backend-go

