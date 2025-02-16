export class DownloadRequest {
    constructor(id, torrentPath, downloadPath, ipAddress, encryptionLevel) {
        this.id = id;
        this.torrentPath = torrentPath;
        this.downloadPath = downloadPath;
        this.ipAddress = ipAddress
        this.encryptionLevel = encryptionLevel
    }
}