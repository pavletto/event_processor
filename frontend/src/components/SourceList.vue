<template>
  <aside>
    <div class="empty-list" v-if="sources.length === 0">
      <h4>No sources available.</h4><p>Please upload resources.</p>
    </div>
    <ul v-else class="list">
      <li
          class="item"
          v-for="source in sources"
          :key="source.id"
          :class="{ 'item__selected': source.id === selectedSourceId }"
      >
        <a class="item__link" href="#" @click.prevent="selectSource(source.id)">
          <span class="item__event-id"> {{ source.id }}</span>
          <span class="item__event-type">{{ source.event_type }}</span>
          <span
              v-if="source.events.length > 0"
              class="item__events-count"
          >
            Events:{{ source.events.length }}
          </span>
        </a>
      </li>
    </ul>
    <div
        @dragover.prevent
        @drop.prevent="handleDrop"
        class="dropzone"
    >
      <p>Drag and drop files here or click the button to upload.</p>

      <input
          type="file"
          multiple
          @change="handleFilesUpload"
          ref="fileInput"
          style="display: none;"
      />

      <button @click="triggerFileSelect">Upload Files</button>
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
      selectedSourceId: '' as string,    };
  },
  created() {
    this.fetchSources();
  },
  methods: {
    fetchSources() {
      sources<Source[]>()
          .then((response) => {
            this.sources = response.data;
          })
          .catch((error) => {
            console.error('Error fetching sources:', error);
          });
    },
    selectSource(id: string) {
      this.selectedSourceId = id;      this.$emit('source-selected', id);
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

<style scoped>
aside {
  display: flex;
  flex-direction: column;
  min-width: 200px;
  max-width: 330px;
}

.empty-list {
  display: flex;
  flex-direction: column;
  flex:1;
  color:#666;
  justify-content: center;
  align-items: center;
}
.list {
  margin: 0;
  padding: 0;
  flex: 1;
  overflow: auto;
}

.item {
  list-style: none;
  margin: 0;
  padding: 0;
}

.item__link {
  display: flex;
  position: relative;
  flex-direction: column;
  padding: 10px;
  text-decoration: none;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.3s;
}

.item__link:hover {
  background-color: #f9f9f9;
}

.item__selected .item__link {
  background-color: #e0f7fa; }

.item__events-count {
  color: white;
  background: darkred;
  border-radius: 1rem;
  padding: 0 4px;
  position: absolute;
  right: 12px;
  font-size: 0.6rem;
}

.item__event-id {
  font-weight: bold;
  float: left;
}

.item__event-type {
  color: #666;
  font-size: 0.8rem;
}

.dropzone {
  border: 2px dashed #ccc;
  padding: 20px;
  position: relative;
  border-radius: 8px;
  margin: 4px;
}

.dropzone p {
  margin: 0;
  text-align: center;
  color: #666;
}

.dropzone button {
  display: block;
  margin: 10px auto;
}

.item__selected .item__link {
  transition: background-color 0.3s;
}
</style>
