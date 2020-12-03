.PHONY: build init plan apply destroy release

VERSION?=$(shell cat version)

build:
	go mod download && go build -o plugins/terraform-provider-etcd_v$(VERSION)

init:
	cd example && terraform init -plugin-dir=$(PWD)/plugins

plan:
	cd example && terraform plan

apply:
	cd example && terraform apply

destroy:
	cd example && terraform destroy

release:
	git add version && \
	git commit -m "release v$(VERSION)" && \
	git tag v$(VERSION) && \
	git push origin v$(VERSION) && \
	git push --follow-tags origin
