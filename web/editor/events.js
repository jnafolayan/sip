let slidingOffset = 0;

function initSliding(x, y) {
    const { slider, sliderWidth, sliderHookRadius } = editorState;
    const sliderLeft = slider * editorCanvas.width - sliderWidth * 0.5;
    const sliderRight = slider * editorCanvas.width + sliderWidth * 0.5;

    const dx = x - slider * editorCanvas.width;
    const dy = y - 0.5 * editorCanvas.height;
    const inHook = dx ** 2 + dy ** 2 <= sliderHookRadius ** 2;
    if ((x >= sliderLeft && x <= sliderRight) || inHook) {
        editorState.sliding = true;
        editorCanvas.style.cursor = "ew-resize";
        slidingOffset = dx;
        return true;
    }

    return false;
}

function moveSlider(x) {
    const { sliding } = editorState;
    if (!sliding) return false;
    editorState.slider = (x - slidingOffset) / editorCanvas.width;
    return true;
}

function tryStartMovingSlider(evt) {
    evt.stopPropagation();
    initSliding(evt.pageX, evt.pageY);
}

function tryMoveSlider(evt) {
    moveSlider(evt.pageX);
}

function endSliding(_evt) {
    const { sliding } = editorState;
    if (!sliding) return;
    editorState.sliding = false;
    editorCanvas.style.cursor = "inherit";
    slidingOffset = 0;
}

// MOBILE SLIDING
function tryStartMovingSliderMobile(evt) {
    const [{ pageX, pageY }] = evt.touches;
    evt.stopPropagation();
    initSliding(pageX, pageY);
}

function tryMoveSliderMobile(evt) {
    if (evt.changedTouches.length) {
        moveSlider(evt.changedTouches[0].pageX);
    }
}

function initPanning(x, y) {
    // Don't pan if user is sliding
    if (editorState.sliding) return false;
    editorState.panning = true;
    editorState.pan.oldX = x;
    editorState.pan.oldY = y;
    return true;
}

function panEditor(x, y) {
    const { pan, panning, sliding } = editorState;
    if (!panning) return false;
    // Don't pan if user is sliding
    if (sliding) return false;

    const dx = x - pan.oldX;
    const dy = y - pan.oldY;
    pan.x += dx;
    pan.y += dy;

    pan.oldX = x;
    pan.oldY = y;
    return true;
}

// DESKTOP PANNING
function tryStartPanning(evt) {
    initPanning(evt.pageX, evt.pageY);
}

function tryPanEditor(evt) {
    evt.preventDefault();
    panEditor(evt.pageX, evt.pageY);
}

// MOBILE PANNING
function tryStartPanningMobile(evt) {
    const [{ pageX, pageY }] = evt.touches;
    initPanning(pageX, pageY);
}

function tryPanEditorMobile(evt) {
    if (!evt.changedTouches.length) return;
    evt.preventDefault();

    const [{ pageX, pageY }] = evt.changedTouches;
    panEditor(pageX, pageY);
}

function endImagePanning(_evt) {
    if (!editorState.panning) return;
    editorState.panning = false;
}

// DESKTOP ZOOM
function zoomEditor(delta, pageX, pageY) {
    editorState.scale += delta;
    editorState.scale = Math.min(Math.max(0.15, editorState.scale), 3);

    EventEditorZoom.fire({
        pageX,
        pageY,
        delta,
    });
}

function tryZoomEditor(evt) {
    evt.preventDefault();
    zoomEditor(evt.deltaY * -0.005, evt.pageX, evt.pageY);
}

// MOBILE ZOOM
let initialFingersDistApart = 0;
let mobileZooming = false;
function tryStartMobileZoom(evt) {
    if (!evt.touches || evt.touches.length != 2) return;
    if (mobileZooming) return;
    evt.preventDefault();

    const [a, b] = evt.touches;
    initialFingersDistApart = Math.hypot(a.pageX - b.pageX, a.pageY - b.pageY);
    mobileZooming = true;
}

function tryMobileZoom(evt) {
    if (!evt.changedTouches || evt.changedTouches.length != 2) return;
    if (!mobileZooming) return;
    evt.preventDefault();

    const [a, b] = evt.changedTouches;
    const curFingersDistApart = Math.hypot(
        a.pageX - b.pageX,
        a.pageY - b.pageY
    );
    const newScale = curFingersDistApart / initialFingersDistApart;
    const oldScale = editorState.scale;
    const delta = (newScale - oldScale) * 0.2;

    const midX = a.pageX + (b.pageX - a.pageX) / 2;
    const midY = a.pageY + (b.pageY - a.pageY) / 2;
    zoomEditor(delta, midX, midY);
}

function endMobileZoom(_evt) {
    if (mobileZooming) {
        mobileZooming = false;
        initialFingersDistApart = 0;
    }
}
