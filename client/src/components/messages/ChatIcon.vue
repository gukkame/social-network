<script setup>
import ChatBoxComp from "./ChatBox.vue"
import { connectToWS, ws } from "../../common-js/messages.js"
</script>

<template>
    <ChatBoxComp v-if="chatWindowStatus" />
    <div class="chatbuttondiv">
        <div class>
            <i :class="{ chatBoxActive: chatWindowStatus }" @click="toggleChatbox"
                class="chatbutton bi bi-chat-dots"></i>
        </div>
    </div>
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import router from "../../router/index.js"
import axios from "axios"
import { createCookie } from "../../common-js/cookies.js";
import { delay } from "../../common-js/time.js";

export default {
    data() {
        return {
            chatWindowStatus: false
        }
    },
    
    methods: {
        async toggleChatbox() {
            let token = document.cookie
            if (token.length === 0) {
                return router.go("/")
            }
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let data = {

            }

            await delay(50).then(() => {
                axios.post("http://localhost:8080/authedchat", data, config)
                    .then((res) => {
                        if (res.data.message === "User not authenticated") {
                            document.cookie = 'Token=; Max-Age=-99999999;';
                            return router.go("/")
                        }
                        this.chatWindowStatus = !this.chatWindowStatus
                        this.info = res.data
                        createCookie(res.data.Username, res.data.Token)
                    })
                    .then(() => {
                        if (this.chatWindowStatus == false) {
                            let msg = {
                                "Type": "closeClient"
                            }
                            ws.send(JSON.stringify(msg));
                        }
                    })
                    .catch((error) => { });
            })

        },

    }
}






</script>

<style>
.chatBoxActive {
    color: #FF9D5A;
}

.chatbutton {
    background-color: rgb(253, 252, 252);
    transition: 0.25s;
    font-size: 45px;
    background: rgba(0, 0, 0, 0);
}

.chatbutton:hover {
    transform: scale(1.13);
    color: #FF9D5A;
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