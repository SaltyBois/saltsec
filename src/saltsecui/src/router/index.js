import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Login from "../views/Login.vue";
import CertificateForm from "../views/CertificateForm.vue";
import AdminPage from "../views/AdminPage.vue";
import UserPage from "@/views/UserPage";

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/getCertificate',
    name: 'CertificateFrom',
    component: CertificateForm
  },
  {
    path: '/admin',
    name: 'AdminPage',
    component: AdminPage
  },
  {
    path: '/user/:username',
    name: 'UserPage',
    component: UserPage
  },
//  {
//    path: '/about',
//    name: 'About',
//    // route level code-splitting
//    // this generates a separate chunk (about.[hash].js) for this route
//    // which is lazy-loaded when the route is visited.
//    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
//  }

]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
