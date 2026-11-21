package main

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const (
	wordmarkText    = "specforce"
	wordmarkSplit   = 4 // glyphs 1-4 are "spec"
	wordmarkEm      = 40.0
	wordmarkOriginX = 88.0
	wordmarkBase    = 50.0
	loadPPEM        = 1000
)

// wordmarkShapes outlines "specforce" in Go Mono Bold (design 2.3), one shape per glyph.
func wordmarkShapes(p palette) ([]shape, error) {
	f, err := sfnt.Parse(gomonobold.TTF)
	if err != nil {
		return nil, fmt.Errorf("parse Go Mono Bold: %w", err)
	}
	var buf sfnt.Buffer
	ppem := fixed.I(loadPPEM)
	k := wordmarkEm / loadPPEM
	x := wordmarkOriginX
	shapes := make([]shape, 0, len(wordmarkText))
	for i, r := range wordmarkText {
		idx, err := f.GlyphIndex(&buf, r)
		if err != nil {
			return nil, fmt.Errorf("glyph %q: %w", r, err)
		}
		segs, err := f.LoadGlyph(&buf, idx, ppem, nil)
		if err != nil {
			return nil, fmt.Errorf("load glyph %q: %w", r, err)
		}
		adv, err := f.GlyphAdvance(&buf, idx, ppem, font.HintingNone)
		if err != nil {
			return nil, fmt.Errorf("advance %q: %w", r, err)
		}
		fill := p.Mint
		if i < wordmarkSplit {
			fill = p.SpecText
		}
		shapes = append(shapes, shape{Fill: fill, Path: glyphPath(segs, x, k)})
		x += fixedToFloat(adv) * k
	}
	return shapes, nil
}

func fixedToFloat(v fixed.Int26_6) float64 { return float64(v) / 64 }

// glyphPath converts sfnt segments to a path scaled by k at origin (x, baseline).
func glyphPath(segs sfnt.Segments, x, k float64) []seg {
	conv := func(p fixed.Point26_6) pt {
		return pt{x + fixedToFloat(p.X)*k, wordmarkBase + fixedToFloat(p.Y)*k}
	}
	var path []seg
	for i, s := range segs {
		if s.Op == sfnt.SegmentOpMoveTo && i > 0 {
			path = append(path, seg{Op: opClose})
		}
		n, o := segArity(s.Op)
		pts := make([]pt, n)
		for j := 0; j < n; j++ {
			pts[j] = conv(s.Args[j])
		}
		path = append(path, seg{Op: o, Pts: pts})
	}
	if len(path) > 0 {
		path = append(path, seg{Op: opClose})
	}
	return path
}

func segArity(o sfnt.SegmentOp) (int, op) {
	switch o {
	case sfnt.SegmentOpMoveTo:
		return 1, opMove
	case sfnt.SegmentOpLineTo:
		return 1, opLine
	case sfnt.SegmentOpQuadTo:
		return 2, opQuad
	default:
		return 3, opCube
	}
}
