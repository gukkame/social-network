<template>
    <div class="profile container col justify-content-center d-flex">
        <div class="d-flex justify-content-center align-items-center">
            <div class="notFollowingHeader d-flex justify-content-center row">
                <img class="profileimage col" src="../assets/images/lock.svg" />
                <div class="d-flex justify-content-center profilename">
                    This Account is Private
                </div>
                <div class="d-flex justify-content-center profilename" style="font-size: 16px;">Follow {{ returnName }}
                    to see their profile.
                </div>
                <div class="d-flex">
                    <div class="d-flex justify-content-center col">
                        <button v-if="notRequested" @click="performFollow" class="followProfile">Follow</button>
                        <button v-if="Requested" @click="performUnrequest" class="unfollowProfile">Requested
                        </button>
                    </div>
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
            notRequested: false,
            Requested: false,
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.path,
            () => {
                this.fetchFollowStatus()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchFollowStatus() {
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

            await delay(20).then(() => {
                axios.post("http://localhost:8080/checkfollow", data, config)
                    .then((res) => {

                        if (res.data.message == "Profile does not exist") {
                            return
                        }

                        if (res.data.message == "Not requested") {
                            this.Requested = false
                            this.notRequested = true
                            return
                        }


                        if (res.data.message == "Requested") {
                            this.notRequested = false
                            this.Requested = true
                            return
                        }

                    })
                    .catch((error) => { });
            })

        },

        performFollow() {
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
                    this.fetchFollowStatus()
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
            let person = path.split("/")

            let data = {
                Username: person[2]
            }
            axios.post("http://localhost:8080/cancelrequest", data, config)
                .then((res) => {

                    if (res.data.message == "User not authenticated") {
                        return router.push(`"${path}"`)
                    }

                    if (res.data.message == "Profile does not exist") {
                        return router.push(`"${path}"`)
                    }
                    this.fetchFollowStatus()
                })
                .catch((error) => { });
        }

    },

    computed: {
        returnName() {
            let path = this.$route.path
            let person = path.split("/")
            return person[2]
        }
    },
}
</script>