install-ui:
	cd ./web/ui && npm i

dev-ui: install-ui
	cd ./web/ui && npm run build -- -w

dev-web:
	reflex -d none -R web/ui. -r \.go -s -- go run ./cmd/brankasd/main.go
