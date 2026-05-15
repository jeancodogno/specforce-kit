package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// FragmentType defines the kind of memory fragment.
type FragmentType string

const (
	FragmentAction   FragmentType = "Action"
	FragmentLesson   FragmentType = "Lesson"
	FragmentDecision FragmentType = "Decision"
	FragmentContext  FragmentType = "Context"
)

// Fragment represents a single piece of project memory.
type Fragment struct {
	Date    time.Time    `yaml:"date"`
	Scope   string       `yaml:"scope"`
	Author  string       `yaml:"author"`
	Type    FragmentType `yaml:"type"`
	Title   string       `yaml:"title"`
	Content string       `yaml:"content"`
}

// MemorialService defines the operations for managing the distributed memorial.
type MemorialService interface {
	// Record adds a new fragment to the memorial directory.
	Record(ctx context.Context, fragment Fragment) error

	// Consolidate aggregates the Rules of Engagement and the latest fragments.
	Consolidate(ctx context.Context, limit int) (string, error)

	// Initialize sets up the memorial directory and initial ROUTING.md.
	Initialize(ctx context.Context) error
}

type memorialService struct {
	projectRoot string
}

// NewMemorialService creates a new instance of the memorial service.
func NewMemorialService(projectRoot string) MemorialService {
	return &memorialService{
		projectRoot: projectRoot,
	}
}

func (s *memorialService) Initialize(ctx context.Context) error {
	memorialDir, err := core.SecurePath(s.projectRoot, filepath.Join(".specforce", "memorial"))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(memorialDir, 0750); err != nil {
		return fmt.Errorf("failed to create memorial directory: %w", err)
	}

	// Migrate legacy monolithic file if it exists
	if err := s.migrateLegacy(memorialDir); err != nil {
		// Log but don't fail initialization
		fmt.Fprintf(os.Stderr, "Warning: Failed to migrate legacy memorial: %v\n", err)
	}

	routingPath := filepath.Join(memorialDir, "ROUTING.md")
	if _, err := os.Stat(routingPath); os.IsNotExist(err) {
		routingContent := `# Project Memorial (Distributed)

> **FOR AI AGENTS: RULES OF ENGAGEMENT**
> This is your cross-session memory. You MUST obey these rules:
> 1. **READ THIS CONTEXT FIRST** before starting any task.
> 2. **TRIAGE:** Read the most recent fragments to understand recent changes and established precedents.
> 3. **STRICT LIMITS:** The CLI will provide a consolidated view of the most relevant fragments.
> 4. **DISTILLATION:** If a lesson or decision becomes a permanent rule, it MUST be moved to the official Constitution files.
`
		if err := os.WriteFile(routingPath, []byte(routingContent), 0600); err != nil {
			return fmt.Errorf("failed to create ROUTING.md: %w", err)
		}
	}

	return nil
}

func (s *memorialService) migrateLegacy(memorialDir string) error {
	legacyPath := filepath.Join(s.projectRoot, ".specforce", "docs", "memorial.md")
	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		return nil
	}

	targetPath := filepath.Join(memorialDir, "legacy.md")
	if _, err := os.Stat(targetPath); err == nil {
		// Already migrated or legacy.md exists
		return nil
	}

	// #nosec G304 - Path is within project root and validated by secure components
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return err
	}

	// #nosec G306, G703 - permissions are restricted to owner; path is internal to .specforce and trusted
	if err := os.WriteFile(targetPath, data, 0600); err != nil {
		return err
	}

	// Rename old file to mark as deprecated (don't delete to be safe)
	deprecatedPath := legacyPath + ".deprecated"
	return os.Rename(legacyPath, deprecatedPath)
}

func (s *memorialService) Record(ctx context.Context, f Fragment) error {
	if f.Date.IsZero() {
		f.Date = time.Now()
	}

	memorialDir, err := core.SecurePath(s.projectRoot, filepath.Join(".specforce", "memorial"))
	if err != nil {
		return err
	}

	// Filename format: YYYYMMDD-HHMM-slug.md
	slug := strings.ToLower(strings.ReplaceAll(f.Scope, " ", "-"))
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(f.Title, " ", "-"))
	}
	// Sanitize slug
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	filename := fmt.Sprintf("%s-%s.md", f.Date.Format("20060102-1504"), slug)
	filePath := filepath.Join(memorialDir, filename)

	var builder strings.Builder
	builder.WriteString("---\n")
	fmt.Fprintf(&builder, "date: %s\n", f.Date.Format("2006-01-02"))
	fmt.Fprintf(&builder, "scope: %s\n", f.Scope)
	fmt.Fprintf(&builder, "author: %s\n", f.Author)
	fmt.Fprintf(&builder, "type: %s\n", f.Type)
	builder.WriteString("---\n\n")
	fmt.Fprintf(&builder, "# %s\n\n", f.Title)
	builder.WriteString(f.Content)
	builder.WriteString("\n")

	// #nosec G306 - permissions are restricted to owner
	if err := os.WriteFile(filePath, []byte(builder.String()), 0600); err != nil {
		return fmt.Errorf("failed to write fragment: %w", err)
	}

	return nil
}

func (s *memorialService) Consolidate(ctx context.Context, limit int) (string, error) {
	memorialDir, err := core.SecurePath(s.projectRoot, filepath.Join(".specforce", "memorial"))
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	// 1. Read ROUTING.md (Rules of Engagement)
	routingPath := filepath.Join(memorialDir, "ROUTING.md")
	// #nosec G304 - Path is within project root and validated by secure components
	routingData, err := os.ReadFile(routingPath)
	if err == nil {
		builder.Write(routingData)
		builder.WriteString("\n---\n\n")
	}

	// 2. Read Fragments
	entries, err := os.ReadDir(memorialDir)
	if err != nil {
		return builder.String(), nil
	}

	var fragmentFiles []string
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "ROUTING.md" || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		fragmentFiles = append(fragmentFiles, entry.Name())
	}

	// Sort fragments newest first (based on filename timestamp)
	sort.Sort(sort.Reverse(sort.StringSlice(fragmentFiles)))

	if limit > 0 && len(fragmentFiles) > limit {
		fragmentFiles = fragmentFiles[:limit]
	}

	for _, filename := range fragmentFiles {
		// #nosec G304 - Path is within project root and validated by secure components
		data, err := os.ReadFile(filepath.Join(memorialDir, filename))
		if err != nil {
			continue
		}
		fmt.Fprintf(&builder, "## Fragment: %s\n\n", filename)
		builder.Write(data)
		builder.WriteString("\n---\n\n")
	}

	return builder.String(), nil
}
