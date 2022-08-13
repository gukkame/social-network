<script setup>
import MessageComp from "./MessageBox.vue"
</script>

<template>
    <div class="chatwindowPerson d-flex" @click="toggleMessenger">
        <div>
            <img v-if="data.Avatar_image == ``" class="bubble3 col" src="../../assets/images/profile.svg" />
            <div v-else v-bind:id="data.Username" class="bubble3 col-2"></div>
        </div>
        <div class="col messageDetails justify-content-start">
            <div class="col-10 chatwindowUser">{{ data.Username }}</div>
            <div class="col chatwindowTime">Active</div>
        </div>
        <div v-if="notification" class="notification">
            <i class="bi bi-envelope"></i>
        </div>
        <div class="col d-flex justify-content-center align-items-center">
            <div class="activityDot"></div>
        </div>
    </div>
    <MessageComp :data=data v-if="messageWindowStatus" />
</template>

<script>
import { ws } from "../../common-js/messages.js"

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
            messageWindowStatus: false
        }
    },

    mounted: function () {
        let bubble = document.getElementById(this.data.Username)
        if (bubble != null) {
            this.img = "url(http://localhost:8080" + this.data.Avatar_image + ")"
            bubble.style.backgroundImage = this.img
        }
        ws.addEventListener('message', (event) => { this.handleUsers(event) });
    },

    methods: {
        handleUsers(event) {
            let incomingData = JSON.parse(event.data.toString());
            if (incomingData.Type == "privateMSG" && !document.querySelector(".messagewindowdiv") && incomingData.Content.Sender == this.data.Username) {
                this.notification = true
            }

            if (incomingData.Type == "privateMSG" && document.querySelector(".messagewindowdiv") && incomingData.Content.Receiver == this.data.Username) {
                this.notification = false
            }

        },
        toggleMessenger() {
            if (document.contains(document.querySelector(".messagewindowdiv")) && this.messageWindowStatus == false) {
                return
            }
            this.messageWindowStatus = !this.messageWindowStatus

        }
    },

}
</script>

<style>
.notification {
    position: absolute;
    right: 10px;
    top: 7px;
}

.messageDetails {
    padding: 2px;
    margin-left: 10px;
    max-width: px;
    display: flex;
    flex-wrap: wrap;
}

.chatwindowPerson {
    transition: 0.3s;
    cursor: pointer;
    padding: 4px;
    border-bottom: solid 1px rgb(219, 219, 219);
    flex-wrap: wrap;
    max-width: 220px;
    position: relative;
}

.chatwindowPerson:hover {
    transition: 0.3s;
    transform: scale(0.99);
    background-color: #D3D3D3;
    padding: 4px;
    border-bottom: solid 1px rgb(219, 219, 219);
}

.chatwindowImg {
    margin-top: 3px;
    margin-left: 10px;
    width: 30px;
    height: 30px;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    -o-user-select: none;
    user-select: none;
}

.chatwindowUser {
    font-size: 12px;
    color: #2E343D;
}

.chatwindowTime {
    color: rgb(165, 159, 159);
    font-size: 10px;
    margin-top: -3px;
    padding-left: 2px;
}
</style>