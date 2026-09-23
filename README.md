# ahab

Docker の状態とコンテナ同士のつながりをブラウザで一望する常駐ダッシュボード。

    docker run -d --name ahab \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -p 18510:8080 \
      ghcr.io/tommykey-apps/ahab

http://localhost:18510

ローカル専用。認証機構は無いのでネットに公開しないこと。
Docker socket にアクセスできる = ホストの root 相当の権限を持つ。

## 画面

- コンテナ一覧 (状態 / CPU / メモリ / 公開ポート / ネットワーク)、start / stop / restart、ログ追従
- イメージ一覧と削除、ボリューム一覧と削除
- ネットワーク構成図 (コンテナ / ネットワーク / ボリューム / 公開ポート)
- 日本語 / English、明 / 暗 (OS の設定に追随、画面で切替)

## 開発

    go run .

`DOCKER_HOST` が `unix://` で始まればその socket、無ければ `/var/run/docker.sock` を使う。
