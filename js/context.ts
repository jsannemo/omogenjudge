type Context = {
    contest_start_timestamp: number;
    contest_duration: number;
    contest_started: boolean;
    contest_ended: boolean
};

let _context: Context | null = null;

window.addEventListener("load", function () {
    const contextEl = document.getElementById("js_context");
    if (contextEl) {
        _context = JSON.parse(contextEl.innerText);
        console.log("Loaded context", _context);
    } else {
        console.error("Failed loading context!");
    }
    window.dispatchEvent(new Event('context'));
});

function context(): Context {
    return <Context>_context;
}

function time() {
    return new Date().getTime() / 1000;
}

function isStarted() {
    return time() >= context().contest_start_timestamp;
}

function isEnded() {
    let ctx = context();
    return time() >= ctx.contest_start_timestamp + ctx.contest_duration;
}

function timeLeft() {
    let ctx = context();
    return Math.floor(Math.max((ctx.contest_start_timestamp + ctx.contest_duration) - new Date().getTime() / 1000, 0));
}

function timeUntil() {
    let ctx = context();
    return Math.floor(Math.max(ctx.contest_start_timestamp - new Date().getTime() / 1000, 0));
}
