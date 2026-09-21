const es = new EventSource("/events");
es.onmessage = (e) => render(JSON.parse(e.data));

function render(containers) {
  document.querySelector("thody").innerHTML = containers.map(c => `
    <tr>
      <td>${c.Name}<td>${c,Image}<td>${s.State}
      <td>${(c.Ports || []).map(p=> p.Public + "->"+p.Private).join(" ")}
      <td>${(c.Networks || []).join(" ")}
    </tr>`).join("");
}
