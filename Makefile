V?=0
GRACE_PERIOD?=0

deploy:
	mkdir -p ~/.kube/plugins/kill
	go build -o ~/.kube/plugins/kill/kill
	cp plugin.yaml ~/.kube/plugins/kill/
	@echo

create-pod:
	kubectl create -f test.yaml
	@echo
	kubectl get pod hello-world -o jsonpath='{.metadata.finalizers}'
	@echo
	@echo

kill-pod:
	kubectl plugin kill hello-world -v=${V} --grace-period=${GRACE_PERIOD}

test: create-pod kill-pod
