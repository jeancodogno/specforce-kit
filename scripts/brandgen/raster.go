package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"

	"golang.org/x/image/vector"
)

// rasterSize returns the pixel size of an artwork rendered at the given height.
func rasterSize(a artwork, height int) (int, int) {
	w := int(math.Round(float64(height) * a.Width / Grid))
	return w, height
}

// addPath feeds a path to the rasterizer with every coordinate multiplied by k.
func addPath(r *vector.Rasterizer, path []seg, k float64) {
	f := func(p pt) (float32, float32) { return float32(p.X * k), float32(p.Y * k) }
	for _, s := range path {
		switch s.Op {
		case opMove:
			r.MoveTo(f(s.Pts[0]))
		case opLine:
			r.LineTo(f(s.Pts[0]))
		case opQuad:
			x1, y1 := f(s.Pts[0])
			x2, y2 := f(s.Pts[1])
			r.QuadTo(x1, y1, x2, y2)
		case opCube:
			x1, y1 := f(s.Pts[0])
			x2, y2 := f(s.Pts[1])
			x3, y3 := f(s.Pts[2])
			r.CubeTo(x1, y1, x2, y2, x3, y3)
		case opClose:
			r.ClosePath()
		}
	}
}

// renderPNG rasterizes the artwork at the given pixel height and encodes it.
func renderPNG(a artwork, height int) ([]byte, error) {
	w, h := rasterSize(a, height)
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	k := float64(height) / Grid
	r := vector.NewRasterizer(w, h)
	for _, s := range a.Shapes {
		r.Reset(w, h)
		addPath(r, s.Path, k)
		r.Draw(dst, dst.Bounds(), image.NewUniform(s.Fill), image.Point{})
	}
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, dst); err != nil {
		return nil, fmt.Errorf("encode %s-%d.png: %w", a.Name, height, err)
	}
	return buf.Bytes(), nil
}
