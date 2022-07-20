<script setup>
import { RouterView } from "vue-router";
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
      <ChatComp v-if="cookieExists == true"/>
    </div>
    <FooterComp />
  </div>
</template>


<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"

export default {
  data() {
    return {
      cookieExists: false,
      info: {
        Username: "",
      },
    };
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
    }
  },
}

</script>
<style>
@import "@/assets/css/base.css";
</style>
