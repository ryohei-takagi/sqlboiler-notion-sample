FROM golang:1.21.1-alpine

WORKDIR /workspace

RUN go install github.com/volatiletech/sqlboiler/v4@v4.15.0 && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.15.0
