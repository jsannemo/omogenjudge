function pad(t: number) {
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
    return pad(h) + ":" + pad(m) + ":" + pad(s);
}