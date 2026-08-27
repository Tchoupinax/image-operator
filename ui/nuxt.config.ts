// https://nuxt.com/docs/api/configuration/nuxt-config
import { fileURLToPath } from "node:url";

const monorepoRoot = fileURLToPath(new URL("..", import.meta.url));

export default defineNuxtConfig({
  compatibilityDate: "2024-11-02",
  devtools: { enabled: false },
  css: ["~/assets/css/main.css"],
  modules: ["@nuxt/eslint", "@nuxtjs/tailwindcss"],
  runtimeConfig: {
    public: {
      graphqlApiUrl: "http://localhost:9090/graphql",
      operatorVersion: "v0.0.0"
    }
  },
  vite: {
    server: {
      fs: {
        allow: [monorepoRoot],
      },
    },
  },
  // Due to nuxt-build-cache
  workspaceDir: process.cwd(),
})