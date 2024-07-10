# GMMIT



https://github.com/golang-standards/project-layout/tree/master/internal

## Getting started

### Build
Build it on docker

```bash
docker run --rm --name gmmit-dev -v $(pwd):/go -it golang:1.22-alpine sh
```

```bash
apk add --update git
env GOOS=darwin GOARCH=arm64 go build -o build/gmmit ./cmd/gmmit/
```

### Run

1. Get a Gemini API Key: 

2. Export it on your env:

```bash
export GMMIT_API_KEY="<API_KEY>"
```

Here's how: https://geminiforwork.gwaddons.com/setup-api-keys/set-up-geminiai-api-key

3. Run the binary on a folder inside a git repository:

```bash
./build/gmmit
```

4. Profit