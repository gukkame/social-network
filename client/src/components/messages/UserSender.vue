<template>
  <div v-if="checkSender()" class="d-flex flex-row senderMessage justify-content-end">
    <div class="d-flex flex-column" style="overflow: hidden">
      <div class="senderMSG d-flex justify-content-end text-break col" @mouseover="hover = true"
        @mouseleave="hover = false">
        {{ data.Message }}
      </div>
      <div v-if="hover" class="dateTime col">
        {{ datetime }}
      </div>
    </div>
  </div>
  <div v-else class="d-flex flex-row message">
    <div v-if="data.IsGroup == 0">
      <img v-if="data1 == `` && data.IsGroup == 0" class="receiverwindowImg" src="../../assets/images/profile.svg" />
      <div v-else-if="data.IsGroup == 0" v-bind:id="data.Sender + [1]" class="receiverwindowImg"></div>
    </div>

    <div v-if="data.IsGroup == 1" class="d-flex flex-column" style="overflow: hidden; margin-left: 7%;">
      <div v-if="data.IsGroup == 1" class="msger-header">&nbsp; {{ data.Sender }}</div>
      <div class="receiverMSG" @mouseover="hover = true" @mouseleave="hover = false">
        {{ data.Message }}
      </div>
      <div v-if="hover" class="dateTime col">{{ datetime }}<br /></div>
    </div>

    <div v-else class="d-flex flex-column" style="overflow: hidden;">
      <div v-if="data.IsGroup == 1" class="msger-header">&nbsp; {{ data.Sender }}</div>
      <div class="receiverMSG" @mouseover="hover = true" @mouseleave="hover = false">
        {{ data.Message }}
      </div>
      <div v-if="hover" class="dateTime col">{{ datetime }}<br /></div>
    </div>
    
  </div>
</template>

<script>

export default {
  data() {
    return {
      hover: false,
    };
  },

  props: {
    data: {
      type: Object,
      required: true,
    },
    data1: {
      type: String,
      required: true,
    },
  },
  mounted: function () {
    const bubble1 = document.querySelectorAll(".receiverwindowImg")
    bubble1.forEach((bubble) => {
      if (bubble != null) {
        this.img = "url(http://localhost:8080" + this.data1 + ")"
        bubble.style.backgroundImage = this.img
      }
    });

  },

  methods: {
    checkSender() {
      if (this.data == undefined) {
        return false;
      }

      if (this.data.Sender == undefined) {
        return false;
      }
      let token = document.cookie.split(":");
      if (this.data.Sender == token[1]) {
        return true;
      }
      return false;
    },
  },

  computed: {
    datetime() {
      var d = new Date(this.data.Datetime);
      d = d.toLocaleString();
      return d;
    },
  },
};
</script>

<style>
.msger-header {
  display: flex;
  justify-content: space-between;
  font-size: 75%;
  border-bottom: var(--border);
  color: #666;
}

.dateTime {
  font-size: 10px;
}

.message {
  max-width: 180px;
  margin-top: 10px;
  margin-bottom: 10px;
}

.receiverMSG {
  padding: 8px;
  font-size: 13px;
  overflow-y: auto;
  flex-wrap: wrap;
  height: fit-content;
  width: auto;
  max-width: 180px;
  background-color: rgb(219, 219, 219);
  border-radius: 5px;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  -o-user-select: none;
  user-select: none;
}

.receiverwindowImg {
  border-radius: 50%;
  margin-top: 7px;
  margin-left: 10px;
  margin-right: 5px;
  width: 25px;
  height: 25px;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  -o-user-select: none;
  user-select: none;
  border: solid 1px rgb(219, 219, 219);
  display: flex;
  justify-content: center;
  align-items: center;
  background-size: cover;
}

.senderMessage {
  width: 100%;
  margin-top: 10px;
  margin-bottom: 10px;
}

.senderMSG {
  margin-right: 10px;
  width: auto;
  max-width: 180px;
  padding: 8px;
  font-size: 13px;
  overflow-y: auto;
  flex-wrap: wrap;
  height: fit-content;
  background-color: #ff9d5a;
  border-radius: 5px;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  -o-user-select: none;
  user-select: none;
}
</style>
