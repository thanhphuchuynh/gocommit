---
title: "Why Do You Push Code During Work Hours?" - How an Interview Question Led Me to Build a Delayed Commit Feature
published: false
description: An interview question about Git commit timestamps sparked an idea for a feature that helps developers maintain work-life boundaries while keeping honest contribution graphs.
tags: git, golang, productivity, devtools
cover_image: https://dev-to-uploads.s3.amazonaws.com/uploads/articles/random.jpg
canonical_url:
---

## The Interview That Changed My Perspective

Picture this: I'm sitting in a technical interview, feeling pretty confident about my answers so far. Then the interviewer glances at my GitHub profile and asks a seemingly innocent question:

> "I notice you push code mostly during work hours. Can you tell me about your development workflow?"

I froze.

Not because I was doing anything wrong, but because I suddenly realized how much our commit timestamps reveal about our work habits. In that moment, I understood that my GitHub contribution graph wasn't just showing my coding activity—it was inadvertently broadcasting my work schedule to the world.

This interview question got me thinking: **What if developers want more control over when their commits appear to be made?**

## The Problem We Don't Talk About

As developers, we live in a transparent world. Our GitHub profiles are our resumes, our contribution graphs are our performance metrics, and our commit histories are open books. But here's the uncomfortable truth:

- **Work-life boundaries are blurring**: That green square at 2 AM? It might just be because you finally fixed that bug, but it signals "always available" to some employers.
- **Timezone confusion**: Working from Asia but applying to European companies? Your commit times might look odd to recruiters.
- **Privacy concerns**: Not everyone wants their daily routine visible to the world.
- **Unconscious bias**: Some people judge productivity by when commits happen, not what they contain.

After that interview, I decided to do something about it. I built the **Delayed Commit** feature for [GoCommit](https://github.com/thanhphuchuynh/gocommit), my AI-powered commit message generator.

## What Is Delayed Commit?

The concept is simple but powerful: **prevent commits during your defined work hours and let you schedule them for later.**

Here's how it works:

1. You configure your "restricted hours" (e.g., 9 AM - 5 PM)
2. When you try to commit during those hours, GoCommit intercepts it
3. A friendly UI appears, suggesting alternative times
4. You pick a time (or enter a custom one)
5. Your commit is created with that timestamp

No scripts, no manual git commands, no fumbling with `--date` flags. Just a smooth, integrated experience.

## The Architecture Behind the Magic

Building this feature was fascinating. Here's what happens under the hood:

### 1. Time Validation
```go
func isTimeInRestrictedRange(currentTime time.Time, config DelayedCommitConfig) (bool, time.Time, error) {
    hour := currentTime.Hour()

    if hour >= config.RestrictedStartHour && hour < config.RestrictedEndHour {
        // Calculate when restrictions end
        endTime := time.Date(
            currentTime.Year(), currentTime.Month(), currentTime.Day(),
            config.RestrictedEndHour, 0, 0, 0, currentTime.Location(),
        )
        return true, endTime, nil
    }

    return false, time.Time{}, nil
}
```

### 2. Smart Time Suggestions
The system generates suggested commit times at configurable intervals (default: 20 minutes) starting from when your restrictions end:

```
┌──────────────────────────────────────────────────────────┐
│ Commit during work hours (09:00-17:00) detected         │
│ Current time: 14:30                                      │
│                                                          │
│ Select commit time:                                      │
│  → 17:20 (5:20 PM) - Today                              │
│  • 17:40 (5:40 PM) - Today                              │
│  • 18:00 (6:00 PM) - Today                              │
│  • 18:20 (6:20 PM) - Today                              │
│  • Enter custom time...                                  │
│                                                          │
│ ↑↓: Move  Enter: Select  Esc: Use current time anyway   │
└──────────────────────────────────────────────────────────┘
```

### 3. Git Integration
Git actually supports custom timestamps natively! The trick is using both the `--date` flag and the `GIT_COMMITTER_DATE` environment variable:

```go
func executeDelayedCommit(message string, timestamp time.Time) *exec.Cmd {
    dateStr := timestamp.Format(time.RFC3339)

    cmd := exec.Command("git", "commit", "-m", message, "--date", dateStr)
    cmd.Env = append(os.Environ(),
        fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))

    return cmd
}
```

This ensures both the **author date** (when the code was written) and **committer date** (when the commit was recorded) match your chosen time.

## Real-World Use Cases

Since implementing this feature, I've discovered several legitimate use cases:

### 1. Maintaining Work-Life Boundaries
Developers who want to separate personal projects from work hours without actually changing when they code.

### 2. Timezone Management
Remote workers who want their contributions to appear during "normal" hours in their employer's timezone.

### 3. Privacy Protection
Freelancers who don't want clients tracking their exact working hours.

### 4. Contribution Graph Aesthetics
Let's be honest—some people just want a nicer-looking contribution graph. And that's okay!

### 5. Retroactive Commits
Sometimes you work on code offline or forget to commit. Being able to set accurate timestamps helps maintain an honest history.

## The Ethics Question

Now, I know what you're thinking: *"Isn't this... dishonest?"*

Let's talk about it.

**Here's my take:** Git timestamps have always been flexible. You can manually set them, rebase changes, or use any number of tools to modify history. This feature doesn't enable anything new—it just makes an existing capability more accessible and user-friendly.

The real question is: **Should your coding schedule be public information?**

I argue it shouldn't. Your value as a developer comes from:
- The quality of your code
- Your problem-solving skills
- Your collaboration and communication
- The impact of your work

Not from whether you pushed at 2 PM or 2 AM.

This tool helps you control your narrative without lying. You're still writing the code, still making the commits—you're just choosing when they appear in your history.

## Getting Started

Want to try it? Here's how to set up GoCommit with delayed commits:

### 1. Install GoCommit
```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh | bash

# Or download from releases
# https://github.com/thanhphuchuynh/gocommit/releases
```

### 2. Configure Your API Key
GoCommit uses AI to generate commit messages, so you'll need an API key:
```bash
gocommit --config
```

Choose between Google Gemini or OpenRouter (which gives you access to Claude, GPT-4, Llama, and more).

### 3. Enable Delayed Commits
```bash
gocommit --config-delayed
```

You'll be prompted to set:
- Restricted start hour (e.g., 9 for 9 AM)
- Restricted end hour (e.g., 17 for 5 PM)
- Suggestion interval (e.g., 20 minutes)

### 4. Use It Naturally
```bash
git add .
gocommit
```

That's it! If you're committing during restricted hours, you'll see the time selection UI. Otherwise, it works exactly like before.

## Configuration File

Under the hood, it's just a simple JSON config at `~/.gocommit.json`:

```json
{
  "api_key": "your-api-key",
  "logging_enabled": true,
  "icon_mode": false,
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 9,
    "restricted_end_hour": 17,
    "suggestion_interval_min": 20
  }
}
```

## The Technical Challenges

Building this feature taught me a few things:

### Challenge 1: Time Zones Are Hard
Handling time zones correctly is always harder than you think. I ended up always using the local system timezone to avoid confusion.

### Challenge 2: UI in the Terminal
Creating an intuitive terminal UI with `termbox-go` was tricky. I studied the existing message selection UI in GoCommit and built something consistent with that pattern.

### Challenge 3: Edge Cases Everywhere
What if someone commits at 11:59 PM? What if they enter a time in the past? What if there are no available times today? Each edge case needed thoughtful handling.

### Challenge 4: Git Command Complexity
Git has two separate timestamps (author and committer), and they behave differently. Learning to set both properly took some trial and error.

## What's Next?

I'm considering adding:

- **Multi-day scheduling**: Schedule commits for future days, not just today
- **Smart suggestions**: Learn from your patterns and suggest times automatically
- **Batch processing**: Apply delayed commits to multiple staged changes at once
- **Template presets**: Quick configurations for common scenarios (remote work, freelancing, etc.)

## The Bigger Picture

This project started from an uncomfortable interview question, but it's really about something bigger: **developer autonomy**.

We should have control over what information we share and when. Our contribution graphs should reflect our **work**, not our **schedules**. And tools should empower us to make these choices easily.

Whether you use this feature or not, I hope it sparks a conversation about privacy, work-life balance, and what we expect from developers in our industry.

## Try It Yourself

The delayed commit feature is available now in GoCommit:
- **GitHub**: [github.com/thanhphuchuynh/gocommit](https://github.com/thanhphuchuynh/gocommit)
- **Docs**: [Full documentation](https://thanhphuchuynh.github.io/gocommit/)
- **Architecture**: [Technical deep-dive](https://github.com/thanhphuchuynh/gocommit/blob/main/docs/DELAYED_COMMIT_ARCHITECTURE.md)

And if you're curious about the code, it's all open source. Feel free to explore, contribute, or fork it for your own projects.

## Final Thoughts

That interview question—*"Why do you push code during work hours?"*—wasn't trying to catch me doing something wrong. But it made me realize that in our quest for transparency and collaboration, we might have given up more privacy than we intended.

This feature is my answer: a way to be honest about your work while maintaining control over your personal information.

What do you think? Would you use something like this? What other privacy concerns do you have as a developer? Drop your thoughts in the comments—I'd love to hear your perspective!

---

*P.S. If you found this useful, give [GoCommit](https://github.com/thanhphuchuynh/gocommit) a star on GitHub! It helps more developers discover the tool. Also, GoCommit generates your commit messages with AI, so you get both better messages AND schedule control in one tool.*

---

## About the Project

**GoCommit** is an AI-powered commit message generator that:
- Analyzes your staged changes automatically
- Generates meaningful, conventional commit messages
- Supports multiple AI providers (Gemini, Claude, GPT-4, Llama, etc.)
- Includes built-in logging for prompt improvement
- Now features delayed commits for schedule control

It's free, open source, and designed to make your Git workflow smoother. Check it out!
