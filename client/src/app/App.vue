<script setup>
import { RouterView } from "vue-router";
import Notifications from "../components/page_parts/Notifications.vue";
import HeaderComp from "../components/page_parts/Header.vue";
import SidebarComp from "../components/page_parts/Sidebar.vue"
import FooterComp from "../components/page_parts/Footer.vue"
import ScrollTopComp from "../components/page_parts/ScrollTop.vue"
import ChatComp from "../components/messages/ChatIcon.vue"

</script>

<template>
  <HeaderComp :LoggedIn="cookieExists" :data="info.Username" />
  <div class="body-app">
    <SidebarComp :LoggedIn="cookieExists" />
    <div class="body-app-main">
      <RouterView :LoggedIn="cookieExists" :data="info.Username" />
      <ScrollTopComp />
      <ChatComp v-if="cookieExists == true" />
      <Notifications :data="data" />
    </div>
    <FooterComp />
  </div>
</template>


<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import { connectToWS, ws } from "../common-js/messages.js";


export default {
  data() {
    return {
      cookieExists: false,
      info: {
        Username: "",
      },
      data: {},
    };
  },

  mounted() {
    if (this.cookieExists == true) {
      ws.addEventListener("message", (event) => {
        this.notifications(event);
      })
    }
  },

  beforeMount() {
    if (this.cookieExists == true) {
      connectToWS();
    }
  },

  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        this.checkAuth()
      },
      { immediate: true }
    )
  },

  methods: {
    checkAuth() {
      let token = document.cookie
      if (token.length === 0) {
        return this.cookieExists = false
      }
      let correctToken = token.split(":")
      this.cookieExists = true
      this.info.Username = correctToken[1]
    },
    notifications(event) {
      let incomingData = JSON.parse(event.data.toString());
      if (incomingData.Type == "followNotif" || incomingData.Type == "GroupInvNotif" || incomingData.Type == "GroupJoinReqNotif" || incomingData.Type == "NewEventNotif") {

        this.data = incomingData
      }
    },
  },
}

</script>
<style>
@import "@/assets/css/base.css";
</style>
