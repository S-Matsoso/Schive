<template>
  <v-app>
  <v-layout>
    <SideBar />
    <v-main>
      <div class="page-container">
        <router-view />
      </div>

    </v-main>
  </v-layout>
  </v-app>
</template>

<script setup>
  import SideBar from './components/SideBar.vue';
</script>

<style>
.page-container {
  padding: 24px;
  height: 100%;
  width: 100%;
}

</style>