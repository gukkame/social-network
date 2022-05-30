<template>
  <div class="header">
    <div class="header-item-logo">
      <img alt="kood/ logo" class="logo" src="../assets/images/logo.svg" />
    </div>

    <div class="header-item-nav">
      <div v-if="LoggedIn">
        <span class="nav-items-0">k<span style="color: #FF9D5A">/{{ data }}</span></span>
      </div>
      <div v-else>
        <button class="nav-button-1">
          <RouterLink class="nav-items-1 stretched-link" to="/login">Log In</RouterLink>
        </button>
        <button class="nav-button-2">
          <RouterLink class="nav-items-2 stretched-link" to="/signup">Sign Up</RouterLink>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Header',
  props: {
    LoggedIn: {
      type: Boolean,
      required: true
    },
    data: {
      type: String,
      required: true
    }
  },
}
</script>
