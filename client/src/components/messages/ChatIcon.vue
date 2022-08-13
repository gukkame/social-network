<script setup>
import ChatBoxComp from "./ChatBox.vue";
import { connectToWS, ws } from "../../common-js/messages.js";
</script>

<template>
  <ChatBoxComp v-if="chatWindowStatus" />
  <div class="chatbuttondiv">
    <div class>
      <i
        :class="{ chatBoxActive: chatWindowStatus }"
        @click="toggleChatbox"
        class="chatbutton bi bi-chat-dots"
      ></i>
    </div>
  </div>
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import "bootstrap/dist/js/bootstrap.js";

export default {
  data() {
    return {
      chatWindowStatus: false,
    };
  },

  methods: {
    async toggleChatbox() {
      this.chatWindowStatus = !this.chatWindowStatus;
    },
  },
};
</script>

<style>
.chatBoxActive {
  color: #ff9d5a;
}

.chatbutton {
  background-color: rgb(253, 252, 252);
  transition: 0.25s;
  font-size: 45px;
  background: rgba(0, 0, 0, 0);
}

.chatbutton:hover {
  transform: scale(1.13);
  color: #ff9d5a;
  cursor: pointer;
}

.chatbuttondiv {
  position: fixed;
  bottom: 15px;
  right: 48px;
  z-index: 1;
  opacity: 1;
  transition: 0.5s;
  transition: opacity 0.3s linear;
  background: rgba(0, 0, 0, 0);
}
</style>