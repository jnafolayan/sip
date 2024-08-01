if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

async function compressSourceImage() {
    const { compressionOptions, source } = appState;
    console.log(compressionOptions);

    try {
        const instance = await getWasmModule();
        const { image, width, height } = source;
        const imageData = getImagePixels(image, width, height);
        const { Compressed, Result } = Sip_CompressImage(
            imageData,
            width,
            height,
            compressionOptions
        );
        appState.compressed = {
            image: createImageFromPixels(Compressed, width, height),
            width,
            height,
        };
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
