<template>
  <v-container class="pa-4">
    <v-row>
      <!-- Hero Section -->
      <v-col cols="12" lg="6">
        <h1 class="text-h3 font-weight-bold mb-4">Meet {{ client.user.displayName }}</h1>
        <p class="text-body-1 mb-6">
          The learning AI that makes your Discord server dangerously fun.
        </p>
        <v-img :src="client.user.avatarURL" alt="AI Discord Bot" class="rounded-circle" height="200px"
          width="200px"></v-img>
        <p class="text-body-2 text-gray-500 mt-4">
          Warning: This AI bot learns from messages and can generate NSFW content.
          <br />User discretion is advised.
        </p>
        <h6 class="mt-6">Currently part of {{ guildsCount }} servers</h6>
        <v-btn :href="inviteLink" color="#5865F2" class="mt-4 " elevation="2">
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
          <v-list-item v-for="(command, index) in commands" :key="index">
            <v-list-item-content>
              <v-list-item-title class="font-weight-bold text-h6">
                /{{ command.name }}
              </v-list-item-title>
              <v-list-item-subtitle class="text-body-2">{{ command.description }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>

  </v-container>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
const commandsList = useTemplateRef<HTMLLIElement>("commandsList")
export default defineComponent({
  name: 'IndexPage',
  data() {
    return {
      client: {
        user: {
          displayName: 'Rolando',
          avatarURL: 'https://cdn.discordapp.com/avatars/1100311970428756039/d5ffa17a735e6da1f34a39ffbaddd826.webp',
        },
      },
      guildsCount: 100,
      inviteLink: 'https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&redirect_uri=YOUR_REDIRECT_URI&response_type=code&scope=identify',
      commands: [
        { name: 'help', description: 'Displays help information' },
        { name: 'ping', description: 'Checks the bot\'s response time' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
        { name: 'info', description: 'Shows information about the bot' },
      ],
      commandsList,
    };
  },
  methods: {
    shuffleList() {
      const list = document.getElementById('commandsList')!;
      for (let i = list.children.length; i >= 0; i--) {
        list?.appendChild(list.children[Math.random() * i | 0]);
      }
    },
  },
});
</script>
