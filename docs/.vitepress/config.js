export default {
  title: 'GoCommit Documentation',
  description: 'AI-powered Git commit message generator',

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
        text: 'Documentation',
        items: [
          { text: 'Overview', link: '/README' },
          { text: 'Delayed Commit', link: '/README_DELAYED_COMMIT' },
          { text: 'Architecture', link: '/DELAYED_COMMIT_ARCHITECTURE' },
          { text: 'Refactoring', link: '/REFACTORING_SUMMARY' },
          { text: 'Design', link: '/REFACTORING_DESIGN' },
          { text: 'Parsing Fix', link: '/PARSING_FIX_SUMMARY' }
        ]
      }
    ]
  }
}