if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const codec = new Worker("codec.worker.js");
let taskID = 0;

codec.onmessage = handleWorkerMessage;

async function handleWorkerMessage(e) {
    if (e.data.taskID != taskID) return;
    // If this happens for whatever reason
    if (!appState.compressing) return;

    if (!e.data.error) {
        const { compressed, result, width, height } = e.data;
        const {
            source: { size: sourceFileSize, name: sourceFileName },
        } = appState;

        const image = createCanvasFromPixels(compressed, width, height);
        const fileName = `${sourceFileName || "sip" + taskID}-compressed.jpg`;
        const compressedFile = await exportCanvasToJPEG(image, fileName, 0.75);

        result.Ratio = sourceFileSize / compressedFile.size * 100;

        ratioElement.innerText = result.Ratio.toFixed(2);
        psnrElement.innerText = result.PSNR.toFixed(2);

        console.log(`Took ${result.Time}s to compress.`);

        if (
            appState.compressed != null &&
            appState.compressed.downloadURL != ""
        ) {
            // Revoke any previous compressed file url
            window.URL.revokeObjectURL(appState.compressed.objectURL);
        }

        appState.compressed = {
            image,
            result,
            width,
            height,
            size: compressedFile.size,
            objectURL: window.URL.createObjectURL(compressedFile),
        };
    }

    appState.compressing = false;
    compressButton.innerText = "Compress";
    compressButton.removeAttribute("disabled");
}

async function compressSourceImage() {
    const { compressionOptions, source, compressing } = appState;
    if (compressing) return;

    taskID++;
    appState.compressing = true;
    compressButton.innerText = "Compressing...";
    compressButton.setAttribute("disabled", true);

    try {
        const { image, width, height } = source;
        const imageData = getCanvasPixels(image);
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
