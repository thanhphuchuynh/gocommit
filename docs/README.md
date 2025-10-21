# GoCommit Documentation

Welcome to the GoCommit documentation! This collection provides comprehensive guides for using and understanding GoCommit.

## 📖 Documentation Index

### Getting Started
- Overview, installation, and basic usage
- **[Installation Guide](INSTALL.md)** - Detailed installation instructions

### Features & Architecture
- **[JSON Schema Implementation](JSON_SCHEMA_IMPLEMENTATION.md)** - LangChain-style structured outputs for AI providers ⭐
- **[Delayed Commit Architecture](DELAYED_COMMIT_ARCHITECTURE.md)** - Schedule commits with smart time suggestions
- **[Application Flow](FLOW.md)** - Architecture and workflow diagrams
- **[Refactoring Summary](REFACTORING_SUMMARY.md)** - Codebase improvements and structure

### Configuration & Integration
- **[GitHub Actions Integration](README_GITHUB_ACTIONS.md)** - Automate commit generation in CI/CD
- **[Logging Guide](README_LOGGING.md)** - Debug and monitor your workflow
- **[Logging Configuration](README_LOGGING_CONFIG.md)** - Configure logging settings
- **[Delayed Commit Setup](README_DELAYED_COMMIT.md)** - Configure time-based commits

### Additional Resources
- **[Parsing Fix Summary](PARSING_FIX_SUMMARY.md)** - Technical details on parsing improvements
- **[Refactoring Design](REFACTORING_DESIGN.md)** - Design decisions and patterns
- **[Dev.to Blog Post](gocommit-devto-blog.md)** - Blog article about GoCommit

## 🌟 Highlighted Features

### JSON Schema & Structured Outputs
GoCommit implements LangChain-style structured outputs for consistent JSON responses from AI providers. See [JSON_SCHEMA_IMPLEMENTATION.md](JSON_SCHEMA_IMPLEMENTATION.md) for details.

**Important**: Only certain models support strict JSON schema validation:
- ✅ OpenAI GPT-4o and newer
- ✅ Fireworks-provided models
- ⚠️ Other models rely on enhanced prompting + fallback parsing

### Multi-Provider Support
- **Gemini** via Google AI API
- **OpenRouter** for access to Claude, GPT-4, Llama, and more

### Delayed Commits
Schedule commits for later with intelligent time suggestions based on your work hours.

## 🚀 Quick Links

- [GitHub Repository](https://github.com/thanhphuchuynh/gocommit)
- [Issue Tracker](https://github.com/thanhphuchuynh/gocommit/issues)
- [Releases](https://github.com/thanhphuchuynh/gocommit/releases)

## 📝 Contributing

Contributions are welcome! Please see the main repository for contribution guidelines.

## 📄 License

MIT License - see repository for details
