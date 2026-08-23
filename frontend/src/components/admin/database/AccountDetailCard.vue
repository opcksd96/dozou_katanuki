<!-- frontend/src/components/admin/database/AccountDetailCard.vue (100行以下) -->
<script setup lang="ts">
import { ref, watch } from 'vue';
import Avatar from '../../article/Avatar.vue';

const props = defineProps<{
  account: any;
  postCount: number;
  isEditing: boolean;
  availableAvatars?: string[];
}>();

const emit = defineEmits<{
  (e: 'startEdit'): void;
  (e: 'cancelEdit'): void;
  (e: 'save', payload: { displayName: string; username: string; avatarUrl: string; description: string }): void;
  (e: 'viewPosts'): void;
  (e: 'viewMedia'): void;
}>();

const editForm = ref({ displayName: '', username: '', avatarUrl: '', description: '' });
const showPicker = ref(false);

watch(() => props.account, (acc) => {
  if (acc) {
    editForm.value = {
      displayName: acc.display_name || '', username: acc.username || '',
      avatarUrl: acc.avatar_url || '', description: acc.description || '',
    };
  }
}, { immediate: true });

const save = () => emit('save', { ...editForm.value });
</script>

<template>
  <div v-if="account" class="p-4 bg-slate-900 border border-slate-800 rounded-xl space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Avatar :avatar-url="isEditing ? editForm.avatarUrl : (account?.avatar_url || '')" :handle="account?.username || ''" size-class="w-14 h-14" />
        <div>
          <h3 class="text-base font-bold text-slate-100">{{ account?.display_name || '名称未設定' }}</h3>
          <p class="text-xs text-slate-400">@{{ account?.username || '' }} <span class="text-slate-500">({{ account?.numeric_id || '' }})</span></p>
        </div>
      </div>
      <div class="flex gap-2">
        <button v-if="!isEditing" @click="emit('startEdit')" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-blue-400 rounded-lg text-xs font-bold cursor-pointer">編集</button>
        <template v-else>
          <button @click="save" class="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-xs font-bold cursor-pointer">保存</button>
          <button @click="emit('cancelEdit')" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-xs cursor-pointer">キャンセル</button>
        </template>
      </div>
    </div>

    <!-- 編集フォーム -->
    <div v-if="isEditing" class="space-y-3 pt-2 border-t border-slate-800 text-xs">
      <div><label class="text-slate-400 block mb-1">表示名</label><input v-model="editForm.displayName" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-slate-100" /></div>
      <div><label class="text-slate-400 block mb-1">ユーザー名 (@)</label><input v-model="editForm.username" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-slate-100" /></div>
      <div>
        <div class="flex justify-between items-center mb-1"><label class="text-slate-400">アバターURL</label><button @click="showPicker = !showPicker" class="text-blue-400 text-[11px] cursor-pointer">プリセット選択</button></div>
        <input v-model="editForm.avatarUrl" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-slate-100" />
      </div>
      <div v-if="showPicker && availableAvatars?.length" class="p-2 bg-slate-950 rounded border border-slate-800 flex gap-2 flex-wrap max-h-32 overflow-y-auto">
        <img v-for="av in availableAvatars" :key="av" :src="av" @click="editForm.avatarUrl = av; showPicker = false" class="w-8 h-8 rounded-full border border-slate-700 cursor-pointer hover:scale-110" />
      </div>
      <div><label class="text-slate-400 block mb-1">自己紹介</label><textarea v-model="editForm.description" rows="2" class="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-slate-100"></textarea></div>
    </div>
    <div v-else class="text-xs text-slate-300 whitespace-pre-wrap">{{ account?.description || '（自己紹介未設定）' }}</div>

    <div class="flex gap-2 pt-2 border-t border-slate-800 text-xs">
      <button @click="emit('viewPosts')" class="flex-1 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded font-mono cursor-pointer">📝 投稿一覧 ({{ postCount }})</button>
      <button @click="emit('viewMedia')" class="flex-1 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded font-mono cursor-pointer">🖼️ メディア一覧</button>
    </div>
  </div>
</template>

