<!-- frontend/src/components/admin/AuditReportView.vue (100行以下) -->
<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useAdminAudit } from '../../composables/admin/useAdminAudit';
import AuditReportHeader from './audit/AuditReportHeader.vue';
import AuditReportOrphans from './audit/AuditReportOrphans.vue';
import AuditReportRestore from './audit/AuditReportRestore.vue';

const adminAudit = useAdminAudit();

const props = withDefaults(defineProps<{
  report?: any | null;
  loading?: boolean;
  purgingFiles?: boolean;
  purgingDb?: boolean;
  purgingDB?: boolean;
  restoring?: boolean;
  statusMessage?: { success: boolean; message: string } | null;
  restoreStatus?: { success: boolean; message: string } | null;
}>(), {
  report: null,
  loading: false,
  purgingFiles: false,
  purgingDb: false,
  purgingDB: false,
  restoring: false,
  statusMessage: null,
  restoreStatus: null,
});

const emit = defineEmits<{
  (e: 'runAudit', purgeFiles?: boolean, purgeDB?: boolean): void;
  (e: 'purgeOrphanFiles', paths?: string[]): void;
  (e: 'purgeOrphanDBMedia', mediaIDs?: string[]): void;
  (e: 'triggerRestore', resetDB?: boolean): void;
}>();

const currentReport = computed(() => props.report || adminAudit.auditReport.value);
const currentLoading = computed(() => props.loading || adminAudit.isAuditing.value);
const currentPurgingFiles = computed(() => props.purgingFiles || adminAudit.isPurgingFiles.value);
const currentPurgingDB = computed(() => props.purgingDb || props.purgingDB || adminAudit.isPurgingDB.value);
const currentStatusMessage = computed(() => props.statusMessage || adminAudit.auditStatusMessage.value);

const handleRunAudit = (pf = false, pdb = false) => {
  emit('runAudit', pf, pdb);
  adminAudit.runAudit(pf, pdb);
};

const handlePurgeFiles = (paths?: string[]) => {
  emit('purgeOrphanFiles', paths);
  adminAudit.purgeOrphanFiles(paths);
};

const handlePurgeDB = (ids?: string[]) => {
  emit('purgeOrphanDBMedia', ids);
  adminAudit.purgeOrphanDBMedia(ids);
};

onMounted(() => {
  if (!currentReport.value) handleRunAudit(false, false);
});
</script>

<template>
  <div class="space-y-4">
    <AuditReportHeader
      :report="currentReport"
      :loading="currentLoading"
      :purging-files="currentPurgingFiles"
      :status-message="currentStatusMessage"
      @run-audit="handleRunAudit"
    />

    <AuditReportOrphans
      :report="currentReport"
      :purging-files="currentPurgingFiles"
      :purging-db="currentPurgingDB"
      @purge-files="handlePurgeFiles"
      @purge-db="handlePurgeDB"
    />

    <AuditReportRestore
      :restoring="restoring"
      :restore-status="restoreStatus"
      @trigger-restore="(resetDB) => emit('triggerRestore', resetDB)"
    />
  </div>
</template>
