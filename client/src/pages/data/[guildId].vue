<template>
  <v-container>
    <template v-if="!isLoading && !isError && data">
      <v-card flat :prepend-avatar="guildIconUrl(data.guild.id, data.guild.icon)">
        <template v-slot:title>
          <span class="font-weight-light">{{ data.guild.name }}</span>
        </template>
        <template v-slot:subtitle>
          <span class="text-sm"><b>{{ data.guild.members }}</b> members</span>
        </template>
        <template v-slot:text v-if="!!chain">
          <v-row justify="center" class="pa-3 pb-0">
            <span>{{ formatBytes(chain?.bytes ?? 0) }} / {{ formatBytes(1024 ** 2 *
              (chain?.max_size_mb ?? 0)) }}</span>
          </v-row>
          <v-row justify="space-between" class="pa-3">
            <v-col cols="12">
              <v-row v-for="(field, key) in getChainAnalytics()" :key="key" justify="space-between">
                <span class="text-xs">{{ key }}</span>
                <span class="text-xs">{{ formatNumber(field) }}</span>
              </v-row>
            </v-col>
          </v-row>
        </template>
      </v-card>
      <v-divider class="my-4" />
      <h2 class="mb-4">All the learned messages</h2>
      <v-list density="compact">
        <v-list-item v-for="(text, i) in data.messages" :key="i" :class="{ 'bg-dark': i % 2 !== 0 }">
          <template v-slot:title>
            <p class="message-content" v-html="renderedMessages[i]" />
          </template>
        </v-list-item>
      </v-list>
    </template>
    <v-skeleton-loader v-else-if="!isError" type="card-avatar" />
    <v-alert v-else type="error" class="text-body-2">
      Oops, big error occured, please report it to the creator on <a :href="discordServerInvite">the discord</a>
    </v-alert>
  </v-container>
</template>

<style scoped>
.message-content {
  white-space: normal;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.bg-dark {
  background-color: #09090989;
}
</style>

<script lang="ts">
import { useGetChainAnalytics } from '@/api/analytics';
import { useGetGuildData } from '@/api/data';
import { useAuthStore } from '@/stores/auth';
import { formatBytes, formatNumber, guildIconUrl } from '@/utils/format';
import DOMPurify from 'dompurify';
import { marked } from 'marked';
import { defineComponent, ref } from 'vue';
import { useRoute } from 'vue-router';

export default defineComponent({
  setup() {
    const route = useRoute();
    const auth = useAuthStore();
    const guildId = (route.params as { guildId: string }).guildId;
    const guildDataQuery = useGetGuildData(auth.token!, guildId);
    const chainQuery = useGetChainAnalytics(auth.token!, guildId);
    const renderMarkdown = (text: string) => {
      const rawHtml = marked(text, { async: false });
      return DOMPurify.sanitize(rawHtml);
    };

    return {
      discordServerInvite: import.meta.env.VITE_DISCORD_SERVER_INVITE,
      data: guildDataQuery.data,
      chain: chainQuery.data,
      isLoading: guildDataQuery.isLoading,
      isError: guildDataQuery.isError,
      renderedMessages: computed(() => guildDataQuery.data?.value?.messages.map(renderMarkdown) ?? []),
    };
  },
  methods: {
    formatBytes,
    formatNumber,
    guildIconUrl,
    getChainAnalytics() {
      const chain = this.chain;
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
  }
});
</script>
