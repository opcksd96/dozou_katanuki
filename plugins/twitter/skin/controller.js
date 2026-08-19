// plugins/twitter/skin/controller.js (100行以下)
export function initTwitterSkin(context) {
  console.log('[Skin:Twitter] Initialized layout controller for', context?.platform || 'twitter');
}

export function onPostRendered(postElement) {
  if (!postElement) return;
  postElement.classList.add('twitter-skin-card');
}
