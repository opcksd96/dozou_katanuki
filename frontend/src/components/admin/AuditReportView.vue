<!-- frontend/src/components/admin/AuditReportView.vue (100行以下) -->
<script setup lang="ts">
import { onMounted } from 'vue';
import AuditReportHeader from './audit/AuditReportHeader.vue';
import AuditReportOrphans from './audit/AuditReportOrphans.vue';
import AuditReportRestore from './audit/AuditReportRestore.vue';

const props = defineProps<{
  report: any | null;
  loading: boolean;
  purgingFiles: boolean;
  purgingDB: boolean;
  restoring?: boolean;
  statusMessage: { success: boolean; message: string } | null;
  restoreStatus?: { success: boolean; message: string } | null;
}>();

const emit = defineEmits<{
  (e: 'runAudit', purgeFiles?: boolean, purgeDB?: boolean): void;
  (e: 'purgeOrphanFiles', paths?: string[]): void;
  (e: 'purgeOrphanDBMedia', mediaIDs?: string[]): void;
  (e: 'triggerRestore', resetDB?: boolean): void;
}>();

onMounted(() => {
  if (!props.report) emit('runAudit', false, false);
});
</script>

<template>
  <div class="space-y-4">
    <AuditReportHeader
      :report="report"
      :loading="loading"
      :purging-files="purgingFiles"
      :status-message="statusMessage"
      @run-audit="(pf, pdb) => emit('runAudit', pf, pdb)"
    />

    <AuditReportOrphans
      :report="report"
      :purging-files="purgingFiles"
      :purging-db="purgingDB"
      @purge-files="() => emit('purgeOrphanFiles')"
      @purge-db="() => emit('purgeOrphanDBMedia')"
    />

    <AuditReportRestore
      :restoring="restoring"
      :restore-status="restoreStatus"
      @trigger-restore="(resetDB) => emit('triggerRestore', resetDB)"
    />
  </div>
</template>
