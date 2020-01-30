lightserver:
	geth \
		--lightserv=100 \
		--light.maxpeers=1 \
		--datadir ~/datadirs/goerli/fast \
		--rpc \
		--rpcapi=admin,les,web3,eth \
		--rpcaddr=0.0.0.0 \
		--goerli
web:
	# keys for testing from https://developers.google.com/recaptcha/docs/faq
	go run faucet.go --template ./faucet.html \
		--recaptcha.public recaptcha_v2_test_public.txt \
		--recaptcha.secret recaptcha_v2_test_secret.txt

docker-build:
	docker build -t offcode/lightfaucet:latest --cache-from offcode/lightfaucet:dependencies .