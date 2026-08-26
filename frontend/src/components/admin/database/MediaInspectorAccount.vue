<!-- frontend/src/components/admin/database/MediaInspectorAccount.vue (100行以下) -->
<script setup lang="ts">
import { ref } from 'vue';
import Avatar from '../../article/Avatar.vue';

const props = defineProps<{
  media: any;
  hasStash: boolean;
}>();

const emit = defineEmits<{
  (e: 'copyToStashDetails', text: string): void;
}>();

const copiedField = ref<string | null>(null);
const copyToClipboard = (text: string, field: string) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text);
    copiedField.value = field;
    setTimeout(() => { copiedField.value = null; }, 2000);
  }
};
</script>

<template>
  <div class="space-y-3">
    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2">
      <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">👤 投稿者アカウント</div>
      <div class="flex items-center gap-2.5">
        <Avatar :avatar-url="media.avatar_url" :handle="media.username" size-class="w-10 h-10" />
        <div class="min-w-0 flex-1 leading-tight">
          <div class="text-slate-100 font-bold truncate text-[13px]">{{ media.display_name || media.username || 'Unknown' }}</div>
          <div class="text-slate-400 text-[11px] truncate">@{{ media.username }} <span class="text-slate-500">({{ media.account_id || '-' }})</span></div>
        </div>
      </div>
      <div v-if="media.created_at || media.tweet_date" class="text-[10px] text-slate-500 pt-1 border-t border-slate-850 flex items-center justify-between">
        <span>投稿日: {{ media.tweet_date || (media.created_at ? new Date(media.created_at).toLocaleString() : '-') }}</span>
      </div>
    </div>

    <div class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2">
      <div class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">📷 メディア情報</div>
      <div class="space-y-1.5 text-[11px] font-mono">
        <div class="flex items-center justify-between">
          <span class="text-slate-400">メディアID</span>
          <span class="text-slate-200">{{ media.media_id || media.id || '-' }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-slate-400">種別</span>
          <span class="text-slate-200 uppercase">{{ media.type || '-' }}</span>
        </div>
        <div v-if="media.width && media.height" class="flex items-center justify-between">
          <span class="text-slate-400">解像度</span>
          <span class="text-slate-200">{{ media.width }} x {{ media.height }}</span>
        </div>
        <div v-if="media.download_url" class="space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-slate-400">原本URL</span>
            <button @click="copyToClipboard(media.download_url, 'dl')" class="text-blue-400 hover:text-blue-300">
              {{ copiedField === 'dl' ? '✓ コピー済' : '📋 コピー' }}
            </button>
          </div>
          <div class="text-slate-200 break-all leading-relaxed">{{ media.download_url }}</div>
        </div>
        <div v-if="media.wayback_url" class="space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-slate-400">Wayback原本URL</span>
            <button @click="copyToClipboard(media.wayback_url, 'wb')" class="text-blue-400 hover:text-blue-300">
              {{ copiedField === 'wb' ? '✓ コピー済' : '📋 コピー' }}
            </button>
          </div>
          <div class="text-slate-200 break-all leading-relaxed">{{ media.wayback_url }}</div>
        </div>
      </div>
    </div>

    <div v-if="media.full_text || media.full_text_ja" class="p-3 bg-slate-950/80 rounded-xl border border-slate-800 space-y-2">
      <div class="flex items-center justify-between text-[10px]">
        <span class="text-slate-400 font-bold uppercase tracking-wider">📝 親ツイート本文</span>
        <button @click="copyToClipboard(media.full_text_ja || media.full_text, 'tweet')" class="text-blue-400 hover:text-blue-300">
          {{ copiedField === 'tweet' ? '✓ コピー済' : '📋 コピー' }}
        </button>
      </div>
      <div class="p-2.5 bg-slate-900 rounded-lg border border-slate-800 space-y-1.5 select-text">
        <div v-if="media.full_text_ja" class="text-xs text-slate-100 leading-relaxed font-sans">{{ media.full_text_ja }}</div>
        <div v-if="media.full_text && media.full_text !== media.full_text_ja" class="text-[11px] text-slate-400 leading-relaxed font-sans border-t border-slate-800 pt-1">{{ media.full_text }}</div>
      </div>
      <button
        v-if="hasStash"
        @click="emit('copyToStashDetails', media.full_text_ja ? media.full_text_ja + '\n\n' + media.full_text : media.full_text)"
        class="w-full py-1 bg-slate-800 hover:bg-slate-700 text-purple-300 border border-purple-700/40 rounded text-[10px] font-mono flex items-center justify-center gap-1 transition-colors"
      >
        <span>⬇️</span> 本文を Stash 詳細メモへ転記
      </button>
    </div>
  </div>
</template>
