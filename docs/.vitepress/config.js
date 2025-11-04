export default {
  title: 'GoCommit Documentation',
  description: 'AI-powered Git commit message generator',
  base: '/',

  // Ignore dead links to Go source files
  // These are IDE reference links (e.g., ./config/config.go:14) that don't work as web links
  ignoreDeadLinks: [
    // Ignore all .go file references with or without line numbers
    /\.go(:\d+)?$/,
  ],

  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/README' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Overview', link: '/README' },
          { text: 'Installation', link: '/INSTALL' }
        ]
      },
      {
        text: 'Planning & Development',
        items: [
          { text: 'Roadmap & Philosophy', link: '/todo/doto' },
          { text: 'Usage Examples', link: '/todo/EXAMPLES' }
        ]
      },
      {
        text: 'Architecture',
        items: [
          { text: 'Delayed Commit Architecture', link: '/DELAYED_COMMIT_ARCHITECTURE' },
          { text: 'Application Flow', link: '/FLOW' },
          { text: 'JSON Schema', link: '/JSON_SCHEMA_IMPLEMENTATION' },
          { text: 'Parsing Fix', link: '/PARSING_FIX_SUMMARY' }
        ]
      },
      {
        text: 'Refactoring',
        items: [
          { text: 'Summary', link: '/REFACTORING_SUMMARY' },
          { text: 'Design', link: '/REFACTORING_DESIGN' }
        ]
      },
      {
        text: 'Integration',
        items: [
          { text: 'GitHub Actions', link: '/README_GITHUB_ACTIONS' },
          { text: 'Delayed Commit Setup', link: '/README_DELAYED_COMMIT' },
          { text: 'Logging Guide', link: '/README_LOGGING' },
          { text: 'Logging Config', link: '/README_LOGGING_CONFIG' }
        ]
      }
    ]
  }
}