<!-- frontend/src/components/admin/JobController.vue (100行以下 - SPEC-PRINCIPLE-001) -->
<script setup lang="ts">
import { ref } from 'vue';
import JobFormPanel from './jobs/JobFormPanel.vue';
import JobProgressBar from './jobs/JobProgressBar.vue';
import JobTerminalView from './jobs/JobTerminalView.vue';
import MissionReportModal from './jobs/MissionReportModal.vue';

const getApp = () => (window as any)?.go?.app?.App || (window as any)?.go?.main?.App;

defineProps<{
  activeJob: any;
  logs: string[];
  salvageForm: { platform: string; account: string; source?: string; limit: number };
  importForm: { warcPath: string; offline: boolean };
  actionLoading: boolean;
}>();

defineEmits<{
  (e: 'startSalvage'): void;
  (e: 'startImport'): void;
  (e: 'cancelJob', jobId: string): void;
  (e: 'clearLogs'): void;
}>();

const isReportModalOpen = ref(false), currentReport = ref<any>(null);

const openReport = async () => {
  try {
    const rep = await getApp()?.GetLatestMissionReport?.();
    if (rep) { currentReport.value = rep; isReportModalOpen.value = true; }
  } catch (_) {}
};
</script>

<template>
  <div class="space-y-4 font-sans max-w-4xl mx-auto py-1">
    <!-- 1. スクレイパー採取 / WARCインポート 起動パネル -->
    <JobFormPanel
      :salvage-form="salvageForm"
      :import-form="importForm"
      :action-loading="actionLoading"
      :is-job-running="Boolean(activeJob && activeJob.status === 'running')"
      @start-salvage="$emit('startSalvage')"
      @start-import="$emit('startImport')"
    />

    <!-- 2. 実行中リアルタイム・モニタリング & 完遂レポート動線 -->
    <JobProgressBar
      :active-job="activeJob"
      @cancel-job="(id) => $emit('cancelJob', id)"
      @view-report="openReport"
    />

    <!-- 3. ミッション実行ログターミナル -->
    <JobTerminalView
      :logs="logs"
      @clear-logs="$emit('clearLogs')"
    />

    <!-- 4. 5W1H 完遂レポートモーダル -->
    <MissionReportModal
      :is-open="isReportModalOpen"
      :report="currentReport"
      @close="isReportModalOpen = false"
    />
  </div>
</template>
