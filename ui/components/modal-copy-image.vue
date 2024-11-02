<template>
  <div>
    <button @click="showModal = true" class="px-4 py-2 mx-2 font-bold text-white bg-green-500 rounded">
      Copy image
    </button>

    <div v-if="showModal" class="modal-overlay" @click="closeModal"></div>

    <div v-if="showModal" class="modal">
      <div class="modal-header">
        <h2 class="font-bold underline">Copy image</h2>
        <button @click="closeModal" class="close-button">×</button>
      </div>

      <form @submit.prevent="submitForm">
        <div class="form-group">
          <label for="name">Name</label>
          <input v-model="formData.name" type="text" id="name" required placeholder="Alpine"
            class="border-b-2 border-black" />
        </div>

        <div class="form-group">
          <label for="source-repo">Source Repository</label>
          <input v-model="formData.sourceRepository" type="text" id="source-repo" required
            placeholder="e.g., github.com/user/source-repo" class="border-b-2 border-black">
        </div>

        <div class="form-group">
          <label for="source-version">Source Version</label>
          <input v-model="formData.sourceVersion" type="text" id="source-version" required placeholder="e.g., v1.0.0"
            class="border-b-2 border-black" />
        </div>

        <div class="form-group">
          <label for="destination-repo">Destination Repository</label>
          <input v-model="formData.destinationRepository" type="text" id="destination-repo" required
            placeholder="e.g., github.com/user/destination-repo" class="border-b-2 border-black" />
        </div>

        <!-- Destination Version -->
        <div class="form-group">
          <label for="destination-version">Destination Version</label>
          <input v-model="formData.destinationVersion" type="text" id="destination-version"
            class="border-b-2 border-black" required placeholder="e.g., v1.0.1" />
        </div>

        <!-- Recurrent Checkbox -->
        <div class="form-group checkbox-group">
          <label for="recurrent">
            <input v-model="formData.recurrent" type="checkbox" id="recurrent" />
            Recurrent
          </label>
        </div>

        <button type="submit" class="submit-button">Submit</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  emits: ["create"],
  data() {
    return {
      showModal: false,
      formData: {
        destinationRepository: '',
        destinationVersion: '',
        mode: "OneShot",
        name: "titi",
        sourceRepository: 'quay.io/nginx/nginx-ingress',
        sourceVersion: '3.7-alpine',
      },
    };
  },
  methods: {
    closeModal() {
      this.showModal = false;
    },
    submitForm() {
      this.showModal = false;
      this.$emit("create", this.formData)
    },
  },
};
</script>

<style scoped>
/* Modal styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  padding: 20px;
  width: 90%;
  max-width: 600px;
  z-index: 1001;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.close-button {
  font-size: 1.5rem;
  background: none;
  border: none;
  cursor: pointer;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.checkbox-group label {
  font-weight: normal;
}

.submit-button {
  background-color: #007bff;
  color: white;
  padding: 10px;
  border: none;
  cursor: pointer;
  width: 100%;
  border-radius: 4px;
}

.submit-button:hover {
  background-color: #0056b3;
}
</style>
