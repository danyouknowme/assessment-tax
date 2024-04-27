FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o out/assessment-tax main.go

FROM alpine:3.19.1
COPY --from=builder /app/out/assessment-tax /app/assessment-tax

EXPOSE 8080

CMD ["/app/assessment-tax"]