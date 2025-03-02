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
    setInterval(updateStatus, 500);
}

async function updateStatus() {
    torrents.forEach(async (pair) => {
        const id = pair[0];
        const running = pair[1];
        if (!running)
            return;

        const statusProgressBar = document.querySelector(`#status-progress-bar-${id}`);
        const statusPeers = document.querySelector(`#status-peers-${id}`);
        const statusDownload = document.querySelector(`#status-download-${id}`);
        
        const response = await get(API_URL + "update/" + id);

        console.log(`Update response: \n${response}`);
        if (response.Successful) {
            if (response.Progress == 1) {
                statusDownload.innerHTML = "Seed";
            }

            statusProgressBar.value = response.Progress * 100;
            statusPeers.innerHTML = response.Peers;            
        } else {
            errorMessage.innerHTML = response.ErrorMessage;
        }
    });
}

async function stop(torrentPath) {
    const id = torrents.get(torrentPath)[0];
    torrents.set(torrentPath, [id, false]);

    const statusPeers = document.querySelector(`#status-peers-${id}`);
    
    const response = await get(API_URL + `kill/${id}`);
    
    console.log(`Kill response: \n${response}`);
    if (response.Successful) {
        statusPeers.innerHTML = "-1";
    } else {
        errorMessage.innerHTML = response.ErrorMessage;
    }
}

async function download(torrentPath) {
    const pair = torrents.get(torrentPath);
    let id = undefined;
    let running = undefined;

    if (pair != undefined){
        id = pair[0];
        running = pair[1];

        const previousStatus = statusList.querySelector(`#status-${id}`);
        if (previousStatus != null) {
            statusList.removeChild(previousStatus);
        }
    }

    if (running === true) {
        stop(torrentPath);
        torrents.delete(torrentPath);
    }

    id = randomId();
    const downloadRequest = new DownloadRequest(
        id,
        torrentPath,
        "./data/downloads",
        true,
    );

    const response = await post(API_URL+"download", downloadRequest);
    torrents.set(torrentPath, [id, true]);

    console.log(`Download response: \n${response}`);
    // TODO: Get .torrent file's name
    if (response.Successful) {
        statusList.innerHTML += 
`
<li id="status-${id}">
    <label class="status-name" id="status-name-${id}">${torrentPath.split('/').pop()}</label>
    <progress class="status-progress-bar" id="status-progress-bar-${id}" value="0" max="100"></progress>
    <label class="status-peers" id="status-peers-${id}">0</label>
    <label class="status-download" id="status-download-${id}">Download</label>
    <label class="status-stop" id="status-stop-${id}">Stop</label>
    <label class="status-remove" id="status-remove-${id}">X</label>
</li>
`;
        const statusProgressBar = document.querySelector(`#status-progress-bar-${id}`);
        statusProgressBar.display = "inline";

        const statusDownloads = document.querySelectorAll(".status-download");
        const stopDownloads = document.querySelectorAll(".status-stop");
        const removeDownloads = document.querySelectorAll(".status-remove");

        statusDownloads.forEach(item => {
            const id = item.id.split('-').pop();
            item.addEventListener("click", async () => {
                const path = Array.from(torrents.keys()).find(key => torrents.get(key)[0] === id);
                await download(path);
            });
        });

        stopDownloads.forEach(item => {
            const id = item.id.split('-').pop();
            item.addEventListener("click", async () => {
                const path = Array.from(torrents.keys()).find(key => torrents.get(key)[0] === id);
                await stop(path);
            });
        });

        removeDownloads.forEach(item => {
            let id = item.id.split('-').pop();
            item.addEventListener("click", async () => {
                const path = Array.from(torrents.keys()).find(key => torrents.get(key)[0] === id);
                const pair = torrents.get(path);
                id = pair[0];
                const running = pair[1];

                if (running)
                    await stop(path);
                
                const previousStatus = statusList.querySelector(`#status-${id}`);
                if (previousStatus != null)
                    statusList.removeChild(previousStatus);

                torrents.delete(path);
            });
        });
    } else {
        torrents.delete(torrentPath)
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
