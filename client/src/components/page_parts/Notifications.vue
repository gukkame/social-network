<template>
    <div v-if="this.data.Content != undefined">{{ notif }}</div>
    <div id="toast">
        <RouterLink v-if="this.type=='followNotif'" :to="`/profile/${data.Content.Receiver}/followers`" id="img">Reply!</RouterLink>

        <RouterLink v-if="this.type=='GroupInvNotif'" :to="`/groups/${data.Content.IsGroup}/posts`" id="img">Reply!</RouterLink>

        <RouterLink v-if="this.type=='GroupJoinReqNotif'" :to="`/groups/${data.Content.IsGroup}/admin`" id="img">Reply!</RouterLink>

        <RouterLink v-if="this.type=='NewEventNotif'" :to="`/groups/${data.Content.IsGroup}/events`" id="img">See!</RouterLink>
        <div id="desc">{{ msg }}</div>
    </div>

</template>
<script>


export default {
    props: {
        data: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            msg: "",
            type: "",
        };
    },
    computed: {
        notif() {
            this.msg = this.data.Content.Message
            this.type = this.data.Type
            return this.launch_toast()
        }
    },
    methods: {
        launch_toast() {
            if (document.getElementById("toast")) {
                var x = document.getElementById("toast")
                x.className = "show";
                setTimeout(function () { x.className = x.className.replace("show", ""); }, 5000);

            }

        },
    }
}


</script>
<style>
#toast {
    visibility: hidden;
    max-width: 50px;
    height: 50px;
    margin: auto;
    background-color: #333;
    color: #fff;
    text-align: center;
    border-radius: 2px;
    position: fixed;
    z-index: 1;
    left: 0;
    right: 0;
    bottom: 30px;
    font-size: 17px;
    white-space: nowrap;
}

#toast #img {
    width: 50px;
    height: 50px;
    float: left;
    padding-top: 13px;
    padding-bottom: 13px;
    box-sizing: border-box;
    background-color: #111;
    color: #fff;
}

#toast #desc {
    color: #fff;
    padding: 13px;
    overflow: hidden;
    white-space: nowrap;
}

#toast.show {
    visibility: visible;
    -webkit-animation: fadein 0.5s, expand 0.5s 0.5s, stay 3s 1s, shrink 0.5s 2s, fadeout 0.5s 2.5s;
    animation: fadein 0.5s, expand 0.5s 0.5s, stay 3s 1s, shrink 0.5s 4s, fadeout 0.5s 4.5s;
}

@-webkit-keyframes fadein {
    from {
        bottom: 0;
        opacity: 0;
    }

    to {
        bottom: 30px;
        opacity: 1;
    }
}

@keyframes fadein {
    from {
        bottom: 0;
        opacity: 0;
    }

    to {
        bottom: 30px;
        opacity: 1;
    }
}

@-webkit-keyframes expand {
    from {
        min-width: 50px
    }

    to {
        min-width: 380px
    }
}

@keyframes expand {
    from {
        min-width: 50px
    }

    to {
        min-width: 380px
    }
}

@-webkit-keyframes stay {
    from {
        min-width: 380px
    }

    to {
        min-width: 380px
    }
}

@keyframes stay {
    from {
        min-width: 380px
    }

    to {
        min-width: 380px
    }
}

@-webkit-keyframes shrink {
    from {
        min-width: 380px;
    }

    to {
        min-width: 50px;
    }
}

@keyframes shrink {
    from {
        min-width: 380px;
    }

    to {
        min-width: 50px;
    }
}

@-webkit-keyframes fadeout {
    from {
        bottom: 30px;
        opacity: 1;
    }

    to {
        bottom: 60px;
        opacity: 0;
    }
}

@keyframes fadeout {
    from {
        bottom: 30px;
        opacity: 1;
    }

    to {
        bottom: 60px;
        opacity: 0;
    }
}
</style>