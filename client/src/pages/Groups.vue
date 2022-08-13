<script setup>
import SingleGroupPreviewComp from "../components/SingleGroupPreview.vue";
import CreateGroupComp from "../components/modals/CreateGroup.vue";
</script>

<template>
  <div v-if="LoggedIn" class="catepage container col">
    <div class="catetitle">
      <i class="catetitleicon bi bi-server"></i>My groups ({{
        checkGroups(mygroups)
      }})
    </div>
    <div class="d-flex justify-content-center">
      <div
        v-if="checkGroups(mygroups) == 0"
        class="groupsHeader d-flex justify-content-center col align-items-center overflow-auto"
      >
        <div class="noAvailable d-flex justify-content-center">
          You are not part of any group!
        </div>
      </div>
      <div
        v-else
        class="groupsHeader d-flex col align-items-center overflow-auto"
      >
        <div v-for="data in mygroups">
          <SingleGroupPreviewComp :data="data" />
        </div>
      </div>
    </div>
  </div>
  <div class="groupPage container col">
    <div class="d-flex">
      <div class="catetitle col">
        <i class="catetitleicon bi bi-share"></i>Available groups ({{
          checkGroups(groups)
        }})
      </div>
      <div v-if="LoggedIn">
        <button
          data-bs-toggle="modal"
          data-bs-target="#createGroup"
          class="createGroup col"
        >
          Create a group
        </button>
      </div>
    </div>
    <div class="d-flex">
      <div class="mainheader content col">
        <div
          class="noAvailable d-flex justify-content-center"
          v-if="checkGroups(groups) == 0"
        >
          No groups available
        </div>
        <div v-else>
          <div
            class="d-flex flex-wrap contentinside justify-content-center col"
          >
            <div v-for="data in groups">
              <SingleGroupPreviewComp :data="data" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <CreateGroupComp />
</template>

<script>
import axios from "axios";
import { delay } from "../common-js/time.js";
export default {
  data() {
    return {
      mygroups: [],
      groups: [],
    };
  },

  props: {
    LoggedIn: {
      type: Boolean,
      required: true,
    },
  },

  created() {
    this.$watch(
      () => this.$route.path,
      () => {
        this.fetchGroups();
      },
      { immediate: true }
    );
  },

  methods: {
    allgroups() {
      this.mygroups = this.result.filter(
        (ismember) => ismember.Group_member == 1
      );
      this.groups = this.result.filter(
        (ismember) => ismember.Group_member == 0
      );
    },
    async fetchGroups() {
      let data = {};

      let token = document.cookie;
      let correctToken = token.split(":");

      let config = {
        headers: {
          header1: correctToken[0],
        },
      };

      await delay(200).then(() => {
        axios
          .post("http://localhost:8080/groups", data, config)
          .then((res) => {
            this.result = res.data;
            this.allgroups();
            if (res.data.message === "Post request failed") {
              return router.go(-1);
            }
          })
          .catch((error) => {});
      });
    },

    checkGroups(arg) {
      if (arg == null) {
        return 0;
      } else {
        return arg.length;
      }
    },
  },
};
</script>

<style>
@import "../assets/css/base.css";
</style>
