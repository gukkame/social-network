<script setup>
import UserSenderComp from "./UserSender.vue";
</script>

<template>
  <div class="messagewindowdiv">
    <div class="messagewindowBox">
      <div class="messagewindowBoxH d-flex">
        <div>
          <img v-if="data.Avatar_image == ``" class="bubble3 col" src="../../assets/images/profile.svg" />
          <div v-else v-bind:id="data.Username + [0]" class="bubble3 col"></div>
        </div>
        <div class="col-7 posthDetails">
          <div v-if="data.Username" class="col chatwindowUser">{{ data.Username }}</div>
          <div v-if="data.Status == 1" class="col chatwindowTime">Active</div>
          <div v-if="data.Status == 0" class="col chatwindowTime">Inactive</div>
        </div>
        <div class="col-2 d-flex justify-content-center align-items-center">
          <div v-if="data.Status == 1" class="activityDot"></div>
          <div v-if="data.Status == 0" class="offlineDot"></div>
        </div>
      </div>
      <div class="messagewindowBoxB" v-on:scroll="loadMoreMsgs" ref="messageBody" id="messageBody">
        <div v-for="(message, key) in messages" :key="key">
          <UserSenderComp :data="message" :data1="data.Avatar_image" />
        </div>

        <div class="typing d-flex flex-row" v-if="typing">
          <div class="dots col-3">
            <div class="typing__dot"></div>
            <div class="typing__dot"></div>
            <div class="typing__dot"></div>
          </div>
          <div class="col-12">
            <div class="typer-font">
              <a class="typer typer-font">{{ data.Username }}</a> typing..
            </div>
          </div>
        </div>
      </div>

      <div class="messagewindowBoxF">
        <Form class="chattextForm d-flex">
          <EmojiPicker :hide-search="true"
            disabled-groups="[animals_nature, food_drink, activities,travel_places,objects,symbols]" mode="append"
            @select="onSelectEmoji" pickerType="textarea" @keyup.enter.exact="sendMessage()" class="chattextBox col-8"
            @input="updateTheVariable($event.target.value)" as="textarea" name="description" placeholder=""
            id="textarea" />
          <button type="button" @click="sendMessage()" class="chattextSend">Send</button>
        </Form>
      </div>
    </div>
  </div>
</template>

<script>
import "../../../node_modules/vue3-emoji-picker/dist/style.css";
import EmojiPicker from "vue3-emoji-picker";
import _, { debounce } from "lodash";
import { Form, Field } from "vee-validate";
import { ws } from "../../common-js/messages.js";
export default {
  data() {
    return {
      img: "",
      description: null,
      messages: [],
      typing: false,
      debounce: null,
      divHeight: 0,
      user1: "",
      user2: "",
    };
  },

  props: {
    data: {
      type: Object,
      required: true,
    },
  },

  mounted: function () {
    let bubble = document.getElementById(this.data.Username + [0])
    if (bubble != null) {
      this.img = "url(http://localhost:8080" + this.data.Avatar_image + ")"
      bubble.style.backgroundImage = this.img
    }
    const boxes = document.querySelectorAll(".v3-emoji-picker-textarea");
    boxes.forEach((box) => {
      box.style.minHeight = "30px";
    });
    document.getElementsByName(
      "description"
    )[0].placeholder = `Message ${this.data.Username}`;

    ws.addEventListener("message", (event) => {
      this.handleData(event);
    });
    this.initialChatHistory();
  },

  methods: {
    onSelectEmoji(emoji) {
      if (this.description == null) {
        this.description = emoji.i;
      } else {
        this.description += emoji.i;
      }
    },
    updateTheVariable(value) {
      this.description = value;
    },

    handleData(event) {
      let incData = JSON.parse(event.data);
      this.user2 = this.data.Username
      let token = document.cookie.split(":");
      this.user1 = token[1]
      if (incData.Type == "privateMSG" && ((incData.Content.Receiver == this.user2 && incData.Content.Sender == this.user1) || (incData.Content.Receiver == this.user1 && incData.Content.Sender == this.user2))) {
        this.newMessage(incData);
      }

      if (incData.Type == "initialMessageHistory") {
        this.handleInitialChatHistory(incData);
      }

      if (incData.Type == "messageHistory") {
        this.oldChatHistory(incData);
      }

      if (incData.Type != "typing") {
        this.typing = false;
      }

      if (incData.Type == "typing") {
        this.typing = true;
      }
    },

    newMessage(object) {
      let msg = object.Content;
      this.messages.push(msg);
      this.$nextTick(function () {
        let element = this.$refs.messageBody;

        if (element) {
          element.scrollTop = element.scrollHeight;
        }
      });
    },

    handleInitialChatHistory(object) {
      let msg = object.Content;
      this.messages = [];
      let y = [...msg].reverse();
      for (let i = 0; i <= y.length - 1; i++) {
        this.messages.push(y[i]);
      }
      this.$nextTick(function () {
        let element = this.$refs.messageBody;

        if (element) {
          element.scrollTop = element.scrollHeight;
        }
      });
    },

    oldChatHistory(object) {
      const element = document.getElementById("messageBody");
      let height = element.scrollHeight;
      let msg = object.Content;
      let msglen = this.messages.length;

      let y = [...msg].reverse();
      if (msglen == y.length) {
        return;
      }
      this.messages = [];
      for (let i = 0; i <= y.length - 1; i++) {
        this.messages.push(y[i]);
      }

      this.$nextTick(function () {
        let elements = this.$refs.messageBody;

        if (elements) {
          elements.scrollTop = elements.scrollHeight - height - 50;
        }
      });
    },

    loadMoreMsgs() {
      const element = document.getElementById("messageBody");
      if (element.scrollTop < 5) {
        clearTimeout(this.debounce);
        this.debounce = setTimeout(() => {
          this.initialChatHistory();
        }, 1200);
      }
      return;
    },

    initialChatHistory() {
      let token = document.cookie.split(":");
      let payload = {
        Type: "messageHistory",
        Content: null,
        User1: token[1],
        User2: this.data.Username,
        MsgCount: this.messages.length,
      };
      ws.send(JSON.stringify(payload));
    },

    sendMessage() {
      if (this.description == null) {
        return;
      }
      let regex = /^\s+$/g;
      if (this.description.match(regex)) {
        return;
      }

      let token = document.cookie.split(":");
      let payload = {
        Type: "privateMSG",
        Content: {
          Message: this.description,
          Sender: token[1],
          Receiver: this.data.Username,
        },
        User1: token[1],
        User2: this.data.Username,
        MsgCount: this.messages.length,
      };

      ws.send(JSON.stringify(payload));
      const textarea = document.getElementById("textarea");
      textarea.value = null;
      textarea.placeholder = null;
    },
  },

  components: {
    Form,
    Field,
  },
};
</script>

<style>
.typing__dot {
  float: left;
  width: 4px;
  height: 4px;
  margin: 0 1.5px;
  background: #8d8c91;
  border-radius: 50%;
  opacity: 0;
  animation: loadingFade 1s infinite;
}

.dots {
  padding-top: 10px;
  margin-right: -11px;
}

.typing__dot:nth-child(1) {
  animation-delay: 0s;
}

.typing__dot:nth-child(2) {
  animation-delay: 0.2s;
}

.typing__dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes loadingFade {
  0% {
    opacity: 0;
  }

  50% {
    opacity: 0.8;
  }

  100% {
    opacity: 0;
  }
}

.typing {
  position: absolute;
  bottom: 72px;
  left: 7px;
  transition: 0.1;
  max-height: 6px;
}

.typer {
  font-weight: 600;
  text-decoration: none;
  color: #2e343d;
}

.typer-font {
  font-size: 13px;
}

.chattextForm {
  width: 100%;
  height: 100%;
}

.chattextBox {
  margin-top: 2px;
  font-size: 13px;
  width: 200px;
  max-height: 40px;
  border-radius: 4px;
  outline: 0;
  border-right: none;
  border-left: none;
  border: 0.1px solid #2e343d;
}

.chattextSend {
  width: 50px;
  border: solid 1px rgb(219, 219, 219);
  margin: 10px 12px 10px 4px;
  border-radius: 4px;
  padding: 2px;
  background-color: rgb(253, 252, 252);
  transition: 0.25s;
  position: relative;
}

.chattextSend:hover {
  transform: scale(1.13);
  background-color: #ff9d5a;
  color: rgb(253, 252, 252);
  border: solid 0px rgb(219, 219, 219);
}

.messagewindowBoxH {
  border-bottom: solid 1px rgb(219, 219, 219);
  box-shadow: 0 2px 6px 0 rgb(218 218 253 / 65%), 0 2px 6px 0 rgb(206 206 238 / 54%);
}

.messagewindowBoxB {
  width: 100%;
  height: 195px;
  /* background-color: red; */
  overflow-wrap: break-word;
  overflow-y: auto;
  flex-direction: column-reverse;
}

.messagewindowBoxF {
  margin-top: 14px;
  padding: 4px;
  border-top: solid 1px rgb(219, 219, 219);
}

.messagewindowClose {
  top: -11px;
  right: -1px;
  position: absolute;
  font-size: 25px;
  cursor: pointer;
  transition: 0.1s;
}

.messagewindowClose:hover {
  color: #ff9d5a;
  transition: 0.1s;
}

.messagewindowdiv {
  position: fixed;
  bottom: 69px;
  right: 280px;
  z-index: 999;
  transition: 0.5s;
  width: 265px;
  height: 300px;
}

.messagewindowBox {
  box-shadow: 0 2px 6px 0 rgb(218 218 253 / 65%), 0 2px 6px 0 rgb(206 206 238 / 54%);
  border-radius: 6px;
  width: 100%;
  height: 100%;
  border: solid 1px rgb(219, 219, 219);
  transition: 0.5s;
  background-color: rgb(253, 252, 252);
}
</style>
