# Hastebin Client (in Go)

This small project aims to provide a CLI for hastebin for any architecture. It is intended to be used similarly to the
original [haste client](https://github.com/seejohnrun/haste-client), but also be able to get hastes and have no dependency on
`gem`.

This client also supports client certificates for some reason.

* [Hastebin Client (in Go)](#hastebin-client-in-go)
  * [Usage](#usage)
    * [Creating a haste](#creating-a-haste)
    * [Reading a haste](#reading-a-haste)
    * [Help](#help)
  * [Installation](#installation)
  * [Build](#build)
  * [Author](#author)
  * [License](#license)

## Usage

### Creating a haste

Without any arguments, `haste` reads from the STDIN:

```bash
cat file | haste
```

Provided with a file argument, `haste` will create a haste from the contents of that file:

```bash
haste -f file # equivalent to cat file | haste
```

#### Output

When creating hastes, the following outputs will be printed to STDOUT:

```bash
haste -f file # URL of the haste, e.g. http://hastebin.com/ahuwabaqij
haste -f file -o key # Key of the haste, e.g. ahuwabaqij
haste -f file -o json # Key of the haste in JSON format, e.g. {"key": "ahuwabaqij"}
```

### Reading a haste

`haste` can read a haste with the `--get` flag:

```bash
haste get <key> # prints the haste contents to STDOUT
haste get <key> -f ./file # prints the haste contents to ./file
```

### Help

For more detailed information on how `haste` can be used, use `haste --help` or look here:

<!-- BEGIN:HELP OUTPUT -->
<!-- TODO: generate help and add here -->
<!-- END:HELP OUTPUT -->

## Installation

Navigate to the latest [release]() and download the client for your OS and architecture.
<!-- TODO: add link to releases -->

## Build

_Requires `golang 1.15+`._

If there is no release for your OS and architecture or if you want to build the binary yourself for other reasons,
you can follow these steps:

```bash
git clone git@github.com:jagoe/haste-client-go.git # clone this repo
cd haste-client-go
[GOARCH=<your architecture>] [GOOS=<your OS>] make build # build client
sudo mv ./bin/haste /usr/local/bin/haste
```

## Author

Jakob Göbel (goebel.jakob@gmail.com)

## License

[MIT License](./LICENSE.MD)
