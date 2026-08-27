<template>
  <div>
    <button class="btn-primary" @click="showModal = true">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4">
        <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
      </svg>
      Copy image
    </button>

    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="showModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4"
        >
          <div class="absolute inset-0 bg-slate-900/30 backdrop-blur-sm" @click="closeModal" />

          <div
            class="relative z-10 w-full max-w-lg rounded-2xl border border-slate-200 bg-white p-6 shadow-xl shadow-slate-300/40"
            role="dialog"
            aria-labelledby="modal-title"
            aria-modal="true"
          >
            <div class="mb-6 flex items-start justify-between gap-4">
              <div>
                <h2 id="modal-title" class="text-lg font-semibold text-slate-900">Copy image</h2>
                <p class="mt-1 text-sm text-slate-500">Sync a container image between registries.</p>
              </div>
              <button
                class="rounded-lg p-1.5 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700"
                aria-label="Close modal"
                @click="closeModal"
              >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5">
                  <path d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z" />
                </svg>
              </button>
            </div>

            <form class="space-y-4" @submit.prevent="submitForm">
              <div>
                <label for="name" class="mb-1.5 block text-sm font-medium text-slate-700">Name</label>
                <input
                  id="name"
                  v-model="formData.name"
                  type="text"
                  required
                  placeholder="nginx-alpine"
                  class="input-field"
                  autocomplete="off"
                />
              </div>

              <div class="relative">
                <label for="source-repo" class="mb-1.5 block text-sm font-medium text-slate-700">Source repository</label>
                <div class="relative">
                  <input
                    id="source-repo"
                    v-model="formData.sourceRepository"
                    type="text"
                    required
                    placeholder="quay.io/nginx/nginx-ingress"
                    class="input-field pr-10"
                    autocomplete="off"
                    @input="debouncedSearch"
                    @focus="showSuggestions = true"
                  />
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                  </svg>
                </div>

                <ul
                  v-if="showSuggestions && filteredRepositories.length"
                  class="absolute z-20 mt-1 max-h-56 w-full overflow-y-auto rounded-xl border border-slate-200 bg-white shadow-lg"
                >
                  <li
                    v-for="(repo, index) in filteredRepositories"
                    :key="index"
                    class="cursor-pointer border-b border-slate-100 px-3 py-2.5 last:border-0 hover:bg-slate-50"
                    @click="selectRepository(repo)"
                  >
                    <div class="flex items-center justify-between gap-3">
                      <div class="flex min-w-0 items-center gap-2">
                        <img
                          v-if="repo.registry === 'Quay.io'"
                          class="h-5 w-5 shrink-0"
                          src="https://upload.wikimedia.org/wikipedia/commons/d/d8/Red_Hat_logo.svg"
                          alt=""
                        />
                        <img
                          v-if="repo.registry === 'Amazon ECR'"
                          class="h-5 w-5 shrink-0"
                          src="https://upload.wikimedia.org/wikipedia/commons/9/93/Amazon_Web_Services_Logo.svg"
                          alt=""
                        />
                        <img
                          v-if="repo.registry === 'DockerHub'"
                          class="h-5 w-5 shrink-0"
                          src="https://icon.icepanel.io/Technology/svg/Docker.svg"
                          alt=""
                        />
                        <span class="truncate text-sm text-slate-800">{{ repo.name }}</span>
                      </div>
                      <div class="flex shrink-0 items-center gap-2 text-xs text-slate-500">
                        <span>{{ repo.downloadCount }}</span>
                        <svg
                          v-if="repo.isOfficial"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                          class="h-4 w-4 text-emerald-600"
                        >
                          <path fill-rule="evenodd" d="M8.603 3.799A4.49 4.49 0 0 1 12 2.25c1.357 0 2.573.6 3.397 1.549a4.49 4.49 0 0 1 3.498 1.307 4.491 4.491 0 0 1 1.307 3.497A4.49 4.49 0 0 1 21.75 12a4.49 4.49 0 0 1-1.549 3.397 4.491 4.491 0 0 1-3.497 1.307 4.491 4.491 0 0 1-3.498 1.306 4.49 4.49 0 0 1-3.397 1.549A4.49 4.49 0 0 1 12 21.75a4.49 4.49 0 0 1-3.397-1.549 4.49 4.49 0 0 1-3.498-1.306 4.491 4.491 0 0 1-3.497-1.307A4.49 4.49 0 0 1 2.25 12c0-1.357.6-2.573 1.549-3.397a4.49 4.49 0 0 1 1.307-3.497 4.49 4.49 0 0 1 3.497-1.307c.923-.707 2.04-1.106 3.198-1.149ZM9.748 8.25c.742-.742 1.947-.742 2.689 0l.094.094 3.784 3.784.094.094c.742.742.742 1.947 0 2.689l-.094.094-3.784 3.784-.094.094c-.742.742-1.947.742-2.689 0l-.094-.094-3.784-3.784-.094-.094c-.742-.742-.742-1.947 0-2.689l.094-.094 3.784-3.784.094-.094Z" clip-rule="evenodd" />
                        </svg>
                      </div>
                    </div>
                  </li>
                </ul>
              </div>

              <div>
                <label for="source-version" class="mb-1.5 block text-sm font-medium text-slate-700">Source version</label>
                <input
                  id="source-version"
                  v-model="formData.sourceVersion"
                  type="text"
                  required
                  placeholder="3.7-alpine"
                  class="input-field font-mono"
                  autocomplete="off"
                />
              </div>

              <div>
                <label for="destination-repo" class="mb-1.5 block text-sm font-medium text-slate-700">Destination repository</label>
                <input
                  id="destination-repo"
                  v-model="formData.destinationRepository"
                  type="text"
                  required
                  placeholder="myregistry.io/nginx/nginx-ingress"
                  class="input-field font-mono"
                  autocomplete="off"
                />
              </div>

              <div>
                <label for="destination-version" class="mb-1.5 block text-sm font-medium text-slate-700">Destination version</label>
                <input
                  id="destination-version"
                  v-model="formData.destinationVersion"
                  type="text"
                  required
                  placeholder="3.7-alpine"
                  class="input-field font-mono"
                  autocomplete="off"
                />
              </div>

              <div>
                <label for="mode" class="mb-1.5 block text-sm font-medium text-slate-700">Mode</label>
                <select id="mode" v-model="formData.mode" class="input-field">
                  <option value="OneShot">OneShot</option>
                  <option value="OnceByTag">OnceByTag</option>
                  <option value="Recurrent">Recurrent</option>
                </select>
              </div>

              <div class="flex gap-3 pt-2">
                <button type="button" class="btn-secondary flex-1" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn-primary flex-1">
                  Start copy
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";
import { useFetch } from "#app";
import type { RegistryImage } from "~/server/api/images/search";

const emit = defineEmits(["create"]);

const showModal = ref(false);
const showSuggestions = ref(false);
const filteredRepositories = ref<Array<RegistryImage>>([]);
const formData = ref({
  destinationRepository: "",
  destinationVersion: "",
  mode: "OneShot" as "OneShot" | "OnceByTag" | "Recurrent",
  name: "",
  sourceRepository: "",
  sourceVersion: "",
});
let timeout: NodeJS.Timeout | undefined;

const closeModal = () => {
  showModal.value = false;
  showSuggestions.value = false;
};

const submitForm = () => {
  showModal.value = false;
  showSuggestions.value = false;
  emit("create", formData.value);
};

const debouncedSearch = () => {
  clearTimeout(timeout);
  timeout = setTimeout(() => {
    filterSuggestions();
  }, 500);
};

const filterSuggestions = async () => {
  if (formData.value.sourceRepository) {
    const { data: repos } = await useFetch<Array<RegistryImage>>(
      `/api/images/search?repo=${formData.value.sourceRepository.toLowerCase()}`
    );
    filteredRepositories.value =
      repos.value?.filter((repo) =>
        repo.name.toLowerCase().includes(formData.value.sourceRepository.toLowerCase())
      ) || [];
  } else {
    filteredRepositories.value = [];
  }
};

const selectRepository = (repo: RegistryImage) => {
  formData.value.sourceRepository = repo.name;
  showSuggestions.value = false;
};

watch(() => formData.value.sourceRepository, () => debouncedSearch());
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
