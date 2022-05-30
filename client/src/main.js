import { createApp } from "vue";
import App from "./app/App.vue";
import router from "./router";
import { RouterLink } from "vue-router";
import VueCookies from 'vue-cookies-reactive'

//Mounts App.vue to #app class in index.html
createApp(App).use(router, RouterLink, VueCookies).mount("#app");




