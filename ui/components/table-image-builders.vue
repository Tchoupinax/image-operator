<template>
  <div>
    <h2 class="mb-4 text-2xl font-semibold text-gray-800">Image builders</h2>
    <table class="min-w-full bg-white border border-gray-300 rounded-lg shadow">
      <thead>
        <tr class="text-sm leading-normal text-gray-600 uppercase bg-gray-200">
          <th class="px-6 py-3 text-left">Name</th>
          <th class="px-6 py-3 text-left">Architecture</th>
          <th class="px-6 py-3 text-left">Created At</th>
        </tr>
      </thead>
      <tbody class="text-sm font-light text-gray-700">
        <tr v-for="(imageBuilder, index) in imageBuilders" :key="index"
          class="transition-colors border-b border-gray-200 hover:bg-gray-100">
          <td class="px-6 py-3 text-lg">{{ imageBuilder.name }}</td>
          <td class="px-6 py-3">
            <Tag v-if="imageBuilder.architecture !== 'Arm64'" text="Arm64" color="bg-blue-500"
              borderColor="border-blue-700" />
            <Tag v-if="imageBuilder.architecture !== 'Amd64'" text="Amd64" color="bg-green-500"
              borderColor="border-green-700" />
          </td>
          <td class="px-6 py-3 text-lg">{{ formatDate(imageBuilder.createdAt) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  props: ["imageBuilders"],
  methods: {
    formatDate(dateString) {
      const options = { year: "numeric", month: "long", day: "numeric", hour: "2-digit", minute: "2-digit" };
      return new Date(dateString).toLocaleDateString(undefined, options);
    },
  },
};
</script>
