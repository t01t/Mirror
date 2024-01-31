v = 020
name = mirror
compile:
	npx tailwindcss -i ./frontend/app.css -o ./frontend/style.css --build --minify
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o bin/${name}.${v}.windows-amd64.exe -ldflags "-H=windowsgui"
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o bin/${name}_${v}_linux-amd64

build:
	npx tailwindcss -i ./frontend/app.css -o ./frontend/style.css --build --minify
	go build