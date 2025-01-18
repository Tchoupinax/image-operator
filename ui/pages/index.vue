<template>
  <div class="p-2">
    <div class="max-w-6xl mx-auto mt-8">
      <h1 class="mb-16 text-4xl font-thin">Image Operator <span v-if="version">({{ version }})</span></h1>

      <div class="flex justify-between my-6">
        <div>
          <button :class="{ 'bg-blue-500': displayImages, 'bg-gray-500': !displayImages }"
            class="px-4 py-2 mr-2 font-bold text-white rounded" @click="displayImages = true">Images</button>
          <button :class="{ 'bg-blue-500': !displayImages, 'bg-gray-500': displayImages }"
            class="px-4 py-2 mx-2 font-bold text-white rounded" @click="displayImages = false">Image
            Builders</button>
        </div>

        <ModalCopyImage v-if="displayImages" @create="createImage" />
      </div>

      <div v-if="displayImages">
        <TableImages :images="images" />
      </div>

      <div v-else>
        <TableImageBuilders :imageBuilders="imageBuilders" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { type Image, type ImageBuilder } from "../sdk/backend.generated";

const images = ref<Array<Image>>([]);
const imageBuilders = ref<Array<ImageBuilder>>([]);
const displayImages = ref(true);
const version = ref("");

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

// Method to create an image
const createImage = async (form: any) => {
  await $fetch("/api/image", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: form,
  });
};

onMounted(fetchData);
</script>
