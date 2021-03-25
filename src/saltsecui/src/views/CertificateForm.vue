<template>
  <div id="GetCertificate-page">
    <v-card width="800" class="mx-auto mt-5 mb-5" color="white">
      <v-card-title class="justify-center">
        <h1 class="display-1 ">Apply for certificate</h1>
      </v-card-title>
      <v-card-text>
        <v-form class="mx-auto ml-5 mr-5">
          <v-radio-group>
            <v-radio color="black" @click="ChooseUserOrService(0)" label="User"/>
            <v-radio color="black" @click="ChooseUserOrService(1)" label="Service" />
          </v-radio-group>
        </v-form>
        <v-form class="mx-auto ml-5 mr-5" v-if="isUser">
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
        </v-form>
        <v-form class="mx-auto ml-5 mr-5" v-else>
          <v-text-field
              label="Service name"
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
        </v-form>
        <v-form class="mx-auto ml-5 mr-5">
          <label>Certificate type:</label>
          <v-radio-group>
            <v-radio color="black" @click="ChooseCertificateType(0)" label="Self-issued(root) Certificate"/>
            <v-radio color="black" @click="ChooseCertificateType(1)" label="Intermediate Certificate Authority" />
            <v-radio color="black" @click="ChooseCertificateType(2)" label="End-entity Certificate" />
          </v-radio-group>
          <div v-if="CertificateType !== 0 && CertificateType !== null">
            <v-layout justify-start align-baseline v-if="selectedCertificate !== null">
              <h4>Selected CA</h4>
              <h3>: {{this.selectedCertificate.Cert.EmailAddresses[0]}}</h3>
              <h3>, {{this.selectedCertificate.Type}}</h3>
              <h3>, {{new Date(this.selectedCertificate.Cert.NotAfter).toLocaleString('sr')}}</h3>
              <h3>, {{this.selectedCertificate.Cert.Subject.SerialNumber}}</h3>
            </v-layout>
            <v-data-table label="Choose Certificate Authority" :items="CACertificates" :headers="headers" >
              <template slot="item" slot-scope="data">
                <td><h3>{{data.item.Type}}</h3></td>
                <td>{{data.item.Cert.EmailAddresses[0]}}</td>
                <td><v-btn class="accent" @click="selectCA(data.item)">Select</v-btn></td>
              </template>
            </v-data-table>
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
    CACertificates: [],
    CertificateType: null,
    selectedCertificate: null,
    isUser: true,
    headers: [
      { text: 'Certificate Type', value: 'CertificateType', align: 'center',},
      { text: 'Email/Name', value: 'EmailName', align: 'center', },
      { text: 'Select CA', value: 'SelectCA', align: 'center', },
    ]
  }),
  computed: {
    user() {
      return {'email': this.username, 'password': this.password, 'phoneNumber':this.phoneNumber,'name':this.firstName,'lastName':this.lastName,'address':this.address}
    }
  },
  methods: {
    register() {
      if (this.username === '') {
        alert('You must enter a service name or user email')
        return;
      }
      if (this.password!==this.password2){
        alert("Passwords don't match !!!");
        this.password='';
        this.password2='';
        return;
      }
      this.$http.post('http://localhost:8081/api/', this.user)
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
    ChooseUserOrService(number) {
      this.isUser = number === 0;
    },
    getCACertificates() {
      this.axios.get('http://localhost:8081/api/cert')
      .then(resp => {
        for(let i = 0; i < resp.data.length; ++i) {
          if (resp.data[i].Cert.IsCA) this.CACertificates.push(resp.data[i])
        }
      })
      console.log(this.CACertificates)
    },
    selectCA(certificate) {
      this.selectedCertificate = certificate
      console.log(this.selectedCertificate)
    },
  },
  mounted() {
    this.getCACertificates()
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