<!-- frontend/src/components/admin/database/AccountDetailCard.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, watch } from 'vue';
import Avatar from '../../article/Avatar.vue';
import { resolveAvatarUrl } from '../../../utils/avatar';

const props = defineProps<{ account: any; postCount: number; isEditing: boolean; availableAvatars?: string[]; }>();
const emit = defineEmits<{
  (e: 'startEdit'): void; (e: 'cancelEdit'): void;
  (e: 'save', payload: { displayName: string; username: string; avatarUrl: string; description: string; aliasOf: string; groupName: string }): void;
  (e: 'toggleWhitelist'): void; (e: 'viewPosts'): void; (e: 'viewMedia'): void;
  (e: 'trashAccount', id: string): void; (e: 'restoreAccount', id: string): void;
}>();

const editForm = ref({ displayName: '', username: '', avatarUrl: '', description: '', aliasOf: '', groupName: '' });
const showPicker = ref(false);

watch(() => props.account, (acc) => {
  if (acc) {
    editForm.value = {
      displayName: acc.display_name || '', username: acc.username || '', avatarUrl: acc.avatar_url || '',
      description: acc.description || '', aliasOf: acc.alias_of || '', groupName: acc.group_name || '',
    };
  }
}, { immediate: true });

const save = () => emit('save', { ...editForm.value });
</script>

<template>
  <div v-if="account" class="p-4 bg-slate-900 border border-slate-800 rounded-xl space-y-3 font-sans">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3 min-w-0">
        <Avatar :avatar-url="isEditing ? editForm.avatarUrl : resolveAvatarUrl(account)" :handle="account?.username || ''" size-class="w-12 h-12" />
        <div class="min-w-0">
          <div class="flex items-center gap-1.5">
            <h3 class="text-sm md:text-base font-bold text-slate-100 truncate">{{ account?.display_name || '名称未設定' }}</h3>
            <span v-if="account?.is_trash" class="px-1.5 py-0.2 bg-rose-950 text-rose-300 border border-rose-800/60 rounded text-[9px] font-bold">🗑️ ゴミ箱</span>
          </div>
          <p class="text-xs text-slate-400 font-mono">@{{ account?.username || '' }} <span class="text-slate-500 text-[10px]">({{ account?.numeric_id || '' }})</span></p>
        </div>
      </div>
      <div class="flex items-center gap-2 shrink-0">
        <button @click="emit('toggleWhitelist')" class="px-2.5 py-1 rounded text-xs font-bold font-mono transition-colors cursor-pointer shadow-sm" :class="account?.is_whitelist ? 'bg-emerald-950 text-emerald-300 border border-emerald-700/60' : 'bg-slate-800 text-slate-400 hover:text-slate-200'">
          {{ account?.is_whitelist ? '🛡️ Whitelist ON' : '⚪ Whitelist OFF' }}
        </button>
        <button v-if="!isEditing" @click="emit('startEdit')" class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-blue-400 rounded text-xs font-bold cursor-pointer">編集</button>
        <template v-else>
          <button @click="save" class="px-3 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded text-xs font-bold cursor-pointer">保存</button>
          <button @click="emit('cancelEdit')" class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded text-xs cursor-pointer">取消</button>
        </template>
      </div>
    </div>

    <div v-if="!isEditing" class="flex flex-wrap items-center justify-between gap-1.5 pt-1 border-t border-slate-800 text-[10px] font-mono">
      <div class="flex flex-wrap gap-1.5">
        <span v-if="account?.group_name" class="px-2 py-0.5 bg-indigo-900/50 text-indigo-300 rounded">グループ: {{ account.group_name }}</span>
        <span v-if="account?.alias_of" class="px-2 py-0.5 bg-amber-900/50 text-amber-300 rounded">エイリアス: @{{ account.alias_of }}</span>
        <span v-if="!account?.group_name && !account?.alias_of" class="text-slate-500">（グループ・エイリアス未設定）</span>
      </div>
      <button v-if="!account?.is_trash" @click="emit('trashAccount', account.numeric_id)" class="px-2 py-0.5 bg-rose-950/80 hover:bg-rose-900 text-rose-300 border border-rose-800/60 rounded text-[10px] cursor-pointer" title="このアカウントをゴミ箱へ退避">🗑️ ゴミ箱へ</button>
      <button v-else @click="emit('restoreAccount', account.numeric_id)" class="px-2 py-0.5 bg-emerald-950/80 hover:bg-emerald-900 text-emerald-300 border border-emerald-800/60 rounded text-[10px] cursor-pointer" title="このアカウントを通常復元">♻️ 復元する</button>
    </div>

    <!-- 編集フォーム -->
    <div v-if="isEditing" class="space-y-2.5 pt-2 border-t border-slate-800 text-xs">
      <div><label class="text-slate-400 block mb-1">表示名</label><input v-model="editForm.displayName" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1 text-slate-100" /></div>
      <div><label class="text-slate-400 block mb-1">ユーザー名 (@)</label><input v-model="editForm.username" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1 text-slate-100" /></div>
      <div>
        <div class="flex justify-between items-center mb-1"><label class="text-slate-400">アバターURL</label><button @click="showPicker = !showPicker" class="text-blue-400 text-[11px] cursor-pointer">プリセット選択</button></div>
        <input v-model="editForm.avatarUrl" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1 text-slate-100" />
      </div>
      <div v-if="showPicker && availableAvatars?.length" class="p-2 bg-slate-950 rounded border border-slate-800 flex gap-2 flex-wrap max-h-24 overflow-y-auto">
        <img v-for="av in availableAvatars" :key="av" :src="av" @click="editForm.avatarUrl = av; showPicker = false" class="w-7 h-7 rounded-full border border-slate-700 cursor-pointer hover:scale-110" />
      </div>
      <div><label class="text-slate-400 block mb-1">自己紹介</label><textarea v-model="editForm.description" rows="2" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1 text-slate-100 text-xs"></textarea></div>
      <div class="flex gap-2">
        <div class="flex-1"><label class="text-slate-400 block mb-1">名寄せ先 (@)</label><input v-model="editForm.aliasOf" placeholder="空欄 = 独立" class="w-full bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-100 text-xs" /></div>
        <div class="flex-1"><label class="text-slate-400 block mb-1">グループ名</label><input v-model="editForm.groupName" placeholder="空欄 = なし" class="w-full bg-slate-950 border border-slate-700 rounded px-2 py-1 text-slate-100 text-xs" /></div>
      </div>
    </div>
    <div v-else class="text-xs text-slate-300 whitespace-pre-wrap leading-relaxed bg-slate-950/70 p-2 rounded border border-slate-800/80">{{ account?.description || '（自己紹介未設定）' }}</div>

    <div class="flex gap-2 pt-1 border-t border-slate-800 text-xs">
      <button @click="emit('viewPosts')" class="flex-1 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded font-mono cursor-pointer">📝 投稿一覧 ({{ postCount }})</button>
      <button @click="emit('viewMedia')" class="flex-1 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded font-mono cursor-pointer">🖼️ メディア一覧</button>
    </div>
  </div>
</template>
