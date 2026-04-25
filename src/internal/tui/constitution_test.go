package tui

import (
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/constitution"
)

func TestRenderConstitutionStatus(t *testing.T) {
	status := constitution.ConstitutionStatus{
		Artifacts: []constitution.ArtifactStatus{
			{Name: "principles.md", Description: "Core values", Exists: true},
			{Name: "security.md", Description: "Auth rules", Exists: false},
		},
	}

	got := RenderConstitutionStatus(status)

	if !strings.Contains(got, "principles.md") {
		t.Errorf("RenderConstitutionStatus() = %q, want it to contain %q", got, "principles.md")
	}
	if !strings.Contains(got, "security.md") {
		t.Errorf("RenderConstitutionStatus() = %q, want it to contain %q", got, "security.md")
	}
	if !strings.Contains(got, BulletGlyph) {
		t.Errorf("RenderConstitutionStatus() should contain BulletGlyph for exists:true")
	}
	if !strings.Contains(got, EmptyBulletGlyph) {
		t.Errorf("RenderConstitutionStatus() should contain EmptyBulletGlyph for exists:false")
	}
}
