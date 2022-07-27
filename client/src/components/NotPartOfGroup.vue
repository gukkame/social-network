<template>
    <div class="d-flex" style="margin-top: 90px">
        <div class="content col">
            <div class="d-flex flex-wrap justify-content-around contentinside col">
                <div class="noAvailable d-flex justify-content-center">
                    To see the rest of the content, you need to join the group!
                </div>
            </div>
            <div class="noAvailable d-flex justify-content-center">
                <button v-if="data.User_status == `Not Requested`" @click="performRequest" class="requestGroup">Request
                    to join</button>
                <button v-if="data.User_status == `Requested`" @click="performUnrequest"
                    class="unRequestGroup">Requested
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import axios from "axios"
import { delay } from "../common-js/time.js";
import router from "../router";
export default {
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
        data: {

        }
    },
    data() {
        return {
            notRequested: false,
            Requested: false,
        }
    },

    methods: {
        performRequest() {
            let token = document.cookie
            if (token.length == 0) {
                return router.push("/login")
            }
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let path = this.$route.path
            let groupId = path.split("/")

            let data = {
                GroupID: parseInt(groupId[2])
            }
            axios.post("http://localhost:8080/group/join", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.$emit('requestChanged', true)
                })
                .catch((error) => { });
        },

        performUnrequest() {
            let token = document.cookie
            if (token.length == 0) {
                return router.push("/login")
            }
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let path = this.$route.path
            let groupId = path.split("/")

            let data = {
                GroupID: parseInt(groupId[2])
            }
            axios.post("http://localhost:8080/group/join/cancel", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.$emit('requestChanged', true)

                })
                .catch((error) => { });
        }

    },
}
</script>