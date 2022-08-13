<template>
     <RouterLink v-if="notification.Type == 'followNotif'" :to="`/profile/${getNameFromToken}/followers`" class="notificationa">
    <div class="notificationwindowPerson d-flex overflow-auto">
        <div class="col d-flex">
            <div class="d-flex col-2 justify-content-center">
                <i class="bi bi-envelope-open align-self-center"></i>
            </div>

            <div class="d-flex col justify-content-start flex-column">
                <div class="col">
                    {{ notification.Content }}
                </div>
                <div class="col postTime">
                    {{ humanReadableTime }}
                </div>

            </div>
        </div>
    </div>
    </RouterLink>

    <RouterLink v-if="notification.Type == 'GroupInvNotif'" :to="`/groups/${notification.Group_id}/posts`" class="notificationa">
        <div class="notificationwindowPerson d-flex overflow-auto">
            <div class="col d-flex">
                <div class="d-flex col-2 justify-content-center">
                    <i class="bi bi-envelope-open align-self-center"></i>
                </div>

                <div class="d-flex col justify-content-start flex-column">
                    <div class="col">

                        {{ notification.Content }}
                    </div>
                    <div class="col postTime">
                        {{ humanReadableTime }}
                    </div>

                </div>
            </div>
        </div>


    </RouterLink>

    <RouterLink v-if="notification.Type == 'GroupJoinReqNotif'" :to="`/groups/${notification.Group_id}/admin`" class="notificationa">
        <div class="notificationwindowPerson d-flex overflow-auto">
            <div class="col d-flex">
                <div class="d-flex col-2 justify-content-center">
                    <i class="bi bi-envelope-open align-self-center"></i>
                </div>

                <div class="d-flex col justify-content-start flex-column">
                    <div class="col">

                        {{ notification.Content }}
                    </div>
                    <div class="col postTime">
                        {{ humanReadableTime }}
                    </div>

                </div>
            </div>
        </div>
    </RouterLink>

    <RouterLink v-if="notification.Type == 'NewEventNotif'" :to="`/groups/${notification.Group_id}/events`" class="notificationa">
        <div class="notificationwindowPerson d-flex overflow-auto">
            <div class="col d-flex">
                <div class="d-flex col-2 justify-content-center">
                    <i class="bi bi-envelope-open align-self-center"></i>
                </div>

                <div class="d-flex col justify-content-start flex-column">
                    <div class="col">

                        {{ notification.Content }}
                    </div>
                    <div class="col postTime">
                        {{ humanReadableTime }}
                    </div>

                </div>
            </div>
        </div>
    </RouterLink>
</template>

<script>
import { timeago } from '../common-js/time'
export default {
    props: {
        data: {
            type: Object,
            required: true
        }
    },

    data() {
        return {
            notification: null
        }
    },

    created() {
        this.$watch(
            () => this.data,
            () => {
                this.returnData()
            },
            // fetch the data when the view is created and the data is
            // already being observed
            { immediate: true }
        )
    },
    methods: {
        returnData() {
            this.notification = this.data
        }
    },


    computed: {
        humanReadableTime() {
            return timeago(new Date(Date.now() - new Date(Date.parse(this.notification.Created_at ?? "01 Jan 1970 00:00:00 GMT"))))
        },
        getNameFromToken() {
            let token = document.cookie
            let correctToken = token.split(":")
            return correctToken[1]
        }
    },
}
</script>





<style>
.notificationa {
    text-decoration: none;
    color: #212529
}

.notificationa:hover {
     text-decoration: none;
    color: #212529
}

.notificationwindowPerson {
    transition: 0.3s;
    cursor: pointer;
    padding: 4px;
    border-bottom: solid 1px rgb(219, 219, 219);
    width: 100%;
    height: fit-content;
    word-break: break-word;
    padding: 5px 12px 5px 0px;
    font-size: 14px
}

.notificationwindowPerson:hover {
    transition: 0.3s;
    background-color: #D3D3D3;
    padding: 5px 12px 5px 0px;
    border-bottom: solid 1px rgb(219, 219, 219);
}

.notificationDetails {
    padding: 2px;
    margin-left: 10px;
    max-width: fit-content;
    display: flex;
    flex-wrap: wrap;
}
</style>