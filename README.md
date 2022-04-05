# go-imagequant

This package wraps some (not all) functionality of
[libimagequant](https://pngquant.org/lib/).

Please follow the [instructions](https://github.com/ImageOptim/libimagequant/tree/main/imagequant-sys) to get a working
c lib.

See [cmd/go-imagequant.go](cmd/go-imagequant.go) how to use this package. 

## Using docker to use the command line binary.

The example cli binary supports reading and writing png, jpeg and gif (non animated) and serves as an example how to use this package.

See docker/alpine/Dockerfile for details. 

Requirements:

- docker
- make

call

```bash
make docker-alpine
```

to create a local image. Check if image was created:

```
docker image ls
REPOSITORY      TAG       IMAGE ID       CREATED        SIZE
go-imagequant   latest    6b9f9364ab77   10 hours ago   51.8MB
<none>          <none>    f5648e483897   10 hours ago   64.3MB
<none>          <none>    139c553f6e4a   10 hours ago   51.8MB
```

You should see image go-imagequant.

---
Docker run:

```
docker run go-imagequant                           
  -dst string
        dst image path
  -gamma float
        gamma
  -max.quality uint
        min allowed quality (default 100)
  -min.quality uint
        min allowed quality (default 0)
  -showLibImageQuantVersion
        show lib image quant version and exit
  -speed uint
        speed to to use (default 4)
  -src string
        src image path
no src image

```

Examples:

Converting a png assuming you got image ```benchmark_image_1.png``` in your current path:

```
docker run -v $PWD:/tmp go-imagequant -max.quality 75 -src /tmp/benchmark_image_1.png -dst /tmp/benchmark_image_1.docker.max40.png
```
