<template>
  <div class="border-l border-gray-300 max-w-xs overflow-y-auto flex-1">
    <h3 class="my-2 ml-4 text-lg font-bold">Events</h3>

    <ul v-if="events.length > 0" class="list-none p-0 m-0">
      <li
          v-for="event in events"
          :key="event.id"
          class="p-2.5 border-b border-gray-200"
      >
        <div class="flex justify-between mb-1.5">
          <span class=" text-white px-2 py-0.5 rounded text-sm "  :class="event.tag === 'safe' ? 'bg-green-500' : 'bg-red-500'">
            {{ event.tag }}
          </span>
          <span class="text-xs text-gray-600">
            {{ formatDate(event.start_time * 1000, '{m}:{s}.{ms}') }} - {{ formatDate(event.end_time * 1000, '{m}:{s}.{ms}') }}
          </span>
        </div>
        <p class="text-sm mb-2">{{ event.comment }}</p>
        <div>
          <button
              @click="seekToEvent(event.start_time)"
              class="bg-blue-500 text-white border-none px-2 py-1 rounded text-xs cursor-pointer hover:bg-blue-700"
          >
            Replay
          </button>
        </div>
      </li>
    </ul>
    <p v-else class="text-gray-500 text-center">No events to display.</p>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';
import {Event} from '@/types/Event';
import {formatDate} from "../utils/date.ts";

export default defineComponent({
  name: 'EventList',
  props: {
    events: {
      type: Array as PropType<Event[]>,
      required: true,
    },
    onSeek: {
      type: Function as PropType<(time: number) => void>,
      required: true,
    },
  },
  emits: [],
  methods: {
    formatDate,
    seekToEvent(time: number) {
      this.onSeek(time);
    },
  },
});
</script>
