<template>
    <div class="chatwindowPerson d-flex">
        <div>
            <img class="chatwindowImg col" src="../../assets/images/profile.svg" />
        </div>
        <div class="col posthDetails">
            <div class="col chatwindowUser">robert1234</div>
            <div class="col chatwindowTime">24 minutes ago</div>
        </div>
        <div class="col d-flex justify-content-center align-items-center">
            <div class="offlineDot"></div>
        </div>
    </div>
</template>

<style>
.offlineDot {
    border-radius: 50%;
    width: 15px;
    height: 15px;
    background-color: rgb(253, 252, 252);
    border: solid 1px rgb(219, 219, 219);
}
</style>