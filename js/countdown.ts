let countdowns: Element[] = [];

function registerCountdown(el: Element) {
    countdowns.push(el);
}

function updateCountdowns() {
    let until = timeUntil();
    let left = timeLeft();
    let str = until ? ("Starts in: " + formatTime(timeUntil()))
        : left ? ("Ends in: " + formatTime(timeLeft()))
            : "Contest is over";
    if ((!until && !context().contest_started)) {
        setTimeout(() => window.location.reload(), 1000 + Math.random() * 2000);
        clearInterval(iv);
        str = "Contest is starting!";
    }
    if (!left && !context().contest_ended) {
        setTimeout(() => window.location.reload(), 1000 + Math.random() * 2000);
        clearInterval(iv);
        str = "Contest is ending";
    }
    for (let countdown in countdowns) {
        countdowns[countdown].innerHTML = str;
    }
}

let started = false;
let iv: number | undefined;

window.addEventListener("context", function () {
    Array.from(document.getElementsByClassName("contest-countdown")).forEach(e => {
        registerCountdown(e);
    });
    iv = setInterval(updateCountdowns, 100);
});
