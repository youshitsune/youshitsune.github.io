particlesJS.load('anim', '/static/particles.json', function() {});

function trig(){
    htmx.trigger("#brand", "run");
}

function redirect(url){
    window.location.assign(url);
}

htmx.on("#brand", "click", trig);
htmx.on("#hero", "htmx:afterRequest", prepareSwap);
htmx.on("#hero","htmx:afterSwap", swap);


function unfade(element) {
    var op = 0.1;
    element.style.display = 'block';
    var timer = setInterval(function () {
        if (op >= 1){
            clearInterval(timer);
        }
        element.style.opacity = op;
        element.style.filter = 'alpha(opacity=' + op * 100 + ")";
        op += op * 0.1;
    }, 15);
}

function fade(element){
    var op = 1;
    element.style.display = 'block';
    var timer = setInterval(function() {
        if (op <= 0) {
            clearInterval(timer);
        }
        element.style.opacity = op;
        element.style.filter = 'alpha(opacity=' + op * 100 + ")";
        op -= op * 0.1;
    }, 15);
}

function prepareSwap(){
    var elem = document.getElementById("hero");
    fade(elem);
}

function swap() {
    var elem = document.getElementById("hero");
    unfade(elem);
}

window.onload = async function(){
    var elem = document.body;
    elem.style.opacity = 0;
    unfade(elem);
}
