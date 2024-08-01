if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

async function compress({ imageData, width, height, compressionOptions }) {
    console.log(compressionOptions);

    try {
        const instance = await getWasmModule();
        const { compressed, result } = instance.exports.compressImage(
            imageData,
            width,
            height,
            compressionOptions
        );
        return {
            compressed,
            result,
            width,
            height,
        };
    } catch (err) {
        return console.error(err);
    }
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
