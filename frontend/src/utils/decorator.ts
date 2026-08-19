// frontend/src/utils/decorator.ts (100行以下・純粋関数)

/**
 * XSSを防ぐためのHTML特殊文字エスケープ
 */
export function escapeHtml(str: string): string {
  if (!str) return '';
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

/**
 * 本文テキスト内の URL / ハッシュタグ / メンション / 改行 を安全にDOMリンク化します
 * (SPEC-MIDDLEWARE-001-2 / part3_03_2_data_decorator.md 準拠)
 */
export function decorateText(rawText?: string | null, platform: string = 'twitter'): string {
  if (!rawText) return '';

  // 既に一部HTMLタグが含まれている場合の安全なハンドリング
  const hasHtmlTags = /<[a-z][\s\S]*>/i.test(rawText);
  if (hasHtmlTags) {
    return rawText;
  }

  const safe = escapeHtml(rawText);

  // 1. URLのリンク化 (https?://...)
  const urlRegex = /(https?:\/\/[^\s<]+)/g;
  let decorated = safe.replace(urlRegex, (url) => {
    return `<a href="${url}" target="_blank" rel="noopener noreferrer" class="external-link text-blue-400 hover:text-blue-300 hover:underline inline-flex items-center gap-0.5" data-url="${url}">${url}</a>`;
  });

  // 2. ハッシュタグのリンク化 (#タグ名 / 英数・CJK・アンダースコア対応)
  const hashtagRegex = /(^|[^\w#&;])#([a-zA-Z0-9_\p{L}\p{N}]+)/gu;
  decorated = decorated.replace(hashtagRegex, (_match, prefix, tag) => {
    return `${prefix}<a href="/${platform}/search?q=${encodeURIComponent(tag)}" class="hashtag-link text-sky-400 hover:text-sky-300 hover:underline font-medium cursor-pointer" data-tag="${tag}">#${tag}</a>`;
  });

  // 3. メンションのリンク化 (@ユーザー名)
  const mentionRegex = /(^|[^\w@&;])@([a-zA-Z0-9_]{1,30})/g;
  decorated = decorated.replace(mentionRegex, (_match, prefix, handle) => {
    return `${prefix}<a href="/${platform}/${handle}" class="mention-link text-sky-400 hover:text-sky-300 hover:underline font-medium cursor-pointer" data-mention="${handle}">@${handle}</a>`;
  });

  // 4. 改行コードのDOM展開 (\n -> <br/>)
  decorated = decorated.replace(/\r?\n/g, '<br/>');

  return decorated;
}

