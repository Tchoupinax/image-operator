<template>
  <div>
    <h2 class="mb-4 text-2xl font-semibold text-gray-800">Images</h2>
    <table class="min-w-full bg-white border border-gray-300 rounded-lg shadow">
      <thead>
        <tr class="text-sm leading-normal text-gray-600 uppercase bg-gray-200">
          <th class="px-6 py-3 text-left">Name</th>
          <th class="px-6 py-3 text-left">Source/Destination</th>
          <th class="px-6 py-3 text-left">Status</th>
          <th class="px-6 py-3 text-left">Created At</th>
        </tr>
      </thead>
      <tbody class="text-sm font-light text-gray-700">
        <tr v-for="(image, index) in images" :key="index"
          class="transition-colors border-b border-gray-200 hover:bg-gray-100">
          <td class="px-6 py-3 text-lg">{{ image.name }}</td>
          <td class="flex-col items-end justify-end px-6 py-3 text-lg">
            <p>{{ image.source.name }}:{{ image.source.version }}</p>
            <!--<p>⬇</p>
            <p>{{ image.destination.name }}:{{ image.destination.version }}</p>-->
          </td>
          <td class="px-6 py-3 text-xs font-bold" :class="{ 'text-green-400': image.status === 'COMPLETED' }">{{
            image.status }}
          </td>
          <td class="px-6 py-3 text-lg">{{ formatDate(image.createdAt) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  props: ["images"],
  methods: {
    formatDate(dateString) {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    },
  },
};
</script>
