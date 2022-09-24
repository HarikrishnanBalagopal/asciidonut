// import './wasm_exec.js';

const OUTPUT_PRE = document.querySelector('#output-pre');

// function draw(s) {
// console.log(s);
// OUTPUT_PRE.textContent = s;
// }

async function main() {
    const pre = document.querySelector('#output-pre');
    pre.textContent = 'Hi there!';
    const res = await fetch('assets/wasm/asciidonut.wasm');
    if (!res.ok) return console.error('failed to fetch the wasm module. status:', res.status);
    const moduleBytes = await res.arrayBuffer();
    // const module = await WebAssembly.compile(moduleBytes);
    const go = new Go();
    const module = await WebAssembly.instantiate(moduleBytes, go.importObject);
    console.log('module', module);
    go.run(module.instance);

    let last_t = null;
    const TIME_STEP = 0.1;
    function draw(t) {
        requestAnimationFrame(draw);
        if (!last_t) last_t = t;
        const delta_t = t - last_t;
        if (delta_t < TIME_STEP) return;
        last_t = t;
        const donut = CalculateDonut(0.001*t);
        pre.textContent = donut;
    }
    requestAnimationFrame(draw);
    console.log('done');
}

main().catch(console.error);
