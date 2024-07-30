let uploadButton, fileInput, uploadProgress;

// VIEWS
let uploadView, editorView;
let editorCanvas;

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
let userState;

window.onload = setup;

function setup() {
    userState = createUserState();

    uploadButton = document.getElementById("uploadButton");
    uploadProgress = document.getElementById("uploadProgress");
    fileInput = document.getElementById("imageUpload");

    uploadView = document.querySelector(".view__upload");
    editorView = document.querySelector(".view__editor");

    editorCanvas = document.querySelector("#editorCanvas");
    userState.editor.rendering.ctx = editorCanvas.getContext("2d");


    setupEvents();
    subscribeToAppEvents();

    setupEditor();
}

function setupEvents() {
    fileInput.addEventListener("change", handleImageUpload);

    editorView.addEventListener("wheel", handleEditorMouseWheel);
    editorCanvas.addEventListener("mousedown", startImagePanning);
    editorCanvas.addEventListener("mousemove", panImage);
    editorCanvas.addEventListener("mouseup", endImagePanning);
    editorCanvas.addEventListener("mouseout", (evt) => {
        if (userState.editor.panning) {
            endImagePanning(evt);
        }
    });

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
            EventEditorOpened.fire();
        }, 500);
    });

    // View
    EventEditorOpened.subscribe(setupEditor);
    EventEditorZoom.subscribe(applyImageZoom);
}

function handleImageUpload(evt) {
    const files = evt.target.files;
    if (!files.length) return;

    userState.source = null;

    const [file] = files;
    const fr = new FileReader();
    fr.onload = function () {
        image.src = fr.result;
    };
    fr.onprogress = function (evt) {
        EventFileUploadProgress.fire({ progress: evt.loaded / evt.total });
    };

    const image = new Image();
    image.onload = () => {
        EventFileUploadEnd.fire({ progress: 1 });
        userState.source = {
            image,
            width: image.naturalWidth,
            height: image.naturalHeight,
        };
    };

    EventFileUploadStart.fire({ progress: 0 });
    fr.readAsDataURL(file);
}

// state
function createUserState() {
    return {
        source: null,
        compressionResult: null,
        editor: {
            rendering: {
                ctx: null,
                raf: null,
            },
            scale: 1,
            panning: false,
            pan: {
                oldX: 0,
                oldY: 0,
                x: 0,
                y: 0,
            },
        },
    };
}
