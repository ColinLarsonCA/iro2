.PHONY: web pb

be:
	docker compose up --build -d

web:
	cd web; npm run dev

prettier:
	npx prettier . --write

pb:
	rm -rf api/*
	rm -rf backend/pb/*
	cd protos; buf generate