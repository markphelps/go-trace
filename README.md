# go-trace

Implementing a Path Tracer so that I can teach myself the Go language.

Now with more concurrency! :sparkles:

Following the E-book by Peter Shirley ['Ray Tracing in One Weekend'](http://www.amazon.com/Ray-Tracing-Weekend-Peter-Shirley-ebook/dp/B01B5AODD8), but in Go instead of C++.

I also wrote a series of posts chronicling the process on [my blog](https://www.markphelps.me/2016/03/15/writing-a-ray-tracer-in-go/).

!['Scene'](https://s3.amazonaws.com/markphelps.me/2016/out.png)

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/markphelps/go-trace)
[![Go Report Card](https://goreportcard.com/badge/github.com/markphelps/example?style=flat-square)](https://goreportcard.com/report/github.com/markphelps/go-trace)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/markphelps)

## Options

```
$ ./bin/go-trace -help
Usage of ./bin/go-trace:
  -a value
        camera aperture (default 0.010000)
  -cpus value
        number of CPUs to use (default 8)
  -fov value
        vertical field of view (degrees) (default 75.000000)
  -h value
        height of image (pixels) (default 500)
  -n value
        number of samples per pixel for AA (default 100)
  -o value
        output filename (default out.png)
  -version
        show version and exit
  -w value
        width of image (pixels) (default 600)
  -x float
        look from X (default 10)
  -y float
        look from Y (default 4)
  -z float
        look from Z (default 6)
```

## Running

```
$ make
go build -o bin/go-trace github.com/markphelps/go-trace/cmd/go-trace
```

```
$ ./bin/go-trace

Rendering 600 x 500 pixel scene with 486 objects:
[4 cpus, 100 samples/pixel, 75.00° fov, 0.01 aperture]

[==========----------------------------------------------------------------------] 12.80%
```

## Tests

```
$ make test
```

## Output

The default output format is [PNG](https://en.wikipedia.org/wiki/Portable_Network_Graphics) and the default filename is **out.png**.

This can be changed by using the `-o` flag as described in the help output.

Ex: `$ ./bin/go-trace -o output.jpeg`

Current supported file extensions:

* .png
* .jpg
* .jpeg
