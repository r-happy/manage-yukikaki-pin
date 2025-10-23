# manage-yukikaki-pin プロジェクト概要

## プロジェクト概要

`manage-yukikaki-pin` は、グループベースの地図ピン管理システムです。ユーザーがグループを作成し、メンバーと共に地図上にピンを配置・管理できるWebアプリケーションです。

## システムアーキテクチャ

本プロジェクトは、フロントエンドとバックエンドを分離したモダンなWebアプリケーション構成を採用しています。

```
manage-yukikaki-pin/
├── backend/          # Go (Echo Framework) バックエンドAPI
├── frontend/         # Angular フロントエンド
├── docs/            # ドキュメント
├── compose.yaml     # Docker Compose設定
└── Jenkinsfile      # CI/CD設定
```

### バックエンド (Go + Echo Framework)

- **フレームワーク**: Echo v4
- **データベース**: SQLite (GORM ORM)
- **認証**: JWT (JSON Web Token)
- **ポート**: 1323

#### 主要な技術スタック
- `github.com/labstack/echo/v4` - Webフレームワーク
- `gorm.io/gorm` - ORM
- `github.com/golang-jwt/jwt/v5` - JWT認証
- `github.com/google/uuid` - UUID生成

#### データモデル

1. **User** - ユーザー情報
   - UserID (UUID)
   - UserName, UserEmail, UserPassword
   - メール形式のバリデーション

2. **Group** - グループ情報
   - GroupID (UUID)
   - GroupName, GroupDescription
   - GroupCreatedBy (作成者)

3. **GroupMember** - グループメンバー管理
   - グループとユーザーの関連付け

4. **PinType** - ピンの種類
   - PinTypeID (UUID)
   - カスタマイズ可能なピンタイプ

5. **Pin** - 地図上のピン
   - PinID (UUID)
   - Latitude, Longitude (緯度・経度)
   - PinName, PinTypeDescription
   - NotAllowed フラグ

### フロントエンド (Angular)

- **フレームワーク**: Angular v20
- **UIライブラリ**: Angular Material
- **地図ライブラリ**: Leaflet
- **ポート**: 8080 (コンテナ内は4200)

#### 主要な機能ページ

- `/home` - ホームページ
- `/signin` - サインイン
- `/signup` - サインアップ
- `/dashboard` - ダッシュボード
- `/map` - 地図表示・ピン管理
- `/pin` - ピン管理
- `/profile` - プロフィール

#### 技術スタック
- Angular Material - UIコンポーネント
- Leaflet - インタラクティブ地図
- RxJS - リアクティブプログラミング
- Angular SSR - サーバーサイドレンダリング

### Webサーバー

- **本番環境**: Nginx
- 静的ファイル配信とリバースプロキシ

## API エンドポイント

### 認証不要エンドポイント

- `POST /signup` - ユーザー登録
- `POST /signin` - サインイン (JWT発行)

### 認証必須エンドポイント (JWT)

#### ユーザー
- `GET /api/me` - ログインユーザー情報取得
- `GET /api/me/groups` - ユーザーが所属するグループ一覧

#### グループ
- `POST /api/groups` - グループ作成
- `GET /api/groups/:groupID` - グループ詳細取得

#### グループメンバー
- `POST /api/groups/:groupID/add-members` - メンバー追加（管理者）
- `POST /api/groups/:groupID/join-requests` - 参加リクエスト

#### ピンタイプ
- `POST /api/groups/:groupID/pin-types` - ピンタイプ作成
- `GET /api/groups/:groupID/pin-types` - ピンタイプ一覧

#### ピン
- `POST /api/groups/:groupID/pins` - ピン作成
- `GET /api/groups/:groupID/pins` - グループのピン一覧

## セキュリティ

- **認証方式**: JWT (JSON Web Token)
- **パスワード管理**: ハッシュ化（詳細は実装を確認）
- **CORS**: 有効化
- **環境変数**: SECRET_KEY による秘密鍵管理

## デプロイメント

### Docker Compose

本プロジェクトはDocker Composeを使用した簡単なデプロイが可能です。

```yaml
services:
  - backend: ポート1323で起動
  - frontend: ポート8080で起動
```

### 環境変数

- `SECRET_KEY` - JWT署名用シークレットキー（必須）
- `FRONTEND_BACKEND_URL` - バックエンドAPIのURL（デフォルト: http://localhost:1323）

### CI/CD

- Jenkins を使用した継続的インテグレーション/デプロイメント
- `Jenkinsfile` で設定

## データ永続化

- SQLiteデータベースファイルは `backend/src/db/` にボリュームマウント
- コンテナ再起動後もデータが保持される

## 開発の始め方

1. **環境変数の設定**
   ```bash
   export SECRET_KEY="your-secret-key"
   ```

2. **Docker Composeでの起動**
   ```bash
   docker-compose up -d
   ```

3. **アクセス**
   - フロントエンド: http://localhost:8080
   - バックエンドAPI: http://localhost:1323

## ネットワーク

- Docker外部ネットワーク `manage-yukikaki-net` を使用
- バックエンドとフロントエンドはこのネットワーク上で通信

## 主要機能

1. **ユーザー管理**
   - ユーザー登録・ログイン
   - プロフィール管理

2. **グループ管理**
   - グループ作成
   - メンバー招待・管理
   - 参加リクエスト

3. **ピン管理**
   - 地図上にピンを配置
   - カスタムピンタイプ定義
   - ピン情報の表示・編集
   - グループメンバー間でのピン共有

4. **地図機能**
   - Leafletによるインタラクティブ地図
   - 緯度・経度ベースの位置情報管理

## 制約事項

- 緯度: -90〜90度の範囲
- 経度: -180〜180度の範囲
- メールアドレスの重複不可
- グループへのアクセスはメンバーのみ