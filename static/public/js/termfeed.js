var term = new Terminal({
  fontSize: 14,
  // fontFamily: "'JetBrains Mono', monospace",
  // cursorBlink,
  cols: 115,
  rows: 50,
});
term.open(document.getElementById("terminal"));
const feed = new EventSource("/building_sse");
feed.addEventListener("close", (event) => {
  console.log("closed");
  feed.close();
});
feed.addEventListener("message", (event) => {
  let d;
  try {
    d = atob(JSON.parse(event.data));
  } catch {
    d = event.data.replace(/\\n/g, "\n").replace(/\\r/g, "\r");
  }
  term.write(d);
});
term.reset();
