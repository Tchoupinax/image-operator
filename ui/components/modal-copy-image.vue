<template>
  <div>
    <button @click="showModal = true" class="px-4 py-2 mx-2 font-bold text-white bg-green-500 rounded">
      Copy image
    </button>

    <transition name="modal-fade">
      <div v-if="showModal" class="modal-overlay" @click="closeModal" aria-hidden="true"></div>
    </transition>

    <transition name="modal-fade">
      <div v-if="showModal" class="modal" role="dialog" aria-labelledby="modal-title" aria-modal="true">
        <div class="modal-header">
          <h2 id="modal-title" class="font-bold underline">Copy image</h2>
          <button @click="closeModal" class="close-button" aria-label="Close modal">×</button>
        </div>

        <form @submit.prevent="submitForm">
          <div class="form-group">
            <label for="name">Name</label>
            <input v-model="formData.name" type="text" id="name" required placeholder="Alpine"
              class="border-b-2 border-gray-300 focus:outline-none" autocomplete="off" />
          </div>

          <div class="form-group autocomplete">
            <label for="source-repo">Source Repository</label>
            <div class="flex items-center justify-center">
              <input v-model="formData.sourceRepository" type="text" id="source-repo" required
                placeholder="quay.io/nginx/nginx-ingress" @input="debouncedSearch" @focus="showSuggestions = true"
                class="border-b-2 border-gray-300 focus:outline-none" autocomplete="off" />
              <div>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="size-6">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="m15.75 15.75-2.489-2.489m0 0a3.375 3.375 0 1 0-4.773-4.773 3.375 3.375 0 0 0 4.774 4.774ZM21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                </svg>
              </div>
            </div>

            <ul v-if="showSuggestions && filteredRepositories.length" class="suggestions-dropdown">
              <li v-for="(repo, index) in filteredRepositories" :key="index" @click="selectRepository(repo)"
                class="suggestion-item">
                <div class="flex justify-between">
                  <div class="flex">
                    <img v-if="repo.registry === 'Quay.io'" class="mr-2 size-6"
                      src="https://upload.wikimedia.org/wikipedia/commons/d/d8/Red_Hat_logo.svg" />
                    <img v-if="repo.registry === 'Amazon ECR'" class="mr-2 size-6"
                      src="https://upload.wikimedia.org/wikipedia/commons/9/93/Amazon_Web_Services_Logo.svg" />
                    <img v-if="repo.registry === 'DockerHub'" class="mr-2 size-6"
                      src="https://icon.icepanel.io/Technology/svg/Docker.svg" />

                    <p>
                      {{ repo.name }}
                    </p>
                  </div>

                  <div class="flex">
                    <p class="mr-1">
                      {{ repo.downloadCount }}
                    </p>

                    <svg v-if="repo.isOfficial" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                      stroke-width="1.5" stroke="currentColor" class="text-green-600 size-6">
                      <path stroke-linecap="round" stroke-linejoin="round"
                        d="M9 12.75 11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 0 1-1.043 3.296 3.745 3.745 0 0 1-3.296 1.043A3.745 3.745 0 0 1 12 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 0 1-3.296-1.043 3.745 3.745 0 0 1-1.043-3.296A3.745 3.745 0 0 1 3 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 0 1 1.043-3.296 3.746 3.746 0 0 1 3.296-1.043A3.746 3.746 0 0 1 12 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 0 1 3.296 1.043 3.746 3.746 0 0 1 1.043 3.296A3.745 3.745 0 0 1 21 12Z" />
                    </svg>
                  </div>
                </div>
              </li>
            </ul>
          </div>

          <div class="form-group">
            <label for="source-version">Source Version</label>
            <input v-model="formData.sourceVersion" type="text" id="source-version" required placeholder="v1.2.3"
              class="border-b-2 border-gray-300 focus:outline-none" autocomplete="off" />
          </div>

          <div class="form-group">
            <label for="destination-repo">Destination Repository</label>
            <input v-model="formData.destinationRepository" type="text" id="destination-repo" required
              placeholder="myregistry.io/nginx/nginx-ingress" class="border-b-2 border-gray-300 focus:outline-none"
              autocomplete="off" />
          </div>

          <div class="form-group">
            <label for="destination-version">Destination Version</label>
            <input v-model="formData.destinationVersion" type="text" id="destination-version"
              class="border-b-2 border-gray-300 focus:outline-none" required placeholder="v1.2.3" autocomplete="off" />
          </div>

          <div class="form-group">
            <label>Mode</label>
            <select v-model="formData.mode">
              <option value="OneShot">OneShot</option>
              <option value="OnceByTag">OnceByTag</option>
              <option value="Recurrent">Recurrent</option>
            </select>
          </div>

          <button type="submit" class="submit-button">Submit</button>
        </form>
      </div>
    </transition>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";
import { useFetch } from "#app";
import type { RegistryImage } from "~/server/api/images/search";

type Store = {
  showModal: boolean;
  showSuggestions: boolean;
  repositories: Array<RegistryImage>;
  filteredRepositories: Array<RegistryImage>;
  formData: {
    destinationRepository: string;
    destinationVersion: string;
    mode: "OneShot" | "OnceByTag" | "Recurrent",
    name: string;
    sourceRepository: string;
    sourceVersion: string;
  };
  timeout?: NodeJS.Timeout,
}

const emit = defineEmits(["create"])

const showModal = ref(false);
const showSuggestions = ref(false);
const filteredRepositories = ref<Array<RegistryImage>>([]);
const formData = ref({
  destinationRepository: "",
  destinationVersion: "",
  mode: "OneShot" as "OneShot" | "OnceByTag" | "Recurrent",
  name:  "",
  sourceRepository:  "",
  sourceVersion:  "",
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
    const { data: repos } = await useFetch<Array<RegistryImage>>(`/api/images/search?repo=${formData.value.sourceRepository.toLowerCase()}`);
    filteredRepositories.value = repos.value?.filter((repo) =>
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
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  padding: 20px;
  width: 90%;
  max-width: 600px;
  z-index: 1001;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.close-button {
  font-size: 1.5rem;
  background: none;
  border: none;
  cursor: pointer;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.checkbox-group label {
  font-weight: normal;
}

.submit-button {
  background-color: #007bff;
  color: white;
  padding: 10px;
  border: none;
  cursor: pointer;
  width: 100%;
  border-radius: 4px;
}

.submit-button:hover {
  background-color: #0056b3;
}

/* Autocomplete dropdown styling */
.autocomplete {
  position: relative;
}

.suggestions-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 2px solid #ddd;
  border-top: none;
  box-shadow: inset;
  max-height: 250px;
  overflow-y: auto;
  z-index: 1002;
  list-style: none;
  padding: 0;
  margin: 0;
}

.suggestion-item {
  padding: 8px;
  cursor: pointer;
}

.suggestion-item:hover {
  background-color: #f0f0f0;
}
</style>
