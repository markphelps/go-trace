# go-trace [WIP]

Implementing a Path Tracer so that I can teach myself the Go language.

Following the E-book by Peter Shirley ['Ray Tracing in One Weekend'](http://www.amazon.com/Ray-Tracing-Weekend-Peter-Shirley-ebook/dp/B01B5AODD8), but in Go instead of C++.

## Current Status

!['Purple/Green'](https://s3.amazonaws.com/markphelps.me/2016/purple-green-aa.png)

It currently writes the image as a [PPM file](http://netpbm.sourceforge.net/doc/ppm.html), mainly because it's a super simple file format.

I may add support for more formats later.

## Running

`$ go run main.go`

## Viewing

To view the image in OS X, download [ToyViewer](https://itunes.apple.com/us/app/toyviewer/id414298354?mt=12).
