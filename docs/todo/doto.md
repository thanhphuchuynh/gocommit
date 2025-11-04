# 🧭 GoCommit Roadmap & Philosophy

> _“Every commit is a moment of reflection — a trace of intent, not just change.”_

This roadmap outlines the evolution of **GoCommit**, from a simple AI-powered commit message helper to a mindful, delightful developer companion.

---

## ⚙️ Phase 0 — Custom Templates & Config
### 🎯 Goal: Empower users with personalization and control
- [ ] **Config File (`.gocommitrc` / YAML / JSON)**
  - Define default commit types (`feat`, `fix`, `chore`, etc.)  
  - Map scopes by folder or pattern (`db/*` → `chore(db)`)  
  - Choose tone, language, message length, and AI backend.  
- [ ] **Template Engine**
  - Custom header/body/footer formats using placeholders:  
    ```yaml
    header: "{{type}}({{scope}}): {{summary}}"
    body: |
      Reason: {{why}}
      Impact: {{impact}}
    ```
  - Optional commit body generation with “why” and “intent” sections.  
- [ ] **Repository Overrides**
  - `.gocommitrc` in each repo overrides global config.  
- [ ] **CLI Shortcuts**
  - `gocommit config set tone=zen`  
  - `gocommit config edit` (opens default editor).

🧩 _“Before automation comes awareness — define your own meaning first.”_

---

## 🚀 Phase 1 — Interactive Mode ✅ COMPLETED
### 🎯 Goal: Improve usability & feedback loop
- [x] Interactive prompt after AI suggestion:
  - Accept / Edit / Re-generate (alternate tone or shorter version).
  - Preview diff + generated message side-by-side.
- [x] Optional "auto-commit" flag for fast workflows.

**Implementation Details:**
- Interactive prompt with 5 actions: Accept, Edit, Regenerate, Diff, Quit
- Split-screen diff view showing changes and message side-by-side
- Regeneration loop that generates → select → interactive → repeat
- Color-coded diff display (green +, red -, cyan @@)
- `--auto` and `-y` flags to skip interactive prompt
- Seamless integration in workflow with goto-based flow control

🧠 _"Before you commit, consider if your change has a place in the story."_

---

## 🔍 Phase 2 — Quality & Model Management
### 🎯 Goal: Handle diverse AI models and ensure quality
- [ ] **Quality Validation System**
  - Automatic quality scoring (1-10) for generated messages.
  - Detect generic messages: "update files", "fix stuff" → reject & retry.
  - Auto-retry with enhanced prompts if quality is low.
  - Validation rules: scope required, minimum description length, etc.
- [ ] **Multi-Tier Model Support**
  - **High-end models**: Claude Sonnet, GPT-4 (⭐⭐⭐⭐⭐ quality, $$ cost)
  - **Mid-tier models**: GPT-4o-mini, Claude Haiku (⭐⭐⭐⭐ quality, $ cost)
  - **Local models**: Ollama, CodeLlama (⭐⭐⭐ quality, free)
  - Configure per project or per commit type.
- [ ] **Adaptive Model Selection**
  - Use cheap models for simple changes (docs, README).
  - Use expensive models for complex code changes.
  - Triggers: file count, line changes, file types, presence of tests.
  - Show cost estimation before committing.
- [ ] **Weak Model Enhancement**
  - **Strict templates**: Force format for weaker models.
  - **Progressive generation**: Generate in steps (type → scope → summary → body).
  - **Template wizard**: Interactive fallback when AI fails.
  - **Learned style guide**: Analyze repo history to learn commit patterns.
- [ ] **Hybrid Strategy**
  - Try local model first (free, fast).
  - Auto-upgrade to cloud model if confidence < 60%.
  - User can manually trigger upgrade.
  - Show confidence score and quality indicators.
- [ ] **Human-in-the-Loop**
  - Always show preview for weak models.
  - Highlight potential issues (missing scope, too generic).
  - Suggest improvements inline.
  - Never auto-commit with low confidence scores.

🎯 _"Quality is not an accident; it is always the result of intelligent effort." — John Ruskin_

---

## 🌈 Phase 3 — Personality & Delight
### 🎯 Goal: Make commits fun and human
- [ ] **Tone Modes**
  - `--mode=zen`, `--mode=funny`, `--mode=pirate`, `--mode=minimal`.  
  - Each mode generates commits in its own expressive style.  
- [ ] **Humor & Easter Eggs**
  - Playful verbs or ASCII art feedback.  
  - Optional “fortune mode”: prints a haiku or quote after commit.  

🧘 _"Simplicity is the ultimate sophistication." — Leonardo da Vinci_

---

## 📊 Phase 4 — Reflection & Feedback
### 🎯 Goal: Help users see their commit patterns
- [ ] **Analytics Dashboard**
  - `gocommit stats` → shows commit counts, types, length, etc.  
  - Local summary or optional export for visualization.  
- [ ] **Gamification**
  - Fun badges:  
    - 🏅 “Commit Grandmaster” (100 commits)  
    - 🔥 “Fix Fiend” (50 fixes)  
  - Optional: ASCII fireworks for milestones.

💬 _"Measure what you commit — not for vanity, but for growth."_

---

## ☯ Phase 5 — Philosophy Integration
### 🎯 Goal: Embed mindfulness into the developer’s flow
- [ ] **Reflection Prompt**
  - Before commit: short pause — “Is this change meaningful?”  
  - Display random philosophical or developer quote.  
- [ ] **Mindfulness Mode**
  - `--pause`: a brief breathing space before confirming commit.  
- [ ] **Ethical Reminder**
  - Encourage clarity, honesty, and purpose in messages.

🪶 _"Code fades, intent remains."_

---

## 📖 Phase 6 — Documentation & Legacy
### 🎯 Goal: Tell the story behind the tool
- [ ] **README Update**
  - Add **“Why We Wrote This”** section — the philosophy behind GoCommit.  
- [ ] **Dev.to Blog Post**
  - “GoCommit: Writing Commit Messages Like a Philosopher.”  
  - Share reasoning, architecture, and lessons learned.

🪷 _“We don’t just automate the message — we elevate the meaning.”_

---

## 🛠️ Summary of Priorities
| Priority | Feature | Type | Status |
|-----------|----------|------|--------|
| 1 | Custom templates/config | Personalization | 🏁 In progress |
| 2 | Interactive mode | Core UX | ⏳ Planned |
| 3 | Quality & model management | Reliability | 🔥 Critical |
| 4 | Tone modes (funny/zen/pirate) | Delight | 🔜 Next |
| 5 | Analytics & gamification | Reflection | 🧩 Later |
| 6 | Philosophical layer | Spirit | 🌙 Optional |

### 🎯 Implementation Focus Areas

**Phase 0-2 (MVP - Core Functionality)**
- Config system for customization
- Interactive commit workflow
- Quality validation & multi-model support
- Essential for production use

**Phase 3-4 (Enhanced Experience)**
- Personality modes & delight
- Analytics & reflection
- User engagement features

**Phase 5-6 (Polish & Community)**
- Philosophical integration
- Documentation & outreach
- Long-term vision

---

**Created by:** [@thanhphuchuynh](https://github.com/thanhphuchuynh)  
**Project:** [GoCommit](https://github.com/thanhphuchuynh/gocommit)  
**Last updated:** 2025-11-01  

> “Each commit is a dialogue between you and your future self.”
