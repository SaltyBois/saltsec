<template>
  <div id="home-main">
     <main-nav></main-nav>
    <div id="header-container">
      <h1 id="header" class="primary--text">
        SaltSec
      </h1>
      <p id="header-text" class="primary--text">
        Adding salt to all your security needs
      </p>
      <v-btn class="accent primary--text" @click="routeToGetCertificate">Get your certificate</v-btn>
      <br>
      <v-btn dark class="accent primary--text" @click="routeToLogin">Login</v-btn>
    </div>
    <v-img src="https://image.flaticon.com/icons/png/512/114/114928.png" id="logo-image" contain/>
  </div>
</template>

<script>

export default {
  name: 'Home',
  methods: {
    ping: function() {
      this.axios.get("http://localhost:8081/ping")
        .then(resp => {
          console.log(resp.data);
        })
        .finally(response => {
          console.log(response);
        })
    },
    routeToLogin()
    {
      this.$router.push('/login');
    },
    routeToGetCertificate()
    {
      this.ping();
      this.$router.push('/getCertificate');
    },
  },
}
</script>

<style scoped>
  #home-main {
    display: flex;
    flex-direction: row;
    align-content: center;
    height: 100%;
    background: linear-gradient(45deg, #211010 70%,  #f5efe3 30%)
  }

  #header-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: flex-start;
    margin-left: 4rem;
  }

  #header {
    font-size: 4rem;
    font-weight: 600;
  }

  #header-text {
    font-size: 2rem;
  }

  #logo-image {
    float: right;
    margin-left: 5%;
    margin-top: 1%;
    height: 40%;
  }

</style>
