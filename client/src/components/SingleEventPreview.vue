<template>
    <div class="eventcard">
        <div class="d-flex flex-column flex-wrap justify-content-space">
            <div class="eventcardH d-flex flex-row">
                <div>
                    <img v-if="data.Event.User.Avatar_image == ``" class="postProfileImg col"
                        src="../assets/images/profile.svg" />
                    <div v-else class="bubble3 col" style="margin-top: 7px; margin-left: 2px"></div>
                </div>
                <div class="col posthDetails">
                    <RouterLink :to="{ name: 'profile', params: { id: data.Event.User.Username } }"
                        style="text-decoration: none; width: 100%;" class="href">
                        <div class="col postUser followerLink">{{ data.Event.User.Username }}</div>
                    </RouterLink>
                    <div class="col postTime">{{ humanReadableTime }}</div>
                </div>

            </div>

            <div class="postcardb d-flex flex-row ">
                <div class="col">
                    <div class="postDescp text-wrap text-break text-start col" style="margin-bottom: 5px">
                        {{ getDateTime }}
                    </div>
                    <div class="postTitle text-wrap text-break text-start col">
                        {{ data.Event.Title }}
                    </div>

                    <div class="postDescp text-wrap text-break text-start col">
                        {{ data.Event.Content }}
                    </div>
                </div>

            </div>

            <div class="eventcardF d-flex flex-row">
                <div class="col-6 d-flex justify-content-around">
                    <div class="align-self-center">
                        <a style="font-size: 18px; margin-right: 5px;">{{ countGoing }}</a>
                        <button class="showeventGoingButton" @click="eventGoing" :class="{goingActive: IsGoingActive}">Going</button>
                    </div>
                    <div class="align-self-center ">
                        <a style="font-size: 18px; margin-right: 5px;">{{ countNotGoing }}</a>
                        <button class="showeventNotGoingButton" @click="eventNotGoing" :class="{notGoingActive: IsNotGoingActive}">Not going</button>
                    </div>
                </div>
            </div>


        </div>
    </div>
</template>

<script>
import { timeago } from "../common-js/time.js"
import { delay } from "../common-js/time.js";
import axios from "axios";
export default {
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
        data: {
            type: Boolean,
            required: true
        },
    },

    methods: {
        async eventGoing() {
            console.log('here')
            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let data = {
                Event_id: this.data.Event.Id
            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/group/event/going", data, config)
                    .then((res) => {
                        this.$emit('repliesChanged', true)

                    })
                    .catch((error) => { });
            })

        },

        async eventNotGoing() {
            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let data = {
                Event_id: this.data.Event.Id
            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/group/event/notgoing", data, config)
                    .then((res) => {
                        this.$emit('repliesChanged', true)

                    })
                    .catch((error) => { });
            })

        },

    },

    mounted() {
        let bubble = this.$el.querySelector(".bubble3")
        if (bubble == null) {
            return
        }
        bubble.style.backgroundImage = `url('http://localhost:8080${this.data.Event.User.Avatar_image}')`
    },

    computed: {
        humanReadableTime() {
            return timeago(new Date(Date.now() - new Date(Date.parse(this.data.Event.Created_at ?? "01 Jan 1970 00:00:00 GMT"))))
        },

        getDateTime() {
            let options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
            let date = new Date(this.data.Event.Happening_at);
            return date.toLocaleDateString("en-US", options)
        },

        countGoing() {
            if (this.data.Replies == null) {
                return 0
            }
            if (this.data.Replies == 0) {
                return 0
            }
            const filteredLike = this.data.Replies.filter(function (key) {
                return key.Status == "Going"
            })
            return filteredLike.length
        },

        countNotGoing() {
            if (this.data.Replies == null) {
                return 0
            }
            if (this.data.Replies.length == 0) {
                return 0
            }
            const filteredDislike = this.data.Replies.filter(function (key) {
                return key.Status == "Not Going"
            })
            return filteredDislike.length
        },

        IsGoingActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Replies == null) {
                return false
            }


            const filteredLike = this.data.Replies.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Status == "Going"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },

        IsNotGoingActive() {
                 let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Replies == null) {
                return false
            }

            const filteredLike = this.data.Replies.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Status == "Not Going"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },

    }



}
</script>



<style>
@import "../assets/css/base.css";
</style>