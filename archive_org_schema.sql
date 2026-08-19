CREATE TABLE scrape_logs (
    log_id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_type TEXT, -- 'wayback', 'live', 'hashtag' etc.
    target TEXT, -- アカウント名やハッシュタグ
    step_name TEXT,
    status TEXT, -- 'success', 'skipped', 'error'
    items_processed INTEGER,
    error_count INTEGER,
    http_404_count INTEGER,
    message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE whitelist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT, -- 'account', 'hashtag', 'mention'
    value TEXT UNIQUE, -- 例: 'yike_luo', '#yike_luo'
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE url_redirects (
    short_url TEXT PRIMARY KEY,
    expanded_url TEXT NOT NULL,
    tweet_id TEXT
);
CREATE TABLE accounts (
    numeric_id TEXT PRIMARY KEY, -- 数値のID（例: 1749477300754878464）
    username TEXT NOT NULL, -- 現在/最新の @ハンドル名
    avatar_local_path TEXT, -- assets以下のバイナリ保存パス
    avatar_base64 TEXT, -- フォールバック用
    custom_header_path TEXT, -- ユーザーカスタマイズ用
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0
);
CREATE TABLE account_profile_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    numeric_id TEXT NOT NULL,
    display_name TEXT,
    description TEXT,
    observed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);
CREATE TABLE tweets (
    tweet_id TEXT PRIMARY KEY,
    numeric_id TEXT NOT NULL, -- accounts.numeric_id
    conversation_id TEXT, -- 会話（スレッド）全体を束ねるコンテナID
    created_at DATETIME NOT NULL,
    full_text TEXT NOT NULL,
    reply_to_tweet_id TEXT, -- 直近の親ツイートID
    is_retweet BOOLEAN DEFAULT 0,
    retweet_target_id TEXT, -- 元ツイートのID
    reply_count INTEGER DEFAULT 0,
    retweet_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    bookmark_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    source_type TEXT,
    is_liked BOOLEAN DEFAULT 0,
    wayback_url TEXT,
    status TEXT,
    FOREIGN KEY(numeric_id) REFERENCES accounts(numeric_id)
);
CREATE TABLE system_raw_tweets (
    tweet_id TEXT PRIMARY KEY,
    raw_content TEXT,
    FOREIGN KEY(tweet_id) REFERENCES tweets(tweet_id) ON DELETE CASCADE
);
CREATE TABLE media (
    media_id TEXT PRIMARY KEY,
    tweet_id TEXT NOT NULL,
    type TEXT NOT NULL, -- 'photo', 'video', 'animated_gif'
    
    -- ダウンローダー連携
    source_platform TEXT, -- 'x', 'wayback', 'manual'
    download_url TEXT, -- 取得元URL
    download_status TEXT DEFAULT 'PENDING', -- PENDING, REQUESTED, MOTRIX_SENDED, THUNDER_SENDED, DOWNLOAD_PROGRESS, SUCCESS, DEAD_404
    
    -- Stashネイティブ連携
    stash_scene_id INTEGER,
    stash_image_id INTEGER,
    
    width INTEGER,
    height INTEGER,
    
    FOREIGN KEY(tweet_id) REFERENCES tweets(tweet_id)
);
CREATE TABLE media_performers (
    media_id TEXT NOT NULL,
    performer_id INTEGER NOT NULL, -- StashのPerformer ID
    PRIMARY KEY (media_id, performer_id),
    FOREIGN KEY(media_id) REFERENCES media(media_id) ON DELETE CASCADE
);
