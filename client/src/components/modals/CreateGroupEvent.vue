<template>
    <div class="modal fade" id="groupEvent" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog  modal-dialog-centered modal-dialog-scrollable modal-xl" style="width:500px">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="createGroupEventModal">
                        <Form @submit="createEvent" v-slot="{ errors }" :validation-schema="errorSchema">
                            <div class="container">
                                <div class="user-details row">
                                    <div class="d-flex flex-row">
                                        <div class="modal-header-title col">Create an event</div>
                                    </div>
                                    <div class="modal-sec1 col">
                                        <div class="modalTitles">
                                            Title
                                            <i class="formstar">*</i>
                                        </div>
                                        <Field class="modal-title-input" :class="{ modalinputerror: errors.title }"
                                            as="input" name="title" placeholder="60 characters" v-model="data.Title" />
                                        <br />
                                        <span class="formErrors">{{ errors.title }}</span>
                                    </div>
                                    <div class="modal-sec2 col">
                                        <div class="modalTitles">
                                            Description
                                            <i class="formstar">*</i>
                                        </div>
                                        <Field class="modal-description-input"
                                            :class="{ modalinputerror: errors.description }" as="textarea"
                                            name="description" style="resize: none;" placeholder="2000 characters"
                                            v-model="data.Description" />
                                        <br />
                                        <span class="formErrors">{{ errors.description }}</span>
                                    </div>
                                    <div class="modal-sec1 col">
                                        <div class="col">
                                            <div class="modalTitles">
                                                Date
                                                <i class="formstar">*</i>
                                            </div>
                                            <Field class="form-input" :class="{ forminputerror: errors.date }"
                                                type="date" name="date" data-date-inline-picker="true" />

                                            <br />
                                            <span class="formErrors">{{ errors.date }}</span>
                                        </div>
                                    </div>
                                </div>
                                <span class="formErrors">{{ errormsg }}</span>
                            </div>
                            <div class="d-flex">
                                <button type="button" data-bs-dismiss="modal"
                                    class="modal-delete ms-auto">Cancel</button>
                                <button class="modal-create">Confirm</button>
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
        ""
        const errorSchema = yup.object().shape({
            checkUsername: yup.boolean(),
            title: yup.string()
                .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
                .required("Required")
                .max(60, "Maximum length for a title is 60 characters"),

            description: yup.string()
                .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
                .required("Required")
                .max(500, "Maximum length for a description is 500 characters"),
            date: yup
                .string()
                .test("date", "Event can't be happening in the past!", (date) => {
                    const today = new Date();

                    if (new Date(date) < today) {
                        return false
                    }
                    return true

                })
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
        createEvent(values) {
            let currentRouter = this.$route.path
            if (document.cookie.length == 0) {
                return router.go(`"${currentRouter}"`)
            }

            let token = document.cookie
            let correctToken = token.split(":")
            let groupId = currentRouter.split("/")

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let payload = {
                Title: values.title,
                Content: values.description,
                Happening_at: values.date,
                Group_id: parseInt(groupId[2])
            }

            axios.post("http://localhost:8080/group/event/new", payload, config)
                .then((res) => {
                    if (res.data.message === "Malicious user detected") {
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
                        return router.go(`"${currentRouter}"`)
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
    },
};
</script>
