# ahab

Docker の状態と**コンテナ同士のつながり**をブラウザで一望する常駐ダッシュボード。

一覧とログだけでなく、どのコンテナがどのネットワーク・ボリュームに繋がっているかを図で示すことを主眼に置く。

## 前提

ローカル専用。認証なし。ネットに公開しない。単一ホストの Docker Engine のみ。

## やらないこと

認証、リモートホスト接続、prune、イメージの pull / build、compose の編集やデプロイ、履歴の永続化。
コンテナ操作は start / stop / restart のみ。

## 起動

```
docker run -d --name ahab \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -p 18510:8080 \
  ghcr.io/tommykey-apps/ahab
```

開発時は `go run .`。socket は `DOCKER_HOST` があればそれ、無ければ `/var/run/docker.sock`。

## 機能

- **ネットワーク構成図**（Mermaid）— コンテナ / ネットワーク / ボリューム / ホスト公開ポートの関係
- コンテナ一覧（名前 / イメージ / 状態 / 公開ポート / ネットワーク / CPU / メモリ）
- start / stop / restart
- ログ追従
- イメージ一覧（タグ / ID / サイズ / 作成日時 / 使用中のコンテナ）と削除
- ボリューム一覧（名前 / ドライバ / 作成日時 / 使用中）と削除
- 日本語 / English の切替、明 / 暗の切替（OS 設定に追随）

Docker の `/events` と `/stats` をストリームで読み、Server-Sent Events でブラウザに push する。ポーリングしない。

## 図のルール

- コンテナ同士は直接結ばない。ネットワークノードを経由する（直接結線は N²/2 本になるため）
- bind mount は描かない。名前付きボリュームのみ（図が破綻するため）
- 名前でソートしてから描く。出力を決定的にするため

## 技術方針

- 標準ライブラリのみ。外部モジュールを追加しない
- フロントはビルド不要。`html/template` + SSE + 素の JavaScript。部品はデジタル庁デザインシステムの HTML 版を `internal/web/static/dads/` に複製、書体と mermaid も同梱
- 状態を持つのは単一の goroutine。他は channel 経由
- `go:embed` で単一バイナリ。ベースイメージは `scratch`

## 配布

`ghcr.io/tommykey-apps/ahab`（amd64 / arm64）。副で `go install`。
