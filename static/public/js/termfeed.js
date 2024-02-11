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
  let d;
  try {
    d = atob(JSON.parse(event.data));
  } catch {
    d = event.data.replace(/\\n/g, "\n").replace(/\\r/g, "\r");
  }
  term.write(d);
});
// feed.addEventListener("message", (event) => {
//   term.writeln(event.data);
// });
term.reset();
