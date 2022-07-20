<script setup>
import CommentComp from "../components/Comment.vue"
import PageNotFoundComp from "./PageNotFound.vue"
import EditPostComp from "../components/modals/EditPost.vue"
import DeletePostComp from "../components/modals/DeletePost.vue"
</script>

<template>
    <div v-if="post.Id != 0" class="showpost">
        <div class="d-flex justify-content-center">
            <div class="showpostinfo">
                <div v-if="IsPostOwner" class="d-flex">
                    <div class="d-flex justify-content-start col-11">
                        <i type="button" class="editPost bi bi-gear align-self-start" data-bs-toggle="modal"
                            data-bs-target="#editpost"></i>
                    </div>
                    <i type="button" data-bs-toggle="modal" data-bs-target="#deletepost"
                        class="deletePost bi bi-trash justify-content-end"></i>
                </div>
                <div class="showpostinfoh d-flex flex-row">
                    <div>
                        <img class="showpostProfileImg col" src="../assets/images/profile.svg" />
                    </div>
                    <div class="col posthDetails">
                        <div class="col showpostUser">{{ post.Username }}</div>
                        <div class="col showpostTime">{{ humanReadableTime }}</div>
                    </div>
                    <div class="col-2 p-1">
                        <img v-if="post.CategoryTitle == `Go`" class="showpostCategoryImg col"
                            src="../assets/images/Go.svg" />
                        <img v-if="post.CategoryTitle == `HTML5`" class="showpostCategoryImg col"
                            src="../assets/images/HTML.svg" />
                        <img v-if="post.CategoryTitle == `CSS`" class="showpostCategoryImg col"
                            src="../assets/images/CSS.svg" />
                        <img v-if="post.CategoryTitle == `JavaScript`" class="showpostCategoryImg col"
                            src="../assets/images/Javascript.svg" />
                        <img v-if="post.CategoryTitle == `Vue.js`" class="showpostCategoryImg col"
                            src="../assets/images/Vue.svg" />
                    </div>
                </div>
                <div class="showpostinfob d-flex flex-row">
                    <div class="col">
                        <div class="showpostTitle text-wrap text-break text-start col">
                            {{ post.Title }}
                        </div>
                        <div class="showpostDescp text-wrap text-break text-start col">
                            {{ post.Description }}
                        </div>
                    </div>
                </div>
                <div class="showpostinfof d-flex flex-row">
                    <div class="col d-flex justify-content-center">
                        <i type="button" @click="likePost(post.Id)" :class="{ likeActive: IsLikeActive }"
                            class="showpostLike align-self-center  bi bi-hand-thumbs-up"></i>
                        <span class="showpostReactCount align-self-center" style="padding-top: 12px">{{
                                countLikes
                        }}</span>
                    </div>
                    <div class="col d-flex justify-content-center">
                        <div class="align-self-center">
                            <i class="showpostComment bi bi-chat-dots"></i>
                            <span class="showpostReactCount" style="padding-left: 2px;">{{ humanReadableAmount }}</span>
                        </div>
                    </div>
                    <div class="col d-flex justify-content-center">
                        <div class="align-self-center ">
                            <i type="button" @click="dislikePost(post.Id)" :class="{ dislikeActive: IsDislikeActive }"
                                class="showpostDislike bi bi-hand-thumbs-down"></i>
                            <span class="showpostReactCount ">{{ countDislikes }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="d-flex justify-content-center">
            <div v-if="LoggedIn == true || Comments != null" class="commentsection d-flex justify-content-center row">
                <div v-if="LoggedIn" class="addComment">
                    <Form @submit="newComment" v-slot="{ errors }" :validation-schema="errorSchema">
                        <div class="user-details row">
                            <div class="commentTitle">Comment</div>
                            <div style="padding: 0px 40px 0px 40px">
                                <div class="row">
                                    <div style="max-height: 115px">
                                        <div class="col">
                                            <Field class="comment-description-input" id="newComment"
                                                :class="{ commentinputerror: errors.description }" as="textarea"
                                                name="description" placeholder="Add a comment" />
                                            <br />
                                        </div>
                                        <div class="d-flex justify-content-end">
                                            <span class="formErrors">{{ errors.description }}</span>
                                            <div class="formErrors">{{ errormsg }}</div>
                                        </div>
                                    </div>
                                    <div class="d-flex justify-content-end">
                                        <button class="comment-create">Comment</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </Form>
                </div>
                <!-- COMMMENTS -->
                <div v-for="comments in Comments">
                    <CommentComp :data="comments" @likesChanged="fetchData" />
                </div>
            </div>
        </div>
    </div>
    <PageNotFoundComp v-else />
    <EditPostComp :data="post" />
    <DeletePostComp :data="post" />
</template>


<script>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import axios from "axios";
import router from "../router";
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import { timeago } from "../common-js/time.js"
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
        ""
        const errorSchema = yup.object().shape({
            checkUsername: yup.boolean(),
            description: yup.string()
                .required("Required")
                .min(1, "Title must be atleast one letter")
                .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
                .max(600, "Maximum length for a comment is 600 characters"),
        });
        return {
            errorSchema,
            post: {
                Id: 1,
                Username: "",
                Created_at: "",
                Title: "",
                Description: "",
                CategoryTitle: "",
                Likes: [],
            },
            Comments: []
        };
    },
    components: {
        Form,
        Field,
    },

    created() {
        // watch the params of the route to fetch the data again
        this.$watch(
            () => this.$route.path,
            () => {
                this.fetchData()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },

    methods: {
        async fetchData() {
            let currentRouter = this.$route.path
            let correctCategory = currentRouter.split("/")
            let data = {
                categoryname: correctCategory[1],
                postid: correctCategory[2]
            }

            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/onepost", data, config)
                    .then((res) => {
                        if (res.data.message === "Post request failed") {
                            return router.go(-1)
                        }
                        this.post = res.data
                        this.Comments = res.data.Comments
                        let Cookie = res.data.Cookie
                        if (Cookie.Id.length != 0 && Cookie.Username.length != 0) {
                            createCookie(Cookie.Id, Cookie.Username)
                        }
                    })
                    .catch((error) => { });
            })
        },
        newComment(value) {
            let token = document.cookie

            let correctToken = token.split(":")
            let currentRouter = this.$route.path
            let correctCategory = currentRouter.split("/")

            if (token.length === 0) {
                return router.go(currentRouter)
            }

            const textarea = document.getElementById('newComment');
            textarea.value = null;

            let data = {
                categoryname: correctCategory[1],
                postid: correctCategory[2],
                Content: value.description,
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/createcomment", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.go(currentRouter)
                    }

                    if (res.data.message === "Malicious user detected") {
                        this.errormsg = res.data.message
                        return
                    }
                    this.Comments = res.data
                })
                .catch((error) => { });

        },

        likePost(Id) {
            let token = document.cookie
            let correctToken = token.split(":")
            let currentRouter = this.$route.path

            if (token.length === 0) {
                return router.push("/login")
            }

            let data = {
                Id: Id,
                type: "like",
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }
            axios.post("http://localhost:8080/likepost", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.push('/login')
                    }
                    if (res.data.message == "No likes") {
                        return this.post.Likes = []
                    }
                    this.post.Likes = res.data
                })
                .catch((error) => { });
        },

        dislikePost(Id) {
            let token = document.cookie
            let correctToken = token.split(":")
            let currentRouter = this.$route.path

            if (token.length === 0) {
                return router.push("/login")
            }

            let data = {
                Id: Id,
                type: "dislike",
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/dislikepost", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.push("/login")
                    }

                    if (res.data.message == "No likes") {
                        return this.post.Likes = []
                    }

                    this.post.Likes = res.data
                })
                .catch((error) => { });
        },
    },

    computed: {
        humanReadableTime() {
            return timeago(new Date(Date.now() - new Date(Date.parse(this.post.Created_at))))
        },

        humanReadableAmount() {
            if (this.Comments == null) {
                return 0
            } else {
                return this.Comments.length
            }
        },

        countLikes() {
            if (this.post.Likes == null) {
                return 0
            }
            if (this.post.Likes.length == 0) {
                return 0
            }
            const filteredLike = this.post.Likes.filter(function (key) {
                return key.Type == "like"
            })
            return filteredLike.length
        },

        countDislikes() {
            if (this.post.Likes == null) {
                return 0
            }
            if (this.post.Likes.length == 0) {
                return 0
            }
            const filteredDislike = this.post.Likes.filter(function (key) {
                return key.Type == "dislike"
            })
            return filteredDislike.length
        },

        IsLikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.post.Likes == null) {
                return false
            }


            const filteredLike = this.post.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "like"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },

        IsDislikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.post.Likes == null) {
                return false
            }

            const filteredLike = this.post.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "dislike"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },

        IsPostOwner() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (token.length == 0) {
                return false
            }

            if (this.post.Username != correctToken[1]) {
                return false
            }

            return true
        }
    },
};
</script>