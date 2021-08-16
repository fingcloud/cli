<h1 align="center">Fing cli</h1>

<p align="center">
  <img width="140" height="140" src="https://fing.ir/images/icon.png" alt="fing icon" />
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/fing-ir/cli"><img src="https://pkg.go.dev/badge/github.com/fing-ir/cli.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/fingcloud/cli"><img src="https://goreportcard.com/badge/github.com/fingcloud/cli" alt="go report card" />
  <a href="https://golang.com"><img src="https://img.shields.io/github/go-mod/go-version/fingcloud/cli?label=version&logo=go" alt="npm package"></a>
  <a href="https://www.npmjs.com/package/@fingcloud/cli"><img src="https://img.shields.io/npm/v/@fingcloud/cli?label=npm&logo=npm" alt="npm package"></a>
  </a>
</p>


- [Fing CLI](#fing-cli)
  - [Installation](#installation)
    - [Homebrew](#homebrew)
    - [NPM](#npm)
    - [CURL](#curl)
      - [Linux nad macOS](#linux-nad-macos)
      - [Windows](#windows)
  - [Quick Start](#quick-start)
  - [Documentation](#documentation)


# Fing CLI
`fing` is a command-line interface for [fing.ir](https://fing.ir?utm_source=github&utm_medium=link&utm_campaign=github_cli&utm_id=ref&utm_content=header_link)

Use this to deploy & manage your applications to Fing Cloud without being worried about your infrastructure and environment.

## Installation

### Homebrew
```bash
brew install fingcloud/tap/fing
```
To Upgrade to the latest version
```bash
brew upgrade fing
```

### NPM
install from npm
```bash
npm install -g @fingcloud/cli
# or
yarn global add @fingcloud/cli
```

### CURL
#### Linux nad macOS
```bash
curl -fsSL https://fing.ir/install.sh | sh
```

#### Windows
```shell
iwr https://fing.ir/install.ps1 -useb | iex
```

## Quick Start
Fing builds your app and deploys it to fing cloud servers.

First you need to login. (if you don't have an account yes, [create an account](https://dashboard.fing.ir/register?utm_source=github&utm_medium=link&utm_campaign=github_cli&utm_id=ref&utm_content=register_link)).
```bash
fing login
````

Go to the root of your project and run the following command:
```bash
cd ~/my-react-app
fing up
```

## Documentation
[View the full documentation](https://docs.fing.ir?utm_source=github&utm_medium=link&utm_campaign=github_cli&utm_id=ref&utm_content=docs_link)
