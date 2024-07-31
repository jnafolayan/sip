let slidingOffset = 0;

function tryStartMovingSlider(evt) {
    const { slider, sliderWidth, sliderHookRadius } = editorState;
    const sliderLeft = slider * editorCanvas.width - sliderWidth * 0.5;
    const sliderRight = slider * editorCanvas.width + sliderWidth * 0.5;

    const dx = evt.pageX - slider * editorCanvas.width;
    const dy = evt.pageY - 0.5 * editorCanvas.height;
    const inHook = dx ** 2 + dy ** 2 <= sliderHookRadius ** 2;
    console.log({ dx, dy });
    if ((evt.pageX >= sliderLeft && evt.pageX <= sliderRight) || inHook) {
        editorState.sliding = true;
        editorCanvas.style.cursor = "ew-resize";
        slidingOffset = dx;
    }

    evt.stopPropagation();
}

function tryMoveSlider(evt) {
    const { sliding } = editorState;
    if (!sliding) return;
    editorState.slider = (evt.pageX - slidingOffset) / editorCanvas.width;
}

function endSliding(_evt) {
    const { sliding } = editorState;
    if (!sliding) return;
    editorState.sliding = false;
    editorCanvas.style.cursor = "inherit";
    slidingOffset = 0;
}

function startImagePanning(evt) {
    // Don't pan if user is sliding
    if (editorState.sliding) return;
    editorState.panning = true;
    editorState.pan.oldX = evt.pageX;
    editorState.pan.oldY = evt.pageY;
}

function panImage(evt) {
    const { pan, panning, sliding } = editorState;
    if (!panning) return;
    // Don't pan if user is sliding
    if (sliding) return;

    evt.preventDefault();
    const dx = evt.pageX - pan.oldX;
    const dy = evt.pageY - pan.oldY;
    pan.x += dx;
    pan.y += dy;

    pan.oldX = evt.pageX;
    pan.oldY = evt.pageY;
}

function endImagePanning(_evt) {
    if (!editorState.panning) return;
    editorState.panning = false;
}

function handleEditorMouseWheel(evt) {
    evt.preventDefault();

    const oldScale = editorState.scale;

    editorState.scale += evt.deltaY * -0.005;
    editorState.scale = Math.min(Math.max(0.15, editorState.scale), 3);

    const delta = editorState.scale - oldScale;

    EventEditorZoom.fire({
        pageX: evt.pageX,
        pageY: evt.pageY,
        scale: editorState.scale,
        delta,
    });
}
