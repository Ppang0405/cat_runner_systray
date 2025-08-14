.PHONY: build clean run dmg all

# Application name
APP_NAME = CatRunner
APP_BUNDLE = $(APP_NAME).app
DMG_FILE = $(APP_NAME).dmg

# Go build flags
GO_FLAGS = -ldflags="-s -w"

all: build

# Build the Go binary
build:
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=1 GOOS=darwin go build $(GO_FLAGS) -o $(APP_NAME) .
	@echo "Build complete: $(APP_NAME)"

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@./$(APP_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(APP_NAME) $(APP_BUNDLE) $(DMG_FILE)
	@echo "Clean complete"

# Create macOS application bundle
bundle: build
	@echo "Creating application bundle..."
	@mkdir -p $(APP_BUNDLE)/Contents/{MacOS,Resources}
	@cp $(APP_NAME) $(APP_BUNDLE)/Contents/MacOS/
	@cp assets/*.png $(APP_BUNDLE)/Contents/Resources/
	@ln -sf ../Resources/0.png $(APP_BUNDLE)/Contents/Resources/AppIcon.png
	@cat > $(APP_BUNDLE)/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$(APP_NAME)</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.catrunner</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>$(APP_NAME)</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF
	@echo "Application bundle created: $(APP_BUNDLE)"

# Create DMG file
dmg: bundle
	@echo "Creating DMG file..."
	@hdiutil create -volname "$(APP_NAME)" -srcfolder $(APP_BUNDLE) -ov -format UDZO $(DMG_FILE)
	@echo "DMG file created: $(DMG_FILE)"

# Help command
help:
	@echo "Available commands:"
	@echo "  make build     - Build the Go binary"
	@echo "  make run       - Build and run the application"
	@echo "  make bundle    - Create a macOS application bundle"
	@echo "  make dmg       - Create a DMG file for distribution"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make all       - Build the application (default)"
	@echo "  make help      - Show this help message"
