<template>
  <v-container>
    <v-form ref="broadcastForm" @submit.prevent="onFormSubmit" class="d-flex justify-around">
      <v-row>
        <v-col cols="12">
          <v-list dense class="grid grid-cols-7 gap-2">
            <v-list-item v-for="guild in guilds" :key="guild.id" :title="guild.name" :data-name="guild.name"
              class="bg-grey lighten-3 w-32 rounded-lg p-4 flex align-center">
              <v-checkbox v-model="selectedGuilds[guild.id]" class="toggle-primary" />
              <v-avatar size="48" class="ml-2">
                <img :src="guild.iconURL || 'assets/noimage.svg'" alt="Guild Icon" />
              </v-avatar>
              <v-list-item-title class="text-truncate">{{ guild.name }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-col>
      </v-row>

      <v-col class="d-flex flex-column sticky">
        <span>Guilds: <b>{{ selectedGuildCount }}</b></span>

        <v-row class="py-4 align-center">
          <v-btn small outlined color="primary" @click="toggleSelection">
            Toggle Select
          </v-btn>

          <v-text-field v-model="searchText" label="Search" outlined dense prepend-icon="mdi-magnify"
            @input="searchGuild" />
        </v-row>

        <v-textarea v-model="message" outlined dense label="Message" rows="6" class="mb-4"></v-textarea>

        <v-switch v-model="keepAfterSubmit" label="Keep after submit" inset></v-switch>

        <v-btn type="submit" color="primary">Submit</v-btn>
      </v-col>
    </v-form>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      guilds: [
        { id: "1", name: "Guild One", iconURL: null },
        { id: "2", name: "Guild Two", iconURL: "assets/noimage.svg" },
      ],
      selectedGuilds: {},
      message: "",
      keepAfterSubmit: false,
      searchText: "",
    };
  },
  computed: {
    selectedGuildCount() {
      return Object.values(this.selectedGuilds).filter((val) => val).length;
    },
  },
  methods: {
    onFormSubmit() {
      if (!this.message.trim()) {
        this.showSnackbar("Broadcast message is empty", "error");
        return;
      }

      const selectedCount = this.selectedGuildCount;
      if (selectedCount === 0) {
        this.showSnackbar("No guilds are selected", "error");
        return;
      }

      this.showSnackbar(
        `Sent "${this.message}" to ${selectedCount} guilds`,
        "success"
      );

      if (!this.keepAfterSubmit) {
        this.message = "";
      }
    },
    toggleSelection() {
      const allSelected = Object.values(this.selectedGuilds).every((val) => val);
      this.guilds.forEach((guild) => {
        this.$set(this.selectedGuilds, guild.id, !allSelected);
      });
    },
    searchGuild() {
      const searchUpper = this.searchText.trim().toUpperCase();
      this.guilds.forEach((guild) => {
        const nameUpper = guild.name.toUpperCase();
        guild.hidden = !nameUpper.includes(searchUpper);
      });
    },
    showSnackbar(message, type) {
      console.log(`${type}: ${message}`);
    },
  },
};
</script>

<style scoped>
.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
