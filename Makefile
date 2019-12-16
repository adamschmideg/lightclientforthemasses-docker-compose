test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

docker-deps:
	docker build --target dependencies -t offcode/lightfaucet:dependencies .

docker-build:
	docker build -t offcode/lightfaucet:latest --cache-from offcode/lightfaucet:dependencies .