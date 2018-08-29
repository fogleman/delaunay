## Delaunay Triangulation

Fast Delaunay triangulation implemented in Go.

This code was ported from [Mapbox's Delaunator project](https://github.com/mapbox/delaunator) (JavaScript).

### Installation

    $ go get -u github.com/fogleman/delaunay

### Documentation

https://godoc.org/github.com/fogleman/delaunay

See https://mapbox.github.io/delaunator/ for more information about the `Triangles` and `Halfedges` data structures.

### Usage

```go
var points []delaunay.Point
// populate points...
triangulation, err := delaunay.Triangulate(points)
// handle err...
// use triangulation.Triangles, triangulation.Halfedges
```

![Example](https://i.imgur.com/xhfW1EV.png)
