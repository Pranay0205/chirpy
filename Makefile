# Build and run the application
run:
	go build -o out && ./out

# Just build without running
build:
	go build -o out

# Clean up built binaries
clean:
	rm -f out

# Run with live reload during development
dev:
	go run .

# Build for production (with optimizations)
build-prod:
	go build -ldflags="-w -s" -o out

.PHONY: run build clean dev build-prod