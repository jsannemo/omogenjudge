function pad00(t: number) {
    if (t < 10) {
        return "0" + t;
    }
    return "" + t;
}

function formatTime(seconds: number): string {
    if (seconds < 0) {
        return "-" + formatTime(seconds);
    }
    let s = seconds % 60;
    seconds = (seconds - s) / 60;
    let m = seconds % 60;
    seconds = (seconds - m) / 60;
    let h = seconds;
    return pad00(h) + ":" + pad00(m) + ":" + pad00(s);
}

(function () {
    for (let element of document.getElementsByClassName("local_date")) {
        let htmlElement = element as HTMLElement;
        let timestamp = parseInt(htmlElement.dataset["timestamp"]!) * 1000;
        htmlElement.innerText = new Date(timestamp).toLocaleString([...navigator.languages], {timeZoneName: "short"});
    }

    for (let element of document.getElementsByClassName("simple_local_date")) {
        let htmlElement = element as HTMLElement;
        let timestamp = parseInt(htmlElement.dataset["timestamp"]!) * 1000;
        htmlElement.innerText = new Date(timestamp).toLocaleString([...navigator.languages]);
    }
}());
