# Instructions for Claude AI Agent

## Project Documentation

Read before starting any task:
- `DESIGN.md` - System architecture and data schemas
- `ROADMAP.md` - Implementation tasks with acceptance criteria
- `PROJECT_STRUCTURE.md` - Directory organization

## Code Style

- **Concise code**: Clear and minimal
- **No comments**: Code should be self-documenting
- **No emoji**: Never use emoji anywhere
- **Follow exact schemas**: Data models in DESIGN.md are contracts

## Implementation Rules

### Task Execution
1. Read the specific task in ROADMAP.md
2. Implement minimum viable functionality
3. Test thoroughly with appropriate test types (unit/integration/e2e)
4. Run code and verify functionality works with no regressions
5. Meet all acceptance criteria

### Service Requirements
- Health endpoints (`/health`, `/ready`)
- Environment variable configuration
- Graceful shutdown
- Structured logging

### Error Handling
- Log with context for debugging
- Fail fast on configuration errors
- Continue on recoverable errors

## What NOT to Do

- Don't add unspecified features
- Don't over-engineer solutions
- Don't write obvious comments
- Don't use placeholder TODOs
- Don't include decorative elements
