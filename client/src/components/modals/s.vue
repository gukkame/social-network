<template>
  <main>
    <Form @submit="newGroup">
 
        <div>Title</div>
        <Field class="form-input" as="input" name="title" />
        <input name="title" v-model="title" type="text" />

        <!-- <span class="formErrors">{{ errors.title }}</span> -->
        <div>Description</div>
        <Field class="form-input" as="input" name="content" />
        <input name="content" v-model="content" type="text" />
        <br />
        <!-- <span class="formErrors">{{ errors.content }}</span> -->

        <button class="">Create group!</button>
  

      <h2>Single File</h2>
      <hr />
      <!-- <label
        >File
        <input
          type="file"
          name="img"
          accept="image/x-png,image/gif,image/jpeg"
          @change="handleFileUpload($event)"
        />
      </label> -->
      <br
    /></Form>
  </main>
</template>

<script>
export default {
  data() {
    return {
      errorMessage: "",
      allGroups: [],
      title: "",
      content: "",
    };
  },
  methods: {
    handleFileUpload(event) {
      this.file = event.target.files[0];
    },
    newGroup(values) {
      let token = document.cookie;
      if (token.length === 0) {
        // return router.push("/");
      }
      let correctToken = token.split(":");

      //formdata title, content, img
      let formData = new FormData();
      formData.append("title", values.title);
      formData.append("content", values.content);
      formData.append("token", correctToken);
      formData.append("image", this.file);
      console.log(values);
      axios
        .post("http://localhost:8080/creategroup", formData)
        .then((res) => {
          if (res.data.message === "Malicious user detected") {
            return (this.errormsg = res.data.message);
          }
          if (res.data.message === "Image is too big!") {
            return (this.errormsg = res.data.message);
          }
          console.log(res.data.message);
          this.errormsg = "";
          // router.push("/group_page");
        })
        .catch((error) => {
          console.log("errorrrrs");
        });
    },
  },
};
</script>