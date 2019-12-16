test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Run this first: docker-compose -f docker-compose.test.yml up 
retest:
	./retest.sh
# This one takes time, you have to do it only once
docker-deps:
	docker build --target dependencies -t offcode/lightfaucet:dependencies .

docker-build:
	docker build -t offcode/lightfaucet:latest --cache-from offcode/lightfaucet:dependencies .