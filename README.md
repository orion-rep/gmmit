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

2. Setup your env:

```bash
echo GMMIT_API_KEY="<API_KEY>" > ~/.gmenv
```

Here's how: https://geminiforwork.gwaddons.com/setup-api-keys/set-up-geminiai-api-key

3. Run the binary on a folder inside a git repository:

```bash
gmmit
```

i.e.:

![gmmit command example](images/commit.gif)

4. Profit

## Troubleshooting

### Message Blocked
```
<date-time> blocked: candidate: FinishReasonSafety
```

This error happen when the AI Model detects the content being sent is potencially dangerous. Review the content of the changes being sent and try again.