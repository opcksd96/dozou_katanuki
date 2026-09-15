// frontend/src/types/pipeline_media_tasks.ts (100行以下 - SPEC-PRINCIPLE-001)

export interface ActiveThunderTask {
  file_name: string;
  target: string;
  status: string;
  progress_text: string;
  speed_text: string;
  detail_text: string;
  task_id: number;
  task_type: number;
  task_status_code: number;
  error_code: number;
  url: string;
  save_path: string;
  file_size: number;
  download_size: number;
  download_speed: number;
  vip_speed: number;
  create_time: number;
  completion_time: number;
  src_total: number;
  src_using: number;
  peers: number;
  cid: string;
  gcid: string;
  origin: string;
  raw_json?: string;
  raw_text?: string;
}

export interface CandidateTask {
  id: string;
  media_id: string;
  article_id?: string;
  resolution_type: string;
  url: string;
  file_name: string;
  status: string;
  summary_size?: string;
  error_reason?: string;
  thunder_task_id?: number;
  file_size?: number;
  download_size?: number;
  save_path?: string;
  task_status_code?: number;
  error_code?: number;
  detail_text?: string;
  last_attempt_at?: string;
  dispatched_at?: string;
  completed_at?: string;
  reaped_at?: string;
  raw_json?: string;
  created_at: string;
  updated_at?: string;
  cdp?: ActiveThunderTask;
}

export interface MediaWithTasks {
  media_id: string;
  article_id: string;
  original_url: string;
  type?: string;
  width?: number;
  height?: number;
  download_status: string;
  failed_reason?: string;
  account_id?: string;
  media_quality?: string;
  candidate_tasks: CandidateTask[];
}
