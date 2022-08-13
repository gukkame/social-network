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
            <div class="col chatwindowTime">Inactive</div>
        </div>
        <div class="col d-flex justify-content-center align-items-center">
            <div class="offlineDot"></div>
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
.offlineDot {
    border-radius: 50%;
    width: 15px;
    height: 15px;
    background-color: rgb(253, 252, 252);
    border: solid 1px rgb(219, 219, 219);
}
</style>