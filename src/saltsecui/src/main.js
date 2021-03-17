import Vue from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import vuetify from './plugins/vuetify';
import MainNavigation from './components/MainNavigation'

Vue.config.productionTip = false
Vue.use(VueAxios, axios)

Vue.component('main-nav', MainNavigation);

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
