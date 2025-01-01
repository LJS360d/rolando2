<template>
  <v-app-bar app>
    <v-app-bar-nav-icon icon="fas fa-bars" @click="drawer = !drawer" />

    <v-toolbar-title></v-toolbar-title>

    <v-spacer />

    <v-btn class="ml-3" @click="loginWithDiscord" color="#5865F2">
      <v-icon class="mr-2" icon="fa-brands fa-discord"></v-icon>
      Login with Discord
    </v-btn>
  </v-app-bar>

  <v-navigation-drawer v-model="drawer" app width="250">
    <v-list>
      <v-list-item v-for="(item, index) in menuItems" :key="index" @click="onMenuItemClick(item)">
        <v-list-item-content>{{ item }}</v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";

export default defineComponent({
  name: "AppHeader",
  setup() {
    const drawer = ref(false);
    const menuItems = ref(["Home", "About", "Settings"]);

    const loginWithDiscord = () => {
      window.location.href = "https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&redirect_uri=YOUR_REDIRECT_URI&response_type=code&scope=identify";
    };

    const onMenuItemClick = (item: string) => {
      console.log(`Clicked: ${item}`);
    };

    return { drawer, menuItems, loginWithDiscord, onMenuItemClick };
  },
});
</script>

<style scoped>
.v-btn {
  color: white;
}
</style>
