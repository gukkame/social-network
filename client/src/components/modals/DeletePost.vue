<template>
    <div class="modal fade" id="deletepost" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog  modal-dialog-centered modal-dialog-scrollable modal-xl" style="width:500px">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="deletePostModal">
                            <div class="container">
                                <div class="user-details row">
                                    <div class="d-flex flex-row">
                                        <div class="modal-header-title col">Are you sure you want to delete this post?</div>
                                    </div>
                                </div>
                                <span class="formErrors">{{ errormsg }}</span>
                            </div>
                            <div class="d-flex">
                                <button type="button" data-bs-dismiss="modal"
                                    class="modal-delete ms-auto">Cancel</button>
                                <button class="modal-create" @click="deletePost">Confirm</button>
                            </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { Form, Field, ErrorMessage } from "vee-validate";
import * as yup from "yup";
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

    data() {
        const errorSchema = yup.object().shape({
            checkUsername: yup.boolean(),
            title: yup.string()
                .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
                .required("Required")
                .max(60, "Maximum length for a title is 60 characters"),

            description: yup.string()
                .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
                .required("Required")
                .max(1000, "Maximum length for a description is 1000 characters"),
        });
        return {
            errorSchema,
        };
    },
    components: {
        Form,
        Field,
    },

    methods: {
        deletePost(values) {
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

            let payload = {
                Id: this.data.Id,
            }

            axios.post("http://localhost:8080/deletepost", payload, config)
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
                    router.push("/")

                })
                .catch((error) => { });
        },
    },
};
</script>
