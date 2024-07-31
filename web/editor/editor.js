const editorState = {
    scale: 1,
    panning: false,
    pan: {
        oldX: 0,
        oldY: 0,
        x: 0,
        y: 0,
    },
    slider: 0.3,
};

let editorRAF;

function setupEditor() {
    // dummy image
    const dummy1 = createDummyImage("#0f0");
    const dummy2 = createDummyImage("#0ff");
    appState.source = {
        image: dummy1,
        width: dummy1.width,
        height: dummy1.height,
    };
    appState.compressed = {
        image: dummy2,
        width: dummy2.width,
        height: dummy2.height,
    };

    editorState.pan.x = -appState.source.width / 2;
    editorState.pan.y = -appState.source.height / 2;
    editorState.scale = 1;

    editorRAF = requestAnimationFrame(editorFrame);
}

function createDummyImage(color) {
    const image = document.createElement("canvas");
    const ctx = image.getContext("2d");
    image.width = 300;
    image.height = 300;
    ctx.fillStyle = color;
    ctx.fillRect(0, 0, image.width, image.height);
    return image;
}

function editorFrame() {
    const { source, compressed } = appState;
    const { pan, scale, slider } = editorState;
    const ctx = editorCtx;

    const halfWidth = editorCanvas.width * 0.5;
    const halfHeight = editorCanvas.height * 0.5;
    const scaledWidth = source.width * scale;

    ctx.clearRect(0, 0, editorCanvas.width, editorCanvas.height);

    let edge =
        (slider * editorCanvas.width - halfWidth - pan.x) /
        scaledWidth;
    edge = Math.min(1, Math.max(edge, 0));

    ctx.save();
    ctx.translate(
        halfWidth + pan.x,
        halfHeight + pan.y
    );
    ctx.scale(scale, scale);
    ctx.drawImage(
        source.image,
        0,
        0,
        source.width * edge,
        source.height,
        0,
        0,
        source.width * edge,
        source.height
    );
    ctx.drawImage(
        compressed.image,
        edge * compressed.width,
        0,
        compressed.width * (1 - edge),
        compressed.height,
        source.width * edge,
        0,
        compressed.width * (1 - edge),
        source.height
    );
    ctx.restore();

    // slider width
    const w = 4;
    ctx.fillStyle = "rgba(40,40,40,.8)";
    const k = mapRange(slider, 0, 1, 0, editorCanvas.width);
    ctx.fillRect(k - w / 2, 0, w, editorCanvas.height);

    editorRAF = requestAnimationFrame(editorFrame);
}

function applyImageZoom({ pageX, pageY, scale, delta }) {
    const { pan } = editorState;

    const centerX = editorCanvas.width / 2;
    const centerY = editorCanvas.height / 2;

    const mouseOffsetX = pageX - centerX;
    const mouseOffsetY = pageY - centerY;

    const pivotX = mouseOffsetX - pan.x;
    const pivotY = mouseOffsetY - pan.y;

    const offsetX = -pivotX * delta;
    const offsetY = -pivotY * delta;

    pan.x += offsetX;
    pan.y += offsetY;
}
