FROM golang:1.20.3-alpine3.17

RUN apk --no-cache add make
WORKDIR /strava
COPY . ./
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /strava/bin/strava ./
RUN chmod +x ./strava
ENTRYPOINT ["./strava", "server"]
