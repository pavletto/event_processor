<template>
    <SourceList @source-selected="handleSourceSelected" />
    <SourceDetail :sourceId="selectedSourceId" />
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import SourceList from './components/SourceList.vue';
import SourceDetail from './components/SourceDetail.vue';

export default defineComponent({
  name: 'App',
  components: {
    SourceList,
    SourceDetail,
  },
  emits: [],
  data() {
    return {
      selectedSourceId: undefined as string | undefined,
    };
  },
  methods: {
    handleSourceSelected(id: string): void {
      this.selectedSourceId = id;
    },
  },
});
</script>

<style scoped>
.wrapper {
  display: flex;
  height: 100vh;
  width: 100vw;
  flex: 1;
}
</style>
