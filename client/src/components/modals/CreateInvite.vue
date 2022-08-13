<template>
    <div class="modal fade" id="inviteGroup" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog  modal-dialog-centered modal-dialog-scrollable modal-xl" style="width:500px">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="createGroupInviteModal">
                        <Form @submit="createInvite" v-slot="{ errors }" :validation-schema="errorSchema">
                            <div class="container">
                                <div class="user-details row">
                                    <div class="d-flex flex-row">
                                        <div class="modal-header-title col">Invite a friend</div>
                                    </div>
                                    <div class="modal-sec1 col">
                                        <div class="modalTitles">
                                            <i class="formstar">*</i>
                                        </div>
                                        <Field class="modal-title-input" :class="{ modalinputerror: errors.name }"
                                            as="input" name="name" placeholder="Name" />
                                        <br />
                                        <span class="formErrors">{{ errors.name }}</span>
                                    </div>

                                </div>
                                <span class="formErrors">{{ errormsg }}</span>
                            </div>
                            <div class="d-flex">
                                <button type="button" data-bs-dismiss="modal" class="modal-delete ms-auto"
                                    @click="deleteModalData">Cancel</button>
                                <button class="modal-create">Invite</button>
                            </div>
                        </Form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import axios from "axios";
import router from "../../router";
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import $ from 'jquery'
import {  ws } from "../../common-js/messages.js";


export default {
    data() {
        const errorSchema = yup.object().shape({
            name: yup.string()
                .required("Required")
                .test(
                    "Unique username",
                    "User does not exist!",
                    async function (value) {
                        const payload = {
                            username: `${value}`,
                        };

                        try {
                            const res = await axios.post(
                                "http://localhost:8080/available/username",
                                payload
                            );
                            if (res.data.value === "true") {
                                return false;
                            }
                            return true;
                        } catch (error) { }
                    }
                ),
        });
        return {
            errorSchema,
            errormsg: ""
        };
    },
    components: {
        Form,
        Field,
    },

    methods: {
        createInvite(values) {
            this.errormsg = ""
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
                username: values.name,
                group_id: parseInt(person[2])
            }

            axios.post("http://localhost:8080/group/invite", data, config)
                .then((res) => {
                    if (res.data.message === "Already Member") {
                        return this.errormsg = res.data.message
                    }

                    if (res.data.messag === "Already Invited") {
                        return this.errormsg = res.data.message
                    }

                    if (res.data.message === "User not authenticated") {
                        $('body').removeClass('modal-open');
                        $('.modal-backdrop').hide()
                        async function removeAllAttrs(element) {
                            for (var i = element.attributes.length; i-- > 0;)
                                await element.removeAttributeNode(element.attributes[i]);
                        }
                        removeAllAttrs(document.body);
                        $('body').css('overflow', 'auto');
                        return router.push(`"${currentRouter}"`)
                    }

                    let payload = {
                        Type: "GroupInvNotif",
                        Content: {
                            Message: "Invite to group",
                            Sender: correctToken[1],
                            Receiver: values.name,
                            IsGroup: parseInt(person[2]),
                        },

                    };
                    ws.send(JSON.stringify(payload));
                    if (ws.readyState === WebSocket.CLOSED) {
                        clearInterval(this.timer);
                        return;
                    }
                    this.errormsg = ""
                    $('body').removeClass('modal-open');
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

        deleteModalData() {
            $('#inviteGroup form')[0].reset();
        }
    },
};
</script>
