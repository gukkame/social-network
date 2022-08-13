
<template>
    <div class="modal fade" id="staticBackdrop" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
        aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog  modal-dialog-centered modal-dialog-scrollable modal-xl" style="width:500px">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="createPostModal2">
                        <Form @submit="createPost" v-slot="{ errors }" :validation-schema="errorSchema">
                            <div class="container">
                                <div class="user-details row">
                                    <div class="d-flex flex-row">
                                        <div class="modal-header-title col">Create a post</div>
                                        <img v-if="this.$route.path == `/Go`" class="modal-icon col-1"
                                            src="../../assets/images/Go.svg" />
                                        <img v-if="this.$route.path == `/HTML5`" class="modal-icon col-1"
                                            src="../../assets/images/HTML.svg" />
                                        <img v-if="this.$route.path == `/CSS`" class="modal-icon col-1"
                                            src="../../assets/images/CSS.svg" />
                                        <img v-if="this.$route.path == `/JavaScript`" class="modal-icon col-1"
                                            src="../../assets/images/Javascript.svg" />
                                        <img v-if="this.$route.path == `/Vue.js`" class="modal-icon col-1"
                                            src="../../assets/images/Vue.svg" />
                                    </div>
                                    <div class="modal-sec1 col">
                                        <div class="modalTitles">
                                            Title
                                            <i class="formstar">*</i>
                                        </div>
                                        <Field class="modal-title-input" :class="{ modalinputerror: errors.title }"
                                            as="input" name="title" placeholder="60 characters" />
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
                                            name="description" placeholder="2000 characters" style="resize: none;" />
                                        <br />
                                        <span class="formErrors">{{ errors.description }}</span>
                                    </div>
                                    <div>
                                        <div>Privacy</div>
                                        <Field class="form-input" v-model="selected" style="width: 50px"
                                            :class="{ forminputerror: errors.privacy }" as="select"
                                            data-date-inline-picker="true" name="privacy" id="privacy">
                                            <option value="Public">Public</option>
                                            <option value="Followers">Followers only</option>
                                            <option value="Private">Private</option>
                                        </Field>
                                        <br />
                                        <span class="formErrors">{{ errors.privacy }}</span>
                                    </div>

                                    <div v-if="selected == `Private` && filteredUsers != 0"
                                        style="margin-top: 20px; width: 50px; position: absolute; bottom: 80px; right: 140px;">
                                        <Field v-model="selected2" :class="{ forminputerror: errors.allowedUsers }"
                                            as="select" name="allowedUsers" id="allowedUsers" multiple>
                                            <option v-for="option in  filteredUsers" :value="option.Username">{{
                                                    option.Username
                                            }}</option>
                                        </Field>
                                        <br />
                                        <span class="formErrors">{{ errors.privateUsers }}</span>
                                    </div>
                                    <div v-if="selected == `Private` && filteredUsers == 0"
                                        style="margin-top: 20px; width: 50px; position: absolute; bottom: 102px; right: 140px;">
                                        No users available
                                    </div>

                                    <div class="col d-flex align-items-center justify-content-around"
                                        style="height: 30px; margin-top: 5px;">
                                        <div class="d-flex col">
                                            <input type="file" name="img" id="img" accept="image/x-png,image/jpeg"
                                                @change="handleFileUpload($event)" style="display: none" />
                                            <label class="uploadProf d-flex justify-content-center align-items-center "
                                                for="img">Upload image</label>
                                            <label v-if="file != null" @click="deleteImage"
                                                class="deleteProf d-flex justify-content-center align-items-center">

                                                <i class="bi bi-trash"></i></label>
                                        </div>
                                        <div class="d-flex align-items-center col" style="padding-top: 25px"
                                            v-if="file != null">{{ this.file.name }}
                                        </div>

                                    </div>
                                </div>
                                <span class="formErrors">{{ errormsg }}</span>
                            </div>
                            <div class="d-flex" style="margin-top: 25px">
                                <button type="button" data-bs-dismiss="modal" class="modal-delete ms-auto"
                                    @click="deleteModalData">Cancel</button>
                                <button class="modal-create">Create</button>
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
import { delay } from "../../common-js/time.js";
import $ from 'jquery'

export default {
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
            file: null,
            selected: "Public",
            selected2: [],
            privateUsers: [],
        };
    },
    components: {
        Form,
        Field,
    },

    created() {
        this.selected = "Public"
        this.$watch(
            () => this.$route.path,
            () => {
                this.fetchAvailableUsers()
            },
            { immediate: true }
        )
    },

    computed: {
        filteredUsers() {
            let token = document.cookie
            let correctToken = token.split(":")
            if (this.privateUsers == null) {
                return this.privateUsers
            }

            const filteredUsers = this.privateUsers.filter(function (key) {
                return key.Username != correctToken[1]
            })

            return filteredUsers

        },
    },

    methods: {
        async fetchAvailableUsers() {
            let token = document.cookie
            let correctToken = token.split(":")
            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let data = {

            }

            await delay(200).then(() => {
                axios.post("http://localhost:8080/users", data, config)
                    .then((res) => {
                        if (res.data.message === "User not authenticated") {
                            return router.go(-1)
                        }
                        this.privateUsers = res.data
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

        createPost(values) {
            let currentRouter = this.$route.path
            if (document.cookie.length == 0) {
                return router.go(`"${currentRouter}"`)
            }

            let token = document.cookie
            let correctToken = token.split(":")
            let correctCategory = currentRouter.split("/")

            let config = {
                headers: {
                    header1: correctToken[0],
                }
            }

            let formData = new FormData()
            formData.append("title", values.title)
            formData.append("description", values.description)
            formData.append("title", values.title)
            formData.append("categoryname", correctCategory[1])
            formData.append("privacy", values.privacy)
            if (values.privacy == "Public" || values.privacy == "Followers") {
                formData.append("allowedUsers", null)
            } else {
                formData.append("allowedUsers", values.allowedUsers)
            }
            if (this.file == null) {
                formData.append("img", "");
            } else {
                formData.append("img", this.file);
            }


            axios.post("http://localhost:8080/createpost", formData, config)
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
                        return router.push("/")
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
                    router.push(`/${correctCategory[1]}/${res.data.Id}`)

                })
                .catch((error) => { });
        },

        deleteModalData() {
            $('#staticBackdrop form')[0].reset();
        },
    },

};
</script>

<style>
@import "../../assets/css/base.css";
</style>
