<template>
<div id="login-page" class="align-center">
  <v-card id="card">
    <v-card-title class="justify-center" >Login</v-card-title>
    <v-card-text class="label">Email:</v-card-text>
    <v-text-field type="username" class="textfield" v-model="email"/>
    <v-card-text class="label">Password:</v-card-text>
    <v-text-field type="password" class="textfield" v-model="password"/>
    <v-btn class="accent mb-3" @click="login">Login</v-btn>
  </v-card>
</div>
</template>

<script>
export default {
  name: "Login",
  data: () => ({
    email: '',
    password: ''
  }),
  methods: {
    login() {
      console.log(this.email + ' --- ' + this.password)
      if (this.email === 'admin@email.com') {
        if (this.password === 'admin1') {
          this.$router.push('/admin')
        }
        else {
          alert('Incorrect login attempt. Try Again.')
        }
      }
      else {
        alert('Incorrect login attempt. Try Again.')
      }
    }
  }
}
</script>


<style scoped>

#login-page {
  display: flex;
  flex-direction: column;
  align-content: center;
  height: 100%;
  background: linear-gradient(45deg, #211010 70%,  #f5efe3 30%)
}

#card {
  margin: 10px;
  width: 50%;
}

.textfield {
  margin-left: 20%;
  margin-right: 20%;
  margin-top: 10px;
  font-size: 20pt;
}

.label {
  font-size: 20pt;
}

.title {
}



</style>