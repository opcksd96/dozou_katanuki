<!-- frontend/src/components/admin/skin/SkinCardPreview.vue (100行以下) -->
<script setup lang="ts">
defineProps<{ platform: string }>();
</script>

<template>
  <div class="bg-slate-900/60 p-4 border border-slate-800 rounded-xl space-y-3">
    <h4 class="text-xs font-bold text-slate-200 flex items-center gap-2"><span>👁️</span> リアルタイムカードプレビュー</h4>
    <div class="border border-slate-800 rounded-xl overflow-hidden bg-slate-950 p-4">
      <div class="twitter-card">
        <div class="twitter-avatar-col">
          <div class="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-sm font-bold text-slate-300">M</div>
        </div>
        <div class="twitter-content-col">
          <div class="twitter-header-row">
            <span class="twitter-display-name">マシュ・キリエライト</span>
            <span class="twitter-handle">@mash_dev</span>
            <span class="twitter-timestamp">10分前</span>
          </div>
          <div class="twitter-body-text">
            先輩！design.css の編集内容がここにリアルタイムで反映されます。<a href="#" class="external-link">#NES_APU</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
