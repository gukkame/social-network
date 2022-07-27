<script setup>
import PageNotFoundComp from "../PageNotFound.vue"
import NotPartOfGroupComp from "../../components/NotPartOfGroup.vue"
import AlreadyInvitedComp from "../../components/AlreadyInvited.vue"
import GroupPageSelectorComp from "../../components/GroupPageSelector.vue"
import OneGroupComp from "../../components/OneGroup.vue"
</script>

<template>
    <PageNotFoundComp v-if="notExist" />
    <div v-else>
        <OneGroupComp v-if="fetched" :data="group" />
        <div class="d-flex flex-wrap justify-content-around col" style="margin-top: 100px">
            <GroupPageSelectorComp :selected="selected" :data="group" />
        </div>
    </div>
    <NotPartOfGroupComp :data="group" v-if="group.User_status == `Not Requested` || group.User_status == `Requested`"
        @requestChanged="fetchGroup" />
    <AlreadyInvitedComp v-else-if="group.User_status == `Invited`"   @inviteChanged="fetchGroup"/>
    <router-view v-else />
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import axios from "axios"
import { delay } from "../../common-js/time.js";

export default {
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
    },
    data() {
        return {
            group: {},
            notExist: false,
            fetched: false,
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.params.id,
            () => {
                this.fetchGroup()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },
    computed: {
        selected() {
            return this.$route.name
        }
    },
    methods: {
        async fetchGroup() {
            this.notExist = false
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
                axios.post("http://localhost:8080/group", data, config)
                    .then((res) => {

                        if (res.data.message == "Group does not exist") {
                            this.notExist = true
                            return
                        }
                        this.group = res.data
                        this.fetched = true

                    })
                    .catch((error) => { });
            })
        },
    },
}
</script>

<style>
@import "@/assets/css/base.css";
</style>