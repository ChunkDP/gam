<!-- src/components/Permission.vue -->
<template>
    <span style="margin-right:5px" v-if="hasPermission">
      <slot></slot>
    </span>
  </template>
  
  <script setup>
  import { computed } from 'vue';
  import { useUserPermissionsStore } from '@/stores/userPermissions';
  import { useRoute } from 'vue-router';
  
  const props = defineProps({
    permission: {
      type: String,
      required: true
    }
  });
  
  const userPermissionsStore = useUserPermissionsStore();

  const hasPermission = computed(() => {
    return userPermissionsStore.permissions.includes(props.permission);
  });
  </script>