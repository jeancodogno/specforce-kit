package agent

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/installer"
	"gopkg.in/yaml.v3"
)

// Transformer defines how to transform a blueprint into a target format.
type Transformer func(bp *core.Blueprint, mapping core.MappingConfig) string

var transformers = map[string]Transformer{
	".toml": func(bp *core.Blueprint, mapping core.MappingConfig) string {
		desc := bp.Metadata.Description
		if desc == "" {
			desc = mapping.Name
		}
		// Gemini TOML format: REQ-2 AC-5
		return fmt.Sprintf("description = %q\nprompt = \"\"\"\n%s\n\"\"\"\n", desc, bp.Content)
	},
}

// LoadKitConfig parses the kit.yaml configuration file. It prefers embedded defaults from kitFS
// and applies optional overrides from the project root if kit.yaml exists there.
func LoadKitConfig(kitFS fs.FS, root string) (*core.KitConfig, error) {
	// 1. Load from kitFS (embedded defaults)
	data, err := fs.ReadFile(kitFS, "kit.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded kit.yaml: %w", err)
	}

	var config core.KitConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse embedded kit.yaml: %w", err)
	}

	// 2. Load from root/kit.yaml (user overrides) - OPTIONAL
	overridePath, err := core.SecurePath(root, "kit.yaml")
	if err == nil {
		// #nosec G304 - Path is secured by SecurePath
		overrideData, err := os.ReadFile(overridePath)
		if err == nil {
			var override core.KitConfig
			if err := yaml.Unmarshal(overrideData, &override); err == nil {
				// Basic merge logic: replace tools defined in override
				if override.Tools != nil {
					if config.Tools == nil {
						config.Tools = make(map[string]core.ToolRoute)
					}
					for id, route := range override.Tools {
						config.Tools[id] = route
					}
				}
			}
		}
	}

	return &config, nil
}

// AdaptArtifacts reads blueprints from the given kitFS and adapts them for the specified agent.
func AdaptArtifacts(ctx context.Context, root string, kitFS fs.FS, agent string, ui core.UI, opts installer.Options) error {
	if ui != nil {
		ui.SubTask(fmt.Sprintf("Adapting artifacts for %s...", agent))
	}

	kitConfig, err := LoadKitConfig(kitFS, root)
	if err != nil {
		return err
	}

	tracker := newNameTracker()

	return fs.WalkDir(kitFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Only parse blueprints (YAML files with structured metadata and content)
		matched, _ := filepath.Match("*.yaml", d.Name())
		if !matched || d.Name() == "mapping.yaml" || d.Name() == "manifest.yaml" {
			return nil
		}

		// REQ-3: Validate name uniqueness
		if err := validateBlueprintHeaderName(kitFS, kitConfig, path, agent, tracker); err != nil {
			return err
		}

		return processBlueprint(ctx, root, kitFS, kitConfig, path, agent, opts)
	})
}

func validateBlueprintHeaderName(kitFS fs.FS, kitConfig *core.KitConfig, path string, agent string, tracker *nameTracker) error {
	data, err := fs.ReadFile(kitFS, path)
	if err != nil {
		return nil // Non-critical for uniqueness check if it fails here; processBlueprint will handle it
	}

	bp, err := core.ParseBlueprint(path, data)
	if err != nil {
		return nil
	}

	mappings, err := resolveMappings(kitConfig, path, bp, agent)
	if err != nil {
		return nil
	}

	category := strings.Split(filepath.ToSlash(path), "/")[0]
	for _, m := range mappings {
		if m.Ext == ".md" {
			name := resolveHeaderName(bp, m, category)
			if err := tracker.validate(name, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func processBlueprint(ctx context.Context, root string, kitFS fs.FS, kitConfig *core.KitConfig, path string, agent string, opts installer.Options) error {
	data, err := fs.ReadFile(kitFS, path)
	if err != nil {
		return fmt.Errorf("failed to read blueprint %s: %w", path, err)
	}

	bp, err := core.ParseBlueprint(path, data)
	if err != nil {
		return fmt.Errorf("failed to parse blueprint %s: %w", path, err)
	}

	mappings, err := resolveMappings(kitConfig, path, bp, agent)
	if err != nil {
		return err
	}

	toolRoute := kitConfig.Tools[agent]

	for _, mapping := range mappings {
		// REQ-2: Check if this artifact should be installed based on the target path
		if !installer.ShouldInstall(mapping.Path, opts) {
			continue
		}

		var targetDir, targetFile string
		if toolRoute.Security.GlobalWrite {
			targetDir = mapping.Path
			if !filepath.IsAbs(targetDir) {
				targetDir = filepath.Join(root, targetDir)
			}
			targetFile = filepath.Join(targetDir, mapping.Name+mapping.Ext)
		} else {
			var err error
			targetDir, err = core.SecurePath(root, mapping.Path)
			if err != nil {
				return fmt.Errorf("security: %w", err)
			}
			targetFile, err = core.SecurePath(root, filepath.Join(mapping.Path, mapping.Name+mapping.Ext))
			if err != nil {
				return fmt.Errorf("security: %w", err)
			}
		}

		if err := os.MkdirAll(targetDir, 0750); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
		}

		category := strings.Split(filepath.ToSlash(path), "/")[0]
		content := applyTransformation(bp, mapping, path, category)

		// #nosec G306 - Path is secured by SecurePath
		if err := os.WriteFile(targetFile, []byte(content), 0600); err != nil {
			return fmt.Errorf("failed to write adapted artifact %s: %w", targetFile, err)
		}
	}

	return nil
}

func resolveMappings(kitConfig *core.KitConfig, path string, bp *core.Blueprint, agent string) ([]core.MappingConfig, error) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) == 0 {
		return nil, nil
	}

	category := parts[0]
	slug := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	toolRoute, ok := kitConfig.Tools[agent]
	if !ok {
		return nil, core.ErrToolMappingNotFound
	}

	rawMappings := getRawMappings(kitConfig, bp, toolRoute, agent, category)
	if len(rawMappings) == 0 {
		return nil, nil
	}

	results := make([]core.MappingConfig, 0, len(rawMappings))
	for _, m := range rawMappings {
		mapping := m
		applyBlueprintOverrides(&mapping, bp, agent)

		if err := finalizeMapping(&mapping, toolRoute, slug, parts); err != nil {
			return nil, err
		}
		results = append(results, mapping)
	}

	return results, nil
}

func getRawMappings(kitConfig *core.KitConfig, bp *core.Blueprint, toolRoute core.ToolRoute, agent, category string) []core.MappingConfig {
	if specific, ok := toolRoute.Mappings[category]; ok {
		return specific
	}
	if bp.Metadata.Mapping != nil {
		if override, ok := bp.Metadata.Mapping[agent]; ok {
			return []core.MappingConfig{override}
		}
	}
	if defaults, ok := kitConfig.Defaults[category]; ok {
		return defaults
	}
	return nil
}

func finalizeMapping(mapping *core.MappingConfig, toolRoute core.ToolRoute, slug string, parts []string) error {
	initialPath := mapping.Path
	applyWildcards(mapping, slug, parts, initialPath)

	// Default naming logic: If mapping.Name is missing, resolve it based on source filename
	if mapping.Name == "" {
		mapping.Name = slug
	}

	// Resolve final target path
	rawTarget := mapping.Target
	if rawTarget == "" {
		rawTarget = toolRoute.Target
	}

	target := core.ExpandPath(rawTarget)
	if strings.HasPrefix(target, "~") {
		return fmt.Errorf("failed to resolve home directory for path: %s", target)
	}

	// Use filepath.Join to combine target and the mapping's internal path
	mapping.Path = filepath.Clean(filepath.Join(target, mapping.Path))
	return nil
}

func applyBlueprintOverrides(mapping *core.MappingConfig, bp *core.Blueprint, agent string) {
	if bp.Metadata.Mapping != nil {
		if override, ok := bp.Metadata.Mapping[agent]; ok {
			if override.Path != "" {
				mapping.Path = override.Path
			}
			if override.Name != "" {
				mapping.Name = override.Name
			}
			if override.Ext != "" {
				mapping.Ext = override.Ext
			}
		}
	}
}

func applyWildcards(mapping *core.MappingConfig, slug string, parts []string, initialPath string) {
	hasWildcardInPath := strings.Contains(mapping.Path, "*")
	hasWildcardInName := strings.Contains(mapping.Name, "*")
	mapping.Path = strings.ReplaceAll(mapping.Path, "*", slug)
	mapping.Name = strings.ReplaceAll(mapping.Name, "*", slug)

	if !hasWildcardInPath && !hasWildcardInName && len(parts) > 2 {
		subDir := filepath.Join(parts[1 : len(parts)-1]...)
		mapping.Path = filepath.Join(mapping.Path, subDir)
	}

	// If name has wildcard but path is fixed, we might want flat export.
	// This specifically handles the Antigravity/Codex case where Path was intended to be "."
	// but it was actually initialized with some value from the mapping.
	if hasWildcardInName && !hasWildcardInPath && initialPath == "." {
		mapping.Path = "."
	}
}

func applyTransformation(bp *core.Blueprint, mapping core.MappingConfig, sourcePath string, category string) string {
	if transformer, ok := transformers[mapping.Ext]; ok {
		return transformer(bp, mapping)
	}

	// REQ-2: Standardized YAML Frontmatter Headers for Markdown
	if mapping.Ext == ".md" {
		return injectYAMLHeader(bp, mapping, sourcePath, category)
	}

	return bp.Content
}

func injectYAMLHeader(bp *core.Blueprint, mapping core.MappingConfig, sourcePath string, category string) string {
	var header strings.Builder
	header.WriteString("---\n")

	// Resolve the display name: prefer metadata name, fallback to mapping name
	displayName := resolveHeaderName(bp, mapping, category)

	// REQ-2 AC-1: Agent header
	if category == "agents" {
		fmt.Fprintf(&header, "name: %s\n", displayName)
		fmt.Fprintf(&header, "description: %s\n", bp.Metadata.Description)
	}

	// REQ-2 AC-2 & AC-3: Skill header (only for primary SKILL.md)
	if category == "skills" || mapping.Name == "SKILL" {
		if mapping.Name == "SKILL" {
			fmt.Fprintf(&header, "name: %s\n", displayName)
			fmt.Fprintf(&header, "description: %s\n", bp.Metadata.Description)
			if bp.Metadata.Version != "" {
				fmt.Fprintf(&header, "version: %s\n", bp.Metadata.Version)
			}
			if bp.Metadata.Priority != "" {
				fmt.Fprintf(&header, "priority: %s\n", bp.Metadata.Priority)
			}
		} else if category == "skills" {
			// Skills secondary files don't get headers
			return bp.Content
		}
	}

	// REQ-2 AC-4: Command header (only if not already handled as SKILL)
	if category == "commands" && mapping.Name != "SKILL" {
		fmt.Fprintf(&header, "name: %s\n", displayName)
		fmt.Fprintf(&header, "description: %s\n", bp.Metadata.Description)
	}

	// If no header fields were added (unknown type), just return content
	if header.Len() == 4 { // Only "---\n"
		return bp.Content
	}

	header.WriteString("---\n\n")
	header.WriteString(bp.Content)
	return header.String()
}

// ResolveTemplatePath returns the path to a template, preferring the project-level
func ResolveTemplatePath(root, templateName string) string {
	projectPath, err := core.SecurePath(root, filepath.Join(".specforce/templates", templateName))
	if err != nil {
		return ""
	}
	if _, err := os.Stat(projectPath); err == nil {
		return projectPath
	}

	return projectPath
}
