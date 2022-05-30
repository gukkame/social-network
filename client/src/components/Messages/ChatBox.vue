<script setup>
import ChatBoxUserCompO from "./ChatBoxUserO.vue"
</script>

<template>
    <div class="chatwindowdiv  d-flex">
        <div class="chatwindow overflow-auto flex-wrap justify-content-space">
            <div class="chatwindowH d-flex justify-content-center align-items-center">
                Messenger
            </div>
            <div class="chatwindowB d-flex flex-column">
                <div class="chatwindowA  d-flex justify-content-center align-items-center flex-row">Active users
                    ({{ amountofClients }})
                </div>
                <div v-for="clients in { realAmountofUsers }">
                    <div v-for="username in clients" :key="username.Username">
                        <ChatBoxUserCompO :data="username" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import _ from 'lodash';
import { connectToWS, ws } from "../../common-js/messages.js"

export default {

    data() {
        return {
            users: [],
            timer: null,
        }
    },

    beforeMount() {
        connectToWS()
    },

    mounted: function () {
        ws.addEventListener('message', (event) => { this.handleUsers(event) });
        this.timer = window.setInterval(() => {
            this.getAllUser();
        }, 1500)
    },

    methods: {
        getAllUser() {
            let payload = {
                "Type": "allClients"
            }
            if (ws.readyState === WebSocket.CLOSED) {
                clearInterval(this.timer)
                return
            }
            ws.send(JSON.stringify(payload));
        },

        handleUsers(event) {
            let incomingData = JSON.parse(event.data.toString());
            if (incomingData.Type == "allClients") {
                this.users = incomingData.Clients
            }
        }
    },

    computed: {
        amountofClients() {
            if (this.users.Clients == undefined) {
                return 0
            } else {
                return this.users.Clients.length
            }
        },
        realAmountofUsers() {
            if (this.users.Clients == undefined) {
                return
            }
            let correctToken = document.cookie.split(":")
            let array = this.users.Clients

            let newArray = array.filter((item) => item.Username != correctToken[1]);
            if (newArray.some(e => e.LastMessage.length != 0)) {
                let times = newArray.filter((item) => item.LastMessage != "no date");

                times.sort(function compare(a, b) {
                    var dateA = new Date(a.LastMessage);
                    var dateB = new Date(b.LastMessage);
                    return dateA - dateB;
                });
                times.reverse()
                let timesNull = newArray.filter((item) => item.LastMessage == "no date");
                newArray = [...times, ...timesNull];

            } else {
                newArray.sort((a, b) => {
                    let fa = a.Username.toLowerCase(),
                        fb = b.Username.toLowerCase();

                    if (fa < fb) {
                        return -1;
                    }
                    if (fa > fb) {
                        return 1;
                    }
                    return 0;
                });
            }

            return this.users.Clients = newArray
        }
    }
}
</script>



<style>
.chatwindowA {
    border-bottom: solid 1px rgb(219, 219, 219);
}

.chatwindowO {
    margin-top: 10px;
    border-bottom: solid 1px rgb(219, 219, 219);
}

.activityDot {
    border-radius: 50%;
    width: 15px;
    height: 15px;
    background-color: #00FF7F;
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