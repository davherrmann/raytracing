let canvas = document.querySelector("canvas");
let ctx = canvas.getContext("2d");

function draw(x, y, r, g, b) {
  ctx.fillStyle = `rgba(${r},${g},${b},1)`;
  ctx.fillRect(x, y, 1, 1);
}

async function readStream() {
  const chunkSize = 7;
  let res = await fetch("/stream");
  let reader = res.body.getReader();
  let buffer = new Uint8Array(0);

  while (true) {
    let { done, value } = await reader.read();
    if (done) {
      return;
    }

    // append to buffer
    let newBuffer = new Uint8Array(buffer.length + value.length);
    newBuffer.set(buffer, 0);
    newBuffer.set(value, buffer.length);
    buffer = newBuffer;

    // iterate buffer
    let position = 0;
    while (position + chunkSize < buffer.length) {
      let dataView = new DataView(buffer.buffer, position);
      let littleEndian = true;
      let x = dataView.getUint16(0, littleEndian);
      let y = dataView.getUint16(2, littleEndian);
      let r = dataView.getUint8(4, littleEndian);
      let g = dataView.getUint8(5, littleEndian);
      let b = dataView.getUint8(6, littleEndian);

      if (x === 0 && y === 0) {
        console.log("AAAA");
        console.log(buffer.buffer);
      }

      draw(x, y, r, g, b);
      position += chunkSize;
    }
    buffer = buffer.slice(position);
  }
}

async function retryReadStream() {
  try {
    await readStream();
  } catch (err) {
    console.log(err);
    setTimeout(retryReadStream, 1000);
  } finally {
  }
}

async function handleChange(event) {
  let form = event.target.closest("form");
  let data = new FormData(form);

  await fetch("/change", { method: "POST", body: data });
}

async function randomizeColors(event) {
  await fetch("/randomize", { method: "POST" });
}

document.addEventListener("DOMContentLoaded", () => {
  retryReadStream();
});

window.rt = {
  handleChange,
  randomizeColors,
};
