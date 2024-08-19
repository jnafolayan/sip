const VERSION = "0.3";

let uploadButton, fileInput, uploadProgress;
let compressButton, downloadButton, backButton;

// VIEWS
let uploadView, editorView;
let editorCanvas, editorCtx;
let psnrElement, ratioElement;

// STATE
let appState = createAppState();
let DEBUG = "";

window.onload = setup;

function setup() {
    console.log(`Sip v${VERSION}.`);

    setupDOM();

    setupEvents();
    subscribeToAppEvents();

    if (DEBUG == "editor") {
        uploadView.classList.add("hide");
        editorView.classList.remove("hide");
        setupEditor();
    }
}

function setupDOM() {
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
}

function setupEvents() {
    fileInput.addEventListener("change", handleImageUpload);

    editorCanvas.addEventListener("wheel", tryZoomEditor);
    editorCanvas.addEventListener("touchstart", tryStartMobileZoom);
    editorCanvas.addEventListener("touchmove", tryMobileZoom);
    editorCanvas.addEventListener("touchend", endMobileZoom);

    editorCanvas.addEventListener("mousedown", tryStartPanning);
    editorCanvas.addEventListener("touchstart", tryStartPanningMobile);
    editorCanvas.addEventListener("mousemove", tryPanEditor);
    editorCanvas.addEventListener("touchmove", tryPanEditorMobile);
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
    useControl("thresholdStrategy", controlCompressionOption("thresholdStrategy"));

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
    EventEditorZoom.subscribe(applyEditorZoom);
}

// state
function createAppState() {
    return {
        source: null,
        compressionOptions: {
            waveletFamily: "haar",
            decompLevel: 1,
            threshold: 10,
            thresholdStrategy: "hard",
        },
        compressed: null,
        compressing: false,
    };
}
