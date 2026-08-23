<!-- frontend/src/components/admin/AuditReportView.vue (100行以下) -->
<script setup lang="ts">
import { onMounted } from 'vue';
import { useAdminAudit } from '../../composables/admin/useAdminAudit';
import AuditReportHeader from './audit/AuditReportHeader.vue';
import AuditReportOrphans from './audit/AuditReportOrphans.vue';
import AuditReportRestore from './audit/AuditReportRestore.vue';

const props = defineProps<{ restoring?: boolean }>();
const emit = defineEmits<{ (e: 'triggerRestore', resetDB?: boolean): void }>();

const {
  auditReport, isAuditing, isPurgingFiles, isPurgingDB, auditStatusMessage,
  runAudit, purgeOrphanFiles, purgeOrphanDBMedia,
} = useAdminAudit();

onMounted(() => {
  if (!auditReport.value && !isAuditing.value) {
    runAudit(false, false);
  }
});
</script>

<template>
  <div class="space-y-4">
    <AuditReportHeader
      :report="auditReport"
      :loading="isAuditing"
      :purging-files="isPurgingFiles"
      :status-message="auditStatusMessage"
      @run-audit="(pf, pdb) => runAudit(pf, pdb)"
    />

    <AuditReportOrphans
      :report="auditReport"
      :purging-files="isPurgingFiles"
      :purging-db="isPurgingDB"
      @purge-files="(paths) => purgeOrphanFiles(paths)"
      @purge-db="(ids) => purgeOrphanDBMedia(ids)"
    />

    <AuditReportRestore
      :restoring="restoring"
      :restore-status="null"
      @trigger-restore="(resetDB) => emit('triggerRestore', resetDB)"
    />
  </div>
</template>
