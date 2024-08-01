if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

let wasmLoaded = false;
let moduleInstance = null;

async function compressSourceImage() {
    const { compressionOptions, source } = appState;
    console.log(compressionOptions);

    try {
        const instance = await getWasmModule();
        const { image, width, height } = source;
        const imageData = getImageData(image, width, height);
        const { Compressed, Result } = Sip_CompressImage(
            Array.from(imageData),
            width,
            height,
            compressionOptions
        );
        appState.compressed = {
            image: imageFromImageData(Compressed, width, height),
            width,
            height,
        };
    } catch (err) {
        return console.error(err);
    }
}

function getImageData(image, width, height) {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    canvas.width = width;
    canvas.height = height;
    ctx.drawImage(image, 0, 0);
    return ctx.getImageData(0, 0, width, height).data;
}

function imageFromImageData(data, width, height) {
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
        if (!moduleInstance) {
            loadWasm("main.wasm").then((instance) => {
                moduleInstance = instance;
                resolve(instance);
            });
        } else {
            resolve(moduleInstance);
        }
    });
}

function loadWasm(path) {
    const go = new Go();
    go.importObject.env["syscall/js.finalizeRef"] = () => {}
    return new Promise((resolve, reject) => {
        WebAssembly.instantiateStreaming(fetch(path), go.importObject)
            .then(({ instance }) => {
                go.run(instance);
                resolve(instance);
            })
            .catch(reject);
    });
}
