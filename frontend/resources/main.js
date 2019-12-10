function formatTimeHHMMSS(timeInSeconds){
    h = Math.floor(timeInSeconds/3600);
    m = Math.floor((timeInSeconds/60)%60);
    s = Math.floor(timeInSeconds%60);

    mm = (m<10?"0":"")+m;
    ss = (s<10?"0":"")+s;
    return ""+h+":"+mm+":"+s;
}

function initTimers(){
    var pageLoadedTime = (new Date()).getTime()/1000;
    var allTimers = document.getElementsByClassName("timer");

    var updateTimers = function(){
        for(var i = 0; i<allTimers.length; i++){
            let currentTime = (new Date()).getTime()/1000;
            let countDir = parseInt(allTimers[i].getAttribute("data-countdir"))
            let time = parseInt(allTimers[i].getAttribute("data-time"))+countDir*(currentTime-pageLoadedTime);
            allTimers[i].innerText = formatTimeHHMMSS(time);
        };
    }

    updateTimers();

    setInterval(updateTimers,1000);
}

initTimers();
