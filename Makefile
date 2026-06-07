.PHONY: build test testacc vet fmt install clean

build:
	go build -o terraform-provider-workos .

test:
	go test -v -cover ./internal/...

testacc:
	TF_ACC=1 go test -v -cover -count=1 ./internal/provider/

vet:
	go vet ./...

fmt:
	go fmt ./...

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/workos/workos/1.0.0/linux_amd64/
	cp terraform-provider-workos ~/.terraform.d/plugins/registry.terraform.io/workos/workos/1.0.0/linux_amd64/

clean:
	rm -f terraform-provider-workos
