<template>
<div id="login-page" class="align-center">
  <v-card id="card">
    <v-card-title class="justify-center" >Login</v-card-title>
    <v-card-text class="label">Username:</v-card-text>
    <v-text-field type="username" class="textfield"/>
    <v-card-text class="label">Password:</v-card-text>
    <v-text-field type="password" class="textfield"/>
    <v-btn class="accent mb-3">Login</v-btn>
  </v-card>
</div>
</template>

<script>
export default {
name: "Login"
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