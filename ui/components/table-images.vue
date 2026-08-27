<template>
  <div>
    <div
      v-if="!images?.length"
      class="flex flex-col items-center justify-center rounded-xl border border-dashed border-slate-300 bg-slate-50 px-6 py-16 text-center"
    >
      <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-500">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="h-6 w-6">
          <path
            fill-rule="evenodd"
            d="M6.75 2.25A.75.75 0 0 1 7.5 3v1.5h9V3A.75.75 0 0 1 18 3v1.5h.75a3 3 0 0 1 3 3v11.25a3 3 0 0 1-3 3H5.25a3 3 0 0 1-3-3V7.5a3 3 0 0 1 3-3H6V3a.75.75 0 0 1 .75-.75Zm13.5 9a1.5 1.5 0 0 0-1.5-1.5H5.25a1.5 1.5 0 0 0-1.5 1.5v7.5a1.5 1.5 0 0 0 1.5 1.5h13.5a1.5 1.5 0 0 0 1.5-1.5v-7.5Z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
      <h3 class="text-base font-semibold text-slate-900">No images yet</h3>
      <p class="mt-1 max-w-sm text-sm text-slate-500">
        Create an image sync job to copy container images between registries.
      </p>
    </div>

    <div v-else class="overflow-hidden rounded-xl border border-slate-200">
      <div class="overflow-x-auto">
        <table class="min-w-[960px] w-full table-fixed">
          <colgroup>
            <col class="w-[18%]" />
            <col class="w-[42%]" />
            <col class="w-[12%]" />
            <col class="w-[14%]" />
            <col class="w-[14%]" />
          </colgroup>
          <thead>
            <tr class="border-b border-slate-200 bg-slate-50/90">
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Name
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Sync path
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Status
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Created
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500">
                Last run
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 bg-white">
            <tr
              v-for="(image, index) in images"
              :key="image.name || index"
              class="group transition hover:bg-slate-50/80"
            >
              <td class="px-4 py-3 align-top">
                <p
                  class="truncate text-sm font-medium text-slate-900"
                  :title="image.name || undefined"
                >
                  {{ image.name }}
                </p>
              </td>

              <td class="px-4 py-3 align-top">
                <div class="max-w-xl space-y-1.5 rounded-lg border border-slate-200/80 bg-slate-50/70 p-2.5">
                  <CopyableRef
                    label="Source"
                    variant="source"
                    :name="image.source?.name"
                    :version="image.source?.version"
                  />
                  <div class="flex items-center gap-2 px-1">
                    <div class="h-px flex-1 bg-slate-200" />
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      class="h-3.5 w-3.5 shrink-0 text-slate-400"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M10 3a.75.75 0 0 1 .75.75v10.638l3.96-4.158a.75.75 0 1 1 1.08 1.04l-5.25 5.5a.75.75 0 0 1-1.08 0l-5.25-5.5a.75.75 0 1 1 1.08-1.04l3.96 4.158V3.75A.75.75 0 0 1 10 3Z"
                        clip-rule="evenodd"
                      />
                    </svg>
                    <div class="h-px flex-1 bg-slate-200" />
                  </div>
                  <CopyableRef
                    label="Destination"
                    variant="destination"
                    :name="image.destination?.name"
                    :version="image.destination?.version"
                  />
                </div>
              </td>

              <td class="px-4 py-3 align-top">
                <StatusBadge :status="image.status" />
              </td>

              <td class="px-4 py-3 align-top text-sm text-slate-600">
                {{ formatDate(image.createdAt) }}
              </td>

              <td class="px-4 py-3 align-top text-sm text-slate-600">
                {{ image.lastExecution ? format(image.lastExecution) : "Never" }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { format } from "timeago.js";

export default {
  props: ["images"],
  methods: {
    format,
    formatDate(dateString: string) {
      const options: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "short",
        day: "numeric",
      };
      return new Date(dateString).toLocaleDateString(undefined, options);
    },
  },
};
</script>
