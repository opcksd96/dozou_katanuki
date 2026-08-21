// frontend/src/composables/useAdmin.ts (100行以下 - Facade)
import { useAdminJob } from './admin/useAdminJob';
import { useAdminConfig } from './admin/useAdminConfig';
import { useAdminWhitelist } from './admin/useAdminWhitelist';
import { useAdminDatabase } from './admin/useAdminDatabase';
import { useAdminAudit } from './admin/useAdminAudit';
import { useAdminSkin } from './admin/useAdminSkin';
import { useAdminBroadcast } from './admin/useAdminBroadcast';

export * from './admin/useAdminJob';
export * from './admin/useAdminConfig';
export * from './admin/useAdminWhitelist';
export * from './admin/useAdminDatabase';
export * from './admin/useAdminAudit';
export * from './admin/useAdminSkin';
export * from './admin/useAdminBroadcast';

export function useAdmin() {
  const job = useAdminJob();
  const config = useAdminConfig();
  const whitelist = useAdminWhitelist();
  const database = useAdminDatabase();
  const audit = useAdminAudit();
  const skin = useAdminSkin();
  const broadcast = useAdminBroadcast();

  return {
    ...job,
    ...config,
    ...whitelist,
    ...database,
    ...audit,
    ...skin,
    ...broadcast,
  };
}
