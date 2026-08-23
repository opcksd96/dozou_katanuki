# scripts/frontend_creator.ps1
# 1. 必要なディレクトリの作成
New-Item -ItemType Directory -Force -Path "frontend/src/models", "frontend/src/mock", "frontend/src/utils", "frontend/src/composables", "frontend/src/components/article", "frontend/src/components/media"

# 2. models/RenderTree.ts (仕様書 2.3節 完全準拠)
@'
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
    original: string;
  };
  width?: number;
  height?: number;
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
}
'@ | Out-File -Encoding utf8 "frontend/src/models/RenderTree.ts"

# 3. mock/sample_render_trees.json (COMPLETED, QUEUED, DEAD_404 検証用)
@'
[
  {
    "id": "1879382757924868404",
    "conversation_id": "1879382757924868404",
    "created_at": "2024-01-15T12:30:00Z",
    "source_url": "https://web.archive.org/web/https://twitter.com/msluo14/status/1879382757924868404",
    "is_liked": false,
    "is_pinned": true,
    "author": {
      "numeric_id": "1001",
      "handle": "msluo14",
      "display_name": "Yike Luo",
      "avatar_url": "https://api.dicebear.com/7.x/bottts/svg?seed=msluo14",
      "bio": "Archival preservation target."
    },
    "content": {
      "original": "Archived snapshot from Wayback Machine. Deep persistence verified.",
      "ja": "Wayback Machineからの魚拓スナップショット。完全オフライン動態保存の検証に成功しました。",
      "en": "Archived snapshot from Wayback Machine. Deep persistence verified.",
      "zh": "来自 Wayback Machine 的归档快照。深度持久化验证成功。"
    },
    "media": [
      {
        "id": "F8wZ1abXYAAY7kL.jpg",
        "type": "image",
        "download_status": "COMPLETED",
        "urls": {
          "stream": "",
          "image": "https://picsum.photos/800/600?random=10",
          "thumbnail": "https://picsum.photos/400/300?random=10",
          "original": "https://web.archive.org/..."
        }
      },
      {
        "id": "sample_video_clip.mp4",
        "type": "video",
        "download_status": "OUTSOURCED",
        "failed_reason": "外部Motrix(aria2)にて並列ダウンロード中...",
        "urls": {
          "stream": "",
          "image": "",
          "thumbnail": "",
          "original": "https://web.archive.org/..."
        }
      }
    ],
    "metrics": {
      "replies": 12,
      "retweets": 45,
      "likes": 128,
      "views": 3200
    }
  },
  {
    "id": "1879382757924868405",
    "conversation_id": "1879382757924868405",
    "created_at": "2024-01-16T08:15:00Z",
    "source_url": "https://web.archive.org/...",
    "is_liked": true,
    "is_pinned": false,
    "author": {
      "numeric_id": "1001",
      "handle": "msluo14",
      "display_name": "Yike Luo",
      "avatar_url": "https://api.dicebear.com/7.x/bottts/svg?seed=msluo14",
      "bio": "Archival preservation target."
    },
    "content": {
      "original": "Testing dead link SVG filler fallback and animated GIF.",
      "ja": "404消失メディア用SVGフィラーとGIFアニメーションの描画フォールバック検証。"
    },
    "media": [
      {
        "id": "dead_asset.jpg",
        "type": "image",
        "download_status": "DEAD_404",
        "failed_reason": "Wayback CDN 404: 元サーバーから消滅しています",
        "urls": { "stream": "", "image": "", "thumbnail": "", "original": "" }
      },
      {
        "id": "reaction.gif",
        "type": "gif",
        "download_status": "QUEUED",
        "urls": { "stream": "", "image": "", "thumbnail": "", "original": "" }
      }
    ],
    "metrics": {
      "replies": 3,
      "retweets": 8,
      "likes": 56
    }
  }
]
'@ | Out-File -Encoding utf8 "frontend/src/mock/sample_render_trees.json"
