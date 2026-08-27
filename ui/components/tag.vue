<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset"
    :class="variantClasses"
  >
    {{ text }}
  </span>
</template>

<script lang="ts">
export default {
  props: {
    text: {
      type: String,
      required: true,
    },
    variant: {
      type: String,
      default: "sky",
    },
  },
  computed: {
    variantClasses(): string {
      const variants: Record<string, string> = {
        sky: "bg-sky-50 text-sky-700 ring-sky-200",
        emerald: "bg-emerald-50 text-emerald-700 ring-emerald-200",
        violet: "bg-violet-50 text-violet-700 ring-violet-200",
      };

      return variants[this.variant] || variants.sky;
    },
  },
};
</script>
