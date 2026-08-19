// frontend/src/utils/formatters.ts (100行以下・純粋関数)

// 統計数値のフォーマット (1.2K, 3.4M 等)
export function formatStatNumber(val?: number | string | null): string {
  if (val === undefined || val === null || val === '') {
    return '0';
  }
  const num = typeof val === 'string' ? parseFloat(val) : val;
  if (isNaN(num)) {
    return '0';
  }
  if (num >= 1_000_000) {
    return (num / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M';
  }
  if (num >= 1_000) {
    return (num / 1_000).toFixed(1).replace(/\.0$/, '') + 'K';
  }
  return num.toLocaleString();
}

// 投稿日時のフォーマット (YYYY/MM/DD HH:mm)
export function formatDate(dateStr?: string | null): string {
  if (!dateStr) return '';
  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return '';

  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');

  return `${year}/${month}/${day} ${hours}:${minutes}`;
}
