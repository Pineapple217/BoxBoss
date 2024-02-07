var term = new Terminal({
  fontSize: 14,
  // fontFamily: "'JetBrains Mono', monospace",
  // cursorBlink,
  cols: 115,
  rows: 50,
});
term.open(document.getElementById("terminal"));
const feed = new EventSource("/h/building_sse");
// term.writeln("Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ");
feed.addEventListener("close", (event) => {
  console.log("closed");
  feed.close();
});
feed.addEventListener("message_encoded", (event) => {
  const d = atob(JSON.parse(event.data));
  // let uint8Array = new Uint8Array(d.length);
  // for (var i = 0; i < d.length; i++) {
  //   uint8Array[i] = d.charCodeAt(i);
  // }
  // console.log(uint8Array);
  term.write(d);
});
feed.addEventListener("message", (event) => {
  term.writeln(event.data);
});
