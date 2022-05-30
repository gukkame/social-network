<template>
    <div class="container errorpage d-flex justify-content-center">
        <div class="d-flex justify-content-center">
            <div class="erroricon d-flex justify-content-center row">
                <img class="errorimage col" src="../assets/images/nopage.svg" />
                <div class="d-flex justify-content-center errormsg">Oops! This page does not exist!</div>
                <div class="errorbox">
                    <a class="errorsneak col d-flex justify-content-center" @click="$router.go(-1)">Sneak back..</a>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=Passero+One&display=swap");
</style>