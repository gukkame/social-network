<script setup>
import ActiviySelectorComp from "../components/ActiviySelector.vue"
import PageNotFoundComp from "./PageNotFound.vue"
import NotFollowingComp from "../components/NotFollowing.vue"
import OneFollowerComp from "../components/OneFollower.vue"
import OneFollowingComp from "../components/OneFollowing.vue"
import OneFollowReqComp from "../components/FollowerReq.vue"
</script>

<template>
    <PageNotFoundComp v-if="notExist" />
    <NotFollowingComp v-else-if="notFollowing" />
    <div v-else class="followers container col">
        <ActiviySelectorComp :selected="selected" />
        <div class="d-flex row flex-wrap justify-content-around">
            <div class="followerSect col">
                <div class="col followSectTitle d-flex justify-content-center align-items-center">
                    <div class="followSectTitle2">Following ({{ returnLength(following) }})</div>
                </div>
                <div class="followSectBody overflow-auto flex-wrap justify-content-space">
                    <div v-if="following != null" v-for="onefollowing in following" :key="onefollowing.Username">
                        <OneFollowingComp :data="onefollowing" :status="status" @dataChanged="fetchFollowers"/>
                    </div>
                    <div v-else class="d-flex justify-content-center align-items-center"
                        style="width: 100%; height: 91%; font-size: 22px">Not following anyone</div>
                </div>
            </div>
            <div class="followerSect col">
                <div class="col followSectTitle d-flex justify-content-center align-items-center">
                    <div class="followSectTitle2">Followers ({{ returnLength(followers) }})</div>
                </div>

                <div class="followSectBody overflow-auto flex-wrap justify-content-space">
                    <div v-if="followers != null" v-for="onefollower in followers" :key="onefollower.Username">
                        <OneFollowerComp :data="onefollower" :status="status" @dataChanged="fetchFollowers"/>
                    </div>
                    <div v-else class="d-flex justify-content-center align-items-center"
                        style="width: 100%; height: 91%; font-size: 22px">No followers</div>
                </div>
            </div>
            <div v-if="status == `owner`" class="followerSect col">
                <div class="col followSectTitle d-flex justify-content-center align-items-center">
                    <div class="followSectTitle2">Follower requests ({{ returnLength(followReq) }})</div>
                </div>
                <div class="followSectBody overflow-auto flex-wrap justify-content-space">
                    <div v-if="followReq != null" v-for="onefollowReq in followReq" :key="onefollowReq.Username">
                        <OneFollowReqComp :data="onefollowReq" @dataChanged="fetchFollowers"/>
                    </div>
                    <div v-else class="d-flex justify-content-center align-items-center"
                        style="width: 100%; height: 91%; font-size: 22px">No follower requests</div>
                </div>
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
    },
    data() {
        return {
            following: [],
            followers: [],
            followReq: [],
            status: "",
            notFollowing: false,
            notExist: false,
            selected: "followers"
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.path,
            () => {
                this.fetchFollowers()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchFollowers() {
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
                Username: person[2]
            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/followers", data, config)
                    .then((res) => {

                        if (res.data.message == "Profile does not exist") {
                            this.notExist = true
                            return
                        }

                        if (res.data.message == "Not following") {
                            this.notFollowing = true
                            return
                        }

                        this.followers = res.data.Followers
                        this.following = res.data.Following
                        this.followReq = res.data.FollowReq
                        this.status = res.data.Status
                    })
                    .catch((error) => { });
            })

        },

        returnLength(arg) {
            if (arg == null) {
                return 0
            }
            return arg.length
        }

    },
}
</script>

<style>
@import "@/assets/css/base.css";
</style>
