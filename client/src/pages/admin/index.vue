<template>
  <v-container class="pa-2" min-width="100%">
    <v-card flat :prepend-avatar="botUser?.avatar_url">
      <template v-slot:title>
        <span class="font-weight-light">{{ botUser?.global_name }}</span>
      </template>
      <template v-slot:subtitle>
        <span class="text-sm mr-4">Uptime: <b>{{ uptime }}</b></span>
        <span class="text-sm">Currently part of <b>{{ guilds?.length }}</b> guilds</span>
      </template>
      <template v-slot:text>
        <memory-usage-bar class="pa-3" v-if="chains && botUser" :max-bytes="String(botUser.mem_usage_max)"
          :peak-bytes="String(botUser.mem_usage_peak)" :blocks="computedBlocks" />
      </template>
    </v-card>
    <v-divider class="my-4"></v-divider>
    <div class="d-flex flex-wrap ga-4">
      <v-card width="250" flat v-for="guild in guilds" :key="guild.id"
        :prepend-avatar="'https://cdn.discordapp.com/icons/' + guild.id + '/' + guild.icon + '.webp?size=64'">
        <template v-slot:title>
          <span class="font-weight-light">{{ guild.name }}</span>
        </template>
        <template v-slot:subtitle>
          <span class="text-sm"><b>{{ guild.approximate_member_count }}</b> members</span>
        </template>
        <template v-slot:text>
          <v-row justify="center" class="pa-3 pb-0">
            <span>{{ formatBytes(getChain(guild.id)?.bytes ?? 0) }} / {{ formatBytes(1024 ** 2 *
              (getChain(guild.id)?.max_size_mb ?? 0)) }}</span>
          </v-row>
          <v-row justify="space-between" class="pa-3">
            <v-col cols="12">
              <v-row v-for="(field, key) in getAnalyticsForGuild(guild.id)" :key="key" justify="space-between">
                <span class="text-xs">{{ key }}</span>
                <span class="text-xs">{{ formatNumber(field) }}</span>
              </v-row>
            </v-col>
          </v-row>
        </template>
        <template v-slot:actions>
          <v-row justify="space-between">
            <v-col cols="8">
              <v-tooltip v-slot:activator="{ props }" text="Invite to server" location="bottom">
                <guild-invite-btn :guild-id="guild.id" v-bind="props"></guild-invite-btn>
              </v-tooltip>
              <v-tooltip v-slot:activator="{ props }" text="Check data" location="bottom">
                <v-btn v-bind="props" @click="$router.push(`/data/${guild.id}`)" icon="far fa-file-lines"
                  size="small"></v-btn>
              </v-tooltip>
              <v-tooltip v-slot:activator="{ props }" text="Copy ID" location="bottom">
                <v-btn v-bind="props" @click="copyToClipboard(guild.id)" icon="far fa-copy" size="small"></v-btn>
              </v-tooltip>
            </v-col>
            <v-col cols="3">
              <v-tooltip v-slot:activator="{ props }" text="Leave" location="bottom">
                <v-btn v-bind="props" class="justify-self-end" color="red" icon="fas fa-right-from-bracket"
                  size="small"></v-btn>
              </v-tooltip>
            </v-col>
          </v-row>
        </template>

      </v-card>
    </div>
    <v-snackbar v-model="snackbar.visible" :color="snackbar.color" :timeout="3000" bottom>
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { useGetAllChainsAnalytics } from '@/api/analytics';
import { useGetBotGuilds, useGetBotUser } from '@/api/bot';
import { useAuthStore } from '@/stores/auth';
import { formatBytes, formatNumber, formatTime } from '@/utils/format';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

export default {
  setup() {
    const auth = useAuthStore();
    const botUserQuery = useGetBotUser();
    const botGuildsQuery = useGetBotGuilds(auth.token!);
    const chainsQuery = useGetAllChainsAnalytics(auth.token!);
    const snackbar = ref({
      visible: false,
      message: 'This is a notification!',
      color: "",
    });
    const elapsedSeconds = ref(0);
    botUserQuery.data?.value?.mem_usage_max;
    botUserQuery.data?.value?.mem_usage_peak;
    // Watch for changes in the startup time and update elapsedSeconds accordingly
    watch(() => botUserQuery.data?.value?.startup_timestamp_unix, (newTime) => {
      if (newTime) {
        elapsedSeconds.value = Math.floor(Date.now() / 1000) - newTime;
      }
    }, { immediate: true });

    onMounted(() => {
      const interval = setInterval(() => {
        elapsedSeconds.value += 1;
      }, 1000);

      onBeforeUnmount(() => clearInterval(interval));
    });

    const uptime = computed(() => formatTime(elapsedSeconds.value));

    return {
      botUser: botUserQuery.data,
      inviteLink: import.meta.env.VITE_DISCORD_SERVER_INVITE,
      guilds: botGuildsQuery.data,
      chains: chainsQuery.data,
      uptime,
      snackbar,
      botUserLoading: botUserQuery.isLoading,
      botUserError: botUserQuery.isError,
      guildsLoading: botGuildsQuery.isLoading,
      guildsError: botGuildsQuery.isError,
      chainsLoading: chainsQuery.isLoading,
      chainsError: chainsQuery.isError,
    };
  },
  computed: {
    computedBlocks() {
      return this.chains?.map(c => BigInt(c.bytes)) || [];
    },
  },
  methods: {
    formatBytes,
    formatNumber,
    copyToClipboard(text: string) {
      navigator.clipboard.writeText(text)
        .then(() => {
          this.snackbar.visible = true;
          this.snackbar.message = "Copied to clipboard";
          this.snackbar.color = "success";
        })
        .catch(() => {
          this.snackbar.visible = true;
          this.snackbar.message = "Failed to copy to clipboard";
          this.snackbar.color = "error";
        });
    },
    getChain(guildId: string) {
      return this.chains?.find(c => c.id === guildId);
    },
    getAnalyticsForGuild(guildId: string) {
      const chain = this.chains?.find(c => c.id === guildId);
      if (!chain) return null;
      return {
        Gifs: chain.gifs,
        Images: chain.images,
        Videos: chain.videos,
        Messages: chain.messages,
        Words: chain.words,
        Complexity: chain.complexity_score,
      };
    },
  },
};
</script>