// frontend/src/utils/avatar.ts (100行以下)
export function getAvatarInitial(nameOrHandle?: string): string {
  if (!nameOrHandle) return 'A';
  const clean = nameOrHandle.replace(/^@/, '').trim();
  return (clean.charAt(0) || 'A').toUpperCase();
}

export function resolveHistoryAvatarUrl(history?: any, platform = 'twitter'): string {
  if (!history) return '';
  // ① DBに格納された Base64 イメージを最優先
  if (history.avatar_base64) {
    return history.avatar_base64;
  }
  if (history.avatar_virtual_key) {
    return `/avatars/${platform}/${history.avatar_virtual_key}.jpg`;
  }
  return history.avatar_original_url || '';
}

export function resolveAvatarUrl(
  accountOrAuthor?: any,
  histories?: any[],
  platform = 'twitter'
): string {
  if (!accountOrAuthor) return '';

  // ① アカウント自身に Base64 がある場合を最優先
  if (accountOrAuthor.avatar_base64 || accountOrAuthor.avatarBase64) {
    return accountOrAuthor.avatar_base64 || accountOrAuthor.avatarBase64;
  }

  // ② 世代履歴の最新に Base64 がある場合
  if (histories && histories.length > 0) {
    const latest = histories[histories.length - 1];
    if (latest?.avatar_base64) {
      return latest.avatar_base64;
    }
    if (latest?.avatar_virtual_key) {
      return `/avatars/${platform}/${latest.avatar_virtual_key}.jpg`;
    }
  }

  // ③ 埋め込み履歴の最新に Base64 がある場合
  const embeddedHistories = accountOrAuthor.profile_history || accountOrAuthor.ProfileHistory;
  if (Array.isArray(embeddedHistories) && embeddedHistories.length > 0) {
    const latest = embeddedHistories[embeddedHistories.length - 1];
    if (latest?.avatar_base64) {
      return latest.avatar_base64;
    }
    if (latest?.avatar_virtual_key) {
      return `/avatars/${platform}/${latest.avatar_virtual_key}.jpg`;
    }
  }

  // ④ アカウントのアバターURL
  const rawUrl = accountOrAuthor.avatar_url || accountOrAuthor.avatarUrl || '';
  if (rawUrl) {
    return rawUrl;
  }

  return '';
}
