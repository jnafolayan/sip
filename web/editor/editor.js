function setupEditor() {
    const {
        source,
        editor: { pan, rendering },
    } = userState;

    // dummy image
    // const image = document.createElement("canvas");
    // const ctx = image.getContext("2d");
    // image.width = 300;
    // image.height = 300;
    // ctx.fillStyle = "#0f0";
    // ctx.fillRect(0, 0, image.width, image.height);
    // userState.source = {
    //     image,
    //     width: image.width,
    //     height: image.height,
    // };

    pan.x = -source.width / 2;
    pan.y = -source.height / 2;
    userState.editor.scale = 1;

    rendering.raf = requestAnimationFrame(editorFrame);
}

function editorFrame() {
    const { source, editor } = userState;
    const {
        pan,
        scale,
        rendering: { ctx },
    } = editor;

    ctx.clearRect(0, 0, editorCanvas.width, editorCanvas.height);

    ctx.save();
    ctx.translate(
        editorCanvas.width / 2 + pan.x,
        editorCanvas.height / 2 + pan.y
    );
    ctx.scale(scale, scale);
    ctx.drawImage(source.image, 0, 0, source.width, source.height);
    ctx.restore();

    editor.rendering.raf = requestAnimationFrame(editorFrame);
}

function applyImageZoom({ pageX, pageY, scale, delta }) {
    const { pan } = userState.editor;

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
