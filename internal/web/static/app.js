const es = new EventSource("/events");
es.onmessage = (e) => render(JSON.parse(e.data));

function render(containers) {
  document.querySelector("tbody").innerHTML = containers.map(c => `
    <tr>
      <td>${c.Name}<td>${c.Image}<td>${c.State}
      <td>${(c.Ports || []).map(p=> p.Public + "->"+p.Private).join(" ")}
      <td>${(c.Networks || []).join(" ")}
      <td>${c.CPU.toFixed(1)}%<td>${(c.Memory / 1048576).toFixed(0)} MB
    </tr>`).join("");
}
