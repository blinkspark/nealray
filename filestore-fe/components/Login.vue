<template>
  <div class="modal" ref="modal">
    <div class="modal-background" @click="toggleModal"></div>
    <div class="modal-content">
      <div class="panel is-primary">
        <div class="panel-heading">
          <p class="panel-title">Login</p>
        </div>
        <div class="panel-block has-background-dark">
          <form @submit.prevent="login">
            <div class="field">
              <label class="label has-text-light">Email</label>
              <div class="control">
                <input class="input" type="email" placeholder="Email" v-model="email">
              </div>
            </div>
            <div class="field">
              <label class="label has-text-light">Password</label>
              <div class="control">
                <input class="input" type="password" placeholder="Password" v-model="password">
              </div>
            </div>
            <div class="field">
              <div class="control">
                <button class="button is-primary">Login</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
    <button class="modal-close" aria-label="close" @click="toggleModal"></button>
  </div>
</template>

<script>
export default {
  // name: 'Login',
  data () {
    return {
      email: '',
      password: ''
    }
  },
  methods: {
    login () {
      this.$axios.get('/api/test').then(res => {
        console.log(res)
      })
      this.toggleModal()
    },
    toggleModal () {
      this.$refs.modal.classList.toggle('is-active')
    }
  }
}
</script>

<style>
form {
  width: 100%;
  margin: 1rem 1rem;
}
</style>