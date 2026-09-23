async function act(id, action) {
  const res = await fetch(`/api/containers/${id}/${action}`, { method: "POST" });
  if (!res.ok) alert(await res.text());
}

let logSource;
function showLogs(id) {
  logSource?.close();
  const pre = document.querySelector("#logs");
  pre.textContent = "";
  logSource = new EventSource(`/api/containers/${id}/logs`);
  logSource.onmessage = (e) => { pre.textContent += e.data + "\n"; };
  logSource.addEventListener("logerror", (e) => {
    pre.textContent += "[エラー] " + e.data + "\n";
    logSource.close();
  });
}

const es = new EventSource("/events");
es.onmessage = (e) => render(JSON.parse(e.data));

function render(containers) {
  document.querySelector("tbody").innerHTML = containers.map(c => `
    <tr>
      <td>${c.Name}<td>${c.Image}<td>${c.State}
      <td>${c.CPU.toFixed(1)}%<td>${(c.Memory / 1048576).toFixed(0)} MB
      <td>${(c.Ports || []).map(p => p.Public + "→" + p.Private).join(" ")}
      <td>${(c.Networks || []).join(" ")}
      <td>
        <button onclick="act('${c.ID}','start')">start</button>
        <button onclick="act('${c.ID}','stop')">stop</button>
        <button onclick="act('${c.ID}','restart')">restart</button>
        <button onclick="showLogs('${c.ID}')">logs</button>
    </tr>`).join("");
}
