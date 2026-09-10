# plugins/base/scraper/core/translator.py (SPEC-PLUGIN-001 / 100行以下)
import os, time, requests, re
from typing import Dict, List, Optional


class Translator:
    """DeepL / Google Translate API による 3大主要言語一括翻訳エンジン"""
    def __init__(self, provider: Optional[str] = None, delay_sec: float = 0.05):
        self.deepl_key = os.environ.get("DEEPL_API_KEY", "")
        self.google_key = os.environ.get("GOOGLE_TRANSLATE_API_KEY", "")
        self.provider = provider or ("deepl" if self.deepl_key else ("google" if self.google_key else "none"))
        self.delay_sec = delay_sec
        self.session = requests.Session()

    def detect_lang(self, text: str) -> str:
        if not text: return "ja"
        clean = re.sub(r"https?://\S+|@\w+|#\w+", "", text).strip()
        s = clean if clean else text
        ja_kana = len(re.findall(r"[\u3040-\u309f\u30a0-\u30ff]", s))
        cjk_kanji = len(re.findall(r"[\u4e00-\u9fff]", s))
        latin_words = len(re.findall(r"[a-zA-Z]+", s))
        hangul = len(re.findall(r"[\uac00-\ud7a3]", s))
        if hangul > 0 and hangul >= ja_kana and hangul >= latin_words: return "ko"
        if ja_kana > 0:
            if latin_words > 10 and ja_kana <= 2: return "en"
            return "ja"
        if cjk_kanji > 0:
            return "en" if latin_words > cjk_kanji * 2 else "zh"
        if latin_words > 0: return "en"
        return "ja"

    def determine_targets(self, source_lang: str, supported: Optional[List[str]] = None) -> List[str]:
        pool = supported or ["ja", "en", "zh"]
        return [lang for lang in pool if lang != source_lang]

    def _call_deepl(self, text: str, target: str) -> Optional[str]:
        target_map = {"ja": "JA", "en": "EN-US", "zh": "ZH", "ko": "KO"}
        t_lang = target_map.get(target, target.upper())
        url = "https://api-free.deepl.com/v2/translate" if self.deepl_key.endswith(":fx") else "https://api.deepl.com/v2/translate"
        for retry in range(3):
            try:
                r = self.session.post(url, headers={"Authorization": f"DeepL-Auth-Key {self.deepl_key}"},
                                      data={"text": text, "target_lang": t_lang}, timeout=5.0)
                if r.status_code == 200:
                    return r.json().get("translations", [{}])[0].get("text")
                if r.status_code == 429: time.sleep(1.0 * (2 ** retry))
            except Exception: pass
        return None

    def _call_google(self, text: str, target: str) -> Optional[str]:
        target_map = {"ja": "ja", "en": "en", "zh": "zh-CN", "ko": "ko"}
        t_lang = target_map.get(target, target)
        url = f"https://translation.googleapis.com/language/translate/v2?key={self.google_key}"
        for retry in range(3):
            try:
                r = self.session.post(url, json={"q": text, "target": t_lang, "format": "text"}, timeout=5.0)
                if r.status_code == 200:
                    return r.json().get("data", {}).get("translations", [{}])[0].get("translatedText")
                if r.status_code == 429: time.sleep(1.0 * (2 ** retry))
            except Exception: pass
        return None

    def translate(self, text: str, target_lang: str, source_lang: Optional[str] = None) -> Optional[str]:
        if not text: return None
        if target_lang == source_lang: return text
        res = None
        if self.provider == "deepl" and self.deepl_key:
            res = self._call_deepl(text, target_lang)
        elif self.provider == "google" and self.google_key:
            res = self._call_google(text, target_lang)
        if self.delay_sec > 0: time.sleep(self.delay_sec)
        return res

    def translate_article(self, full_text: str, source_lang: Optional[str] = None,
                          targets: Optional[List[str]] = None) -> Dict[str, Optional[str]]:
        if not full_text or not full_text.strip():
            return {"lang": "ja", "ja": None, "en": None, "zh": None}
        src = source_lang or self.detect_lang(full_text)
        target_langs = targets if targets is not None else self.determine_targets(src)
        res: Dict[str, Optional[str]] = {"lang": src, "ja": None, "en": None, "zh": None}
        if src in res: res[src] = full_text
        for target in target_langs:
            if target != src:
                t = self.translate(full_text, target, src)
                if t: res[target] = t
        return res
