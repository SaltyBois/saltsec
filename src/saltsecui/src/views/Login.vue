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
      if (this.email === 'admin@email.com') {
        if (this.password === 'admin1') {
          this.$router.push('/admin')
        }
        else {
          alert('Incorrect login attempt. Try Again Admin.')
        }
      }
      else {
        this.axios.get('http://localhost:8081/api/uos')
            // eslint-disable-next-line no-unused-vars
        .then(resp => {
          this.$router.push('/admin/')
          let arr = resp.data
          for (let i = 0; i < arr.length; ++i) {
            if (arr[i].username === this.email && arr[i].password === this.password) {
              // this.$router.push('/admin/' + this.email)
            }
          }

        })
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