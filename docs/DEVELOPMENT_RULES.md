# AgriSmart Development Rules

## Working model
AgriSmart is developed milestone-by-milestone. Each milestone has a bounded scope and explicit verification criteria.

## Workflow
1. Inspect the current repository state.
2. Read the relevant project documentation.
3. Define the requirement before coding.
4. Implement the smallest complete change.
5. Run automated tests and static checks.
6. Manually verify user-facing behavior when applicable.
7. Record findings and limitations.
8. Commit with a focused message.

## Branching
- `main` is the stable integration branch.
- Feature/foundation work uses focused branches.
- Avoid direct development on `main` for non-trivial changes.

## Scope control
- Do not add unrelated features to an active milestone.
- Defer ideas to the roadmap or an issue.
- Prefer incremental changes over rewrites.

## Data and knowledge quality
Agricultural knowledge is treated as product data, not casual copy.
Each important record should have provenance and a validation state.

## Safety
Content concerning agricultural chemicals must be handled conservatively. Never invent label instructions, doses, PHI, compatibility, or legal/registration status.

## API standards
- Consistent JSON responses.
- Explicit validation errors.
- Authentication/authorization boundaries are documented.
- Business logic remains separate from HTTP transport where practical.

## Testing baseline
Every milestone should establish or preserve:
- unit tests for domain logic
- integration tests for critical persistence/API paths
- static analysis/linting appropriate to the stack
- deterministic local development setup

## Completion criteria
No task is complete solely because code compiles. Completion requires tests, verification, documentation updates when behavior or architecture changes, and a Git commit.
