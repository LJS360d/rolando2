/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import 'vuetify/styles';

// Composables
import '@fortawesome/fontawesome-free/css/all.css';
import { createVuetify } from 'vuetify';
import { fa } from 'vuetify/iconsets/fa';

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'dark',
  },
  icons: {
    defaultSet: 'fa',
    sets: {
      fa,
    },
  },
})
