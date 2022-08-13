<template>
    <div class="d-flex justify-content-center col">
        <RouterLink v-if="data.User_status == `Owner`" class="row text" :to="{ name: 'admin', params: { id: returnPath } }"
            style="text-decoration:none">
            <div :class="[selected == 'admin' ? 'active' : '']"
                class="selectorCard d-flex justify-content-center row">
                <div class="d-flex justify-content-center align-middle selectorText">Admin</div>
            </div>
        </RouterLink>

        <RouterLink class="row text" :to="{ name: 'groupposts', params: { id: returnPath } }" style="text-decoration:none">
            <div :class="[selected == 'groupposts' ? 'active' : '']"
                class="selectorCard d-flex justify-content-center row">
                <div class="d-flex justify-content-center align-middle selectorText">Posts</div>
            </div>
        </RouterLink>

        <RouterLink class="row text" :to="{ name: 'events', params: { id: returnPath } }"
            style="text-decoration:none">
            <div :class="[selected == 'events' ? 'active' : '']" class="selectorCard d-flex justify-content-center">
                <div class="d-flex justify-content-center align-middle selectorText">Events</div>
            </div>
        </RouterLink>
    </div>
</template>


<script>
export default {
    props: {
        selected: {
            type: String,
            required: true
        },
        data: {
            type: Object,
            required: true
        }
    },
    computed: {
        returnPath() {
            let path = this.$route.path.split("/")
            return path[2]
        }
    }
}


</script>


<style>
.selectorCard {
    border-radius: 5px;
    margin: -52px 23px 23px 23px;
    width: 200px;
    height: 50px;
    box-shadow: 0 2px 6px 0 rgb(218 218 253 / 65%), 0 2px 6px 0 rgb(206 206 238 / 54%);
    display: inline-block;
    transition: 0.3s;
}

.selectorCard:hover {
    transform: scale(1.13);
    transition: 0.3s;
    cursor: pointer;
}

.selectorText {
    text-decoration: none;
    margin-top: 4px;
    color: #2E343D;
    font-size: 25px;
    font-weight: 100px;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    -o-user-select: none;
    user-select: none;
}

.active {
    border: solid 2px #FF9D5A;
}
</style>