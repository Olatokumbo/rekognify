deploy:
	sam build
	sam deploy

dev:
	sam local start-api --warm-containers eager

sync:
	sam sync