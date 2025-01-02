<template>
  <v-app-bar app>
    <v-app-bar-nav-icon icon="fas fa-bars" @click="drawer = !drawer" />

    <v-toolbar-title></v-toolbar-title>

    <v-spacer />

    <div class="cursor-pointer pa-1 mr-2 rounded" title="Click to logout" @click="logout" v-if="isLoggedIn">
      <v-avatar class="mr-3" :size="40">
        <img :src="userAvatarUrl" alt="User Avatar" />
      </v-avatar>
      <span>{{ userDisplayName }}</span>
    </div>
    <v-btn v-else class="ml-3" @click="loginWithDiscord" color="#5865F2">
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
import { defineComponent, ref, computed } from "vue";
import { useAuthStore } from "@/stores/auth";

const OAUTH2_URL =
  "https://discord.com/api/oauth2/authorize?client_id=1100311970428756039&response_type=token&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Flogin&scope=identify";

export default defineComponent({
  name: "AppHeader",
  setup() {
    const authStore = useAuthStore();
    const drawer = ref(false);
    const menuItems = ref(["Home", "About", "Settings"]);

    const loginWithDiscord = () => {
      window.location.href = OAUTH2_URL;
    };

    const logout = () => {
      authStore.logout();
    }

    const onMenuItemClick = (item: string) => {
      console.log(`Clicked: ${item}`);
    };

    const isLoggedIn = computed(() => !!authStore.token);
    const userAvatarUrl = computed(() =>
      authStore.user?.avatar
        ? `https://cdn.discordapp.com/avatars/${authStore.user.id}/${authStore.user.avatar}.png?size=40`
        : "https://cdn.discordapp.com/embed/avatars/0.png"
    );
    const userDisplayName = computed(() =>
      authStore.user?.username || authStore.user?.global_name || "User"
    );

    return {
      drawer,
      menuItems,
      loginWithDiscord,
      logout,
      onMenuItemClick,
      isLoggedIn,
      userAvatarUrl,
      userDisplayName,
    };
  },
});
</script>

<style scoped>
.v-btn {
  color: white;
}

.v-avatar {
  cursor: pointer;
}
</style>
