# Defer Panic CLI

Toolkit for managing Defer Panic articles and other stuff.

## Requirements

* [pandoc](https://pandoc.org)
* xelatex (for exporting to PDF)
* [Computer Modern Unicode fonts](https://cm-unicode.sourceforge.io/index.html) (if you want to use them for exporting to PDF) 

## Install

Go to Releases page and download archive with latest version for your platform.

If you want to build Defer Panic CLI manually, make sure you have Go 1.19+ on your machine, then run:

```shell
go install github.com/defer-panic/dp-cli@latest
```

## Usage

You can get help with `--help` flag or just by running `dp-cli` without any arguments:

```
Toolkit for managing Defer Panic articles and other stuff

Usage:
  dp-cli [command]

Available Commands:
  article     Manage articles
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  login       Login to your account on https://dfrp.cc (default) or other instance of dfrp-like infrastructure
  url         Shorten given URL
  version     Get dp-cli version

Flags:
  -h, --help   help for dp-cli

Use "dp-cli [command] --help" for more information about a command.
```

