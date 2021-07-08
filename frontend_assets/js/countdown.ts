import {getContext} from "./context";

function contestTimeLeft() {
    let ctx = getContext();
    return Math.floor(Math.max((ctx.contest_start_timestamp!! + ctx.contest_duration) - new Date().getTime() / 1000, 0));
}

function timeUntil() {
    let ctx = getContext();
    return Math.floor(Math.max(ctx.contest_start_timestamp!! - new Date().getTime() / 1000, 0));
}

let countdowns: Element[] = [];

function registerCountdown(el: Element) {
    countdowns.push(el);
}

function updateCountdowns() {
    let str = "";
    if (getContext().contest_start_timestamp === null) {
        if (getContext().only_virtual) {
            str = "Practice contest";
        } else {
            str = "Wait for start";
        }
    } else {
        let until = timeUntil();
        let left = contestTimeLeft();
        str = until ? ("Starts in: " + formatTime(timeUntil()))
            : left ? ("Ends in: " + formatTime(contestTimeLeft()))
                : "Contest is over";
        // Refresh the page on contest start with somee jitter to spread out requests a bit.
        if ((!until && !getContext().contest_started)) {
            setTimeout(() => window.location.reload(), 1000 + Math.random() * 3000);
            clearInterval(iv);
            str = "Contest is starting - loading the problems...";
        }
    }
    for (let countdown in countdowns) {
        countdowns[countdown].innerHTML = str;
    }
}

let iv: number | undefined;

window.addEventListener("context", function () {
    Array.from(document.getElementsByClassName("contest-countdown")).forEach(e => {
        registerCountdown(e);
    });
    iv = setInterval(updateCountdowns, 100);
});
