const es = new EventSource("/events");
es.onmessage = (e) => render(JSON.parse(e.data));

function render(containers) {
  document.querySelector("tbody").innerHTML = containers.map(c => `
    <tr>
      <td>${c.Name}<td>${c.Image}<td>${c.State}
      <td>${c.CPU.toFixed(1)}%<td>${(c.Memory / 1048576).toFixed(0)} MB
      <td>${(c.Ports || []).map(p=> p.Public + "->"+p.Private).join(" ")}
      <td>${(c.Networks || []).join(" ")}
      <td>
        <button onclick="act('${c.ID}','start')">start</button>
        <button onclick="act('${c.ID}','stop')">stop</button>
        <button onclick="act('${c.ID}','restart')">restsrt</button>
        <button onclick="showLogs('${c.ID}')">logs</button>
    </tr>`).join("");
}
