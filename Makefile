gen: swagger-gen

.PHONY: gen

swagger-gen:
	sudo rm -rf ./strava/stravalib
	docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate \
	-i https://developers.strava.com/swagger/swagger.json \
	-l go \
	-o /local/strava/stravalib \
	--additional-properties packageName=stravalib

.PHONY: swagger-gen