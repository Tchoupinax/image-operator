<template>
  <div class="min-h-screen bg-slate-50 text-slate-900 antialiased">
    <div class="pointer-events-none fixed inset-0 overflow-hidden">
      <div class="absolute -left-32 top-0 h-96 w-96 rounded-full bg-sky-200/40 blur-3xl" />
      <div class="absolute -right-32 top-1/3 h-96 w-96 rounded-full bg-violet-200/40 blur-3xl" />
    </div>
    <NuxtPage />
  </div>
</template>

<script setup lang="ts">
useHead({
  link: [
    {
      rel: "preconnect",
      href: "https://fonts.googleapis.com",
    },
    {
      rel: "preconnect",
      href: "https://fonts.gstatic.com",
      crossorigin: "",
    },
    {
      rel: "stylesheet",
      href: "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap",
    },
  ],
});
</script>
