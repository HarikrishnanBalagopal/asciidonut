import './wasm_exec.js';

function fallbackCopyTextToClipboard(text) {
    var textArea = document.createElement("textarea");
    textArea.value = text;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        var successful = document.execCommand('copy');
        var msg = successful ? 'successful' : 'unsuccessful';
        console.log('Fallback: Copying text command was ' + msg);
    } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
    }

    document.body.removeChild(textArea);
}

// https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript
function copyTextToClipboard(text) {
    if (!navigator.clipboard) {
        fallbackCopyTextToClipboard(text);
        return;
    }
    navigator.clipboard.writeText(text).then(function () {
        console.log('Async: Copying to clipboard was successful!');
    }, function (err) {
        console.error('Async: Could not copy text: ', err);
    });
}

async function main() {
    const button_play_pause = document.querySelector('#button-play-pause');
    let paused = false;
    let first = null;
    let offset = 2000.0;
    button_play_pause.addEventListener('click', () => {
        paused = !paused;
        if (paused) {
            console.log('pausing');
            button_play_pause.textContent = '▶️ Play';
        } else {
            console.log('playing');
            requestAnimationFrame(draw);
            button_play_pause.textContent = '⏸ Pause';
        }
    })

    const pre = document.querySelector('#output-pre');
    pre.textContent = 'Hi there!';

    const copy_alert = document.querySelector('#copy-alert');
    const button_copy = document.querySelector('#button-copy');
    button_copy.addEventListener('click', () => {
        copyTextToClipboard(pre.textContent);
        copy_alert.classList.remove('hidden');
        setTimeout(() => copy_alert.classList.add('hidden'), 1000);
    });

    const res = await fetch('assets/wasm/asciidonut.wasm');
    if (!res.ok) return console.error('failed to fetch the wasm module. status:', res.status);
    const moduleBytes = await res.arrayBuffer();
    const go = new Go();
    const module = await WebAssembly.instantiate(moduleBytes, go.importObject);
    // console.log('module', module);
    go.run(module.instance);

    const decoder = new TextDecoder();
    const address = module.instance.exports.GetBufferAddress();
    // console.log('address', address);
    const mem = new Uint8Array(module.instance.exports.memory.buffer, address, 32 * (64 + 1));
    // console.log('mem', mem);

    let last_t = null;
    const TIME_STEP = 1;
    function draw(t) {
        if (paused) { first = null; offset = last_t; return; }
        if (!first) first = t;
        requestAnimationFrame(draw);
        t = t - first + offset;
        if (!last_t) last_t = t;
        const delta_t = t - last_t;
        if (delta_t < TIME_STEP) return;
        last_t = t;
        module.instance.exports.Step(0.01 * t);
        pre.textContent = decoder.decode(mem);
    }
    requestAnimationFrame(draw);
    console.log('done');
}

main().catch(console.error);
