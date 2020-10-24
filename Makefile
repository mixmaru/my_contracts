docker-build:
	cd ./development && \
		docker-compose up -d --build
docker-clean:
	cd ./development && \
		docker-compose stop && docker-compose rm -f
docker-rebuild: docker-clean docker-build
test:
	cd ./development && docker-compose exec dev go test ./... -count=1 -cover -p 1
