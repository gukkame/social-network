<template>
    <nav class="main-menu">
        <ul>
            <li>
                <RouterLink to="/" class="href">
                    <i class="fa fa-home fa-2x"></i>
                    <span class="nav-text">
                        Home
                    </span>
                </RouterLink>

            </li>
            <li>
                <RouterLink to="/groups" class="href">
                    <i class="fa fa-solid fa-users fa-2x"></i>
                    <span class="nav-text">
                        Groups
                    </span>
                </RouterLink>
            </li>
            <li v-if="LoggedIn">
                <RouterLink :to="{ name: 'profile', params: { id: returnPath } }" class="href">
                    <i class="fa fa-solid fa-user fa-2x"></i>
                    <span class="nav-text">
                        Profile
                    </span>
                </RouterLink>

            </li>
        </ul>

        <ul class="logout" v-if="LoggedIn">
            <li @click="logOut">
                <a class="href">
                    <i class="fa fa-power-off fa-2x"></i>
                    <span class="nav-text">
                        Logout
                    </span>
                </a>
            </li>
        </ul>
    </nav>
</template>

<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import router from "../../router";
import { closeWsConnection } from "../../common-js/messages.js";
export default {

    name: 'Header',
    props: {
        LoggedIn: {
            type: Boolean,
            required: true
        },
    },

    methods: {
        logOut() {
            document.cookie = 'Token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            closeWsConnection()
            router.go("/")
        }
    },

    computed: {
        returnPath() {
            let path = document.cookie.split(":")
            return path[1]
        }
    }
}
</script>