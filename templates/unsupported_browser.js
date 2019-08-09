function onHover(element, prefix) {
    //element.classList.add('container-hover');
    element.className += " container-hover";
    var unHoverElement = element.querySelector('.' + prefix + ':not(.hidden)');
    var hoverElement = element.querySelector('.' + prefix + '-hover.hidden');
    if (unHoverElement && hoverElement) {
        unHoverElement.className += " hidden";
        //unHoverElement.classList.add('hidden');
        //hoverElement.classList.remove('hidden');
        hoverElement.className = hoverElement.className.replace(/\bhidden\b/g, "");
    }
}
function onUnHover(element, prefix) {
    element.className = element.className.replace(/\bcontainer-hover\b/g, "");
    //element.classList.remove('container-hover');
    var unHoverElement = element.querySelector('.' + prefix + '.hidden');
    var hoverElement = element.querySelector('.' + prefix + '-hover:not(.hidden)');
    if (unHoverElement && hoverElement) {
        hoverElement.className += " hidden";
        //hoverElement.classList.add('hidden');
        //unHoverElement.classList.remove('hidden');
        unHoverElement.className = unHoverElement.className.replace(/\bhidden\b/g, "");
    }
}

document.addEventListener('DOMContentLoaded', function () {
    var hovers = document.querySelectorAll("div[data-mattermost-hover]");
    for (var i = 0; i < hovers.length; i++) {
        var element = hovers[i];
        element.addEventListener("mouseover", function(e) {
            onHover(e.currentTarget, e.currentTarget.getAttribute("data-mattermost-hover"));
        });
        element.addEventListener("mouseout", function(e) {
            onUnHover(e.currentTarget, e.currentTarget.getAttribute("data-mattermost-hover"))
        });
    }
    var clicks = document.querySelectorAll("div[data-mattermost-click], button[data-mattermost-click]");
    for (var i = 0; i < clicks.length; i++) {
        var element = clicks[i];
        element.addEventListener("click", function(e) {
            window.location.href = e.currentTarget.getAttribute("data-mattermost-click");
        });
    };
});