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
              v-model="username"/>
          <v-text-field
              type="password"
              label="Password"
              v-model="password2"/>
          <v-text-field
              type="password"
              label="Confirm Password"
              v-model="password"/>
          <v-text-field
              label="Certificate Common Name"
              v-model="commonName"/>
        </v-form>
        <v-form class="mx-auto ml-5 mr-5" v-else>
          <v-text-field
              label="Service name"
              v-model="username"/>
          <v-text-field
              type="password"
              label="Password"
              v-model="password2"/>
          <v-text-field
              type="password"
              label="Confirm Password"
              v-model="password"/>
          <v-text-field
              label="Certificate Common Name"
              v-model="commonName"/>
        </v-form>
        <v-form class="mx-auto ml-5 mr-5">
          <label>Certificate type:</label>
          <v-radio-group>
            <v-radio color="black" @click="ChooseCertificateType('Root')" label="Self-issued(root) Certificate"/>
            <v-radio color="black" @click="ChooseCertificateType('Intermediary')" label="Intermediate Certificate" />
            <v-radio color="black" @click="ChooseCertificateType('EndEntity')" label="End-entity Certificate" />
          </v-radio-group>
          <div v-if="CertificateType !== 'Root'">
            <v-layout justify-start align-baseline v-if="selectedCertificate !== null">
              <h4>Selected CA</h4>
              <h3>: {{this.selectedCertificate.Cert.EmailAddresses[0]}}</h3>
              <h3>, {{this.selectedCertificate.Type}}</h3>
              <h3>, {{new Date(this.selectedCertificate.Cert.NotAfter).toLocaleString('sr')}}</h3>
              <h3>, {{this.selectedCertificate.Cert.Subject.CommonName}}</h3>
            </v-layout>
            <v-data-table label="Choose Certificate Authority" :items="CACertificates" :headers="headers" >
              <template slot="item" slot-scope="data">
                <td><h3>{{data.item.Type}}</h3></td>
                <td>{{data.item.Cert.Subject.CommonName}}</td>
                <td>{{data.item.Cert.EmailAddresses[0]}}</td>
                <td>{{new Date(data.item.Cert.NotAfter).toLocaleString('sr')}}</td>
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
    isCA: false,
    issuerSerial: '',
    commonName: '',
    issuer: {
      username: "",
      password: "",
      commonName: "",
    },
    headers: [
      { text: 'Certificate Type', value: 'CertificateType', align: 'center',},
      { text: 'Common Name', value: 'EmailName', align: 'center', },
      { text: 'Email/Name', value: 'EmailName', align: 'center', },
      { text: 'Date Of Expire', value: 'DateOfExpire', align: 'center', },
      { text: 'Select CA', value: 'SelectCA', align: 'center', },
    ]
  }),
  computed: {
    user() {
      return {'username': this.username, 'password': this.password, 'parentCommonName': this.parentCommonName}
    },
    certDTO() {
      return {
              'type': this.CertificateType,
              'issuer': this.Issuer,
              'isCA': this.isCA,
              'commonName': this.commonName,
              'emailAddress': this.username,
              'password': this.password}
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
      if (this.CertificateType !== 'Root' && !this.selectedCertificate) {
        alert('You must select Certificate Authority')
        return;
      }
      if (!this.commonName) {
        alert('You must enter Certificate Common Name')
        return;
      }

      if (this.CertificateType !== 'Root'){
        this.Issuer = {
          username: this.selectedCertificate.Cert.EmailAddresses[0],
          password: "",
          commonName: "",
        }
      }
      this.$http.post('http://localhost:8081/api/uos/add', this.user)
          // eslint-disable-next-line no-unused-vars
          .then(resp => {
            if (this.CertificateType === 'Root') {
              this.$http.post('http://localhost:8081/api/cert/root', this.certDTO)
                  // eslint-disable-next-line no-unused-vars
                  .then(resp2 => {
                    this.$router.push('/')
                  }).catch(err => {
                    console.log(err.response)
                  })
            } else if (this.CertificateType === 'Intermediary') {
              this.$http.post('http://localhost:8081/api/cert/intermediary', this.certDTO)
                  // eslint-disable-next-line no-unused-vars
                  .then(resp2 => {
                    this.$router.push('/')
                  })
            }
            else if (this.CertificateType === 'EndEntity') {
              this.$http.post('http://localhost:8081/api/cert/end-entity', this.certDTO)
                  // eslint-disable-next-line no-unused-vars
                  .then(resp2 => {
                    this.$router.push('/')
                  })
            }

          })
          .catch(er => {
            console.log('Error while registering in');
            console.log(er.response.data);
          })
    },
    ChooseCertificateType(number) {
        this.CertificateType = number;
        this.isCA = number !== 1;
        if (number === 10) {
          this.issuerSerial = null
        }
        else {
          this.issuerSerial = this.selectedCertificate.Cert.Issuer.SerialNumber;
        }
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
      this.parentCommonName = this.selectedCertificate.Cert.Subject.CommonName
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