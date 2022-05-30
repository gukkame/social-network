<script setup>
import CategoryComp from "../components/Categories.vue"
import SingePostPreviewComp from "../components/SinglePostPreview.vue"
</script>

<template>
    <div class="d-flex justify-content-center">
        <div class="homeWelcome" v-if="LoggedIn">Welcome, <a style="color: #FF9D5A">{{ data }}</a>!</div>
        <div class="homeWelcome" v-else>Welcome!</div>
    </div>
    <CategoryComp :category="category" />
    <div class="mainpagec container col">
        <div class="catetitle">
            <i class="catetitleicon bi bi-calendar"></i>Latest posts
        </div>
        <div class="d-flex">
            <div class="mainheader content col">
                <div class="noAvailable d-flex justify-content-center"
                    v-if="checkPosts(go) == 0 && checkPosts(html) == 0 && checkPosts(css) == 0 && checkPosts(javascript) == 0 && checkPosts(vue) == 0">
                    No new posts available
                </div>
                <div v-if="checkPosts(go) != 0">
                    <div class="d-flex">
                        <div class="catetitle col">
                            <i class="bi bi-bookmark"></i>
                            Go
                        </div>
                        <div>
                            <RouterLink to="/Go">
                                <button class="seeall col">See all</button>
                            </RouterLink>
                        </div>
                    </div>
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in go">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>

                <div v-if="checkPosts(html) != 0">
                    <div class="d-flex">
                        <div class="catetitle col">
                            <i class="bi bi-bookmark"></i>
                            HTML5
                        </div>
                        <div>
                            <RouterLink to="/HTML5">
                                <button class="seeall col">See all</button>
                            </RouterLink>
                        </div>
                    </div>
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in html">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>

                <div v-if="checkPosts(css) != 0">
                    <div class="d-flex">
                        <div class="catetitle col">
                            <i class="bi bi-bookmark"></i>
                            CSS
                        </div>
                        <div>
                            <RouterLink to="/CSS">
                                <button class="seeall col">See all</button>
                            </RouterLink>
                        </div>
                    </div>
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in css">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>

                <div v-if="checkPosts(javascript) != 0">
                    <div class="d-flex">
                        <div class="catetitle col">
                            <i class="bi bi-bookmark"></i>
                            JavaScript
                        </div>
                        <div>
                            <RouterLink to="/JavaScript">
                                <button class="seeall col">See all</button>
                            </RouterLink>
                        </div>
                    </div>
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in javascript">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>

                <div v-if="checkPosts(vue)">
                    <div class="d-flex">
                        <div class="catetitle col">
                            <i class="bi bi-bookmark"></i>
                            Vue.js
                        </div>
                        <div>
                            <RouterLink to="/Vue.js">
                                <button class="seeall col">See all</button>
                            </RouterLink>
                        </div>
                    </div>
                    <div class="d-flex flex-wrap justify-content-around contentinside col">
                        <div v-for="item in vue">
                            <SingePostPreviewComp :data="item" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script >
import axios from "axios";
import { createCookie } from "../common-js/cookies.js";
import { delay } from "../common-js/time.js";
export default {
    data() {
        return {
            category: "",
            go: [],
            html: [],
            css: [],
            javascript: [],
            vue: [],
        }

    },

    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
        data: {
            type: String,
            required: true
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
            if (this.$route.path != "/") {
                return
            }

            let data = {

            }

            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }


            await delay(200).then(() => {
                axios.post("http://localhost:8080/allcategory", data, config)
                    .then((res) => {
                        if (res.data.message === "Post request failed") {
                            return router.go(-1)
                        }
                        this.go = res.data.Posts.Go
                        this.html = res.data.Posts.Html
                        this.css = res.data.Posts.Css
                        this.javascript = res.data.Posts.JavaScript
                        this.vue = res.data.Posts.Vuejs
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


}
</script>


<style>
@import "../assets/css/base.css";
</style>