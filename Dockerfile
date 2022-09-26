FROM golang:1.18-alpine as builder
WORKDIR /build/
COPY . .
RUN apk add --no-cache make
RUN make build

FROM alpine:latest
LABEL author=Nathan13888

WORKDIR /app
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/America/Toronto /etc/localtime

COPY --from=builder /build/bin/ipstat /app/ipstat
RUN chmod +x /app/ipstat

EXPOSE 3000

ENTRYPOINT [ "/app/ipstat" ]


