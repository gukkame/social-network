<script setup>
import CommentComp from "../../components/GroupComment.vue"
import PageNotFoundComp from "../PageNotFound.vue"
</script>

<template>
    <div v-if="post.Id != 0" class="showpost">
        <div class="d-flex justify-content-center">
            <div class="showpostinfo">
                <div class="showpostinfoh d-flex flex-row">
                    <div>
                        <div v-if="post_user.Avatar_image != ``" class="bubble2 col"
                            style="margin-top: 7px; margin-left: 2px"></div>
                        <img v-else class="showpostProfileImg col" src="../../assets/images/profile.svg" />

                    </div>
                    <div class="col posthDetails">
                        <RouterLink :to="`/profile/${post_user.Username}`" style="text-decoration: none; width: 100%;"
                            class="href">
                            <div class="col showpostUser followerLink">{{ post_user.Username }}</div>
                        </RouterLink>
                        <div class="col showpostTime">{{ humanReadableTime }}</div>
                    </div>
                </div>
                <div class="showpostinfob d-flex flex-row">
                    <div class="col">
                        <div class="showpostTitle text-wrap text-break text-start col">
                            {{ post.title }}
                        </div>
                        <div class="showpostDescp text-wrap text-break text-start col">
                            {{ post.description }}
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
                            <span class="showpostReactCount" style="padding-left: 2px;">{{ post.Comments_amount
                            }}</span>
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
                                        <div class="col d-flex align-items-center justify-content-around"
                                            style="height: 30px;">
                                            <div class="d-flex col">
                                                <input type="file" name="img" id="img" accept="image/x-png,image/jpeg"
                                                    @change="handleFileUpload($event)" style="display: none" />
                                                <label style=""
                                                    class="uploadProf d-flex justify-content-center align-items-center "
                                                    for="img">Upload image</label>
                                                <label v-if="file != null" @click="deleteImage"
                                                    class="deleteProf d-flex justify-content-center align-items-center">

                                                    <i class="bi bi-trash"></i></label>
                                            </div>
                                            <div class="d-flex align-items-center col" style="padding-top: 25px"
                                                v-if="file != null">{{ displayFile }}
                                            </div>
                                        </div>
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
</template>


<script>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import axios from "axios";
import router from "../../router";
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import { timeago } from "../../common-js/time.js"
import { createCookie } from "../../common-js/cookies.js";
import { delay } from "../../common-js/time.js";

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
                Id: 0
            },
            post_likes: [],
            post_comments: [],
            post_user: {},
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
                GroupID: parseInt(person[2]),
                Id: parseInt(person[4])
            }

            await delay(100).then(() => {
                axios.post("http://localhost:8080/group/post", data, config)
                    .then((res) => {
                        if (res.data.message === "Post request failed") {
                            return router.go(-1)
                        }
                        this.post = res.data.Post
                        this.post_user = res.data.User
                        this.post_likes = res.data.Likes
                        this.Comments = res.data.Comments


                        this.$nextTick(() => {

                            let bubble = this.$el.querySelector(".bubble2")
                            if (bubble == null) {
                                return
                            }
                            bubble.style.backgroundImage = `url('http://localhost:8080${this.post_user.Avatar_image}')`
                        })
                    })
                    .catch((error) => { });
            })
        },

        handleFileUpload(event) {
            this.file = event.target.files[0];
        },

        deleteImage() {
            this.file = null;
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
            formData.append("Content", value.description)
            formData.append("Post_id", this.post.Id)

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

            axios.post("http://localhost:8080/group/post/comment/new", formData, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.go(currentRouter)
                    }

                    if (res.data.message === "Malicious user detected") {
                        this.errormsg = res.data.message
                        return
                    }
                    this.fetchData()
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
                Post_id: Id,
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }
            axios.post("http://localhost:8080/group/post/like", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.push('/login')
                    }
                    this.fetchData()
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
                Post_id: Id,
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/group/post/dislike", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.push("/login")
                    }

                    this.fetchData()
                })
                .catch((error) => { });
        },
    },

    computed: {
        yes() {
            if (this.post_user) {
                return true
            }
            return false
        },

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

        countLikes() {
            if (this.post_likes == null) {
                return 0
            }
            if (this.post_likes.length == 0) {
                return 0
            }
            const filteredLike = this.post_likes.filter(function (key) {
                return key.Type == "Like"
            })
            return filteredLike.length
        },

        countDislikes() {
            if (this.post_likes == null) {
                return 0
            }
            if (this.post_likes.length == 0) {
                return 0
            }
            const filteredDislike = this.post_likes.filter(function (key) {
                return key.Type == "Dislike"
            })
            return filteredDislike.length
        },

        IsLikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.post_likes == null) {
                return false
            }


            const filteredLike = this.post_likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "Like"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },

        IsDislikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.post_likes == null) {
                return false
            }

            const filteredLike = this.post_likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "Dislike"
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

            if (this.post_user.Username != correctToken[1]) {
                return false
            }

            return true
        },

        displayFile() {
            if (this.file != null) {
                return this.file.name
            }
        }
    },
};
</script>