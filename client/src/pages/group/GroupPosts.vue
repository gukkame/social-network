<script setup>
import CreateGroupPostModal from "../../components/modals/CreateGroupPost.vue"
import SingleGroupPostPreviewComp from "../../components/SingleGroupPostPreview.vue"
</script>

<template>

    <div class="d-flex justify-content-end" style="margin-bottom: 15px">
        <div>
            <button class="createPost col" data-bs-toggle="modal" data-bs-target="#createGroupPost">Create a
                post</button>
        </div>
    </div>
    <div class="d-flex">
        <div class="mainheader content col">
            <div class="d-flex flex-wrap justify-content-around contentinside col">
                <div class="noAvailable d-flex justify-content-center" v-if="checkPosts == 0">
                    No posts available
                </div>
                <div v-for="item in posts">
                    <SingleGroupPostPreviewComp :data="item" />
                </div>
            </div>
        </div>
        <CreateGroupPostModal/>
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
            posts: []
        }
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.params.id,
            () => {
                this.fetchGroupPosts()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchGroupPosts() {
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
                axios.post("http://localhost:8080/group/posts", data, config)
                    .then((res) => {
                         this.posts = res.data
                         console.log(this.posts)
                        
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
            if (this.posts == null) {
                return 0
            } else {
                return this.posts.length
            }
        }
    }
}
</script>

<style>
@import "@/assets/css/base.css";
</style>