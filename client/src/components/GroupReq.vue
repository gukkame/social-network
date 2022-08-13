<template>
    <div class="oneFollowing row">
        <div class="col-7 d-flex align-items-center">
            <img v-if="data.Avatar_image == ``" class="oneFollowingPic col" src="../assets/images/profile.svg" />
            <div v-else class="bubble2 col-2"></div>
            <div class="col-10 d-flex flex-wrap" style="padding-left: 5px">
                <RouterLink :to="{ name: 'profile', params: { id: data.Username } }"
                    style="text-decoration: none; width: 100%;" class="href">
                    <div class="followerLink d-flex" style="font-size: 16px; color:#2E343D ;">{{ data.Username }}</div>
                </RouterLink>
            </div>

            <div class="col d-flex">
                <button class="followProfile col" @click="acceptRequest(data.Id)"
                    style="margin-right: 6px">Confirm</button>

                <div class="d-flex justify-content-center align-items-center denyFollow"><svg
                        @click="removeRequest(data.Id)" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                        fill="red" class="bi bi-x-lg" viewBox="0 0 16 16">
                        <path
                            d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z" />
                    </svg></div>
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
        data: {
            type: Object,
            required: true
        },
    },

    mounted() {
        let bubble = this.$el.querySelector(".bubble2")
        if (bubble == null) {
            return
        }
        bubble.style.backgroundImage = `url('http://localhost:8080${this.data.Avatar_image}')`
    },


    methods: {
        removeRequest(user) {
            let token = document.cookie
            if (token.length == 0) {
                return router.go("/login")
            }
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let path = this.$route.path
            let person = path.split("/")

            let data = {
                User_id: user,
                Group_id: parseInt(person[2])
            }
            axios.post("http://localhost:8080/group/join/deny", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.$emit('dataChanged', true)
                })
                .catch((error) => { });
        },

        acceptRequest(user) {
            let token = document.cookie
            if (token.length == 0) {
                return router.go("/login")
            }
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let path = this.$route.path
            let person = path.split("/")

            let data = {
                User_Id: user,
                Group_Id: parseInt(person[2])
            }

            axios.post("http://localhost:8080/group/join/accept", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.$emit('dataChanged', true)
                })
                .catch((error) => { });
        }

    },
}

</script>