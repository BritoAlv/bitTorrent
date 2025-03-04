export class DownloadRequest {
    constructor(id, torrentPath, downloadPath, encryptionLevel) {
        this.id = id;
        this.torrentPath = torrentPath;
        this.downloadPath = downloadPath;
    }
}