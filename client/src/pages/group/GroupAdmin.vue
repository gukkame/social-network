<script setup>
import OneGroupReqComp from "../../components/GroupReq.vue"
</script>

<template>
    <div class="d-flex row flex-wrap justify-content-around">
        <div class="followerSect col">
            <div class="col followSectTitle d-flex justify-content-center align-items-center">
                <div class="followSectTitle2">Join requests ({{ returnLength(joinRequests) }})
                </div>
            </div>
            <div class="followSectBody overflow-auto flex-wrap justify-content-space">
                <div v-if="joinRequests != null" v-for="oneFollowReq in joinRequests" :key="oneFollowReq.Username">
                    <OneGroupReqComp :data="oneFollowReq" @dataChanged="fetchRequests" />
                </div>
                <div v-else class="d-flex justify-content-center align-items-center"
                    style="width: 100%; height: 91%; font-size: 22px">No requests available</div>
            </div>
        </div>
    </div>
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import axios from "axios"
import { delay } from "../../common-js/time.js";
import router from "../../router";

export default {
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
    },
    data() {
        return {
            joinRequests: []
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.params.id,
            () => {
                this.fetchRequests()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchRequests() {
            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let path = this.$route.path
            let person = path.split("/")

            let data = {
                GroupID: parseInt(person[2])
            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/group/join/requests", data, config)
                    .then((res) => {
                        this.joinRequests = res.data
                    })
                    .catch((error) => { });
            })

        },

        returnLength(arg) {
            if (arg == null) {
                return 0
            }
            return arg.length
        },
    },

    computed: {
        returnTime() {
            const d = new Date(this.info.Date);
            let text = d.toLocaleDateString();
            return text
        }
    }
}
</script>

<style>
@import "@/assets/css/base.css";
</style>