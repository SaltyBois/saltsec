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
         <h3>{{admin.Username}}</h3>
       </v-layout>
       <v-layout align-baseline justify-center class="mb-5">
         <h3><b>Admin email:</b></h3>
         <h3>{{admin.Email}}</h3>
       </v-layout>

       <h2> Certificates info</h2>
       <v-layout justify-center>
         <v-data-table v-bind:items="certificates">
<!--           <template >-->
<!--           </template>-->
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
    certificates: null,
    admin: null
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
            console.log(this.certificates)
          })
          .catch(err => {
            console.log(err.response.data())
          })
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