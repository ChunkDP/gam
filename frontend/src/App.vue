<!-- App.vue -->
<template>

  <ErrorBoundary>
    <router-view></router-view>
  </ErrorBoundary>

</template>

<script setup>

import ErrorBoundary from './components/ErrorBoundary.vue';


</script>
<style lang="scss">
// 可以直接使用scss语法
$primary-color: #409EFF;

.container {
  color: $primary-color;
}
</style>
