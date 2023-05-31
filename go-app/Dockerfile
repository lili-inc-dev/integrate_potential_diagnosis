FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY services/lili/etc /app/etc
COPY services/lili/static /app/static
RUN go build -ldflags="-s -w" -o /app/lili services/lili/lili.go


FROM golang:1.18-alpine

ENV TZ Asia/Tokyo

WORKDIR /app
COPY --from=builder /app/lili /app/lili
COPY --from=builder /app/etc /app/etc
COPY --from=builder /app/static /app/static
COPY --from=builder /build/services/lili/database /app/database
COPY --from=builder /build/services/lili/service-account-file.json /app/service-account-file.json

EXPOSE 8080

CMD ["./lili", "-f", "etc/lili-api.yaml"]
