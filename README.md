[![Go Report Card](https://goreportcard.com/badge/github.com/zuuby/zuuby-ipfs)](https://goreportcard.com/report/github.com/zuuby/zuuby-ipfs)
# zuuby-ipfs
Persistent storage implementation using IPFS

# Quickstart
First clone the repository and build the Docker image.

```bash
$ git clone https://github.com/zuuby/zuuby-ipfs.git
$ cd zuuby-ipfs
$ docker build -t zuuby/zuuby-ipfs:1.0.0 .
```

A sample Docker run command is provided as an bash script

```bash
$ ./cmd/cli/docker_run
```

Now you can `curl` the api endpoints at `localhost:5000/api/file?hash=<file-hash>`
