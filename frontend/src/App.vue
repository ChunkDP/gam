<!-- App.vue -->
<template>

    <ErrorBoundary>
    <router-view></router-view>
  </ErrorBoundary>

</template>

<script>
import { defineComponent } from 'vue';
import ErrorBoundary from './components/ErrorBoundary.vue';

export default defineComponent({
  name: 'App',
  components: {
    ErrorBoundary
  }
});
</script>
<style lang="scss">
// 可以直接使用scss语法
$primary-color: #409EFF;

.container {
  color: $primary-color;
}
</style>
