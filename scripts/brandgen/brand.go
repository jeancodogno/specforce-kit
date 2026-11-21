package main

import "image/color"

// Grid is the side of the symbol's square unit grid (design 2.2).
const Grid = 72.0

// Lockup viewBox width (design 2.3).
const lockupWidth = 312.0

// pt is a point in grid units, y down.
type pt struct{ X, Y float64 }

// op is a path operation.
type op byte

const (
	opMove  op = 'M'
	opLine  op = 'L'
	opQuad  op = 'Q'
	opCube  op = 'C'
	opClose op = 'Z'
)

// seg is one path segment; Pts holds the control and end points.
type seg struct {
	Op  op
	Pts []pt
}

// shape is one filled path in a single color.
type shape struct {
	Fill color.NRGBA
	Path []seg
}

// palette holds the role colors of design 2.1.
type palette struct {
	Mint, Cyan, Ice, SpecText color.NRGBA
}

func rgb(v uint32) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xFF}
}

var (
	darkPalette  = palette{Mint: rgb(0x00FA9A), Cyan: rgb(0x00FFFF), Ice: rgb(0xE0FFFF), SpecText: rgb(0xFFFFFF)}
	lightPalette = palette{Mint: rgb(0x007A4D), Cyan: rgb(0x00727F), Ice: rgb(0x3F5F6B), SpecText: rgb(0x1A1A1A)}
	tileColor    = rgb(0x0B0B0B)
)

var leftBracket = []pt{
	{6.5, 6.5}, {25.5, 6.5}, {25.5, 13.5}, {13.5, 13.5},
	{13.5, 58.5}, {25.5, 58.5}, {25.5, 65.5}, {6.5, 65.5},
}

// polygon turns a closed point list into a path.
func polygon(pts []pt) []seg {
	path := []seg{{Op: opMove, Pts: []pt{pts[0]}}}
	for _, p := range pts[1:] {
		path = append(path, seg{Op: opLine, Pts: []pt{p}})
	}
	return append(path, seg{Op: opClose})
}

func mirror(pts []pt) []pt {
	out := make([]pt, len(pts))
	for i, p := range pts {
		out[i] = pt{Grid - p.X, p.Y}
	}
	return out
}

// checkMark is the check polygon of design 2.2 (mitered join, square caps).
var checkMark = []pt{
	{19.05, 37.00}, {33.22, 51.17}, {53.93, 26.58},
	{48.58, 22.07}, {32.78, 40.83}, {24.00, 32.05},
}

// gateShapes returns the Gate symbol in draw order: left, right, check.
func gateShapes(p palette) []shape {
	return []shape{
		{Fill: p.Mint, Path: polygon(leftBracket)},
		{Fill: p.Cyan, Path: polygon(mirror(leftBracket))},
		{Fill: p.Ice, Path: polygon(checkMark)},
	}
}

// roundedTile is the app-icon tile of design 2.4.
func roundedTile() shape {
	const s, r, k = Grid, 16.0, 0.5523
	c := r * k
	path := []seg{
		{opMove, []pt{{r, 0}}},
		{opLine, []pt{{s - r, 0}}},
		{opCube, []pt{{s - r + c, 0}, {s, r - c}, {s, r}}},
		{opLine, []pt{{s, s - r}}},
		{opCube, []pt{{s, s - r + c}, {s - r + c, s}, {s - r, s}}},
		{opLine, []pt{{r, s}}},
		{opCube, []pt{{r - c, s}, {0, s - r + c}, {0, s - r}}},
		{opLine, []pt{{0, r}}},
		{opCube, []pt{{0, r - c}, {r - c, 0}, {r, 0}}},
		{opClose, nil},
	}
	return shape{Fill: tileColor, Path: path}
}

// transform applies scale(k) then translate(t) to every point of the shapes.
func transform(shapes []shape, k, t float64) []shape {
	out := make([]shape, len(shapes))
	for i, s := range shapes {
		path := make([]seg, len(s.Path))
		for j, sg := range s.Path {
			pts := make([]pt, len(sg.Pts))
			for m, p := range sg.Pts {
				pts[m] = pt{p.X*k + t, p.Y*k + t}
			}
			path[j] = seg{Op: sg.Op, Pts: pts}
		}
		out[i] = shape{Fill: s.Fill, Path: path}
	}
	return out
}

func appIconShapes() []shape {
	return append([]shape{roundedTile()}, transform(gateShapes(darkPalette), 0.72, 10.08)...)
}

// artwork is one master of table 2.5 and the raster heights derived from it.
type artwork struct {
	Name   string
	Width  float64 // viewBox width; height is always Grid
	Shapes []shape
	Sizes  []int // raster heights in pixels
}

var squareSizes = []int{16, 32, 64, 128, 256, 512}

// artworks returns the masters in the fixed order of table 2.5.
func artworks() ([]artwork, error) {
	darkWord, err := wordmarkShapes(darkPalette)
	if err != nil {
		return nil, err
	}
	lightWord, err := wordmarkShapes(lightPalette)
	if err != nil {
		return nil, err
	}
	return []artwork{
		{"gate-dark", Grid, gateShapes(darkPalette), squareSizes},
		{"gate-light", Grid, gateShapes(lightPalette), squareSizes},
		{"app-icon", Grid, appIconShapes(), squareSizes},
		{"lockup-dark", lockupWidth, append(gateShapes(darkPalette), darkWord...), []int{128, 256}},
		{"lockup-light", lockupWidth, append(gateShapes(lightPalette), lightWord...), []int{128, 256}},
	}, nil
}
