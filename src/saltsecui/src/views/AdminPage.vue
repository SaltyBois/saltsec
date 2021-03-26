<template>
 <div id="Admin-page">
   <v-layout justify-center align-baseline>
     <span class="bg"/>
     <v-card width="70%" elevation="10" justify-center class="text-center mt-3">
       <v-card-title style="margin: 10px" >
         <v-layout justify-center style="margin: 5px">
           <h1>Admin Profile</h1>
         </v-layout>
       </v-card-title>
       <h2>Admin Info</h2>
       <!--   Place for basic Admin information    -->
       <v-layout align-baseline justify-center>
         <h3><b>Admin name: </b></h3>
         <h3>{{name}}</h3>
       </v-layout>
       <v-layout align-baseline justify-center class="mb-5">
         <h3><b>Admin email:</b></h3>
         <h3>{{email}}</h3>
       </v-layout>

       <h2> Certificates info</h2>
       <v-layout justify-center>
         <v-data-table v-bind:items="certificates" v-bind:headers="headers">
           <template v-slot:item="row">
             <tr>
               <td>Very very big number</td>
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
     </v-card>
   </v-layout>
 </div>
</template>

<script>
export default {
  name: "AdminPage",
  data: () => ({
    certificates: [],
    admin: null,
    name: '',
    email: '',
    headers: [
      { text: 'Serial Number', value: 'SerialNumber'},
      { text: 'Certificate Type', value: 'CertificateType' },
      { text: 'Is CA', value: 'IsCA' },
      { text: 'Not Before', value: 'NotBefore' },
      { text: 'Not After', value: 'NotAfter' },
      { text: 'Issued By', value: 'Issuer' },
      { text: 'Is Valid', value: 'IsValid' },
      { text: 'Disable Certifiacate', value: 'DisableCertificate' },
    ],
  }),
  mounted() {
    this.getAdmin()
    this.getCertificates()
  },
  methods: {
    getAdmin() {
      this.axios.get("http://localhost:8081/api/admin")
        .then(resp => {
          this.admin = resp.data
          this.name = this.admin.Username
          this.email = this.admin.Email
          console.log(this.admin)
        })
        .catch(err => {
          console.log(err.response.data())
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
    }
  }
}
</script>

<style scoped>

#Admin-page {
  display: flex;
  flex-direction: row;
  align-content: center;
  height: 100%;
  background: linear-gradient(45deg, #211010 70%,  #f5efe3 30%)
}

</style>