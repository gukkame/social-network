<template>
    <div class="postcard">
        <RouterLink class="d-flex flex-column flex-wrap justify-content-space"
            :to="{ name: `onegrouppost`, params: { postid: data.Post.Id } }">
            <div class="postcardh d-flex flex-row">
                <div>
                    <img v-if="data.User.Avatar_image == ``" class="postProfileImg col" src="../assets/images/profile.svg" />
                    <div v-else class="bubble3 col" style="margin-top: 7px; margin-left: 2px"></div>
                </div>
                <div class="col posthDetails">
                    <RouterLink :to="{ name: 'profile', params: { id: data.User.Username } }"
                        style="text-decoration: none; width: 100%;" class="href">
                        <div class="col postUser followerLink">{{ data.User.Username }}</div>
                    </RouterLink>
                    <div class="col postTime">{{ humanReadableTime }}</div>
                </div>
            </div>
            <div class="postcardb d-flex flex-row">
                <div class="col">
                    <div class="postTitle text-wrap text-break text-start col">
                        {{ data.Post.title }}
                    </div>
                    <div class="postDescp text-wrap text-break text-start col">
                        {{ maxlengthdesc }}
                    </div>
                </div>
            </div>
            <div class="postcardf d-flex flex-row">
                <div class="col d-flex justify-content-center">
                    <i :class="{ likeActive: IsLikeActive, postVoteColor: NoLikeActive }"
                        class="postLike align-self-center  bi bi-hand-thumbs-up"></i>
                    <span class="postReactCount align-self-center" style="padding-top: 6px">{{ countLikes }}</span>
                </div>
                <div class="col d-flex justify-content-center">
                    <div class="align-self-center">
                        <i class="postComment bi bi-chat-dots"></i>
                        <span class="postReactCount" style="padding-left: 2px;">{{ data.Post.Comments_amount }}</span>
                    </div>
                </div>
                <div class="col d-flex justify-content-center">
                    <div class="align-self-center ">
                        <i :class="{ dislikeActive: IsDislikeActive, postVoteColor: NoDislikeActive }"
                            class="postDislike bi bi-hand-thumbs-down"></i>
                        <span class="postReactCount ">{{ countDislikes }}</span>
                    </div>
                </div>
            </div>
        </RouterLink>
    </div>
</template>

<script>
import { timeago } from "../common-js/time.js"

export default {
    props: {
        data: {
            type: Object,
            required: true,
            default: {
                Created_at: "01 Jan 1970 00:00:00 GMT"
            }
        },
    },

    mounted() {
        let bubble = this.$el.querySelector(".bubble3")
        if (bubble == null) {
            return
        }
        bubble.style.backgroundImage = `url('http://localhost:8080${this.data.User.Avatar_image}')`
    },

    computed: {
        humanReadableTime() {
            return timeago(new Date(Date.now() - new Date(Date.parse(this.data.Post.Created_at ?? "01 Jan 1970 00:00:00 GMT"))))
        },

        maxlengthdesc() {
            if (this.data.Post.description.length > 300) {
                return this.data.Post.description.substring(0, 300) + "..."
            }
            return this.data.Post.description
        },

        countLikes() {
            if (this.data.Post.Likes == null) {
                return 0
            }
            const filteredLike = this.data.Post.Likes.filter(function (key) {
                return key.Type == "Like"
            })
            return filteredLike.length
        },

        countDislikes() {
            if (this.data.Post.Likes == null) {
                return 0
            }
            const filteredDislike = this.data.Post.Likes.filter(function (key) {
                return key.Type == "Dislike"
            })
            return filteredDislike.length
        },

        IsLikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Post.Likes == null) {
                return false
            }
            const filteredLike = this.data.Post.Likes.filter(function (key) {
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
            if (this.data.Post.Likes == null) {
                return false
            }
            const filteredLike = this.data.Post.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "Dislike"
            })

            if (filteredLike.length == 0) {
                return false
            }

            return true
        },

        NoLikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Post.Likes == null) {
                return true
            }
            const filteredLike = this.data.Post.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "Like"
            })

            if (filteredLike.length != 0) {
                return false
            }

            return true
        },

        NoDislikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Post.Likes == null) {
                return true
            }
            const filteredLike = this.data.Post.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "Dislike"
            })

            if (filteredLike.length != 0) {
                return false
            }

            return true
        },
    },
}


</script>

<style>
@import "../assets/css/base.css";
</style>