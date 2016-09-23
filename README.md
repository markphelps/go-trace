# go-trace

Implementing a Path Tracer so that I can teach myself the Go language.

Now with more concurrency! :sparkles:

Following the E-book by Peter Shirley ['Ray Tracing in One Weekend'](http://www.amazon.com/Ray-Tracing-Weekend-Peter-Shirley-ebook/dp/B01B5AODD8), but in Go instead of C++.

I also wrote a series of posts chronicling the process on [my blog](http://www.markphelps.me/2016/03/15/writing-a-ray-tracer-in-go.html).

!['Scene'](https://s3.amazonaws.com/markphelps.me/2016/out.png)

It currently writes the image as a [PNG file](https://en.wikipedia.org/wiki/Portable_Network_Graphics), but it would be easy to support JPG or GIF files too.

## Options

```
$ ./go-trace -help
Usage of ./go-trace:
  -a float
        camera aperture (default 0.01)
  -cpus int
        number of CPUs to use (default runtime.NumCPU())
  -fov float
        vertical field of view (degrees) (default 75)
  -h int
        height of image (pixels) (default 500)
  -o string
        output filename (default "out.png")
  -n int
        number of samples per pixel for AA (default 100)
  -w int
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
$ ./go-trace

Rendering 600 x 500 pixel scene with 486 objects:
[4 cpus, 100 samples/pixel, 75.00° fov, 0.01 aperture]

[==========----------------------------------------------------------------------] 12.80%
```
