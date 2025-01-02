/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from './vuetify'
import pinia from '../stores'
import router from '../router'
import piniaPersistedState from 'pinia-plugin-persistedstate';

// Types
import type { App } from 'vue'
import { autoAnimatePlugin } from '@formkit/auto-animate/vue'

export function registerPlugins(app: App) {
  pinia.use(piniaPersistedState)
  app
    .use(autoAnimatePlugin)
    .use(vuetify)
    .use(router)
    .use(pinia)
}
