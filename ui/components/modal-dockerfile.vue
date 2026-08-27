<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click="closeModal"
      >
        <div class="absolute inset-0 bg-slate-900/30 backdrop-blur-sm" />

        <div
          class="relative z-10 flex max-h-[85vh] w-full max-w-4xl flex-col overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-xl shadow-slate-300/40"
          @click.stop
        >
          <div class="flex items-center justify-between border-b border-slate-200 px-5 py-4">
            <div>
              <h3 class="text-base font-semibold text-slate-900">Dockerfile</h3>
              <p class="text-xs text-slate-500">Build source definition</p>
            </div>
            <button
              class="rounded-lg p-1.5 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700"
              @click="closeModal"
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-5 w-5">
                <path d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z" />
              </svg>
            </button>
          </div>

          <div class="overflow-auto bg-slate-50 p-5">
            <pre class="rounded-xl border border-slate-200 bg-white p-4"><code ref="codeBlock" class="language-dockerfile text-sm leading-relaxed">{{ code }}</code></pre>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script>
import hljs from "highlight.js";
import "highlight.js/styles/github.min.css";

export default {
  name: "CodeModal",
  props: {
    code: {
      type: String,
      required: true,
    },
    visible: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        this.$nextTick(() => {
          hljs.highlightElement(this.$refs.codeBlock);
        });
      }
    },
  },
  methods: {
    closeModal() {
      this.$emit("close");
    },
  },
};
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
