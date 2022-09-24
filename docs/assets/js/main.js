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

    // console.log('module.instance.exports.CalculateDonut(0.001 * t)', module.instance.exports.CalculateDonut(0.001));
    const decoder = new TextDecoder();
    const address = module.instance.exports.GetBufferAddress();
    console.log('address', address);
    const mem = new Uint8Array(module.instance.exports.memory.buffer, address, 32 * (64 + 1));
    console.log('mem', mem);

    let last_t = null;
    const TIME_STEP = 1;
    function draw(t) {
        requestAnimationFrame(draw);
        if (!last_t) last_t = t;
        const delta_t = t - last_t;
        if (delta_t < TIME_STEP) return;
        last_t = t;
        module.instance.exports.Step(0.001 * t);
        pre.textContent = decoder.decode(mem);
    }
    requestAnimationFrame(draw);
    console.log('done');
}

main().catch(console.error);
