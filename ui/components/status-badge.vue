<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold ring-1 ring-inset"
    :class="statusClasses"
  >
    {{ status || "UNKNOWN" }}
  </span>
</template>

<script lang="ts">
export default {
  props: {
    status: {
      type: String,
      default: "",
    },
  },
  computed: {
    statusClasses(): string {
      const normalized = (this.status || "").toUpperCase();

      if (normalized === "COMPLETED") {
        return "bg-emerald-50 text-emerald-700 ring-emerald-200";
      }
      if (normalized === "FAILED" || normalized === "ERROR") {
        return "bg-rose-50 text-rose-700 ring-rose-200";
      }
      if (normalized === "RUNNING" || normalized === "PENDING") {
        return "bg-amber-50 text-amber-700 ring-amber-200";
      }

      return "bg-slate-100 text-slate-600 ring-slate-200";
    },
  },
};
</script>
