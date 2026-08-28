# plugins/base/scraper/core/source_orchestrator.py (SPEC-PLUGIN-001 / 100行以下)
from typing import Any, Callable, Dict, List, Optional
from .base_source import BaseSource


class SourceOrchestrator:
    """マルチソース・スクレイパー統括オーケストレーター（剪定なし完全収集）"""

    def __init__(self):
        self._sources: Dict[str, BaseSource] = {}

    def register(self, source: BaseSource) -> None:
        self._sources[source.name.lower()] = source

    def get_source(self, name: str) -> Optional[BaseSource]:
        return self._sources.get(name.lower())

    def list_sources(self) -> List[str]:
        return sorted(self._sources.keys(), key=lambda k: self._sources[k].priority)

    def collect(
        self, account: str, limit: int = 0, source_filter: str = "all", log_fn: Optional[Callable[[str], None]] = None
    ) -> List[Dict[str, Any]]:
        def _log(m: str):
            if log_fn: log_fn(m)
            print(f"[Orchestrator] {m}", flush=True)

        target_names = [source_filter.lower()] if source_filter.lower() in self._sources else self.list_sources()
        _log(f"Starting multi-source collection for @{account} (sources: {', '.join(target_names)}, limit={limit})")

        all_records: List[Dict[str, Any]] = []

        for name in target_names:
            src = self._sources.get(name)
            if not src or not src.is_available():
                _log(f"Source [{name}] is unavailable or not registered. Skipping.")
                continue

            try:
                _log(f"Querying source [{src.name}] (priority={src.priority})...")
                records = src.fetch_account(account, limit=limit, log_fn=log_fn)
                _log(f"Source [{src.name}] returned {len(records)} records.")
                all_records.extend(records)
                if limit > 0 and len(all_records) >= limit:
                    _log(f"Reached collection limit ({len(all_records)}/{limit}). Stopping chain.")
                    return all_records[:limit]
            except Exception as e:
                _log(f"Error executing source [{name}]: {type(e).__name__}: {e}")

        _log(f"Collection complete. Total retrieved records: {len(all_records)}")
        return all_records
