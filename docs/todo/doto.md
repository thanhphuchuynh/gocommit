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

## 🚀 Phase 1 — Interactive Mode
### 🎯 Goal: Improve usability & feedback loop
- [ ] Interactive prompt after AI suggestion:
  - Accept / Edit / Re-generate (alternate tone or shorter version).  
  - Preview diff + generated message side-by-side.  
- [ ] Optional “auto-commit” flag for fast workflows.

🧠 _“Before you commit, consider if your change has a place in the story.”_

---

## 🌈 Phase 2 — Personality & Delight
### 🎯 Goal: Make commits fun and human
- [ ] **Tone Modes**
  - `--mode=zen`, `--mode=funny`, `--mode=pirate`, `--mode=minimal`.  
  - Each mode generates commits in its own expressive style.  
- [ ] **Humor & Easter Eggs**
  - Playful verbs or ASCII art feedback.  
  - Optional “fortune mode”: prints a haiku or quote after commit.  

🧘 _“Simplicity is the ultimate sophistication.” — Leonardo da Vinci_

---

## 📊 Phase 3 — Reflection & Feedback
### 🎯 Goal: Help users see their commit patterns
- [ ] **Analytics Dashboard**
  - `gocommit stats` → shows commit counts, types, length, etc.  
  - Local summary or optional export for visualization.  
- [ ] **Gamification**
  - Fun badges:  
    - 🏅 “Commit Grandmaster” (100 commits)  
    - 🔥 “Fix Fiend” (50 fixes)  
  - Optional: ASCII fireworks for milestones.

💬 _“Measure what you commit — not for vanity, but for growth.”_

---

## ☯ Phase 4 — Philosophy Integration
### 🎯 Goal: Embed mindfulness into the developer’s flow
- [ ] **Reflection Prompt**
  - Before commit: short pause — “Is this change meaningful?”  
  - Display random philosophical or developer quote.  
- [ ] **Mindfulness Mode**
  - `--pause`: a brief breathing space before confirming commit.  
- [ ] **Ethical Reminder**
  - Encourage clarity, honesty, and purpose in messages.

🪶 _“Code fades, intent remains.”_

---

## 📖 Phase 5 — Documentation & Legacy
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
| 3 | Tone modes (funny/zen/pirate) | Delight | 🔜 Next |
| 4 | Analytics & gamification | Reflection | 🧩 Later |
| 5 | Philosophical layer | Spirit | 🌙 Optional |

---

**Created by:** [@thanhphuchuynh](https://github.com/thanhphuchuynh)  
**Project:** [GoCommit](https://github.com/thanhphuchuynh/gocommit)  
**Last updated:** 2025-11-01  

> “Each commit is a dialogue between you and your future self.”
