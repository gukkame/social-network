<template>
    <div class="modal fade" id="leaveGroup" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog  modal-dialog-centered modal-dialog-scrollable modal-xl" style="width:500px">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="deletePostModal">
                        <div class="container">
                            <div class="user-details row">
                                <div class="d-flex flex-row">
                                    <div class="modal-header-title col">Are you sure you want to leave this group?</div>
                                </div>
                            </div>
                            <span class="formErrors">{{ errormsg }}</span>
                        </div>
                        <div class="d-flex">
                            <button type="button" data-bs-dismiss="modal" class="modal-delete ms-auto">Cancel</button>
                            <button class="modal-create" @click="leaveGroup">Confirm</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from "axios";
import router from "../../router";
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import $ from 'jquery'

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

    methods: {
        leaveGroup() {
            let currentRouter = this.$route.path
            if (document.cookie.length == 0) {
                return router.go(`"${currentRouter}"`)
            }

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
                GroupID: parseInt(person[2])
            }

            axios.post("http://localhost:8080/group/leave", data, config)
                .then((res) => {
                    if (res.data.message === "Malicious user detected") {
                        return this.errormsg = res.data.message
                    }

                    if (res.data.message === "User not authenticated") {
                        $('body').removeClass('modal-open');
                        /* $('#staticBackdrop').hide() */
                        $('.modal-backdrop').hide()
                        async function removeAllAttrs(element) {
                            for (var i = element.attributes.length; i-- > 0;)
                                await element.removeAttributeNode(element.attributes[i]);
                        }
                        removeAllAttrs(document.body);
                        $('body').css('overflow', 'auto');
                        return router.go(`"${currentRouter}"`)
                    }

                    this.errormsg = ""
                    $('body').removeClass('modal-open');
                    /* $('#staticBackdrop').hide() */
                    $('.modal-backdrop').hide()
                    async function removeAllAttrs(element) {
                        for (var i = element.attributes.length; i-- > 0;)
                            await element.removeAttributeNode(element.attributes[i]);
                    }
                    removeAllAttrs(document.body);
                    $('body').css('overflow', 'auto');
                    router.go(`"${currentRouter}"`)

                })
                .catch((error) => { });
        },
    },
};
</script>
