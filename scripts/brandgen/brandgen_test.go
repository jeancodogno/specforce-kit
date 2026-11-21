package main

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildOrFail(t *testing.T) []output {
	t.Helper()
	outs, err := build()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	return outs
}

func findOutput(t *testing.T, outs []output, name string) []byte {
	t.Helper()
	for _, o := range outs {
		if o.name == name {
			return o.data
		}
	}
	t.Fatalf("output %s not found", name)
	return nil
}

func TestOutputList(t *testing.T) {
	want := []string{"gate-dark.svg"}
	sizes := []string{"16", "32", "64", "128", "256", "512"}
	add := func(base string, sz []string) {
		for _, s := range sz {
			want = append(want, base+"-"+s+".png")
		}
	}
	add("gate-dark", sizes)
	want = append(want, "gate-light.svg")
	add("gate-light", sizes)
	want = append(want, "app-icon.svg")
	add("app-icon", sizes)
	want = append(want, "lockup-dark.svg")
	add("lockup-dark", []string{"128", "256"})
	want = append(want, "lockup-light.svg")
	add("lockup-light", []string{"128", "256"})

	outs := buildOrFail(t)
	if len(outs) != 27 || len(want) != 27 {
		t.Fatalf("got %d outputs, want 27", len(outs))
	}
	for i, o := range outs {
		if o.name != want[i] {
			t.Errorf("output %d = %q, want %q", i, o.name, want[i])
		}
	}
}

func TestDeterministicBuild(t *testing.T) {
	a := buildOrFail(t)
	b := buildOrFail(t)
	if len(a) != len(b) {
		t.Fatalf("output count differs: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i].name != b[i].name || !bytes.Equal(a[i].data, b[i].data) {
			t.Errorf("output %s differs between builds", a[i].name)
		}
	}
}

func TestSVGHasNoTextStrokeOrGradient(t *testing.T) {
	for _, o := range buildOrFail(t) {
		if !strings.HasSuffix(o.name, ".svg") {
			continue
		}
		s := string(o.data)
		for _, bad := range []string{"<text", "stroke", "Gradient"} {
			if strings.Contains(s, bad) {
				t.Errorf("%s contains %q", o.name, bad)
			}
		}
		if !strings.Contains(s, "<path fill=\"#") {
			t.Errorf("%s has no filled path", o.name)
		}
	}
}

func TestPNGSizes(t *testing.T) {
	cases := map[string][2]int{
		"gate-dark-16.png":     {16, 16},
		"app-icon-512.png":     {512, 512},
		"lockup-dark-128.png":  {555, 128},
		"lockup-light-256.png": {1109, 256},
	}
	outs := buildOrFail(t)
	for name, wh := range cases {
		img, err := png.Decode(bytes.NewReader(findOutput(t, outs, name)))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if b := img.Bounds(); b.Dx() != wh[0] || b.Dy() != wh[1] {
			t.Errorf("%s size = %dx%d, want %dx%d", name, b.Dx(), b.Dy(), wh[0], wh[1])
		}
	}
}

// nearTile reports whether a pixel is visually the tile background (#0B0B0B).
func nearTile(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	const bg, tol = 0x0B, 0x30
	for _, c := range []uint32{r >> 8, g >> 8, b >> 8} {
		if c > bg+tol {
			return false
		}
	}
	return true
}

func TestAppIcon16BracketAndCheckStayDistinct(t *testing.T) {
	img, err := png.Decode(bytes.NewReader(findOutput(t, buildOrFail(t), "app-icon-16.png")))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	// Symbol x -> pixel: (0.72*x + 10.08) * 16/72.
	toPx := func(x float64) int { return int((0.72*x + 10.08) * 16 / 72) }
	innerEdge := toPx(13.5) // left bracket inner edge
	checkLeft := toPx(24)   // check centerline start A(24,37)
	const row = 7           // middle row of a 16 px icon: (16-1)/2
	found := false
	for x := innerEdge; x <= checkLeft; x++ {
		if nearTile(img, x, row) {
			found = true
		}
	}
	if !found {
		t.Errorf("no background pixel between bracket (x=%d) and check (x=%d) on row %d", innerEdge, checkLeft, row)
	}
}

func TestWriteIntoReadOnlyDirFails(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	dir := t.TempDir()
	existing := filepath.Join(dir, "gate-dark.svg")
	if err := os.WriteFile(existing, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	var stdout, stderr bytes.Buffer
	code := run([]string{"-out", dir}, &stdout, &stderr)
	if code == 0 {
		t.Fatalf("exit code = 0, want non-zero")
	}
	msg := stderr.String()
	if !strings.Contains(msg, "brandgen: cannot write ") || !strings.Contains(msg, "gate-dark.svg") {
		t.Errorf("stderr = %q, want message naming gate-dark.svg", msg)
	}
	got, err := os.ReadFile(existing)
	if err != nil || string(got) != "old" {
		t.Errorf("pre-existing file changed: %q, %v", got, err)
	}
	if _, err := os.Stat(existing + ".tmp"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("temp file left behind: %v", err)
	}
}
