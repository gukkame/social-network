<template>
    <div class="catepage container col">
        <div class="catetitle"><i class="catetitleicon bi bi-folder"></i>Categories</div>
        <div class="d-flex justify-content-center">
            <div class="cateheader d-flex justify-content-center col flex-wrap">
                <RouterLink class="row text" to="/Go">
                    <div :class="[category == 'go' ? 'active' : '']" class="catecard d-flex justify-content-center row">
                        <img class="col cateiconGO" src="../assets/images/Go.svg" />
                        <div class="d-flex justify-content-center catetextGO catetext col-2">Go</div>
                    </div>
                </RouterLink>

                <RouterLink class="row text" to="/HTML5">
                    <div :class="[category == 'html' ? 'active' : '']"
                        class="catecard d-flex justify-content-center row">
                        <img class="col cateiconHTML" src="../assets/images/HTML.svg" />
                        <div class="d-flex justify-content-center catetextHTML catetext">HTML5</div>
                    </div>
                </RouterLink>


                <RouterLink class="row text" to="/CSS">
                    <div :class="[category == 'css' ? 'active' : '']"
                        class="catecard d-flex justify-content-center row">
                        <img class="col cateiconCSS" src="../assets/images/CSS.svg" />
                        <div class="d-flex justify-content-center catetextCSS catetext">CSS</div>
                    </div>
                </RouterLink>

                <RouterLink class="row text" to="/JavaScript">
                    <div :class="[category == 'javascript' ? 'active' : '']"
                        class="catecard d-flex justify-content-center row">
                        <img class="col cateiconJS" src="../assets/images/Javascript.svg" />
                        <div class="d-flex justify-content-center catetextJS catetext">JavaScript</div>
                    </div>
                </RouterLink>

                <RouterLink class="row text" to="/Vue.js">
                    <div :class="[category == 'vue' ? 'active' : '']"
                        class="catecard d-flex justify-content-center row">
                        <img class="col cateiconVUE" src="../assets/images/Vue.svg" />
                        <div class="d-flex justify-content-center catetextVUE catetext">Vue.js</div>
                    </div>
                </RouterLink>
            </div>
        </div>
    </div>
</template>


<script>
export default {
    props: {
        category: {
            type: String,
            required: true
        },
    },
}
</script>

<style>
.active {
    border: solid 2px #FF9D5A;
}

.text {
    text-decoration: none;
}
</style>