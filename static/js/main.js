function trig(){
    htmx.trigger("#brand", "run");
}

function redirect(url){
    window.location.assign(url);
}

htmx.on("#brand", "click", trig);

particlesJS.load('anim', '/static/particles.json', function() {});
