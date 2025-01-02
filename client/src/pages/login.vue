<template>
  <v-col cols="12" centered>
    <v-progress-circular color="primary" indeterminate :size="64" :width="6"></v-progress-circular>
  </v-col>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const router = useRouter();
const fragment = window.location.hash.substring(1);
if (!fragment) {
  router.replace('/');
}

const params = new URLSearchParams(fragment);
const accessToken = params.get('access_token');
const authStore = useAuthStore();

if (accessToken) {
  fetch('https://discord.com/api/v10/users/@me', {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  })
    .then(async (res) => {
      if (res.ok) {
        const body = await res.json();
        authStore.setAuth(accessToken, body);
        router.replace('/'); // Redirect after login
      } else {
        router.replace('/'); // Redirect on error
      }
    })
    .catch(() => {
      router.replace('/'); // Redirect on fetch failure
    });
} else {
  router.replace('/');
}
</script>
