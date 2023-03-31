module.exports = {
  theme: "cosmos",
  title: "Mineplex Chain",
  locales: {
    "/": {
      lang: "en-US"
    },
  },
  base: process.env.VUEPRESS_BASE || "/",
  head: [
    ['link', { rel: "apple-touch-icon", sizes: "180x180", href: "/apple-touch-icon.png" }],
    ['link', { rel: "icon", type: "image/png", sizes: "32x32", href: "/favicon-32x32.png" }],
    ['link', { rel: "icon", type: "image/png", sizes: "16x16", href: "/favicon-16x16.png" }],
    ['link', { rel: "manifest", href: "/site.webmanifest" }],
    ['meta', { name: "msapplication-TileColor", content: "#2e3148" }],
    ['meta', { name: "theme-color", content: "#ffffff" }],
    ['link', { rel: "icon", type: "image/svg+xml", href: "/favicon-svg.svg" }],
    ['link', { rel: "apple-touch-icon-precomposed", href: "/apple-touch-icon-precomposed.png" }],
  ],
  themeConfig: {
    repo: "cosmos/cosmos-sdk",
    docsRepo: "cosmos/cosmos-sdk",
    docsBranch: "release/v0.46.x",
    docsDir: "docs",
    editLinks: true,
    label: "chain",
    versions: [
      {
        "label": "v1.0",
        "key": "v1.0"
      }
    ],
    topbar: {
      banner: false
    },
    sidebar: {
      auto: true,
      nav: [

      ]
    },
  },
  plugins: [
    // [
    //   "@vuepress/google-analytics",
    //   {
    //     ga: "UA-51029217-2"
    //   }
    // ],
    // [
    //   "sitemap",
    //   {
    //     hostname: "https://docs.cosmos.network"
    //   }
    // ]
  ]
};
