<template>
  <div id="GetCertificate-page">
    <v-card width="400" class="mx-auto mt-5" color="white">
      <v-card-title class="justify-center">
        <h1 class="display-1 ">Apply for certificate</h1>
      </v-card-title>
      <v-card-text>
        <v-form class="mx-auto ml-5 mr-5">
          <v-text-field
              label="Username/Email"
              v-model="username"
              prepend-icon="mdi-account-circle"/>
          <v-text-field
              :type="showPassword ? 'text' : 'password'"
              label="Password"
              v-model="password2"
              prepend-icon="mdi-lock"
              :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
              @click:append="showPassword = !showPassword"/>
          <v-text-field
              :type="showPassword ? 'text' : 'password'"
              label="Confirm Password"
              v-model="password"
              prepend-icon="mdi-lock"
              :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
              @click:append="showPassword = !showPassword"/>
          <v-text-field
              label="First name"
              v-model="firstName"
              prepend-icon="mdi-name-circle"/>
          <v-text-field
              label="Last name"
              v-model="lastName"
              prepend-icon="mdi-address-circle"/>
          <v-text-field
              label="Address"
              v-model="address"
              prepend-icon="mdi-address-circle"/>
          <v-text-field
              label="Phone number"
              v-model="phoneNumber"
              prepend-icon="mdi-address-circle"/>
          <label>Certificate type:</label>
          <v-radio-group>
            <v-radio color="black" @click="ChooseCertificateType(0)" label="Self-issued(root) Certificate"/>
            <v-radio color="black" @click="ChooseCertificateType(1)" label="Intermediate Certificate Authority" />
            <v-radio color="black" @click="ChooseCertificateType(2)" label="End-entity Certificate" />
          </v-radio-group>
          <div v-if="CertificateType !== 0 && CertificateType !== null">
            <v-combobox label="Choose Certificate Authority" items="null"/>
          </div>
        </v-form>
      </v-card-text>
      <v-card-actions class="justify-center mb-5">
        <v-btn color="red mb-5" dark v-on:click="register">
          Apply
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
export default {
  name: "CertificateForm",
  data: () => ({
    showPassword: false,
    username: '',
    password: '',
    password2:'',
    phoneNumber:'',
    firstName: '',
    lastName: '',
    address: '',
    users: [],
    CertificateType: null,
  }),
  computed: {
    user() {
      return {'email': this.username, 'password': this.password, 'phoneNumber':this.phoneNumber,'name':this.firstName,'lastName':this.lastName,'address':this.address}
    }
  },
  methods: {
    register() {
      if(!this.ValidateEmail()){
        return;
      }
      if (this.password!==this.password2){
        alert("Passwords don't match !!!");
        this.password='';
        this.password2='';
        return;
      }
      this.$http.post('http://localhost:8081/users/register', this.user)
          .then(resp => {
            console.log(resp.data);
            window.location.href = 'http://localhost:8080/login';
          })
          .catch(er => {
            console.log('Error while registering in');
            console.log(er.response.data);
          })
    },
    ChooseCertificateType(number) {
        this.CertificateType = number;
        console.log("Cerificate Type: " + this.CertificateType);
    },
    ValidateEmail()
    {
      if (/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(this.username))
      {
        return (true)
      }
      alert("You have entered an invalid email address!")
      return (false)
    }
  }
};

</script>

<style scoped>

#GetCertificate-page {
  display: flex;
  flex-direction: row;
  align-content: center;
  height: 100%;
  background: linear-gradient(45deg, #211010 70%,  #f5efe3 30%)
}

</style>