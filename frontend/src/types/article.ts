export interface Author {
  platform: string;
  username: string;
  displayName: string;
  avatarUrl: string; // 解決済み相対パス (/assets/twitter/username_avatar_001.jpg)
}

export interface MediaItem {
  mediaId: string; // URL BaseName (例: F8wZ1abXYAAY7kL.jpg)
  type: 'image' | 'video';
  url: string; // /stash-proxy/scenes/xxx または /media-local/xxx
  thumbnailUrl?: string;
  width?: number;
  height?: number;
}

export interface TranslatedText {
  original: string;
  ja: string;
  en: string;
  zh: string;
}

export interface ArticleRenderTree {
  id: string; // post_id
  platform: string;
  createdAt: string;
  originalUrl: string;
  author: Author;
  content: TranslatedText;
  media: MediaItem[];
  stats: {
    likes: number;
    reposts: number;
  };
}
