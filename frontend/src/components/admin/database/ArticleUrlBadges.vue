<!-- frontend/src/components/admin/database/ArticleUrlBadges.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
  domain?: string; originalUrl?: string; waybackUrl?: string;
  sotweUrl?: string; nitterUrl?: string; twistalkerUrl?: string;
}>();

const copied = ref<string | null>(null);

const copy = (val?: string, key?: string) => {
  if (!val || !key) return;
  navigator.clipboard.writeText(val);
  copied.value = key;
  setTimeout(() => { if (copied.value === key) copied.value = null; }, 1500);
};
</script>

<template>
  <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-2.5 space-y-1.5 font-sans">
    <div class="flex items-center justify-between text-[11px] font-bold text-slate-300">
      <span class="flex items-center gap-1.5">
        <span>🔗 外部参照 & アーカイブリンク</span>
        <span v-if="domain" class="px-1.5 py-0.2 rounded text-[10px] font-mono" :class="domain.includes('x.com') ? 'bg-zinc-800 text-zinc-200 border border-zinc-700' : 'bg-sky-950 text-sky-300 border border-sky-800'">
          {{ domain }}
        </span>
      </span>
      <span v-if="copied" class="text-emerald-400 text-[10px] animate-pulse">📋 {{ copied }} をコピーしました!</span>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-3 gap-1.5 text-[11px]">
      <span v-if="!originalUrl && !waybackUrl && !sotweUrl && !nitterUrl && !twistalkerUrl" class="text-slate-500 text-[10px] col-span-full py-0.5">
        登録された外部アーカイブURLはありません
      </span>

      <a v-if="originalUrl" :href="originalUrl" target="_blank" rel="noopener noreferrer" class="flex items-center justify-between p-1.5 bg-slate-900/90 hover:bg-blue-950/60 border border-slate-800 hover:border-blue-700/60 rounded-lg text-slate-300 hover:text-blue-300 transition-colors group" :title="originalUrl">
        <span class="truncate font-mono text-[10px]">🌐 元投稿 (SNS)</span>
        <button @click.prevent.stop="copy(originalUrl, '元投稿URL')" class="text-[10px] text-slate-500 hover:text-white px-1">📋</button>
      </a>

      <a v-if="waybackUrl" :href="waybackUrl" target="_blank" rel="noopener noreferrer" class="flex items-center justify-between p-1.5 bg-slate-900/90 hover:bg-amber-950/60 border border-slate-800 hover:border-amber-700/60 rounded-lg text-slate-300 hover:text-amber-300 transition-colors group" :title="waybackUrl">
        <span class="truncate font-mono text-[10px]">🏛️ Wayback</span>
        <button @click.prevent.stop="copy(waybackUrl, 'Wayback URL')" class="text-[10px] text-slate-500 hover:text-white px-1">📋</button>
      </a>

      <a v-if="sotweUrl" :href="sotweUrl" target="_blank" rel="noopener noreferrer" class="flex items-center justify-between p-1.5 bg-slate-900/90 hover:bg-purple-950/60 border border-slate-800 hover:border-purple-700/60 rounded-lg text-slate-300 hover:text-purple-300 transition-colors group" :title="sotweUrl">
        <span class="truncate font-mono text-[10px]">⚡ Sotwe</span>
        <button @click.prevent.stop="copy(sotweUrl, 'Sotwe URL')" class="text-[10px] text-slate-500 hover:text-white px-1">📋</button>
      </a>

      <a v-if="nitterUrl" :href="nitterUrl" target="_blank" rel="noopener noreferrer" class="flex items-center justify-between p-1.5 bg-slate-900/90 hover:bg-emerald-950/60 border border-slate-800 hover:border-emerald-700/60 rounded-lg text-slate-300 hover:text-emerald-300 transition-colors group" :title="nitterUrl">
        <span class="truncate font-mono text-[10px]">🐦 Nitter</span>
        <button @click.prevent.stop="copy(nitterUrl, 'Nitter URL')" class="text-[10px] text-slate-500 hover:text-white px-1">📋</button>
      </a>

      <a v-if="twistalkerUrl" :href="twistalkerUrl" target="_blank" rel="noopener noreferrer" class="flex items-center justify-between p-1.5 bg-slate-900/90 hover:bg-indigo-950/60 border border-slate-800 hover:border-indigo-700/60 rounded-lg text-slate-300 hover:text-indigo-300 transition-colors group" :title="twistalkerUrl">
        <span class="truncate font-mono text-[10px]">🔍 Twistalker</span>
        <button @click.prevent.stop="copy(twistalkerUrl, 'Twistalker URL')" class="text-[10px] text-slate-500 hover:text-white px-1">📋</button>
      </a>
    </div>
  </div>
</template>
