<template>
 <div>
   <v-layout justify-center align-baseline>
     <span class="bg"/>
     <v-card width="70%" elevation="10" justify-center class="text-center">
       <v-card-title style="margin: 10px">
         <v-layout justify-center style="margin: 5px">
           <h1>All Medicine Orders</h1>
         </v-layout>
         <v-layout justify-center style="margin: 5px">
           <h2>Click on each order to see more details</h2>
         </v-layout>
         <v-layout justify-center style="margin: 5px">
         </v-layout>
       </v-card-title>
       <v-layout justify-center>
         <v-data-table :items="ordersInPharmacy" :headers="headers"  >
           <template v-slot:item="row" >
             <tr>
               <td>{{row.item.id}}</td>
               <td>{{new Date(row.item.deadline).toLocaleString('sr')}}</td>
               <td>
                 <v-dialog v-model="row.item.dialog" >
                   <template v-slot:activator="{ on, attrs }">
                     <v-btn :color="colorFun(row.item)" dark v-bind="attrs" v-on="on" @click="getListOfOffers(row.item)">
                       More Details
                     </v-btn>
                   </template>
                   <v-card>
                     <v-card-title class="headline grey lighten-2">
                       Offers for order ID {{ row.item.id }}
                     </v-card-title>
                     <v-data-table :items="listOfOffers" :headers="headers2" show-select single-select v-model="selectedOffer">
                     </v-data-table>
                     <v-divider/>
                     <v-card-actions>
                       <v-spacer/>
                       <v-btn color="#24f232" text @click="acceptOffer" v-if="admin.id === row.item.pharmacyAdminId" >
                         Accept selected offer
                       </v-btn>
                       <v-btn color="#e6052a" text @click="row.item.dialog = false">
                         Exit dialog
                       </v-btn>
                     </v-card-actions>
                   </v-card>
                 </v-dialog>
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
  name: "AdminPage"
}
</script>

<style scoped>

</style>