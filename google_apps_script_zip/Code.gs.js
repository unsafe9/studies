function doGet(e) {
    // you can collect your spreadsheet data here. below is just for testing
    const csvs = [
        "a,b,c\n1,2,3",
        "d,e,f\n4,5,6"
    ];

    let blobs = [];
    for (let i = 0; i < csvs.length; i++) {
        blobs.push(Utilities.newBlob(csvs[i], "text/csv", `test${i}.csv`));
    }
    
    var zipBlob = Utilities.zip(blobs, "sheets.zip");
    var b64Encoded = Utilities.base64Encode(zipBlob.getBytes());
    return ContentService.createTextOutput(b64Encoded).setMimeType(ContentService.MimeType.TEXT);
}
