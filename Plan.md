# Feature Implementation Plan: Support Fallback to Init Containers on `-c` Option

## Goal
Enhance the tool so that when the `-c` option is specified (selecting a container), if the specified container is not found in the main containers list, the tool should also search in `initContainers` as a fallback.

This is especially useful for scenarios where helper/side containers are defined as init containers with a restart policy of Always.

---

## Breakdown of Implementation Steps

### 1. Analyze Current `-c` Handling
- Identify where in the codebase the container selection takes place when using the `-c` flag/parameter.
- Determine where the search is done in the pod/container spec (likely only in main containers currently).

### 2. Update the Container Lookup Logic
- Modify the lookup to:
  - First, search in the standard containers list for the specified container name.
  - If not found, search in the `initContainers` list.
- Ensure the logic is efficient and does not introduce regressions for existing functionality.

### 3. Add/Update Tests
- Add or update unit/integration tests that:
  - Test for main container selection.
  - Test for successful fallback to `initContainers` when main container does not exist.
  - Test error behavior when no such container in either list.

### 4. Documentation
- Update command-line help and user documentation to note new lookup fallback behavior.

### 5. Edge/Compatibility Cases
- Ensure behavior is compatible with pods/specs:
  - With no `initContainers`.
  - Where both main and init containers have same name (define priority/order clearly).

### 6. Pre-commit/Lint/CI
- Run pre-commit hooks and CI to validate changes pass code quality checks.

---

## Deliverables
- [ ] Updated code for container lookup logic
- [ ] Added/updated tests
- [ ] Documentation updates (inline help and/or project docs)
- [ ] Plan.md with this strategy

