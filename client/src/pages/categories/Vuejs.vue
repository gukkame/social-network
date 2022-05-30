<script setup>
import CategoryComp from "../../components/Categories.vue"
import SingePostPreviewComp from "../../components/SinglePostPreview.vue"
import CreatePostComp from "../../components/CreatePost.vue"
</script>

<template>
    <CategoryComp :category="category" />
    <div class="mainpagec container col">
        <div class="d-flex">
            <div class="catetitle col">
                <i class="bi bi-bookmark"></i>
                Vue.js ({{ checkPosts }})
            </div>
            <div>
                <button v-if="LoggedIn" class="createPost col" data-bs-toggle="modal"
                    data-bs-target="#staticBackdrop">Create a post</button>
            </div>
        </div>
        <div class="d-flex">
            <div class="mainheader content col">
                <div class="d-flex flex-wrap justify-content-around contentinside col">
                    <div class="noAvailable d-flex justify-content-center" v-if="checkPosts == 0">
                        No posts available
                    </div>
                    <div v-for="item in posts">
                        <SingePostPreviewComp :data="item" />
                    </div>
                </div>
            </div>
        </div>
    </div>
    <CreatePostComp />
</template>

<script >
import axios from "axios";
import { createCookie } from "../../common-js/cookies.js";
import { delay } from "../../common-js/time.js";
export default {
    name: 'Header',
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
    },

    data() {
        return {
            posts: [],
            category: "vue"
        }
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
            if (this.$route.path != "/Vue.js") {
                return
            }
            let currentRouter = this.$route.path
            let correctCategory = currentRouter.split("/")
            let data = {
                categoryname: correctCategory[1]
            }

            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            await delay(200).then(() => {
                axios.post("http://localhost:8080/onecategory", data, config)
                    .then((res) => {
                        if (res.data.message === "Post request failed") {
                            return router.go(-1)
                        }
                        this.posts = res.data.Posts
                        let Cookie = res.data.Cookie
                        if (Cookie.Id.length != 0 && Cookie.Username.length != 0) {
                            createCookie(Cookie.Id, Cookie.Username)
                        }
                    })
                    .catch((error) => { });
            })

        },

    },

    computed: {
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