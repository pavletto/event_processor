<template>
  <div class="flex flex-col flex-1 ">
  <main v-if="source" class="flex h-full">
    <div class="flex-1 flex flex-col items-stretch overflow-y-auto">
      <h3 v-if="source" class="my-2 ml-4 text-xl font-bold">Source {{source.id}} details</h3>
      <h3 v-else class="my-2 ml-4 text-lg font-bold">Source  details</h3>
      <video
          :src="videoUrl"
          controls
          @timeupdate="updateCurrentTime"
          ref="videoPlayer"
          class="object-cover w-full max-w-screen-lg m-auto	"
      ></video>
      <RiskGraph
          class="h-30 mt-4"
          :data="source.risk"
          :events="source.events"
          :currentTime="currentTime"
          @seek-video="handleSeekVideo"
          @current-data="handleCurrentRiskData"
      />
      <div class="flex h-36 mt-4">
        <AccelerometerGraph
            class="flex-1 mr-2"
            :data="source.accelerometer"
            :currentTime="currentTime"
            :events="source.events"
            @seek-video="handleSeekVideo"
            @current-data="handleCurrentAccelData"
        />
        <SpeedGraph
            class="flex-1 ml-2"
            :data="source.location"
            :currentTime="currentTime"
            :events="source.events"
            @seek-video="handleSeekVideo"
            @current-data="handleCurrentLocationData"
        />
      </div>

    </div>
    <div class="flex flex-col ">
    <EventList
        :events="source.events"
        :onSeek="handleSeekAndPlayVideo"
        class="w-72 border-l border-gray-300 overflow-y-auto"
    />
      <CreateEventButton
          :videoCurrentTime="currentTime"
          :sourceId="sourceId"
          @event-created="handleEventCreated"
      />
    </div>
  </main>
  <div v-else class="flex justify-center items-center h-full w-full text-gray-500">
    <h4 class="text-lg font-semibold">Select a source from the list.</h4>
  </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue';
import {Source} from '@/types/Source';
import AccelerometerGraph from './AccelerometerGraph.vue';
import SpeedGraph from './SpeedGraph.vue';
import RiskGraph from './RiskGraph.vue';
import {source} from '../api.ts';
import {LocationData} from '../types/LocationData.ts';
import {AccelerometerData} from '../types/AccelerometerData.ts';
import CreateEventButton from './CreateEventButton.vue';
import EventList from './EventList.vue';

export default defineComponent({
  name: 'SourceDetail',
  components: {
    AccelerometerGraph,
    SpeedGraph,
    RiskGraph,
    CreateEventButton,
    EventList,
  },
  props: {
    sourceId: {
      type: String,
      required: false,
    },
  },
  emits: [],
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
            this.source = response.data as Source;
            this.videoUrl = `/api/video/${id}`;
          })
          .catch(error => {
            console.error('Error fetching source data:', error);
          });
    },
    handleSeekVideo(seekTime: number) {
      const videoElement = this.$refs.videoPlayer as HTMLVideoElement;
      if (videoElement) {
        videoElement.currentTime = seekTime;
        videoElement.pause();
      }
    },
    handleSeekAndPlayVideo(seekTime: number) {
      const videoElement = this.$refs.videoPlayer as HTMLVideoElement;
      if (videoElement) {
        videoElement.currentTime = seekTime;
        videoElement.play();
      }
    },
    handleCurrentAccelData(data: AccelerometerData) {
      console.log('Current Accelerometer Data:', data);
    },
    handleCurrentLocationData(data: LocationData) {
      console.log('Current Location Data:', data);
    },
    handleCurrentRiskData(data: LocationData) {
      console.log('Current Risk Data:', data);
    },
    handleEventCreated(newEvent: Event) {
      if (this.source) {
        this.source.events.push(newEvent);
      }
      console.log('Event created:', newEvent);
    },
  },
});
</script>
