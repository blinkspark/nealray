<template>
  <div class="panel is-primary">
    <div class="panel-heading">Login</div>
    <div class="panel-block">
      <form @submit.prevent="login">
        <div class="field">
          <label class="label">Email</label>
          <div class="control">
            <input
              class="input"
              type="email"
              v-model="email"
              placeholder="Email"
            />
          </div>
        </div>
        <div class="field">
          <label class="label">Password</label>
          <div class="control">
            <input
              class="input"
              type="password"
              v-model="password"
              placeholder="Password"
            />
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
</template>

<script>
export default {
  data() {
    return {
      email: "",
      password: "",
    };
  },
  methods: {
    login() {
      console.log(this.email, this.password);
    },
  },
};
</script>

<style>
</style>