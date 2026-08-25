// frontend/src/utils/avatar.ts (100行以下)
export function getAvatarInitial(nameOrHandle?: string): string {
  if (!nameOrHandle) return 'A';
  const clean = nameOrHandle.replace(/^@/, '').trim();
  return (clean.charAt(0) || 'A').toUpperCase();
}

export function resolveHistoryAvatarUrl(history?: any, platform = 'twitter'): string {
  if (!history) return '';
  if (history.avatar_base64 || history.avatarBase64) {
    return history.avatar_base64 || history.avatarBase64;
  }
  if (history.avatar_original_url || history.avatarOriginalUrl) {
    return history.avatar_original_url || history.avatarOriginalUrl;
  }
  return '';
}

export function resolveAvatarUrl(
  accountOrAuthor?: any,
  histories?: any[],
  platform = 'twitter'
): string {
  if (!accountOrAuthor) return '';

  // ① アカウント自身の Base64
  if (accountOrAuthor.avatar_base64 || accountOrAuthor.avatarBase64) {
    return accountOrAuthor.avatar_base64 || accountOrAuthor.avatarBase64;
  }

  // ② 世代履歴の Base64 / 原本URL
  if (histories && histories.length > 0) {
    const latest = histories[histories.length - 1];
    if (latest?.avatar_base64 || latest?.avatarBase64) {
      return latest.avatar_base64 || latest.avatarBase64;
    }
    if (latest?.avatar_original_url || latest?.avatarOriginalUrl) {
      return latest.avatar_original_url || latest.avatarOriginalUrl;
    }
  }

  // ③ 埋め込み履歴の Base64 / 原本URL
  const embeddedHistories = accountOrAuthor.profile_history || accountOrAuthor.ProfileHistory;
  if (Array.isArray(embeddedHistories) && embeddedHistories.length > 0) {
    const latest = embeddedHistories[embeddedHistories.length - 1];
    if (latest?.avatar_base64 || latest?.avatarBase64) {
      return latest.avatar_base64 || latest.avatarBase64;
    }
    if (latest?.avatar_original_url || latest?.avatarOriginalUrl) {
      return latest.avatar_original_url || latest.avatarOriginalUrl;
    }
  }

  // ④ アカウントの avatar_url
  const rawUrl = accountOrAuthor.avatar_url || accountOrAuthor.avatarUrl || '';
  if (rawUrl) {
    return rawUrl;
  }

  return '';
}

