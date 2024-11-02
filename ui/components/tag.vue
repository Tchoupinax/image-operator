<template>
  <span :class="tagClasses" class="inline-flex items-center px-3 py-1 text-sm font-bold rounded-full">
    {{ text }}
  </span>
</template>

<script>
export default {
  props: {
    text: {
      type: String,
      required: true,
    },
    color: {
      type: String,
      default: "bg-blue-500",
    },
    borderColor: {
      type: String,
      default: "border-blue-700",
    },
  },
  computed: {
    tagClasses() {
      return [
        this.color,
        this.borderColor,
        "border-2 border-solid",
        "text-white",
        "mx-1"
      ];
    },
  },
};
</script>
