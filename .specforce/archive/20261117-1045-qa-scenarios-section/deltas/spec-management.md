# Delta: spec-management
## ADDED
- `[NEW:mitigations-by-title]` The implementation status identifies a roadmap's pre-emptive mitigations section by its title, whatever section number it carries, exposes its text up to the next top-level section, and leaves the mitigations empty when the roadmap has no such section; roadmaps numbering it 3 keep the same output.
### [NEW:us-mitigations-by-title] Mitigations Shown Regardless of Section Number
- **Scenario:** Renumbered section
  - **GIVEN** a roadmap whose mitigations section is numbered 4
  - **WHEN** the implementation status is requested
  - **THEN** the mitigations text is shown.
- **Scenario:** Legacy numbering
  - **GIVEN** a roadmap whose mitigations section is numbered 3
  - **WHEN** the implementation status is requested
  - **THEN** the output is unchanged.
- **Scenario:** Followed by another section
  - **GIVEN** a mitigations section followed by another top-level section
  - **WHEN** the implementation status is requested
  - **THEN** only the mitigations text is exposed.
## MODIFIED
## REMOVED
