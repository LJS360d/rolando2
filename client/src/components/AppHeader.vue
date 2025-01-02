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
    <v-btn v-else class="ml-3 bg-discord" :href="OAUTH2_URL" target="_self">
      <v-icon class="mr-2" icon="fa-brands fa-discord"></v-icon>
      Login with Discord
    </v-btn>
  </v-app-bar>

  <v-navigation-drawer v-model="drawer" app width="250">
    <v-list>
      <v-list-item :href="item.href" v-for="(item, index) in links" :key="index">
        <span>
          <v-icon :v-if="item.icon" class="mx-2" size="16" :icon="item.icon"></v-icon>
          <span class="text-body-2" :href="item.href">{{ item.name }}</span>
        </span>
      </v-list-item>
      <v-divider :v-if="isOwner"></v-divider>
      <v-list-item :v-if="isOwner" :href="item.href" v-for="(item, index) in adminLinks" :key="index">
        <span>
          <v-icon :v-if="item.icon" class="mx-2" size="16" :icon="item.icon"></v-icon>
          <span class="text-body-2" :href="item.href">{{ item.name }}</span>
        </span>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from "vue";
import { useAuthStore } from "@/stores/auth";

interface MenuItem {
  name: string;
  href: string;
  icon?: string;
}

const OAUTH2_URL: string = import.meta.env.VITE_OAUTH2_URL;

export default defineComponent({
  name: "AppHeader",
  setup() {
    const authStore = useAuthStore();
    const drawer = ref(false);
    const links = ref<MenuItem[]>([
      { name: "Home", icon: "fas fa-home", href: "/" },
      { name: "Privacy Policy", icon: "fas fa-user-shield", href: "/privacy-policy" },
      { name: "Terms of Service", icon: "fas fa-file-signature", href: "/terms-of-service" },
    ]);
    const adminLinks = ref<MenuItem[]>([
      { name: "Admin Panel", icon: "fas fa-cogs", href: "/admin" },
      { name: "Message Broadcast", icon: "fas fa-bullhorn", href: "/admin/broadcast" },
    ]);

    const isLoggedIn = computed(() => !!authStore.token);
    const userAvatarUrl = computed(() =>
      authStore.user?.avatar
        ? `https://cdn.discordapp.com/avatars/${authStore.user.id}/${authStore.user.avatar}.png?size=40`
        : "https://cdn.discordapp.com/embed/avatars/0.png"
    );

    const userDisplayName = computed(() =>
      authStore.user?.username || authStore.user?.global_name || ""
    );
    // TODO call server for isOwner
    return {
      drawer,
      links,
      adminLinks,
      OAUTH2_URL,
      logout: authStore.logout,
      isLoggedIn,
      isOwner: isLoggedIn,
      userAvatarUrl,
      userDisplayName,
    };
  },
});
</script>
