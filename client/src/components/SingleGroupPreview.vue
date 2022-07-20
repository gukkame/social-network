<template>
    <div class="singleGroup">
        <RouterLink class="d-flex flex-row flex-wrap justify-content-space"
            :to="{ name: 'group', params: { id: data.GroupID } }">
            <div class="singleGroupCol col-5 d-flex flex-row">
                <div class="d-flex flex-column col">
                    <div class="d-flex justify-content-center align-items-center" style="height: 180px;">
                        <img v-if="data.Image == ``" class="profileimage col" src="../assets/images/groups.svg" />
                         <div v-else class="row bubble" style="margin: 25px 0px 0px 0px"></div>
                    </div>
                    <div class="d-flex justify-content-center align-items-center" style="height: 60px;">
                        <div style="font-size: 18px;">{{ data.Title }}</div>
                    </div>
                </div>
            </div>
            <div class="singleGroupCol col d-flex flex-row">
                <div class="d-flex flex-column col">
                    <div class="d-flex flex-wrap col" style="height: 180px; padding: 10px; word-break:break-all;">
                        {{ data.Content }}
                    </div>
                    <div class="d-flex justify-content-start align-items-center flex-row" style="height: 60px;">
                        <div class="col" style="font-size: 18px;">
                            <RouterLink :to="{ name: 'profile', params: { id: data.Creator_name } }"
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
            </div>
        </RouterLink>
    </div>
</template>

<script >
export default {
    props: {
        data: {
            type: Object,
            required: true,
        }
    },

    mounted() {
        console.log(this.data)
        let bubble = this.$el.querySelector(".bubble")
        if (bubble == null) {
            return
        }
        bubble.style.backgroundImage = `url('http://localhost:8080${this.data.Image}')`
    },
}
</script>