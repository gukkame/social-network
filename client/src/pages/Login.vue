<template>
    <div class="login">
        <div class="formIcon">
            <img class="loginformImg" src="../assets/images/logo.svg" />
        </div>
        <div class="loginform">
            <Form v-slot="{ errors }" @submit="login" :validation-schema="errorSchema">
                <div class="container">
                    <div class="account-details la row">
                        <div class="col">
                            <div>Username/Email<i class="formstar">*</i></div>
                            <Field class="form-input" :class="{ forminputerror: errors.username }" as="input"
                                name="username" />
                            <br />
                            <span class="formErrors">{{ errors.username }}</span>
                        </div>
                    </div>

                    <div class="passwords la row">
                        <div class="col">
                            <div>Password<i class="formstar">*</i></div>
                            <Field class="form-input" :class="{ forminputerror: errors.password }" as="input"
                                name="password" type="password" />
                            <br />
                            <span class="formErrors">{{ errors.password }}</span>
                        </div>
                    </div>
                    <span class="formErrors">{{ errormsg }}</span>
                    <button class="nav-button-2 nav-items-2 slogin">Log In</button>
                </div>
            </Form>
            <RouterLink class="loginformDirectText" to="/signup">
                <i>New here? Sign Up</i>
            </RouterLink>
        </div>
    </div>
</template>

<script>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import axios from "axios";
import router from "../router";
import { createCookie } from "../common-js/cookies.js"

export default {
    components: {
        Form,
        Field,
    },

    setup() {
        const errorSchema = yup.object().shape({
            username: yup.string()
                .required("Required"),
            password: yup.string()
                .required('Required')
        });
        return {
            errorSchema,
        };
    },

    methods: {
        async login(values) {
            try {
                const res = await axios.post('http://localhost:8080/login', values)
                if (res.data.message === "Invalid credentials") {
                    return this.errormsg = res.data.message
                }

                createCookie(res.data.Id, res.data.Username)
                this.errormsg = ""
                router.push("/")
            } catch (error) {
                console.log(error)
            }
        }
    },

    mounted() {
        let token = document.cookie
        if (token.length != 0) {
            return router.push("/")
        }
    }
};
</script>