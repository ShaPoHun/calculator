build:
	go build -o testApp

mac:
	fyne package -os darwin -icon Icon.png

windows:
	fyne package -os windows -icon Icon.png

linux:
	fyne package -os linux -icon Icon.png

help:
	@echo "Available targets:"
	@echo "  build - Build the application"
	@echo "  help  - Display this help message"
	@echo "  mac   - Build the application for macOS"
	@echo "  windows - Build the application for Windows"
	@echo "  linux - Build the application for Linux"