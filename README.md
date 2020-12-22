# Hastebin Client (in Go)

Yet another CLI for hastebin, this one written in Go. It is intended to be used similarly to the original
[haste client](https://github.com/seejohnrun/haste-client), but also be able to get/read hastes.\
As an added bonus, there's also cross-platform support without dependencies on other tools.

It also supports client certificates for mTLS for some reason.

* [Hastebin Client (in Go)](#hastebin-client-in-go)
  * [Installation](#installation)
  * [Usage](#usage)
    * [Creating a haste](#creating-a-haste)
    * [Reading a haste](#reading-a-haste)
    * [Help](#help)
  * [Build](#build)
  * [Author](#author)
  * [License](#license)

## Installation

Navigate to the latest [release](https://github.com/jagoe/haste-client-go/releases) and download the client for your OS and architecture.

## Usage

### Creating a haste

Without any arguments, `haste` reads from the STDIN:

```bash
cat file | haste
```

Provided with a file argument, `haste` will create a haste from the contents of that file and return the URL pointing to it:

```bash
url=$(haste file) # equivalent to cat file | haste
echo $url         # e.g. https://hastebin.com/ogoquyocaq
```

### Reading a haste

`haste` can read a haste using the `get` command:

```bash
haste get <key>           # prints the haste contents to STDOUT
haste get <key> -o ./file # prints the haste contents to ./file
```

### Help

For more detailed information on how `haste` can be used, use `haste --help` or look here:

```plaintext
A hastebin client that can create hastes from files and STDIN and read hastes from a configurable server.

Usage:
  haste [file] [flags]
  haste [command]

Examples:
echo Test | haste
cat ./file | haste
haste ./file

Available Commands:
  get         Get a haste from the server
  help        Help about any command

Flags:
      --client-cert string       Client certificate path
      --client-cert-key string   Client certificate key path
  -c, --config string            Config file [$HOME/.haste-client-go.yaml]
  -h, --help                     help for haste
  -s, --server string            Server URL (default "https://hastebin.com")

Use "haste [command] --help" for more information about a command.
```

## Build

_Requires [`golang 1.15+`](https://golang.org/doc/install)._

If there is no release for your OS and architecture or if you want to build the binary yourself for other reasons,
you can follow these steps:

```bash
git clone git@github.com:jagoe/haste-client-go.git        # clone this repo
cd haste-client-go
[GOARCH=<your architecture>] [GOOS=<your OS>] make build  # build client
sudo mv ./bin/haste /usr/local/bin/haste                  # move to a bin directory in your PATH
```

## Author

Jakob Göbel (goebel.jakob@gmail.com)

## License

[Apache License 2.0](./LICENSE)
