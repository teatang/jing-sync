<template>
  <Layout v-if="token" />
  <LoginView v-if="!token" />
</template>

<script setup lang="ts">
import Layout from "@/views/Layout.vue";
import LoginView from "@/views/LoginView.vue";
import {getToken} from "@/utils/token";

const token = getToken();
</script>
