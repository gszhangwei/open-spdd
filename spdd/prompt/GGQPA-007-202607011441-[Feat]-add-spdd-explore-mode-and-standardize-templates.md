1. **New Explore Mode Workflow**:
   - Create `spdd-explore.md` to establish a safe, read-only thinking phase ("Explore Mode") for AI coding assistants.
   - Enforce an absolute read-only policy, prohibiting code generation or file mutation in this mode.
   - Embed the "Distorted Mirror Technique" to surface user assumptions and validate initial scopes.
   - Design micro action segments (M1-M4) and clear transition criteria before initiating `/spdd-analysis`.

2. **Configuration & Paths Alignment**:
   - Update `Antigravity` configuration directory path mapping to `.agents/workflows` in `internal/detector/types.go`.

3. **Template Metadata Standardization**:
   - Standardize core Markdown templates (`spdd-analysis`, `spdd-generate`, `spdd-prompt-update`, `spdd-reasons-canvas`, `spdd-sync`) by stripping the leading slash (`/`) from their frontmatter `name` fields.
