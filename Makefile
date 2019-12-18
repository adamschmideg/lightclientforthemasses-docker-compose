test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# I wish this worked:
# docker-compose -f docker-compose.test.yml run lightserver
# but the "bridge" network and port forwarding may not work together
lightserver:
	geth \
		--lightserv=100 \
		--datadir ./datadir \
		--rpc \
		--rpcapi=admin,les,web3,eth \
		--rpcaddr=0.0.0.0 \
		--goerli
web:
	# keys for testing from https://developers.google.com/recaptcha/docs/faq
	go run faucet.go --template ./faucet.html \
		--recaptcha.public '6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI' \
		--recaptcha.secret '6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe'

# Run this first: docker-compose -f docker-compose.test.yml up 
retest:
	./retest.sh
# This one takes time, you have to do it only once
docker-deps:
	docker build --target dependencies -t offcode/lightfaucet:dependencies .

docker-build:
	docker build -t offcode/lightfaucet:latest --cache-from offcode/lightfaucet:dependencies .