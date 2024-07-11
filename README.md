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

### Usage

| option | description |
| ------ | ----------- |
| --no-verify | Add the option to skip git hooks and commit even if the pre-commit hooks fail. This is helpful when you are sure that your changes are correct and you don't want to wait for the hooks to run. |

## Troubleshooting

If you're having issues running **gmmit**, try executing it in *debug*  mode, and checking the following common issues before submiting a new issue.

### Debug Mode

Run the command as follow to get debug output logs:

```bash
GMMIT_DEBUG=true gmmit
```

### FinishReasonSafety - Message Blocked
```
<date-time> blocked: candidate: FinishReasonSafety
```

This error happen when the AI Model detects the content being sent is potencially dangerous. Review the content of the changes being sent and try again.

### Error 429 - Quota Exceeded

```
<date-time> googleapi: Error 429:
```

Error 429 generally indicates that you have exceeded the limit of allowed requests in a specific time period when interacting with a Google API. 

This can occur when making many requests to an API in a short period of time.