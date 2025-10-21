package ai

const promptTemplate = `You are an expert at writing conventional commit messages. Analyze the git diff and generate 3 high-quality commit messages.

## Conventional Commit Format
type(optional-scope): description

## Available Types
- feat: new feature for users
- fix: bug fix for users
- docs: documentation changes
- style: formatting, missing semicolons, etc (no code change)
- refactor: code change that neither fixes bug nor adds feature
- perf: code change that improves performance
- test: adding/updating tests
- chore: updating build tasks, package manager configs, etc

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end
4. Scope should reflect the actual changed component
5. Description should explain WHAT changed, not WHY
6. Be specific about the actual changes in the diff

## Examples
- feat(auth): add OAuth2 login support
- fix(api): handle null response in user endpoint
- refactor(utils): extract validation logic to separate module
- docs(readme): update installation instructions
- style(components): fix indentation in header component

Generate exactly 3 different commit messages in this EXACT JSON format:

{
  "messages": [
    "feat(scope): commit message 1",
    "fix(scope): commit message 2",
    "refactor(scope): commit message 3"
  ]
}

CRITICAL:
- Return ONLY the JSON object above
- Each message must be a complete string in the "messages" array
- Do NOT return an array of objects
- Do NOT add any text before or after the JSON
- The JSON must have exactly one key "messages" containing an array of 3 strings`

const detailedPromptTemplate = `You are an expert at writing conventional commit messages. Analyze the git diff and generate 3 high-quality, detailed commit messages.

## Conventional Commit Format
type(optional-scope): description

Detailed description explaining what was changed and why.

## Available Types
- feat: new feature for users
- fix: bug fix for users
- docs: documentation changes
- style: formatting, missing semicolons, etc (no code change)
- refactor: code change that neither fixes bug nor adds feature
- perf: code change that improves performance
- test: adding/updating tests
- chore: updating build tasks, package manager configs, etc

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Subject line: max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end of subject line
4. Scope should reflect the actual changed component
5. Subject should explain WHAT changed
6. Body should explain WHY and HOW in more detail
7. Be specific about the actual changes in the diff
8. Include technical details and context in the body

## Examples
feat(auth): add OAuth2 login support

Implements OAuth2 authentication flow with Google and GitHub providers.
Adds token validation, refresh mechanism, and user profile fetching.
Includes comprehensive error handling for failed authentication attempts.

fix(api): handle null response in user endpoint

Prevents application crash when user data is missing from database.
Adds null checks and default values for required user fields.
Improves error messaging for better debugging experience.

Generate exactly 3 different detailed commit messages with body text. Format them as:

feat: enhance commit message generation with JSON output

Refactors the commit message generation process to return responses in JSON format.
This change ensures a structured and parsable output, improving integration with other tools.
Updates prompt templates to explicitly request JSON formatted messages and removes parsing logic.

---

fix: correct JSON parsing in commit message generation

Addresses an issue where the JSON output from the AI model was not correctly parsed.
Improves JSON extraction from the response by handling potential code blocks.
Adds more robust error handling and logging for debugging JSON parsing failures.

---

chore: update dependencies and improve error handling

Updates the go.mod and go.sum files to include the latest dependencies.
Improves error handling throughout the application, providing more informative error messages.
Includes changes to gracefully handle API request failures.

Return only the 3 commit messages in the format shown above with no additional text.`

const iconPromptTemplate = `You are an expert at writing conventional commit messages with emoji icons. Analyze the git diff and generate 3 high-quality commit messages using emoji icons.

## Conventional Commit Format with Icons
emoji type(optional-scope): description

## Available Types with Icons
- ✨ feat: new feature for users
- 🐛 fix: bug fix for users
- 📖 docs: documentation changes
- 💄 style: formatting, missing semicolons, etc (no code change)
- 🛠 refactor: code change that neither fixes bug nor adds feature
- ⚡️ perf: code change that improves performance
- ✅ test: adding/updating tests
- 📦 build: changes that affect the build system or external dependencies
- ⚙️ ci: changes to CI configuration files and scripts
- 🚀 chore: other changes that don't modify src or test files
- 🗑 revert: reverts a previous commit
- 🤞 try: add untested to production
- 🎉 init: project init

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end
4. Scope should reflect the actual changed component
5. Description should explain WHAT changed, not WHY
6. Be specific about the actual changes in the diff
7. Always start with the appropriate emoji icon

## Examples
- ✨ feat(auth): add OAuth2 login support
- 🐛 fix(api): handle null response in user endpoint
- 🛠 refactor(utils): extract validation logic to separate module
- 📖 docs(readme): update installation instructions
- 💄 style(components): fix indentation in header component

Generate exactly 3 different commit messages in this EXACT JSON format:

{
  "messages": [
    "✨ feat(scope): commit message 1",
    "🐛 fix(scope): commit message 2",
    "🛠 refactor(scope): commit message 3"
  ]
}

CRITICAL:
- Return ONLY the JSON object above
- Each message must be a complete string starting with an emoji in the "messages" array
- Do NOT return an array of objects
- Do NOT add any text before or after the JSON
- The JSON must have exactly one key "messages" containing an array of 3 strings`

const iconDetailedPromptTemplate = `You are an expert at writing conventional commit messages with emoji icons. Analyze the git diff and generate 3 high-quality, detailed commit messages using emoji icons.

## Conventional Commit Format with Icons
emoji type(optional-scope): description

Detailed description explaining what was changed and why.

## Available Types with Icons
- ✨ feat: new feature for users
- 🐛 fix: bug fix for users
- 📖 docs: documentation changes
- 💄 style: formatting, missing semicolons, etc (no code change)
- 🛠 refactor: code change that neither fixes bug nor adds feature
- ⚡️ perf: code change that improves performance
- ✅ test: adding/updating tests
- 📦 build: changes that affect the build system or external dependencies
- ⚙️ ci: changes to CI configuration files and scripts
- 🚀 chore: other changes that don't modify src or test files
- 🗑 revert: reverts a previous commit
- 🤞 try: add untested to production
- 🎉 init: project init

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Subject line: max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end of subject line
4. Scope should reflect the actual changed component
5. Subject should explain WHAT changed
6. Body should explain WHY and HOW in more detail
7. Be specific about the actual changes in the diff
8. Include technical details and context in the body
9. Always start with the appropriate emoji icon

## Examples
✨ feat(auth): add OAuth2 login support

Implements OAuth2 authentication flow with Google and GitHub providers.
Adds token validation, refresh mechanism, and user profile fetching.
Includes comprehensive error handling for failed authentication attempts.

🐛 fix(api): handle null response in user endpoint

Prevents application crash when user data is missing from database.
Adds null checks and default values for required user fields.
Improves error messaging for better debugging experience.

Generate exactly 3 different detailed commit messages with body text in this EXACT JSON format:

{
  "messages": [
    "✨ feat: enhance commit message generation with JSON output\n\nRefactors the commit message generation process to return responses in JSON format.\nThis change ensures a structured and parsable output, improving integration with other tools.\nUpdates prompt templates to explicitly request JSON formatted messages and removes parsing logic.",
    "🐛 fix: correct JSON parsing in commit message generation\n\nAddresses an issue where the JSON output from the AI model was not correctly parsed.\nImproves JSON extraction from the response by handling potential code blocks.\nAdds more robust error handling and logging for debugging JSON parsing failures.",
    "🚀 chore: update dependencies and improve error handling\n\nUpdates the go.mod and go.sum files to include the latest dependencies.\nImproves error handling throughout the application, providing more informative error messages.\nIncludes changes to gracefully handle API request failures."
  ]
}

CRITICAL:
- Return ONLY the JSON object above
- Each message must be a complete string with emoji, title, blank line (\\n\\n), and body text
- Do NOT return an array of objects with separate fields
- Do NOT add any text before or after the JSON
- The JSON must have exactly one key "messages" containing an array of 3 strings
- Each string in the array contains the full commit message including body separated by \\n\\n`

// GetPromptTemplate returns the appropriate prompt template based on options
func GetPromptTemplate(detailed, useIcons bool) string {
	if useIcons {
		if detailed {
			return iconDetailedPromptTemplate
		}
		return iconPromptTemplate
	}
	if detailed {
		return detailedPromptTemplate
	}
	return promptTemplate
}
