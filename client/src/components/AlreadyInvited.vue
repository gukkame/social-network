<template>
    <div class="d-flex" style="margin-top: 90px">
        <div class="content col">
            <div class="d-flex flex-wrap justify-content-around contentinside col">
                <div class="noAvailable d-flex justify-content-center">
                    You have been invited to this group!
                </div>
            </div>
            <div class="noAvailable d-flex justify-content-center">
                <button @click="deny" class="unRequestGroup">Deny</button>
                <button @click="accept" style="width: 135px" class="requestGroup">Accept</button>
            </div>
        </div>
    </div>
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import axios from "axios"
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

    methods: {
        deny() {
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
                Group_id: parseInt(groupId[2])
            }

            axios.post("http://localhost:8080/group/invite/deny", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.$emit('inviteChanged', true)
                })
                .catch((error) => { });
        },

        accept() {
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
                Group_id: parseInt(groupId[2])
            }

            axios.post("http://localhost:8080/group/invite/accept", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                     this.$emit('inviteChanged', true)

                })
                .catch((error) => { });
        }

    },
}
</script>