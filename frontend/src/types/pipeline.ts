// frontend/src/types/pipeline.ts (100行以下 - SPEC-PRINCIPLE-001)
export interface PipelineMetrics {
  media_count: number;
  total_media_count: number;
  total_tasks: number;
  registered_tasks: number;
  active_tasks: number;
  pending_tasks: number;
  success_tasks: number;
  failed_tasks: number;
}

export interface CheckpointStatus {
  name: string;
  key: string;
  is_online: boolean;
  active_count: number;
  total_count: number;
  speed_text?: string;
  status_text: string;
  summary_text?: string;
}

export interface PipelineOverview {
  checkpoints: CheckpointStatus[];
  total_media: number;
  completed: number;
  escalated: number;
  outsourced: number;
  retained: number;
  overall_progress: number;
  is_auto_engine_running: boolean;
  metrics?: PipelineMetrics;
}
