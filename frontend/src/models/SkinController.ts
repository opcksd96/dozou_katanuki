// frontend/src/models/SkinController.ts (SPEC-FRONTEND-001 / SPEC-PLUGIN-001)

import type { RenderTree, RenderMedia } from './RenderTree';

export interface SkinContext {
  platform: string;
  router?: {
    push: (path: string) => void;
  };
  api?: {
    fetchRelated?: (id: string) => Promise<RenderTree[]>;
    toggleLike?: (id: string) => Promise<boolean>;
  };
  showToast?: (msg: string) => void;
  state?: any;
}

export interface SkinController {
  init(ctx: SkinContext): void;
  onMount?(containerElement: HTMLElement): void;
  onUnmount?(): void;
  handleItemClick?(item: RenderTree, event: Event): void;
  handleMediaClick?(media: RenderMedia, index: number): void;
  actions?: Record<string, (item: RenderTree, ...args: any[]) => Promise<any> | void>;
}
