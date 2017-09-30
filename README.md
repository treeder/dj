# dj

```sh
make install
dj install devo/ruby --name ruby2 --force
ruby2 -v
```

## TODO

* Make script to build and deploy the dockers, rather than a build.sh in each of them. Use directory name for name.
* Add insaller script, could we use docker for it?  mnt /usr/local/bin/fn and throw it there (make generic installer, `docker run --rm -v /usr/local/bin/fn:/install devo/installer URL_TO_BIN`)
  * Try to use for this as a test too: https://github.com/justjanne/powerline-go/issues/13 . `dj install powerline && powerline update` or `dj install-only powerline` which won't add it to the bin dir. Run it from home dir and it will create .powerline, download bin and add the lines to zshrc (cp .bak, load it, search for powerline, if not there, add it to end and write it). 

