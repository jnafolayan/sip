if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const codec = new Worker("codec.worker.js");
let taskID = 0;

codec.onmessage = (e) => {
    if (e.data.taskID != taskID) return;
    // If this happens for whatever reason
    if (!appState.compressing) return;

    if (!e.data.error) {
        const { compressed, result, width, height } = e.data;
        psnrElement.innerText = result.PSNR.toFixed(2);

        const image = createImageFromPixels(compressed, width, height);
        const url = image.toDataURL("image/jpeg");
        var base64str = url.substring(23);
        var decoded = atob(base64str);

        result.Ratio = appState.source.size / decoded.length;
        ratioElement.innerText = result.Ratio.toFixed(1);

        appState.compressed = {
            image,
            result,
            width,
            height,
            url,
            size: decoded.length,
        };
    }

    appState.compressing = false;
    compressButton.innerText = "Compress";
    compressButton.removeAttribute("disabled");
};

async function compressSourceImage() {
    const { compressionOptions, source, compressing } = appState;
    if (compressing) return;

    taskID++;
    appState.compressing = true;
    compressButton.innerText = "Compressing...";
    compressButton.setAttribute("disabled", true);

    try {
        const { image, width, height } = source;
        const imageData = getImagePixels(image, width, height);
        codec.postMessage({
            taskID,
            imageData,
            width,
            height,
            compressionOptions,
        });
    } catch (err) {
        return console.error(err);
    }
}

function getImagePixels(image, width, height) {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    canvas.width = width;
    canvas.height = height;
    ctx.drawImage(image, 0, 0);
    return ctx.getImageData(0, 0, width, height).data;
}

function createImageFromPixels(data, width, height) {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    canvas.width = width;
    canvas.height = height;

    const imageData = ctx.createImageData(width, height);
    imageData.data.set(data);
    ctx.putImageData(imageData, 0, 0);

    return canvas;
}

function getWasmModule() {
    return new Promise((resolve) => {
        // Fetch a new instance every time
        loadWasm("main.wasm").then((instance) => {
            resolve(instance);
        });
    });
}

function loadWasm(path) {
    const memory = new WebAssembly.Memory({
        initial: Math.pow(2, 16),
    });
    const go = new Go();
    go.importObject.env["syscall/js.finalizeRef"] = () => {};
    return new Promise((resolve, reject) => {
        WebAssembly.instantiateStreaming(fetch(path), {
            ...go.importObject,
            js: { mem: memory },
        })
            .then(({ instance }) => {
                go.run(instance);
                resolve(instance);
            })
            .catch(reject);
    });
}
