const editorState = {
    scale: 1,
    panning: false,
    pan: {
        oldX: 0,
        oldY: 0,
        x: 0,
        y: 0,
    },
    slider: 0.5,
    sliding: false,
    sliderWidth: 10,
    sliderHookRadius: 20,
};

let editorRAF;

function setupEditor() {
    // dummy image
    if (DEBUG == "editor") {
        const dummy1 = createDummyImage("#0f0");
        appState.source = {
            image: dummy1,
            width: dummy1.width,
            height: dummy1.height,
        };
    }

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
    const step = 50;
    for (let i = 0; i < image.width; i += step) {
        ctx.fillStyle = `hsl(${Math.floor(Math.random()*360)}, 100%, 40%)`;
        ctx.fillRect(i, 0, step, image.height);
    }
    return image;
}

function editorFrame() {
    const { source, compressed } = appState;
    const { pan, scale, slider, sliderWidth, sliderHookRadius } = editorState;
    const ctx = editorCtx;

    const halfWidth = editorCanvas.width * 0.5;
    const halfHeight = editorCanvas.height * 0.5;
    const scaledWidth = source.width * scale;

    ctx.clearRect(0, 0, editorCanvas.width, editorCanvas.height);

    let edge = (slider * editorCanvas.width - halfWidth - pan.x) / scaledWidth;
    edge = Math.min(1, Math.max(edge, 0));

    ctx.save();
    ctx.translate(halfWidth + pan.x, halfHeight + pan.y);
    ctx.scale(scale, scale);
    
    // Always render the full width original image
    ctx.drawImage(
        source.image,
        0,
        0,
        source.width,
        source.height,
        0,
        0,
        source.width,
        source.height
    );
    if (compressed) {
        // Render the visible compressed region
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
    }
    ctx.restore();

    // slider width
    ctx.fillStyle = "rgba(40,40,40,.8)";
    const k = mapRange(slider, 0, 1, 0, editorCanvas.width);
    ctx.fillRect(k - sliderWidth / 2, 0, sliderWidth, editorCanvas.height);

    // Slider hook
    ctx.fillStyle = "#000";
    ctx.beginPath();
    ctx.arc(k, halfHeight, sliderHookRadius, 0, 2 * Math.PI);
    ctx.fill();

    // Slider left triangle
    const triW = 10;
    const triH = 20;
    const gap = 4;
    ctx.beginPath();
    ctx.moveTo(k - gap, halfHeight - triH / 2);
    ctx.lineTo(k - gap - triW, halfHeight);
    ctx.lineTo(k - gap, halfHeight + triH / 2);
    ctx.closePath();
    ctx.fillStyle = "hsl(20, 100%, 55%)";
    ctx.fill();

    // Slider right triangle
    ctx.beginPath();
    ctx.moveTo(k + gap, halfHeight - triH / 2);
    ctx.lineTo(k + gap + triW, halfHeight);
    ctx.lineTo(k + gap, halfHeight + triH / 2);
    ctx.closePath();
    ctx.fillStyle = "hsl(200, 100%, 55%)";
    ctx.fill();

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
