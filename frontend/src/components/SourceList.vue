<template>
  <aside class="flex flex-col border-r border-gray-300 max-w-xs">
    <h3 class="my-2 ml-4 text-lg font-bold">Sources</h3>
    <div v-if="sources.length === 0" class="flex flex-col flex-1 justify-center items-center text-gray-500">
      <h4 class="text-lg font-semibold">No sources available.</h4>
      <p>Please upload resources.</p>
    </div>
    <ul v-else class="list flex-1 overflow-auto m-0 p-0">
      <li
          v-for="source in sources"
          :key="source.id"
          :class="[
          'list-none m-0 p-0',
          source.id === selectedSourceId ? 'bg-gray-200' : ''
        ]"
      >

        <a
            href="#"
            @click.prevent="selectSource(source.id)"
            class="flex flex-col relative p-2 border-b border-gray-200 transition-colors duration-300 hover:bg-gray-100"
        >
          <span class="font-bold">{{ source.id }}</span>
          <span class="text-gray-600 text-sm">{{ source.event_type }}</span>
          <span
              v-if="source.events.length > 0"
              class="text-white bg-red-700 rounded-full px-2 py-0.5 text-xs absolute right-3 bottom-2"
          >
            Events: {{ source.events.length }}
          </span>
        </a>
      </li>
    </ul>
    <div
        @dragover.prevent
        @drop.prevent="handleDrop"
        class="border-2 border-dashed border-gray-300 p-5 relative rounded-lg m-1"
    >
      <p class="text-center text-gray-500">Drag and drop files here or click the button to upload.</p>

      <input
          type="file"
          multiple
          @change="handleFilesUpload"
          ref="fileInput"
          class="hidden"
      />

      <button
          @click="triggerFileSelect"
          class="mt-4 mx-auto block bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors duration-300"
      >
        Upload Files
      </button>
    </div>
  </aside>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {Source} from '../types/Source';
import {sources, upload} from '../api.ts';

export default defineComponent({
  name: 'SourceList',
  data() {
    return {
      sources: [] as Source[],
      uploadFiles: [] as File[],
      selectedSourceId: '' as string,
    };
  },
  emits: ['source-selected'],
  created() {
    this.fetchSources();
  },
  methods: {
    fetchSources() {
      sources()
          .then((response) => {
            this.sources = response.data as Source[];
          })
          .catch((error) => {
            console.error('Error fetching sources:', error);
          });
    },
    selectSource(id: string) {
      this.selectedSourceId = id;
      this.$emit('source-selected', id);
    },
    triggerFileSelect() {
      const fileInput = this.$refs.fileInput as HTMLInputElement;
      if (fileInput) {
        fileInput.click();
      }
    },
    handleFilesUpload(event: Event) {
      const target = event.target as HTMLInputElement;
      if (target.files) {
        this.uploadFiles = Array.from(target.files);
        this.uploadFilesToServer();
      }
    },
    handleDrop(event: DragEvent) {
      event.preventDefault();
      event.stopPropagation();
      if (event.dataTransfer) {
        const files = event.dataTransfer.files;
        this.uploadFiles = Array.from(files);
        this.uploadFilesToServer();
      }
    },
    uploadFilesToServer() {
      const formData = new FormData();
      const videoFiles: File[] = [];
      const imageFiles: File[] = [];
      const jsonFiles: File[] = [];
      for (const file of this.uploadFiles) {
        const fileName = file.name.toLowerCase();
        if (
            file.type.startsWith('video/') ||
            fileName.endsWith('.mp4') ||
            fileName.endsWith('.avi') ||
            fileName.endsWith('.mov')
        ) {
          videoFiles.push(file);
        } else if (
            file.type.startsWith('image/') ||
            fileName.endsWith('.jpg') ||
            fileName.endsWith('.jpeg') ||
            fileName.endsWith('.png')
        ) {
          imageFiles.push(file);
        } else if (
            file.type === 'application/json' ||
            fileName.endsWith('.json')
        ) {
          jsonFiles.push(file);
        } else {
          console.warn('Unsupported file type:', file.name);
        }
      }

      videoFiles.forEach((file) => {
        formData.append('videos', file);
      });

      imageFiles.forEach((file) => {
        formData.append('images', file);
      });

      jsonFiles.forEach((file) => {
        formData.append('json', file);
      });

      upload(formData)
          .then((response) => {
            console.log('Files uploaded successfully:', response.data);
            this.fetchSources();
          })
          .catch((error) => {
            console.error('Error uploading files:', error);
            this.fetchSources();
          });
    },
  },
});
</script>
