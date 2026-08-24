"""
relation_extractor.py
Whitelist対象アカウントのツイートからリプライ・メンション先を抽出するモジュール。
Go Job Orchestrator互換のCLIオプションインターフェースを実装。
"""
import argparse
import os
import re
import sqlite3
from typing import Dict, List, Optional, Set, Tuple

MENTION_PATTERN = re.compile(r"@([a-zA-Z0-9_]{1,15})")


def get_existing_article_ids(conn: sqlite3.Connection) -> Set[str]:
    cursor = conn.cursor()
    cursor.execute("SELECT id FROM articles")
    return {row[0] for row in cursor.fetchall()}


def get_active_whitelist_numeric_ids(
    conn: sqlite3.Connection, target_account: Optional[str] = None
) -> List[str]:
    cursor = conn.cursor()
    if target_account:
        cursor.execute(
            """
            SELECT a.numeric_id FROM accounts a
            JOIN whitelists w ON LOWER(w.value) = LOWER(a.username)
            WHERE LOWER(a.username) = LOWER(?) AND w.is_active = 1
        """,
            (target_account,),
        )
    else:
        cursor.execute("""
            SELECT a.numeric_id FROM accounts a
            JOIN whitelists w ON LOWER(w.value) = LOWER(a.username)
            WHERE w.is_active = 1
        """)
    return [row[0] for row in cursor.fetchall()]


def extract_target_relations(
    db_path: str = "archive.db", target_account: Optional[str] = None
) -> Tuple[Set[str], Set[str], Dict[str, Set[str]]]:
    if not os.path.exists(db_path):
        raise FileNotFoundError(f"Database not found at: {db_path}")

    conn = sqlite3.connect(db_path)
    existing_ids = get_existing_article_ids(conn)
    target_author_ids = get_active_whitelist_numeric_ids(conn, target_account)

    if not target_author_ids:
        conn.close()
        return set(), set(), {}

    placeholders = ",".join(["?"] * len(target_author_ids))
    cursor = conn.cursor()
    cursor.execute(
        f"""
        SELECT reply_to_id, reply_to_handle, full_text 
        FROM articles
        WHERE account_id IN ({placeholders})
    """,
        target_author_ids,
    )

    missing_parent_ids: Set[str] = set()
    related_handles: Set[str] = set()
    user_to_missing_ids: Dict[str, Set[str]] = {}

    for reply_to_id, reply_to_handle, full_text in cursor.fetchall():
        if reply_to_id and reply_to_id not in existing_ids:
            missing_parent_ids.add(reply_to_id)
            if reply_to_handle:
                handle = reply_to_handle.lower().lstrip("@")
                related_handles.add(handle)
                user_to_missing_ids.setdefault(handle, set()).add(reply_to_id)

        if full_text:
            mentions = {m.lower() for m in MENTION_PATTERN.findall(full_text)}
            related_handles.update(mentions)

    conn.close()
    return missing_parent_ids, related_handles, user_to_missing_ids


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Extract missing relations")
    parser.add_argument(
        "--account", "-a", default=None, help="Target account username"
    )
    parser.add_argument(
        "--db",
        "-d",
        default="archive.db",
        help="Path to SQLite3 archive.db (default: archive.db)",
    )
    args = parser.parse_args()

    ids, handles, mapping = extract_target_relations(args.db, args.account)
    target_label = f"@{args.account}" if args.account else "ALL ACTIVE WHITELIST"
    print(
        f"[+] スコープ: {target_label} | 未取得親ID: {len(ids)}件 | 関連垢: {len(handles)}件"
    )
