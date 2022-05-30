<template>
    <RouterLink :to="{ name: data.Category_title, params: { id: data.Post_id } }" style="text-decoration: none">
        <div class="activityComment">
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
                        <i type="button" :class="{ likeActive: IsLikeActive, postVoteColor: NoLikeActive }"
                            class="activitycommentLike align-self-center  bi bi-hand-thumbs-up"></i>
                        <span class="commentReactCount align-self-center">{{ countLikes }}</span>
                    </div>
                    <div style="margin-left: 20px">
                        <i type="button" :class="{ dislikeActive: IsDislikeActive, postVoteColor: NoDislikeActive }"
                            class="activitycommentLike bi bi-hand-thumbs-down"></i>
                        <span class="commentReactCount ">{{ countDislikes }}</span>
                    </div>
                </div>
            </div>
        </div>
    </RouterLink>
</template>

<script>
import { timeago } from "../common-js/time.js"
export default {
    props: {
        data: {
            type: Object,
            required: true
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

                NoLikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Likes == null) {
                return true
            }
            const filteredLike = this.data.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "like"
            })

            if (filteredLike.length != 0) {
                return false
            }

            return true
        },

        NoDislikeActive() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.data.Likes == null) {
                return true
            }
            const filteredLike = this.data.Likes.filter(function (key) {
                return key.Username == correctToken[1] &&
                    key.Type == "dislike"
            })

            if (filteredLike.length != 0) {
                return false
            }

            return true
        },
    }
}
</script>