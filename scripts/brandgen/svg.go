package main

import (
	"fmt"
	"strconv"
	"strings"
)

func num(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }

// pathData serializes a path to SVG "d" syntax.
func pathData(path []seg) string {
	var b strings.Builder
	for _, s := range path {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteByte(byte(s.Op))
		for _, p := range s.Pts {
			b.WriteByte(' ')
			b.WriteString(num(p.X))
			b.WriteByte(',')
			b.WriteString(num(p.Y))
		}
	}
	return b.String()
}

// renderSVG serializes the shapes as a deterministic SVG document.
func renderSVG(width, height float64, shapes []shape) []byte {
	var b strings.Builder
	w, h := num(width), num(height)
	fmt.Fprintf(&b, "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 %s %s\" width=\"%s\" height=\"%s\">\n", w, h, w, h)
	for _, s := range shapes {
		fmt.Fprintf(&b, "<path fill=\"#%02X%02X%02X\" d=\"%s\"/>\n", s.Fill.R, s.Fill.G, s.Fill.B, pathData(s.Path))
	}
	b.WriteString("</svg>\n")
	return []byte(b.String())
}
