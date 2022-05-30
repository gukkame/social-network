<script setup>
import ActiviySelectorComp from "../components/ActiviySelector.vue"
import { create } from "yup/lib/Reference"

</script>

<template>
    <div class="profile container col">
        <ActiviySelectorComp :selected="selected" />
        <div class="d-flex justify-content-center">
            <div class="profileheader d-flex justify-content-center row">
                <img class="profileimage col" src="../assets/images/profile.svg" />
                <div class="d-flex justify-content-center profilename">{{ info.FirstName }} {{ info.LastName }}</div>
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
                <div class="sect row align-items-center">
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
            </div>
        </div>
    </div>
</template>


<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import router from "../router/index.js"
import axios from "axios"
import { createCookie } from "../common-js/cookies.js";
import { delay } from "../common-js/time.js";

export default {
    data() {
        return {
            info: {
                Username: "",
                Email: "",
                Age: 0,
                Gender: "",
                FirstName: "",
                LastName: "",
            },
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
            if (this.$route.path != "/profile") {
                return
            }

            let token = document.cookie
            if (token.length === 0) {
                return router.push("/")
            }
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let data = {

            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/profile", data, config)
                    .then((res) => {
                        if (res.data.message === "User not authenticated") {
                            console.log("wo")
                            document.cookie = 'Token=; Max-Age=-99999999;';
                            return router.push("/")
                        }
                        this.info = res.data
                        createCookie(res.data.Token, res.data.Username)
                    })
                    .catch((error) => { });
            })

        }

    },
}
</script>

<style>
@import "@/assets/css/base.css";
</style>