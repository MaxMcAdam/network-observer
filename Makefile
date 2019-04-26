
-include horizon/hzn.cfg

build:
	docker build -t $(DOCKER_USER)/$(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) .

dev:
	docker run --net=host -it $(DOCKER_USER)/$(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) /bin/sh

run:
	docker run -d --env-file servenv --net=host $(DOCKER_USER)/$(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) --name $(SERVICE_NAME)

publish-service:
	hzn exchange service publish -f horizon/service.definition.json -k $(HZN_PRIVATE_KEY_FILE) -K $(HZN_PUBLIC_KEY_FILE)

publish-pattern:
	: $${HZN_ORG_ID:?} $${HZN_EXCHANGE_USER_AUTH:?}
	hzn exchange pattern publish -p network-observer -f horizon/pattern/network-observer.json

clean:
	-docker rmi $(DOCKER_USER)$(DOCKER_IMAGE_BASE)_$(ARCH):$(SERVICE_VERSION) 2> /dev/null || :

.PHONY: build clean
