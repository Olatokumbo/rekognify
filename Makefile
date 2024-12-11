deploy:
	sam build
	sam deploy

watch:
	sam --image-repository 626635403708.dkr.ecr.eu-west-2.amazonaws.com/rekognifye82d090a/imageuploadhandlera77b1b71rep rekognify --watch