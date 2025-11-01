# 🎯 GoCommit Examples & Usage Guide

> Practical examples for all features described in the roadmap

## ⚙️ Phase 0 — Custom Templates & Config

### Example 1: Basic Config File (`.gocommitrc`)

**YAML format:**
```yaml
# .gocommitrc
ai:
  provider: "anthropic"  # or "openai", "local"
  model: "claude-sonnet-4"
  temperature: 0.7

commit:
  # Define allowed commit types
  types:
    - feat      # New feature
    - fix       # Bug fix
    - docs      # Documentation
    - style     # Formatting
    - refactor  # Code restructuring
    - test      # Testing
    - chore     # Maintenance

  # Default scope based on file patterns
  scopes:
    "src/api/*": "api"
    "src/db/*": "database"
    "src/ui/*": "ui"
    "docs/*": "docs"
    "*.test.ts": "test"

  # Message preferences
  tone: "professional"  # zen, funny, pirate, minimal
  language: "en"
  maxLength: 72
  includeBody: true
  includeEmoji: true

template:
  header: "{{emoji}} {{type}}({{scope}}): {{summary}}"
  body: |
    {{why}}

    {{impact}}
  footer: |
    Related: {{issues}}
```

**JSON format:**
```json
{
  "ai": {
    "provider": "anthropic",
    "model": "claude-sonnet-4",
    "temperature": 0.7
  },
  "commit": {
    "types": ["feat", "fix", "docs", "style", "refactor", "test", "chore"],
    "scopes": {
      "src/api/*": "api",
      "src/db/*": "database",
      "docs/*": "docs"
    },
    "tone": "professional",
    "maxLength": 72,
    "includeBody": true,
    "includeEmoji": true
  }
}
```

### Example 2: Template Engine

**Custom template with placeholders:**
```yaml
template:
  header: "{{type}}({{scope}}): {{summary}}"
  body: |
    ## Why this change?
    {{why}}

    ## Impact
    {{impact}}

    ## Technical details
    {{technical}}
  footer: |
    Closes: {{issues}}
    Co-authored-by: {{coauthors}}
```

**Generated output:**
```
feat(api): add user authentication endpoint

## Why this change?
Users need to authenticate before accessing protected resources.
This implements JWT-based authentication to secure our API.

## Impact
- New /api/auth/login endpoint
- Protected routes now require valid JWT token
- User sessions managed with 24hr expiration

## Technical details
- Uses bcrypt for password hashing
- JWT tokens signed with RS256 algorithm
- Redis cache for token blacklisting

Closes: #123, #145
Co-authored-by: Jane Doe <jane@example.com>
```

### Example 3: CLI Config Commands

```bash
# View current config
$ gocommit config show
Current configuration:
  AI Provider: anthropic
  Tone: professional
  Max Length: 72
  Include Body: true

# Set individual values
$ gocommit config set tone=zen
✓ Config updated: tone = zen

$ gocommit config set maxLength=50
✓ Config updated: maxLength = 50

# Edit config in default editor
$ gocommit config edit
# Opens .gocommitrc in $EDITOR

# Reset to defaults
$ gocommit config reset
? Are you sure you want to reset to default config? (y/N) y
✓ Configuration reset to defaults

# Show config file location
$ gocommit config path
Global: ~/.config/gocommit/.gocommitrc
Local:  /path/to/project/.gocommitrc (active)
```

### Example 4: Repository Override

**Global config (~/.config/gocommit/.gocommitrc):**
```yaml
commit:
  tone: "professional"
  includeEmoji: false
```

**Project-specific config (./myproject/.gocommitrc):**
```yaml
commit:
  tone: "funny"
  includeEmoji: true
  # This overrides global config for this repo only
```

---

## 🚀 Phase 1 — Interactive Mode

### Example 5: Interactive Prompt Flow

```bash
$ gocommit

🔍 Analyzing changes...
   • src/api/auth.ts (+45, -12)
   • src/types/user.ts (+8, -0)
   • tests/auth.test.ts (+67, -0)

🤖 Generated commit message:

┌────────────────────────────────────────────────────┐
│ feat(api): implement JWT authentication system     │
│                                                    │
│ Add secure authentication endpoints with JWT      │
│ token generation and validation. Includes user    │
│ login, logout, and token refresh functionality.   │
│                                                    │
│ - bcrypt password hashing                         │
│ - Redis session management                        │
│ - Comprehensive test coverage                     │
└────────────────────────────────────────────────────┘

What would you like to do?
[A] Accept  [E] Edit  [R] Regenerate  [S] Shorter  [D] Show diff  [Q] Quit

> r

🔄 Regenerating with different approach...

┌────────────────────────────────────────────────────┐
│ feat(api): add JWT-based user authentication      │
│                                                    │
│ Implement secure login/logout with JWT tokens,    │
│ bcrypt hashing, and Redis session store.          │
└────────────────────────────────────────────────────┘

What would you like to do?
[A] Accept  [E] Edit  [R] Regenerate  [S] Shorter  [D] Show diff  [Q] Quit

> a

✓ Committed successfully!
  SHA: a3f8b2c
  Files: 3 changed, 120 insertions(+), 12 deletions(-)
```

### Example 6: Side-by-Side Preview

```bash
$ gocommit --preview

╔═══════════════════════ DIFF ═══════════════════════╗   ╔══════════ MESSAGE ══════════╗
║ src/api/auth.ts                                    ║   ║ feat(api): add user auth    ║
║ @@ -1,5 +1,18 @@                                   ║   ║                             ║
║ +import jwt from 'jsonwebtoken';                   ║   ║ Implement JWT authentication║
║ +import bcrypt from 'bcrypt';                      ║   ║ with secure token handling  ║
║ +                                                  ║   ║ and session management.     ║
║ +export async function login(                     ║   ║                             ║
║ +  username: string,                              ║   ║ Changes:                    ║
║ +  password: string                               ║   ║ - JWT token generation      ║
║ +): Promise<AuthResponse> {                       ║   ║ - bcrypt password hashing   ║
║ +  const user = await findUser(username);         ║   ║ - Redis session store       ║
╚════════════════════════════════════════════════════╝   ╚═════════════════════════════╝
```

### Example 7: Auto-commit Flag

```bash
# Skip interactive mode and commit directly
$ gocommit --auto
🤖 Generating commit message...
✓ Committed: feat(api): add user authentication
  SHA: a3f8b2c

# Alias for quick commits
$ gocommit -y
✓ Committed: fix(db): resolve connection timeout issue
  SHA: b4e9d1a
```

---

## 🌈 Phase 2 — Personality & Delight

### Example 8: Tone Modes

**Professional mode (default):**
```bash
$ gocommit --mode=professional
feat(api): implement user authentication system

Add JWT-based authentication with secure token handling
and session management capabilities.
```

**Zen mode:**
```bash
$ gocommit --mode=zen
🧘 feat(api): add authentication flow

Like a river finding its path, authentication now flows
through our API. JWT tokens carry the essence of identity,
while bcrypt guards the sacred passwords.

In simplicity lies security.
```

**Funny mode:**
```bash
$ gocommit --mode=funny
😄 feat(api): teach the API who's boss (authentication)

Finally, the API learned to check IDs at the door! 🚪
JWT tokens are now the VIP passes, and bcrypt is our
tough bouncer who never forgets a password hash.

No more freeloaders in this API club! 🎉
```

**Pirate mode:**
```bash
$ gocommit --mode=pirate
⚓ feat(api): hoist the authentication sails!

Arrr! We've added JWT treasure maps fer validatin'
scallywags tryin' to board our API ship. Bcrypt be
guardin' the passwords like Davy Jones' locker!

No landlubber gets past without proper credentials! 🏴‍☠️
```

**Minimal mode:**
```bash
$ gocommit --mode=minimal
feat(api): add auth

JWT + bcrypt authentication.
```

### Example 9: Easter Eggs & Humor

```bash
$ gocommit

🎲 Rolling the AI dice...
✓ Committed: feat(ui): add dark mode toggle

🔮 Your commit fortune:
   "A well-crafted commit message is worth a thousand debugs"

───────────────────────────────────────────
   _____ ____  __  __ __  __ _____ _______
  / ____/ __ \|  \/  |  \/  |_   _|__   __|
 | |   | |  | | \  / | \  / | | |    | |
 | |   | |  | | |\/| | |\/| | | |    | |
 | |___| |__| | |  | | |  | |_| |_   | |
  \_____\____/|_|  |_|_|  |_|_____|  |_|

───────────────────────────────────────────
```

**Milestone celebrations:**
```bash
$ gocommit

✓ Committed: chore(deps): update dependencies

🎊 Achievement Unlocked! 🎊
You've made your 100th commit with GoCommit!

      *    *    *
   *    \   |   /    *
  *      \  |  /      *
 *        \ | /        *
  *        \|/        *
   *       /|\       *
    *     / | \     *
         /  |  \
        /   |   \

🏆 Title earned: "Commit Centurion"
```

---

## 📊 Phase 3 — Reflection & Feedback

### Example 10: Analytics Dashboard

```bash
$ gocommit stats

╔══════════════════════ COMMIT STATISTICS ═══════════════════════╗
║                                                                 ║
║  📊 Summary (Last 30 days)                                     ║
║  ────────────────────────────────────────────────────────────  ║
║  Total commits:        127                                      ║
║  Average per day:      4.2                                      ║
║  Longest streak:       12 days                                  ║
║  Current streak:       5 days 🔥                                ║
║                                                                 ║
║  📝 By Type                                                     ║
║  ────────────────────────────────────────────────────────────  ║
║  feat        ████████████████████░░░░░░ 42 (33%)              ║
║  fix         ████████████████░░░░░░░░░░ 35 (28%)              ║
║  chore       ██████████░░░░░░░░░░░░░░░░ 20 (16%)              ║
║  docs        ████████░░░░░░░░░░░░░░░░░░ 15 (12%)              ║
║  refactor    ██████░░░░░░░░░░░░░░░░░░░░ 10 (8%)               ║
║  test        ████░░░░░░░░░░░░░░░░░░░░░░  5 (4%)               ║
║                                                                 ║
║  ⏱️  Commit Timing                                              ║
║  ────────────────────────────────────────────────────────────  ║
║  Morning   (6-12):   ████████░░░░ 23%                          ║
║  Afternoon (12-18):  ████████████████░░░░ 45%                  ║
║  Evening   (18-24):  ██████████░░ 28%                          ║
║  Night     (0-6):    ██░░ 4%                                   ║
║                                                                 ║
║  📏 Message Quality                                             ║
║  ────────────────────────────────────────────────────────────  ║
║  Avg message length: 58 characters                             ║
║  With body:          45%                                        ║
║  With scope:         89%                                        ║
║  With emoji:         67%                                        ║
║                                                                 ║
║  🎯 Most Active Scopes                                          ║
║  ────────────────────────────────────────────────────────────  ║
║  1. api          32 commits                                    ║
║  2. ui           28 commits                                    ║
║  3. database     19 commits                                    ║
║  4. auth         15 commits                                    ║
║  5. docs         12 commits                                    ║
║                                                                 ║
╚═════════════════════════════════════════════════════════════════╝

💡 Tip: You're committing frequently! Consider breaking down large
   changes into smaller, more focused commits.
```

### Example 11: Export Analytics

```bash
$ gocommit stats --export=json > stats.json
✓ Exported statistics to stats.json

$ gocommit stats --export=csv
✓ Exported statistics to gocommit-stats-2025-11-01.csv

$ gocommit stats --range="last-year"
📊 Generating annual report...
✓ Stats for 2024: 1,247 commits across 365 days
```

### Example 12: Gamification & Badges

```bash
$ gocommit badges

╔═══════════════════════ YOUR ACHIEVEMENTS ══════════════════════╗
║                                                                 ║
║  🏅 Unlocked Badges (12/25)                                    ║
║                                                                 ║
║  ✅ 🌱 First Steps                                             ║
║     Made your first commit with GoCommit                       ║
║                                                                 ║
║  ✅ 🔥 On Fire                                                 ║
║     5-day commit streak                                        ║
║                                                                 ║
║  ✅ 💯 Centurion                                               ║
║     100 commits                                                ║
║                                                                 ║
║  ✅ 🐛 Bug Squasher                                            ║
║     Fixed 50 bugs                                              ║
║                                                                 ║
║  ✅ 🎨 Style Master                                            ║
║     10 style/formatting commits                                ║
║                                                                 ║
║  🔒 🚀 Feature Factory          [Progress: 73/100]             ║
║     Create 100 features                                        ║
║                                                                 ║
║  🔒 📚 Documentation Hero       [Progress: 12/50]              ║
║     Write 50 documentation commits                             ║
║                                                                 ║
║  🔒 🎯 Perfectionist           [Progress: 45/100]              ║
║     100 commits with detailed body messages                    ║
║                                                                 ║
║  🔒 🌙 Night Owl               [Not yet achieved]              ║
║     Make 25 commits between midnight and 6am                   ║
║                                                                 ║
║  🔒 ⚡ Speed Demon             [Not yet achieved]              ║
║     10 commits in one day                                      ║
║                                                                 ║
╚═════════════════════════════════════════════════════════════════╝

Next milestone: 27 more commits to unlock "Feature Factory" 🚀
```

---

## ☯ Phase 4 — Philosophy Integration

### Example 13: Reflection Prompt

```bash
$ gocommit

🧘 Take a moment to reflect...

  Before we commit, consider:

  💭 "Is this change meaningful?"
  💭 "Does it serve a clear purpose?"
  💭 "Will your future self understand why?"

  ╔═══════════════════════════════════════════════╗
  ║  "Code is like humor. When you have to       ║
  ║   explain it, it's bad."                     ║
  ║                    — Cory House              ║
  ╚═══════════════════════════════════════════════╝

[Press Enter to continue, or Ctrl+C to reconsider]

🤖 Generating commit message...
```

### Example 14: Mindfulness Mode

```bash
$ gocommit --pause

🌸 Mindfulness Mode Activated

  Take three deep breaths...

  Breathe in... ○ ○ ○ ○ ○  [1/3]

  [5 seconds]

  Breathe out... ● ● ● ● ●  [1/3]

  [Continue for 3 cycles]

  🧘 Now, let's commit with intention.

🤖 Analyzing changes...
```

### Example 15: Ethical Reminder

```bash
$ gocommit

✨ Remember: Clarity is kindness ✨

  A good commit message:
  ✓ Explains the "why", not just the "what"
  ✓ Helps others understand your reasoning
  ✓ Shows respect for future maintainers
  ✓ Builds trust through transparency

🤖 Generating commit message...

feat(api): add rate limiting to prevent abuse

Protect our API from excessive requests by implementing
token-bucket rate limiting. This ensures fair usage and
maintains service quality for all users.

- 100 requests per minute per IP
- Configurable limits per endpoint
- Clear error messages for rate-limited requests

[A] Accept  [E] Edit  [R] Regenerate  [Q] Quit
```

### Example 16: Philosophical Quotes Collection

```bash
$ gocommit

🎲 Developer Wisdom:

  "Programs must be written for people to read,
   and only incidentally for machines to execute."
                             — Harold Abelson

$ gocommit

🎲 Developer Wisdom:

  "Simplicity is prerequisite for reliability."
                             — Edsger Dijkstra

$ gocommit --mode=zen

🎲 Zen Wisdom:

  "The quieter you become, the more you can hear.
   The cleaner your code, the more it can speak."
```

---

## 🔧 Advanced Usage Examples

### Example 17: Custom AI Provider

```yaml
# .gocommitrc
ai:
  provider: "custom"
  endpoint: "http://localhost:11434/api/generate"  # Ollama
  model: "codellama"
  headers:
    Authorization: "Bearer ${CUSTOM_API_KEY}"
```

### Example 18: Multi-Language Support

```yaml
commit:
  language: "es"  # Spanish commits

# Output:
# feat(api): agregar autenticación de usuario
#
# Implementar autenticación basada en JWT con
# manejo seguro de tokens y gestión de sesiones.
```

### Example 19: Hook Integration

```bash
# .git/hooks/prepare-commit-msg
#!/bin/bash
# Auto-generate commit message if none provided
if [ -z "$(cat $1)" ]; then
  gocommit --auto --output=$1
fi
```

### Example 20: CI/CD Integration

```yaml
# .github/workflows/commit-check.yml
name: Validate Commit Messages
on: [push, pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install GoCommit
        run: curl -sSfL https://gocommit.sh/install | sh
      - name: Validate commits
        run: gocommit validate --strict
```

---

## 🔍 Handling Weak Models & Improving Quality

### Example 21: Model Quality Tiers

**Understanding different model capabilities:**

```yaml
# .gocommitrc - High Quality (Expensive)
ai:
  provider: "anthropic"
  model: "claude-sonnet-4"
  temperature: 0.7

# Output quality: ⭐⭐⭐⭐⭐
# - Deep understanding of code context
# - Meaningful commit messages
# - Proper conventional commit format
# Cost: $$ - Higher API costs
```

```yaml
# .gocommitrc - Medium Quality (Balanced)
ai:
  provider: "openai"
  model: "gpt-4o-mini"
  temperature: 0.5

# Output quality: ⭐⭐⭐⭐
# - Good understanding of changes
# - Consistent formatting
# - May miss subtle context
# Cost: $ - Moderate API costs
```

```yaml
# .gocommitrc - Low Quality (Cheap/Local)
ai:
  provider: "local"
  endpoint: "http://localhost:11434"
  model: "llama3:8b"
  temperature: 0.3

# Output quality: ⭐⭐⭐
# - Basic commit message generation
# - May need manual editing
# - Generic descriptions
# Cost: Free - Run locally
```

### Example 22: Quality Validation & Auto-Retry

**Implement quality checks before accepting commits:**

```bash
$ gocommit --validate

🔍 Analyzing changes...
🤖 Generating commit message...

Generated message:
┌────────────────────────────────────────────┐
│ feat: update files                         │
└────────────────────────────────────────────┘

❌ Quality Check Failed:
  - Message too generic (score: 2/10)
  - Missing scope
  - No context in description

🔄 Auto-retry with enhanced prompt...

Generated message (attempt 2/3):
┌────────────────────────────────────────────┐
│ feat(api): add user authentication system │
│                                            │
│ Implement JWT-based authentication with   │
│ secure token handling and validation.     │
└────────────────────────────────────────────┘

✅ Quality Check Passed (score: 8/10)
```

### Example 23: Constraining Weak Models with Templates

**Use strict templates to guide weaker models:**

```yaml
# .gocommitrc - For weak models
ai:
  provider: "local"
  model: "codellama:7b"

commit:
  # Enforce strict structure
  enforce_conventional: true
  require_scope: true
  require_description: true
  min_description_length: 20

  # Provide clear constraints
  types: ["feat", "fix", "docs", "chore"]  # Limited options

template:
  # Force specific format
  header: "{{type}}({{scope}}): {{summary}}"
  rules:
    - "Summary must be < 50 chars"
    - "Use imperative mood (add, not added)"
    - "No period at end of summary"
    - "Scope must match changed files"

# Example prompt enhancement for weak models
prompt_template: |
  Analyze these changes and generate a commit message.

  STRICT RULES:
  1. Format: type(scope): summary
  2. Type must be one of: feat, fix, docs, chore
  3. Scope must match the file path
  4. Summary: imperative mood, < 50 chars
  5. Be specific, not generic

  Changed files: {{files}}
  Diff summary: {{diff}}

  Good example: "feat(api): add JWT authentication"
  Bad example: "update files" ❌
```

### Example 24: Manual Fallback for Poor Results

**When AI fails, provide easy manual editing:**

```bash
$ gocommit

🤖 Generating commit message...

Generated message:
┌────────────────────────────────────────────┐
│ chore: update stuff                        │
└────────────────────────────────────────────┘

⚠️  Low confidence: This message seems generic

What would you like to do?
[A] Accept anyway
[E] Edit manually
[R] Regenerate with stricter prompt
[T] Use template wizard
[M] Switch to manual mode
[Q] Quit

> t

📝 Template Wizard Started

Step 1/4: What type of change is this?
  1. feat     - New feature
  2. fix      - Bug fix
  3. docs     - Documentation
  4. refactor - Code restructuring
  5. chore    - Maintenance
> 1

Step 2/4: What scope (component/area)?
  Detected scopes from changed files:
  1. api      (src/api/auth.ts)
  2. database (src/db/users.ts)
> 1

Step 3/4: Brief summary (imperative mood):
> add JWT authentication

Step 4/4: Would you like to add a body? (y/N)
> y

Why did you make this change?
> Users need secure authentication to access protected endpoints

Generated message:
┌────────────────────────────────────────────┐
│ feat(api): add JWT authentication          │
│                                            │
│ Users need secure authentication to       │
│ access protected endpoints                │
└────────────────────────────────────────────┘

✅ Quality score: 9/10
[A] Accept  [E] Edit  [R] Retry wizard  [Q] Quit
```

### Example 25: Progressive Enhancement Strategy

**Start simple, add detail incrementally:**

```yaml
# Strategy for weak models: Layer improvements
ai:
  provider: "local"
  model: "weak-model"

commit:
  strategy: "progressive"

  # Step 1: Get basic structure
  pass_1:
    prompt: "Generate commit type and scope only"
    example: "feat(api)"

  # Step 2: Add summary
  pass_2:
    prompt: "Add brief summary based on diff"
    example: "feat(api): add user authentication"

  # Step 3: Add body (optional)
  pass_3:
    prompt: "Add why/impact if significant changes"
    enabled: ${ENABLE_BODY}
```

**In practice:**
```bash
$ gocommit --progressive

Pass 1/3: Determining commit type...
✓ Type: feat, Scope: api

Pass 2/3: Generating summary...
✓ Summary: add user authentication

Pass 3/3: Adding context...
✓ Body: Implement JWT-based auth for secure access

Final result:
┌────────────────────────────────────────────┐
│ feat(api): add user authentication         │
│                                            │
│ Implement JWT-based auth for secure       │
│ access to protected endpoints.            │
└────────────────────────────────────────────┘
```

### Example 26: Cost vs Quality Optimization

**Smart model selection based on change complexity:**

```yaml
# .gocommitrc - Adaptive model selection
ai:
  strategy: "adaptive"

  # Use cheap model for simple changes
  simple_changes:
    model: "gpt-4o-mini"
    triggers:
      - files_changed: "<= 3"
      - lines_changed: "<= 50"
      - file_types: ["*.md", "*.txt"]
    cost_per_commit: "$0.001"

  # Use better model for complex changes
  complex_changes:
    model: "claude-sonnet-4"
    triggers:
      - files_changed: "> 3"
      - lines_changed: "> 50"
      - file_types: ["*.go", "*.ts", "*.rs"]
      - has_tests: true
    cost_per_commit: "$0.01"

  # Use local model for docs only
  documentation:
    model: "local/llama3"
    triggers:
      - all_files_match: "docs/*"
    cost_per_commit: "$0"
```

**Usage:**
```bash
$ gocommit

🔍 Analyzing changes...
  - 1 file changed (README.md)
  - 15 lines added
  - Change type: documentation

💡 Using gpt-4o-mini (cost-optimized for simple docs)
✓ Generated commit in 2s
💰 Cost: $0.0008

$ gocommit

🔍 Analyzing changes...
  - 8 files changed (auth.go, db.go, ...)
  - 347 lines added
  - Change type: major feature

💡 Using claude-sonnet-4 (quality-optimized for complex changes)
✓ Generated commit in 5s
💰 Cost: $0.012
```

### Example 27: Learning from Good Commits

**Train weak models with your repository's style:**

```bash
# Build a style guide from your best commits
$ gocommit learn --from-history

🔍 Analyzing last 100 commits...
📊 Detected patterns:
  - Preferred type: feat (35%), fix (28%), chore (20%)
  - Average summary length: 48 chars
  - Scope usage: 89% of commits
  - Body inclusion: 45% of commits
  - Common phrases:
    • "Implement" (42 times)
    • "Fix issue with" (28 times)
    • "Update" (31 times)

✅ Created .gocommit-style.json

# Use learned style with weak model
$ gocommit --style=learned

🤖 Using learned style guide to improve results...
✓ Commit generated with 92% style match
```

### Example 28: Human-in-the-Loop for Weak Models

**Always review with weak models:**

```yaml
# .gocommitrc - Safety for weak models
ai:
  provider: "local"
  model: "tinyllama:1b"  # Very weak model

commit:
  # Force human review
  auto_commit: false
  always_preview: true
  require_confirmation: true

  # Show quality indicators
  show_confidence_score: true
  highlight_issues: true
  suggest_improvements: true
```

**Example interaction:**
```bash
$ gocommit

🤖 Generating with TinyLlama:1b...

Generated message:
┌────────────────────────────────────────────┐
│ update code                                │
└────────────────────────────────────────────┘

📊 Quality Analysis:
  Confidence: ⭐⭐ (2/5)

  ⚠️  Issues detected:
  - Missing commit type (feat, fix, etc.)
  - Missing scope
  - Too generic ("update code")
  - No description of what changed

  💡 Suggestions:
  - Add type based on changes (looks like "feat")
  - Add scope from changed files (api, ui, db?)
  - Be specific: what feature was added?

Would you like to:
[E] Edit and fix issues manually
[R] Regenerate with better model
[T] Use template wizard
[Q] Quit and try again later

> e

Manual editing mode:
Type (feat/fix/docs/chore): feat
Scope: api
Summary: add JWT authentication endpoints
Add body? (y/n): y
Body:
> Implement login, logout, and token refresh endpoints
> with bcrypt password hashing and Redis session storage

Final message:
┌────────────────────────────────────────────┐
│ feat(api): add JWT authentication endpoints│
│                                            │
│ Implement login, logout, and token refresh│
│ endpoints with bcrypt password hashing and │
│ Redis session storage                      │
└────────────────────────────────────────────┘

✅ Commit? (y/n): y
```

### Example 29: Hybrid Approach - Local + Cloud

**Use local model first, upgrade if needed:**

```yaml
# .gocommitrc - Hybrid strategy
ai:
  strategy: "hybrid"

  # Try local first (free)
  primary:
    provider: "local"
    model: "codellama:7b"
    timeout: 10s

  # Fallback to cloud if quality is low
  fallback:
    provider: "anthropic"
    model: "claude-sonnet-4"
    triggers:
      - confidence: "< 0.6"
      - quality_score: "< 5"
      - user_rejected: true
```

**Usage:**
```bash
$ gocommit

🔄 Trying local model first...
✓ Generated in 3s

Generated message (local):
┌────────────────────────────────────────────┐
│ feat: add auth                             │
└────────────────────────────────────────────┘

📊 Confidence: 45% (LOW)

💡 This looks generic. Upgrade to cloud model? (y/n): y

🔄 Upgrading to Claude Sonnet 4...
✓ Generated in 4s (cost: $0.008)

Generated message (cloud):
┌────────────────────────────────────────────┐
│ feat(api): implement JWT authentication    │
│                                            │
│ Add secure login/logout with JWT tokens,  │
│ bcrypt hashing, and Redis sessions.       │
└────────────────────────────────────────────┘

📊 Confidence: 95% (HIGH)

[A] Accept  [E] Edit  [Q] Quit
```

---

## 📝 Complete Workflow Example

```bash
# 1. Make changes to your code
$ vim src/api/auth.ts
$ vim tests/auth.test.ts

# 2. Check what changed
$ git status
$ git diff

# 3. Stage your changes
$ git add src/api/auth.ts tests/auth.test.ts

# 4. Generate commit with GoCommit
$ gocommit --mode=zen --preview

# 5. Review the generated message
# [Interactive prompt appears]

# 6. Accept or edit as needed
> a

# 7. Check your stats
$ gocommit stats

# 8. Celebrate milestone if achieved
🎊 Achievement Unlocked: "Bug Squasher"

# 9. Push your changes
$ git push origin main
```

---

**Happy Committing! 🚀**

> "Each commit is a dialogue between you and your future self."
