<template>
  <div id="home-main">
    <!-- <main-nav></main-nav> -->
    <div id="header-container">
      <h1 id="header" class="primary--text">
        SaltSec
      </h1>
      <p id="header-text" class="primary--text">
        Adding salt to all your security needs
      </p>
      <v-btn class="accent primary--text" @click="ping">Get your certificate</v-btn>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Home',
  methods: {
    ping: function() {
      this.axios.get("http://localhost:8081/ping")
        .finally(response => {
          console.log(response);
        })
    },
  },
}
</script>

<style scoped>
  #home-main {
    display: flex;
    flex-direction: row;
    align-content: center;
    height: 100vh;
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

</style>
