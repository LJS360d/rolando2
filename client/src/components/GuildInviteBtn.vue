<template>
  <v-btn v-bind="props" :icon="buttonIcon" :color="buttonColor" :href="inviteLink" target="_blank" size="small"
    @click="getInvite"></v-btn>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth';
import { ref, defineProps } from 'vue';

const props = defineProps({
  guildId: {
    type: String,
    required: true,
  },
});
const auth = useAuthStore();
const buttonIcon = ref('fas fa-door-closed'); // Initially the door-closed icon
const buttonColor = ref(''); // Default color (no color)
const inviteLink = ref(''); // Initially no invite link

const getInvite = async () => {
  try {
    // Construct the API URL using the passed guildId prop
    const response = await fetch(`/api/bot/guilds/${props.guildId}/invite`,
      {
        headers: {
          Authorization: auth.token
        }
      }
    );

    // Check if the response is successful
    if (!response.ok) {
      throw new Error('Failed to fetch invite');
    }

    const data = await response.json();

    if (data && data.invite) {
      // Update button to reflect the invite
      buttonIcon.value = 'fas fa-door-open';
      buttonColor.value = 'green'; // Set color to green
      inviteLink.value = data.invite; // Set the href to the invite link
    }
  } catch (error) {
    console.error('Error fetching invite:', error);
    buttonColor.value = 'red';
  }
};
</script>
