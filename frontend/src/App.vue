<!-- frontend/src/App.vue (100行以下 - SPEC-FRONTEND-001) -->
<script setup lang="ts">
// ルーティングのベースコンテナとして動作します
</script>

<template>
  <router-view />
</template>

<style>
/* 共通のフェードインアニメーションなどが必要な場合はここに記述 */
.router-view-enter-active,
.router-view-leave-active {
  transition: opacity 0.3s ease;
}
.router-view-enter-from,
.router-view-leave-to {
  opacity: 0;
}
</style>
