const VERSION = "0.2";

let uploadButton, fileInput, uploadProgress;
let compressButton, downloadButton, backButton;

// VIEWS
let uploadView, editorView;
let editorCanvas, editorCtx;
let psnrElement, ratioElement;

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
let DEBUG = "editor";

window.onload = setup;

function setup() {
    console.log(`Sip v${VERSION}.`);

    uploadButton = document.getElementById("uploadButton");
    uploadProgress = document.getElementById("uploadProgress");
    fileInput = document.getElementById("imageUpload");
    compressButton = document.getElementById("compressButton");
    backButton = document.getElementById("backButton");
    downloadButton = document.getElementById("downloadButton");
    psnrElement = document.getElementById("psnr");
    ratioElement = document.getElementById("compressionRatio");

    uploadView = document.querySelector(".view__upload");
    editorView = document.querySelector(".view__editor");

    editorCanvas = document.querySelector("#editorCanvas");
    editorCtx = editorCanvas.getContext("2d");

    setupEvents();
    subscribeToAppEvents();

    if (DEBUG == "editor") {
        uploadView.classList.add("hide");
        editorView.classList.remove("hide");
        setupEditor();
    }
}

function setupEvents() {
    fileInput.addEventListener("change", handleImageUpload);

    editorCanvas.addEventListener("wheel", handleEditorMouseWheel);
    editorCanvas.addEventListener("touchstart", tryStartMobileZoom);
    editorCanvas.addEventListener("touchmove", tryMobileZoom);
    editorCanvas.addEventListener("touchend", endMobileZoom);

    editorCanvas.addEventListener("mousedown", startImagePanning);
    editorCanvas.addEventListener("touchstart", startImagePanningMobile);
    editorCanvas.addEventListener("mousemove", panImage);
    editorCanvas.addEventListener("touchmove", panImageMobile);
    editorCanvas.addEventListener("mouseup", endImagePanning);
    editorCanvas.addEventListener("mouseout", endImagePanning);
    editorCanvas.addEventListener("touchend", endImagePanning);

    editorCanvas.addEventListener("mousedown", tryStartMovingSlider);
    editorCanvas.addEventListener("touchstart", tryStartMovingSliderMobile);
    editorCanvas.addEventListener("mousemove", tryMoveSlider);
    editorCanvas.addEventListener("touchmove", tryMoveSliderMobile);
    editorCanvas.addEventListener("mouseup", endSliding);
    editorCanvas.addEventListener("mouseout", endSliding);
    editorCanvas.addEventListener("touchend", endSliding);

    compressButton.addEventListener("click", compressSourceImage);
    backButton.addEventListener("click", () => location.reload());
    downloadButton.addEventListener("click", downloadCompressedImage);

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
            compressSourceImage().then(() => {
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
            threshold: 10,
        },
        compressed: null,
        compressing: false,
    };
}
