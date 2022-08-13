<script setup>
import LeaveGroupModal from "../components/modals/LeaveGroup.vue"
import CreateGroupInviteModal from "../components/modals/CreateInvite.vue"

</script>

<template>
    <div class="profile container col">
        <div class="d-flex justify-content-center">
            <div class="oneGroup d-flex flex-row flex-wrap justify-content-space">
                <div class="singleGroupCol col-5 d-flex flex-row">
                    <div class="d-flex flex-column col">
                        <div class="d-flex justify-content-center align-items-center" style="height: 180px;">
                            <div v-if="data.Image != ``" class="row bubble" style="margin: 25px 0px 0px 0px"></div>
                            <img v-else class="profileimage col" src="../assets/images/groups.svg" />
                        </div>
                        <div class="d-flex justify-content-center align-items-center" style="height: 60px;">
                            <div style="font-size: 18px;">{{ data.Title }}</div>
                        </div>
                    </div>
                </div>
                <div class="singleGroupCol col d-flex flex-row" style="position: relative">
                    <div class="d-flex flex-column col">
                        <div class="d-flex flex-wrap col"
                            style="height: 180px; padding: 35px 10px 10px 10px ; word-break:break-all;">
                            {{ data.Content }}
                        </div>
                        <div class="d-flex justify-content-start align-items-center flex-row" style="height: 60px;">
                            <div class="col" style="font-size: 18px;">
                                <RouterLink v-if="data.Creator_name"
                                    :to="{ name: 'profile', params: { id: data.Creator_name } }"
                                    style="text-decoration: none; width: 100%;" class="href">
                                    <div class="followerLink d-flex" style="font-size: 16px; color:#2E343D ;">
                                        <i class="bi bi-diamond" style="margin-right:5px"></i>{{ data.Creator_name }}
                                    </div>
                                </RouterLink>
                            </div>
                            <div class="col" style="margin-bottom: 5px">
                                <i style="font-size: 23px;" class="bi bi-person-check"></i> {{ data.Member_count }}
                            </div>
                        </div>
                    </div>
                    <div v-if="data.User_status == `Member`" class="d-flex flex-column"
                        style="position: absolute; top: 1px; right: 1px;">
                        <button data-bs-toggle="modal" data-bs-target="#leaveGroup" class="unfollowProfile col"
                            style="margin-right: 6px">Leave</button>
                    </div>
                    <div v-if="data.User_status == `Member`|| data.User_status == `Owner`" class="d-flex flex-column"
                        style="position: absolute; bottom: 1px; right: 1px;">
                        <button data-bs-toggle="modal" data-bs-target="#inviteGroup" class="followProfile col"
                        style="margin-right: 6px">Invite</button>
                    </div>
                </div>
            </div>
        </div>
        <LeaveGroupModal />
        <CreateGroupInviteModal/>
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

    mounted() {
        let bubble = this.$el.querySelector(".bubble")
        if (bubble == null) {
            return
        }
        bubble.style.backgroundImage = `url('http://localhost:8080${this.data.Image}')`
    },
}



</script>