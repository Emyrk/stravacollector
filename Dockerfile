FROM golang:1.20.3

WORKDIR /strava
COPY . ./
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /strava/bin/strava ./
RUN chmod +x ./strava
ENTRYPOINT ["./strava", "server"]
