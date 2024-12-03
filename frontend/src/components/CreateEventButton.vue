<template>
  <div>
    <button @click="handleButtonClick" class="btn btn-primary">
      {{ buttonLabel }}
    </button>

    <!-- Modal -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h3>Create Event</h3>
        <form @submit.prevent="submitEvent">
          <div class="form-group">
            <label for="tag">Tag:</label>
            <select v-model="tag" id="tag" required>
              <option disabled value="">Please select one</option>
              <option value="danger">Danger</option>
              <option value="safe">Safe</option>
            </select>
          </div>
          <div class="form-group">
            <label for="comment">Comment:</label>
            <textarea v-model="comment" id="comment" rows="3"></textarea>
          </div>
          <div class="modal-actions">
            <button type="submit" class="btn btn-success">Submit</button>
            <button type="button" @click="cancelEvent" class="btn btn-secondary">Cancel</button>
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
      required: true,
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
        start_time: startTime.value * 1000,
        end_time: endTime.value * 1000,
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

<style scoped>
.btn {
  padding: 8px 16px;
  margin: 5px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  width: 300px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
}

.form-group select,
.form-group textarea {
  width: 100%;
  padding: 5px;
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
}

.modal-actions .btn {
  margin-left: 10px;
}
</style>
