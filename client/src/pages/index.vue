<template>
  <v-container class="pa-4">
    <v-row v-if="!isLoading">
      <!-- Hero Section -->
      <v-col cols="12" lg="6">
        <h1 class="text-h3 font-weight-bold mb-4">
          Meet {{ botUser?.username }}
        </h1>
        <p class="text-body-1 mb-6">
          The learning AI that makes your Discord server dangerously fun.
        </p>
        <v-img :src="botUser?.avatar_url" alt="AI Discord Bot" class="rounded-circle" height="200px" width="200px" />
        <p class="text-body-2 text-gray-500 mt-4">
          Warning: This AI bot learns from messages and can generate NSFW content.
          <br />User discretion is advised.
        </p>
        <h6 class="mt-6">
          Currently part of {{ "100+" }} servers
        </h6>
        <v-btn :href="botUser?.invite_url" class="mt-4 bg-discord" elevation="2">
          <v-icon left class="mr-2" icon="fa-brands fa-discord"></v-icon>
          Invite to your server
        </v-btn>
      </v-col>

      <!-- Commands Section -->
      <v-col cols="12" lg="6">
        <v-row class="d-flex align-center">
          <v-col cols="12" sm="6">
            <h2 class="text-h5 font-weight-semibold">Commands</h2>
          </v-col>
          <v-col cols="12" sm="6" class="text-right">
            <v-btn icon @click="shuffleList()" title="Shuffle the commands, why not?">
              <v-icon icon="fa-solid fa-shuffle"></v-icon>
            </v-btn>
          </v-col>
        </v-row>

        <v-list v-auto-animate ref="commandsList" id="commandsList" class="text-lg">
          <v-list-item v-for="(command, index) in botUser?.slash_commands" class="mb-2" :key="index">
            <v-list-item-title class="font-weight-bold text-h6">
              /{{ command.name }}
            </v-list-item-title>
            <v-list-item-subtitle class="text-body-2">
              {{ command.description }}
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
    <v-progress-circular indeterminate v-else-if="!isError" color="primary" size="64" />
    <v-alert v-else type="error" class="text-body-2">
      Oops
    </v-alert>
  </v-container>
</template>

<script lang="ts">
import { useQuery } from '@tanstack/vue-query';
import { defineComponent } from 'vue';

export interface BotUser {
  accent_color: number;
  avatar_url: string;
  discriminator: string;
  global_name: string;
  id: string;
  invite_url: string;
  slash_commands: SlashCommand[];
  username: string;
  verified: boolean;
}

export interface SlashCommand {
  id: string;
  application_id: string;
  version: string;
  type: number;
  name: string;
  dm_permission: boolean;
  nsfw: boolean;
  description: string;
  options: Option[] | null;
}

export interface Option {
  type: number;
  name: string;
  description: string;
  channel_types: null;
  required: boolean;
  options: null;
  autocomplete: boolean;
  choices: null;
}


export default defineComponent({
  name: 'IndexPage',
  setup() {
    const { data: botUser, isLoading, isError } = useQuery({
      queryKey: ["/bot/user"],
      queryFn: async () => {
        const response = await fetch(`/api/bot/user`);
        if (!response.ok) throw new Error("Failed to fetch bot user");
        return response.json() as Promise<BotUser>;
      },
    });

    return {
      botUser,
      isLoading,
      isError
    };
  },
  methods: {
    shuffleList() {
      const list = document.getElementById('commandsList')!;
      for (let i = list.children.length; i >= 0; i--) {
        list?.appendChild(list.children[Math.random() * i | 0]);
      }
    }
  }
});
</script>