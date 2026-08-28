<!-- frontend/src/components/admin/database/MediaTableView.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import Avatar from '../../article/Avatar.vue';

defineProps<{ items: any[] }>();

const emit = defineEmits<{
  (e: 'select', m: any): void; (e: 'retry', mediaId: string): void; (e: 'purge', mediaId: string): void;
  (e: 'openExplorer', mediaId: string): void; (e: 'openDefault', mediaId: string): void; (e: 'copy', media: any): void; (e: 'toggleBookmark', mediaId: string): void;
}>();

const parseUrls = (val: any) => {
  if (Array.isArray(val)) return val;
  if (typeof val === 'string' && val.startsWith('[')) {
    try { return JSON.parse(val); } catch { return []; }
  }
  return [];
};
</script>

<template>
  <div class="overflow-x-auto border border-slate-800 rounded-xl bg-slate-950">
    <table class="w-full text-left text-xs font-mono">
      <thead class="bg-slate-900 text-slate-400 border-b border-slate-800 text-[11px]">
        <tr>
          <th class="p-2.5 w-8 text-center">⭐</th><th class="p-2.5 w-16 text-center">プレビュー</th>
          <th class="p-2.5">タイトル / 記事詳細</th><th class="p-2.5">アカウント</th><th class="p-2.5">種別</th>
          <th class="p-2.5">解像度</th><th class="p-2.5">ステータス</th><th class="p-2.5">Stash ID</th><th class="p-2.5 text-right">アクション</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-850">
        <tr v-for="m in items" :key="m.media_id || m.id" @click="emit('select', m)" class="hover:bg-slate-900/60 cursor-pointer transition-colors">
          <td class="p-2.5 text-center" @click.stop="emit('toggleBookmark', m.media_id || m.id)">
            <button class="hover:scale-125 transition-transform text-sm">{{ m.is_bookmarked ? '⭐' : '☆' }}</button>
          </td>
          <td class="p-2 text-center">
            <div class="w-14 h-14 bg-slate-900 rounded-lg overflow-hidden flex items-center justify-center mx-auto border border-slate-800">
              <img v-if="m.urls?.thumbnail" :src="m.urls.thumbnail" class="w-full h-full object-contain" loading="lazy" />
              <span v-else class="text-sm">{{ m.type === 'video' ? '🎬' : '🖼️' }}</span>
            </div>
          </td>
          <td class="p-2.5 max-w-[280px]">
            <div class="font-bold text-slate-200 truncate text-xs flex items-center gap-1" :title="m.title || m.media_id">
              <span>{{ m.title || `X (@${m.username}): Tweet ${m.article_id}` }}</span>
              <span v-if="parseUrls(m.tweet_urls).length > 1" class="px-1 py-0.2 rounded bg-indigo-950 text-indigo-300 border border-indigo-700/50 text-[9px]" :title="parseUrls(m.tweet_urls).join('\n')">+{{ parseUrls(m.tweet_urls).length - 1 }}件</span>
            </div>
            <div v-if="m.tweet_date" class="text-[10px] text-slate-400 font-semibold">{{ m.tweet_date }}</div>
            <div v-if="m.full_text || m.full_text_ja" class="text-[11px] text-slate-400 font-sans truncate" :title="m.full_text_ja || m.full_text">
              {{ m.full_text_ja || m.full_text }}
            </div>
          </td>
          <td class="p-2.5 text-slate-300">
            <div class="flex items-center gap-1.5"><Avatar :avatar-url="m.avatar_url" :handle="m.username" size-class="w-4 h-4" /><span>@{{ m.username }}</span></div>
          </td>
          <td class="p-2.5 text-slate-400 uppercase text-[10px]">{{ m.type }}</td>
          <td class="p-2.5 text-slate-400 text-[10px]">
            <div class="flex items-center gap-1">
              <span v-if="m.media_quality" class="px-1 py-0.2 rounded bg-blue-950 text-blue-300 font-bold border border-blue-700/50 uppercase text-[9px]">{{ m.media_quality }}</span>
              <span>{{ m.width && m.height ? `${m.width}x${m.height}` : '-' }}</span>
            </div>
          </td>
          <td class="p-2.5">
            <span class="px-2 py-0.5 rounded text-[10px] font-bold" :class="{
              'bg-emerald-950 text-emerald-300 border border-emerald-700': m.download_status === 'COMPLETED',
              'bg-purple-950 text-purple-300 border border-purple-700': m.download_status === 'OUTSOURCED',
              'bg-amber-950 text-amber-300 border border-amber-600': m.download_status === 'RETAINED',
              'bg-slate-800 text-slate-300 border border-slate-700': m.download_status === 'QUEUED',
              'bg-rose-950 text-rose-300 border border-rose-700': m.download_status === 'DEAD_404',
              'bg-red-950 text-red-300 border border-red-700': m.download_status === 'FAILED',
              'bg-zinc-800 text-zinc-400': m.download_status === 'EXCLUDED'
            }">{{ m.download_status }}</span>
          </td>
          <td class="p-2.5 text-[11px] text-purple-400 font-bold">{{ m.stash_scene_id || m.stash_image_id ? (m.stash_scene_id || m.stash_image_id).slice(0, 8) : '-' }}</td>
          <td class="p-2.5 text-right space-x-1" @click.stop>
            <button @click.stop="emit('openExplorer', m.media_id || m.id)" class="p-1.5 bg-slate-900 hover:bg-slate-800 text-slate-300 rounded-lg text-xs" title="エクスプローラー">📂</button>
            <button @click.stop="emit('openDefault', m.media_id || m.id)" class="p-1.5 bg-slate-900 hover:bg-slate-800 text-slate-300 rounded-lg text-xs" title="既定アプリ">🚀</button>
            <button @click.stop="emit('copy', m)" class="p-1.5 bg-slate-900 hover:bg-slate-800 text-slate-300 rounded-lg text-xs" title="コピー">📋</button>
            <button @click="emit('retry', m.media_id || m.id)" class="p-1.5 bg-blue-950/80 hover:bg-blue-900 border border-blue-700/60 text-blue-300 rounded-lg text-xs" title="再取得 (リトライ)">🔄</button>
            <button @click="emit('purge', m.media_id || m.id)" class="p-1.5 bg-rose-950/80 hover:bg-rose-900 border border-rose-700/60 text-rose-300 rounded-lg text-xs" title="DBからパージ">🗑️</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
