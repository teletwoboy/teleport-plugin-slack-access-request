# Cursor Prompt

This project, `teleport-plugin-slack-access-request`, follows the guidelines below:

- All code must be written in **Go**.
- External dependencies should be **abstracted via interfaces** and injected where needed.
- Error handling must use the format: `fmt.Errorf("context: %w", err)`.
- Follow the code style defined in the `.golangci.yaml` file located in the root directory.

Commit messages must follow the **Conventional Commit** style and be written in **English**:
- feat: Add login handler
- fix: Handle token expiration
- chore: Update dependencies