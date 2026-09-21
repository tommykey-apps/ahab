const es = new EventSource("/events");
es.onmessage = (e) => render(JSON.parse(e.data));

function render(containers) {
  document.querySelector("tbody").innerHTML = containers.map(c => `
    <tr>
      <td>${c.Name}<td>${c.Image}<td>${c.State}
      <td>${(c.Ports || []).map(p=> p.Public + "->"+p.Private).join(" ")}
      <td>${(c.Networks || []).join(" ")}
    </tr>`).join("");
}
