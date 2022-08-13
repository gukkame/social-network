<script setup>
import MessageComp from "./MessageBoxGroup.vue";
</script>

<template>
  <div class="chatwindowPerson d-flex" @click="toggleMessenger">
    <div>
      <img
        v-if="data.Image == ``"
        class="bubble3 col"
        src="../../assets/images/profile.svg"
      />
      <div v-else v-bind:id="data.Name" class="bubble3 col-2"></div>
    </div>
    <div class="col messageDetails justify-content-start">
      <div class="col-10 chatwindowUser">{{ data.Name }}</div>
    </div>
    <div v-if="notification" class="notification">
      <i class="bi bi-envelope"></i>
    </div>
  </div>
  <MessageComp :data="data" v-if="messageWindowStatus" />
</template>

<script>
import { ws } from "../../common-js/messages.js";

export default {
  props: {
    data: {
      type: Object,
      required: true,
      notification: false,
    },
  },
  data() {
    return {
      notification: false,
      img: "",
      messageWindowStatus: false,
    };
  },

  mounted: function () {
      let bubble = document.getElementById(this.data.Name);
      if (bubble != null) {
            this.img = "url(http://localhost:8080" + this.data.Image + ")"
            bubble.style.backgroundImage = this.img
        }
      ws.addEventListener("message", (event) => {
          this.handleUsers(event);
      });
  },

  methods: {
    handleUsers(event) {
      let incomingData = JSON.parse(event.data.toString());
      if (
        incomingData.Type == "privateMSG" &&
        !document.querySelector(".messagewindowdiv") &&
        incomingData.Content.Sender == this.data.Username
      ) {
        this.notification = true;
      }

      if (
        incomingData.Type == "privateMSG" &&
        document.querySelector(".messagewindowdiv") &&
        incomingData.Content.Receiver == this.data.Username
      ) {
        this.notification = false;
      }
    },
    toggleMessenger() {
      if (
        document.contains(document.querySelector(".messagewindowdiv")) &&
        this.messageWindowStatus == false
      ) {
        return;
      }
      this.messageWindowStatus = !this.messageWindowStatus;
    },
  },
};
</script>
