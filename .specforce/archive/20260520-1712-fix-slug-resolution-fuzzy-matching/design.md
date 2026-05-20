---
slug: 20260520-1712-fix-slug-resolution-fuzzy-matching
lens: Bugfix
---

# Technical Design: Fuzzy Slug Resolution (Fix Blueprint)

## 1. Code Path Inventory
- `src/internal/spec/slug.go` -> Implement `ResolveSlug(projectRoot, slug string) string` to handle the resolution logic.
- `src/internal/spec/status.go` -> Update `GetStatus` to use `ResolveSlug`.
- `src/internal/spec/service.go` -> Update `GetImplementationStatus` and `UpdateTaskStatus` to use `ResolveSlug`.
- `src/internal/spec/archive.go` -> Update `ArchiveSpec` to use `ResolveSlug`.

## 2. Regression Strategy (Verification Plan)
- **Unit Tests:** 
    - Add `TestResolveSlug` to `src/internal/spec/slug_test.go` covering exact matches, timestamped matches, sub-paths, and archive resolution.
    - Test the "Prioritize Newest" logic by creating two timestamped directories for the same base slug.
- **Manual Verification:** 
    - `specforce spec status <base-slug>` (active and archived).
    - `specforce implementation status <base-slug>`.
    - `specforce spec archive <base-slug>`.

## 3. Side Effects & Risks
- **Performance:** Scanning directories (`os.ReadDir`) is O(N) where N is the number of specs. For typical projects (dozens to hundreds of specs), this is negligible.
- **Compatibility:** No changes to existing files or external APIs. The resolution is transparent to the user.

## 4. Proposed Fix (Abstract Logic)

```go
func ResolveSlug(projectRoot string, slug string) string {
    // 1. Exact match
    if exists(slug) return slug

    // 2. Fuzzy match in specs/
    matches := findFuzzy(projectRoot, "specs", slug)
    if len(matches) > 0 {
        sort(matches) // Newest first
        return matches[0]
    }

    // 3. Fuzzy match in archive/
    matches = findFuzzy(projectRoot, "archive", slug)
    if len(matches) > 0 {
        sort(matches) // Newest first
        return matches[0]
    }

    return slug // Fallback
}
```

### Ambiguitiy Resolution
The `findFuzzy` logic will use the `timestampRegex` (`^\d{8}-\d{4}-`) to identify the timestamp prefix and compare the remaining part with the user-provided slug. If multiple matches are found, lexicographical sorting on the full directory name ensures that higher timestamps come last (or first if reversed).
