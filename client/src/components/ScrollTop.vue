<template>
    <div ref="scrollTopButton" class="invisible scrollbuttondiv">
        <div class>
            <i @click="scrollToTop" class="scrollbutton bi bi-arrow-up-circle"></i>
        </div>
    </div>
</template>


<script>
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"

export default {
    mounted() {
        window.addEventListener("scroll", this.handleScroll);
    },

    beforeUnmount() {
        window.removeEventListener("scroll", this.handleScroll);
    },

    methods: {
        handleScroll() {
            const scrollBtn = this.$refs.scrollTopButton;

            if (window.scrollY > 200) {
                scrollBtn.classList.remove("invisible");
            } else {
                scrollBtn.classList.add("invisible");
            }
        },
        scrollToTop() {
            window.scrollTo({ top: 0, behavior: "smooth" });
        },
    },
};
</script>


<style>
@import "@/assets/css/base.css";

.scrollbuttondiv {
    position: fixed;
    top: 50px;
    right: 25px;
    z-index: 1;
    transition: 0.5s;
    visibility: visible;
    opacity: 1;
    transition: opacity 0.3s linear;
}

.scrollbutton {
    color: #2E343D;
    cursor: pointer;
    font-size: 50px;
    transition: 0.2s;
}

.scrollbutton:hover {
    color: #FF9D5A;
    cursor: pointer;
    transition: 0.2s;
}

.invisible {
    visibility: hidden;
    opacity: 0;
    transition: visibility 0s 2s, opacity 0.3s linear;
}
</style>