gen: swagger-gen

.PHONY: gen

swagger-gen:
	#sudo rm -rf ./strava/stravalib
	docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli:latest generate \
	-i https://developers.strava.com/swagger/swagger.json \
	-l go \
	-o /local/strava/stravalib \
	--additional-properties packageName=stravalib
	# Some types are wrong, so I override them
	sudo cp ./strava/custom/* ./strava/stravalib/

.PHONY: swagger-gen

swagger-config:
	sudo rm -rf ./strava/stravalib
	docker run --rm -v ${PWD}:/local swaggerapi/swagger-codegen-cli:latest config-help -l=go

.PHONY: swagger-config