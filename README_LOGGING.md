# Logging Feature for gocommit

## Overview
The gocommit tool now includes comprehensive logging functionality to track all requests and responses with the AI service. This data can be used to analyze and improve prompts over time.

## What Gets Logged

Each interaction with the AI service logs the following information:

- **Timestamp**: When the request was made
- **Git Diff**: The staged changes being committed
- **Last Commit Message**: The previous commit message used as reference
- **Prompt Sent**: The complete prompt sent to the AI
- **AI Response**: The raw response from the AI service
- **Generated Messages**: The 3 cleaned commit message options
- **Selected Message**: The final commit message chosen by the user
- **Success/Error Status**: Whether the request succeeded or failed
- **Error Details**: If applicable, what went wrong

## Log File Location

Logs are stored in JSON format at:
```
~/.gocommit/gocommit_requests.log
```

Each log entry is a single line of JSON, making it easy to parse programmatically.

## Example Log Entry

```json
{
  "timestamp": "2025-05-29T21:23:00Z",
  "git_diff": "diff --git a/main.go b/main.go\nindex 1234567..abcdefg 100644\n--- a/main.go\n+++ b/main.go\n@@ -1,3 +1,4 @@\n package main\n \n import \"fmt\"\n+import \"log\"\n",
  "last_commit_msg": "feat: add initial CLI structure",
  "prompt_sent": "As an AI commit message generator, analyze the following git diff...",
  "ai_response": "1. feat(import): add log import for debugging\n2. chore(deps): import log package\n3. style(imports): add log import",
  "generated_messages": [
    "feat(import): add log import for debugging",
    "chore(deps): import log package", 
    "style(imports): add log import"
  ],
  "selected_message": "feat(import): add log import for debugging",
  "success": true
}
```

## Usage

The logging happens automatically whenever you use gocommit. After each commit, you'll see a message showing where the request was logged:

```
Successfully created commit with message: feat: add logging functionality
Request logged to: /Users/username/.gocommit/gocommit_requests.log
```

## Benefits

1. **Prompt Analysis**: Review what prompts work best for different types of changes
2. **AI Response Quality**: Track the quality and consistency of AI-generated messages
3. **User Preferences**: See which types of commit messages users prefer to select
4. **Error Tracking**: Monitor any issues with the AI service
5. **Prompt Improvement**: Use the data to refine and improve the prompt template

## Privacy Note

The log file is stored locally on your machine and contains your git diffs and commit messages. Be aware of this if working with sensitive code.
