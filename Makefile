IMAGE=garugaru/breeze
publish:
	docker build --build-arg ARCH=amd64 -t $(IMAGE):amd64-latest .
	docker build --build-arg ARCH=arm -t $(IMAGE):arm-latest .
	docker push $(IMAGE):amd64-latest
	docker push $(IMAGE):arm-latest
	#docker manifest create $(IMAGE) $(IMAGE):amd64-latest $(IMAGE):arm-latest
	docker manifest annotate --arch amd64 $(IMAGE):latest $(IMAGE):amd64-latest
	docker manifest annotate --arch arm $(IMAGE):latest $(IMAGE):arm-latest
	docker push garugaru/breeze
	docker manifest push garugaru/breeze

