importScripts("assets/wasm_exec.js");

if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

onmessage = (e) => compress(e.data);

async function compress({ taskID, imageData, width, height, compressionOptions }) {
    try {
        await getWasmModule();
        const { Compressed, Result } = Sip_CompressImage(
            imageData,
            width,
            height,
            compressionOptions
        );
        postMessage({
            taskID,
            compressed: Compressed,
            result: Result,
            width,
            height,
        });
    } catch (err) {
        postMessage({ taskID, error: err });
    }
}

function getWasmModule() {
    return new Promise((resolve) => {
        // Fetch a new instance every time
        loadWasm("assets/main.wasm").then((instance) => {
            resolve(instance);
        });
    });
}

function loadWasm(path) {
    const memory = new WebAssembly.Memory({
        initial: Math.pow(2, 8), // 16mb
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
