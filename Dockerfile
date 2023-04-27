FROM node:16.19.1-alpine

RUN apk --no-cache add make

WORKDIR /strava
COPY . ./

RUN npm install react-scripts@3.4.1 -g --silent
RUN make site-install
RUN make site

FROM golang:1.20.3-alpine3.17

RUN apk --no-cache add make

COPY --from=0 /strava /strava

WORKDIR /strava
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /strava/site/strava-frontend/build /strava/site/strava-frontend/build
COPY --from=1 /strava/bin/strava ./
RUN chmod +x ./strava
ENTRYPOINT ["./strava", "server"]
