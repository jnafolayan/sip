let uploadButton, fileInput, uploadProgress;
let compressButton;

// VIEWS
let uploadView, editorView;
let editorCanvas, editorCtx;

// EVENTS
const [EventFileUploadStart, EventFileUploadProgress, EventFileUploadEnd] = [
    new AppEvent("FILE_UPLOAD_START"),
    new AppEvent("FILE_UPLOAD_PROGRESS"),
    new AppEvent("FILE_UPLOAD_END"),
];
const EventEditorOpened = new AppEvent("EDITOR_OPENED");
const EventEditorZoom = new AppEvent("EDITOR_ZOOM");
const EventEditorMouseDown = new AppEvent("EDITOR_MOUSE_DOWN");

// STATE
let appState = createAppState();

window.onload = setup;

function setup() {
    uploadButton = document.getElementById("uploadButton");
    uploadProgress = document.getElementById("uploadProgress");
    fileInput = document.getElementById("imageUpload");
    compressButton = document.getElementById("compressButton");

    uploadView = document.querySelector(".view__upload");
    editorView = document.querySelector(".view__editor");

    editorCanvas = document.querySelector("#editorCanvas");
    editorCtx = editorCanvas.getContext("2d");

    setupEvents();
    subscribeToAppEvents();

    // setupEditor();
}

function setupEvents() {
    fileInput.addEventListener("change", handleImageUpload);

    editorView.addEventListener("wheel", handleEditorMouseWheel);
    editorCanvas.addEventListener("mousedown", startImagePanning);
    editorCanvas.addEventListener("mousemove", panImage);
    editorCanvas.addEventListener("mouseup", endImagePanning);
    editorCanvas.addEventListener("mouseout", endImagePanning);

    editorCanvas.addEventListener("mousedown", tryStartMovingSlider);
    editorCanvas.addEventListener("mousemove", tryMoveSlider);
    editorCanvas.addEventListener("mouseup", endSliding);
    editorCanvas.addEventListener("mouseout", endSliding);

    compressButton.addEventListener("click", compressSourceImage);

    useControl("waveletFamily", controlCompressionOption("waveletFamily"));
    useControl("decompLevel", controlCompressionOption("decompLevel"));
    useControl("threshold", controlCompressionOption("threshold"));

    window.addEventListener("resize", handleWindowResize);
    handleWindowResize();
}

function handleWindowResize() {
    editorCanvas.width = window.innerWidth;
    editorCanvas.height = window.innerHeight;
}

function subscribeToAppEvents() {
    // File upload
    EventFileUploadStart.subscribe(stepFileUploadAnimation);
    EventFileUploadProgress.subscribe(stepFileUploadAnimation);
    EventFileUploadEnd.subscribe(stepFileUploadAnimation);
    EventFileUploadEnd.subscribe(() => {
        setTimeout(() => {
            uploadView.classList.add("hide");
            editorView.classList.remove("hide");
            compressSourceImage()
                .then(() => {
                    EventEditorOpened.fire();
                });
        }, 500);
    });

    // View
    EventEditorOpened.subscribe(setupEditor);
    EventEditorZoom.subscribe(applyImageZoom);
}

// state
function createAppState() {
    return {
        source: null,
        compressionOptions: {
            waveletFamily: "haar",
            decompLevel: 1,
            threshold: 50,
        },
        compressed: null,
    };
}
