1. DB（SQLite）の仕様：保存場所・スキーマ・マイグレーション
■ 保存場所

USB内の backend/db/airgate.db に固定します。

ONE-AIR/
  backend/
    db/
      airgate.db   ← ここ。sqlite単独ファイル。


USB直下は混乱する

backend/db が最も明確でバックアップもしやすい

ロック問題が起きてもファイル単体で扱える

■ SQLite スキーマ（確定版）

ONE-AIR の性質上、必要最小に絞ります。

▼ officers（審査官アカウント）
CREATE TABLE officers (
    phone TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    role TEXT NOT NULL,   -- IMMIGRATION / SECONDARY
    pin_hash TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

▼ entries（一時審査ログ）
CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TEXT NOT NULL,
    officer_phone TEXT NOT NULL,
    mrz TEXT NOT NULL,
    risk_level TEXT NOT NULL,     -- LOW/MED/HIGH
    action TEXT NOT NULL,         -- NORMAL / SECONDARY / HOLD
    note TEXT,
    device_id TEXT NOT NULL,
    ip TEXT NOT NULL,
    user_agent TEXT,
    FOREIGN KEY (officer_phone) REFERENCES officers(phone)
);

▼ secondary（二次審査詳細）
CREATE TABLE secondary (
    entry_id INTEGER PRIMARY KEY,
    detail TEXT,
    photo BLOB,
    FOREIGN KEY (entry_id) REFERENCES entries(id)
);

■ マイグレーション方式

「複雑なマイグレーション禁止」なので シンプルなバージョン管理のみにします。

backend/db/schema_version.txt          ← 例： "1"
backend/db/migrate_1_to_2.sql          ← 次期拡張が必要な場合のみ


※短期運用なので、実質 version=1 のまま固定運用で問題なし。

✅ 2. ログ仕様：フォーマット・ローテーション・署名方式
■ 保存場所

USB内：

ONE-AIR/logs/
   airgate_2025-11-14.jsonl

■ フォーマット

JSON Lines（1行＝1イベント）
→ 監査性・grep性・後処理の容易さすべて兼ねる

▼ 行フォーマット（確定）
{
  "ts": "2025-11-14T13:22:55+09:00",
  "event": "ENTRY_RECORDED",
  "officer": "09012345678",
  "entry_id": 124,
  "mrz": "P<JPNYAMADA<<TARO<<<<<<<<<<<<<<<<<<<<<<<<<",
  "risk": "MED",
  "action": "SECONDARY",
  "device_id": "WIN-ABC123",
  "ip": "127.0.0.1",
  "ua": "Mozilla/5.0"
}

■ ローテーション

1日1ファイル。サイズ上限は 50MB で強制ローテーション。

時間基準：日次（UTCではなく JST）

サイズ基準：50MB を超えたら自動的に …_01.jsonl を作成

■ 署名方式

「短期用途」「USB単体」なので HMAC-SHA256 だけで十分。

ONE-AIR/logs/hmac_key.txt


ビルド時にランダム生成

USB内に平文配置（機密扱いだが短期用途のため許容）

各行に "sig": "..." を追加

ISE不要（明白に要求されていない）

✅ 3. 認証：CSV→DB同期・PIN方式・暗号化
■ CSV カラム構成（確定）
phone,name,role,pin
09012345678,山田太郎,IMMIGRATION,4829
08099998888,佐藤花子,SECONDARY,9351

■ CSVの扱い

USB直下の /officers.csv

起動時に読み込んで SQLite に上書き反映

officers.csv が無ければ DB の既存レコードを使う

運用者は Excel で編集 → CSV 書き出し → USBにコピー

■ PIN の暗号化方式

SHA-256（ソルト無し） に固定

PBKDF2/Bcrypt は使用しない（WindowsのUSB環境でトラブル要因）

DBには pin_hash に SHA256 を保存。

✅ 4. MRZ の取得方式（優先順位とGoとの整合性）
■ 結論：ONE-AIR は MRZ を“文字列として受け取るだけ”

MRZ取得処理そのものは 外部依存で良い。

■ 優先順位
1位：パスポートリーダ（USB HID・キーボードエミュ）  
2位：手入力（非常用）  
3位：スマホOCR（旅客対応でも使える）  
4位：NFC（今回は優先しない。実装コストが高い）


GoでNFC/OCRを内製するのは 短期では非現実。
ONE-AIR側は以下のように対応：

MRZ入力欄に貼り付ければ処理可能

UI側で「OCR起動ボタン」を出す → 外部exeを叩く（任意）

NFCは今回は 非対応でOK（利用環境のばらつきが大きい）

✅ 5. フロントエンド：全部Goにするのか？静的HTMLは許容か？

答え：静的HTMLは全面的に許容。
Goテンプレート不要。WebView不要。

理由：

Goの単体バイナリで HTTPサーバを立てる

/frontend をそのまま静的配信

Windows/Kiosk/スマホからブラウザで開ける

学習不要・依存最小

ONE-AIR の目的は「UIの美しさ」ではなく 10秒で理解できる操作性。

よって：

/frontend/index.html

/frontend/triage.html

/frontend/secondary.html

こういう HTML直書きで十分。
JSも必要最小のバニラJSのみ。
