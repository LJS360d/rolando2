<template>
  <v-container>
    <v-row class="py-0">
      <v-col class="pa-0" cols="12">
        <!-- Check if data is available -->
        <v-progress-linear v-if="isDataAvailable" :value="100" height="20" color="grey lighten-2" class="rounded-lg">
          <!-- Memory blocks -->
          <div v-for="(block, index) in blocks" :key="index" class="memory-block" :style="{
            width: (block / max) * BigInt(100) + '%',
            backgroundColor: getBlockColor(index),
          }" />
          <!-- Peak line -->
          <div class="peak-line" :style="{ left: (peak / max) * BigInt(100) + '%' }"></div>
        </v-progress-linear>

        <!-- Skeleton loader when data is unavailable -->
        <v-progress-linear v-else indeterminate height="20" color="grey lighten-3" class="rounded-lg" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';

export default defineComponent({
  name: 'MemoryUsageBar',
  props: {
    maxBytes: {
      type: String as PropType<string>, // Accept as string
      required: false,
    },
    peakBytes: {
      type: String as PropType<string>, // Accept as string
      required: false,
    },
    blocks: {
      type: Array as PropType<bigint[]>,
      required: false,
    },
  },
  data() {
    return {
      max: this.maxBytes ? BigInt(this.maxBytes) : BigInt(0),
      peak: this.peakBytes ? BigInt(this.peakBytes) : BigInt(0),
    };
  },
  computed: {
    isDataAvailable(): boolean {
      return !!this.maxBytes && !!this.peakBytes && !!this.blocks?.length;
    },
  },
  methods: {
    getBlockColor(index: number): string {
      const colors = ['#4caf50', '#ffeb3b', '#f44336', '#2196f3', '#9c27b0'];
      return colors[index % colors.length];
    },
  },
});
</script>

<style scoped>
.memory-block {
  height: 100%;
  transition: all 0.3s ease;
}

.peak-line {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: red;
  z-index: 10;
}
</style>
