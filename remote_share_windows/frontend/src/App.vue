<template>
  <div id="app">
    <guacamole-client></guacamole-client>
  </div>
</template>

<script>

import GuacamoleClient from "@/components/GuacClient";
export default {
  name: 'App',
  components: {
    GuacamoleClient: GuacamoleClient,
  }
}
</script>

<style>
body {
  margin: 0;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
