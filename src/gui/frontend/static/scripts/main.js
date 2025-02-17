import { get, post } from "./httpMethods.js";
import { DownloadRequest } from "./contracts.js";
import { randomId } from "./utils.js";
import { API_URL } from "./constants.js";

const errorMessage = document.querySelector("#error-message");
const torrentPathInput = document.querySelector("#add-torrent-input");
const addTorrentButton = document.querySelector("#add-torrent-button");
const statusList = document.querySelector("#status-list");

const torrents = new Map();

function updateStatusPeriodically() {
    setInterval(updateStatus, 10000);
}

async function updateStatus() {
    torrents.forEach(async (id) => {
        if (id == undefined)
            return;

        const statusProgressBar = document.querySelector(`#status-progress-bar-${id}`);
        const statusPeers = document.querySelector(`#status-peers-${id}`);
        const statusDownload = document.querySelector(`#status-download-${id}`);
        
        const response = await get(API_URL + "update?id=" + id);

        if (response.Successful) {
            if (response.Progress == 1) {
                statusDownload.innerHTML = "Seed";
            }

            statusProgressBar.value = response.Progress * 100;
            statusPeers.innerHTML = response.Peers;            
        } else {
            errorMessage.innerHTML += " " + response.ErrorMessage;
        }
        console.log(response)
    });
}

async function stop(torrentPath) {
    const id = torrents.get(torrentPath);
    torrents.set(torrentPath, undefined);
    const statusPeers = document.querySelector(`#status-peers-${id}`);
    
    const response = await get(API_URL + `kill?id=${id}`);
    
    if (response.Successful) {
        statusPeers.innerHTML = "-1";
    } else {
        errorMessage.innerHTML = response.ErrorMessage;
    }
}

async function download(torrentPath) {
    let id = torrents.get(torrentPath);
    const previousStatus = statusList.querySelector(`#status-${id}`);

    if (previousStatus != null) {
        statusList.removeChild(previousStatus);
    }

    id = randomId();
    torrents.set(torrentPath, id);
    const downloadRequest = new DownloadRequest(
        id,
        torrentPath,
        "./data/downloads",
        "127.0.0.1",
        false,
    );

    const response = await post(API_URL+"download", downloadRequest);
    
    // TODO: Get .torrent file's name
    if (response.Successful) {
        statusList.innerHTML += 
`
<li id="status-${id}">
    <label class="status-name" id="status-name-${id}">${torrentPath.substring(torrentPath.length - 30, torrentPath.length - 10)}</label>
    <progress class="status-progress-bar" id="status-progress-bar-${id}" value="0" max="100"></progress>
    <label class="status-peers" id="status-peers-${id}">0</label>
    <label class="status-download" id="status-download-${id}">Download</label>
    <label class="status-stop" id="status-stop-${id}">Stop</label>
    <label class="status-remove" id="status-remove-${id}">X</label>
</li>
`;
        const statusProgressBar = document.querySelector(`#status-progress-bar-${id}`);
        statusProgressBar.display = "inline";

        const statusDownload = document.querySelector(`#status-download-${id}`);
        const statusStop = document.querySelector(`#status-stop-${id}`);
        const statusRemove = document.querySelector(`#status-remove-${id}`);

        statusDownload.addEventListener("click", async () => {
            const path = torrentPath;
            await download(path);
        });

        statusStop.addEventListener("click", async () => {
            const path = torrentPath;
            await stop(path);
        });

        statusRemove.addEventListener("click", async () => {
            const path = torrentPath;
            await stop(path);
            const id = torrents.get(path);
            const previousStatus = statusList.querySelector(`#status-${id}`);
            
            if (previousStatus != null)
                statusList.removeChild(previousStatus);
        });
    } else {
        errorMessage.innerHTML = response.ErrorMessage;
    }
}

torrentPathInput.addEventListener("click", () => {
    errorMessage.innerHTML = ""
});

addTorrentButton.addEventListener("click", async () => {
    const torrentPath = torrentPathInput.value;
    torrentPathInput.value = ""

    await download(torrentPath)
});

updateStatusPeriodically()

// let s = 
//         // statusList.removeChild(s)