module.exports = {
  // Types of commits that trigger version bumps
  types: [
    { type: "feat", section: "🚀 Features" },
    { type: "fix", section: "🐛 Bug Fixes" },
    { type: "docs", section: "📄 Documentation" },
    { type: "style", section: "🎨 Style & Formatting" },
    { type: "refactor", section: "♻️ Enhancement" },
    { type: "perf", section: "⚡ Performance" },
    { type: "test", section: "🧪 Test" },
    { type: "build", section: "🧩 Dependency Updates" },
    { type: "ci", section: "👷 CI" },
    { type: "chore", section: "🧰 Maintenance" },
    { type: "revert", section: "⏪ Revert" },
    { type: "breaking", section: "⚠️ Breaking Changes" },
  ],

  // Commit scopes for our Go Food Delivery project
  releaseCommitMessageFormat: "chore(release): 📦 {{currentTag}}",

  // Issue prefixes for linking commits to issues
  issuePrefixes: ["#", "GH-", "Fixes #", "Closes #"],

  // Custom commit message format
  commitUrlFormat:
    "https://github.com/DavidReque/go-food-delivery/commit/{{hash}}",
  compareUrlFormat:
    "https://github.com/DavidReque/go-food-delivery/compare/{{previousTag}}...{{currentTag}}",
  issueUrlFormat:
    "https://github.com/DavidReque/go-food-delivery/issues/{{id}}",

  // Changelog configuration
  changelogFile: "CHANGELOG.md",
  changelogTitle:
    "# 📋 Changelog\n\nAll notable changes to this project will be documented in this file.\n\n",

  // Version bump rules
  releaseRules: [
    { type: "breaking", release: "major" },
    { type: "feat", release: "minor" },
    { type: "fix", release: "patch" },
    { type: "docs", release: "patch" },
    { type: "style", release: "patch" },
    { type: "refactor", release: "patch" },
    { type: "perf", release: "patch" },
    { type: "test", release: "patch" },
  ],

  // Pre-release configuration
  prerelease: false,

  // Git configuration
  gitTag: true,
  gitCommit: true,
  gitPush: false, // Set to true if you want automatic pushing

  // Skip bumping if no changes
  skip: {
    bump: false,
    changelog: false,
    commit: false,
    tag: false,
  },
};
