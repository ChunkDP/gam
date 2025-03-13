<template>
  <div v-if="error">
    <h3>出错了</h3>
    <p>{{ error }}</p>
    <button @click="reset">重试</button>
  </div>
  <slot v-else></slot>
</template>

<script setup>
import { ref, onErrorCaptured } from 'vue'

const error = ref(null)

onErrorCaptured((err) => {
  error.value = err
  return false
})

const reset = () => {
  error.value = null
}
</script>