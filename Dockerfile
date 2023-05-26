FROM golang:1.20-alpine as build-stage

RUN apk --no-cache add ca-certificates
WORKDIR /go/delivery/bitcoin-rate-app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -buildvcs=false -o /bitcoin-rate-app ./cmd/bitcoin-rate-app

#
# final build stage
#
FROM scratch

# Copy ca-certs for app web access
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /bitcoin-rate-app /bitcoin-rate-app

# app uses port 3333
EXPOSE 3333

ENTRYPOINT ["/bitcoin-rate-app"]
