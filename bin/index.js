// Copyright © 2023 fing Corp
// SPDX-License-Identifier: Apache-2.0

const binwrap = require('binwrap')
const path = require('path')

const pkg = require(path.join(__dirname, '..', 'package.json'))
const binary = 'fing'
const root = `https://github.com/fingcloud/cli/releases/download/v${pkg.version}/${binary}-`

module.exports = binwrap({
  dirname: __dirname,
  binaries: [binary],
  urls: {
    'darwin-x64': root + 'darwin-amd64.tar.gz',
    'darwin-arm64': root + 'darwin-arm64.tar.gz',
    'linux-x64': root + 'linux-amd64.tar.gz',
    'linux-arm64': root + 'linux-arm64.tar.gz',
    'win32-x32': root + 'windows-386.tar.gz',
    'win32-x64': root + 'windows-amd64.tar.gz',
    'win32-arm64': root + 'windows-arm64.tar.gz',
  },
})