// plugins/base/skin/controller.js (SPEC-PLUGIN-001 / Base Controller Skeleton / 100行以下)
export class BaseSkinController {
  constructor(options = {}) {
    this.platform = options.platform || 'base';
    this.container = options.container || null;
    this.eventBus = options.eventBus || null;
  }

  init() {
    this.bindEvents();
  }

  bindEvents() {
    if (!this.container) return;
    this.container.addEventListener('click', (e) => {
      const card = e.target.closest('.base-timeline-card');
      if (!card) return;

      const media = e.target.closest('.base-media-container img, .base-media-container video');
      if (media) {
        this.onMediaClick(media, card.dataset.articleId);
        return;
      }

      const likeBtn = e.target.closest('.action-like');
      if (likeBtn) {
        this.onLikeClick(card.dataset.articleId);
      }
    });
  }

  onMediaClick(mediaEl, articleId) {
    if (this.eventBus) {
      this.eventBus.emit('media:open', { articleId, src: mediaEl.src || mediaEl.currentSrc });
    }
  }

  onLikeClick(articleId) {
    if (this.eventBus) {
      this.eventBus.emit('article:toggle-like', { articleId });
    }
  }

  destroy() {
    this.container = null;
    this.eventBus = null;
  }
}
