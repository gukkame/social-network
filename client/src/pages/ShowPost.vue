<script setup>
import CommentComp from "../components/Comment.vue"
import PageNotFoundComp from "./PageNotFound.vue"
import EditPostComp from "../components/modals/EditPost.vue"
import DeletePostComp from "../components/modals/DeletePost.vue"
</script>

<template>
    <PageNotFoundComp v-if="post.Id == 0" />
    <div v-else class="showpost">
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
                        <div v-if="post_user.Avatar_image != ``" class="bubble22 col"
                            style="margin-top: 7px; margin-left: 2px"></div>
                        <img v-else class="showpostProfileImg col" src="../assets/images/profile.svg" />
                    </div>
                    <div class="col posthDetails">
                          <RouterLink :to="`/profile/${post.Username}`" style="text-decoration: none; width: 100%;"
                            class="href">
                            <div class="col showpostUser followerLink">{{ post.Username }}</div>
                        </RouterLink>
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
                        <div class="showpostDescp text-wrap text-break text-start col">
                            <img v-if="displayImage" style="width: 250px; height: auto"
                                :src="`http://localhost:8080${post.Image}`">
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
                                            <div class="col d-flex align-items-center justify-content-around"
                                                style="height: 30px;">
                                                <div class="d-flex col">
                                                    <input type="file" name="img" id="img"
                                                        accept="image/x-png,image/jpeg"
                                                        @change="handleFileUpload($event)" style="display: none" />
                                                    <label style=""
                                                        class="uploadProf d-flex justify-content-center align-items-center "
                                                        for="img">Upload image</label>
                                                    <label v-if="file != null" @click="deleteImage"
                                                        class="deleteProf d-flex justify-content-center align-items-center">

                                                        <i class="bi bi-trash"></i></label>
                                                </div>
                                                <div class="d-flex align-items-center col" style="padding-top: 25px"
                                                    v-if="file != null">{{ this.file.name }}
                                                </div>
                                            </div>
                                            <button class="comment-create">Comment</button>
                                        </div>
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
        <EditPostComp :data="post" />
        <DeletePostComp :data="post" />
    </div>
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
            post: {},
            post_user: {},
            Comments: [],
            file: null
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

    updated() {

    },

    methods: {

        handleFileUpload(event) {
            this.file = event.target.files[0];
        },

        deleteImage() {
            this.file = null;
        },

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
                        /*  if (res.data.message === "Post request failed") {
                             return router.go(-1)
                         } */
                        this.post = res.data
                        this.Comments = res.data.Comments
                        this.post_user = res.data.User
                        let bubble = this.$el.querySelector(".bubble22")
                        if (bubble == null) {
                            return
                        }

                        bubble.style.backgroundImage = `url('http://localhost:8080${this.post_user.Avatar_image}')`

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

            let formData = new FormData()
            formData.append("categoryname", correctCategory[1])
            formData.append("postid", correctCategory[2])
            formData.append("content", value.description)
            if (this.file == null) {
                formData.append("img", "");
            } else {
                formData.append("img", this.file);
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/createcomment", formData, config)
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
        displayImage() {
            if (this.post.Image) {
                if (this.post.Image.length == 0) {
                    return false
                }
                return true
            }
        },
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