be:
	docker compose up --build -d

web:
	npm run web

ios:
	npm run ios

android:
	npm run android

protogen:
	rm -rf backend/pb/*
	mkdir -p backend/pb
	protoc --go_out=backend/pb --go-grpc_out=backend/pb --proto_path=protos protos/*.proto
