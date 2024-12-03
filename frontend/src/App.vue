<template>
  <div class="wrapper">
      <SourceList @source-selected="handleSourceSelected"/>
      <SourceDetail :sourceId="selectedSourceId"/>
  </div>
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
  data() {
    return {
      selectedSourceId: null as string | null,
    };
  },
  methods: {
    handleSourceSelected(id: string) {
      this.selectedSourceId = id;
    },
  },
});
</script>
<style scoped>
.wrapper {
  display: flex;
  height: 100vh;
  flex: 1;
}



</style>
