CREATE TABLE IF NOT EXISTS `accounts` (`numeric_id` text,`username` text NOT NULL,`display_name` text NOT NULL,`avatar_url` text NOT NULL,`avatar_base64` text,`description` text,`group_name` text DEFAULT '',`alias_of` text DEFAULT '',`updated_at` datetime NOT NULL,PRIMARY KEY (`numeric_id`));
CREATE INDEX IF NOT EXISTS `idx_accounts_username` ON `accounts`(`username`);
CREATE TABLE IF NOT EXISTS `account_profile_histories` (`id` integer PRIMARY KEY AUTOINCREMENT,`account_id` text NOT NULL,`display_name` text NOT NULL,`avatar_original_url` text NOT NULL,`avatar_seq` integer NOT NULL,`avatar_virtual_key` text NOT NULL,`observed_at` datetime NOT NULL,CONSTRAINT `fk_accounts_profile_history` FOREIGN KEY (`account_id`) REFERENCES `accounts`(`numeric_id`));
CREATE INDEX IF NOT EXISTS `idx_history_lookup` ON `account_profile_histories`(`account_id`, `avatar_seq` DESC);
CREATE TABLE IF NOT EXISTS `articles` (`id` text,`account_id` text NOT NULL,`conversation_id` text NOT NULL,`reply_to_id` text,`reply_to_handle` text,`created_at` datetime NOT NULL,`full_text` text NOT NULL,`lang` text NOT NULL DEFAULT "ja",`full_text_ja` text,`full_text_en` text,`full_text_zh` text,`via` text NOT NULL,`is_repost` boolean NOT NULL DEFAULT false,`is_liked` boolean NOT NULL DEFAULT false,`wayback_url` text NOT NULL,PRIMARY KEY (`id`),CONSTRAINT `fk_accounts_articles` FOREIGN KEY (`account_id`) REFERENCES `accounts`(`numeric_id`));
CREATE INDEX IF NOT EXISTS `idx_articles_is_liked_created` ON `articles`(`is_liked`, `created_at` DESC) WHERE `is_liked` = 1;
CREATE INDEX IF NOT EXISTS `idx_articles_account_created` ON `articles`(`account_id`, `created_at` DESC);
CREATE INDEX IF NOT EXISTS `idx_articles_conversation` ON `articles`(`conversation_id`, `created_at` ASC);
CREATE INDEX IF NOT EXISTS `idx_articles_reply_to` ON `articles`(`reply_to_id`);
CREATE INDEX IF NOT EXISTS `idx_articles_created_at` ON `articles`(`created_at` DESC);
CREATE TABLE IF NOT EXISTS `media` (`media_id` text,`article_id` text NOT NULL,`type` text NOT NULL,`download_url` text NOT NULL,`width` integer NOT NULL,`height` integer NOT NULL,`download_status` text NOT NULL DEFAULT "QUEUED",`failed_reason` text,`stash_scene_id` text,`stash_image_id` text,PRIMARY KEY (`media_id`),CONSTRAINT `fk_articles_media` FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`));
CREATE INDEX IF NOT EXISTS `idx_media_article` ON `media`(`article_id`);
CREATE UNIQUE INDEX IF NOT EXISTS `idx_media_stash_scene` ON `media`(`stash_scene_id`) WHERE `stash_scene_id` IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS `idx_media_stash_image` ON `media`(`stash_image_id`) WHERE `stash_image_id` IS NOT NULL;
CREATE TABLE IF NOT EXISTS `url_redirects` (`short_url` text,`expanded_url` text NOT NULL,`article_id` text NOT NULL,PRIMARY KEY (`short_url`));
CREATE INDEX IF NOT EXISTS `idx_url_redirects_article_id` ON `url_redirects`(`article_id`);
CREATE TABLE IF NOT EXISTS `whitelists` (`id` integer PRIMARY KEY AUTOINCREMENT,`type` text NOT NULL,`value` text NOT NULL,`group_name` text DEFAULT '',`alias_of` text DEFAULT '',`is_active` boolean NOT NULL DEFAULT true);
CREATE UNIQUE INDEX IF NOT EXISTS `idx_whitelists_value` ON `whitelists`(`value`);

