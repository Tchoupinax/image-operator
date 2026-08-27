<template>
  <div class="relative mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
    <header class="mb-10 flex flex-col gap-6 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <div class="mb-3 inline-flex items-center gap-2 rounded-full border border-sky-200 bg-sky-50 px-3 py-1 text-xs font-medium text-sky-700">
          <span class="h-1.5 w-1.5 rounded-full bg-sky-500" />
          Kubernetes operator
        </div>
        <h1 class="text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl">
          Image Operator
        </h1>
        <p class="mt-2 max-w-xl text-sm leading-relaxed text-slate-600">
          Monitor image sync jobs and builders across your cluster.
          <span v-if="version" class="font-mono text-slate-500">· {{ version }}</span>
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          class="btn-secondary"
          :disabled="refreshing"
          @click="refresh"
        >
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4" :class="{ 'animate-spin': refreshing }">
            <path fill-rule="evenodd" d="M15.312 11.424a5.5 5.5 0 0 1-9.201 2.466 5.5 5.5 0 0 1 9.201-2.466ZM15.312 11.424V8.5a.75.75 0 0 0-1.5 0v3.5a.75.75 0 0 0 .22.53l2.5 2.5a.75.75 0 1 0 1.06-1.06l-2.28-2.28Z" clip-rule="evenodd" />
          </svg>
          Refresh
        </button>
        <ModalCopyImage v-if="displayImages" @create="createImage" />
      </div>
    </header>

    <div class="mb-8 grid gap-4 sm:grid-cols-3">
      <div class="panel p-5">
        <p class="text-xs font-medium uppercase tracking-wider text-slate-500">Images</p>
        <p class="mt-2 text-3xl font-semibold text-slate-900">{{ images.length }}</p>
      </div>
      <div class="panel p-5">
        <p class="text-xs font-medium uppercase tracking-wider text-slate-500">Builders</p>
        <p class="mt-2 text-3xl font-semibold text-slate-900">{{ imageBuilders.length }}</p>
      </div>
      <div class="panel p-5">
        <p class="text-xs font-medium uppercase tracking-wider text-slate-500">Completed</p>
        <p class="mt-2 text-3xl font-semibold text-emerald-600">{{ completedCount }}</p>
      </div>
    </div>

    <div class="panel overflow-hidden">
      <div class="flex flex-col gap-4 border-b border-slate-200 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
        <div class="inline-flex rounded-xl border border-slate-200 bg-slate-100 p-1">
          <button
            class="tab-btn"
            :class="displayImages ? 'tab-btn-active' : 'tab-btn-inactive'"
            @click="displayImages = true"
          >
            Images
            <span class="ml-1.5 rounded-md bg-slate-200 px-1.5 py-0.5 text-xs text-slate-600">{{ images.length }}</span>
          </button>
          <button
            class="tab-btn"
            :class="!displayImages ? 'tab-btn-active' : 'tab-btn-inactive'"
            @click="displayImages = false"
          >
            Image Builders
            <span class="ml-1.5 rounded-md bg-slate-200 px-1.5 py-0.5 text-xs text-slate-600">{{ imageBuilders.length }}</span>
          </button>
        </div>
      </div>

      <div class="p-5">
        <TableImages v-if="displayImages" :images="images" />
        <TableImageBuilders v-else :imageBuilders="imageBuilders" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from "vue";
import { type Image, type ImageBuilder } from "../sdk/backend.generated";

const images = ref<Array<Image>>([]);
const imageBuilders = ref<Array<ImageBuilder>>([]);
const displayImages = ref(true);
const version = ref("");
const refreshing = ref(false);

const completedCount = computed(
  () => images.value.filter((image) => image.status === "COMPLETED").length
);

const { data } = await useFetch("/api/data", { credentials: "include" });
if (data?.value) {
  images.value = data.value.images;
  imageBuilders.value = data.value.imageBuilders;
}

const { data: dataVersion } = await useFetch("/api/version", { credentials: "include" });
if (dataVersion?.value) {
  version.value = dataVersion.value;
}

const fetchData = async () => {
  version.value = await $fetch("/api/version");

  const { images: fetchedImages, imageBuilders: fetchedBuilders } = await $fetch(
    "/api/data"
  );
  images.value = fetchedImages;
  imageBuilders.value = fetchedBuilders;
};

const refresh = async () => {
  refreshing.value = true;
  try {
    await fetchData();
  } finally {
    refreshing.value = false;
  }
};

const createImage = async (form: Record<string, string>) => {
  try {
    await $fetch("/api/image", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: form,
    });

    await fetchData();
  } catch {
    alert("Failed to launch the copy of the image");
  }
};

onMounted(fetchData);
</script>
