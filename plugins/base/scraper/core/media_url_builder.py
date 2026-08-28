# plugins/base/scraper/core/media_url_builder.py (SPEC-PLUGIN-001 / 100行以下)
import os, re
from typing import List, Tuple

QUALITY_TAGS = ["orig", "large", "medium", "small", "thumb", "tiny", "1200x1200", "900x900", "4096x4096"]
WAYBACK_PREFIXES = ["https://web.archive.org/web/2id_/", "https://web.archive.org/web/id_/"]


class MediaUrlBuilder:
    """Twitter / X メディアのあらゆる解像度・フォーマット・アーカイブ候補URLを網羅生成"""

    @classmethod
    def clean_base_url(cls, raw_url: str) -> Tuple[str, str, str]:
        if not raw_url: return "", "", "jpg"
        u = raw_url.split("?")[0]
        for sfx in [":orig", ":large", ":medium", ":small", ":thumb", ":tiny", ":900x900", ":1200x1200"]:
            if u.endswith(sfx): u = u[:-len(sfx)]; break
        base_no_ext, ext = os.path.splitext(u)
        fmt = ext.lstrip(".").lower() or "jpg"
        return u, base_no_ext, fmt

    @classmethod
    def build_all_candidates(cls, raw_url: str) -> List[Tuple[str, str]]:
        if not raw_url: return []
        clean_u, base_no_ext, fmt = cls.clean_base_url(raw_url)
        is_video = "video.twimg.com" in raw_url or fmt in ("mp4", "m3u8", "webm")
        cands: List[Tuple[str, str]] = []

        if is_video:
            cands.append((raw_url, "video_raw"))
            cands.append((clean_u, "video_clean"))
            for pfx in WAYBACK_PREFIXES:
                cands.append((f"{pfx}{raw_url}", "wb_video_raw"))
                if clean_u != raw_url: cands.append((f"{pfx}{clean_u}", "wb_video_clean"))
        else:
            # 1. 最高品質 orig (Direct & Wayback)
            for qual in QUALITY_TAGS:
                u_param = f"{base_no_ext}?format={fmt}&name={qual}"
                u_colon = f"{clean_u}:{qual}"
                u_dot = f"{base_no_ext}.{fmt}?name={qual}"
                # 直接候補
                cands.extend([(u_param, qual), (u_colon, qual), (u_dot, qual)])
                # Wayback 候補 (2id_ & id_)
                for pfx in WAYBACK_PREFIXES:
                    cands.append((f"{pfx}{u_param}", f"wb_{qual}"))
                    cands.append((f"{pfx}{u_colon}", f"wb_{qual}"))
                    cands.append((f"{pfx}{u_dot}", f"wb_{qual}"))

            # 2. クリーンなベースURL & 生URL
            cands.append((clean_u, "standard"))
            cands.append((raw_url, "raw"))
            for pfx in WAYBACK_PREFIXES:
                cands.append((f"{pfx}{clean_u}", "wb_standard"))
                cands.append((f"{pfx}{raw_url}", "wb_raw"))

        # 重複排除 (順序維持)
        seen = set()
        deduped: List[Tuple[str, str]] = []
        for url_str, q in cands:
            if url_str and url_str not in seen:
                seen.add(url_str)
                deduped.append((url_str, q))
        return deduped

    @classmethod
    def build_url_list(cls, raw_url: str) -> List[str]:
        return [u for u, _ in cls.build_all_candidates(raw_url)]
