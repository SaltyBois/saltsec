<template>
 <div id="Admin-page">
   <v-layout justify-center align-baseline>
     <span class="bg"/>
     <v-card width="75%" elevation="10" justify-center class="text-center mt-3">
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
               <td>{{row.item.Cert.SerialNumber}}</td>
               <td>{{row.item.Type}}</td>
               <td>{{row.item.Cert.IsCA}}</td>
               <td>{{new Date(row.item.Cert.NotBefore).toLocaleString('sr')}}</td>
               <td>{{new Date(row.item.Cert.NotAfter).toLocaleString('sr')}}</td>
               <td>{{row.item.Cert.Issuer.SerialNumber}}</td>
               <td v-if="!isArchived(row.item.Cert.SerialNumber)">
                 <v-btn dark class="accent primary--text" @click="archiveCertificate(row.item)">Archive</v-btn>
               </td>
               <td v-else>
                 <v-btn disabled class="accent primary--text">Archived</v-btn>
               </td>
               <td>
                 <v-btn dark class="info primary--text" @click="downloadCert(row.item.Cert.SerialNumber)">Download</v-btn>
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
    archived: [],
    certificates: [],
    admin: null,
    name: '',
    email: '',
    commonName: '',
    headers: [
      { text: 'Serial Number', value: 'SerialNumber'},
      { text: 'Certificate Type', value: 'CertificateType' },
      { text: 'Is CA', value: 'IsCA' },
      { text: 'Not Before', value: 'NotBefore' },
      { text: 'Not After', value: 'NotAfter' },
      { text: 'Issued By', value: 'Issuer' },
      { text: 'Disable Certificate', value: 'DisableCertificate' },
      { text: 'Download Certificate', value: 'Download' },
    ],
  }),
  mounted() {
    this.getAdmin()
    this.getCertificates()
    this.getArchived()
  },
  computed: {
    userDn() {
      return {'username': this.username, 'password': this.password, 'commonName': this.commonName}
    }
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
            this.certificates.forEach(element => {
              this.checkIfArchived(element.Cert.serialNumber);
            });
            console.log("Certificates:")
            console.log(this.certificates)
          })
          .catch(err => {
            console.log(err.response)
          })
    },
    archiveCertificate(cert) {
      this.commonName = cert.Cert.Subject.CommonName
      let dto = {
        username: cert.Cert.EmailAddresses[0],
        password: "",
        commonName: cert.Cert.Subject.CommonName
      }
      this.axios.post("http://localhost:8081/api/cert/archive/add", dto)
      this.$router.go()
    },
    getArchived() {
      this.axios.get("http://localhost:8081/api/cert/archive")
        .then(response => {
          console.log(response);
          this.archived = response.data;
        })
        .catch(response => {
          console.log(response)
        });
    },
    downloadCert(serialNumber) {
      this.$router.push("http:/localhost:8081/api/cert/download/" + serialNumber)
    },
    isArchived(serialNumber) {
      let retval = false;
      this.archived.forEach(a => {
        if (a.serialNumber == serialNumber){
          retval = true;
          return retval
        }
      })
      return retval;
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