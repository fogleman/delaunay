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

### Performance

3.3 GHz Intel Core i5

| # of Points | Time |
| ---: | ---: |
| 10 | 1.559µs |
| 100 | 37.645µs |
| 1,000 | 485.625µs |
| 10,000 | 5.552ms |
| 100,000 | 79.895ms |
| 1,000,000 | 1.272s |
| 10,000,000 | 23.481s |

![Example](https://i.imgur.com/xhfW1EV.png)
