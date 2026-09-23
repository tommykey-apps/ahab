import { t, lang, detectLang, setLang, fmtBytes, fmtDate, fmtNumber } from "./i18n.js";

const $ = (sel, root = document) => root.querySelector(sel);
const $$ = (sel, root = document) => Array.from(root.querySelectorAll(sel));
const esc = (s) => String(s ?? "").replace(/[&<>"']/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]));

const THEME_KEY = "ahab.theme";
const darkMedia = matchMedia("(prefers-color-scheme: dark)");

function currentTheme() {
  return document.documentElement.dataset.theme || (darkMedia.matches ? "dark" : "light");
}

function applyTheme(theme, persist) {
  document.documentElement.dataset.theme = theme;
  if (persist) {
    try { localStorage.setItem(THEME_KEY, theme); } catch {}
  }
  for (const b of $$("[data-theme-option]")) b.setAttribute("aria-checked", String(b.dataset.themeOption === theme));
  document.dispatchEvent(new CustomEvent("themechange", { detail: { theme } }));
}

function syncThemeSwitch() {
  const theme = currentTheme();
  for (const b of $$("[data-theme-option]")) b.setAttribute("aria-checked", String(b.dataset.themeOption === theme));
}

function initTheme() {
  syncThemeSwitch();
  $("#theme-switch").addEventListener("change", (e) => {
    const b = e.target.closest("[data-theme-option]");
    if (b && b.getAttribute("aria-checked") === "true") applyTheme(b.dataset.themeOption, true);
  });
  darkMedia.addEventListener("change", () => {
    if (document.documentElement.dataset.theme) return;
    syncThemeSwitch();
    document.dispatchEvent(new CustomEvent("themechange", { detail: { theme: currentTheme() } }));
  });
}

function initLang() {
  setLang(detectLang());
  for (const a of $$("[data-lang]")) {
    a.addEventListener("click", (e) => {
      e.preventDefault();
      setLang(a.dataset.lang);
    });
  }
  document.addEventListener("langchange", () => {
    markLang();
    $$(".dads-tab__panel").forEach((panel, i) => {
      const h = panel.querySelector(":scope > h2[tabindex='-1']");
      const tab = $$("[data-js-tab]")[i];
      if (h && tab) h.textContent = tab.textContent.trim();
    });
    renderAll();
  });
  markLang();
}

function markLang() {
  for (const a of $$("[data-lang]")) {
    if (a.dataset.lang === lang) { a.setAttribute("aria-current", "true"); a.setAttribute("data-current", ""); }
    else { a.removeAttribute("aria-current"); a.removeAttribute("data-current"); }
  }
}

const ICONS = {
  error: '<path d="M8.25 21 3 15.75v-7.5L8.25 3h7.5L21 8.25v7.5L15.75 21h-7.5Z" fill="currentcolor"/><path d="m12 13.4-2.85 2.85-1.4-1.4L10.6 12 7.75 9.15l1.4-1.4L12 10.6l2.85-2.85 1.4 1.4L13.4 12l2.85 2.85-1.4 1.4L12 13.4Z" fill="Canvas"/>',
  success: '<circle cx="12" cy="12" r="10" fill="currentcolor"/><path d="m17.6 9.6-7 7-4.3-4.3L7.7 11l2.9 2.9 5.7-5.6 1.3 1.4Z" fill="Canvas"/>',
  warning: '<path d="M1 21 12 2l11 19H1Z" fill="currentcolor"/><path d="M13 15h-2v-5h2v5Z" fill="Canvas"/><circle cx="12" cy="17" r="1" fill="Canvas"/>',
  "info-1": '<circle cx="12" cy="12" r="10" fill="currentcolor"/><circle cx="12" cy="8" r="1" fill="Canvas"/><path d="M11 11h2v6h-2z" fill="Canvas"/>',
};

function bannerHTML(type, title, body, { closable = true, actions = "" } = {}) {
  return `
    <div class="dads-notification-banner" data-style="standard" data-type="${type}" role="${type === "error" ? "alert" : "status"}">
      <h2 class="dads-notification-banner__heading">
        <svg class="dads-notification-banner__icon" width="24" height="24" viewBox="0 0 24 24" aria-hidden="true">${ICONS[type]}</svg>
        <span class="dads-notification-banner__heading-text">${esc(title)}</span>
      </h2>
      ${closable ? `<button class="dads-notification-banner__close" type="button" data-close>
        <svg class="dads-notification-banner__close-icon" width="24" height="24" viewBox="0 0 24 24" aria-hidden="true"><path d="m6.4 18.6-1-1 5.5-5.6-5.6-5.6 1.1-1 5.6 5.5 5.6-5.6 1 1.1L13 12l5.6 5.6-1 1L12 13l-5.6 5.6Z" fill="currentcolor"/></svg>
        <span class="dads-notification-banner__close-label">${esc(t("action.close"))}</span>
      </button>` : ""}
      ${body ? `<div class="dads-notification-banner__body"><p>${esc(body)}</p></div>` : ""}
      ${actions ? `<div class="dads-notification-banner__actions">${actions}</div>` : ""}
    </div>`;
}

function notify(type, title, body, { timeout = 0, id = "" } = {}) {
  const box = $("#notices");
  if (id) box.querySelector(`[data-notice-id="${id}"]`)?.remove();
  const wrap = document.createElement("div");
  if (id) wrap.dataset.noticeId = id;
  wrap.innerHTML = bannerHTML(type, title, body);
  wrap.querySelector("[data-close]")?.addEventListener("click", () => wrap.remove());
  box.prepend(wrap);
  if (timeout) setTimeout(() => wrap.remove(), timeout);
  return wrap;
}

function dismissNotice(id) {
  $("#notices").querySelector(`[data-notice-id="${id}"]`)?.remove();
}

function confirmDialog({ title, body, okLabel }) {
  const dlg = $("#confirm");
  $("#confirm-title").textContent = title;
  $("#confirm-body").textContent = body;
  $("[data-confirm-ok]", dlg).textContent = okLabel;
  return new Promise((resolve) => {
    const done = (v) => { dlg.close(); resolve(v); };
    const ok = () => done(true);
    const cancel = () => done(false);
    $("[data-confirm-ok]", dlg).addEventListener("click", ok, { once: true });
    dlg.addEventListener("close", () => { $("[data-confirm-ok]", dlg).removeEventListener("click", ok); resolve(false); }, { once: true });
    for (const b of $$("[data-confirm-cancel]", dlg)) b.addEventListener("click", cancel, { once: true });
    dlg.showModal();
  });
}

async function api(method, url) {
  const res = await fetch(url, { method });
  if (res.ok) return res;
  const message = (await res.text()).trim() || res.statusText;
  const err = new Error(message);
  err.status = res.status;
  throw err;
}

function errorBody(err) {
  if (err.status === 404 || err.status === 405) return t("error.notImplemented", { status: err.status });
  if (err.status) return t("error.http", { status: err.status, message: err.message });
  return err.message;
}

function makePanel(name, spec) {
  const root = $(`[data-resource="${name}"]`);
  const p = {
    name, spec, root,
    data: null,
    error: null,
    filter: "",
    stateFilter: "",
    sort: spec.defaultSort,
    page: 1,
    pageSize: 25,
    busy: new Set(),
  };
  $("[data-filter]", root)?.addEventListener("input", (e) => { p.filter = e.target.value.trim().toLowerCase(); p.page = 1; render(p); });
  $("[data-state-filter]", root)?.addEventListener("change", (e) => { p.stateFilter = e.target.value; p.page = 1; render(p); });
  $("[data-page-size]", root)?.addEventListener("change", (e) => { p.pageSize = Number(e.target.value); p.page = 1; render(p); });
  $("[data-pager-prev]", root)?.addEventListener("click", () => { p.page--; render(p); });
  $("[data-pager-next]", root)?.addEventListener("click", () => { p.page++; render(p); });
  $("[data-reload]", root)?.addEventListener("click", () => load(p));
  $("[data-head]", root)?.addEventListener("click", (e) => {
    const th = e.target.closest("[data-sort-key]");
    if (!th) return;
    const key = th.dataset.sortKey;
    p.sort = { key, dir: p.sort.key === key && p.sort.dir === "asc" ? "desc" : "asc" };
    render(p);
  });
  $("[data-body]", root)?.addEventListener("click", (e) => {
    const b = e.target.closest("[data-action]");
    if (b) spec.onAction?.(p, b.dataset.action, b.dataset.id, b);
  });
  return p;
}

function show(p, which) {
  for (const s of ["loading", "error", "empty", "table"]) {
    const el = p.root.querySelector(`[data-state-${s}]`);
    if (el) el.hidden = s !== which;
  }
  const pager = $("[data-pager]", p.root);
  if (pager) pager.hidden = which !== "table";
}

async function load(p) {
  if (!p.spec.url) return;
  p.error = null;
  if (p.data === null) show(p, "loading");
  try {
    const res = await api("GET", p.spec.url);
    if ((res.headers.get("content-type") || "").includes("text/html")) {
      throw Object.assign(new Error("not implemented"), { status: 404 });
    }
    p.data = p.spec.parse ? p.spec.parse(await res.text()) : await res.json();
  } catch (err) {
    p.error = err;
  }
  render(p);
}

function render(p) {
  if (p.spec.render) return p.spec.render(p);
  if (p.error) {
    $("[data-state-error]", p.root).innerHTML = bannerHTML("error", t("error.fetch.title", { what: t(`what.${p.name}`) }), errorBody(p.error), {
      closable: false,
      actions: `<button class="dads-button" type="button" data-size="md" data-type="outline" data-retry>${esc(t("action.retry"))}</button>`,
    });
    $("[data-retry]", p.root).addEventListener("click", () => { if (p.spec.url) { p.data = null; load(p); } else location.reload(); });
    $("[data-count]", p.root) && ($("[data-count]", p.root).textContent = "");
    return show(p, "error");
  }
  if (p.data === null) return show(p, "loading");

  const all = p.data;
  const rows = all.filter((r) => p.spec.match(r, p.filter, p.stateFilter));
  const cols = p.spec.columns();
  const col = cols.find((c) => c.key === p.sort.key) || cols[0];
  const dir = p.sort.dir === "asc" ? 1 : -1;
  rows.sort((a, b) => {
    const va = col.sortValue ? col.sortValue(a) : col.value(a);
    const vb = col.sortValue ? col.sortValue(b) : col.value(b);
    if (typeof va === "number" && typeof vb === "number") return (va - vb) * dir;
    return String(va).localeCompare(String(vb), lang) * dir;
  });

  const count = $("[data-count]", p.root);
  if (count) count.textContent = rows.length === all.length ? t("count", { n: fmtNumber(all.length) }) : t("count.filtered", { n: fmtNumber(rows.length), total: fmtNumber(all.length) });

  if (all.length === 0 || rows.length === 0) {
    $("[data-state-empty]", p.root).textContent = t(all.length === 0 ? `empty.${p.name}` : "empty.filtered");
    return show(p, "empty");
  }

  const pages = p.pageSize ? Math.max(1, Math.ceil(rows.length / p.pageSize)) : 1;
  p.page = Math.min(Math.max(1, p.page), pages);
  const pageRows = p.pageSize ? rows.slice((p.page - 1) * p.pageSize, p.page * p.pageSize) : rows;

  $("[data-head]", p.root).innerHTML = cols.map((c) => {
    const sorted = c.key === p.sort.key;
    const hcls = (c.cls || "").split(" ").filter((x) => x === "cell-num").join(" ");
    if (!c.sortable) return `<th class="dads-table__col-header ${hcls}" scope="col">${esc(c.label)}</th>`;
    const icon = sorted
      ? (p.sort.dir === "asc"
        ? '<path d="M17 18.12L21.27 14L22 14.7L16.5 20L11 14.7L11.73 14L16 18.12V4H17V18.12ZM14 8.92L11.73 11L9 8.52V20H6V8.52L3.27 11L1 8.93L7.5 3L14 8.93Z"/>'
        : '<path d="M8 5.88L12.27 10L13 9.3L7.5 4L2 9.3L2.73 10L7 5.88V20H8V5.88ZM10 15.08L12.27 13L15 15.48V4H18V15.48L20.73 13L23 15.07L16.5 21L10 15.07Z"/>')
      : '<path d="M17 18.11L21.27 14L22 14.7L16.5 20L11 14.7L11.73 14L16 18.12V4H17V18.12ZM8 5.88L12.27 10L13 9.3L7.5 4L2 9.3L2.73 10L7 5.88V20H8V5.88Z"/>';
    return `<th class="dads-table__sort-header ${hcls}" scope="col" data-sort-key="${c.key}" ${sorted ? `aria-sort="${p.sort.dir === "asc" ? "ascending" : "descending"}"` : ""}>
      <div class="dads-table__sort-inner"><div class="dads-table__sort-label">
        <button class="dads-table__sort-button" type="button" aria-label="${esc(t("sort.by", { col: c.label }))}">${esc(c.label)}
          <span class="dads-table__sort-icon"><svg class="dads-table__sort-svg" width="24" height="24" viewBox="0 0 24 24" fill="currentcolor" aria-hidden="true">${icon}</svg></span>
        </button>
      </div></div></th>`;
  }).join("");

  $("[data-body]", p.root).innerHTML = pageRows.map((r) => `<tr>${cols.map((c) => `<td class="${c.cls || ""}">${c.html ? c.html(r, p) : esc(c.value(r))}</td>`).join("")}</tr>`).join("");

  const pager = $("[data-pager]", p.root);
  if (pager) {
    $("[data-pager-pos]", pager).textContent = t("pager.pos", { page: p.page, pages });
    $("[data-pager-prev]", pager).disabled = p.page <= 1;
    $("[data-pager-next]", pager).disabled = p.page >= pages;
  }
  show(p, "table");
}

function button(label, action, id, { danger = false, disabled = false, size = "xs" } = {}) {
  return `<button class="dads-button" type="button" data-size="${size}" data-type="outline" ${danger ? "data-danger" : ""} data-action="${action}" data-id="${esc(id)}" ${disabled ? "disabled" : ""}>${esc(label)}</button>`;
}

function stateChip(state) {
  const color = { running: "green", paused: "yellow", restarting: "yellow", dead: "red", removing: "red" }[state] || "gray";
  const label = t(`state.${state}`) === `state.${state}` ? state : t(`state.${state}`);
  return `<span class="dads-chip-label" data-style="outlined" data-color="${color}">${esc(label)}</span>`;
}

const containers = makePanel("containers", {
  defaultSort: { key: "Name", dir: "asc" },
  match: (c, f, st) => (!st || c.State === st) && (!f || c.Name.toLowerCase().includes(f) || c.Image.toLowerCase().includes(f)),
  columns: () => [
    { key: "Name", label: t("col.name"), sortable: true, cls: "cell-name cell-nowrap", value: (c) => c.Name },
    { key: "Image", label: t("col.image"), sortable: true, cls: "cell-mono", value: (c) => c.Image },
    { key: "State", label: t("col.state"), sortable: true, value: (c) => c.State, html: (c) => stateChip(c.State) },
    { key: "CPU", label: t("col.cpu"), sortable: true, cls: "cell-num cell-nowrap", value: (c) => (c.State === "running" ? c.CPU : -1), html: (c) => (c.State === "running" ? `${fmtNumber(c.CPU, 1)} %` : '<span class="cell-muted">-</span>') },
    { key: "Memory", label: t("col.memory"), sortable: true, cls: "cell-num cell-nowrap", value: (c) => (c.State === "running" ? c.Memory : -1), html: (c) => (c.State === "running" ? esc(fmtBytes(c.Memory)) : '<span class="cell-muted">-</span>') },
    { key: "Ports", label: t("col.ports"), sortable: false, cls: "cell-mono cell-nowrap", value: (c) => (c.Ports || []).map((p) => `${p.Public}:${p.Private}`).join(" "), html: (c) => (c.Ports || []).length ? esc((c.Ports || []).map((p) => `${p.Public}:${p.Private}`).join(" ")) : '<span class="cell-muted">-</span>' },
    { key: "Networks", label: t("col.networks"), sortable: false, value: (c) => (c.Networks || []).join(" ") },
    { key: "actions", label: t("col.actions"), sortable: false, cls: "cell-actions", value: () => "", html: (c, p) => {
      const busy = p.busy.has(c.ID);
      const size = mobile.matches ? "sm" : "xs";
      const on = c.State === "running" || c.State === "paused" || c.State === "restarting";
      return `<div class="actions">${[
        on ? button(t("action.stop"), "stop", c.ID, { disabled: busy, size }) : button(t("action.start"), "start", c.ID, { disabled: busy, size }),
        button(t("action.restart"), "restart", c.ID, { disabled: busy, size }),
        button(t("action.logs"), "logs", c.ID, { size }),
      ].join("")}</div>`;
    } },
  ],
  onAction: async (p, action, id, btn) => {
    const c = p.data.find((x) => x.ID === id);
    if (!c) return;
    if (action === "logs") return openLogs(c);
    p.busy.add(id);
    render(p);
    try {
      await api("POST", `/api/containers/${encodeURIComponent(id)}/${action}`);
    } catch (err) {
      notify("error", t("error.action.title"), `${c.Name}: ${errorBody(err)}`);
    } finally {
      p.busy.delete(id);
      render(p);
    }
  },
});

function connectEvents() {
  const status = $("#conn-status");
  const setStatus = (key, color) => { status.dataset.i18n = key; status.textContent = t(key); status.dataset.color = color; };
  const es = new EventSource("/events");
  es.onopen = () => { setStatus("status.connected", "green"); dismissNotice("connection"); };
  es.onmessage = (e) => {
    containers.data = JSON.parse(e.data) || [];
    containers.error = null;
    render(containers);
    if (graph.root.offsetParent !== null) scheduleGraph();
  };
  es.onerror = () => {
    setStatus("status.reconnecting", "yellow");
    if (containers.data === null) {
      containers.error = new Error(t("error.connection.body"));
      render(containers);
    }
    if (!$("#notices [data-notice-id='connection']")) notify("error", t("error.connection.title"), t("error.connection.body"), { id: "connection" });
  };
}

let logSource = null;
let logLines = 0;
const LOG_MAX = 2000;

function openLogs(c) {
  closeLogs();
  const panel = $("#logs");
  panel.hidden = false;
  panel.dataset.name = c.Name;
  $("#logs-heading").textContent = t("logs.title", { name: c.Name });
  const out = $("#logs-out");
  out.innerHTML = `<span class="cell-muted" data-waiting>${esc(t("logs.waiting"))}</span>`;
  logLines = 0;
  logSource = new EventSource(`/api/containers/${encodeURIComponent(c.ID)}/logs`);
  logSource.onmessage = (e) => appendLog(e.lastEventId, e.data, false);
  logSource.addEventListener("logerror", (e) => { appendLog("", t("logs.error", { message: e.data }), true); closeLogs(false); });
  panel.scrollIntoView({ block: "nearest" });
}

function appendLog(ts, text, isError) {
  const out = $("#logs-out");
  out.querySelector("[data-waiting]")?.remove();
  const line = document.createElement("div");
  if (isError) line.className = "err";
  const time = ts ? fmtLogTime(ts) : "";
  if (time) {
    const span = document.createElement("span");
    span.className = "ts";
    span.textContent = time + " ";
    line.append(span);
  } else if (ts) {
    text = ts + " " + text;
  }
  line.append(document.createTextNode(text));
  out.append(line);
  if (++logLines > LOG_MAX) { out.firstElementChild?.remove(); logLines--; }
  if ($("#logs-follow").checked) out.scrollTop = out.scrollHeight;
}

function fmtLogTime(ts) {
  if (!/^\d{4}-\d{2}-\d{2}T/.test(ts)) return "";
  const d = new Date(ts);
  return Number.isNaN(d.getTime()) ? "" : d.toLocaleTimeString(lang, { hour12: false });
}

function closeLogs(hide = true) {
  logSource?.close();
  logSource = null;
  if (hide) $("#logs").hidden = true;
}

$("#logs-close").addEventListener("click", () => closeLogs());
$("#logs-clear").addEventListener("click", () => { $("#logs-out").textContent = ""; logLines = 0; });

const images = makePanel("images", {
  url: "/api/images",
  defaultSort: { key: "Tags", dir: "asc" },
  match: (i, f) => !f || (i.Tags || []).some((tg) => tg.toLowerCase().includes(f)) || i.ID.includes(f),
  columns: () => [
    { key: "Tags", label: t("col.tags"), sortable: true, cls: "cell-name", value: (i) => (i.Tags || []).join("\n") || "￿", html: (i) => (i.Tags || []).length ? `<ul class="tag-list">${i.Tags.map((x) => `<li>${esc(x)}</li>`).join("")}</ul>` : `<span class="cell-muted">${esc(t("noTag"))}</span>` },
    { key: "ID", label: t("col.id"), sortable: false, cls: "cell-mono cell-nowrap", value: (i) => i.ID.replace("sha256:", "").slice(0, 12) },
    { key: "Size", label: t("col.size"), sortable: true, cls: "cell-num cell-nowrap", value: (i) => i.Size, html: (i) => esc(fmtBytes(i.Size)) },
    { key: "Created", label: t("col.created"), sortable: true, cls: "cell-nowrap", value: (i) => i.Created, html: (i) => esc(fmtDate(i.Created)) },
    { key: "UsedBy", label: t("col.usedBy"), sortable: true, value: (i) => (i.UsedBy || []).join(" "), html: (i) => (i.UsedBy || []).length ? esc(i.UsedBy.join(", ")) : '<span class="cell-muted">-</span>' },
    { key: "actions", label: t("col.actions"), sortable: false, cls: "cell-actions", value: () => "", html: (i, p) => `<div class="actions">${button(t("action.delete"), "delete", i.ID, { danger: true, disabled: p.busy.has(i.ID), size: mobile.matches ? "sm" : "xs" })}</div>` },
  ],
  onAction: async (p, action, id) => {
    const img = p.data.find((x) => x.ID === id);
    if (!img || action !== "delete") return;
    const name = (img.Tags || [])[0] || id.replace("sha256:", "").slice(0, 12);
    if (!(await confirmDialog({ title: t("confirm.deleteImage.title"), body: t("confirm.deleteImage.body", { name }), okLabel: t("action.delete") }))) return;
    await removeResource(p, id, name, `/api/images/${encodeURIComponent(id)}`, true);
  },
});

const volumes = makePanel("volumes", {
  url: "/api/volumes",
  defaultSort: { key: "Name", dir: "asc" },
  match: (v, f) => !f || v.Name.toLowerCase().includes(f),
  columns: () => [
    { key: "Name", label: t("col.name"), sortable: true, cls: "cell-name", value: (v) => v.Name },
    { key: "Driver", label: t("col.driver"), sortable: true, cls: "cell-nowrap", value: (v) => v.Driver },
    { key: "Created", label: t("col.created"), sortable: true, cls: "cell-nowrap", value: (v) => v.Created, html: (v) => esc(fmtDate(v.Created)) },
    { key: "InUse", label: t("col.inUse"), sortable: true, value: (v) => (v.InUse ? 1 : 0), html: (v) => (v.InUse ? `<span class="dads-chip-label" data-style="outlined" data-color="green">${esc(t("yes"))}</span>` : '<span class="cell-muted">-</span>') },
    { key: "actions", label: t("col.actions"), sortable: false, cls: "cell-actions", value: () => "", html: (v, p) => `<div class="actions">${button(t("action.delete"), "delete", v.Name, { danger: true, disabled: v.InUse || p.busy.has(v.Name), size: mobile.matches ? "sm" : "xs" })}</div>` },
  ],
  onAction: async (p, action, name) => {
    if (action !== "delete") return;
    if (!(await confirmDialog({ title: t("confirm.deleteVolume.title"), body: t("confirm.deleteVolume.body", { name }), okLabel: t("action.delete") }))) return;
    await removeResource(p, name, name, `/api/volumes/${encodeURIComponent(name)}`, false);
  },
});

async function removeResource(p, id, name, url, allowForce) {
  p.busy.add(id);
  render(p);
  try {
    try {
      await api("DELETE", url);
    } catch (err) {
      if (allowForce && err.status === 409 && (await confirmDialog({ title: t("confirm.force.title"), body: t("confirm.force.body", { message: err.message }), okLabel: t("action.forceDelete") }))) {
        await api("DELETE", `${url}?force=1`);
      } else {
        throw err;
      }
    }
    notify("success", t("success.deleted", { name }), "", { timeout: 6000 });
    p.busy.delete(id);
    await load(p);
  } catch (err) {
    p.busy.delete(id);
    notify("error", t("error.action.title"), `${name}: ${errorBody(err)}`);
    render(p);
  }
}

const graph = makePanel("graph", {
  url: "/api/graph",
  defaultSort: { key: "", dir: "asc" },
  parse: (text) => text,
  render: renderGraph,
});
let mermaidLoading = null;
let graphTimer = 0;

function loadMermaid() {
  if (window.mermaid) return Promise.resolve();
  if (mermaidLoading) return mermaidLoading;
  mermaidLoading = new Promise((resolve, reject) => {
    const s = document.createElement("script");
    s.src = "/static/mermaid.min.js";
    s.onload = resolve;
    s.onerror = () => { mermaidLoading = null; reject(new Error(t("graph.unavailable"))); };
    document.head.append(s);
  });
  return mermaidLoading;
}

function cssVar(name) {
  const el = document.createElement("span");
  el.style.color = `var(${name})`;
  document.body.append(el);
  const v = getComputedStyle(el).color;
  el.remove();
  return v;
}

async function renderGraph(p) {
  if (p.error) {
    $("[data-state-error]", p.root).innerHTML = bannerHTML("error", t("error.fetch.title", { what: t("what.graph") }), errorBody(p.error), { closable: false, actions: `<button class="dads-button" type="button" data-size="md" data-type="outline" data-retry>${esc(t("action.retry"))}</button>` });
    $("[data-retry]", p.root).addEventListener("click", () => { p.data = null; load(p); });
    return show(p, "error");
  }
  if (p.data === null) return show(p, "loading");
  const text = p.data.trim();
  if (text.split("\n").length <= 1) {
    $("[data-state-empty]", p.root).textContent = t("empty.graph");
    return show(p, "empty");
  }
  try {
    await loadMermaid();
    await document.fonts.load("16px 'Noto Sans JP'", text);
    await document.fonts.ready;
    window.mermaid.initialize({
      startOnLoad: false,
      securityLevel: "strict",
      theme: "base",
      fontFamily: getComputedStyle(document.documentElement).getPropertyValue("--font-family-sans").trim(),
      themeVariables: {
        primaryColor: cssVar("--color-key-50"),
        primaryTextColor: cssVar("--color-neutral-solid-gray-900"),
        primaryBorderColor: cssVar("--color-key-900"),
        secondaryColor: cssVar("--color-neutral-solid-gray-50"),
        tertiaryColor: cssVar("--color-neutral-white"),
        lineColor: cssVar("--color-neutral-solid-gray-600"),
        textColor: cssVar("--color-neutral-solid-gray-900"),
        background: cssVar("--color-neutral-white"),
        edgeLabelBackground: cssVar("--color-neutral-white"),
      },
    });
    const { svg } = await window.mermaid.render(`graph-${Date.now()}`, text);
    $("#graph").innerHTML = svg;
    show(p, "table");
  } catch (err) {
    p.error = err;
    renderGraph(p);
  }
}

function scheduleGraph() {
  clearTimeout(graphTimer);
  graphTimer = setTimeout(() => load(graph), 500);
}

const mobile = matchMedia("(max-width: 47.99rem)");
const panels = { containers, images, volumes, graph };

function renderAll() {
  for (const p of Object.values(panels)) if (p.data !== null || p.error) render(p);
  const s = $("#conn-status");
  if (s.dataset.i18n) s.textContent = t(s.dataset.i18n);
  const logs = $("#logs");
  if (!logs.hidden) $("#logs-heading").textContent = t("logs.title", { name: logs.dataset.name });
}

function applyDensity() {
  for (const tbl of $$(".dads-table")) {
    if (mobile.matches) delete tbl.dataset.size; else tbl.dataset.size = "dense";
  }
}

$("#tabs").addEventListener("tab-change", (e) => {
  const name = e.detail.selectedPanel.dataset.resource;
  const p = panels[name];
  if (p && p.spec.url && p.data === null) load(p);
});

document.addEventListener("themechange", () => { if (graph.data !== null) renderGraph(graph); });
mobile.addEventListener("change", () => { applyDensity(); renderAll(); });

initTheme();
initLang();
applyDensity();
connectEvents();
