<template>
  <div>
    <div
      v-if="!imageBuilders?.length"
      class="flex flex-col items-center justify-center rounded-xl border border-dashed border-slate-300 bg-slate-50 px-6 py-16 text-center"
    >
      <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-500">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="h-6 w-6">
          <path
            d="M11.25 4.533A9.707 9.707 0 0 0 6 3a9.735 9.735 0 0 0-3.25.555.75.75 0 0 0-.5.707v14.25a.75.75 0 0 0 1 .707A8.237 8.237 0 0 1 6 18.75c1.995 0 3.823.707 5.25 1.886V4.533ZM12.75 20.636A8.214 8.214 0 0 1 18 18.75c.966 0 1.89.166 2.75.47a.75.75 0 0 0 1-.708V4.262a.75.75 0 0 0-.5-.707A9.735 9.735 0 0 0 18 3a9.707 9.707 0 0 0-5.25 1.533v16.103Z"
          />
        </svg>
      </div>
      <h3 class="text-base font-semibold text-slate-900">No image builders</h3>
      <p class="mt-1 max-w-sm text-sm text-slate-500">
        ImageBuilder resources will appear here when configured in the cluster.
      </p>
    </div>

    <div v-else class="overflow-hidden rounded-xl border border-slate-200">
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="border-b border-slate-200 bg-slate-50/90">
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Name
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Architecture
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Created
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Dockerfile
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr
              v-for="(imageBuilder, index) in imageBuilders"
              :key="imageBuilder.name || index"
              class="transition hover:bg-slate-50/80"
            >
              <td class="px-4 py-3">
                <p
                  class="max-w-xs truncate text-sm font-medium text-slate-900"
                  :title="imageBuilder.name"
                >
                  {{ imageBuilder.name }}
                </p>
              </td>
              <td class="px-4 py-3">
                <div class="flex flex-wrap gap-1.5">
                  <Tag v-if="imageBuilder.architecture !== 'Arm64'" text="Arm64" variant="sky" />
                  <Tag v-if="imageBuilder.architecture !== 'Amd64'" text="Amd64" variant="emerald" />
                </div>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-sm text-slate-600">
                {{ formatDate(imageBuilder.createdAt) }}
              </td>
              <td class="px-4 py-3">
                <button
                  class="btn-secondary !px-3 !py-1.5 text-xs"
                  @click="showCode(imageBuilder.name)"
                >
                  View source
                </button>
                <ModalDockerfile
                  :code="imageBuilder.source"
                  :visible="showCodeModalName === imageBuilder.name"
                  @close="showCodeModalName = undefined"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: ["imageBuilders"],
  data(): { showCodeModalName?: string } {
    return {
      showCodeModalName: undefined,
    };
  },
  methods: {
    showArchitecture(architecture: string, target: string) {
      const normalized = (architecture || "").toLowerCase();
      if (normalized === "both") {
        return true;
      }
      return normalized === target;
    },
    showCode(name: string) {
      this.showCodeModalName = name;
    },
    formatDate(dateString: string) {
      const options: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      };
      return new Date(dateString).toLocaleDateString(undefined, options);
    },
  },
};
</script>
