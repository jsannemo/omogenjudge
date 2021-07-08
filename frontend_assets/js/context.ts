/**
 * The Context contains values exported from Django on each request, written as JSON to the HTML.
 *
 * The schema and values are documented on the Python side in js_context.py.
 */
export type Context = {
    contest_start_timestamp?: number;
    contest_duration: number;
    contest_started: boolean;
    contest_ended: boolean
    only_virtual: boolean;
};

let context: Context | null = null;

window.addEventListener("load", function () {
    const contextEl = document.getElementById("js_context");
    if (contextEl) {
        context = JSON.parse(contextEl.innerText);
        console.log("Loaded context", context);
    } else {
        console.error("Failed loading context!");
    }
    window.dispatchEvent(new Event("context"));
});

/**
 * Returns the page context.
 *
 * This should only be called after the context has been loaded. For scripts that need to execute actions on load but
 * required the context to exist, wait for the "context" event to be dispatched.
 */
export function getContext(): Context {
    return context!;
}