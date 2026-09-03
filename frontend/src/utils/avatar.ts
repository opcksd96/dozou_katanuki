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

  // The backend (middleware.ResolveAccountAvatar) is the Single Source of Truth.
  // It guarantees that avatar_url is always a Base64 string or the default SVG.
  // We no longer need makeshift fallback logic here.
  const rawUrl = accountOrAuthor.avatar_url || accountOrAuthor.avatarUrl || '';
  if (rawUrl && rawUrl.startsWith('data:image/')) {
    return rawUrl;
  }

  return DEFAULT_HUMAN_AVATAR_SVG;
}

