# dj

```sh
make install
dj install devo/ruby --name ruby2 --force
ruby2 -v
```

## TODO

* Make script to build and deploy the dockers, rather than a build.sh in each of them. Use directory name for name.
* Add insaller script, could we use docker for it?  mnt /usr/local/bin/fn and throw it there (make generic installer, `docker run --rm -v /usr/local/bin/fn:/install devo/installer URL_TO_BIN`)
