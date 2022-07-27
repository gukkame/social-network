<script setup>
import SingleEventPreviewComp from "../../components/SingleEventPreview.vue"
import CreateGroupEventModal from "../../components/modals/CreateGroupEvent.vue"
</script>

<template>

    <div class="d-flex justify-content-end" style="margin-bottom: 15px">
        <div>
            <button class="createPost col" style="width: 160px" data-bs-toggle="modal" data-bs-target="#groupEvent">Create an event</button>
        </div>
    </div>
    <div class="d-flex">
        <div class="mainheader content col">
            <div class="d-flex flex-wrap justify-content-around contentinside col">
                <div class="noAvailable d-flex justify-content-center" v-if="checkPosts == 0">
                    No events available
                </div>
                <div v-for="event in events">
                    <SingleEventPreviewComp :data="event" @repliesChanged="fetchGroupEvents" />
                </div>
            </div>
        </div>
        <CreateGroupEventModal/>
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
            events: []
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.params.id,
            () => {
                this.fetchGroupEvents()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchGroupEvents() {
            let token = document.cookie
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

            await delay(100).then(() => {
                axios.post("http://localhost:8080/group/events", data, config)
                    .then((res) => {
                         this.events = res.data
                        
                    })
                    .catch((error) => { });
            })

        },

    },

    computed: {
        returnTime() {
            const d = new Date(this.info.Date);
            let text = d.toLocaleDateString();
            return text
        },

        checkPosts() {
            if (this.events == null) {
                return 0
            } else {
                return this.events.length
            }
        }
    }
}
</script>

<style>
@import "@/assets/css/base.css";
</style>