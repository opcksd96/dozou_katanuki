# 基礎アーキテクチャに関する詳細レポート

## 概要

| 項目 | 内容 |
|------|------|
| **ドキュメントID** | ARCH-001 |
| **バージョン** | 2.0 |
| **作成日** | 2026年8月16日 |
| **最終更新** | 2026年8月16日 |
| **作成者** | Kilo Assistant |
| **ステータス** | ドラフト（レビュー待ち） |
| **対象プロジェクト** | x_timeline_app |

本ドキュメントは、x_timeline_app における **MVVM（モデル・ビュー・ビュー・モデル）モデル** に基づくアーキテクチャの実装方針、各オペレーション層の責務・権限定義、およびユーザー種別ごとのアクセス権限精査結果をまとめたものです。

---

## 詳細内容

### 1. MVVM モデルに基づくアーキテクチャ概要

#### 1.1 採用理由
- **データバインディングの自動化**: ViewとModel間の状態同期を自動化し、UI更新の効率化を実現
- **関心の分離**: Presentation（View）、ビジネスロジック（ViewModel）、データアクセス（Model）を明確に分離
- **テスト容易性**: ViewModelのロジックをUIから独立させて単体テストを可能に
- **リユーザビリティ**: 同じViewModelを複数のViewで共有可能

#### 1.2 レイヤー構成図（MVVMベース）

```
┌─────────────────────────────────────────────────────────────────┐
│                           View (View)                           │
│  - UIコンポーネント (Vue.js テンプレート)                       │
│  - ユーザーインタラクションハンドリング                        │
│  - データバインディング（v-model, v-bind等）                    │
└─────────────────────┬───────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────────┐
│                       ViewModel (ViewModel)                     │
│  - ビジネスロジックのカプセル化                                  │
│  - 状態管理と変更通知                                            │
│  - Modelとのデータ同期                                           │
│  - コマンド実行（ユーザーアクションのハンドリング）             │
└─────────────────────┬───────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────────┐
│                           Model (Model)                         │
│  - エンティティ定義 (Tweet, Account, Media)                     │
│  - ドメインロジックとバリデーション                             │
│  - データ永続化インターフェース                                 │
│  - 外部サービス連携（Twitter, Wayback等）                       │
└─────────────────────┬───────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────────┐
│                   Infrastructure Layer (永続化層)               │
│  - データベースアクセス (SQLite3 + GORM)                        │
│  - 外部API連携実装                                                │
│  - キャッシュ機構                                                 │
│  - ファイルシステム操作                                           │
└─────────────────────────────────────────────────────────────────┘
```

### 2. 各オペレーション層の実装方針（MVVM視点）

#### 2.1 View (ビュー・プレゼンテーション層)
- **実装技術**: Vue.js 3 (Composition API + テンプレート構文)
- **責務**:
  - ユーザーインターフェースの描画（HTMLテンプレート）
  - ユーザー入力の受け取り（フォーム、クリック、タッチなど）
  - ViewModelとのデータバインディング（v-model, v-bind, v-on）
  - ローディング状態とエラー表示の視覚化
  - アニメーションと遷移効果の実装
- **禁止事項**:
  - ビジネスロジックの直接実装
  - Model層への直接アクセス
  - 永続化処理の実装
  - 複雑な状態管理ロジック（ViewModelに委譲）

#### 2.2 ViewModel (ビュー・モデル層)
- **実装技術**: TypeScript サービスクラス + Vue.js Composables
- **責務**:
  - ビジネスロジックの実装と状態管理
  - Viewへの状態変更通知（リアクティブプロパティ）
  - Model層へのデータ取得・更新要求
  - ユーザーアクションのコマンドとしてのハンドリング
  - エラーハンドリングとバリデーションロジック
  - 外部サービス連携の調整（非同期処理、タイムアウト）
  - ローディング状態とエラー状態の管理
- **許可される操作**:
  - Model層のインターフェース呼び出し
  - Infrastructure層へのアクセス（抽象化経由）
  - 他のViewModelとの連携
  - UI状態の管理（フラグ、カウンタなど）

#### 2.3 Model (モデル層)
- **実装技術**: TypeScript クラス/インターフェース、Go言語エンティティ（バックエンド）
- **責務**:
  - ビジネスエンティティの定義（Tweet, Account, Media, UserSessionなど）
  - ドメインルールのカプセル化（バリデーション、状態遷移、不変性）
  - ドメインサービスによる複雑なビジネスロジックの実装
  - データ永続化のためのインターフェース定義
  - 外部サービス連携のためのポート定義
- **禁止事項**:
  - UI固有のコードへの依存
  - View固有の状態管理
  - フレームワーク固有の実装への依存
  - 永続化詳細への踏み込み（SQLクエリなど）

#### 2.4 Infrastructure Layer (インフラストラクチャ層)
- **実装技術**: 
  - フロントエンド: Axios HTTPクライアント、localStorage、IndexedDB
  - バックエンド: Go言語 + GORM + SQLite3
- **責務**:
  - データベース接続とクエリ実行（Model層のインターフェース実装）
  - 外部APIとの実際の通信実装（Model層のポート実装）
  - ファイルシステム操作（メディアダウンロード保存）
  - キャッシュ実装（Redis、メモリキャッシュなど）
  - ロギングインフラストラクチャ
  - トランザクション管理
- **許可される操作**:
  - データベースCRUD操作
  - 外部HTTPリクエストの送受信
  - ファイルI/O操作
  - キャッシュ読み書き
  - トランザクション境界管理

### 3. MVVM特有のデータフローとバインディング

#### 3.1 双方向データバインディングの流れ
```
[ユーザー操作] 
        ↓ (View → ViewModel: イベントハンドリング)
[ViewModel: 状態更新・ビジネスロジック実行]
        ↓ (ViewModel → Model: データ更新要求)
[Model: ドメインロジック実行・永続化]
        ↓ (Model → ViewModel: 結果通知)
[ViewModel → View: 状態変更通知 → UI自動更新]
```

#### 3.2 リアクティブストリームの実装
- **Vue.js 3のReactivity System**: ref()、reactive()、computed() を使用
- **非同期処理の扱い**: async/await + try/catch パターン
- **状態の不変性**: 必要に応じてimmutableデータ構造を採用
- **変更検知の最適化**:  shallowReactive や custom getters/setter の活用

#### 3.3 コマンドパターンの実装
```typescript
// ViewModel例: コマンドとしてのユーザーアクション
class TimelineViewModel {
  private readonly tweetService: TweetService;
  
  // コマンドプロパティ
  readonly loadTimelineCommand = async (accountId: string) => {
    this.isLoading = true;
    try {
      const tweets = await this.tweetService.getTimeline(accountId);
      this.timeline = tweets;
    } catch (error) {
      this.error = error.message;
    } finally {
      this.isLoading = false;
    }
  };
  
  // 状態プロパティ（Viewにバインディング）
  isLoading = ref(false);
  error = ref<string | null>(null);
  timeline = ref<Tweet[]>([]);
}
```

### 4. ユーザー種別ごとの権限内容（MVVM視点）

#### 4.1 一般ユーザー (General User)
| レイヤー | 許可される操作 | 禁止される操作 |
|----------|----------------|----------------|
| View | UI操作、データ表示、イベント発行 | - |
| ViewModel | タイムライン取得コマンド実行、プロフィール編集コマンド | システム設定変更コマンド、バッチジョブ実行コマンド |
| Model | 自身のエンティティ参照・更新 | 他ユーザーエンティティの直接操作、システム設定変更 |
| Infrastructure | キャッシュ読み取り（自分のデータに限る） | データベース書き込み、システムファイルアクセス |

#### 4.2 管理者ユーザー (Administrator)
| レイヤー | 許可される操作 | 禁止される操作 |
|----------|----------------|----------------|
| View | 管理者UIアクセス、システム監視画面表示 | - |
| ViewModel | システム設定変更コマンド実行、バッチジョブ実行コマンド、ユーザー管理コマンド | セキュリティポリシー変更コマンド（特別権限必要） |
| Model | 全エンティティ参照・更新、システム設定管理 | 暗号鍵の直接アクセス、監査ログ改ざん |
| Infrastructure | システムバックアップ、リストア、パフォーマンスチューニング | セキュリティ設定の直接変更、監査ログ削除 |

#### 4.3 外部連携システム (External System)
| レイヤー | 許可される操作 | 禁止される操作 |
|----------|----------------|----------------|
| View | - | すべてのUI操作 |
| ViewModel | APIエンドポイント経由のデータ送受信コマンド（認証済み） | 直接データベースアクセス、管理者機能コマンド |
| Model | 指定されたエンティティ参照・更新（APIスコープに基づく） | システム設定変更、他の外部システムへのアクセス |
| Infrastructure | データベース参照（読み取り専用レプリカ）、ファイルアップロード | データベース書き込み（特別権限必要）、システム設定変更 |

### 5. セキュリティと権限管理の実装方針（MVVM視点）

#### 5.1 認証・認可フロー
1. **認証**: JWTトークンまたはセッションクッキーによるユーザー識別（ViewModel層で処理）
2. **認可**: ユーザー役割（ROLE_USER、ROLE_ADMIN）に基づくポリシー評価（ViewModel層）
3. **スコープベースアクセス**: 外部システム連携時はOAuthスコープに基づく制限付与（Model層）
4. **監査ログ**: すべてのデータ変更操作をログ記録（Model層またはInfrastructure層）
5. **ViewModelガード**: コマンド実行前に認可チェックを実装

#### 5.2 権限チェックの実装位置
- **ViewModel層でのプレチェック**: ユーザーアクションハンドリングの開始時に権限確認
- **Model層でのビジネスルールチェック**: エンティティ更新時のドメインルール検証
- **Infrastructure層での最後の砦**: データベースレベルの制約とROW LEVEL SECURITY（可能な場合）

### 6. パフォーマンスとスケーラビリティへの配慮（MVVM視点）

#### 6.1 Viewのパフォーマンス最適化
- **v-once ディレクティブ**: 静的コンテンツの再評価を防止
- **キー属性による効率的なリスト更新**: v-for にユニークキーを指定
- **コンポーネントの遅延ロード**: 動的インポートによるコードスプリッティング
- **計算プロパティの活用**: 高頻度アクセスデータのキャッシュ相当

#### 6.2 ViewModelの効率化
- **ストリームベースの状態管理**: RxJS または Vue の reactive streams 利用
- **メモ化の適用**: 高コストな計算結果のキャッシュ
- **バッチ更新**: 関連する状態変更をまとめて適用
- **サブスクリプション管理**: 不要なウォッチャーのクリア防止

#### 6.3 ModelとInfrastructureの最適化
- **クエリ最適化**: インデックス設計とEXPLAIN分析の実施
- **バッチ処理の活用**: 大量データ操作時のバッチインサート/アップデート
- **非同期処理の分離**: 重い操作はバックグラウンドジョブまたはワーカーに委譲
- **接続プーリング**: データベース接続の効率的な再利用

### 7. 移行と互換性の考慮事項（MVVM視点）

#### 7.1 既存コードの段階的MVVM化
- **段階的リファクタリング**: 既存コンポーネントを順次MVVM構造に変換
- **ラッパーコンポーネントの活用**: 既存ロジックをViewModelに移行中はラッパーで対応
- **段階的ロールアウト**: フィーチャーフラグによる新旧実装の共存

#### 7.2 ライブラリとフレームワークの選択
- **Vue.js 3のComposition API**: モジュール化と再利用性の向上
- **TypeScriptのstrictモード**: 型安全性の向上とリファクタリング支援
- **状態管理ライブラリの検討**: Pinia または Vuex 4 の導入評価
- **テストフレームワーク**: Vitest または Jest によるViewModel単体テスト

---

## 参考資料

1. **MVVMパターン**
   - Martin Fowler. "Patterns of Enterprise Application Architecture." Addison-Wesley, 2002. (MVVM章)
   - John Gossman. "Introduction to Model/View/ViewModel pattern for building WPF apps." Microsoft Blog, 2005.
   - "Vue.js Documentation: Reactivity in Depth" - https://vuejs.org/guide/extras/reactivity-in-depth.html

2. **アーキテクチャパターン**
   - Evans, Eric. "Domain-Driven Design: Tackling Complexity in the Heart of Software." Addison-Wesley, 2003.
   - Fowler, Martin. "Patterns of Enterprise Application Architecture." Addison-Wesley, 2002.
   - Hohpe, Gregor, and Bobby Woolf. "Enterprise Integration Patterns." Addison-Wesley, 2004.

3. **プロジェクト固有資料**
   - `implementation_plan.md` - フロントエンド実装方針（Vue.js Composition APIベース）
   - `backend/go.mod` - バックエンド依存関係
   - `frontend/package.json` - フロントエンド依存関係
   - `.env.example` - 環境変数テンプレート

4. **セキュリティガイドライン**
   - OWASP Application Security Verification Standard (ASVS)
   - NIST Special Publication 800-53: Security and Privacy Controls for Information Systems
   - OWASP Cheat Sheet Series: https://cheatsheetseries.owasp.org/

---

*このドキュメントはレビュー中のドラフトです。*  
*内容に関するご意見や修正指示がありましたら、お知らせください。*  
*次のステップに進む前に、このドラフトに対するフィードバックをお待ちしています。*