<script setup>
import ActiviySelectorComp from "../components/ActiviySelector.vue"
import ActivityCommentComp from "../components/ActivityComment.vue"
import SingePostPreviewComp from "../components/SinglePostPreview.vue"
</script>

<template>
    <div class="profile container col">
        <ActiviySelectorComp :selected="selected" />
        <div class="noAvailable d-flex justify-content-center"
            v-if="posts == null && Comments == null && VotedPosts == null">
            No recent activity
        </div>
        <div v-if="posts != null">
            <div class="d-flex">
                <div class="catetitle col">
                    <i class="bi bi-bookmark"></i>
                    Your posts
                    ({{ checkPosts(posts) }})
                </div>
            </div>
            <div class="d-flex">
                <div class="mainheader content col">
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in posts">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="VotedPosts != null">
            <div class="d-flex">
                <div class="catetitle col">
                    <i class="bi bi-bookmark-star"></i>
                    Your reacted posts
                    ({{ checkPosts(VotedPosts) }})
                </div>
            </div>
            <div class="d-flex">
                <div class="mainheader content col">
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in VotedPosts">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="Comments != null">
            <div class="d-flex">
                <div class="catetitle col">
                    <i class="bi bi-chat-dots"></i>
                    Your comments
                    ({{ checkPosts(Comments) }})
                </div>
            </div>
            <div class="d-flex">
                <div class="mainheader content col">
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in Comments">
                            <ActivityCommentComp :data="item" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
import { Form, Field } from "vee-validate";
import axios from "axios";
import router from "../router";
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import { createCookie } from "../common-js/cookies.js";
import { delay } from "../common-js/time.js";

export default {
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
    },

    data() {
        return {
            posts: [],
            VotedPosts: [],
            Comments: [],
            selected: "activity"
        };
    },
    components: {
        Form,
        Field,
    },

    created() {
        this.$watch(
            () => this.$route.path,
            () => {
                this.fetchData()
            },
            { immediate: true }
        )
    },

    methods: {
        async fetchData() {
            let token = document.cookie
            let correctToken = token.split(":")

            if (token.length === 0) {
                return router.push("/")
            }

            let data = {

            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            await delay(200).then(() => {
                axios.post("http://localhost:8080/activity", data, config)
                    .then((res) => {
                        if (res.data.message === "User not authenticated") {
                            return router.push("/")
                        }
                        this.posts = res.data.Package.Posts
                        this.VotedPosts = res.data.Package.VotedPosts
                        this.Comments = res.data.Package.Comments
                        let Cookie = res.data.Cookie
                        if (Cookie.Id.length != 0 && Cookie.Username.length != 0) {
                            createCookie(Cookie.Id, Cookie.Username)
                        }
                    })
                    .catch((error) => { });
            })
        },
        checkPosts(arg) {
            if (arg == null) {
                return 0
            } else {
                return arg.length
            }
        }
    },



};
</script>

