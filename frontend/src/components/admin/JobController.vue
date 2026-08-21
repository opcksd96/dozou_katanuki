<!-- frontend/src/components/admin/JobController.vue (100行以下) -->
<script setup lang="ts">
import JobFormPanel from './jobs/JobFormPanel.vue';
import JobProgressBar from './jobs/JobProgressBar.vue';
import JobTerminalView from './jobs/JobTerminalView.vue';
import JobHistoryTable from './jobs/JobHistoryTable.vue';

defineProps<{
  activeJob: any;
  jobList: any[];
  logs: string[];
  salvageForm: { platform: string; account: string; limit: number };
  importForm: { warcPath: string; offline: boolean };
  actionLoading: boolean;
  loadingJobs: boolean;
}>();

defineEmits<{
  (e: 'startSalvage'): void;
  (e: 'startImport'): void;
  (e: 'cancelJob', jobId: string): void;
  (e: 'fetchJobs'): void;
  (e: 'clearLogs'): void;
}>();
</script>

<template>
  <div class="space-y-6">
    <JobFormPanel
      :salvage-form="salvageForm"
      :import-form="importForm"
      :action-loading="actionLoading"
      :is-job-running="Boolean(activeJob && activeJob.status === 'running')"
      @start-salvage="$emit('startSalvage')"
      @start-import="$emit('startImport')"
    />

    <JobProgressBar
      :active-job="activeJob"
      @cancel-job="(id) => $emit('cancelJob', id)"
    />

    <JobTerminalView
      :logs="logs"
      @clear-logs="$emit('clearLogs')"
    />

    <JobHistoryTable
      :job-list="jobList"
      :loading-jobs="loadingJobs"
      @fetch-jobs="$emit('fetchJobs')"
      @cancel-job="(id) => $emit('cancelJob', id)"
    />
  </div>
</template>
