<template>
  <div id="User-page">
    <v-layout justify-center align-baseline>
      <span class="bg"/>
      <v-card width="70%" elevation="10" justify-center class="text-center mt-3">
        <v-card-title style="margin: 10px" >
          <v-layout justify-center style="margin: 5px">
            <h1>User/Service Profile</h1>
          </v-layout>
        </v-card-title>
        <h2>Info</h2>
        <v-layout align-baseline justify-center>
          <h3><b>User/Service name</b></h3>
          <h3>: {{email}}</h3>
        </v-layout>
        <v-layout align-baseline justify-center class="mb-5">
          <h3><b>User/Service password</b></h3>
          <h3>: {{password}}</h3>
        </v-layout>
        <h2> Certificate info</h2>
        <v-layout justify-center>
          <v-data-table v-bind:items="certificates" v-bind:headers="headers">
            <template v-slot:item="row">
              <tr>
                <td>{{row.item.Cert.SerialNumber}}</td>
                <td>{{row.item.Type}}</td>
                <td>{{row.item.Cert.IsCA}}</td>
                <td>{{new Date(row.item.Cert.NotBefore).toLocaleString('sr')}}</td>
                <td>{{new Date(row.item.Cert.NotAfter).toLocaleString('sr')}}</td>
                <td>{{row.item.Cert.Issuer.SerialNumber}}</td>
                <td>
                  <v-checkbox readonly color="accent" v-model="row.item.isValid"/>
                </td>
                <td>
                  <v-btn dark class="accent primary--text" @click="disableCertificate(row.item)">Disable</v-btn>
                </td>
              </tr>
            </template>
          </v-data-table>
        </v-layout>
        <v-layout align-baseline justify-center class="mb-5">
          <v-btn dark class="accent" @click="downloadCert(this.certificates[0].Cert.SerialNumber)">Download</v-btn>
        </v-layout>
      </v-card>
    </v-layout>
  </div>
</template>

<script>
export default {
name: "UserPage",
  data: () => ({
    certificates: [],
    uos: null,
    name: '',
    email: '',
    password: '',
    headers: [
      { text: 'Serial Number', value: 'SerialNumber'},
      { text: 'Certificate Type', value: 'CertificateType' },
      { text: 'Is CA', value: 'IsCA' },
      { text: 'Not Before', value: 'NotBefore' },
      { text: 'Not After', value: 'NotAfter' },
      { text: 'Issued By', value: 'Issuer' },
      { text: 'Is Valid', value: 'IsValid' },
      { text: 'Disable Certificate', value: 'DisableCertificate' },
    ],
  }),
  mounted() {
    this.init()
    //this.getCertificates()
  },
  methods: {
    init() {
      this.email = this.$route.params.username;
      console.log("email: " + this.email)
      this.$http.get('http://localhost:8081/api/uos/' + this.email)
        .then(resp => {
          console.log(resp.data);
          this.uos = resp.data;
          this.password = this.uos.password
          this.axios.get('http://localhost:8081/api/cert')
              // eslint-disable-next-line no-unused-vars
          .then(resp => {
            console.log(resp.data);
            for (let i = 0; i < resp.data.length; ++i) {
              if (resp.data[i].Cert.EmailAddresses[0] === this.email) {
                this.certificates.push(resp.data[i])
                break
              }
            }
          })
      }).catch(err => {
        console.log("Ne radi");
        console.log(err.response.data);
      })
    },
    getCertificates() {
      this.axios.get("http://localhost:8081/api/cert")
          .then(resp => {
            this.certificates = resp.data
            console.log("Certificates:")
            console.log(this.certificates)
          })
          .catch(err => {
            console.log(err.response)
          })
    },
    disableCertificate(cert) {
      console.log(cert)
    },
    downloadCert(serialNumber) {
      this.$router.push("http:/localhost:8081/api/cert/download/" + serialNumber)
    }
  }
}
</script>

<style scoped>

#User-page {
  display: flex;
  flex-direction: row;
  align-content: center;
  height: 100%;
  background: linear-gradient(45deg, #211010 70%,  #f5efe3 30%)
}

</style>