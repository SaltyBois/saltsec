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
       <h2> Admin Info</h2>
       <!--   Place for basic Admin information    -->
       <h2> Certificates info</h2>
       <v-layout justify-center>
         <v-data-table>
           <template >
           </template>
         </v-data-table>
       </v-layout>
     </v-card>
   </v-layout>
 </div>
</template>

<script>
export default {
  name: "AdminPage"
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