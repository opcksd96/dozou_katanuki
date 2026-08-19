<script setup lang="ts">
defineProps<{
  avatarUrl: string;
  handle: string;
}>();
</script>

<template>
  <div class="relative flex-shrink-0">
    <img
      :src="avatarUrl"
      :alt="handle"
      class="w-10 h-10 rounded-full border border-slate-700 object-cover bg-slate-800"
      loading="lazy"
    />
  </div>
</template>
