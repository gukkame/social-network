<template>
    <div class="comment">
        <div class="commenth d-flex flex-row">
            <div>
                <img class="commentProfileImg col" src="../assets/images/profile.svg" />
            </div>
            <div class="col posthDetails">
                <div class="col commentUser">{{ data.Username }}</div>
                <div class="col commentTime">{{ humanReadableTime }}</div>
            </div>
        </div>
        <div class="commentb d-flex flex-row">
            <div class="col">
                <div class="commentDescp text-wrap text-break text-start col">
                    {{ data.Content }}
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

    methods: {
        likeComment(Id) {
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

            axios.post("http://localhost:8080/likecomment", data, config)
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
                Id: Id,
                type: "dislike",
            }

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            axios.post("http://localhost:8080/dislikecomment", data, config)
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
        humanReadableTime() {
            return timeago(new Date(Date.now() - new Date(Date.parse(this.data.Created_at))))
        },

        countLikes() {
            if (this.data.Likes == null) {
                return 0
            }
            const filteredLike = this.data.Likes.filter(function (key) {
                return key.Type == "like"
            })
            return filteredLike.length
        },

        countDislikes() {
            if (this.data.Likes == null) {
                return 0
            }
            const filteredDislike = this.data.Likes.filter(function (key) {
                return key.Type == "dislike"
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
            if (this.data.Likes == null) {
                return false
            }
            const filteredLike = this.data.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "dislike"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },
    }
}

</script>