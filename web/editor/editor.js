const editorState = {
    scale: 1,
    panning: false,
    pan: {
        oldX: 0,
        oldY: 0,
        x: 0,
        y: 0,
    },
};

let editorRAF;

function setupEditor() {
    const { source } = appState;

    // dummy image
    // const dummy = createDummyImage();
    // appState.source = {
    //     image,
    //     width: image.width,
    //     height: image.height,
    // };

    editorState.pan.x = -source.width / 2;
    editorState.pan.y = -source.height / 2;
    editorState.scale = 1;

    editorRAF = requestAnimationFrame(editorFrame);
}

function createDummyImage() {
    const image = document.createElement("canvas");
    const ctx = image.getContext("2d");
    image.width = 300;
    image.height = 300;
    ctx.fillStyle = "#0f0";
    ctx.fillRect(0, 0, image.width, image.height);
    return image;
}

function editorFrame() {
    const { source } = appState;
    const {
        pan,
        scale,
    } = editorState;
    const ctx = editorCtx;

    ctx.clearRect(0, 0, editorCanvas.width, editorCanvas.height);

    ctx.save();
    ctx.translate(
        editorCanvas.width / 2 + pan.x,
        editorCanvas.height / 2 + pan.y
    );
    ctx.scale(scale, scale);
    ctx.drawImage(source.image, 0, 0, source.width, source.height);
    ctx.restore();

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
