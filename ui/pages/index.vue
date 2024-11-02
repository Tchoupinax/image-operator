<template>
  <div>
    <div class="max-w-4xl mx-auto mt-8">
      <h1 class="mb-16 text-4xl font-thin">Image Operator</h1>

      <div class="my-6">
        <button :class="{ 'bg-blue-500': displayImages, 'bg-gray-500': !displayImages }"
          class="px-4 py-2 mx-2 font-bold text-white rounded" @click="displayImages = true">Images</button>
        <button :class="{ 'bg-blue-500': !displayImages, 'bg-gray-500': displayImages }"
          class="px-4 py-2 mx-2 font-bold text-white rounded" @click="displayImages = false">Image
          Builders</button>
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

<script lang="ts">
import TableImageBuilders from '~/components/table-image-builders.vue';
import { type Image, type ImageBuilder } from '../sdk/backend.generated';

type Store = {
  images: Array<Image>;
  imageBuilders: Array<ImageBuilder>;
  displayImages: boolean;
}

export default {
  data(): Store {
    return {
      displayImages: true,
      images: [],
      imageBuilders: []
    }
  },
  async mounted() {
    const { images, imageBuilders } = await $fetch("/api/data");
    this.images = images;
    this.imageBuilders = imageBuilders;
  }
};
</script>
