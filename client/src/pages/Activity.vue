<script setup>
import ActiviySelectorComp from "../components/ActiviySelector.vue"
import ActivityCommentComp from "../components/ActivityComment.vue"
import SingePostPreviewComp from "../components/SinglePostPreview.vue"
import PageNotFoundComp from "./PageNotFound.vue"
import NotFollowingComp from "../components/NotFollowing.vue"
</script>

<template>
    <PageNotFoundComp v-if="notExist" />
    <NotFollowingComp v-else-if="notFollowing" />
    <div v-else class="profile container col">
        <ActiviySelectorComp :selected="selected" />
        <div class="noAvailable d-flex justify-content-center"
            v-if="posts == null && Comments == null && VotedPosts == null">
            No recent activity
        </div>
        <div v-if="posts != null">
            <div class="d-flex">
                <div v-if="owner" class="catetitle col">
                    <i class="bi bi-bookmark"></i>
                    Your posts
                    ({{ checkPosts(posts) }})
                </div>
                <div v-else class="catetitle col">
                    <i class="bi bi-bookmark"></i>
                    {{ getOwner }} posts
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
                <div v-if="owner" class="catetitle col">
                    <i class="bi bi-bookmark-star"></i>
                    Your reacted posts
                    ({{ checkPosts(VotedPosts) }})
                </div>
                <div v-else class="catetitle col">
                    <i class="bi bi-bookmark-star"></i>
                    {{ getOwner }} reacted posts
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
                <div v-if="owner" class="catetitle col">
                    <i class="bi bi-chat-dots"></i>
                    Your comments
                    ({{ checkPosts(Comments) }})
                </div>
                <div v-else class="catetitle col">
                    <i class="bi bi-chat-dots"></i>
                    {{ getOwner }} comments
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
            selected: "activity",
            notFollowing: false,
            notExist: false,
            owner: false,
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

            await delay(200).then(() => {
                axios.post("http://localhost:8080/activity", data, config)
                    .then((res) => {

                        if (res.data.message == "Profile does not exist") {
                            this.notExist = true
                            return
                        }

                        if (res.data.message == "Not following") {
                            this.notFollowing = true
                            return
                        }

                        if (res.data.Status == "owner") {
                            this.owner = true
                        }

                        this.posts = res.data.Package.Posts
                        this.VotedPosts = res.data.Package.VotedPosts
                        this.Comments = res.data.Package.Comments
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
        },
    },
    computed: {
        getOwner() {
            let path = this.$route.path.split("/")
            return path[2]
        }
    }
};
</script>

