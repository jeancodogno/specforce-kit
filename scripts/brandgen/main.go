// Command brandgen generates the official Specforce brand assets
// (assets/brand) from the brand definition in this package.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// output is one generated file held in memory.
type output struct {
	name string
	data []byte
}

// build renders every output in the fixed order of design table 2.5.
func build() ([]output, error) {
	arts, err := artworks()
	if err != nil {
		return nil, err
	}
	var outs []output
	for _, a := range arts {
		outs = append(outs, output{a.Name + ".svg", renderSVG(a.Width, Grid, a.Shapes)})
		for _, h := range a.Sizes {
			data, err := renderPNG(a, h)
			if err != nil {
				return nil, err
			}
			outs = append(outs, output{fmt.Sprintf("%s-%d.png", a.Name, h), data})
		}
	}
	return outs, nil
}

// writeAtomic writes data to path via a sibling temp file and rename.
func writeAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// writeAll writes the outputs into dir, stopping at the first failure.
func writeAll(dir string, outs []output, stdout io.Writer) (string, error) {
	for _, o := range outs {
		path := filepath.Join(dir, o.name)
		if err := writeAtomic(path, o.data); err != nil {
			return path, err
		}
		_, _ = fmt.Fprintf(stdout, "wrote %s\n", path)
	}
	return "", nil
}

// run is the testable entry point; it returns the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("brandgen", flag.ContinueOnError)
	fs.SetOutput(stderr)
	out := fs.String("out", "assets/brand", "output directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	outs, err := build()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "brandgen: %v\n", err)
		return 1
	}
	if err := os.MkdirAll(*out, 0o750); err != nil {
		_, _ = fmt.Fprintf(stderr, "brandgen: cannot write %s: %v\n", *out, err)
		return 1
	}
	if path, err := writeAll(*out, outs, stdout); err != nil {
		_, _ = fmt.Fprintf(stderr, "brandgen: cannot write %s: %v\n", path, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
