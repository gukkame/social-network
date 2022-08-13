<template>
    <div class="oneFollowing row">
        <div class="col-7 d-flex align-items-center">
            <img v-if="data.Avatar_image == ``" class="oneFollowingPic col" src="../assets/images/profile.svg" />
            <div  v-else class="bubble2 col-2"></div>
            <div class="col-10 d-flex flex-wrap" style="padding-left: 5px">
                <RouterLink :to="{ name: 'profile', params: { id: data.Username } }"
                    style="text-decoration: none; width: 100%;" class="href">
                    <div class="followerLink d-flex" style="font-size: 16px; color:#2E343D ;">{{ data.Username }}</div>
                </RouterLink>
            </div>

            <div v-if="status == `owner`" class="col d-flex">
                <button class="unfollowProfile" @click="performFollow(data.Username)">Unfollow</button>
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
import { nextTick } from "vue";
export default {
    props: {
        data: {
            type: Object,
            required: true
        },
        status: {
            type: String,
            required: true
        }
    },

    mounted() {
            let bubble = this.$el.querySelector(".bubble2")
            if (bubble == null) {
                return
            }
            bubble.style.backgroundImage = `url('http://localhost:8080${this.data.Avatar_image}')` 
    },

    methods: {
        performFollow(user) {
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

            let data = {
                Username: user
            }
            axios.post("http://localhost:8080/follow", data, config)
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