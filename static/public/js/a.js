var term = new Terminal({
  fontSize: 14,
  cols: 115,
  rows: 50,
});
term.open(document.getElementById("terminal"));
// const d = atob(JSON.parse(event.data));
// term.write(d);

let jsonData;
fetch("./static/js/output.json")
  .then((response) => {
    // Check if the response is successful
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    // Parse the JSON response
    return response.json();
  })
  .then((jsonData) => {
    // Do something with the JSON data
    console.log(jsonData);
    loopThroughArray(jsonData);
  })
  .catch((error) => {
    // Handle any errors
    console.error("There was a problem fetching the JSON file:", error);
  });

async function waitForJKeyPress() {
  return new Promise((resolve) => {
    document.addEventListener(
      "keydown",
      function (event) {
        if (event.key === "j") {
          resolve();
        }
      },
      { once: true }
    );
  });
}

async function loopThroughArray(inputArray) {
  let currentIndex = 0;

  for (let i = currentIndex; i < inputArray.length; i++) {
    console.log("Press 'j' key to continue...");
    await waitForJKeyPress();
    console.log(inputArray[i]);
    try {
      const d = atob(JSON.parse(inputArray[i]));
      term.write(d);
    } catch {
      console.log("NOT JSON" + inputArray[i]);
    }
    currentIndex++;
  }
  console.log("End of array reached.");
}
