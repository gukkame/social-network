<script setup>
import OneNotificationBoxVueComp from '../OneNotificationBox.vue';
</script>


<template>
  <div class="header">
    <div class="header-item-logo">
      <img alt="kood/ logo" class="logo" src="../../assets/images/logo.svg" />
    </div>

    <div class="header-item-nav">
      <div v-if="LoggedIn">
        <span class="nav-items-0"><span><i :class="{ chatBoxActive: displayNotificationHistory }"
              class="bi bi-bell notificationBell" @click="getNotifications()"></i></span></span>
        <span class="nav-items-0">k<span style="color: #FF9D5A">/{{ data }}</span></span>
      </div>
      <div v-else>
        <button class="nav-button-1">
          <RouterLink class="nav-items-1 stretched-link" to="/login">Log In</RouterLink>
        </button>
        <button class="nav-button-2">
          <RouterLink class="nav-items-2 stretched-link" to="/signup">Sign Up</RouterLink>
        </button>
      </div>
    </div>
  </div>
  <div class="notificationDiv d-flex" v-if="displayNotificationHistory">
    <div class="chatwindow overflow-auto flex-wrap justify-content-space">
      <div class="chatwindowB d-flex flex-column">
        <div v-if="afterDataFetch == null" class="noNotifications d-flex justify-content-center align-items-center">
          No notifications available</div>
        <div v-else v-for="notification in afterDataFetch">
            <OneNotificationBoxVueComp :data="notification" />

        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ws } from '../../common-js/messages';
export default {
  name: 'Header',
  props: {
    LoggedIn: {
      type: Boolean,
      required: true
    },
    data: {
      type: String,
      required: true
    }
  },

  data() {
    return {
      notifcations: null,
      displayNotificationHistory: false
    }
  },

  mounted() {
    if (this.LoggedIn == true) {
      ws.addEventListener("message", (event) => {
        this.handleEvents(event);
      })
    }
  },



  methods: {
    handleEvents(event) {
      let incomingData = JSON.parse(event.data.toString());
      if (incomingData.Type == "ListAllNotif") {
        this.notifcations = incomingData.Notif
      }
    },

    getNotifications() {
      this.displayNotificationHistory = !this.displayNotificationHistory
      let payload = {
        Type: "allNotifications",
      };

      ws.send(JSON.stringify(payload))
    }
  },

  computed: {
      afterDataFetch() {
        if (!this.notifcations) {
          return null
        }
        if (this.notifcations) {
          return this.notifcations
        }
      }
  }
}
</script>


<style>
.notificationBell {
  transition: 0.25s;
}

.notificationBell:hover {
  transform: scale(1.13);
  color: #FF9D5A;
  cursor: pointer;
}

.notificationDiv {
  position: fixed;
  top: 55px;
  right: 10px;
  z-index: 998;
  transition: 0.5s;
  max-width: 260px;
  width: 260px;
  height: 35%;
}

.notifcationWindow {
  box-shadow: 0 2px 6px 0 rgb(218 218 253 / 65%), 0 2px 6px 0 rgb(206 206 238 / 54%);
  border-radius: 6px;
  /*  background-color: blue; */
  width: 100%;
  height: 100%;
  border: solid 1px rgb(219, 219, 219);
  transition: 0.5s;
  background-color: rgb(253, 252, 252);
}

.noNotifications {
  width: 100%;
  height: 50px;
}

</style>
