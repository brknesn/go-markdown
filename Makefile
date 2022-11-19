BINARY_NAME=MarkDownEditor.app
APP_NAME=MarkDownEditor
VERSION=1.0.0

build:
	rm -rf ${BINARY_NAME}
	rm -f fyne-md
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -release

run:
	go run .

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ${BINARY_NAME}
	@echo "Cleaned!"

test:
	go test -v ./...