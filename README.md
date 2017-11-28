# dj

## Install dj

Install for Mac:

```sh
docker run --rm -v ${HOME}/tmp/dj:/dj/bin devo/dj install --bin --name dj --to /dj/bin https://github.com/devo/dj/releases/download/{LATEST}/dj_mac && mv ${HOME}/tmp/dj/dj /usr/local/bin/dj
```

Install for Linux:

```sh
docker run --rm -v ${HOME}/tmp/dj:/dj/bin devo/dj install --bin --name dj --to /dj/bin https://github.com/devo/dj/releases/download/{LATEST}/dj_linux && mv ${HOME}/tmp/dj/dj /usr/local/bin/dj
```

For other OS's:

TODO

## Usage

Install Docker based tools with `dj install`:

```sh
dj install devo/ruby --name ruby2 --force
ruby2 -v
```

Install any ol' binary tool:

```sh
dj install --bin https://dl.google.com/gactions/updates/bin/darwin/amd64/gactions/gactions
```

## Contributing

```sh
make install
```

## TODO

* Make script to build and deploy the dockers, rather than a build.sh in each of them. Use directory name for name.
* Add installer script, could we use docker for it?  mnt /usr/local/bin/fn and throw it there (make generic installer, `docker run --rm -v /usr/local/bin/fn:/install devo/installer URL_TO_BIN`)
  * Try to use for this as a test too: https://github.com/justjanne/powerline-go/issues/13 . `dj install powerline && powerline update` or `dj install-only powerline` which won't add it to the bin dir. Run it from home dir and it will create .powerline, download bin and add the lines to zshrc (cp .bak, load it, search for powerline, if not there, add it to end and write it).
