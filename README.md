# GMMIT

This a small terminal tool to help developers or any person interacting with Git to write their commit messages.

---

You'd certainly agree that inspecting a properly maintained git history log that follows the [Conventional Commits](https://www.conventionalcommits.org) standard is nice and simple. And you'd also agree that having to dedicate time to come up with those neat and descriptive messages is a bummer. Well, here's where **gmmit** comes to saves the day.

Run the `gmmit` command on a folder with a git repository, and it will check your staged changes, the current branch name, and generate a tidy commit message for you, and even run the "git commit" command. After that, just execute `git push`, and it's done.

![gmmit command example](images/commit.gif)

## Getting started

### Dependencies

| Dependency | Version |
| ---------- | ------- |
| Golagn     | 1.22.*  |

### Build

1. You need to have Golang installed, or use a docker container with it.

```bash
docker run --rm --name gmmit-dev -v $(pwd):/go -it golang:1.22-alpine sh
```

2. To test it you'll need `git` as well. If you're using the docker container, install it with the following command:

```bash
apk add --update git
```

3. Run you code:

```bash
go run ./cmd/gmmit/
```

4. Optionally, run the build command

```bash
go build -o build/gmmit ./cmd/gmmit/
```

**NOTE**: If working on MacOS, you may need to build the arm version in order to test it on you local.

```bash
env GOOS=darwin GOARCH=arm64 go build -o build/gmmit ./cmd/gmmit/
```

### Run

Gmmit uses `Gemmini AI` models to generate the messages. For it to work you need to provide an `API KEY` first.

1. Get a Gemini API Key. You'll need a Google Account (Gmail) to do this, then follow the steps described [here](https://geminiforwork.gwaddons.com/setup-api-keys/create-geminiai-api-key).

2. Setup your env:

```bash
echo GMMIT_API_KEY="<API_KEY>" > ~/.gmenv
```

3. Download the binary for you OS from the release section, and install it with your apps.

    If you're using Linux or MacOS you wanna move it to `/usr/local/bin`, or any other folder on your `PATH`.

4. Now move to a folder with a git repository, add some files to the staging area, and run the command.

```bash
gmmit
```

4. Profit.

## Usage

| option | description |
| ------ | ----------- |
| --no-verify | Add the option to skip git hooks and commit even if the pre-commit hooks fail. This is helpful when you are sure that your changes are correct and you don't want to wait for the hooks to run. |

## Troubleshooting

If you're having issues running **gmmit**, try executing it in *debug* mode, and checking the following common issues before submiting a new issue.

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

## Contributing

Contributors are more than welcome! Here's how you can propose and submit changes to the project.

(These steps are based on [First Contributions](https://github.com/firstcontributions/first-contributions/blob/main/README.md) documents.)

### Fork this repository

Fork this repository by clicking on the fork button on the top of this page.
This will create a copy of this repository in your account.

### Clone the repository

<img align="right" width="300" src="https://firstcontributions.github.io/assets/Readme/clone.png" alt="clone this repository" />

Now clone the forked repository to your machine. Go to your GitHub account, open the forked repository, click on the code button and then click the _copy to clipboard_ icon.

Open a terminal and run the following git command:

```bash
git clone "url you just copied"
```

where "url you just copied" (without the quotation marks) is the url to this repository (your fork of this project). See the previous steps to obtain the url.

<img align="right" width="300" src="https://firstcontributions.github.io/assets/Readme/copy-to-clipboard.png" alt="copy URL to clipboard" />

For example:

```bash
git clone git@github.com:this-is-you/first-contributions.git
```

where `this-is-you` is your GitHub username. Here you're copying the contents of the first-contributions repository on GitHub to your computer.

### Create a branch

Change to the repository directory on your computer (if you are not already there):

```bash
cd first-contributions
```

Now create a branch using the `git switch` command:

```bash
git switch -c your-new-branch-name
```

For example:

```bash
git switch -c add-alonzo-church
```

### Make necessary changes and commit those changes

Now open `Contributors.md` file in a text editor, add your name to it. Don't add it at the beginning or end of the file. Put it anywhere in between. Now, save the file.

<img align="right" width="450" src="https://firstcontributions.github.io/assets/Readme/git-status.png" alt="git status" />

If you go to the project directory and execute the command `git status`, you'll see there are changes.

Add those changes to the branch you just created using the `git add` command:

```bash
git add Contributors.md
```

Of course every commit in this repo should follow the [Conventional Commits](https://www.conventionalcommits.org) standard (c'mon, just use the tool you're building).

```bash
gmmit
```

or

```bash
git commit -m "<type>[optional scope]: Add your-name to Contributors list"
```

replacing `your-name` with your name.

### Push changes to GitHub

Push your changes using the command `git push`:

```bash
git push -u origin your-branch-name
```

replacing `your-branch-name` with the name of the branch you created earlier.

<details>
<summary> <strong>If you get any errors while pushing, click here:</strong> </summary>

- ### Authentication Error
     <pre>remote: Support for password authentication was removed on August 13, 2021. Please use a personal access token instead.
  remote: Please see https://github.blog/2020-12-15-token-authentication-requirements-for-git-operations/ for more information.
  fatal: Authentication failed for 'https://github.com/<your-username>/first-contributions.git/'</pre>
  Go to [GitHub's tutorial](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account) on generating and configuring an SSH key to your account.

</details>

### Submit your changes for review

If you go to your repository on GitHub, you'll see a `Compare & pull request` button. Click on that button.

<img style="float: right;" src="https://firstcontributions.github.io/assets/Readme/compare-and-pull.png" alt="create a pull request" />

Now submit the pull request.

<img style="float: right;" src="https://firstcontributions.github.io/assets/Readme/submit-pull-request.png" alt="submit pull request" />

Soon I'll be merging all your changes into the main branch of this project. You will get a notification email once the changes have been merged.

### Where to go from here?

Congrats! You just completed the standard _fork -> clone -> edit -> pull request_ workflow!

### Repository Internal Structure 

This repository adopts the Golang Project structure as described by Golang-Standards [here](https://github.com/golang-standards/project-layout?tab=readme-ov-file#go-directories).

## License

This project is under [Apache License v2](LICENSE.md).