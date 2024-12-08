<template>
  <div class=" py-4 my-auto flex justify-center border-l border-gray-300 ">
    <button @click="handleButtonClick" class="bg-blue-500 text-white px-4 py-2 rounded">
      {{ buttonLabel }}
    </button>

    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
      <div class="bg-white p-6 rounded shadow-lg w-80">
        <h3 class="text-lg font-semibold mb-4">Create Event</h3>
        <form @submit.prevent="submitEvent">
          <div class="mb-4">
            <label for="tag" class="block mb-1">Tag:</label>
            <select v-model="tag" id="tag" required class="w-full border border-gray-300 rounded px-3 py-2">
              <option disabled value="">Please select one</option>
              <option value="danger">Danger</option>
              <option value="safe">Safe</option>
            </select>
          </div>
          <div class="mb-4">
            <label for="comment" class="block mb-1">Comment:</label>
            <textarea v-model="comment" id="comment" rows="3" class="w-full border border-gray-300 rounded px-3 py-2"></textarea>
          </div>
          <div class="flex justify-end">
            <button type="submit" class="bg-green-500 text-white px-4 py-2 rounded">Submit</button>
            <button type="button" @click="cancelEvent" class="bg-gray-500 text-white px-4 py-2 rounded ml-2">Cancel</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'CreateEventButton',
  props: {
    videoCurrentTime: {
      type: Number,
      required: true,
    },
    sourceId: {
      type: String,
      required: false,
    },
  },
  emits: ['event-created'],
  setup(props, { emit }) {
    const isSettingStart = ref<boolean>(false);
    const isSettingEnd = ref<boolean>(false);
    const startTime = ref<number | null>(null);
    const endTime = ref<number | null>(null);
    const showModal = ref<boolean>(false);
    const comment = ref<string>('');
    const tag = ref<string>('');

    const buttonLabel = ref<string>('Create Event');

    const handleButtonClick = () => {
      if (!isSettingStart.value && !isSettingEnd.value) {
        isSettingStart.value = true;
        startTime.value = props.videoCurrentTime;
        buttonLabel.value = 'Set End Time';
      } else if (isSettingStart.value && !isSettingEnd.value) {
        isSettingEnd.value = true;
        endTime.value = props.videoCurrentTime;
        isSettingStart.value = false;
        buttonLabel.value = 'Create Event';
        showModal.value = true;
      }
    };

    const cancelEvent = () => {
      isSettingStart.value = false;
      isSettingEnd.value = false;
      startTime.value = null;
      endTime.value = null;
      comment.value = '';
      tag.value = '';
      showModal.value = false;
      buttonLabel.value = 'Create Event';
    };

    const submitEvent = async () => {
      if (startTime.value === null || endTime.value === null) {
        alert('Start time or end time not set.');
        return;
      }

      if (endTime.value < startTime.value) {
        alert('End time must be after start time.');
        return;
      }

      if (!tag.value) {
        alert('Please select a tag.');
        return;
      }

      const payload = {
        source_id: props.sourceId,
        start_time: startTime.value,
        end_time: endTime.value,
        tag: tag.value,
        comment: comment.value,
      };

      try {
        const response = await axios.post('/api/events', payload);
        emit('event-created', response.data);
        cancelEvent();
      } catch (error) {
        console.error('Failed to create event:', error);
        alert('Failed to create event.');
      }
    };

    return {
      isSettingStart,
      isSettingEnd,
      handleButtonClick,
      showModal,
      comment,
      tag,
      submitEvent,
      cancelEvent,
      buttonLabel,
    };
  },
});
</script>
