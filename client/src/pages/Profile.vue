<script setup>
import ActiviySelectorComp from "../components/ActiviySelector.vue"
import PageNotFoundComp from "./PageNotFound.vue"
import NotFollowingComp from "../components/NotFollowing.vue"
import ChangeProfileComp from "../components/ChangeProfile.vue"

</script>

<template>
    <PageNotFoundComp v-if="notExist" />
    <NotFollowingComp v-else-if="notFollowing" />
    <div v-else class="profile container col">
        <ActiviySelectorComp :selected="selected" />
        <div class="d-flex justify-content-center">
            <div class="profileheader d-flex justify-content-center row">
                <div v-if="info.Status == `owner`" class="d-flex">
                    <div class="d-flex justify-content-end col">
                        <div style="padding-top:4.1px; margin-right: 10px; color: #5c6166;">Visibility: {{
                                info.Profile_status
                        }}</div>
                        <i type="button" class="editPost bi bi-gear align-self-start" data-bs-toggle="modal"
                            data-bs-target="#changeProfile"></i>
                    </div>
                </div>
                <img class="profileimage col" src="../assets/images/profile.svg" />
                <div class="d-flex justify-content-center profilename">{{ info.FirstName }} {{ info.LastName }}
                    ({{ info.NickName }})</div>
                <div class="d-flex">
                    <div class="d-flex justify-content-center col">
                        <button v-if="info.Status == `notfollowing` && LoggedIn" @click="performFollow"
                            class="followProfile">Follow</button>
                        <button v-if="info.Status == `following` && LoggedIn" @click="performFollow"
                            class="unfollowProfile">Unfollow</button>
                    </div>
                </div>
            </div>
        </div>
        <div class="d-flex justify-content-center">

            <div class="profileAbout d-flex justify-content-center align-items-start row">
                <div class="sect row align-items-center">
                    <div class="col-1"></div>
                    <div class="col d-flex pro-1" style="margin-bottom: 10px">About me</div>
                    <div class="lined"></div>
                    <div class="d-flex justify-content-start pro-2" style="padding:30px; margin-top: -14px;">{{
                            info.AboutMe
                    }}</div>
                </div>
            </div>
        </div>
        <div class="d-flex justify-content-center">

            <div class="profileinfo d-flex justify-content-center row">
                <div class="sect row align-items-center">
                    <div class="col-1"></div>
                    <div class="col d-flex pro-1">Username</div>
                    <div class="col d-flex justify-content-start pro-2">{{ info.Username }}</div>
                    <div class="lined"></div>
                </div>
                <div v-if="info.Status == `owner`" class="sect row align-items-center">
                    <div class="col-1"></div>
                    <div class="col d-flex pro-1">Email</div>
                    <div class="col d-flex justify-content-start pro-2">{{ info.Email }}</div>
                    <div class="lined"></div>
                </div>
                <div class="sect row align-items-center">
                    <div class="col-1"></div>
                    <div class="col d-flex pro-1">Age</div>
                    <div class="col d-flex justify-content-start pro-2">{{ info.Age }}</div>
                    <div class="lined"></div>
                </div>
                <div class="sect row align-items-center">
                    <div class="col-1"></div>
                    <div class="col d-flex pro-1">Gender</div>
                    <div class="col d-flex justify-content-start pro-2">{{ info.Gender }}</div>
                    <div class="lined"></div>
                </div>
                <div class="sect row align-items-center">
                    <div class="col-1"></div>
                    <div class="col d-flex pro-1">Account created</div>
                    <div class="col d-flex justify-content-start pro-2">{{ returnTime }}</div>
                    <div class="lined"></div>
                </div>
            </div>
        </div>
    </div>
    <ChangeProfileComp v-if="info.Status == `owner`" :data="info.Profile_status" />
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
            info: {
                Username: "",
                Email: "",
                Age: 0,
                Gender: "",
                FirstName: "",
                LastName: "",
                NickName: "",
                AboutMe: "",
                Avatar_image: "",
                Date: "",
                Status: "",
                Profile_status: "",
            },
            notFollowing: false,
            notExist: false,
            selected: "profile"
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.path,
            () => {
                this.fetchProfile()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchProfile() {
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
                axios.post("http://localhost:8080/profile", data, config)
                    .then((res) => {

                        if (res.data.message == "Profile does not exist") {
                            this.notExist = true
                            return
                        }
                        if (res.data.Status == "owner") {
                            this.info = res.data
                            return
                        }


                        if (res.data.message == "Not following") {
                            this.notFollowing = true
                            return
                        }

                        this.info = res.data

                    })
                    .catch((error) => { });
            })

        },

        performFollow() {
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
                Username: person[2]
            }
            axios.post("http://localhost:8080/follow", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.fetchProfile()
                })
                .catch((error) => { });
        }

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