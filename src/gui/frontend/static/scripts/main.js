import { get, post } from "./httpMethods.js";
import { DownloadRequest } from "./contracts.js";
import { randomId } from "./utils.js";
import { API_URL } from "./constants.js";

const errorMessage = document.querySelector("#error-message");
const torrentPathInput = document.querySelector("#add-torrent-input");
const addTorrentButton = document.querySelector("#add-torrent-button");
const statusList = document.querySelector("#status-list");

torrentPathInput.addEventListener("click", () => {
    errorMessage.innerHTML = ""
});

addTorrentButton.addEventListener("click", async () => {
    const torrentPath = torrentPathInput.value;
    const downloadRequest = new DownloadRequest(
        randomId(),
        torrentPath,
        "./data/downloads",
        "127.0.0.1",
        false,
    );

    const response = await post(API_URL+"download", downloadRequest);
    
    if (response.Successful) {
        console.log("Successful");
    } else {
        errorMessage.innerHTML = response.ErrorMessage;
    } 

    console.log("Message");
});

// console.log(statusList.innerHTML)
// get("asdsad")