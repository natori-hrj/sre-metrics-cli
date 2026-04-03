# sre-metrics-cli

SLI/SLO の計算・可視化を行う CLI ツール。サービスごとの SLI 記録、SLO 達成率の判定、エラーバジェットの残量管理ができます。

Vercel Deployments API と UptimeRobot からのデータ自動取得にも対応しています（いずれも無料枠で利用可能）。

## インストール

```bash
go install github.com/natori-hrj/sre-metrics-cli@latest
```

または手動ビルド:

```bash
git clone https://github.com/natori-hrj/sre-metrics-cli.git
cd sre-metrics-cli
go build -o slo .
```

## 使い方

### SLI を手動で記録する

```bash
slo record --service natorium --total 10000 --success 9992
```

### SLO 達成状況を表示する

```bash
slo status --service natorium --target 99.9
```

```
╭───────────────────────────────────╮
│              SLO Status           │
│                                   │
│  Service:        natorium         │
│  SLI:            99.92%           │
│  SLO Target:     99.90%           │
│  Status:         MEETING SLO      │
│  Error Budget:   16.7% remaining  │
│                                   │
│  Requests:       29975 / 30000    │
╰───────────────────────────────────╯
```

### エラーバジェットの残量を確認する

```bash
slo budget --service natorium --target 99.9 --window 30d
```

```
╭────────────────────────────────────╮
│               Error Budget         │
│                                    │
│  Service:         natorium         │
│  Window:          30d              │
│  SLO Target:      99.90%           │
│  Current SLI:     99.92%           │
│  Allowed Errors:  29               │
│  Actual Errors:   25               │
│  Budget Left:     4 requests       │
│  Budget Percent:  16.7% remaining  │
╰────────────────────────────────────╯
```

### 外部サービスからデータを自動取得する

#### 1. サービスを初期化する

```bash
slo init --service natorium
```

対話形式で以下を設定します:

| 項目 | 説明 |
|------|------|
| Display name | サービスの表示名（例: `natorium.dev`） |
| URL | サイトの URL（例: `https://natorium.dev`） |
| Vercel API token | [Vercel Account Tokens](https://vercel.com/account/tokens) で作成 |
| Vercel Project ID | Vercel ダッシュボードの Settings > General で確認 |
| Vercel Team ID | チームの場合のみ。Hobby プランは空欄 |
| UptimeRobot API key | [UptimeRobot](https://uptimerobot.com/) の Settings > API Settings で取得 |
| UptimeRobot Monitor ID | モニター作成後、Dashboard で確認 |

#### 2. データを取得する

```bash
slo fetch --service natorium --window 30d
```

#### 3. 結果を確認する

```bash
# デプロイ成功率
slo status --service natorium-deploy --target 99.9

# 応答成功率
slo status --service natorium-uptime --target 99.9
```

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `slo init` | サービスの初期設定（API キー等） |
| `slo record` | SLI データを手動記録 |
| `slo fetch` | Vercel / UptimeRobot からデータ取得 |
| `slo status` | SLO 達成状況を表示 |
| `slo budget` | エラーバジェットの残量を表示 |

## 共通フラグ

| フラグ | デフォルト | 説明 |
|--------|-----------|------|
| `--service` | (必須) | サービス識別子 |
| `--target` | `99.9` | SLO ターゲット (%) |
| `--window` | `30d` | 集計期間（`7d` / `30d` / `90d`） |

## データ保存先

すべてのデータは `~/.slo/` に保存されます。

```
~/.slo/
├── natorium.json           # サービス設定（API キー等）
├── natorium.csv            # 手動記録の SLI データ
├── natorium-deploy.csv     # Vercel デプロイ記録
└── natorium-uptime.csv     # UptimeRobot 応答記録
```

## 無料枠での利用

| サービス | 無料枠 |
|---------|--------|
| Vercel Hobby | API アクセス無料 |
| UptimeRobot Free | 50 モニター、5 分間隔監視、API 無制限 |

## 技術スタック

- Go
- [cobra](https://github.com/spf13/cobra) — CLI フレームワーク
- [lipgloss](https://github.com/charmbracelet/lipgloss) — ターミナル UI スタイリング

## ライセンス

MIT
