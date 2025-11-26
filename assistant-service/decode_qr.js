const { createCanvas, loadImage } = require("canvas");
const jsQR = require("jsqr");

const filePath = process.argv[2]; // đường dẫn file ảnh QR

loadImage(filePath).then((img) => {
  const canvas = createCanvas(img.width, img.height);
  const ctx = canvas.getContext("2d");
  ctx.drawImage(img, 0, 0);
  const imageData = ctx.getImageData(0, 0, img.width, img.height);

  const code = jsQR(imageData.data, imageData.width, imageData.height);
  if (code) {
    console.log(code.data); // in ra stdout
  } else {
    console.log(""); // không đọc được
  }
});
