<!-- frontend/src/components/admin/database/MediaInspectorPanel.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useAdminDatabase } from '../../../composables/admin/useAdminDatabase';
import MediaInspectorAccount from './MediaInspectorAccount.vue';
import MediaInspectorStash from './MediaInspectorStash.vue';
import MediaInspectorSqlite from './MediaInspectorSqlite.vue';

const props = defineProps<{ media: any }>();
const emit = defineEmits<{
  (e: 'saveMetadata', payload: { mediaId: string; downloadStatus: string; stashSceneId: string; stashImageId: string; failedReason: string }): void;
  (e: 'retry', mediaId: string): void;
  (e: 'purge', mediaId: string): void;
  (e: 'escalateThunder', m: any): void;
  (e: 'viewPost', articleId: string): void;
  (e: 'viewPostTimeline', articleId: string): void;
  (e: 'openExplorer', mediaId: string): void;
  (e: 'openDefault', mediaId: string): void;
}>();

const { fetchStashMetadata, updateStashMetadata } = useAdminDatabase();

const editStatus = ref(props.media.download_status || props.media.raw_status || 'QUEUED');
const editReason = ref(props.media.failed_reason || '');
const stashData = ref<any>(null);
const editTitle = ref(''), editDetails = ref(''), editRating = ref(0);
const isMutating = ref(false);
const undoSnapshot = ref<{ title: string; details: string; rating100: number } | null>(null);

const isStashModified = computed(() => {
  if (!stashData.value) return false;
  return editTitle.value !== (stashData.value.title || '') ||
         editDetails.value !== (stashData.value.details || '') ||
         editRating.value !== (stashData.value.rating100 || 0);
});

const loadStashInfo = async () => {
  if (!props.media.stash_scene_id && !props.media.stash_image_id) { stashData.value = null; return; }
  const res = await fetchStashMetadata(props.media.stash_scene_id || '', props.media.stash_image_id || '');
  if (res) {
    stashData.value = res;
    editTitle.value = res.title || ''; editDetails.value = res.details || ''; editRating.value = res.rating100 || 0;
  }
};

watch(() => props.media, (m) => {
  if (m) {
    editStatus.value = m.download_status || m.raw_status || 'QUEUED';
    editReason.value = m.failed_reason || '';
    undoSnapshot.value = null;
    loadStashInfo();
  }
}, { immediate: true });

const handleSaveSqlite = () => {
  emit('saveMetadata', {
    mediaId: props.media.media_id || props.media.id, downloadStatus: editStatus.value,
    stashSceneId: props.media.stash_scene_id || '', stashImageId: props.media.stash_image_id || '',
    failedReason: editReason.value.trim(),
  });
};

const handleSaveStash = async () => {
  const targetId = props.media.stash_scene_id || props.media.stash_image_id;
  if (!targetId) return;
  isMutating.value = true;
  if (stashData.value) undoSnapshot.value = { title: stashData.value.title || '', details: stashData.value.details || '', rating100: stashData.value.rating100 || 0 };
  const updated = await updateStashMetadata(!!props.media.stash_scene_id, targetId, editTitle.value.trim(), editDetails.value.trim(), editRating.value);
  if (updated) { stashData.value = updated; editTitle.value = updated.title || ''; editDetails.value = updated.details || ''; editRating.value = updated.rating100 || 0; }
  isMutating.value = false;
};

const handleUndoStash = async () => {
  const targetId = props.media.stash_scene_id || props.media.stash_image_id;
  if (!targetId || !undoSnapshot.value) return;
  isMutating.value = true;
  const snap = undoSnapshot.value;
  const restored = await updateStashMetadata(!!props.media.stash_scene_id, targetId, snap.title, snap.details, snap.rating100);
  if (restored) { stashData.value = restored; editTitle.value = restored.title || ''; editDetails.value = restored.details || ''; editRating.value = restored.rating100 || 0; undoSnapshot.value = null; }
  isMutating.value = false;
};
</script>

<template>
  <div class="w-80 md:w-[420px] bg-slate-900/95 border-l border-slate-800 flex flex-col p-4 space-y-3.5 overflow-y-auto text-xs font-mono text-slate-300 shrink-0">
    <MediaInspectorAccount :media="media" :has-stash="!!(media.stash_scene_id || media.stash_image_id)" @copy-to-stash-details="editDetails = $event" />
    <MediaInspectorStash :media="media" :stash-data="stashData" v-model:edit-title="editTitle" v-model:edit-details="editDetails" v-model:edit-rating="editRating" :is-modified="isStashModified" :is-mutating="isMutating" :has-undo="!!undoSnapshot" @request-save="handleSaveStash" @undo="handleUndoStash" />
    <MediaInspectorSqlite :media="media" v-model:edit-status="editStatus" v-model:edit-reason="editReason" @save="handleSaveSqlite" @retry="emit('retry', media.media_id || media.id)" @purge="emit('purge', media.media_id || media.id)" @escalate-thunder="emit('escalateThunder', media)" @view-post="emit('viewPost', media.article_id)" @view-post-timeline="emit('viewPostTimeline', media.article_id)" @open-explorer="emit('openExplorer', media.media_id || media.id)" @open-default="emit('openDefault', media.media_id || media.id)" />
  </div>
</template>
