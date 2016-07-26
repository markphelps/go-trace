# go-trace

Implementing a Path Tracer so that I can teach myself the Go language.

Following the E-book by Peter Shirley ['Ray Tracing in One Weekend'](http://www.amazon.com/Ray-Tracing-Weekend-Peter-Shirley-ebook/dp/B01B5AODD8), but in Go instead of C++.

Also writing a blog series chronicling the process on [my blog](http://www.markphelps.me/2016/03/15/writing-a-ray-tracer-in-go.html).

## Current Status

!['Scene'](https://s3.amazonaws.com/markphelps.me/2016/scene2.png)

It currently writes the image as a [PNG file](https://en.wikipedia.org/wiki/Portable_Network_Graphics), but it would be easy to support JPG or GIF files too.

## Running

```
$ go build

$ ./go-trace --help
Usage of ./go-trace:
  -aperture float
        camera aperture (default 0.01)
  -fov float
        vertical field of view (default 75)
  -height int
        height of image (default 500)
  -out string
        output filename (default "out.png")
  -samples int
        number of samples for anti-aliasing (default 100)
  -width int
        width of image (default 600)
  -x float
        look from X (default 10)
  -y float
        look from Y (default 4)
  -z float
        look from Z (default 6)
```
