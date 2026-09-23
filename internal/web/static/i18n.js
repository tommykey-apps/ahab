const dict = {
  ja: {
    "app.title": "Docker の状態",
    "theme.label": "表示",
    "theme.light": "明",
    "theme.dark": "暗",
    "tab.containers": "コンテナ",
    "tab.images": "イメージ",
    "tab.volumes": "ボリューム",
    "tab.graph": "構成図",
    "filter.label": "絞り込み",
    "filter.containers": "名前・イメージで絞り込み",
    "filter.images": "タグ・ID で絞り込み",
    "filter.volumes": "名前で絞り込み",
    "state.label": "状態",
    "state.all": "すべての状態",
    "state.running": "実行中",
    "state.exited": "停止",
    "state.paused": "一時停止",
    "state.created": "作成済み",
    "state.restarting": "再起動中",
    "state.dead": "異常終了",
    "state.removing": "削除中",
    "count": "{n} 件",
    "count.filtered": "{n} 件 (全 {total} 件)",
    "pageSize.label": "表示件数",
    "pageSize.all": "すべて",
    "pager.prev": "前へ",
    "pager.next": "次へ",
    "pager.pos": "{page} / {pages} ページ",
    "col.name": "名前",
    "col.image": "イメージ",
    "col.state": "状態",
    "col.cpu": "CPU",
    "col.memory": "メモリ",
    "col.ports": "公開ポート",
    "col.networks": "ネットワーク",
    "col.actions": "操作",
    "col.tags": "タグ",
    "col.id": "ID",
    "col.size": "サイズ",
    "col.created": "作成日時",
    "col.usedBy": "使用中のコンテナ",
    "col.driver": "ドライバ",
    "col.inUse": "使用中",
    "sort.by": "{col}で並び替え",
    "action.start": "起動",
    "action.stop": "停止",
    "action.restart": "再起動",
    "action.logs": "ログ",
    "action.delete": "削除",
    "action.forceDelete": "強制削除",
    "action.cancel": "キャンセル",
    "action.reload": "再読み込み",
    "action.close": "閉じる",
    "action.clear": "消去",
    "action.retry": "再試行",
    "action.busy": "処理中",
    "yes": "はい",
    "none": "なし",
    "loading": "読み込み中",
    "empty.containers": "コンテナがありません。docker run で作成すると、ここに表示されます。",
    "empty.images": "イメージがありません。docker pull で取得すると、ここに表示されます。",
    "empty.volumes": "ボリュームがありません。docker volume create で作成すると、ここに表示されます。",
    "empty.graph": "描画する対象がありません。コンテナを起動すると表示されます。",
    "empty.filtered": "条件に一致するものがありません。絞り込みを解除してください。",
    "error.connection.title": "Docker との接続が切れました",
    "error.connection.body": "自動で再接続します。ahab が起動しているか、Docker socket にアクセスできるか確認してください。",
    "error.fetch.title": "{what}を取得できませんでした",
    "error.notImplemented": "この機能はまだサーバー側に実装されていません (HTTP {status})。",
    "error.http": "HTTP {status}: {message}",
    "error.action.title": "操作を完了できませんでした",
    "success.deleted": "{name} を削除しました",
    "confirm.deleteImage.title": "イメージを削除",
    "confirm.deleteImage.body": "{name} を削除します。コンテナが使用している場合は失敗します。",
    "confirm.deleteVolume.title": "ボリュームを削除",
    "confirm.deleteVolume.body": "{name} を削除します。中のデータは戻せません。",
    "confirm.force.title": "強制削除",
    "confirm.force.body": "Docker の応答: {message}",
    "logs.title": "ログ: {name}",
    "logs.follow": "自動スクロール",
    "logs.waiting": "ログを待っています",
    "logs.error": "ログの取得に失敗しました: {message}",
    "graph.unavailable": "構成図の描画ライブラリ (/static/mermaid.min.js) を読み込めませんでした。issue #11 を参照。",
    "noTag": "タグなし",
    "what.containers": "コンテナ一覧",
    "what.images": "イメージ一覧",
    "what.volumes": "ボリューム一覧",
    "what.graph": "構成図",
    "legend.title": "凡例",
    "legend.container": "コンテナ",
    "legend.network": "ネットワーク",
    "legend.volume": "ボリューム",
    "legend.host": "ホスト",
    "legend.connect": "ネットワーク接続",
    "legend.mount": "ボリュームのマウント",
    "legend.hostNetwork": "host ネットワーク",
  },
  en: {
    "app.title": "Docker status",
    "theme.label": "Theme",
    "theme.light": "Light",
    "theme.dark": "Dark",
    "tab.containers": "Containers",
    "tab.images": "Images",
    "tab.volumes": "Volumes",
    "tab.graph": "Network map",
    "filter.label": "Filter",
    "filter.containers": "Filter by name or image",
    "filter.images": "Filter by tag or ID",
    "filter.volumes": "Filter by name",
    "state.label": "State",
    "state.all": "All states",
    "state.running": "Running",
    "state.exited": "Exited",
    "state.paused": "Paused",
    "state.created": "Created",
    "state.restarting": "Restarting",
    "state.dead": "Dead",
    "state.removing": "Removing",
    "count": "{n} items",
    "count.filtered": "{n} of {total} items",
    "pageSize.label": "Rows per page",
    "pageSize.all": "All",
    "pager.prev": "Previous",
    "pager.next": "Next",
    "pager.pos": "Page {page} of {pages}",
    "col.name": "Name",
    "col.image": "Image",
    "col.state": "State",
    "col.cpu": "CPU",
    "col.memory": "Memory",
    "col.ports": "Ports",
    "col.networks": "Networks",
    "col.actions": "Actions",
    "col.tags": "Tags",
    "col.id": "ID",
    "col.size": "Size",
    "col.created": "Created",
    "col.usedBy": "Used by",
    "col.driver": "Driver",
    "col.inUse": "In use",
    "sort.by": "Sort by {col}",
    "action.start": "Start",
    "action.stop": "Stop",
    "action.restart": "Restart",
    "action.logs": "Logs",
    "action.delete": "Delete",
    "action.forceDelete": "Force delete",
    "action.cancel": "Cancel",
    "action.reload": "Reload",
    "action.close": "Close",
    "action.clear": "Clear",
    "action.retry": "Retry",
    "action.busy": "Working",
    "yes": "Yes",
    "none": "None",
    "loading": "Loading",
    "empty.containers": "No containers. Create one with docker run and it will appear here.",
    "empty.images": "No images. Pull one with docker pull and it will appear here.",
    "empty.volumes": "No volumes. Create one with docker volume create and it will appear here.",
    "empty.graph": "Nothing to draw. Start a container and it will appear here.",
    "empty.filtered": "Nothing matches the filter. Clear the filter.",
    "error.connection.title": "Lost connection to Docker",
    "error.connection.body": "Reconnecting automatically. Check that ahab is running and can reach the Docker socket.",
    "error.fetch.title": "Could not load {what}",
    "error.notImplemented": "The server does not provide this yet (HTTP {status}).",
    "error.http": "HTTP {status}: {message}",
    "error.action.title": "The action failed",
    "success.deleted": "Deleted {name}",
    "confirm.deleteImage.title": "Delete image",
    "confirm.deleteImage.body": "Delete {name}. This fails if a container is using it.",
    "confirm.deleteVolume.title": "Delete volume",
    "confirm.deleteVolume.body": "Delete {name}. The data inside cannot be recovered.",
    "confirm.force.title": "Force delete",
    "confirm.force.body": "Docker said: {message}",
    "logs.title": "Logs: {name}",
    "logs.follow": "Follow",
    "logs.waiting": "Waiting for log output",
    "logs.error": "Could not read logs: {message}",
    "graph.unavailable": "Could not load the diagram library (/static/mermaid.min.js). See issue #11.",
    "noTag": "untagged",
    "what.containers": "the container list",
    "what.images": "the image list",
    "what.volumes": "the volume list",
    "what.graph": "the network map",
    "legend.title": "Legend",
    "legend.container": "Container",
    "legend.network": "Network",
    "legend.volume": "Volume",
    "legend.host": "Host",
    "legend.connect": "Network link",
    "legend.mount": "Volume mount",
    "legend.hostNetwork": "Host network",
  },
};

const KEY = "ahab.lang";
export let lang = "ja";

export function t(key, vars = {}) {
  const s = dict[lang][key] ?? dict.ja[key] ?? key;
  return s.replace(/\{(\w+)\}/g, (m, k) => (k in vars ? String(vars[k]) : m));
}

export function detectLang() {
  const q = new URLSearchParams(location.search).get("lang");
  if (q === "ja" || q === "en") return q;
  try {
    const saved = localStorage.getItem(KEY);
    if (saved === "ja" || saved === "en") return saved;
  } catch {}
  return (navigator.language || "ja").startsWith("en") ? "en" : "ja";
}

export function setLang(l) {
  lang = l;
  document.documentElement.lang = l;
  try { localStorage.setItem(KEY, l); } catch {}
  applyStatic();
  document.dispatchEvent(new CustomEvent("langchange", { detail: { lang: l } }));
}

export function applyStatic(root = document) {
  for (const el of root.querySelectorAll("[data-i18n]")) el.textContent = t(el.dataset.i18n);
  for (const el of root.querySelectorAll("[data-i18n-attr]")) {
    for (const pair of el.dataset.i18nAttr.split(/\s+/)) {
      const [attr, key] = pair.split(":");
      if (attr && key) el.setAttribute(attr, t(key));
    }
  }
}

export function fmtBytes(n) {
  if (n == null) return "-";
  const units = ["B", "KB", "MB", "GB", "TB"];
  let i = 0;
  let v = n;
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++; }
  return `${new Intl.NumberFormat(lang, { maximumFractionDigits: i >= 2 ? 1 : 0 }).format(v)} ${units[i]}`;
}

export function fmtDate(iso) {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime()) || d.getTime() === 0) return "-";
  return new Intl.DateTimeFormat(lang, { dateStyle: "medium", timeStyle: "short" }).format(d);
}

export function fmtNumber(n, digits = 0) {
  return new Intl.NumberFormat(lang, { minimumFractionDigits: digits, maximumFractionDigits: digits }).format(n);
}
