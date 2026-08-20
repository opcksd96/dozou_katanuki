// plugins/twitter/skin/controller.js (SPEC-PLUGIN-001 / 100行以下)

export default {
  init(ctx) {
    this.ctx = ctx;
    console.log('[Skin:Twitter] Initialized layout controller for platform:', ctx?.platform || 'twitter');
  },

  onMount(containerElement) {
    if (!containerElement) return;
    containerElement.classList.add('skin-twitter-active');
  },

  onUnmount() {
    console.log('[Skin:Twitter] Controller unmounted');
  },

  handleItemClick(item, event) {
    if (!item) return;
    const target = event?.target;
    if (target && (target.closest('button') || target.closest('a') || target.closest('.media-click-target'))) {
      return;
    }
    if (this.ctx?.router?.push) {
      this.ctx.router.push(`/${item.author?.handle || 'user'}/status/${item.id}`);
    }
  },

  handleMediaClick(media, index) {
    console.log('[Skin:Twitter] Media clicked index:', index, media);
  },

  actions: {
    async copyLink(item) {
      const url = item.source_url || `https://twitter.com/${item.author?.handle}/status/${item.id}`;
      try {
        await navigator.clipboard.writeText(url);
        return true;
      } catch {
        return false;
      }
    },
    async toggleLike(item) {
      console.log('[Skin:Twitter] Toggling like for:', item.id);
    }
  }
};
