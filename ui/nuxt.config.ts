// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-11-02",
  devtools: { enabled: false },
  modules: ["@nuxt/eslint", "@nuxtjs/tailwindcss"],
  runtimeConfig: {
    public: {
      graphqlApiUrl: "http://localhost:9090/graphql",
    }
  },
  // Due to nuxt-build-cache
  workspaceDir: process.cwd(),
})