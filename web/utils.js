function mapRange(x, a, b, m, n) {
    return m + ((x - a) / (b - a)) * (n - m);
}

function deprecated_exportCanvasToJPEG(canvas) {
    const dataURI = canvas.toDataURL("image/jpeg");
    const base64Str = dataURI.substring(23);
    return [dataURI, atob(base64Str).length];
}

function exportCanvasToJPEG(canvas, fileName, quality) {
    if (!fileName.endsWith(".jpg") || !fileName.endsWith(".jpeg")) {
        fileName += ".jpg";
    }

    return new Promise((resolve) => {
        canvas.toBlob(
            (blob) => {
                const file = new File([blob], fileName, {
                    type: "application/octet-stream",
                });
                resolve(file);
            },
            "image/jpeg",
            quality === undefined ? 0.75 : quality
        );
    });
}

function image2Canvas(img) {
    const canvas = document.createElement("canvas");
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0);
    return canvas;
}

function getImagePixels(image) {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    canvas.width = image.width;
    canvas.height = image.height;
    ctx.drawImage(image, 0, 0);
    return ctx.getImageData(0, 0, image.width, image.height).data;
}

function getCanvasPixels(canvas) {
    const ctx = canvas.getContext("2d");
    return ctx.getImageData(0, 0, canvas.width, canvas.height).data;
}

function createCanvasFromPixels(data, width, height) {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    canvas.width = width;
    canvas.height = height;

    const imageData = ctx.createImageData(width, height);
    imageData.data.set(data);
    ctx.putImageData(imageData, 0, 0);

    return canvas;
}
