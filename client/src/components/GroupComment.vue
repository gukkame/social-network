<template>
    <div class="comment">
        <div class="commenth d-flex flex-row">
            <div>
                <div v-if="data.User.Avatar_image != ``" class="bubble4 col" style="margin-top: 7px; margin-left: 2px;">
                </div>
                <img v-else class="commentProfileImg col" src="../assets/images/profile.svg" />
            </div>
            <div class="col posthDetails">
                <RouterLink  v-if="data"  :to="{ name: 'profile', params: { id: data.User.Username } }"
                    style="text-decoration: none; width: 100%;" class="href">
                    <div class="col commentUser followerLink">{{ data.User.Username }}</div>
                </RouterLink>

                <div class="col commentTime">{{ humanReadableTime }}</div>
            </div>
        </div>
        <div class="commentb d-flex flex-row">
            <div class="col">
                <div class="commentDescp text-wrap text-break text-start col">
                    {{ data.Content }}
                </div>
                <div class="commentDescp text-wrap text-break text-start col">
                    <img v-if="displayImage" style="width: 250px; height: auto"
                        :src="`http://localhost:8080${data.Image}`">
                </div>
            </div>
        </div>
        <div class="commentf d-flex flex-row allign-items-start">
            <div class="d-flex col">
                <div style="margin-left: 10px">
                    <i type="button" @click="likeComment(data.Id)" :class="{ likeActive: IsLikeActive }"
                        class="commentLike align-self-center  bi bi-hand-thumbs-up"></i>
                    <span class="commentReactCount align-self-center">{{ countLikes }}</span>
                </div>
                <div style="margin-left: 20px">
                    <i type="button" @click="dislikeComment(data.Id)" :class="{ dislikeActive: IsDislikeActive }"
                        class="commentDislike bi bi-hand-thumbs-down"></i>
                    <span class="commentReactCount ">{{ countDislikes }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { timeago } from "../common-js/time.js"
import axios from "axios";
import router from "../router";
export default {
    props: {
        data: {
            type: Object,
            required: true
        }
    },

    mounted() {
        let bubble = this.$el.querySelector(".bubble4")
        if (bubble == null) {
            return
        }
        bubble.style.backgroundImage = `url('http://localhost:8080${this.data.User.Avatar_image}')`
    },

    methods: {
        likeComment(Id) {
            let token = document.cookie
            let correctToken = token.split(":")
            let currentRouter = this.$route.path

            if (token.length === 0) {
                return router.push("/login")
            }

            let data = {
                Comment_id: Id,
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/group/post/comment/like", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.push('/login')
                    }
                    this.$emit('likesChanged', true)
                })
                .catch((error) => { });
        },

        dislikeComment(Id) {
            let token = document.cookie
            let correctToken = token.split(":")
            let currentRouter = this.$route.path

            if (token.length === 0) {
                return router.push("/login")
            }

            let data = {
                Comment_id: Id,
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/group/post/comment/dislike", data, config)
                .then((res) => {
                    if (res.data.message === "User not authenticated") {
                        return router.push("/login")
                    }

                    this.$emit('likesChanged', true)
                })
                .catch((error) => { });
        }
    },

    computed: {
        displayImage() {
            if (this.data.Image.length == 0) {
                return false
            }
            return true
        },

        humanReadableTime() {
            return timeago(new Date(Date.now() - new Date(Date.parse(this.data.Created_at))))
        },

        countLikes() {
            if (this.data.Likes == null) {
                return 0
            }
            const filteredLike = this.data.Likes.filter(function (key) {
                return key.Type == "Like"
            })
            return filteredLike.length
        },

        countDislikes() {
            if (this.data.Likes == null) {
                return 0
            }
            const filteredDislike = this.data.Likes.filter(function (key) {
                return key.Type == "Dislike"
            })
            return filteredDislike.length
        },

        IsLikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Likes == null) {
                return false
            }
            const filteredLike = this.data.Likes.filter(function (key) {
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
            if (this.data.Likes == null) {
                return false
            }
            const filteredLike = this.data.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "Dislike"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },
    }
}

</script>