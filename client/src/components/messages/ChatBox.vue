<script setup>
import ChatBoxUserCompO from "./ChatBoxUserO.vue";
import ChatBoxUserCompOf from "./ChatBoxUserOf.vue";
import ChatBoxGroups from "./ChatBoxGroups.vue";
</script>

<template>
  <div class="chatwindowdiv d-flex">
    <div class="chatwindow overflow-auto flex-wrap justify-content-space">
      <div class="chatwindowH d-flex justify-content-center align-items-center">
        Messenger
      </div>
      <div class="chatwindowB d-flex flex-column">
        <div class="chatwindowA d-flex justify-content-center align-items-center flex-row">
          Active users ({{ onlineUsers }})
        </div>
        <div v-for="clients in { listOnlineUsers }">
          <div v-for="username in clients" :key="username.Username">
            <ChatBoxUserCompO :data="username" />
          </div>
        </div>
        <div class="chatwindowA d-flex justify-content-center align-items-center flex-row">
          Offline ({{ offlineUsers }})
        </div>


        <div v-for="clients in { listOfflineUsers }">
          <div v-for="username in clients" :key="username.Username">
            <ChatBoxUserCompOf :data="username" />
          </div>
        </div>
        <div class="chatwindowA d-flex justify-content-center align-items-center flex-row">
          Groups ({{ groups }})
        </div>
        <div v-for="group in { listGroups }">
          <div v-for="name in group" :key="name.Name">

            <ChatBoxGroups :data="name" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import _, { intersection } from "lodash";
import { ws } from "../../common-js/messages.js";


export default {
  data() {
    return {
      onlineUsers: 0,
      offlineUsers: 0,
      groups: 0,
      users: [],
      allgroups: [],
      timer: null,
    };
  },
  props: {
    data: {
      type: Object,

    },
  },

  mounted: function () {
    ws.addEventListener("message", (event) => {
      this.handleUsers(event);
    });

    this.timer = window.setInterval(() => {
      this.getAllUser();
      this.getAllGroups();
    }, 1500);
  },

  methods: {
    getAllUser() {
      let payload = {
        Type: "allClients",
      };
      if (ws.readyState === WebSocket.CLOSED) {
        clearInterval(this.timer);
        return;
      }
      ws.send(JSON.stringify(payload));
    },
    getAllGroups() {
      let payload = {
        Type: "allGroups",
      };
      if (ws.readyState === WebSocket.CLOSED) {
        clearInterval(this.timer);
        return;
      }
      ws.send(JSON.stringify(payload));
    },

    handleUsers(event) {
      let incomingData = JSON.parse(event.data.toString());
      if (incomingData.Type == "allClients") {
        this.users = incomingData.Clients;
      }
      if (incomingData.Type == "allGroups") {
        this.allgroups = incomingData.Groups;
      }
    },
  },

  computed: {
    listOnlineUsers() {
      if (this.users.Clients == undefined) {
        return;
      }
      let correctToken = document.cookie.split(":");
      if (this.users != null) {
        let filtered = this.users.Clients.filter(
          (user) => user.Username != correctToken[1]
        );
        let filtered1 = filtered.filter((user) => user.Status == "1");
        this.onlineUsers = filtered1.length;
        let allUsersO = _.orderBy(filtered1, ["LastMessage", "Username"], ["asc", "asc"]);
        return allUsersO;
      }
    },
    listOfflineUsers() {
      if (this.users.Clients == undefined) {
        return;
      }
      let correctToken = document.cookie.split(":");
      if (this.users != null) {
        let filtered = this.users.Clients.filter(
          (user) => user.Username != correctToken[1]
        );
        let filtered1 = filtered.filter((user) => user.Status == "0");
        this.offlineUsers = filtered1.length;
        let allUsersO = _.orderBy(filtered1, ["LastMessage", "Username"], ["asc", "asc"]);
        return allUsersO;
      }
    },
    listGroups() {
      if (this.allgroups.Groups == undefined) {
        return;
      }
      if (this.allgroups != null) {
        this.groups = this.allgroups.Groups.length;
        return this.allgroups.Groups;
      }
    },
  },
};
</script>

<style>
.chatwindowA {
  border-bottom: solid 1px rgb(219, 219, 219);
}

.chatwindowO {
  margin-top: 10px;
  border-bottom: solid 1px rgb(219, 219, 29);
}

.activityDot {
  border-radius: 50%;
  width: 15px;
  height: 15px;
  background-color: #00ff7f;
  border: solid 1px rgb(219, 219, 219);
}

.chatwindowdiv {
  position: fixed;
  bottom: 69px;
  right: 48px;
  z-index: 999;
  transition: 0.5s;
  max-width: 220px;
  width: 220px;
  height: 57%;
}

.chatwindow {
  box-shadow: 0 2px 6px 0 rgb(218 218 253 / 65%), 0 2px 6px 0 rgb(206 206 238 / 54%);
  border-radius: 6px;
  /*  background-color: blue; */
  width: 100%;
  height: 100%;
  border: solid 1px rgb(219, 219, 219);
  transition: 0.5s;
  background-color: rgb(253, 252, 252);
}

.chatwindowH {
  width: 100%;
  height: 50px;
  /*  background-color: green; */
  border-radius: 6px 6px 0px 0px;
  font-size: 24px;
  border-bottom: solid 1px rgb(219, 219, 219);
  margin-bottom: 10px;
}

.chatwindowB {
  width: 100%;
}
</style>
