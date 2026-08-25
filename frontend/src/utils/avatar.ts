// frontend/src/utils/avatar.ts (100行以下)
export function getAvatarInitial(nameOrHandle?: string): string {
  if (!nameOrHandle) return 'A';
  const clean = nameOrHandle.replace(/^@/, '').trim();
  return (clean.charAt(0) || 'A').toUpperCase();
}

export const DEFAULT_HUMAN_AVATAR_SVG = "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%2364748b'><path d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/></svg>";

export function resolveHistoryAvatarUrl(history?: any, platform = 'twitter'): string {
  if (!history) return DEFAULT_HUMAN_AVATAR_SVG;
  if (history.avatar_base64 || history.avatarBase64) {
    return history.avatar_base64 || history.avatarBase64;
  }
  return DEFAULT_HUMAN_AVATAR_SVG;
}

export function resolveAvatarUrl(
  accountOrAuthor?: any,
  histories?: any[],
  platform = 'twitter'
): string {
  if (!accountOrAuthor) return DEFAULT_HUMAN_AVATAR_SVG;

  // ① アカウント自身の Base64
  const b64 = accountOrAuthor.avatar_base64 || accountOrAuthor.avatarBase64;
  if (b64 && b64.startsWith('data:image/')) {
    return b64;
  }

  // ② 世代履歴の Base64
  if (histories && histories.length > 0) {
    const latest = histories[histories.length - 1];
    const hb64 = latest?.avatar_base64 || latest?.avatarBase64;
    if (hb64 && hb64.startsWith('data:image/')) {
      return hb64;
    }
  }

  // ③ 埋め込み履歴の Base64
  const embeddedHistories = accountOrAuthor.profile_history || accountOrAuthor.ProfileHistory;
  if (Array.isArray(embeddedHistories) && embeddedHistories.length > 0) {
    const latest = embeddedHistories[embeddedHistories.length - 1];
    const eb64 = latest?.avatar_base64 || latest?.avatarBase64;
    if (eb64 && eb64.startsWith('data:image/')) {
      return eb64;
    }
  }

  // ④ アカウントの avatar_url が Base64 の場合
  const rawUrl = accountOrAuthor.avatar_url || accountOrAuthor.avatarUrl || '';
  if (rawUrl && rawUrl.startsWith('data:image/')) {
    return rawUrl;
  }

  return DEFAULT_HUMAN_AVATAR_SVG;
}

