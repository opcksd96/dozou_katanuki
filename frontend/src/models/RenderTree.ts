export type DownloadStatus = 'QUEUED' | 'COMPLETED' | 'DEAD_404' | 'OUTSOURCED' | 'RETAINED';

export interface RenderMetrics {
  replies: number;
  retweets: number;
  likes: number;
  views?: number;
}

export interface RenderAuthor {
  numeric_id: string;
  handle: string;
  display_name: string;
  avatar_url: string;
  bio: string;
  group_name?: string;
  alias_of?: string;
}


export interface RenderMedia {
  id: string;
  type: 'image' | 'video' | 'gif';
  download_status: DownloadStatus;
  failed_reason?: string;
  urls: {
    stream: string;
    image: string;
    thumbnail: string;
    preview?: string;
    vtt?: string;
    original: string;
  };
  width?: number;
  height?: number;
  stash_scene_id?: string;
  stash_image_id?: string;
}

export interface RenderTree {
  id: string;
  conversation_id: string;
  created_at: string;
  content: {
    original: string;
    ja?: string;
    en?: string;
    zh?: string;
  };
  author: RenderAuthor;
  media: RenderMedia[];
  metrics: RenderMetrics;
  source_url: string;
  is_liked: boolean;
  is_pinned: boolean;
  parent_id?: string;
  reply_to_handle?: string;
}

