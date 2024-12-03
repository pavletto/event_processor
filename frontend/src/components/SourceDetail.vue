<template>
  <main v-if="source">
    <div class="video-container">
      <video :src="videoUrl" controls
             @timeupdate="updateCurrentTime" ref="videoPlayer"></video>
      <RiskGraph class="risk-graph"
                 :data="source.risk"
                 :events="source.events"
                 :currentTime="currentTime"
                 @seek-video="handleSeekVideo"
                 @current-data="handleCurrentRiskData"
      />

    </div>
    <div class="charts-container">
      <AccelerometerGraph class="accel-graph"
                          :data="source.accelerometer"
                          :currentTime="currentTime"
                          @seek-video="handleSeekVideo"
                          @current-data="handleCurrentAccelData"/>
      <SpeedGraph class="speed-graph"
                  :data="source.location"
                  :currentTime="currentTime"
                  @seek-video="handleSeekVideo"
                  @current-data="handleCurrentLocationData"
      />
    </div>
    <CreateEventButton
        :videoCurrentTime="currentTime"
        :sourceId="sourceId"
        @event-created="handleEventCreated"
    />
  </main>
  <div v-else class="empty-details">
    <h4>Select a source from the list.</h4>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {Source} from '../types/Source';
import AccelerometerGraph from './AccelerometerGraph.vue';
import SpeedGraph from './SpeedGraph.vue';
import RiskGraph from './RiskGraph.vue';
import {source} from "../api.ts";
import {LocationData} from "../types/LocationData.ts";
import {AccelerometerData} from "../types/AccelerometerData.ts";
import CreateEventButton from "./CreateEventButton.vue";

export default defineComponent({
  name: 'SourceDetail',
  components: {
    AccelerometerGraph,
    SpeedGraph,
    RiskGraph,
    CreateEventButton,
  },
  props: {
    sourceId: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      source: null as Source | null,
      videoUrl: '',
      currentTime: 0,
    };
  },

  watch: {
    sourceId(newId: string | undefined) {
      if (newId) {
        this.fetchSource(newId);
      }
    },
  },
  mounted() {
    if (this.sourceId) {
      this.fetchSource(this.sourceId);
    }
  },
  methods: {
    updateCurrentTime(event: Event) {
      const video = event.target as HTMLVideoElement;
      this.currentTime = video.currentTime;
    },
    fetchSource(id: string) {
      source(id)
          .then(response => {
            this.source = response.data;
            this.videoUrl = `http://localhost:8080/video/${id}`;
          })
          .catch(error => {
            console.error('Ошибка при получении данных источника:', error);
          });
    }, handleSeekVideo(seekTime: number) {
      const videoElement = this.$refs.videoPlayer as HTMLVideoElement;
      if (videoElement) {
        videoElement.currentTime = seekTime;
        videoElement.play();      }
    }, handleCurrentAccelData(data: AccelerometerData) {
      console.log('Current Accel Data:', data);
    },
    handleCurrentLocationData(data: LocationData) {
      console.log('Current Speed Data:', data);
    },
    handleCurrentRiskData(data: LocationData) {
      console.log('Current Risk Data:', data);
    },
    handleEventCreated() {
        console.log('Event created:', event);
    }

  },
});
</script>

<style scoped>
main {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.video-container {
  max-width: 960px;
}
.risk-graph{
  height: 120px

}
.charts-container {
  display: flex;
  height: 150px

}

video {
  object-fit: cover;
  width: 100%;
}
.empty-details {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 100%;
  color: #666;
}
</style>
