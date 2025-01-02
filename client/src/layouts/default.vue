<template>
  <v-app>
    <v-main>
      <AppHeader />
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import AppHeader from '@/components/AppHeader.vue';

</script>
