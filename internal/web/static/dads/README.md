# デジタル庁デザインシステム 部品の複製

出典: https://github.com/digital-go-jp/design-system-example-components-html (MIT)
複製元コミット: af8b6656c8d864a22ef444d088e5568f3416f6aa (2026-09-09)
トークン: https://github.com/digital-go-jp/design-tokens v2.0.1 `dist/tokens.css` (MIT)

| ファイル | 元ファイル |
|---|---|
| tokens.css | design-tokens `dist/tokens.css` |
| global.css | `src/global.css` の `html {` 以降 (トークン定義を除く) |
| button.css | `src/components/button/button.css` |
| table.css, scroll-shadow.js | `src/components/table/` |
| chip-label.css | `src/components/chip-label/chip-label.css` |
| notification-banner.css | `src/components/notification-banner/notification-banner.css` |
| tab.css, tab.js | `src/components/tab/` |
| switch-mode.css, switch-mode.js | `src/components/switch/` |
| language-selector.css, language-selector.js | `src/components/language-selector/` |
| menu-list-box.css | `src/components/menu-list-box/menu-list-box.css` |
| menu-list.css | `src/components/menu-list/menu-list.css` |
| modal-dialog.css | `src/components/modal-dialog/modal-dialog.css` |
| select.css | `src/components/select/select.css` |
| input-text.css | `src/components/input-text/input-text.css` |
| progress-indicator.css, progress-indicator.js | `src/components/progress-indicator/` |
| heading.css | `src/components/heading/heading.css` |
| link.css | `src/components/link/link.css` |
| description-list.css | `src/components/description-list/description-list.css` |

複製した部品は案件内で固定し、上流には追随しない。変更するときはこの表の出典を見て差分を確認する。
このディレクトリのファイルは編集しない。暗色表示は `../theme.css` でトークンの値を付け替えて実現する。

書体 (`../fonts/`) は fontsource の `@fontsource/noto-sans-jp` と `@fontsource/noto-sans-mono` から複製 (SIL Open Font License 1.1、同梱の LICENSE を参照)。
