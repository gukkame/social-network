<template>
  <div class="modal fade" id="createGroup" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
    aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog  modal-dialog-centered modal-dialog-scrollable modal-xl" style="width:500px">
      <div class="modal-content">
        <div class="modal-body">
          <div class="createGroupModal">
            <Form @submit="createGroup" v-slot="{ errors }" :validation-schema="errorSchema">
              <div class="container">
                <div class="user-details row">
                  <div class="d-flex flex-row">
                    <div class="modal-header-title col">Create a group</div>
                  </div>
                  <div class="row bubble" style="margin: 15px 0px 20px 25px"></div>
                  <div class="col d-flex align-items-center justify-content-center">
                    <input type="file" name="img" id="img" accept="image/x-png,image/jpeg"
                      @change="handleFileUpload($event)" style="display: none" />
                    <label class="
                  uploadProf
                  d-flex
                  justify-content-center
                  align-items-center
                " for="img">Upload image</label>
                    <label v-if="file != null" @click="deleteProf" class="
                  deleteProf
                  d-flex
                  justify-content-center
                  align-items-center
                "><i class="bi bi-trash"></i></label>
                  </div>
                  <span class="formErrors">{{ giferror }}</span>
                  <div class="modal-sec1 col">
                    <div class="modalTitles">
                      Name
                      <i class="formstar">*</i>
                    </div>
                    <Field class="modal-title-input" :class="{ modalinputerror: errors.title }" as="input" name="title"
                      placeholder="60 characters" />
                    <br />
                    <span class="formErrors">{{ errors.title }}</span>
                  </div>
                  <div class="modal-sec2 col">
                    <div class="modalTitles">
                      Description
                      <i class="formstar">*</i>
                    </div>
                    <Field class="modal-description-input" :class="{ modalinputerror: errors.description }"
                      as="textarea" name="description" placeholder="2000 characters" style="resize: none;" />
                    <br />
                    <span class="formErrors">{{ errors.description }}</span>
                  </div>
                </div>
                <span class="formErrors">{{ errormsg }}</span>
              </div>
              <div class="d-flex">
                <button type="button" data-bs-dismiss="modal" class="modal-delete ms-auto"
                  @click="deleteModalData">Cancel</button>
                <button class="modal-create">Create</button>
              </div>
            </Form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import axios from "axios";
import router from "../../router";
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import "bootstrap/dist/js/bootstrap.js"
import $ from 'jquery'

export default {
  data() {
    const errorSchema = yup.object().shape({
      checkUsername: yup.boolean(),
      title: yup.string()
        .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
        .required("Required")
        .max(30, "Maximum length for a group name is 30 characters"),

      description: yup.string()
        .matches(/^[A-Za-z0-9\s\.,;:!?()"'%\-]+$/, "Only characters are allowed")
        .required("Required")
        .max(1000, "Maximum length for a description is 1000 characters"),
    });
    return {
      file: null,
      img: null,
      errormsg: "",
      giferror: "",
      errorSchema,
    };
  },
  components: {
    Form,
    Field,
  },

  methods: {
    createGroup(values) {
      let currentRouter = this.$route.path
      if (document.cookie.length == 0) {
        return router.go(`"${currentRouter}"`)
      }

      let token = document.cookie
      let correctToken = token.split(":")

      let formData = new FormData();
      if (this.file == null) {
        formData.append("image", "");
      } else {
        formData.append("image", this.file);
      }
      formData.append("title", values.title);
      formData.append("content", values.description);
      formData.append("token", correctToken[0]);

      axios.post("http://localhost:8080/creategroup", formData)
      
        .then((res) => {
          if (res.data.message === "Malicious user detected") {
            return this.errormsg = res.data.message
          }

          if (res.data.message === "User not authenticated") {
            $('body').removeClass('modal-open');
            /* $('#staticBackdrop').hide() */
            $('.modal-backdrop').hide()
            async function removeAllAttrs(element) {
              for (var i = element.attributes.length; i-- > 0;)
                await element.removeAttributeNode(element.attributes[i]);
            }
            removeAllAttrs(document.body);
            $('body').css('overflow', 'auto');
            return router.push("/")
          }
          this.errormsg = ""
          $('body').removeClass('modal-open');
          /* $('#staticBackdrop').hide() */
          $('.modal-backdrop').hide()
          async function removeAllAttrs(element) {
            for (var i = element.attributes.length; i-- > 0;)
              await element.removeAttributeNode(element.attributes[i]);
          }
          removeAllAttrs(document.body);
          $('body').css('overflow', 'auto');
          router.push(`/groups/${res.data}`)

        })
        .catch((error) => { });
    },

    deleteModalData() {
      $('#createGroup form')[0].reset();
    },
    // Handle file upload!
    handleFileUpload(event) {
      this.file = event.target.files[0];
      if (this.file.type == "image/gif") {
        this.file = null
        let bubble = this.$el.querySelector(".bubble")
        bubble.style.backgroundImage = `none`
        this.giferror = "GIFs are not allowed"
        return
      }
      this.giferror = ""
      this.img = URL.createObjectURL(this.file);
      let bubble = this.$el.querySelector(".bubble")
      bubble.style.backgroundImage = `url('${this.img}')`
    },

    deleteProf() {
      this.img = null;
      this.file = null;
      let bubble = this.$el.querySelector(".bubble")
      bubble.style.backgroundImage = `none`
    },
  },
};
</script>