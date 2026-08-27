<template>
  <div class="group/ref flex min-w-0 items-start gap-1.5">
    <div class="min-w-0 flex-1">
      <p
        v-if="label"
        class="mb-0.5 text-[10px] font-semibold uppercase tracking-wider"
        :class="labelClass"
      >
        {{ label }}
      </p>
      <p
        class="truncate font-mono text-xs leading-5"
        :class="textClass"
        :title="fullRef"
      >
        {{ fullRef }}
      </p>
    </div>
    <button
      type="button"
      class="mt-0.5 shrink-0 rounded-md p-1 text-slate-400 opacity-0 transition hover:bg-white hover:text-slate-700 group-hover/ref:opacity-100 focus:opacity-100"
      :title="copied ? 'Copied' : 'Copy'"
      @click="copy"
    >
      <svg
        v-if="!copied"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="h-3.5 w-3.5"
      >
        <path
          fill-rule="evenodd"
          d="M13.887 3.182c.396-.037.79.008 1.168.133 1.55.469 2.345 2.237 1.762 3.778l-.496 1.315a.75.75 0 0 0 1.404.528l.496-1.314c1.033-2.744-.406-5.83-3.15-6.864a4.5 4.5 0 0 0-1.674-.267 4.5 4.5 0 0 0-2.17.503l-7.502 3.75a3.75 3.75 0 0 0-1.695 5.002l.992 1.648a.75.75 0 1 0 1.272-.768l-.992-1.648A2.25 2.25 0 0 1 3.6 7.352l7.502-3.75a3 3 0 0 1 1.785-.42Zm-4.106 8.576a.75.75 0 0 0-1.062-.853l-1.106 1.106a2.25 2.25 0 0 0-.393 2.701l.992 1.648a2.25 2.25 0 0 0 3.002.393l1.106-1.106a.75.75 0 0 0-.853-1.062l-1.106 1.106a.75.75 0 0 1-1.001.131l-.992-1.648a.75.75 0 0 1 .131-1.001l1.106-1.106Z"
          clip-rule="evenodd"
        />
      </svg>
      <svg
        v-else
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="h-3.5 w-3.5 text-emerald-600"
      >
        <path
          fill-rule="evenodd"
          d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
          clip-rule="evenodd"
        />
      </svg>
    </button>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";

const props = defineProps<{
  label?: string;
  name?: string | null;
  version?: string | null;
  variant?: "source" | "destination";
}>();

const copied = ref(false);

const fullRef = computed(() => {
  const name = props.name || "";
  const version = props.version || "";
  return version ? `${name}:${version}` : name;
});

const labelClass = computed(() =>
  props.variant === "source" ? "text-sky-600" : "text-violet-600"
);

const textClass = computed(() =>
  props.variant === "source" ? "text-sky-800" : "text-violet-800"
);

const copy = async () => {
  try {
    await navigator.clipboard.writeText(fullRef.value);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 1500);
  } catch {
    // ignore clipboard errors
  }
};
</script>
