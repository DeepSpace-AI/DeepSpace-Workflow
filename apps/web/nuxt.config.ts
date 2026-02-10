// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  modules: [
    "@nuxt/eslint",
    "@nuxt/fonts",
    "@nuxt/ui",
    "@nuxtjs/mdc",
    "@pinia/nuxt",
  ],
  css: ["~/assets/css/main.css"],
  runtimeConfig: {
    aiGateway: {
      url: "",
      apiKey: "",
    },
  },
  routeRules: {
    "/": {
      appLayout: false,
    },
    "/sign-in": {
      appLayout: false,
    },
    "/sign-up": {
      appLayout: false,
    },
    "/privacy-policy": {
      appLayout: false,
    },
    "/terms-of-service": {
      appLayout: false,
    },
  },
  vite: {
    optimizeDeps: {
      include: [
        "@nuxt/ui > prosemirror-state",
        "@nuxt/ui > prosemirror-transform",
        "@nuxt/ui > prosemirror-model",
        "@nuxt/ui > prosemirror-view",
        "@nuxt/ui > prosemirror-gapcursor",
      ],
    },
  },
  app: {
    head: {
      title: "DeepSpaceWorkflow",
    },
  },
  nitro: {
    preset: "cloudflare-module",
    cloudflare: {
      deployConfig: true,
      nodeCompat: true,
    },
  },
});
